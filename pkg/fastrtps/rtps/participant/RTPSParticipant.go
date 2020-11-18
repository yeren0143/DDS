package participant

import (
	common "github.com/yeren0143/DDS/common"
	attributes "github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	discovery "github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery"
	flowctl "github.com/yeren0143/DDS/fastrtps/rtps/flowcontrol"
	reader "github.com/yeren0143/DDS/fastrtps/rtps/reader"
	resources "github.com/yeren0143/DDS/fastrtps/rtps/resources"
	writer "github.com/yeren0143/DDS/fastrtps/rtps/writer"
)

//RTPSParticipant allows the creation and removal of writers and readers. It manages the send and receive threads.
type RTPSParticipant struct {
	DomaninID        uint32
	Att              *attributes.RTPSParticipantAttributes
	GUID             common.GUIDT
	PersistenceGUID  *common.GUIDT // Persistence guid of the RTPSParticipant
	EventThr         *resources.ResourceEvent
	BuiltinProtocols *discovery.BuiltinProtocols //Resource_semaphore   *Semaphore // Semaphore to wait for the listen thread creation.
	IDCounter        uint32                      // Id counter to correctly assign the ids to writers and readers.
	AllWriterList    []*writer.RTPSWriter
	AllReaderList    []*reader.RTPSReader
	Controllers      []*flowctl.FlowController
	Listener         *RTPSParticipantListener
	hasShmTransport  bool
}

// NewParticipant create new rtps participant
func NewParticipant(domainID uint32, pparam *attributes.RTPSParticipantAttributes, guidP, perstGUID *common.GUIDPrefixT, listen *RTPSParticipantListener) *RTPSParticipant {
	return &RTPSParticipant{
		DomaninID:       domainID,
		Att:             pparam,
		GUID:            common.GUIDT{Prefix: *guidP, EntID: *common.CEidRTPSParticipant},
		Listener:        listen,
		hasShmTransport: false,
	}
}
