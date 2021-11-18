package qos

import (
	"dds/core/policy"

	"github.com/golang/glog"
)

/**
 * Class WriterQos, containing all the possible Qos that can be set for a determined Publisher.
 * Although these values can be set and are transmitted
 * during the Endpoint Discovery Protocol, not all of the behaviour associated with them has been implemented in the library.
 * Please consult each of them to check for implementation details and default values.
 * @ingroup FASTRTPS_ATTRIBUTES_MODULE
 */
type WriterQos struct {
	Durability        policy.DurabilityQosPolicy
	DurabilityService policy.DurabilityServiceQosPolicy
	Reliability       policy.ReliabilityQosPolicy
	Liveliness        policy.LivelinessQosPolicy
}

func (qos *WriterQos) Clear() {
	glog.Warning("not impl")
}
