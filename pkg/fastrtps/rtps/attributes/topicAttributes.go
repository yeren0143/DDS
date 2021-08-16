package attributes

import (
	. "dds/common"
	. "dds/core/policy"
)

type TopicAttributes struct {
	TopicKind     TopicKindT
	TopicName     string
	TopicDataType string
	HistoryQos    HistoryQosPolicy
}
