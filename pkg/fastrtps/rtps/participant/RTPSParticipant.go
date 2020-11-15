package participant

import (
	. "github.com/yeren0143/DDS/common"
	// . "github.com/yeren0143/DDS/fastrtps/participant"
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

// NewParticipant create new rtps participant
func NewParticipant(domainId uint32, useProtocol bool, patt *RTPSParticipantAttributes, listen *RTPSParticipantListener) *RTPSParticipant {
	return &RTPSParticipant{
		DomaninId: domainId,
		Att:       patt,
		Listener:  listen,
	}
}
