package attributes

import (
	"dds/common"
	"dds/core/policy"
	"dds/fastrtps/utils"
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

	LivelinessKind          policy.LivelinessQosPolicyKind
	LivelinessLeaseDuration common.DurationT
	// Indicates if the reader expects Inline qos, default value 0.
	ExpectsInlineQos bool

	DisablePositiveAcks      bool
	MatchedWritersAllocation utils.ResourceLimitedContainerConfig
}

func NewReaderAttributes() *ReaderAttributes {
	readerAtt := ReaderAttributes{
		LivelinessKind:           policy.AUTOMATIC_LIVELINESS_QOS,
		LivelinessLeaseDuration:  common.KTimeInfinite,
		ExpectsInlineQos:         false,
		DisablePositiveAcks:      false,
		EndpointAtt:              KDefaultEndpointAttributes,
		MatchedWritersAllocation: utils.KDefaultResourceLimitedContainerConfig,
	}

	readerAtt.EndpointAtt.EndpointKind = common.KReader
	readerAtt.EndpointAtt.DurabilityKind = common.KVolatile
	readerAtt.EndpointAtt.ReliabilityKind = common.KBestEffort
	return &readerAtt
}
