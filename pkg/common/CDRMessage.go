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

func (msg *CDRMessage) AddInt32(lo int32) bool {
	if msg.Pos+4 > msg.MaxSize {
		return false
	}

	o0 := Octet(lo & 0x000000FF)
	o1 := Octet(lo & 0x0000FF00)
	o2 := Octet(lo & 0x00FF0000)
	o3 := Octet(lo >> 30)

	if msg.MsgEndian == KDefaultEndian {
		msg.Buffer[msg.Pos] = o3
		msg.Buffer[msg.Pos+1] = o2
		msg.Buffer[msg.Pos+2] = o1
		msg.Buffer[msg.Pos+3] = o0
	} else {
		msg.Buffer[msg.Pos] = o0
		msg.Buffer[msg.Pos+1] = o1
		msg.Buffer[msg.Pos+2] = o2
		msg.Buffer[msg.Pos+3] = o3
	}

	msg.Pos += 4
	msg.Length += 4
	return true
}

func (msg *CDRMessage) AddUInt32(lo uint32) bool {
	if msg.Pos+4 > msg.MaxSize {
		return false
	}

	o0 := Octet(lo & 0x000000FF)
	o1 := Octet(lo & 0x0000FF00)
	o2 := Octet(lo & 0x00FF0000)
	o3 := Octet(lo & 0xFF000000)

	if msg.MsgEndian == KDefaultEndian {
		msg.Buffer[msg.Pos] = o3
		msg.Buffer[msg.Pos+1] = o2
		msg.Buffer[msg.Pos+2] = o1
		msg.Buffer[msg.Pos+3] = o0
	} else {
		msg.Buffer[msg.Pos] = o0
		msg.Buffer[msg.Pos+1] = o1
		msg.Buffer[msg.Pos+2] = o2
		msg.Buffer[msg.Pos+3] = o3
	}

	msg.Pos += 4
	msg.Length += 4
	return true
}

func (msg *CDRMessage) AddInt64(lolo int64) bool {
	if msg.Pos+8 > msg.MaxSize {
		return false
	}

	o0 := Octet(lolo & 0x00000000000000FF)
	o1 := Octet(lolo & 0x000000000000FF00)
	o2 := Octet(lolo & 0x0000000000FF0000)
	o3 := Octet(lolo & 0x00000000FF000000)
	o4 := Octet(lolo & 0x000000FF00000000)
	o5 := Octet(lolo & 0x0000FF0000000000)
	o6 := Octet(lolo & 0x00FF000000000000)
	o7 := Octet(lolo >> 62)

	if msg.MsgEndian == KDefaultEndian {
		msg.Buffer[msg.Pos] = o7
		msg.Buffer[msg.Pos+1] = o6
		msg.Buffer[msg.Pos+2] = o5
		msg.Buffer[msg.Pos+3] = o4
		msg.Buffer[msg.Pos] = o3
		msg.Buffer[msg.Pos+1] = o2
		msg.Buffer[msg.Pos+2] = o1
		msg.Buffer[msg.Pos+3] = o0
	} else {
		msg.Buffer[msg.Pos] = o0
		msg.Buffer[msg.Pos+1] = o1
		msg.Buffer[msg.Pos+2] = o2
		msg.Buffer[msg.Pos+3] = o3
		msg.Buffer[msg.Pos] = o4
		msg.Buffer[msg.Pos+1] = o5
		msg.Buffer[msg.Pos+2] = o6
		msg.Buffer[msg.Pos+3] = o7
	}

	msg.Pos += 8
	msg.Length += 8
	return true
}

func (msg *CDRMessage) AddString(instr string) bool {
	strSiz := len(instr) + 1
	valid := msg.AddUInt32(uint32(strSiz))
	valid = valid && msg.AddData([]byte(instr))
	for ; (strSiz & 3) > 0; strSiz++ {
		valid = valid && msg.AddOctet('0')
	}
	return valid
}

func (msg *CDRMessage) AddUInt16(us uint16) bool {
	if msg.Pos+2 > msg.MaxSize {
		return false
	}
	left := Octet(us & 0x00FF)
	right := Octet(us & 0xFF00)
	if msg.MsgEndian == KDefaultEndian {
		msg.Buffer[msg.Pos] = left
		msg.Buffer[msg.Pos+1] = right
	} else {
		msg.Buffer[msg.Pos] = right
		msg.Buffer[msg.Pos+1] = left
	}
	msg.Pos += 2
	msg.Length += 2
	return true
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

func (message *CDRMessage) AddLocator(loc *Locator) bool {
	message.AddInt32(int32(loc.Kind))
	message.AddUInt32(loc.Port)
	message.AddData(loc.Address[:16])
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
