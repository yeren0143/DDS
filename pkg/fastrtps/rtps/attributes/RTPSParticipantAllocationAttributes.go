package attributes

import (
	. "fastrtps/utils"
)

type RemoteLocatorsAllocationAttributes struct {
	/** Maximum number of unicast locators per remote entity.
	 *
	 * This attribute controls the maximum number of unicast locators to keep for
	 * each discovered remote entity (be it a participant, reader of writer). It is
	 * recommended to use the highest number of local addresses found on all the systems
	 * belonging to the same domain as this participant.
	 */
	MaxUnicastLocators uint64

	/** Maximum number of multicast locators per remote entity.
	 *
	 * This attribute controls the maximum number of multicast locators to keep for
	 * each discovered remote entity (be it a participant, reader of writer). The
	 * default value of 1 is usually enough, as it doesn't make sense to add more
	 * than one multicast locator per entity.
	 */
	MaxMulticastLocators uint64
}

func NewRemoteLocatorsAllocationAttributes() RemoteLocatorsAllocationAttributes {
	return RemoteLocatorsAllocationAttributes{
		MaxUnicastLocators:   4,
		MaxMulticastLocators: 1,
	}
}

type SendBuffersAllocationAttributes struct {
	/** Initial number of send buffers to allocate.
	 *
	 * This attribute controls the initial number of send buffers to be allocated.
	 * The default value of 0 will perform an initial guess of the number of buffers
	 * required, based on the number of threads from which a send operation could be
	 * started.
	 */
	PreAllocatedNum uint64

	/** Whether the number of send buffers is allowed to grow.
	 *
	 * This attribute controls how the buffer manager behaves when a send buffer is not
	 * available. When true, a new buffer will be created. When false, it will wait for a
	 * buffer to be returned. This is a tradeoff between latency and dynamic allocations.
	 */
	Dynamic bool
}

func NewSendBuffersAllocationAttributes() SendBuffersAllocationAttributes {
	return SendBuffersAllocationAttributes{
		PreAllocatedNum: 0,
		Dynamic:         false,
	}
}

/**
 * @brief Holds limits for variable-length data.
 */
type VariableLengthDataLimits struct {
	//! Defines the maximum size (in octets) of properties data in the local or remote participant
	MaxProperties uint64

	//! Defines the maximum size (in octets) of user data in the local or remote participant
	MaxUserData uint64

	//! Defines the maximum size (in octets) of partitions data
	MaxPartitions uint64
}

func NewVariableLengthDataLimits() VariableLengthDataLimits {
	return VariableLengthDataLimits{
		MaxProperties: 0,
		MaxUserData:   0,
		MaxPartitions: 0,
	}
}

/**
 * @brief Holds allocation limits affecting collections managed by a participant.
 */
type RTPSParticipantAllocationAttributes struct {
	//! Holds limits for collections of remote locators.
	Locators RemoteLocatorsAllocationAttributes
	//! Defines the allocation behaviour for collections dependent on the total number of participants.
	Participants ResourceLimitedContainerConfig
	//! Defines the allocation behaviour for collections dependent on the total number of readers per participant.
	Readers ResourceLimitedContainerConfig
	//! Defines the allocation behaviour for collections dependent on the total number of writers per participant.
	Writers ResourceLimitedContainerConfig
	//! Defines the allocation behaviour for the send buffer manager.
	SendBuffers SendBuffersAllocationAttributes
	//! Holds limits for variable-length data
	DataLimits VariableLengthDataLimits
}

func NewRTPSParticipantAllocationAttributes() *RTPSParticipantAllocationAttributes {
	return &RTPSParticipantAllocationAttributes{
		Locators:     RemoteLocatorsAllocationAttributes{},
		Participants: ResourceLimitedContainerConfig{},
		Readers:      ResourceLimitedContainerConfig{},
		Writers:      ResourceLimitedContainerConfig{},
		SendBuffers:  SendBuffersAllocationAttributes{},
		DataLimits:   VariableLengthDataLimits{},
	}
}
