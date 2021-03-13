package protocol

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/data"
	"github.com/yeren0143/DDS/fastrtps/rtps/history"
	"github.com/yeren0143/DDS/fastrtps/rtps/network"
	"github.com/yeren0143/DDS/fastrtps/rtps/reader"
	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
	"github.com/yeren0143/DDS/fastrtps/rtps/writer"
)

type IParticipant interface {
	GetAttributes() *attributes.RTPSParticipantAttributes
	GetGuid() *common.GUIDT
	CreateReader(param *attributes.ReaderAttributes, payload history.IPayloadPool,
		hist *history.ReaderHistory, listen reader.IReaderListener,
		entityID *common.EntityIDT, isBuiltin bool, enable bool) (bool, reader.IRTPSReader)
	CreateWriter(param *attributes.WriterAttributes, payload history.IPayloadPool,
		hist *history.WriterHistory, listen writer.IWriterListener,
		entityID *common.EntityIDT, isBuiltin bool) (bool, writer.IRTPSWriter)
	NetworkFactory() *network.NetFactory
	GetEventResource() *resources.ResourceEvent
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

	// Force the sending of our local DPD to all remote RTPSParticipants and multicast Locators.
	AnnounceParticipantState(newChange bool, dispose bool, wparams *common.WriteParamsT)

	// This method returns the BuiltinAttributes of the local participant.
	BuiltinAttributes() *attributes.BuiltinAttributes

	GetRTPSParticipant() IParticipant

	GetLocalParticipantProxyData() *data.ParticipantProxyData

	AddReaderProxyData(readerGUID, participantGUID *common.GUIDT, initializer ReaderProxyDataInitFunc) *data.ReaderProxyData
}
