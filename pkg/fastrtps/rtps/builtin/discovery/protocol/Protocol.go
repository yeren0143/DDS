package protocol

import (
	"sync"

	"dds/common"
	"dds/fastrtps/rtps/attributes"
	"dds/fastrtps/rtps/builtin/data"
	"dds/fastrtps/rtps/history"
	"dds/fastrtps/rtps/network"
	"dds/fastrtps/rtps/reader"
	"dds/fastrtps/rtps/resources"
	"dds/fastrtps/rtps/writer"
)

type IParticipant interface {
	GetAttributes() *attributes.RTPSParticipantAttributes
	GetGuid() *common.GUIDT
	CreateReader(param *attributes.ReaderAttributes, payload history.IPayloadPool,
		hist *history.ReaderHistory, listen reader.IReaderListener,
		entityID *common.EntityIDT, isBuiltin bool, enable bool) (reader.IRTPSReader, bool)
	CreateWriter(param *attributes.WriterAttributes, payload history.IPayloadPool,
		hist *history.WriterHistory, listen writer.IWriterListener,
		entityID *common.EntityIDT, isBuiltin bool) (bool, writer.IRTPSWriter)
	NetworkFactory() *network.NetFactory
	GetEventResource() *resources.ResourceEvent
	EnableReader(areader reader.IRTPSReader) bool
	HasShmTransport() bool
}

type ReaderProxyDataInitFunc = func(*data.ReaderProxyData, bool, *data.ParticipantProxyData) bool
type IPDP interface {
	Init(participant IParticipant) bool

	/**
	 * Creates an initializes a new participant proxy from a DATA(p) raw info
	 * @param p from DATA msg deserialization
	 * @param writer_guid GUID of originating writer
	 * @return new ParticipantProxyData * or nullptr on failure
	 */
	CreateParticipantProxyData(p *data.ParticipantProxyData, writerGuid *common.GUIDT) *data.ParticipantProxyData

	// Create the SPDP Writer and Reader
	// True if correct
	CreatePDPEndpoints() bool

	// This method assigns remote endpoints to the builtin endpoints defined in this protocol.
	// It also calls the corresponding methods in EDP and WLP.
	// * @param pdata Pointer to the RTPSParticipantProxyData object.
	AssignRemoteEndpoints(pdata *data.ParticipantProxyData)

	// Override to match additional endpoints to PDP. Like EDP or WLP.
	// @param pdata Pointer to the ParticipantProxyData object.
	NotifyAboveRemoteEndpoints(pdata *data.ParticipantProxyData)

	// Remove remote endpoints from the participant discovery protocol
	// @param pdata Pointer to the ParticipantProxyData to remove
	RemoveRemoteEndpoints(pdata *data.ParticipantProxyData)

	/**
	 * This method removes a remote RTPSParticipant and all its writers and readers.
	 * @param participant_guid GUID_t of the remote RTPSParticipant.
	 * @param reason Why the participant is being removed (dropped vs removed)
	 * @return true if correct.
	 */
	RemoveRemoteParticipant(partGUID *common.GUIDT, reason uint8) bool

	// Force the sending of our local DPD to all remote RTPSParticipants and multicast Locators.
	AnnounceParticipantState(newChange bool, dispose bool, wparams *common.WriteParamsT)

	// Reset the RTPSParticipantAnnouncement (only used in tests).
	ResetParticipantAnnouncement()

	Enable() bool

	// This method returns the BuiltinAttributes of the local participant.
	BuiltinAttributes() *attributes.BuiltinAttributes

	GetRTPSParticipant() IParticipant

	GetParticipantProxies() []*data.ParticipantProxyData

	GetLocalParticipantProxyData() *data.ParticipantProxyData

	GetPDPReaderHistory() *history.ReaderHistory

	GetMutex() *sync.Mutex

	AddReaderProxyData(readerGUID, participantGUID *common.GUIDT, initializer ReaderProxyDataInitFunc) *data.ReaderProxyData
}
