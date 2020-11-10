package attributes

import (
	. "common"
	. "core/policy"
)

type TopicAttributes struct {
	TopicKind     TopicKind_t
	TopicName     string
	TopicDataType string
	HistoryQos    HistoryQosPolicy
}
