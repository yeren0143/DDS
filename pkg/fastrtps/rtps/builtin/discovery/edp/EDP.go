package edp

import (
	"log"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/core/status"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/data"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/protocol"
	"github.com/yeren0143/DDS/fastrtps/rtps/qos"
	"github.com/yeren0143/DDS/fastrtps/rtps/reader"
	"github.com/yeren0143/DDS/fastrtps/rtps/writer"
)

/**
 * Mask to hold the reasons why two endpoints do not match.
 */
type MatchingFailureMask uint32

const (
	KDifferentTopic    MatchingFailureMask = (0x00000001 << 0)
	KInconsistentTopic MatchingFailureMask = (0x00000001 << 1)
	KIncompatibleQos   MatchingFailureMask = (0x00000001 << 2)
	KPartitions        MatchingFailureMask = (0x00000001 << 3)
)

/**
 * Class EDP, base class for Endpoint Discovery Protocols.
 * It contains generic methods used by the two EDP implemented (EDPSimple and EDPStatic),
 * as well as abstract methods definitions required by the specific implementations.
 */
type IEDP interface {
	InitEDP(att *attributes.BuiltinAttributes) bool

	// Abstract method that assigns remote endpoints when a new RTPSParticipantProxyData is discovered.
	AssignRemoteEndpoints(pdata *data.ParticipantProxyData)

	// Remove remote endpoints from the endpoint discovery protocol
	RemoveRemoteEndpoints(pdata *data.ParticipantProxyData)

	// Verify if the given participant EDP enpoints are matched with us
	AreRemoteEndpointsMatched(pdata *data.ParticipantProxyData) bool

	// Abstract method that removes a local Reader from the discovery method
	RemoveLocalReader(areader *reader.IRTPSReader) bool

	// Abstract method that removes a local Writer from the discovery method
	RemoveLocalWriter(awriter *writer.IRTPSWriter) bool

	// After a new local ReaderProxyData has been created some processing is needed (depends on the implementation).
	ProcessLocalReaderProxyData(areader *reader.IRTPSReader, rdata *data.ReaderProxyData) bool

	// After a new local WriterProxyData has been created some processing is needed (depends on the implementation).
	ProcessLocalWriterProxyData(awriter *writer.IRTPSWriter, wdata *data.WriterProxyData) bool
}

type edpImpl interface {
	ProcessLocalReaderProxyData(areader *reader.IRTPSReader, rdata *data.ReaderProxyData) bool
}

type edpBase struct {
	edpImpl
	tempReaderProxyData *data.ReaderProxyData
	tempWriterProxyData *data.WriterProxyData
	readerStatus        map[common.GUIDT]status.SubscriptionMatchedStatus
	writerStatus        map[common.GUIDT]status.PublicationMatchedStatus
	pPDP                protocol.IPDP
	rtpsParticipant     protocol.IParticipant
}

func NewEDPBase(p protocol.IPDP, part protocol.IParticipant) *edpBase {
	var edp edpBase
	edp.pPDP = p
	edp.rtpsParticipant = part
	alloc := part.GetAttributes().Allocation
	locators := alloc.Locators
	edp.tempReaderProxyData = data.NewReaderProxyData(locators.MaxUnicastLocators,
		locators.MaxMulticastLocators, alloc.DataLimits)
	edp.tempWriterProxyData = data.NewWriterProxyData(locators.MaxUnicastLocators,
		locators.MaxMulticastLocators, alloc.DataLimits)
	edp.readerStatus = make(map[common.GUIDT]status.SubscriptionMatchedStatus, alloc.TotalReaders().Initial)
	edp.writerStatus = make(map[common.GUIDT]status.PublicationMatchedStatus, alloc.TotalReaders().Initial)
	return &edp
}

func (edp *edpBase) RemoveRemoteEndpoints(pdata *data.ParticipantProxyData) {
	log.Println("RemoveRemoteEndpoints do nothing")
}

func (edp *edpBase) AreRemoteEndpointsMatched(pdata *data.ParticipantProxyData) bool {
	return false
}

// Create a new ReaderPD for a local Reader.
func (edp *edpBase) NewLocalReaderProxyData(areader reader.IRTPSReader, att *attributes.TopicAttributes,
	rqos *qos.ReaderQos) bool {
	log.Println("Adding", areader.GetGUID().EntityID, " in topic", att.TopicName)
	initFun := func(rpd *data.ReaderProxyData, updating bool, participantData *data.ParticipantProxyData) bool {
		if updating {
			log.Fatalln("Adding already existent reader ", areader.GetGUID().EntityID, " in topic ", att.TopicName)
			return false
		}

		networkFactory := edp.rtpsParticipant.NetworkFactory()
		rpd.SetAlive(true)
		rpd.ExpectsInlineQos = areader.ExpectsInlineQos()
		rpd.GUID = *areader.GetGUID()
		rpd.Key.InitWithGUID(&rpd.GUID)
		ratt := areader.GetAttributes()
		if ratt.MulticastLocatorList.Length() == 0 && ratt.UnicastLocatorList.Length() == 0 {
			rpd.SetLocators(&participantData.DefaultLocators)
		} else {
			rpd.SetMulticastLocators(&areader.GetAttributes().MulticastLocatorList, networkFactory)
			rpd.SetAnnouncedUnicastLocators(&areader.GetAttributes().UnicastLocatorList)
		}

		log.Fatalln("notImpl")

		return false
	}

	//ADD IT TO THE LIST OF READERPROXYDATA
	var participantGUID common.GUIDT
	readerData := edp.pPDP.AddReaderProxyData(areader.GetGUID(), &participantGUID, initFun)
	if readerData == nil {
		return false
	}

	//PAIRING
	edp.PairingReaderProxyWithAnyLocalWriter(&participantGUID, readerData)
	edp.pairingReader(areader, &participantGUID, readerData)
	//DO SOME PROCESSING DEPENDING ON THE IMPLEMENTATION (SIMPLE OR STATIC)
	edp.edpImpl.ProcessLocalReaderProxyData(&areader, readerData)
	return true
}

// Create a new ReaderPD for a local Writer.
func (edp *edpBase) NewLocalWriterProxyData(awriter writer.IRTPSWriter, att *attributes.TopicAttributes,
	rqos *qos.WriterQos) bool {
	log.Fatalln("notImpl")
	return false
}

func (edp *edpBase) pairingReader(r reader.IRTPSReader, participantGUID *common.GUIDT, rdata *data.ReaderProxyData) bool {
	log.Fatalln("notImpl")
	return false
}

// A previously created Reader has been updated
func (edp *edpBase) UpdateLocalReader(areader reader.IRTPSReader, att *attributes.TopicAttributes,
	rqos *qos.ReaderQos) bool {
	log.Fatalln("notImpl")
	return false
}

// A previously created Writer has been updated
func (edp *edpBase) UpdateLocalWriter(awriter writer.IRTPSWriter, att *attributes.TopicAttributes,
	wqos *qos.WriterQos) bool {
	log.Fatalln("notImpl")
	return false

}

// Check the validity of a matching between a local RTPSWriter and a ReaderProxyData object.
func (edp *edpBase) ValidMatching(wdata *data.WriterProxyData, rdata *data.ReaderProxyData) bool {
	log.Fatalln("notImpl")
	return false
}

// Unpair a WriterProxyData object from all local readers.
func (edp *edpBase) UnpairWriterProxy(participantGUID, writerGUID *common.GUIDT, removedByLease bool) bool {
	log.Fatalln("notImpl")
	return false
}

// Unpair a ReaderProxyData object from all local writers.
func (edp *edpBase) UnpairReaderProxy(participantGUID, readerGUID *common.GUIDT) bool {
	log.Fatalln("notImpl")
	return false
}

func (edp *edpBase) PairingReaderProxyWithAnyLocalWriter(participantGUID *common.GUIDT, rdata *data.ReaderProxyData) bool {
	log.Fatalln("notImpl")
	return false
}

func (edp *edpBase) PairingWriterProxyWithAnyLocalReader(participantGUID *common.GUIDT, wdata *data.WriterProxyData) bool {
	log.Fatalln("notImpl")
	return false
}

func (edp *edpBase) UpdateSubscriptionMatchedStatus(readerGuid, writerGuid *common.GUIDT, achange int) {

}

func (edp *edpBase) UpdatePublicationMatchedStatus(readerGuid, writerGuid *common.GUIDT, achange int) {

}
