package qos

import (
	. "core/policy"
)

type WriterQos struct {
	Durability          DurabilityQosPolicy
	DurabilityService   DurabilityServiceQosPolicy
	Deadline            DeadLineQosPolicy
	LatencyBudget       LatencyBudgetQosPolicy
	Liveliness          LiveLinessQosPolicy
	Reliability         ReliabilityQosPolicy
	LifeSpan            LifeSpanQosPolicy
	UserData            UserDataQosPolicy
	TimeBasedFilter     TimeBasedFilterQosPolicy
	OwnerShip           OwnerShipQosPolicy
	OwnerShipStrength   OwnerShipStrengthQosPolicy
	DestinationOrder    DestinationOrderQosPolicy
	Presentation        PresentationQosPolicy
	Partition           PartitionQosPolicy
	TopicData           TopicDataQosPolicy
	GroupData           GroupDataQosPolicy
	PublishMode         PublishModeQosPolicy
	Representation      DataRepresentationQosPolicy
	DisablePositiveACKs DisablePositiveACKsQosPolicy
}
