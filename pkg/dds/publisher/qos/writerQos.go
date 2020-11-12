package qos

import (
	. "core/policy"
)

type WriterQos struct {
	Durability          DurabilityQosPolicy
	DurabilityService   DurabilityServiceQosPolicy
	Deadline            DeadlineQosPolicy
	LatencyBudget       LatencyBudgetQosPolicy
	Liveliness          LivelinessQosPolicy
	Reliability         ReliabilityQosPolicy
	Lifespan            LifespanQosPolicy
	UserData            UserDataQosPolicy
	TimeBasedFilter     TimeBasedFilterQosPolicy
	Ownership           OwnershipQosPolicy
	OwnershipStrength   OwnershipStrengthQosPolicy
	DestinationOrder    DestinationOrderQosPolicy
	Presentation        PresentationQosPolicy
	Partition           PartitionQosPolicy
	TopicData           TopicDataQosPolicy
	GroupData           GroupDataQosPolicy
	PublishMode         PublishModeQosPolicy
	Representation      DataRepresentationQosPolicy
	DisablePositiveACKs DisablePositiveACKsQosPolicy
}
