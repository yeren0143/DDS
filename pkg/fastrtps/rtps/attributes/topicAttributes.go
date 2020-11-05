package attributes

import (
	. "common"
)

type TopicAttributes struct {
	TopicKind     TopicKind_t
	TopicName     string
	TopicDataType string
	HistoryQos    HistoryQosPolicy
}
