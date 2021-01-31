package policy

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/types"
)

type QosPolicyIDT uint32

const (
	INVALID_QOS_POLICY_ID QosPolicyIDT = iota

	// Standard QosPolicies
	USERDATA_QOS_POLICY_ID            //< UserDataQosPolicy
	DURABILITY_QOS_POLICY_ID          //< DurabilityQosPolicy
	PRESENTATION_QOS_POLICY_ID        //< PresentationQosPolicy
	DEADLINE_QOS_POLICY_ID            //< DeadlineQosPolicy
	LATENCYBUDGET_QOS_POLICY_ID       //< LatencyBudgetQosPolicy
	OWNERSHIP_QOS_POLICY_ID           //< OwnershipQosPolicy
	OWNERSHIPSTRENGTH_QOS_POLICY_ID   //< OwnershipStrengthQosPolicy
	LIVELINESS_QOS_POLICY_ID          //< LivelinessQosPolicy
	TIMEBASEDFILTER_QOS_POLICY_ID     //< TimeBasedFilterQosPolicy
	PARTITION_QOS_POLICY_ID           //< PartitionQosPolicy
	RELIABILITY_QOS_POLICY_ID         //< ReliabilityQosPolicy
	DESTINATIONORDER_QOS_POLICY_ID    //< DestinationOrderQosPolicy
	HISTORY_QOS_POLICY_ID             //< HistoryQosPolicy
	RESOURCELIMITS_QOS_POLICY_ID      //< ResourceLimitsQosPolicy
	ENTITYFACTORY_QOS_POLICY_ID       //< EntityFactoryQosPolicy
	WRITERDATALIFECYCLE_QOS_POLICY_ID //< WriterDataLifecycleQosPolicy
	READERDATALIFECYCLE_QOS_POLICY_ID //< ReaderDataLifecycleQosPolicy
	TOPICDATA_QOS_POLICY_ID           //< TopicDataQosPolicy
	GROUPDATA_QOS_POLICY_ID           //< GroupDataQosPolicy
	TRANSPORTPRIORITY_QOS_POLICY_ID   //< TransportPriorityQosPolicy
	LIFESPAN_QOS_POLICY_ID            //< LifespanQosPolicy
	DURABILITYSERVICE_QOS_POLICY_ID   //< DurabilityServiceQosPolicy

	//XTypes extensions
	DATAREPRESENTATION_QOS_POLICY_ID         //< DataRepresentationQosPolicy
	TYPECONSISTENCYENFORCEMENT_QOS_POLICY_ID //< TypeConsistencyEnforcementQosPolicy

	//eProsima Extensions
	DISABLEPOSITIVEACKS_QOS_POLICY_ID       //< DisablePositiveACKsQosPolicy
	PARTICIPANTRESOURCELIMITS_QOS_POLICY_ID //< ParticipantResourceLimitsQos
	PROPERTYPOLICY_QOS_POLICY_ID            //< PropertyPolicyQos
	PUBLISHMODE_QOS_POLICY_ID               //< PublishModeQosPolicy
	READERRESOURCELIMITS_QOS_POLICY_ID      //< Reader ResourceLimitsQos
	RTPSENDPOINT_QOS_POLICY_ID              //< RTPSEndpointQos
	RTPSRELIABLEREADER_QOS_POLICY_ID        //< RTPSReliableReaderQos
	RTPSRELIABLEWRITER_QOS_POLICY_ID        //< RTPSReliableWriterQos
	TRANSPORTCONFIG_QOS_POLICY_ID           //< TransportConfigQos
	TYPECONSISTENCY_QOS_POLICY_ID           //< TipeConsistencyQos
	WIREPROTOCOLCONFIG_QOS_POLICY_ID        //< WireProtocolConfigQos
	WRITERRESOURCELIMITS_QOS_POLICY_ID      //< WriterResourceLimitsQos

	NEXT_QOS_POLICY_ID //< Keep always the last element. For internal use only
)

var (
	PolicyMask = uint32(NEXT_QOS_POLICY_ID)
)

//Class QosPolicy, base for all QoS policies defined for Writers and Readers.
type QosPolicy struct {
	HasChanged bool
	SendAlways bool
}

/**
 * @brief Controls the behavior of the entity when acting as a factory for other entities. In other words,
 * configures the side-effects of the create_* and delete_* operations.
 * @note Mutable Qos Policy
 */
type EntityFactoryQosPolicy struct {
	AutoEnable_Created_Entities bool
}

func CreateEntityFactoryQosPolicy() EntityFactoryQosPolicy {
	return EntityFactoryQosPolicy{
		AutoEnable_Created_Entities: true,
	}
}

type DurabilityQosPolicyKind uint8

const (
	/**
	 * The Service does not need to keep any samples of data-instances on behalf of any DataReader that is not
	 * known by the DataWriter at the time the instance is written. In other words the Service will only attempt
	 * to provide the data to existing subscribers
	 */
	KVolatileDurabilityQos DurabilityQosPolicyKind = iota
	/**
	 * For TRANSIENT_LOCAL, the service is only required to keep the data in the memory of the DataWriter that
	 * wrote the data and the data is not required to survive the DataWriter.
	 */
	KTransientLocalDurabilityQos
	/**
	 * For TRANSIENT, the service is only required to keep the data in memory and not in permanent storage; but
	 * the data is not tied to the lifecycle of the DataWriter and will, in general, survive it.
	 */
	KTransientDurabilityQos
	/**
	 * Data is kept on permanent storage, so that they can outlive a system session.
	 * @warning Not Supported
	 */
	KPersistentDurabilityQos
)

const (
	PARAMETER_KIND_LENGTH = 4
	PARAMETER_BOOL_LENGTH = 4
)

type DurabilityQosPolicy struct {
	ParameterT
	QosPolicy
	Kind DurabilityQosPolicyKind
}

var DefaultDurabilityQosPolicy = DurabilityQosPolicy{
	ParameterT: *NewParameterT(KPidDurability, PARAMETER_KIND_LENGTH),
	QosPolicy: QosPolicy{
		HasChanged: false,
		SendAlways: true,
	},
	Kind: KVolatileDurabilityQos,
}

/**
 * @brief DataReader expects a new sample updating the value of each instance at least once every deadline period.
 * DataWriter indicates that the application commits to write a new value (using the DataWriter) for each instance managed
 * by the DataWriter at least once every deadline period.
 * @note Mutable Qos Policy
 */
type DeadlineQosPolicy struct {
	ParameterT
	QosPolicy
	Period common.DurationT
}

var KDefaultDeadlineQosPolicy = DeadlineQosPolicy{
	ParameterT: *NewParameterT(KPidDeadline, KParameterTimeLength),
	QosPolicy: QosPolicy{
		HasChanged: false,
		SendAlways: true,
	},
	Period: common.Time{
		Seconds: common.KInfiniteSeconds,
		Nanosec: common.KInfiniteNanoSeconds,
	},
}

/**
 * Specifies the maximum acceptable delay from the time the data is written until the data is inserted in the receiver's
 * application-cache and the receiving application is notified of the fact.This policy is a hint to the Service, not something
 * that must be monitored or enforced. The Service is not required to track or alert the user of any violation.
 * @warning This QosPolicy can be defined and is transmitted to the rest of the network but is not implemented in this version.
 * @note Mutable Qos Policy
 */
type LatencyBudgetQosPolicy struct {
	ParameterT
	QosPolicy
	Duration common.DurationT
}

var KDefaultLatencyBudgetQosPolicy = LatencyBudgetQosPolicy{
	ParameterT: *NewParameterT(KPidLatencyBudget, KParameterTimeLength),
	QosPolicy: QosPolicy{
		HasChanged: false,
		SendAlways: true,
	},
	Duration: common.Time{
		Seconds: 0,
		Nanosec: 0,
	},
}

type LivelinessQosPolicyKind uint8
const (
	// The infrastructure will automatically signal liveliness
	// for the DataWriters at least as often as required by the lease_duration.
	AUTOMATIC_LIVELINESS_QOS LivelinessQosPolicyKind = iota

	// The Service will assume that as long as at least one Entity within
	// the DomainParticipant has asserted its liveliness the other
	// Entities in that same DomainParticipant are also alive.
	MANUAL_BY_PARTICIPANT_LIVELINESS_QOS

	//The Service will only assume liveliness of the DataWriter
	//if the application has asserted liveliness of that DataWriter itself.
	MANUAL_BY_TOPIC_LIVELINESS_QOS
)

/**
 * Determines the mechanism and parameters used by the application to determine whether an Entity is “active” (alive).
 * The “liveliness” status of an Entity is used to maintain instance ownership in combination with the setting of the
 * OwnershipQosPolicy.
 * The application is also informed via listener when an Entity is no longer alive.
 *
 * The DataReader requests that liveliness of the writers is maintained by the requested means and loss of liveliness is
 * detected with delay not to exceed the lease_duration.
 *
 * The DataWriter commits to signaling its liveliness using the stated means at intervals not to exceed the lease_duration.
 * Listeners are used to notify the DataReaderof loss of liveliness and DataWriter of violations to the liveliness contract.
 */
type LivelinessQosPolicy struct {
	ParameterT
	QosPolicy
	Kind               LivelinessQosPolicyKind
	LeaseDuration      common.DurationT
	AnnouncementPeriod common.DurationT
}

var KDefaultLivelinessQosPolicy = LivelinessQosPolicy{
	ParameterT: *NewParameterT(KPidLiveliness, KParameterTimeLength),
	QosPolicy: QosPolicy{
		HasChanged: false,
		SendAlways: true,
	},
	Kind: AUTOMATIC_LIVELINESS_QOS,
	LeaseDuration: common.Time{
		Seconds: common.KInfiniteSeconds,
		Nanosec: common.KInfiniteNanoSeconds,
	},
	AnnouncementPeriod: common.Time{
		Seconds: common.KInfiniteSeconds,
		Nanosec: common.KInfiniteNanoSeconds,
	},
}

type ReliabilityQosPolicyKind uint8

const (
	/**
	 * Indicates that it is acceptable to not retry propagation of any samples. Presumably new values for the samples
	 * are generated often enough that it is not necessary to re-send or acknowledge any samples
	 */
	BEST_EFFORT_RELIABILITY_QOS ReliabilityQosPolicyKind = 0x01

	/**
	 * Specifies the Service will attempt to deliver all samples in its history. Missed samples may be retried.
	 * In steady-state (no modifications communicated via the DataWriter) the middleware guarantees that all samples
	 * in the DataWriter history will eventually be delivered to all the DataReader objects. Outside steady state the
	 * HistoryQosPolicy and ResourceLimitsQosPolicy will determine how samples become part of the history and whether
	 * samples can be discarded from it.
	 */
	RELIABLE_RELIABILITY_QOS = 0x02
)

type ReliabilityQosPolicy struct {
	ParameterT
	QosPolicy
	Kind            ReliabilityQosPolicyKind
	MaxBlockingTime common.DurationT
}

var KDefaultReliabilityQosPolicy = ReliabilityQosPolicy{
	ParameterT: *NewParameterT(KPidReliability, PARAMETER_KIND_LENGTH+KParameterTimeLength),
	QosPolicy: QosPolicy{
		HasChanged: false,
		SendAlways: true,
	},
	Kind: BEST_EFFORT_RELIABILITY_QOS,
	MaxBlockingTime: common.Time{
		Seconds: 0,
		Nanosec: 100000000,
	}, // 100ms
}

type OwnershipQosPolicyKind uint8

const (
	/**
	 * Indicates shared ownership for each instance. Multiple writers are allowed to
	 * update the same instance and all the updates are made available to the readers.
	 * In other words there is no concept of an “owner” for the instances.
	 */
	SHARED_OWNERSHIP_QOS OwnershipQosPolicyKind = iota

	/**
	 * Indicates each instance can only be owned by one DataWriter,
	 * but the owner of an instance can change dynamically.
	 * The selection of the owner is controlled by the setting of the OwnershipStrengthQosPolicy.
	 * The owner is always set to be the highest-strength DataWriter
	 * object among the ones currently “active” (as determined by the LivelinessQosPolicy).
	 */
	EXCLUSIVE_OWNERSHIP_QOS
)

type OwnershipQosPolicy struct {
	ParameterT
	QosPolicy
	Kind OwnershipQosPolicyKind
}

var KDefaultOwnershipQosPolicy = OwnershipQosPolicy{
	ParameterT: *NewParameterT(KPidOwnership, PARAMETER_KIND_LENGTH),
	QosPolicy: QosPolicy{
		HasChanged: false,
		SendAlways: true,
	},
	Kind: SHARED_OWNERSHIP_QOS,
}

type DestinationOrderQosPolicyKind uint8

const (
	/**
	 * Indicates that data is ordered based on the reception time at each Subscriber.
	 * Since each subscriber may receive
	 * the data at different times there is no guaranteed that the changes will be
	 * seen in the same order. Consequently,
	 * it is possible for each subscriber to end up with a different final value for the data.
	 */
	BY_RECEPTION_TIMESTAMP_DESTINATIONORDER_QOS DestinationOrderQosPolicyKind = iota

	/**
	 * Indicates that data is ordered based on a timestamp placed at the source
	 * (by the Service or by the application).
	 * In any case this guarantees a consistent final value for the data in all subscribers.
	 */
	BY_SOURCE_TIMESTAMP_DESTINATIONORDER_QOS
)

type DestinationOrderQosPolicy struct {
	ParameterT
	QosPolicy
	Kind DestinationOrderQosPolicyKind
}

var KDefaultDestinationOrderQosPolicy = DestinationOrderQosPolicy{
	ParameterT: *NewParameterT(KPidDestinationOrder, PARAMETER_KIND_LENGTH),
	QosPolicy: QosPolicy{
		HasChanged: false,
		SendAlways: true,
	},
	Kind: BY_RECEPTION_TIMESTAMP_DESTINATIONORDER_QOS,
}

// GenericDataQosPolicy, base class to transmit user data during the discovery phase.
type GenericDataQosPolicy struct {
	ParameterT
	QosPolicy
}

func NewGenericDataQosPolicy(pid ParameterIDT, length uint32) *GenericDataQosPolicy {
	return &GenericDataQosPolicy{
		ParameterT: *NewParameterT(pid, length),
		QosPolicy: QosPolicy{
			HasChanged: false,
			SendAlways: false,
		},
	}
}

type UserDataQosPolicy struct {
	GenericDataQosPolicy
}

func NewUserDataQosPolicy() *UserDataQosPolicy {
	return &UserDataQosPolicy{}
}

type TopicDataQosPolicy struct {
	GenericDataQosPolicy
}

var KDefaultTopicDataQosPolicy = TopicDataQosPolicy{
	GenericDataQosPolicy: *NewGenericDataQosPolicy(KPidTopicData, 0),
}

type GroupDataQosPolicy struct {
	GenericDataQosPolicy
}

var KGroupDataQosPolicy = GroupDataQosPolicy{
	GenericDataQosPolicy: *NewGenericDataQosPolicy(KPidGroupData, 0),
}

/**
 * Filter that allows a DataReader to specify that it is interested only in (potentially) a subset of
 * the values of the data.
 * The filter states that the DataReader does not want to receive more than one value each minimum_separation,
 * regardless
 * of how fast the changes occur. It is inconsistent for a DataReader to have a minimum_separation longer
 * than its Deadline period.
 * @warning This QosPolicy can be defined and is transmitted to the rest of the network but is not implemented in this version.
 * @note Mutable Qos Policy
 */
type TimeBasedFilterQosPolicy struct {
	ParameterT
	QosPolicy
	Minium_Separation common.DurationT
}

var KDefaultTimeBasedFilterQosPolicy = TimeBasedFilterQosPolicy{
	ParameterT: *NewParameterT(KPidTimeBasedFilter, KParameterTimeLength),
	QosPolicy: QosPolicy{
		HasChanged: false,
		SendAlways: false,
	},
	Minium_Separation: common.Time{
		Seconds: 0,
		Nanosec: 0,
	},
}

type PresentationQosPolicyAccessScopeKind uint8

const (
	/**
	 * Scope spans only a single instance. Indicates that changes to one instance need not be coherent nor
	 * ordered with respect to changes to any other instance. In other words,
	 * order and coherent changes apply to each instanceseparately.
	 */
	INSTANCE_PRESENTATION_QOS PresentationQosPolicyAccessScopeKind = iota

	/**
	 * Scope spans to all instances within the same DataWriter (or DataReader), but not across instances
	 * in different DataWriter (or DataReader).
	 */
	TOPIC_PRESENTATION_QOS

	/**
	 * Scope spans to all instances belonging to DataWriter (or DataReader) entities within
	 * the same Publisher (or Subscriber).
	 */
	GROUP_PRESENTATION_QOS
)

const (
	PARAMETER_PRESENTATION_LENGTH = 8
)

/**
 * Specifies how the samples representing changes to data instances are presented to the subscribing application.
 * This policy affects the application’s ability to specify and receive coherent changes and to see the relative
 * order of changes.access_scope determines the largest scope spanning the entities for which the
 * order and coherency of changes can be preserved. The two booleans control whether coherent access
 * and ordered access are supported within the scope access_scope.
 * @warning This QosPolicy can be defined and is transmitted to the rest of the network but is not
 * implemented in this version.
 * @note Immutable Qos Policy
 */
type PresentationQosPolicy struct {
	ParameterT
	QosPolicy
	Access_Scope PresentationQosPolicyAccessScopeKind

	/**
	 * @brief Specifies support coherent access. That is, the ability to group a set of changes as a unit
	 * on the publishing end such that they are received as a unit at the subscribing end.
	 * by default, false.
	 */
	CoherentAccess bool

	/**
	 * @brief Specifies support for ordered access to the samples received at the subscription end. That is,
	 * the ability of the subscriber to see changes in the same order as they occurred on the publishing end.
	 * By default, false.
	 */
	OrderedAccess bool
}

var KDefaultPresentationQosPolicy = PresentationQosPolicy{
	ParameterT: *NewParameterT(KPidPresentation, PARAMETER_KIND_LENGTH),
	QosPolicy: QosPolicy{
		HasChanged: false,
		SendAlways: false,
	},
	Access_Scope:   INSTANCE_PRESENTATION_QOS,
	CoherentAccess: false,
	OrderedAccess:  false,
}

type Partition_t struct {
	Partition string
}

/**
 * Set of strings that introduces a logical partition among the topics visible by the Publisher and Subscriber.
 * A DataWriter within a Publisher only communicates with a DataReader in a Subscriber if (in addition to matching the
 * Topic and having compatible QoS) the Publisher and Subscriber have a common partition name string.
 *
 * The empty string ("") is considered a valid partition that is matched with other partition names using the same rules of
 * string matching and regular-expression matching used for any other partition name.
 * @note Mutable Qos Policy
 */
type PartitionQosPolicy struct {
	ParameterT
	QosPolicy
	MaxSize    uint32
	Partitions common.SerializedPayloadT
	NPartions  uint32 // Number of partitions. <br> By default, 0.
}

var KDefaultPartitionQosPolicy = PartitionQosPolicy{
	ParameterT: *NewParameterT(KPidPartition, 0),
	QosPolicy: QosPolicy{
		HasChanged: false,
		SendAlways: false,
	},
	MaxSize:   0,
	NPartions: 0,
}

type HistoryQosPolicyKind uint8

const (
	/**
	 * On the publishing side, the Service will only attempt to keep the most recent “depth” samples
	 * of each instance of data (identified by its key) managed by the DataWriter.
	 * On the subscribing side, the DataReader will only attempt to keep the most recent “depth” samples
	 * received for each instance (identified by its key) until the application “takes” them
	 * via the DataReader’s take operation.
	 */
	KEEP_LAST_HISTORY_QOS HistoryQosPolicyKind = iota

	/**
	 * On the publishing side, the Service will attempt to keep all samples (representing each value written)
	 * of each instance of data (identified by its key) managed by the DataWriter until they can be delivered
	 * to all subscribers. On the subscribing side, the Service will attempt to keep all samples of each
	 * instance of data (identified by its key) managed by the DataReader.
	 * These samples are kept until the application “takes” them from the Service via the take operation.
	 */
	KEEP_ALL_HISTORY_QOS
)

/**
 * Specifies the behavior of the Service in the case where the value of a sample changes
 * (one or more times) before it can be successfully communicated to one or more existing subscribers.
 * This QoS policy controls whether the Service should deliver only the most recent value,
 * attempt to deliver all intermediate values, or do something in between.
 * On the publishing side this policy controls the samples that should be maintained by the DataWriter
 * on behalf of existing DataReader entities. The behavior with regards to a DataReaderentities discovered
 * after a sample is written is controlled by the DURABILITY QoS policy.
 * On the subscribing side it controls the samples that should be maintained until the application
 * “takes” them from the Service.
 * @note Immutable Qos Policy
 */
type HistoryQosPolicy struct {
	ParameterT
	QosPolicy
	Kind  HistoryQosPolicyKind
	Depth int32
}

//Specifies the resources that the Service can consume in order to meet the requested QoS
type ResourceLimitsQosPolicy struct {
	ParameterT
	QosPolicy

	/**
	 * @brief Specifies the maximum number of data-samples the DataWriter (or DataReader) can manage across all the
	 * instances associated with it. Represents the maximum samples the middleware can store for any one DataWriter
	 * (or DataReader). <br>
	 * By default, 5000.
	 * @warning It is inconsistent for this value to be less than max_samples_per_instance.
	 */
	MaxSamples int32

	/**
	 * @brief Represents the maximum number of instances DataWriter (or DataReader) can manage. <br>
	 * By default, 10.
	 */
	MaxInstances int32

	/**
	 * @brief Represents the maximum number of samples of any one instance a DataWriter(or DataReader) can manage. <br>
	 * By default, 400.
	 * @warning It is inconsistent for this value to be greater than max_samples.
	 */
	MaxSamples_PerInstance int32

	/**
	 * @brief Number of samples currently allocated. <br>
	 * By default, 100.
	 */
	AllocatedSamples int32
}

/**
 * Specifies the configuration of the durability service. That is, the service that implements the
 * DurabilityQosPolicy kind of TRANSIENT and PERSISTENT.
 * @warning This QosPolicy can be defined and is transmitted to the rest of the network but is not
 * implemented in this version.
 * @note Immutable Qos Policy
 */
type DurabilityServiceQosPolicy struct {
	ParameterT
	QosPolicy

	/**
	 * @brief Control when the service is able to remove all information regarding a data-instance. <br>
	 * By default, c_TimeZero.
	 */
	Service_CleanUp_Delay common.DurationT

	//Controls the HistoryQosPolicy of the fictitious DataReader that stores the data
	// within the durability service.
	HistoryKind HistoryQosPolicyKind

	//Number of most recent values that should be maintained on the History.
	// It only have effect if the history_kind
	HistoryDepth int32

	/* Control the ResourceLimitsQos of the implied DataReader that stores the data within the durability service.
	 * Specifies the maximum number of data-samples the DataWriter (or DataReader) can manage across
	 * all the instances associated with it. Represents the maximum samples the middleware can store for
	 * any one DataWriter (or DataReader). It is inconsistent for this value to be less than max_samples_per_instance.
	* By default, -1 (Length Unlimited).
	*/
	MaxSamples int32

	/**
	 * @brief Control the ResourceLimitsQos of the implied DataReader that stores the data within the durability service.
	 * Represents the maximum number of instances DataWriter (or DataReader) can manage. <br>
	 * By default, -1 (Length Unlimited).
	 */
	MaxInstances int32

	/**
	 * @brief Control the ResourceLimitsQos of the implied DataReader that stores the data within the durability service.
	 * Represents the maximum number of samples of any one instance a DataWriter(or DataReader) can manage.
	 * It is inconsistent for this value to be greater than max_samples. <br>
	 * By default, -1 (Length Unlimited).
	 */
	MaxSamplesPerInstance int32
}

var KDefaultDurabilityServiceQosPolicy = DurabilityServiceQosPolicy{
	ParameterT: *NewParameterT(KPidDurabilityService, KParameterTimeLength+PARAMETER_KIND_LENGTH+4+4+4+4),
	QosPolicy: QosPolicy{
		HasChanged: false,
		SendAlways: false,
	},
	HistoryDepth:          -1,
	MaxSamples:            -1,
	MaxInstances:          -1,
	MaxSamplesPerInstance: -1,
}

//Specifies the maximum duration of validity of the data written by the DataWriter.
type LifespanQosPolicy struct {
	ParameterT
	QosPolicy
	Duration common.DurationT
}

var KDefaultLifespanQosPolicy = LifespanQosPolicy{
	ParameterT: *NewParameterT(KPidLifeSpan, KParameterTimeLength),
	QosPolicy: QosPolicy{
		HasChanged: false,
		SendAlways: true,
	},
	Duration: common.Time{
		Seconds: common.KInfiniteSeconds,
		Nanosec: common.KInfiniteNanoSeconds,
	},
}

/**
 * Specifies the value of the “strength” used to arbitrate among multiple DataWriter objects that attempt to modify the same
 * instance of a data-object (identified by Topic + key).This policy only applies if the OWNERSHIP QoS policy is of kind
 * EXCLUSIVE.
 * @note Mutable Qos Policy
 */
type OwnershipStrengthQosPolicy struct {
	ParameterT
	QosPolicy
	Value uint32
}

/**
 * This policy is a hint to the infrastructure as to how to set the priority of the underlying transport used to send the data.
 * @warning This QosPolicy can be defined and is transmitted to the rest of the network but is not implemented in this version.
 * @note Mutable Qos Policy
 */
type TransportPriorityQosPolicy struct {
	ParameterT
	QosPolicy
	Value uint32
}

type PublishModeQosPolicyKind uint8

const (
	KSynPublishMode PublishModeQosPolicyKind = iota
	KAsynPublishMode
)

type PublishModeQosPolicy struct {
	QosPolicy
	Kind PublishModeQosPolicyKind
}

type DataRepresentationId_t int16

const (
	XCDR_DATA_REPRESENTATION  DataRepresentationId_t = 0
	XML_DATA_REPRESENTATION   DataRepresentationId_t = 1 //TODO
	XCDR2_DATA_REPRESENTATION DataRepresentationId_t = 2
)

/**
 * With multiple standard data Representations available, and vendor-specific extensions possible, DataWriters and
 * DataReaders must be able to negotiate which data representation(s) to use. This negotiation shall occur based on
 * DataRepresentationQosPolicy.
 * @warning If a writer’s offered representation is contained within a reader’s sequence, the offer satisfies the
 * request and the policies are compatible. Otherwise, they are incompatible.
 * @note Immutable Qos Policy
 */
type DataRepresentationQosPolicy struct {
	ParameterT
	QosPolicy
	Values []DataRepresentationId_t
}

var KDefaultDataRepresentationQosPolicy = DataRepresentationQosPolicy{
	ParameterT: *NewParameterT(KPidDataRepresentation, 0),
	QosPolicy: QosPolicy{
		HasChanged: false,
		SendAlways: true,
	},
}

type TypeConsistencyKind uint16

const (
	/**
	 * The DataWriter and the DataReader must support the same data type in order for them to communicate.
	 */
	KDisallowTypeCoercionTypeConsistency = iota
	/**
	 * The DataWriter and the DataReader need not support the same data type in order for them to communicate as long as
	 * the reader’s type is assignable from the writer’s type.
	 */
	KAllowTypeCoercion
)

/**
 * The TypeConsistencyEnforcementQosPolicy defines the rules for determining whether the type used to
 * publish a given data stream is consistent with that used to subscribe to it. It applies to DataReaders.
 * @note Immutable Qos Policy
 */
type TypeConsistencyEnforcementQosPolicy struct {
	Kind TypeConsistencyKind

	/**
	 * @brief This option controls whether sequence bounds are taken into consideration for type assignability. If the
	 * option is set to TRUE, sequence bounds (maximum lengths) are not considered as part of the type assignability.
	 * This means that a T2 sequence type with maximum length L2 would be assignable to a T1 sequence type with maximum
	 * length L1, even if L2 is greater than L1. If the option is set to false, then sequence bounds are taken into
	 * consideration for type assignability and in order for T1 to be assignable from T2 it is required that L1>= L2. <br>
	 * By default, true.
	 */
	IgnoreSequenceBounds bool

	/**
	 * @brief This option controls whether string bounds are taken into consideration for type assignability. If the option
	 *  is set to TRUE, string bounds (maximum lengths) are not considered as part of the type assignability. This means
	 * that a T2 string type with maximum length L2 would be assignable to a T1 string type with maximum length L1, even
	 * if L2 is greater than L1. If the option is set to false, then string bounds are taken into consideration for type
	 * assignability and in order for T1 to be assignable from T2 it is required that L1>= L2. <br>
	 * By default, true.
	 */
	IgnoreStringBounds bool

	/**
	 * @brief This option controls whether member names are taken into consideration for type assignability.
	 * If the option is set to TRUE, member names are considered as part of assignability in addition to member
	 * IDs (so that members with the same ID also have the same name). If the option is set to FALSE,
	 * then member names are not ignored.
	 * By default, false.
	 */
	IgnoreMemberNames bool

	/**
	 * @brief This option controls whether type widening is allowed. If the option is set to FALSE, type widening is
	 * permitted. If the option is set to TRUE,it shall cause a wider type to not be assignable to a narrower type. <br>
	 * By default, false.
	 */
	PreventTypeWidening bool

	/**
	 * @brief This option requires type information to be available in order to complete matching between
	 * a DataWriter and DataReader when set to TRUE, otherwise matching can occur without complete type
	 * information when set to FALSE.
	 * By default, false.
	 */
	ForceTypeValidation bool
}

var KDefaultTypeConsistencyEnforcementQosPolicy = TypeConsistencyEnforcementQosPolicy{
	Kind:                 KAllowTypeCoercion,
	IgnoreSequenceBounds: true,
	IgnoreStringBounds:   true,
	IgnoreMemberNames:    false,
	PreventTypeWidening:  false,
	ForceTypeValidation:  false,
}

type DisablePositiveACKsQosPolicy struct {
	ParameterT
	QosPolicy

	//True if this QoS is enabled.
	Enable bool
	//The duration to keep samples for (not serialized as not needed by reader).
	Duration common.DurationT
}

var KDefaultDisablePositiveACKsQosPolicy = DisablePositiveACKsQosPolicy{
	ParameterT: *NewParameterT(KPidDisablePositiveAcks, PARAMETER_BOOL_LENGTH),
	QosPolicy: QosPolicy{
		HasChanged: false,
		SendAlways: true,
	},
	Enable: false,
	Duration: common.Time{
		Seconds: common.KInfiniteSeconds,
		Nanosec: common.KInfiniteNanoSeconds,
	},
}

type TypeIDV1 struct {
	ParameterT
	QosPolicy

	TypeIdentifier types.TypeIdentifier_t
}

type TypeObjectV1 struct {
	ParameterT
	QosPolicy
	Type_Object types.TypeObject_t
}

type TypeInformation struct {
	ParameterT
	QosPolicy
	//Type_Info TypeInform_t
	Assigned bool
}

type WireProtocolConfigQos struct {
	QosPolicy
	Prefix                      common.GUIDPrefixT
	ParticipantID               int32
	Port                        *common.PortParameters
	DefaultUnicastLocatorList   *common.LocatorList
	DefaultMulticastLocatorList *common.LocatorList
	//Builtin_Attr                   *BuiltinAttributes
	//Throughput_Controller          *ThroughputControllerDescriptor
}

type TransportConfigQos struct {
	QosPolicy
	//User_Transports           []*TransportDescriptorInterface
	UseBuiltinTransport    bool
	SendSocketBufferSize   uint32
	ListenSocketBufferSize uint32
}

type RTPSEndpointQos struct {
}

type WriterResourceLimitsQos struct {
}
