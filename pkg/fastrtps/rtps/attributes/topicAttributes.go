package attributes

import (
	. "github.com/yeren0143/DDS/common"
	. "github.com/yeren0143/DDS/core/policy"
)

type TopicAttributes struct {
	TopicKind     TopicKindT
	TopicName     string
	TopicDataType string
	HistoryQos    HistoryQosPolicy
}
