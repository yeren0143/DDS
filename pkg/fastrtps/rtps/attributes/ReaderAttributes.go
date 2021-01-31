package attributes

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/core/policy"
	"github.com/yeren0143/DDS/fastrtps/utils"
)

// Class ReaderTimes, defining the times associated with the Reliable Readers events.
type ReaderTimes struct {
	// Initial AckNack delay. Default value 70ms.
	InitialAcknackDelay common.DurationT
	// Delay to be applied when a hearbeat message is received, default value 5ms.
	HeartBeatResponseDelay common.DurationT
}

var KDefaultReaderTimes = ReaderTimes{
	InitialAcknackDelay:    common.DurationT{Nanosec: 70 * 1000 * 1000, Seconds: 0},
	HeartBeatResponseDelay: common.DurationT{Nanosec: 5 * 1000 * 1000, Seconds: 0},
}

// Class ReaderAttributes, to define the attributes of a RTPSReader.
type ReaderAttributes struct {
	EndpointAtt EndpointAttributes

	// Times associated with this reader (only for stateful readers)
	Times ReaderTimes

	livelinessKind          policy.LivelinessQosPolicyKind
	livelinessLeaseDuration common.DurationT
	// Indicates if the reader expects Inline qos, default value 0.
	expectsInlineQos bool

	DisablePositiveAcks      bool
	matchedWritersAllocation utils.ResourceLimitedContainerConfig
}

var KDefaultReaderAttributes = ReaderAttributes{
	livelinessKind:           policy.AUTOMATIC_LIVELINESS_QOS,
	livelinessLeaseDuration:  common.KTimeInfinite,
	expectsInlineQos:         false,
	DisablePositiveAcks:      false,
	EndpointAtt:              KDefaultEndpointAttributes,
	matchedWritersAllocation: utils.KDefaultResourceLimitedContainerConfig,
}
