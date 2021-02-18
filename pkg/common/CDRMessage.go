package common

const (
	//KRTPSMessageDefaultSize define max size of rtps message in bytes
	KRTPSMessageDefaultSize = 10500

	//KRTPSMessageCommonRTPSPayloadSize define common payload a rtps message has TODO(Ricardo) It is necessary?
	KRTPSMessageCommonRTPSPayloadSize = 536

	//KRTPSMessageCommonDataPayloadSize define common data size
	KRTPSMessageCommonDataPayloadSize = 1000

	//KRTPSMessageHeaderSize define header size in bytes
	KRTPSMessageHeaderSize = 20

	//KRTPSMessageSubMessageHeaderSize define submessage header size
	KRTPSMessageSubMessageHeaderSize = 4

	//KRTPSMessageDataExtraInlineQosSize ?
	KRTPSMessageDataExtraInlineQosSize = 4

	//KRTPSMessageInfoTSSize ?
	KRTPSMessageInfoTSSize = 12

	//KRTPSMessageOcetsToInlineQosDataSubMsg ...
	KRTPSMessageOcetsToInlineQosDataSubMsg = 16

	//KRTPSMessageOcetsToInlineQosDataFragSubMsg ...
	KRTPSMessageOcetsToInlineQosDataFragSubMsg = 28

	//KRTPSMessageDataMinLength ...
	KRTPSMessageDataMinLength = 24
)

//CDRMessage contains a serialized message.
type CDRMessage struct {
	//Buffer Pointer to the buffer where the data is stored.
	Buffer []Octet

	//Read or write position.
	Pos uint32

	//Max size of the message.
	MaxSize uint32

	//Size allocated on buffer. May be higher than max_size.
	ReservedSize uint32

	//Current length of the message.
	Length uint32

	//Whether this message is wrapping a buffer managed elsewhere.
	Wraps bool

	//Endianness of the message
	MsgEndian Endianness
}

func (message *CDRMessage) AddData(data []Octet) bool {
	if message.Pos+uint32(len(data)) > message.MaxSize {
		return false
	}

	for i := 0; i < len(data); i++ {
		message.AddOctet(data[i])
	}
	message.Pos += uint32(len(data))
	message.Length += uint32(len(data))
	return true
}

func (message *CDRMessage) AddOctet(data Octet) bool {
	if message.Pos+1 > message.MaxSize {
		return false
	}

	message.Buffer[message.Pos] = data
	message.Pos++
	message.Length++
	return true
}

func (message *CDRMessage) InitCDRMsg(msg *CDRMessage, payloadSize uint32) bool {
	if len(msg.Buffer) == 0 {
		msg.Buffer = make([]Octet, payloadSize+KRTPSMessageCommonDataPayloadSize)
		msg.MaxSize = payloadSize + KRTPSMessageCommonRTPSPayloadSize
	}
	msg.Pos = 0
	msg.Length = 0
	msg.MsgEndian = KDefaultEndian
	return true
}

func (message *CDRMessage) Init(buffer []Octet, size uint32) {
	if len(buffer) == 0 {
		return
	}

	message.Wraps = true
	message.Pos = 0
	message.Length = 0
	message.Buffer = buffer
	message.MaxSize = size
	message.MsgEndian = KDefaultEndian
}

func NewCDRMessageWithPayload(payload *SerializedPayloadT) *CDRMessage {
	var msg CDRMessage
	msg.Wraps = true
	if payload.Encapsulation == PL_CDR_BE {
		msg.MsgEndian = BIGEND
	} else {
		msg.MsgEndian = LITTLEEND
	}
	msg.Pos = payload.Pos
	msg.Length = payload.Length
	msg.Buffer = payload.Data
	msg.MaxSize = payload.MaxSize
	msg.ReservedSize = payload.MaxSize
	return &msg
}

func NewCDRMessage(size uint32) *CDRMessage {
	return &CDRMessage{
		Buffer:       make([]Octet, size),
		Pos:          0,
		Length:       0,
		MaxSize:      size,
		ReservedSize: size,
		MsgEndian:    KDefaultEndian,
	}
}
