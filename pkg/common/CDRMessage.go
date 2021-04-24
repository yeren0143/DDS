package common

/*
#include <stdio.h>
#include <string.h>

void readBufferReversed(char* dest, char* src, size_t length)
{
	for (size_t i = 0; i < length; ++i) {
		*(dest + i) = *(src + length - 1 - i);
	}
}

void readeBuffer(char* dest, char* src, size_t length)
{
	memcpy(dest, src, length);
}

*/
import "C"

import (
	"log"
	"unsafe"
)

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
		msg.Buffer[msg.Pos] = o0
		msg.Buffer[msg.Pos+1] = o1
		msg.Buffer[msg.Pos+2] = o2
		msg.Buffer[msg.Pos+3] = o3
	} else {
		msg.Buffer[msg.Pos] = o3
		msg.Buffer[msg.Pos+1] = o2
		msg.Buffer[msg.Pos+2] = o1
		msg.Buffer[msg.Pos+3] = o0
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
		msg.Buffer[msg.Pos] = o0
		msg.Buffer[msg.Pos+1] = o1
		msg.Buffer[msg.Pos+2] = o2
		msg.Buffer[msg.Pos+3] = o3
	} else {
		msg.Buffer[msg.Pos] = o3
		msg.Buffer[msg.Pos+1] = o2
		msg.Buffer[msg.Pos+2] = o1
		msg.Buffer[msg.Pos+3] = o0
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
		msg.Buffer[msg.Pos] = o0
		msg.Buffer[msg.Pos+1] = o1
		msg.Buffer[msg.Pos+2] = o2
		msg.Buffer[msg.Pos+3] = o3
		msg.Buffer[msg.Pos] = o4
		msg.Buffer[msg.Pos+1] = o5
		msg.Buffer[msg.Pos+2] = o6
		msg.Buffer[msg.Pos+3] = o7
	} else {
		msg.Buffer[msg.Pos] = o7
		msg.Buffer[msg.Pos+1] = o6
		msg.Buffer[msg.Pos+2] = o5
		msg.Buffer[msg.Pos+3] = o4
		msg.Buffer[msg.Pos] = o3
		msg.Buffer[msg.Pos+1] = o2
		msg.Buffer[msg.Pos+2] = o1
		msg.Buffer[msg.Pos+3] = o0
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

func (msg *CDRMessage) AddData(data []Octet) bool {
	if msg.Pos+uint32(len(data)) > msg.MaxSize {
		return false
	}

	for i := 0; i < len(data); i++ {
		msg.AddOctet(data[i])
	}
	msg.Pos += uint32(len(data))
	msg.Length += uint32(len(data))
	return true
}

func (msg *CDRMessage) AddOctet(data Octet) bool {
	if msg.Pos+1 > msg.MaxSize {
		return false
	}

	msg.Buffer[msg.Pos] = data
	msg.Pos++
	msg.Length++
	return true
}

func (msg *CDRMessage) AddLocator(loc *Locator) bool {
	msg.AddInt32(int32(loc.Kind))
	msg.AddUInt32(loc.Port)
	msg.AddData(loc.Address[:16])
	return true
}

func (msg *CDRMessage) ReadLocator(loc *Locator) bool {
	if msg.Pos+24 > msg.Length {
		return false
	}
	valid := msg.ReadInt32(&loc.Kind)
	valid = valid && msg.ReadUInt32(&loc.Port)
	dataum, ok := msg.ReadData(16)
	copy(loc.Address[:16], dataum[:])
	return valid && ok
}

func (msg *CDRMessage) ReadOctet(oc *Octet) bool {
	if msg.Pos+1 > msg.Length {
		return false
	}
	*oc = msg.Buffer[msg.Pos]
	msg.Pos++
	return true
}

func ReadDataReversed(msg *CDRMessage, dst *Octet, length uint32) bool {
	cDest := (*C.char)(unsafe.Pointer(dst))
	cSrc := (*C.char)(unsafe.Pointer(&msg.Buffer[msg.Pos]))
	CSize := (*C.size_t)(unsafe.Pointer(&length))
	C.readBufferReversed(cDest, cSrc, *CSize)
	msg.Pos += length
	return true
}

func (msg *CDRMessage) ReadUInt32(lo *uint32) bool {
	if msg.Pos+4 > msg.Length {
		return false
	}

	o1 := msg.Buffer[msg.Pos]
	o2 := msg.Buffer[msg.Pos+1]
	o3 := msg.Buffer[msg.Pos+2]
	o4 := msg.Buffer[msg.Pos+3]

	if msg.MsgEndian == KDefaultEndian {
		*lo = uint32(o1) + uint32(o2<<8) + uint32(o3<<16) + uint32(o4<<24)
	} else {
		*lo = uint32(o1<<24) + uint32(o2<<16) + uint32(o3<<8) + uint32(o1)
	}
	msg.Pos += 4

	return true
}

func (msg *CDRMessage) ReadInt32(lo *int32) bool {
	if msg.Pos+4 > msg.Length {
		return false
	}

	o1 := msg.Buffer[msg.Pos]
	o2 := msg.Buffer[msg.Pos+1]
	o3 := msg.Buffer[msg.Pos+2]
	o4 := msg.Buffer[msg.Pos+3]

	if msg.MsgEndian == KDefaultEndian {
		*lo = int32(o1) + int32(o2<<8) + int32(o3<<16) + int32(o4<<24)
	} else {
		*lo = int32(o1<<24) + int32(o2<<16) + int32(o3<<8) + int32(o1)
	}
	msg.Pos += 4

	return true
}

func (msg *CDRMessage) ReadUInt64(lolo *uint64) bool {
	if msg.Pos+8 > msg.Length {
		return false
	}

	o1 := msg.Buffer[msg.Pos]
	o2 := msg.Buffer[msg.Pos+1]
	o3 := msg.Buffer[msg.Pos+2]
	o4 := msg.Buffer[msg.Pos+3]
	o5 := msg.Buffer[msg.Pos+4]
	o6 := msg.Buffer[msg.Pos+5]
	o7 := msg.Buffer[msg.Pos+6]
	o8 := msg.Buffer[msg.Pos+7]

	if msg.MsgEndian == KDefaultEndian {
		*lolo = uint64(o1) + uint64(o2<<8) + uint64(o3<<16) + uint64(o4<<24) +
			uint64(o5<<32) + uint64(o6<<40) + uint64(o7<<48) + uint64(o8<<56)
	} else {
		*lolo = uint64(o8) + uint64(o7<<8) + uint64(o6<<16) + uint64(o5<<24) +
			uint64(o4<<32) + uint64(o3<<40) + uint64(o2<<48) + uint64(o1<<56)
	}
	msg.Pos += 8
	return true
}

func (msg *CDRMessage) ReadInt64(lolo *int64) bool {
	if msg.Pos+8 > msg.Length {
		return false
	}

	o1 := msg.Buffer[msg.Pos]
	o2 := msg.Buffer[msg.Pos+1]
	o3 := msg.Buffer[msg.Pos+2]
	o4 := msg.Buffer[msg.Pos+3]
	o5 := msg.Buffer[msg.Pos+4]
	o6 := msg.Buffer[msg.Pos+5]
	o7 := msg.Buffer[msg.Pos+6]
	o8 := msg.Buffer[msg.Pos+7]

	if msg.MsgEndian == KDefaultEndian {
		*lolo = int64(o1) + int64(o2<<8) + int64(o3<<16) + int64(o4<<24) +
			int64(o5<<32) + int64(o6<<40) + int64(o7<<48) + int64(o8<<56)
	} else {
		*lolo = int64(o8) + int64(o7<<8) + int64(o6<<16) + int64(o5<<24) +
			int64(o4<<32) + int64(o3<<40) + int64(o2<<48) + int64(o1<<56)
	}

	msg.Pos += 8

	return true
}

func (msg *CDRMessage) ReaderLocator(loc *Locator) bool {
	if msg.Pos+24 > msg.Length {
		return false
	}

	valid := msg.ReadInt32(&loc.Kind)
	valid = valid && msg.ReadUInt32(&loc.Port)
	addr, ok := msg.ReadData(16)
	if !ok {
		log.Fatalln("ReaderLocator failed")
	}
	copy(loc.Address[:], addr[:16])
	valid = valid && ok

	return valid
}

func (msg *CDRMessage) ReadData(length uint32) ([]Octet, bool) {
	var data []Octet
	if msg.Pos+length > msg.Length {
		return data, false
	}

	data = append(data, msg.Buffer[msg.Pos:msg.Pos+length]...)
	msg.Pos += length

	return data, true
}

func (msg *CDRMessage) ReadInt16(i16 *int16) bool {
	if msg.Pos+2 > msg.Length {
		return false
	}

	if msg.MsgEndian == KDefaultEndian {
		*i16 = int16(msg.Buffer[msg.Pos+1]<<8) + int16(msg.Buffer[msg.Pos])
	} else {
		*i16 = int16(msg.Buffer[msg.Pos+1]<<8) + int16(msg.Buffer[msg.Pos])
	}
	msg.Pos += 2
	return true
}

func (msg *CDRMessage) ReadUInt16(i16 *uint16) bool {
	if msg.Pos+2 > msg.Length {
		return false
	}

	if msg.MsgEndian == KDefaultEndian {
		*i16 = uint16(msg.Buffer[msg.Pos+1]<<8) + uint16(msg.Buffer[msg.Pos])
	} else {
		*i16 = uint16(msg.Buffer[msg.Pos+1]) + uint16(msg.Buffer[msg.Pos+1])
	}
	msg.Pos += 2
	return true
}

func (msg *CDRMessage) ReadEntityID(id *EntityIDT) bool {
	if msg.Pos+4 > msg.Length {
		return false
	}

	for i := 0; i < len(id.Value); i++ {
		id.Value[i] = msg.Buffer[msg.Pos+uint32(i)]
	}
	msg.Pos += 4

	return true
}

func (msg *CDRMessage) ReadSequenceNumber(sn *SequenceNumberT) bool {
	if (msg.Pos + 8) > msg.Length {
		return false
	}
	var high int32
	var low uint32
	valid := msg.ReadInt32(&high)
	valid = valid && msg.ReadUInt32(&low)
	sn.Value = (uint64(high) << 32) + uint64(low)

	return valid
}

func (msg *CDRMessage) ReadSequenceNumberSet() *SequenceNumberSet {
	valid := true

	var seqNum SequenceNumberT
	valid = valid && msg.ReadSequenceNumber(&seqNum)

	var numBits uint32
	valid = valid && msg.ReadUInt32(&numBits)
	valid = valid && (numBits <= 256)

	longs := (numBits + 31) / 32
	var bitmap [8]uint32
	for i := uint32(0); valid && (i < longs); i++ {
		valid = valid && msg.ReadUInt32(&bitmap[i])
	}

	sns := NewSequenceNumberSet(KSequenceNumberUnknown)
	if valid {
		sns.Base(seqNum)
		sns.BitmapSet(numBits, &bitmap)
	}
	return sns
}

func (msg *CDRMessage) ReadTimestamp(ts *Time) bool {
	valid := true
	valid = valid && msg.ReadInt32(&ts.Seconds)
	valid = valid && msg.ReadUInt32(&ts.Nanosec)
	return valid
}

func (msg *CDRMessage) ReadFragmentNumberSet(fragmentNumberSet *FragmentNumberSet) bool {
	valid := true

	var base FragmentNumberT
	valid = valid && msg.ReadUInt32(&base)

	var numBits uint32
	valid = valid && msg.ReadUInt32(&numBits)
	valid = valid && (numBits <= 256)

	nlongs := (numBits + 31) / 32
	var bitmap [8]uint32
	for i := uint32(0); valid && (i < nlongs); i++ {
		valid = valid && msg.ReadUInt32(&bitmap[i])
	}

	if valid {
		fragmentNumberSet.Base(base)
		fragmentNumberSet.BitmapSet(numBits, &bitmap)
	}

	return valid
}

func (msg *CDRMessage) ReadString() (string, bool) {
	var strSize uint32
	stri := ""
	valid := msg.ReadUInt32(&strSize)
	if msg.Pos+strSize > msg.Length {
		return stri, false
	}

	if strSize > 1 {
		stri = string(msg.Buffer[msg.Pos : msg.Pos+strSize])
	}
	msg.Pos += strSize
	msg.Pos = uint32(int(msg.Pos+3) & (^3))

	return stri, valid
}

func InitCDRMsg(msg *CDRMessage, payloadSize uint32) bool {
	if len(msg.Buffer) == 0 {
		msg.Buffer = make([]Octet, payloadSize+KRTPSMessageCommonDataPayloadSize)
		msg.MaxSize = payloadSize + KRTPSMessageCommonRTPSPayloadSize
	}
	msg.Pos = 0
	msg.Length = 0
	msg.MsgEndian = KDefaultEndian
	return true
}

func (msg *CDRMessage) Init(buffer []Octet, size uint32) {
	if len(buffer) == 0 {
		return
	}

	msg.Wraps = true
	msg.Pos = 0
	msg.Length = 0
	msg.Buffer = buffer
	msg.MaxSize = size
	msg.MsgEndian = KDefaultEndian
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
