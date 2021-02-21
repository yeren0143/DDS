package participant

import (
	"log"
	"math"
	"os"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/data"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/protocol"
	"github.com/yeren0143/DDS/fastrtps/rtps/history"
	"github.com/yeren0143/DDS/fastrtps/rtps/writer"
)

var _ protocol.IPDP = (*PDPSimple)(nil)
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

func (pdp *PDPSimple) Init(participant protocol.IParticipant) bool {
	// The DATA(p) must be processed after EDP endpoint creation
	if !pdp.initPDP(participant) {
		return false
	}

	//INIT EDP
	if pdp.discovery.DiscoveryConfig.UseStaticEndpoint {
		log.Fatalln("not impl")
	}

	return false
}

func (pdp *PDPSimple) CreatePDPEndpoints() bool {
	log.Println("CreatePDPEndpoints Beginning")

	allocation := pdp.rtpsParticipant.GetAttributes().Allocation

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
	ratt.EndpointAtt.MulticastLocatorList = *pdp.builtin.GetMetatrafficMulticastLocatorList()
	ratt.EndpointAtt.UnicastLocatorList = *pdp.builtin.GetMetatrafficUnicastLocatorList()
	ratt.EndpointAtt.TopicKind = common.KNoKey
	ratt.EndpointAtt.DurabilityKind = common.KTransientLocal
	ratt.EndpointAtt.ReliabilityKind = common.KBestEffort
	ratt.MatchedWritersAllocation = *allocation.Participants

	pdp.listener = newPDPListener(pdp)
	var ok bool
	ok, pdp.reader = pdp.rtpsParticipant.CreateReader(&ratt, pdp.readerPayloadPool,
		pdp.pdpReaderHistory, pdp.listener,
		common.KEidSEDPBuiltinTopicReader, true, false)
	if ok {
		if os.Getenv("HAVE_SECURITY") != "" {
			log.Fatal("todo HAVE_SECURITY CreatePDPEndpoints")
		}
	} else {
		log.Fatal("SimplePDP Reader creation failed")
		pdp.pdpReaderHistory = nil
		pdp.listener = nil
		pdp.readerPayloadPool.ReleaseHistory(&readerPoolCfg, true)
		history.ReleaseTopicPayloadPool(pdp.readerPayloadPool)
		pdp.readerPayloadPool = nil
		return false
	}

	//SPDP BUILTIN RTPSParticipant WRITER
	hatt.PayloadMaxSize = pdp.builtin.GetBuiltinAttributes().WriterPayloadSize
	hatt.InitialReservedCaches = 1
	hatt.MaximumReservedCaches = 1
	hatt.MemoryPolicy = pdp.builtin.GetBuiltinAttributes().WriterHistoryMemoryPolicy

	writerPoolCfg := history.FromHistoryAttributes(&hatt)
	pdp.writerPayloadPool = history.GetTopicPayloadPoolProxy("DCPSParticipant", &writerPoolCfg)
	pdp.writerPayloadPool.ReserveHistory(&writerPoolCfg, false)

	pdp.pdpWriterHistory = history.NewWriterHistory(&hatt)
	watt := attributes.NewWriterAttributes()
	watt.EndpointAtt.EndpointKind = common.KWriter
	watt.EndpointAtt.DurabilityKind = common.KTransientLocal
	watt.EndpointAtt.ReliabilityKind = common.KBestEffort
	watt.EndpointAtt.TopicKind = common.KWithKey
	watt.EndpointAtt.RemoteLocatorList.Set(pdp.discovery.InitialPeersList)
	watt.MatchedReadersAllocation = *allocation.Participants

	control := pdp.rtpsParticipant.GetAttributes().ThroughputController
	if control.BytesPerPeriod != math.MaxUint32 && control.PeriodMillisecs != 0 {
		watt.PubMode = attributes.KAsynchronousWriter
	}

	//var wout writer.IRTPSWriter
	wout, ok := pdp.rtpsParticipant.CreateWriter(watt, pdp.writerPayloadPool, pdp.pdpWriterHistory, nil,
		common.KEntityIDSPDPWriter, true)
	if ok {
		pdp.writer = wout
		if pdp.writer != nil {
			var fixedLocators common.LocatorList
			network := pdp.rtpsParticipant.NetworkFactory()
			for _, loc := range pdp.builtin.GetInitialPeers().Locators {
				var localLocator common.Locator
				if network.TransformRemoteLocator(&loc, &localLocator) {
					fixedLocators.PushBack(&localLocator)
				}
			}
			statelessWriter := wout.(*writer.StatelessWriter)
			statelessWriter.SetFixedLocators(fixedLocators)
		}
	} else {
		log.Fatalln("SimplePDP Writer creation failed")
		pdp.pdpWriterHistory = nil
		pdp.writerPayloadPool.ReleaseHistory(&writerPoolCfg, false)
		history.ReleaseTopicPayloadPool(pdp.writerPayloadPool)
		pdp.writerPayloadPool = nil
		return false
	}
	log.Println("SPDP Endpoints creation finished")

	return true
}

func NewPDPSimple(protocol IPDPParent, att *attributes.RTPSParticipantAllocationAttributes) *PDPSimple {
	var pdpSimple PDPSimple
	pdpSimple.pdpBase = *newPDP(protocol, att, &pdpSimple)

	return &pdpSimple
}
