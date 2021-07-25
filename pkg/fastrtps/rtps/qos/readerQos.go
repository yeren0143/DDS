package qos

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/core/policy"
)

/**
 * Class ReaderQos, contains all the possible Qos that can be set for a determined Subscriber.
 * Although these values can be set and are transmitted
 * during the Endpoint Discovery Protocol, not all of the behaviour associated with them has been implemented in the library.
 * Please consult each of them to check for implementation details and default values.
 */
type ReaderQos struct {
	// Durability Qos, implemented in the library.
	Durability policy.DurabilityQosPolicyKind

	// Deadline Qos, implemented in the library.
	Deadline policy.DeadlineQosPolicy

	// Latency Budget Qos, NOT implemented in the library.
	LatencyBudget policy.LatencyBudgetQosPolicy

	// Liveliness Qos, implemented in the library.
	Liveliness policy.LivelinessQosPolicyKind

	// ReliabilityQos, implemented in the library.
	Reliability policy.ReliabilityQosPolicyKind

	// Ownership Qos, NOT implemented in the library.
	Ownership policy.OwnershipQosPolicyKind

	// Destinatio Order Qos, NOT implemented in the library.
	DestinationOrder policy.DestinationOrderQosPolicyKind

	// UserData Qos, NOT implemented in the library.
	UserData policy.UserDataQosPolicy

	// Time Based Filter Qos, NOT implemented in the library.
	TimeBasedFilter policy.TimeBasedFilterQosPolicy

	// Presentation Qos, NOT implemented in the library.
	Presentation policy.PresentationQosPolicy

	// Partition Qos, implemented in the library.
	Partition policy.PartitionQosPolicy

	// Topic Data Qos, NOT implemented in the library.
	TopicData policy.TopicDataQosPolicy

	// GroupData Qos, NOT implemented in the library.
	GroupData policy.GroupDataQosPolicy

	// Durability Service Qos, NOT implemented in the library.
	DurabilityService policy.DurabilityServiceQosPolicy

	// Lifespan Qos, NOT implemented in the library.
	Lifespan policy.LifespanQosPolicy

	// Data Representation Qos, implemented in the library.
	Representation policy.DataRepresentationQosPolicy

	// Type consistency enforcement Qos, NOT implemented in the library.
	TypeConsistency policy.TypeConsistencyEnforcementQosPolicy

	// Disable positive ACKs QoS
	DisablePositiveACKs policy.DisablePositiveACKsQosPolicy
}

func NewReaderQos() *ReaderQos {
	var qos = ReaderQos{
		Deadline:            policy.KDefaultDeadlineQosPolicy,
		LatencyBudget:       policy.KDefaultLatencyBudgetQosPolicy,
		UserData:            *policy.NewUserDataQosPolicy(policy.KPidUserData, []common.Octet{}),
		TimeBasedFilter:     policy.KDefaultTimeBasedFilterQosPolicy,
		Presentation:        policy.KDefaultPresentationQosPolicy,
		Partition:           policy.KDefaultPartitionQosPolicy,
		TopicData:           policy.KDefaultTopicDataQosPolicy,
		GroupData:           policy.KGroupDataQosPolicy,
		DurabilityService:   policy.KDefaultDurabilityServiceQosPolicy,
		Lifespan:            policy.KDefaultLifespanQosPolicy,
		Representation:      policy.KDefaultDataRepresentationQosPolicy,
		TypeConsistency:     policy.KDefaultTypeConsistencyEnforcementQosPolicy,
		DisablePositiveACKs: policy.KDefaultDisablePositiveACKsQosPolicy,
	}

	return &qos
}
