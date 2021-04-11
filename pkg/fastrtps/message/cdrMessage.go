package message

import (
	"github.com/yeren0143/DDS/common"
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
