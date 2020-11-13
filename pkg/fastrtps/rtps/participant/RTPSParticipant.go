package participant

import (
	. "github.com/yeren0143/DDS/common"
	//. "dds/publisher/qos"
	. "github.com/yeren0143/DDS/fastrtps/participant"
	. "github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	. "github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery"
	. "github.com/yeren0143/DDS/fastrtps/rtps/flowcontrol"
	. "github.com/yeren0143/DDS/fastrtps/rtps/reader"
	. "github.com/yeren0143/DDS/fastrtps/rtps/resources"
	. "github.com/yeren0143/DDS/fastrtps/rtps/writer"
)

type RTPSParticipant struct {
	DomaninId         uint32
	Att               *RTPSParticipantAttributes
	Guid              *GUID_t
	Persistence_guid  *GUID_t // Persistence guid of the RTPSParticipant
	Event_Thr         *ResourceEvent
	Builtin_protocols *BuiltinProtocols
	//Resource_semaphore   *Semaphore // Semaphore to wait for the listen thread creation.
	Id_counter    uint32 // Id counter to correctly assign the ids to writers and readers.
	AllWriterList []*RTPSWriter
	AllReaderList []*RTPSReader
	Controllers   []*FlowController
	Listener      *RTPSParticipantListener
}

// type RTPSParticipant struct {
// 	GetGuid() GUID_t

// 	/**
// 	 * Indicate the Participant that you have discovered a new Remote Writer.
// 	 * This method can be used by the user to implements its own Static Endpoint
// 	 * Discovery Protocol
// 	 * @param pguid GUID_t of the discovered Writer.
// 	 * @param userDefinedId ID of the discovered Writer.
// 	 * @return True if correctly added.
// 	 */
// 	NewRemoteWriterDiscovered(pguid GUID_t, userDefinedId int16) bool

// 	/**
// 	 * Indicate the Participant that you have discovered a new Remote Reader.
// 	 * This method can be used by the user to implements its own Static Endpoint
// 	 * Discovery Protocol
// 	 * @param pguid GUID_t of the discovered Reader.
// 	 * @param userDefinedId ID of the discovered Reader.
// 	 * @return True if correctly added.
// 	 */
// 	NewRemoteReaderDiscovered(pguid GUID_t, userDefinedId int16) bool

// 	GetRTPSParticipantID() uint32

// 	RegisterWriter(writerRTPSWriter *RTPSWriter, topicAtt *TopicAttributes, wqos *WriterQos) bool
// }

func NewRTPSParticipant(patt *ParticipantAttributes, listen *ParticipantListener) *RTPSParticipant {
	return &RTPSParticipant{
		Att:      patt,
		Listener: listen,
	}
}
