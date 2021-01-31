package participant

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/data"
	"github.com/yeren0143/DDS/fastrtps/rtps/history"
	"github.com/yeren0143/DDS/fastrtps/rtps/reader"
	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
	"github.com/yeren0143/DDS/fastrtps/rtps/writer"
	"log"
	"sync"
)

type IRTPSParticipant interface {
	GetAttributes() *attributes.RTPSParticipantAttributes
	GetGuid() *common.GUIDT
}

type IPDP interface {
	Init(participant IRTPSParticipant) bool

	/**
	 * Creates an initializes a new participant proxy from a DATA(p) raw info
	 * @param p from DATA msg deserialization
	 * @param writer_guid GUID of originating writer
	 * @return new ParticipantProxyData * or nullptr on failure
	 */
	CreateParticipantProxyData(p *data.ParticipantProxyData, writer_guid *common.GUIDT) *data.ParticipantProxyData

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
}

type IBuiltinProtocols interface {
	UpdateMetatrafficLocators(loclist *common.LocatorList) bool
	GetBuiltinAttributes() *attributes.BuiltinAttributes
}

type IpdpBaseImpl interface {
	CreatePDPEndpoints() bool
}

type pdpBase struct {
	builtin     IBuiltinProtocols
	discovery   *attributes.BuiltinAttributes
	participant IRTPSParticipant
	writer      writer.IRTPSWriter
	reader      reader.IRTPSReader
	// Number of participant proxy data objects created
	participantProxiesNumber int
	// Registered RTPSParticipants (including the local one, that is the first one.)
	participantProxies []*data.ParticipantProxyData
	// Pool of participant proxy data objects ready for reuse
	participantProxiesPool []*data.ParticipantProxyData
	// Pool of reader proxy data objects ready for reuse
	readerProxiesPool []*data.ReaderProxyData
	// Pool of writer proxy data objects ready for reuse
	writerProxiesPool []*data.WriterProxyData
	// Variable to indicate if any parameter has changed.
	hasChangedLocalPDP bool
	// Listener for the SPDP messages.
	listener          *reader.ReaderListener
	pdpWriterHistory  *history.WriterHistory
	pdpReaderHistory  *history.ReaderHistory
	writePayloadPool  history.ITopicPayloadPool
	readerPayloadPool history.ITopicPayloadPool
	// ReaderProxyData to allow preallocation of remote locators
	tempReaderData *data.ReaderProxyData
	// WriterProxyData to allow preallocation of remote locators
	tempWriterData *data.WriterProxyData
	tempDataMutex  sync.Mutex
	// Participant data atomic access assurance
	mutex sync.Mutex
	// To protect callbacks (ParticipantProxyData&)
	callbackMutex sync.Mutex

	// TimedEvent to periodically resend the local RTPSParticipant information.
	resendParticipantInfoEvent *resources.TimedEvent
	initialAnnouncements       attributes.InitialAnnouncementConfig
	impl                       IpdpBaseImpl
}

func (pdp *pdpBase) initPDP(participant IRTPSParticipant) bool {
	log.Println("Beginning")
	pdp.participant = participant
	pdp.discovery = participant.GetAttributes().Builtin
	pdp.initialAnnouncements = *pdp.discovery.DiscoveryConfig.InitialAnnouncements
	//CREATE ENDPOINTS
	if !pdp.impl.CreatePDPEndpoints() {
		return false
	}
	//UPDATE METATRAFFIC.
	pdp.builtin.UpdateMetatrafficLocators(&pdp.reader.GetAttributes().UnicastLocatorList)

	pdp.mutex.Lock()
	pdata := pdp.addParticipantProxyData(participant.GetGuid(), false)
	pdp.mutex.Unlock()

	if pdata == nil {
		return false
	}
	pdp.initializeParticipantProxyData(pdata)

	log.Fatalln("not impl")
	return false
}

func (pdp *pdpBase) addParticipantProxyData(guid *common.GUIDT, withLeaseDuration bool) *data.ParticipantProxyData {
	var retVal *data.ParticipantProxyData
	log.Fatalln("addParticipantProxyData not impl")

	return retVal
}

func (pdp *pdpBase) initializeParticipantProxyData(*data.ParticipantProxyData) {

}

func newPDP(protocol IBuiltinProtocols, att *attributes.RTPSParticipantAllocationAttributes, impl IpdpBaseImpl) *pdpBase {
	var pdp pdpBase
	pdp.builtin = protocol
	// pdp.participantProxies = att.Participants
	pdp.hasChangedLocalPDP = true

	maxUnicastLocators := uint32(att.Locators.MaxUnicastLocators)
	maxMulticastLocators := uint32(att.Locators.MaxMulticastLocators)

	pdp.tempReaderData = data.NewReaderProxyData(maxUnicastLocators, maxMulticastLocators, att.DataLimits)
	pdp.tempWriterData = data.NewWriterProxyData(maxUnicastLocators, maxMulticastLocators, att.DataLimits)

	for i := uint32(0); i < att.Participants.Initial; i++ {
		proxyData := data.NewParticipantProxyData(att)
		pdp.participantProxiesPool = append(pdp.participantProxiesPool, proxyData)
	}

	for i := uint32(0); i < att.TotalReaders().Initial; i++ {
		proxyData := data.NewReaderProxyData(maxUnicastLocators, maxMulticastLocators, att.DataLimits)
		pdp.readerProxiesPool = append(pdp.readerProxiesPool, proxyData)
	}

	for i := uint32(0); i < att.TotalWriters().Initial; i++ {
		proxyData := data.NewWriterProxyData(maxUnicastLocators, maxMulticastLocators, att.DataLimits)
		pdp.writerProxiesPool = append(pdp.writerProxiesPool, proxyData)
	}

	pdp.impl = impl

	return &pdp
}
