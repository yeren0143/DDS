package message

import (
	"dds/common"
)

// Class RTPSMessageGroup_t that contains the messages used to send multiples changes as one message.
type RTPSMessageGroupT struct {
	SubMessage  common.CDRMessage
	FullMessage common.CDRMessage
}

func (messageGroup *RTPSMessageGroupT) Init(participantGUID *common.GUIDPrefixT) {
	common.InitCDRMsg(&messageGroup.FullMessage, common.KRTPSMessageCommonDataPayloadSize)
	addHeader(&messageGroup.FullMessage, participantGUID, &common.KProtocolVersion, &common.KVendorIDTeProsima)
}

func NewRTPSMessageGroup(buffer []common.Octet, payload uint32, participantGuid *common.GUIDPrefixT) *RTPSMessageGroupT {
	messageGroup := RTPSMessageGroupT{}
	messageGroup.FullMessage.Init(buffer, payload)
	messageGroup.SubMessage.Init(buffer[payload:], payload)

	messageGroup.Init(participantGuid)

	return &messageGroup
}
