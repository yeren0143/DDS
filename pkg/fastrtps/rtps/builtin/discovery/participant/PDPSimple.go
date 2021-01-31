package participant

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/data"
	"github.com/yeren0143/DDS/fastrtps/rtps/history"
	"log"
	"math"
)

var _ IPDP = (*PDPSimple)(nil)
var _ IpdpBaseImpl = (*PDPSimple)(nil)

type PDPSimple struct {
	pdpBase
}

func (pdp *PDPSimple) AssignRemoteEndpoints(pdata *data.ParticipantProxyData) {
	log.Println("For RTPSParticipant: ", pdata.Guid.Prefix)
}

func (pdp *PDPSimple) NotifyAboveRemoteEndpoints(pdata *data.ParticipantProxyData) {

}

func (pdp *PDPSimple) RemoveRemoteEndpoints(pdata *data.ParticipantProxyData) {

}

func (pdp *PDPSimple) CreateParticipantProxyData(p *data.ParticipantProxyData, writer_guid *common.GUIDT) *data.ParticipantProxyData {
	return nil
}

func (pdp *PDPSimple) Init(participant IRTPSParticipant) bool {
	// The DATA(p) must be processed after EDP endpoint creation
	if !pdp.initPDP(participant) {
		return false
	}

	//INIT EDP
	if pdp.discovery.DiscoveryConfig.UseStaticEndpoint {

	}

	return false
}

func (pdp *PDPSimple) CreatePDPEndpoints() bool {
	log.Printf("CreatePDPEndpoints Beginning")

	allocation := pdp.participant.GetAttributes().Allocation

	// SPDP BUILTIN RTPSParticipant READER
	var hatt attributes.HistoryAttributes
	hatt.PayloadMaxSize = pdp.builtin.GetBuiltinAttributes().ReaderPayloadSize
	hatt.MemoryPolicy = pdp.builtin.GetBuiltinAttributes().ReaderHostoryMemoryPolicy
	hatt.InitialReservedCaches = 25
	if allocation.Participants.Initial > 0 {
		hatt.InitialReservedCaches = allocation.Participants.Initial
	}
	if allocation.Participants.Maximum < math.MaxUint32 {
		hatt.MaximumReservedCaches = allocation.Participants.Maximum
	}

	readerPoolCfg := history.FromHistoryAttributes(&hatt)
	pdp.readerPayloadPool = history.GetTopicPayloadPoolProxy("DCPSParticipant", &readerPoolCfg)
	pdp.readerPayloadPool.ReserveHistory(&readerPoolCfg, true)

	pdp.pdpReaderHistory = history.NewReaderHistory(&hatt)
	var ratt attributes.ReaderAttributes
	ratt.EndpointAtt.MulticastLocatorList = pdp.builtin.Metra

	return false
}

func NewPDPSimple(protocol IBuiltinProtocols, att *attributes.RTPSParticipantAllocationAttributes) *PDPSimple {
	var pdpSimple PDPSimple
	pdpSimple.pdpBase = *newPDP(protocol, att, &pdpSimple)

	return &pdpSimple
}
