package message

import (
	"log"
	"sync"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/core/policy"
	"github.com/yeren0143/DDS/fastrtps/rtps/endpoint"
	//"github.com/yeren0143/DDS/fastrtps/rtps/reader"
	//"github.com/yeren0143/DDS/fastrtps/rtps/writer"
)

const (
	KInfoSrcSubmsgLength = 20
)

type IReceiverOwner interface {
	GetGuid() *common.GUIDT
	AssertRemoteParticipantLiveliness(*common.GUIDPrefixT)
}

type processDataMessageFunction func(entityID *common.EntityIDT, change *common.CacheChangeT)
type processDataFragmentMessageFunction func(*common.EntityIDT, *common.CacheChangeT, uint32, uint32, uint16)

// Receiver process the received messages.
type Receiver struct {
	participant      IReceiverOwner
	sourceVersion    common.ProtocolVersionT
	sourceVendorID   common.VendorIDT
	sourceGUIDPrefix common.GUIDPrefixT
	destGUIDPrefix   common.GUIDPrefixT
	haveTimeStamp    bool
	timeStamp        common.Time
	mutex            sync.Mutex
	//associatedWriters []writer.IRTPSWriter
	//associatedReaders map[common.EntityIDT][]reader.IRTPSReader
	associatedWriters []IRtpsMsgWriter
	associatedReaders map[common.EntityIDT][]IRtpsMsgReader
	dataMsgFunc       processDataMessageFunction
	dataFragMsgFunc   processDataFragmentMessageFunction
}

func NewMessageReceiver(holder IReceiverOwner, rcvBufferSize uint32) *Receiver {
	receiver := Receiver{
		participant:       holder,
		sourceVersion:     common.KProtocolVersion,
		sourceVendorID:    common.KVendorIDTUnknown,
		sourceGUIDPrefix:  common.KGuidPrefixUnknown,
		destGUIDPrefix:    common.KGuidPrefixUnknown,
		haveTimeStamp:     false,
		timeStamp:         common.KTimeInvalid,
		associatedReaders: make(map[common.EntityIDT][]IRtpsMsgReader),
	}
	log.Printf("Created with CDRMessage of size: %v", rcvBufferSize)

	return &receiver
}

// Reset the MessageReceiver to process a new message.
func (receiver *Receiver) reset() {
	receiver.sourceVersion = common.KProtocolVersion
	receiver.sourceVendorID = common.KVendorIDTUnknown
	receiver.sourceGUIDPrefix = common.KGuidPrefixUnknown
	receiver.haveTimeStamp = false
	receiver.timeStamp = common.KTimeInvalid
}

func (receiver *Receiver) willAReaderAcceptMsgDirectedTo(readerID *common.EntityIDT) (IRtpsMsgReader, bool) {
	var firstReader IRtpsMsgReader
	if len(receiver.associatedReaders) == 0 {
		log.Printf("Data received when NO readers are listening")
		return nil, false
	}

	if readerID != common.KEIDUnknown {
		if readers, ok := receiver.associatedReaders[*readerID]; ok {
			firstReader = readers[0]
			return firstReader, true
		}
	} else {
		for _, readers := range receiver.associatedReaders {
			if len(readers) == 0 {
				continue
			}

			for _, reader := range readers {
				if reader.AcceptMessagesToUnknownReaders() {
					firstReader = reader
					return firstReader, true
				}
			}
		}
	}
	log.Printf("No Reader accepts this message (directed to: %v )", readerID)

	return nil, false
}

func (receiver *Receiver) AssociateEndpoint(toAdd endpoint.IEndpoint) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()
	if toAdd.GetAttributes().EndpointKind == common.KWriter {
		awriter, ok := toAdd.(IRtpsMsgWriter)
		if !ok {
			log.Fatalln("associateEndpoint with fault type")
		}
		for _, it := range receiver.associatedWriters {
			if awriter == it {
				return
			}
		}
		receiver.associatedWriters = append(receiver.associatedWriters, awriter)
	} else {
		areader, ok := toAdd.(IRtpsMsgReader)
		if !ok {
			log.Fatalln("associateEndpoint with fault type")
		}
		entID := toAdd.GetGUID().EntityID
		if readers, ok := receiver.associatedReaders[entID]; ok {
			for _, item := range readers {
				if item == areader {
					return
				}
			}
			readers = append(readers, areader)
		} else {
			var readers []IRtpsMsgReader
			readers = append(readers, areader)
			receiver.associatedReaders[entID] = readers
		}
	}
}

func (receiver *Receiver) readSubmessageHeader(msg *common.CDRMessage, smh *SubmessageHeaderT) bool {
	if (msg.Length - msg.Pos) < 4 {
		log.Print("SubmessageHeader too short")
		return false
	}

	smh.SubMessageID = msg.Buffer[msg.Pos]
	msg.Pos++
	smh.Flags = msg.Buffer[msg.Pos]
	msg.Pos++

	// set endianness of message
	if (smh.Flags & common.Bit(0)) != 0 {
		msg.MsgEndian = common.LITTLEEND
	} else {
		msg.MsgEndian = common.BIGEND
	}

	var length uint16
	readUInt16(msg, &length)
	if (msg.Pos + uint32(length)) > msg.Length {
		log.Printf("SubMsg of invalid length (%v) with current msg position/length (%v/)", msg.Pos, msg.Length)
		return false
	}

	if length == 0 && smh.SubMessageID != KInfoTs && smh.SubMessageID != KPad {
		// THIS IS THE LAST SUBMESSAGE
		smh.SubmessageLength = msg.Length - msg.Pos
		smh.IsLast = true
	} else {
		smh.SubmessageLength = uint32(length)
		smh.IsLast = false
	}

	return true
}

func (receiver *Receiver) checkRTPSHeader(msg *common.CDRMessage) bool {
	header := []rune("RTPS")
	if string(msg.Buffer[:4]) != string(header) {
		log.Println("Msg received with no RTPS in header, ignoring...")
		return false
	}

	msg.Pos += 4

	// CHECK AND SET protocol version
	if msg.Buffer[msg.Pos] <= common.KProtocolVersion.Major {
		receiver.sourceVersion.Major = msg.Buffer[msg.Pos]
		msg.Pos++

		receiver.sourceVersion.Minor = msg.Buffer[msg.Pos]
		msg.Pos++
	} else {
		log.Fatal("Major RTPS Version not supported")
		return false
	}

	// set source vendor id
	receiver.sourceVendorID.Vendor[0] = msg.Buffer[msg.Pos]
	msg.Pos++
	receiver.sourceVendorID.Vendor[1] = msg.Buffer[msg.Pos]
	msg.Pos++
	// set source guid prefix
	data, ok := readData(msg, common.KGUIDPrefixSize)
	if !ok {
		log.Fatal("receive bad data, prefix not ok")
	}
	for i := 0; i < len(receiver.sourceGUIDPrefix.Value); i++ {
		receiver.sourceGUIDPrefix.Value[i] = data[i]
	}
	receiver.haveTimeStamp = false

	return true
}

type findAllReadersCallback func(IRtpsMsgReader)

func (receiver *Receiver) findAllReaders(readerID *common.EntityIDT, callback findAllReadersCallback) {
	if readerID != common.KEIDUnknown {
		if readers, ok := receiver.associatedReaders[*readerID]; ok {
			for _, reader := range readers {
				callback(reader)
			}
		}
	} else {
		for _, readers := range receiver.associatedReaders {
			for _, reader := range readers {
				if reader.AcceptMessagesToUnknownReaders() {
					callback(reader)
				}
			}
		}
	}
}

func (receiver *Receiver) procSubmsgHeartbeat(msg *common.CDRMessage, smh *SubmessageHeaderT) bool {
	endiannessFlag := (smh.Flags & common.Bit(0)) != 0
	finalFlag := (smh.Flags & common.Bit(1)) != 0
	livelinessFlag := (smh.Flags & common.Bit(2)) != 0
	// Assign message endianness
	if endiannessFlag {
		msg.MsgEndian = common.LITTLEEND
	} else {
		msg.MsgEndian = common.BIGEND
	}
	var readerGUID, writerGUID common.GUIDT
	readerGUID.Prefix = receiver.destGUIDPrefix
	readEntityID(msg, &readerGUID.EntityID)
	writerGUID.Prefix = receiver.sourceGUIDPrefix
	readEntityID(msg, &writerGUID.EntityID)
	var firstSN, lastSN common.SequenceNumberT
	readSequenceNumber(msg, &firstSN)
	readSequenceNumber(msg, &lastSN)
	if lastSN.Less(&firstSN) && (lastSN.Value != (firstSN.Value - 1)) {
		log.Printf("Invalid Heartbeat received (%v) - (%v),  ignoring", firstSN, lastSN)
		return false
	}

	var HBCount uint32
	readUInt32(msg, &HBCount)

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()
	// Look for the correct reader and writers:
	callback := func(reader IRtpsMsgReader) {
		reader.ProcessHeartbeatMsg(&writerGUID, HBCount, &firstSN, &lastSN, finalFlag, livelinessFlag)
	}
	receiver.findAllReaders(&readerGUID.EntityID, callback)
	return true
}

func (receiver *Receiver) procSubmsgDataFrag(msg *common.CDRMessage, smh *SubmessageHeaderT) bool {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	// READ and PROCESS
	if smh.SubmessageLength < KRTPSMessageDataMinLength {
		log.Printf("Too short submessage received, ignoring")
		return false
	}

	// Fill flags bool values
	endiannessFlag := (smh.Flags & common.Bit(0)) != 0
	inlineQosFlag := (smh.Flags & common.Bit(1)) != 0
	keyFlag := (smh.Flags & common.Bit(2)) != 0

	// Assign message endianness
	if endiannessFlag {
		msg.MsgEndian = common.LITTLEEND
	} else {
		msg.MsgEndian = common.BIGEND
	}

	// Extra flags don't matter now. Avoid those bytes
	msg.Pos += 2

	var octetsToInLineQos int16
	valid := readInt16(msg, &octetsToInLineQos)

	// reader and writer ID
	var readerID common.EntityIDT
	readEntityID(msg, &readerID)

	// WE KNOW THE READER THAT THE MESSAGE IS DIRECTED TO SO WE LOOK FOR IT:
	_, ok := receiver.willAReaderAcceptMsgDirectedTo(&readerID)
	if !ok {
		return false
	}

	// FOUND THE READER.
	// We ask the reader for a cachechange to store the information.
	ch := common.NewCacheChangeT()
	ch.WriterGUID.Prefix = receiver.sourceGUIDPrefix
	valid = valid && readEntityID(msg, &ch.WriterGUID.EntityID)

	// Get sequence number
	valid = valid && readSequenceNumber(msg, &ch.SequenceNumber)

	if ch.SequenceNumber.Value >= ^uint64(0) {
		log.Println("Invalid message received, bad sequence Number")
		return false
	}

	// READ FRAGMENT NUMBER
	var fragmentStartingNum uint32
	valid = valid && readUInt32(msg, &fragmentStartingNum)

	// READ FRAGMENTSINSUBMESSAGE
	var fragmentsInSubmessage uint16
	valid = valid && readUInt16(msg, &fragmentsInSubmessage)

	// READ FRAGMENTSIZE
	var fragmentSize uint16
	valid = valid && readUInt16(msg, &fragmentSize)

	// READ SAMPLESIZE
	var sampleSize uint32
	valid = valid && readUInt32(msg, &sampleSize)

	if !valid {
		return false
	}

	// Jump ahead if more parameters are before inlineQos (not in this version, maybe if further minor versions.)
	if octetsToInLineQos > KRtpsMessageOctetStoinlineqosDataFragSubmsg {
		msg.Pos += uint32(octetsToInLineQos - KRtpsMessageOctetStoinlineqosDataFragSubmsg)
		if msg.Pos > msg.Length {
			log.Printf("Invalid jump through msg, msg->pos %v > msg->length %v", msg.Pos, msg.Length)
			return false
		}
	}

	var inlineQosSize uint32
	if inlineQosFlag {
		if !UpdateCacheChangeFromInlineQos(ch, msg, &inlineQosSize) {
			log.Printf("SubMessage Data ERROR, Inline Qos ParameterList error")
			return false
		}
	}

	payloadSize := smh.SubmessageLength -
		(KRTPSMessageDataExtraInlineQosSize + uint32(octetsToInLineQos) + uint32(inlineQosSize))

	if !keyFlag {
		nextPos := msg.Pos + payloadSize
		if (msg.Length >= nextPos) && (payloadSize > 0) {
			ch.Kind = common.KAlive
			ch.SerializedPayload.Data = msg.Buffer[msg.Pos:]
			ch.SerializedPayload.Length = payloadSize
			ch.SerializedPayload.MaxSize = payloadSize
			ch.SetFragmentSize(fragmentSize, false)

			msg.Pos = nextPos
		} else {
			log.Printf("Serialized Payload value invalid or larger than maximum allowed size (%v/%v)",
				payloadSize, msg.Length-msg.Pos)
			return false
		}
	} else if keyFlag {
		/* XXX TODO
		   Endianness_t previous_endian = msg->msg_endian;
		   if (ch->serializedPayload.encapsulation == PL_CDR_BE)
		   msg->msg_endian = BIGEND;
		   else if (ch->serializedPayload.encapsulation == PL_CDR_LE)
		   msg->msg_endian = LITTLEEND;
		   else
		   {
		   logError(RTPS_MSG_IN, IDSTRING"Bad encapsulation for KeyHash and status parameter list");
		   return false;
		   }
		   //uint32_t param_size;
		   if (ParameterList::readParameterListfromCDRMsg(msg, &m_ParamList, ch, false) <= 0)
		   {
		   logInfo(RTPS_MSG_IN, IDSTRING"SubMessage Data ERROR, keyFlag ParameterList");
		   return false;
		   }
		   msg->msg_endian = previous_endian;
		*/
	}

	// Set sourcetimestamp
	if receiver.haveTimeStamp {
		ch.SourceTimestamp = receiver.timeStamp
	}

	log.Printf("from Writer %v ; possible RTPSReader entities: %v", ch.WriterGUID, len(receiver.associatedReaders))
	receiver.dataFragMsgFunc(&readerID, ch, sampleSize, fragmentStartingNum, fragmentsInSubmessage)
	ch.SerializedPayload.Data = nil

	log.Println("Sub Message DATA_FRAG processed")

	return true
}

func (receiver *Receiver) procSubmsgData(msg *common.CDRMessage, smh *SubmessageHeaderT) bool {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	// READ and PROCESS
	if smh.SubmessageLength < KRTPSMessageDataMinLength {
		log.Print("Too short submessage received, ignoring")
		return false
	}
	// Fill flags bool values
	endiannessFlag := ((smh.Flags & common.Bit(0)) != 0)
	inlineQosFlag := ((smh.Flags & common.Bit(1)) != 0)
	dataFlag := ((smh.Flags & common.Bit(2)) != 0)
	keyFlag := ((smh.Flags & common.Bit(3)) != 0)
	if keyFlag && dataFlag {
		log.Println("Message received with Data and Key Flag set, ignoring")
		return false
	}

	// assign message endianness
	if endiannessFlag {
		msg.MsgEndian = common.LITTLEEND
	} else {
		msg.MsgEndian = common.BIGEND
	}

	// Extra flags don't matter now. Avoid those bytes
	msg.Pos += 2

	valid := true
	octetsToInLineQos := int16(0)
	valid = valid && readInt16(msg, &octetsToInLineQos) // it should be 16 in this implementation

	// reader and writer ID
	var readerID common.EntityIDT
	valid = valid && readEntityID(msg, &readerID)

	// WE KNOW THE READER THAT THE MESSAGE IS DIRECTED TO SO WE LOOK FOR IT:
	_, ok := receiver.willAReaderAcceptMsgDirectedTo(&readerID)
	if !ok {
		return false
	}

	// FOUND THE READER.
	// We ask the reader for a cachechange to store the information.
	ch := common.NewCacheChangeT()
	ch.Kind = common.KAlive
	ch.WriterGUID.Prefix = receiver.sourceGUIDPrefix
	valid = valid && readEntityID(msg, &ch.WriterGUID.EntityID)

	// Get sequence number
	valid = valid && readSequenceNumber(msg, &ch.SequenceNumber)

	if !valid {
		return false
	}

	if ch.SequenceNumber.Value >= ^uint64(0) {
		log.Printf("Invalid message received, bad sequence Number")
		return false
	}

	// Jump ahead if more parameters are before inlineQos (not in this version, maybe if further minor versions.)
	if octetsToInLineQos > KRtpsMessageOctetStoinlineqosDataSubMsg {
		msg.Pos += (uint32(octetsToInLineQos) - KRtpsMessageOctetStoinlineqosDataSubMsg)
		if msg.Pos > msg.Length {
			log.Printf("Invalid jump through msg, msg->pos %v > msg->length %v", msg.Pos, msg.Length)
			return false
		}
	}

	inlineQosSize := uint32(0)

	if inlineQosFlag {
		if !UpdateCacheChangeFromInlineQos(ch, msg, &inlineQosSize) {
			log.Println("SubMessage Data ERROR, Inline Qos ParameterList error")
			return false
		}
	}

	if dataFlag || keyFlag {
		payloadSize := smh.SubmessageLength -
			(KRTPSMessageDataExtraInlineQosSize + uint32(octetsToInLineQos) + inlineQosSize)

		if dataFlag {
			nextPos := msg.Pos + payloadSize
			if (msg.Length >= nextPos) && (payloadSize > 0) {
				ch.SerializedPayload.Data = msg.Buffer[msg.Pos:]
				ch.SerializedPayload.Length = payloadSize
				ch.SerializedPayload.MaxSize = payloadSize
				msg.Pos = nextPos
			} else {
				log.Printf("Serialized Payload value invalid or larger than maximum allowed size (%v/%v)",
					payloadSize, (msg.Length - msg.Pos))
				return false
			}
		} else if keyFlag {
			if payloadSize <= 0 {
				log.Printf("Serialized Payload value invalid (%v)", payloadSize)
				return false
			}

			if payloadSize <= policy.KParameterKeyHashLength {
				for i := uint32(0); i < payloadSize; i++ {
					ch.InstanceHandle.Value[i] = msg.Buffer[msg.Pos+i]
				}
			} else {
				log.Printf("Ignoring Serialized Payload for too large key-only data (%v)", payloadSize)
			}
			msg.Pos += payloadSize
		}
	}

	if receiver.haveTimeStamp {
		ch.SourceTimestamp = receiver.timeStamp
	}

	log.Printf("from Writer %v; possible RTPSReader entities: %v", ch.WriterGUID, receiver.associatedReaders)

	// look for the correct reader to add the change
	receiver.dataMsgFunc(&readerID, ch)

	payloadPool := ch.PayloadOwner()
	if payloadPool != nil {
		payloadPool.ReleasePayload(ch)
	}

	ch.SerializedPayload.Data = []common.Octet{}
	log.Printf("Sub Message DATA processed")
	return true
}

func (receiver *Receiver) procSubmsgAcknack(msg *common.CDRMessage, smh *SubmessageHeaderT) bool {
	endiannessFlag := ((smh.Flags & common.Bit(0)) != 0)
	finalFlag := ((smh.Flags & common.Bit(1)) != 0)
	if endiannessFlag {
		msg.MsgEndian = common.LITTLEEND
	} else {
		msg.MsgEndian = common.BIGEND
	}

	var readerGUID common.GUIDT
	readerGUID.Prefix = receiver.sourceGUIDPrefix
	readEntityID(msg, &readerGUID.EntityID)

	var writerGUID common.GUIDT
	writerGUID.Prefix = receiver.destGUIDPrefix
	readEntityID(msg, &writerGUID.EntityID)

	snsSet := readSequenceNumberSet(msg)
	var ackCount uint32
	readUInt32(msg, &ackCount)

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()
	for _, writer := range receiver.associatedWriters {
		result, valid := writer.ProcessAcknack(&writerGUID, &readerGUID, ackCount, snsSet, finalFlag)
		if valid {
			if !result {
				log.Print("Acknack msg to NOT stateful writer ")
			}
			return result
		}
	}

	log.Printf("Acknack msg to UNKNOWN writer (I loooked through %v writers in this ListenResource)",
		len(receiver.associatedWriters))

	return false
}

func (receiver *Receiver) procSubmsgNackFrag(msg *common.CDRMessage, smh *SubmessageHeaderT) bool {
	endianness := (smh.Flags & common.Bit(0)) != 0
	if endianness {
		msg.MsgEndian = common.LITTLEEND
	} else {
		msg.MsgEndian = common.BIGEND
	}

	var readerGUID common.GUIDT
	readerGUID.Prefix = receiver.sourceGUIDPrefix
	readEntityID(msg, &readerGUID.EntityID)

	var writerGUID common.GUIDT
	writerGUID.Prefix = receiver.destGUIDPrefix
	readEntityID(msg, &writerGUID.EntityID)

	var writerSN common.SequenceNumberT
	readSequenceNumber(msg, &writerSN)

	fnState := common.NewFragmentNumberSet()
	readFragmentNumberSet(msg, fnState)

	var ackCount uint32
	readUInt32(msg, &ackCount)

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()
	for _, writer := range receiver.associatedWriters {
		result, _ := writer.ProcessNackFrag(&writerGUID, &readerGUID, ackCount, &writerSN, fnState)
		if !result {
			log.Print("Acknack msg to NOT stateful writer ")
		}
		return result
	}

	log.Printf("Acknack msg to UNKNOWN writer (I looked through %d writers in this ListenResource)",
		len(receiver.associatedWriters))
	return false
}

func (receiver *Receiver) procSubmsgGap(msg *common.CDRMessage, smh *SubmessageHeaderT) bool {
	endiannessFlag := (smh.Flags & common.Bit(0)) != 0
	if endiannessFlag {
		msg.MsgEndian = common.LITTLEEND
	} else {
		msg.MsgEndian = common.BIGEND
	}

	var readerGUID common.GUIDT
	readerGUID.Prefix = receiver.destGUIDPrefix
	readEntityID(msg, &readerGUID.EntityID)

	var writerGUID common.GUIDT
	writerGUID.Prefix = receiver.sourceGUIDPrefix
	readEntityID(msg, &writerGUID.EntityID)

	var gapStart common.SequenceNumberT
	readSequenceNumber(msg, &gapStart)
	gapList := readSequenceNumberSet(msg)
	if gapStart.Value <= 0 {
		return false
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	callback := func(reader IRtpsMsgReader) {
		reader.ProcessGapMsg(&writerGUID, &gapStart, gapList)
	}
	receiver.findAllReaders(&writerGUID.EntityID, callback)

	return true
}

func (receiver *Receiver) procSubmsgHeartbeatFrag(msg *common.CDRMessage, smh *SubmessageHeaderT) bool {
	endiannessFlag := (smh.Flags & common.Bit(0)) != 0
	if endiannessFlag {
		msg.MsgEndian = common.LITTLEEND
	} else {
		msg.MsgEndian = common.BIGEND
	}

	var readerGUID common.GUIDT
	readerGUID.Prefix = receiver.destGUIDPrefix
	readEntityID(msg, &readerGUID.EntityID)

	var writerGUID common.GUIDT
	writerGUID.Prefix = receiver.sourceGUIDPrefix
	readEntityID(msg, &writerGUID.EntityID)

	var writerSN common.SequenceNumberT
	readSequenceNumber(msg, &writerSN)

	var lastFN common.FragmentNumberT
	readUInt32(msg, &lastFN)

	var HBCount uint32
	readUInt32(msg, &HBCount)

	return true
}

func (receiver *Receiver) procSubmsgInfoDst(msg *common.CDRMessage, smh *SubmessageHeaderT) bool {
	endianness := (smh.Flags & common.Bit(0)) != 0
	if endianness {
		msg.MsgEndian = common.LITTLEEND
	} else {
		msg.MsgEndian = common.BIGEND
	}

	var guidP common.GUIDPrefixT
	data, _ := readData(msg, common.KGUIDPrefixSize)
	for i := 0; i < len(data); i++ {
		guidP.Value[i] = data[i]
	}

	if !guidP.Equal(&common.KUnknownGUIDPrefix) {
		receiver.destGUIDPrefix = guidP
		log.Printf("DST RTPSParticipant is now: %v", receiver.destGUIDPrefix)
	}

	return true
}

func (receiver *Receiver) procSubmsgInfoTs(msg *common.CDRMessage, smh *SubmessageHeaderT) bool {
	endiannessFlag := (smh.Flags & common.Bit(0)) != 0
	timeFlag := (smh.Flags & common.Bit(1)) != 0
	if endiannessFlag {
		msg.MsgEndian = common.LITTLEEND
	} else {
		msg.MsgEndian = common.BIGEND
	}

	if !timeFlag {
		receiver.haveTimeStamp = true
		readTimestamp(msg, &receiver.timeStamp)
	} else {
		receiver.haveTimeStamp = false
	}

	return true
}

func (receiver *Receiver) procSubmsgInfoSrc(msg *common.CDRMessage, smh *SubmessageHeaderT) bool {
	endiannessFlag := (smh.Flags & common.Bit(0)) != 0
	if endiannessFlag {
		msg.MsgEndian = common.LITTLEEND
	} else {
		msg.MsgEndian = common.BIGEND
	}

	if smh.SubmessageLength == KInfoSrcSubmsgLength {
		msg.Pos += 4
		readOctet(msg, &receiver.sourceVersion.Major)
		readOctet(msg, &receiver.sourceVersion.Minor)

		vendorID, _ := readData(msg, 2)
		receiver.sourceVendorID.Vendor[0] = vendorID[0]
		receiver.sourceVendorID.Vendor[1] = vendorID[1]

		guidP, _ := readData(msg, common.KGUIDPrefixSize)
		for i := 0; i < common.KGUIDPrefixSize; i++ {
			receiver.sourceGUIDPrefix.Value[i] = guidP[i]
		}

		log.Printf("SRC RTPSParticipant is now: %v", receiver.sourceGUIDPrefix)
		return true
	}

	return false
}

// ProcessCDRMsg process a new CDR message
// loc Locator indicating the sending address.
// msg Pointer to the message
func (receiver *Receiver) ProcessCDRMsg(loc *common.Locator, msg *common.CDRMessage) {
	if msg.Length < common.KRTPSMessageHeaderSize {
		log.Fatal("Received message too short, ignoring")
		return
	}

	receiver.reset()

	participantGUIDPrefix := receiver.participant.GetGuid().Prefix
	receiver.destGUIDPrefix = participantGUIDPrefix

	msg.Pos = 0 // Start reading at 0

	// Once everything is set, the reading begins:
	if !receiver.checkRTPSHeader(msg) {
		return
	}

	// Loop until there are no more submessages
	var valid bool
	count := 0
	var submsgh SubmessageHeaderT

	for msg.Pos < msg.Length {
		submessage := msg
		// First 4 bytes must contain: ID | flags | octets to next header
		if !receiver.readSubmessageHeader(submessage, &submsgh) {
			return
		}
		valid = true
		count++
		nextMsgPos := submessage.Pos
		nextMsgPos += (submsgh.SubmessageLength + 3) & (^uint32(3))
		switch submsgh.SubMessageID {
		case KData:
			if receiver.destGUIDPrefix != participantGUIDPrefix {
				log.Println("Data Submsg ignored, DST is another RTPSParticipant")
			} else {
				log.Println("Data Submsg received, processing.")
				valid = receiver.procSubmsgData(submessage, &submsgh)
			}
		case KDataFrag:
			if receiver.destGUIDPrefix != participantGUIDPrefix {
				log.Println("DataFrag Submsg ignored, DST is another RTPSParticipant")
			} else {
				log.Println("DataFrag Submsg received, processing.")
				valid = receiver.procSubmsgDataFrag(submessage, &submsgh)
			}
		case KGap:
			if receiver.destGUIDPrefix != participantGUIDPrefix {
				log.Println("Gap Submsg ignored, DST is another RTPSParticipant...")
			} else {
				log.Println("Gap Submsg received, processing...")
				valid = receiver.procSubmsgGap(submessage, &submsgh)
			}
		case KAcknack:
			if receiver.destGUIDPrefix != participantGUIDPrefix {
				log.Println("Acknack Submsg ignored, DST is another RTPSParticipant...")
			} else {
				log.Println("Acknack Submsg received, processing...")
				valid = receiver.procSubmsgAcknack(submessage, &submsgh)
			}
		case KNackFrag:
			if receiver.destGUIDPrefix != participantGUIDPrefix {
				log.Println("NackFrag Submsg ignored, DST is another RTPSParticipant...")
			} else {
				log.Println("NackFrag Submsg received, processing...")
				valid = receiver.procSubmsgNackFrag(submessage, &submsgh)
			}
		case KHeartbeat:
			if receiver.destGUIDPrefix != participantGUIDPrefix {
				log.Println("HB Submsg ignored, DST is another RTPSParticipant...")
			} else {
				log.Println("Heartbeat Submsg received, processing...")
				valid = receiver.procSubmsgHeartbeat(submessage, &submsgh)
			}
		case KHeartbeatFrag:
			if receiver.destGUIDPrefix != participantGUIDPrefix {
				log.Println("HBFrag Submsg ignored, DST is another RTPSParticipant...")
			} else {
				log.Println("HeartbeatFrag Submsg received, processing...")
				valid = receiver.procSubmsgHeartbeatFrag(submessage, &submsgh)
			}
		case KPad:
			log.Println("PAD messages not yet implemented, ignoring")
		case KInfoDst:
			log.Println("InfoDST message received, processing...")
			valid = receiver.procSubmsgInfoDst(submessage, &submsgh)
		case KInfoSrc:
			log.Println("InfoSRC message received, processing...")
			valid = receiver.procSubmsgInfoSrc(submessage, &submsgh)
		case KInfoTs:
			log.Println("InfoTS Submsg received, processing...")
			valid = receiver.procSubmsgInfoTs(submessage, &submsgh)
		case KInfoReply:
		case KInfoReplyIP4:
		}

		if !valid || submsgh.IsLast {
			break
		}
		submessage.Pos = nextMsgPos
	}

	receiver.participant.AssertRemoteParticipantLiveliness(&receiver.sourceGUIDPrefix)
}
