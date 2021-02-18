package attributes

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/core/policy"
	"github.com/yeren0143/DDS/fastrtps/rtps/flowcontrol"
	"github.com/yeren0143/DDS/fastrtps/utils"
)

type RTPSWriterPublishMode common.Octet

const (
	KSynchronousWriter RTPSWriterPublishMode = iota
	KAsynchronousWriter
)

// Struct WriterTimes, defining the times associated with the Reliable Writers events.
type WriterTimes struct {
	InitialHeartbeatDelay  common.DurationT
	HeartbeatPeriod        common.DurationT
	NackResponseDelay      common.DurationT
	NackSupressionDuration common.DurationT
}

var KDefaultWriterTimes = WriterTimes{
	InitialHeartbeatDelay: common.DurationT{Nanosec: 11},
	HeartbeatPeriod:       common.DurationT{Seconds: 3},
	NackResponseDelay:     common.DurationT{Nanosec: 5},
}

type WriterAttributes struct {
	EndpointAtt                  EndpointAttributes
	Times                        WriterTimes
	LivelinessKind               policy.LivelinessQosPolicyKind
	LivelinessLeaseDuration      common.DurationT
	LivelinessAnnouncementPeriod common.DurationT
	PubMode                      RTPSWriterPublishMode
	ThroughputController         flowcontrol.ThroughputControllerDescriptor // Throughput controller, always the last one to apply
	DisableHeartbeatPiggyback    bool
	MatchedReadersAllocation     utils.ResourceLimitedContainerConfig
	DisablePositiveAcks          bool
	KeepDuration                 common.DurationT // Keep duration to keep a sample before considering it has been acked
}

func NewWriterAttributes() *WriterAttributes {
	return &WriterAttributes{
		Times:                        KDefaultWriterTimes,
		LivelinessKind:               policy.AUTOMATIC_LIVELINESS_QOS,
		LivelinessLeaseDuration:      common.KTimeInfinite,
		LivelinessAnnouncementPeriod: common.KTimeInfinite,
		PubMode:                      KSynchronousWriter,
		DisableHeartbeatPiggyback:    false,
		DisablePositiveAcks:          false,
		KeepDuration:                 common.KTimeInfinite,
		EndpointAtt: EndpointAttributes{
			EndpointKind:    common.KWriter,
			DurabilityKind:  common.KTransientLocal,
			ReliabilityKind: common.KReliable,
		},
	}
}
