package status

import (
	"github.com/yeren0143/DDS/core/policy"
)

// QosPolicyCount storing the id of the incompatible QoS Policy and the number of times it fails
type QosPolicyCount struct {
	PolicyID policy.QosPolicyIDT
	Count    uint32
}

// CreateQosPolicyCount return QosPolicyCount with INVALID_QOS_POLICY_ID
func CreateQosPolicyCount() QosPolicyCount {
	return QosPolicyCount{
		PolicyID: policy.INVALID_QOS_POLICY_ID,
		Count:    0,
	}
}

// QosPolicyCountSeq is alias of []QosPolicyCount
type QosPolicyCountSeq = []QosPolicyCount

// IncompatibleQosStatus storing the requested incompatible QoS status
type IncompatibleQosStatus struct {
	// Total cumulative number of times the concerned writer discovered a reader for the same topic
	// The requested QoS is incompatible with the one offered by the writer
	TotalCount uint32

	// The change in total_count since the last time the listener was called or the status was read
	TotalCountChange uint32

	LastPolicyID policy.QosPolicyIDT
	Policies     QosPolicyCountSeq
}

// CreateIncompatibleQosStatus return default IncompatibleQosStatus
func CreateIncompatibleQosStatus() IncompatibleQosStatus {
	return IncompatibleQosStatus{
		LastPolicyID: policy.INVALID_QOS_POLICY_ID,
	}
}

// RequestedIncompatibleQosStatus is alias of IncompatibleQosStatus
type RequestedIncompatibleQosStatus = IncompatibleQosStatus

// OfferedIncompatibleQosStatus is alias of IncompatibleQosStatus
type OfferedIncompatibleQosStatus = IncompatibleQosStatus
