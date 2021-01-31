package attributes

import (
	"github.com/yeren0143/DDS/fastrtps/utils"
	"math"
)

//RemoteLocatorsAllocationAttributes specified remote locators allocation attributes
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

//NewRemoteLocatorsAllocationAttributes create locators allocation attributes
func NewRemoteLocatorsAllocationAttributes() *RemoteLocatorsAllocationAttributes {
	return &RemoteLocatorsAllocationAttributes{
		MaxUnicastLocators:   4,
		MaxMulticastLocators: 1,
	}
}

//SendBuffersAllocationAttributes describe send buffers allocation attributes
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

//NewSendBuffersAllocationAttributes create send buffers allocation attributes
func NewSendBuffersAllocationAttributes() *SendBuffersAllocationAttributes {
	return &SendBuffersAllocationAttributes{
		PreAllocatedNum: 0,
		Dynamic:         false,
	}
}

//VariableLengthDataLimits Holds limits for variable-length data.
type VariableLengthDataLimits struct {
	//! Defines the maximum size (in octets) of properties data in the local or remote participant
	MaxProperties uint64

	//! Defines the maximum size (in octets) of user data in the local or remote participant
	MaxUserData uint64

	//! Defines the maximum size (in octets) of partitions data
	MaxPartitions uint64
}

//NewVariableLengthDataLimits create VariableLengthDataLimits with default value
func NewVariableLengthDataLimits() *VariableLengthDataLimits {
	return &VariableLengthDataLimits{
		MaxProperties: 0,
		MaxUserData:   0,
		MaxPartitions: 0,
	}
}

//RTPSParticipantAllocationAttributes Holds allocation limits affecting collections managed by a participant.
type RTPSParticipantAllocationAttributes struct {
	//! Holds limits for collections of remote locators.
	Locators *RemoteLocatorsAllocationAttributes
	//! Defines the allocation behaviour for collections dependent on the total number of participants.
	Participants *utils.ResourceLimitedContainerConfig
	//! Defines the allocation behaviour for collections dependent on the total number of readers per participant.
	Readers *utils.ResourceLimitedContainerConfig
	//! Defines the allocation behaviour for collections dependent on the total number of writers per participant.
	Writers *utils.ResourceLimitedContainerConfig
	//! Defines the allocation behaviour for the send buffer manager.
	SendBuffers *SendBuffersAllocationAttributes
	//! Holds limits for variable-length data
	DataLimits *VariableLengthDataLimits
}

func (att *RTPSParticipantAllocationAttributes) TotalReaders() *utils.ResourceLimitedContainerConfig {
	return att.totalEndpoints(att.Readers)
}

func (att *RTPSParticipantAllocationAttributes) TotalWriters() *utils.ResourceLimitedContainerConfig {
	return att.totalEndpoints(att.Writers)
}

func (att *RTPSParticipantAllocationAttributes) totalEndpoints(endpoints *utils.ResourceLimitedContainerConfig) *utils.ResourceLimitedContainerConfig {
	max := uint32(math.MaxUint32)
	initial := att.Participants.Initial * endpoints.Initial
	var maxium uint32
	if att.Participants.Maximum == max || endpoints.Maximum == max {
		maxium = max
	} else {
		maxium = att.Participants.Maximum * endpoints.Maximum
	}

	var increment uint32
	if att.Participants.Increment > endpoints.Increment {
		increment = att.Participants.Increment
	} else {
		increment = endpoints.Increment
	}

	return &utils.ResourceLimitedContainerConfig{
		Initial:   initial,
		Maximum:   maxium,
		Increment: increment,
	}
}

//NewRTPSParticipantAllocationAttributes create RTPSParticipantAllocationAttributes with default value
func NewRTPSParticipantAllocationAttributes() *RTPSParticipantAllocationAttributes {
	return &RTPSParticipantAllocationAttributes{
		Locators:     NewRemoteLocatorsAllocationAttributes(),
		Participants: utils.NewResourceLimitedContainerConfig(),
		Readers:      utils.NewResourceLimitedContainerConfig(),
		Writers:      utils.NewResourceLimitedContainerConfig(),
		SendBuffers:  NewSendBuffersAllocationAttributes(),
		DataLimits:   NewVariableLengthDataLimits(),
	}
}
