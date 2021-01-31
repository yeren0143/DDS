package message

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
	"github.com/yeren0143/DDS/common"
	"unsafe"
)

const (
	KRTPSMessageDefaultSize                     = 10500 // max size of rtps message in bytes
	KRTPSMessageCommonRTPSPayloadSize           = 536   // common payload a rtps message has TODO(Ricardo) It is necessary?
	KRTPSMessageCommonDataPayloadSize           = 10000 // common data size
	KRTPSMessageHeaderSize                      = 20    // header size in bytes
	KRTPSMessageSubMessageHeaderSize            = 4
	KRTPSMessageDataExtraInlineQosSize          = 4
	KRTPSMessageInfoTSSize                      = 12
	KRtpsMessageOctetStoinlineqosDataSubMsg     = 16
	KRtpsMessageOctetStoinlineqosDataFragSubmsg = 28
	KRTPSMessageDataMinLength                   = 24
)

func initCDRMsg(msg *common.CDRMessage, payloadSize uint32) bool {
	if len(msg.Buffer) == 0 {
		msg.Buffer = make([]common.Octet, payloadSize+KRTPSMessageCommonRTPSPayloadSize)
		msg.MaxSize = payloadSize + KRTPSMessageCommonRTPSPayloadSize
	}
	msg.Pos = 0
	msg.Length = 0
	msg.MsgEndian = common.KDefaultEndian
	return true
}

func readData(msg *common.CDRMessage, length uint32) ([]common.Octet, bool) {
	var data []common.Octet
	if msg.Pos+length > msg.Length {
		return data, false
	}

	data = append(data, msg.Buffer[msg.Pos:msg.Pos+length]...)
	msg.Pos += length

	return data, true
}

func readOctet(msg *common.CDRMessage, data *common.Octet) bool {
	if (msg.Pos + 1) > msg.Length {
		return false
	}

	*data = msg.Buffer[msg.Pos]
	msg.Pos++
	return true
}

func readUInt16(msg *common.CDRMessage, i16 *uint16) bool {
	if (msg.Pos + 2) > msg.Length {
		return false
	}
	if msg.MsgEndian == common.KDefaultEndian {
		low := msg.Buffer[msg.Pos]
		high := msg.Buffer[msg.Pos+1]
		*i16 = uint16(low) + uint16(high<<8)
	} else {
		high := msg.Buffer[msg.Pos]
		low := msg.Buffer[msg.Pos+1]
		*i16 = uint16(low) + uint16(high<<8)
	}
	msg.Pos += 2
	return true
}

func readInt16(msg *common.CDRMessage, i16 *int16) bool {
	if msg.Pos+2 > msg.Length {
		return false
	}
	if msg.MsgEndian == common.KDefaultEndian {
		low := int16(msg.Buffer[msg.Pos])
		high := int16(msg.Buffer[msg.Pos+1])
		*i16 = low + (high << 8)
	} else {
		high := int16(msg.Buffer[msg.Pos])
		low := int16(msg.Buffer[msg.Pos+1])
		*i16 = low + (high << 8)
	}
	msg.Pos += 2
	return true
}

func readInt32(msg *common.CDRMessage, dest *int32) bool {
	if (msg.Pos + 4) > msg.Length {
		return false
	}
	if msg.MsgEndian == common.KDefaultEndian {
		part0 := int32(msg.Buffer[msg.Pos])
		part1 := int32(msg.Buffer[msg.Pos+1]) << 8
		part2 := int32(msg.Buffer[msg.Pos+2]) << 16
		part3 := int32(msg.Buffer[msg.Pos+3]) << 24
		*dest = part0 + part1 + part2 + part3
		msg.Pos += 4
	} else {
		readDataReversed(msg, (*common.Octet)(unsafe.Pointer(dest)), 4)
	}
	return true
}

func readUInt32(msg *common.CDRMessage, dest *uint32) bool {
	if (msg.Pos + 4) > msg.Length {
		return false
	}
	if msg.MsgEndian == common.KDefaultEndian {
		part0 := uint32(msg.Buffer[msg.Pos])
		part1 := uint32(msg.Buffer[msg.Pos+1]) << 8
		part2 := uint32(msg.Buffer[msg.Pos+2]) << 16
		part3 := uint32(msg.Buffer[msg.Pos+3]) << 24
		*dest = part0 + part1 + part2 + part3
		msg.Pos += 4
	} else {
		readDataReversed(msg, (*common.Octet)(unsafe.Pointer(dest)), 4)
	}
	return true
}

func readInt64(msg *common.CDRMessage, dest *int64) bool {
	if (msg.Pos + 8) > msg.Length {
		return false
	}

	if msg.MsgEndian == common.KDefaultEndian {
		part0 := int64(msg.Buffer[msg.Pos])
		part1 := int64(msg.Buffer[msg.Pos+1]) << 8
		part2 := int64(msg.Buffer[msg.Pos+2]) << 16
		part3 := int64(msg.Buffer[msg.Pos+3]) << 24
		part4 := int64(msg.Buffer[msg.Pos+4]) << 32
		part5 := int64(msg.Buffer[msg.Pos+5]) << 40
		part6 := int64(msg.Buffer[msg.Pos+6]) << 48
		part7 := int64(msg.Buffer[msg.Pos+7]) << 56
		*dest = part0 + part1 + part2 + part3 + part4 + part5 + part6 + part7
		msg.Pos += 8
	} else {
		readDataReversed(msg, (*common.Octet)(unsafe.Pointer(dest)), 8)
	}
	return true
}

func readUInt64(msg *common.CDRMessage, dest *uint64) bool {
	if (msg.Pos + 8) > msg.Length {
		return false
	}

	if msg.MsgEndian == common.KDefaultEndian {
		part0 := uint64(msg.Buffer[msg.Pos])
		part1 := uint64(msg.Buffer[msg.Pos+1]) << 8
		part2 := uint64(msg.Buffer[msg.Pos+2]) << 16
		part3 := uint64(msg.Buffer[msg.Pos+3]) << 24
		part4 := uint64(msg.Buffer[msg.Pos+4]) << 32
		part5 := uint64(msg.Buffer[msg.Pos+5]) << 40
		part6 := uint64(msg.Buffer[msg.Pos+6]) << 48
		part7 := uint64(msg.Buffer[msg.Pos+7]) << 56
		*dest = part0 + part1 + part2 + part3 + part4 + part5 + part6 + part7
		msg.Pos += 8
	} else {
		readDataReversed(msg, (*common.Octet)(unsafe.Pointer(dest)), 8)
	}
	return true
}

func readDataReversed(msg *common.CDRMessage, dst *common.Octet, length uint32) bool {
	cDest := (*C.char)(unsafe.Pointer(dst))
	cSrc := (*C.char)(unsafe.Pointer(&msg.Buffer[msg.Pos]))
	CSize := (*C.size_t)(unsafe.Pointer(&length))
	C.readBufferReversed(cDest, cSrc, *CSize)
	msg.Pos += length
	return true
}

func readEntityID(msg *common.CDRMessage, id *common.EntityIDT) bool {
	if msg.Pos+4 > msg.Length {
		return false
	}

	for i := 0; i < len(id.Value); i++ {
		id.Value[i] = msg.Buffer[msg.Pos+uint32(i)]
	}
	msg.Pos += 4

	return true
}

func readSequenceNumber(msg *common.CDRMessage, sn *common.SequenceNumberT) bool {
	if (msg.Pos + 8) > msg.Length {
		return false
	}
	var high int32
	var low uint32
	valid := readInt32(msg, &high)
	valid = valid && readUInt32(msg, &low)
	sn.Value = (uint64(high) << 32) + uint64(low)

	return valid
}

func readSequenceNumberSet(msg *common.CDRMessage) *common.SequenceNumberSet {
	valid := true

	var seqNum common.SequenceNumberT
	valid = valid && readSequenceNumber(msg, &seqNum)

	var numBits uint32
	valid = valid && readUInt32(msg, &numBits)
	valid = valid && (numBits <= 256)

	longs := (numBits + 31) / 32
	var bitmap [8]uint32
	for i := uint32(0); valid && (i < longs); i++ {
		valid = valid && readUInt32(msg, &bitmap[i])
	}

	sns := common.NewSequenceNumberSet(common.KSequenceNumberUnknown)
	if valid {
		sns.Base(seqNum)
		sns.BitmapSet(numBits, &bitmap)
	}
	return sns
}

func readFragmentNumberSet(msg *common.CDRMessage, fragmentNumberSet *common.FragmentNumberSet) bool {
	valid := true

	var base common.FragmentNumberT
	valid = valid && readUInt32(msg, &base)

	var numBits uint32
	valid = valid && readUInt32(msg, &numBits)
	valid = valid && (numBits <= 256)

	nlongs := (numBits + 31) / 32
	var bitmap [8]uint32
	for i := uint32(0); valid && (i < nlongs); i++ {
		valid = valid && readUInt32(msg, &bitmap[i])
	}

	if valid {
		fragmentNumberSet.Base(base)
		fragmentNumberSet.BitmapSet(numBits, &bitmap)
	}

	return valid
}

func readTimestamp(msg *common.CDRMessage, ts *common.Time) bool {
	valid := true
	valid = valid && readInt32(msg, &ts.Seconds)
	valid = valid && readUInt32(msg, &ts.Nanosec)
	return valid
}
