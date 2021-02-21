package EDP

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/core/status"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/data"
)

// type Endpoint struct {
// 	guid  common.GUIDT
// 	att   attributes.EndpointAttributes
// 	mutex sync.Mutex
// }

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
	InitEDP(att attributes.BuiltinAttributes) bool

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

type EDP struct {
	tempReaderProxyData *data.ReaderProxyData
	tempWriterProxyData *data.WriterProxyData
	readerStatus        map[common.GUIDT]status.SubscriptionMatchedStatus
	writerStatus        map[common.GUIDT]status.PublicationMatchedStatus
	pPDP                protocol.IPDP
	rtpsParticipant     protocol.IParticipant
}

func (edp *EDP) RemoveRemoteEndpoints(pdata *data.ParticipantProxyData) {
	log.Printlln("RemoveRemoteEndpoints do nothing")
}

func (edp *EDP) AreRemoteEndpointsMatched(pdata *data.ParticipantProxyData) bool {
	return false
}

// Create a new ReaderPD for a local Reader.
func (edp *EDP) NewLocalReaderProxyData(areader reader.IRTPSReader, att *attributes.TopicAttributes,
	rqos *qos.ReaderQos) bool {
	log.Printlln("Adding", areader.GetGUID().EntityID, " in topic", att.TopicName)
	initFun := func(rpd *data.ReaderProxyData, updating bool, participantData *data.ParticipantProxyData) {
		if updating {
			log.Fatalln("Adding already existent reader ", areader.GetGUID().EntityID, " in topic ", att.TopicName)
			return false
		}

		networkFactory := edp.rtpsParticipant.NetworkFactory()
		rpd.IsAlive(true)
		rpd.ExpectsInlineQos = areader.ExpectsInlineQos()
		rpd.GUID = areader.GetGUID()
		rpd.Key = rpd.GUID
		ratt := areader.GetAttributes()
		if ratt.MulticastLocatorList.Length() == 0 && ratt.UnicastLocatorList.Length() == 0 {
			rpd.SetLocators(participantData.DefaultLocators)
		} else {

		}
	}
}

// Create a new ReaderPD for a local Writer.
func (edp *EDP) NewLocalWriterProxyData(awriter writer.IRTPSWriter, att *attributes.TopicAttributes,
	rqos *qos.WriterQos) bool {

}

// A previously created Reader has been updated
func (edp *EDP) UpdateLocalReader(areader reader.IRTPSReader, att *attributes.TopicAttributes,
	rqos *qos.ReaderQos) bool {
}

// A previously created Writer has been updated
func (edp *EDP) UpdateLocalWriter(awriter writer.IRTPSWriter, att *attributes.TopicAttributes,
	wqos *qos.WriterQos) bool {

}

// Check the validity of a matching between a local RTPSWriter and a ReaderProxyData object.
func (edp *EDP) ValidMatching(wdata *writer.WriterProxyData, rdata *reader.ReaderProxyData) bool {

}

// Unpair a WriterProxyData object from all local readers.
func (edp *EDP) UnpairWriterProxy(participantGUID, writerGUID *common.GUIDT, removedByLease bool) bool {

}

// Unpair a ReaderProxyData object from all local writers.
func (edp *EDP) UnpairReaderProxy(participantGUID, readerGUID *common.GUIDT) bool {

}

func (edp *EDP) PairingReaderProxyWithAnyLocalWriter(participantGUID *common.GUIDT, rdata *data.ReaderProxyData) bool {

}

func (edp *EDP) PairingWriterProxyWithAnyLocalReader(participantGUID *common.GUIDT, wdata *data.WriterProxyData) bool {

}

func (edp *EDP) UpdateSubscriptionMatchedStatus(readerGuid, writerGuid *common.GUIDT, achange int) {

}

func (edp *EDP) UpdatePublicationMatchedStatus(readerGuid, writerGuid *common.GUIDT, achange int) {

}
