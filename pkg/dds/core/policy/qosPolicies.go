package policy

type QosPolicyId uint32

const (
	INVALID_QOS_POLICY_ID QosPolicyId = iota

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
type QosPolicy interface {
	SendAlways() bool
	Clear()
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
	VOLATILE_DURABILITY_QOS DurabilityQosPolicyKind = iota
	/**
	 * For TRANSIENT_LOCAL, the service is only required to keep the data in the memory of the DataWriter that
	 * wrote the data and the data is not required to survive the DataWriter.
	 */
	TRANSIENT_LOCAL_DURABILITY_QOS
	/**
	 * For TRANSIENT, the service is only required to keep the data in memory and not in permanent storage; but
	 * the data is not tied to the lifecycle of the DataWriter and will, in general, survive it.
	 */
	TRANSIENT_DURABILITY_QOS
	/**
	 * Data is kept on permanent storage, so that they can outlive a system session.
	 * @warning Not Supported
	 */
	PERSISTENT_DURABILITY_QOS
)
