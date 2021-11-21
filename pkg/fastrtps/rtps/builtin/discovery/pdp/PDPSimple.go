package pdp

import (
	"log"
	"math"
	"os"

	"dds/common"
	"dds/core/policy"
	"dds/fastrtps/rtps/attributes"
	"dds/fastrtps/rtps/builtin/data"
	"dds/fastrtps/rtps/builtin/discovery/edp"
	"dds/fastrtps/rtps/builtin/discovery/protocol"
	"dds/fastrtps/rtps/history"
	"dds/fastrtps/rtps/writer"
)

var _ protocol.IPDP = (*PDPSimple)(nil)
var _ IpdpBaseImpl = (*PDPSimple)(nil)

type PDPSimple struct {
	pdpBase
}

func (pdp *PDPSimple) AssignRemoteEndpoints(pdata *data.ParticipantProxyData) {
	log.Println("For RTPSParticipant: ", pdata.Guid.Prefix)
	//log.Fatalln("not impl")
	network := pdp.rtpsParticipant.NetworkFactory()
	endp := pdata.AviableBuiltinEndpoints
	useMulticastLocators := pdp.rtpsParticipant.GetAttributes().Builtin.AvoidBuiltinMulticast ||
		len(pdata.MetatrafficLocators.Unicast) == 0
	auxendp := endp & data.DISC_BUILTIN_ENDPOINT_PARTICIPANT_ANNOUNCER
	if auxendp != 0 {
		pdp.tempDataMutex.Lock()
		pdp.tempWriterData.Clear()
		pdp.tempWriterData.Guid().Prefix = pdata.Guid.Prefix
		pdp.tempWriterData.Guid().EntityID = *common.KEntityIDSPDPWriter
		guid := pdata.GetPersistenceGuid()
		pdp.tempWriterData.SetPersistenceGuid(&guid)
		pdp.tempWriterData.SetPersistenceEntityID(common.KEntityIDSPDPWriter)
		pdp.tempWriterData.SetRemoteLocators(&pdata.MetatrafficLocators, network, useMulticastLocators)
		pdp.tempWriterData.Qos.Reliability.Kind = policy.RELIABLE_RELIABILITY_QOS
		pdp.tempWriterData.Qos.Durability.Kind = policy.KTransientLocalDurabilityQos
		pdp.reader.MatchedWriterAdd(pdp.tempWriterData)
		pdp.tempDataMutex.Unlock()
	}
	auxendp = endp
	auxendp &= data.DISC_BUILTIN_ENDPOINT_PARTICIPANT_DETECTOR
	if auxendp != 0 {
		pdp.tempDataMutex.Lock()
		pdp.tempReaderData.Clear()
		pdp.tempReaderData.ExpectsInlineQos = false
		pdp.tempReaderData.GUID.Prefix = pdata.Guid.Prefix
		pdp.tempReaderData.GUID.EntityID = *common.KEntityIDSPDPReader
		pdp.tempReaderData.SetRemoteLocators(&pdata.MetatrafficLocators, network, useMulticastLocators)
		pdp.tempReaderData.Qos.Reliability = policy.BEST_EFFORT_RELIABILITY_QOS
		pdp.tempReaderData.Qos.Durability = policy.KTransientLocalDurabilityQos
		pdp.writer.MatchedReaderAdd(pdp.tempReaderData)
		pdp.tempDataMutex.Unlock()
	}
	// log.Fatalln("not impl")

}

func (pdp *PDPSimple) NotifyAboveRemoteEndpoints(pdata *data.ParticipantProxyData) {
	log.Fatalln("not impl")
}

func (pdp *PDPSimple) RemoveRemoteEndpoints(pdata *data.ParticipantProxyData) {
	log.Fatalln("not impl")
}

func (pdp *PDPSimple) CreateParticipantProxyData(participantData *data.ParticipantProxyData, writer_guid *common.GUIDT) *data.ParticipantProxyData {
	// TODO::
	//pdp.GetMutex().Lock()

	// decide if we dismiss the participant using the ParticipantFilteringFlags
	flags := pdp.discovery.DiscoveryConfig.IgnoreParticipantFlags
	if flags != attributes.KNoFilter {
		log.Fatalln("not Impl")
	}
	pdata := pdp.addParticipantProxyData(&participantData.Guid, true)
	if pdata != nil {
		pdata.Copy(participantData)
		pdata.IsAlive = true
		interval := common.NewTime(&pdata.LeaseDuration)
		pdata.LeaseDurationEvent.UpdateInterval(*interval)
		pdata.LeaseDurationEvent.RestartTimer()
	}

	return pdata
}

func (pdp *PDPSimple) Init(participant protocol.IParticipant) bool {
	// The DATA(p) must be processed after EDP endpoint creation
	if !pdp.initPDP(participant) {
		return false
	}

	//INIT EDP
	if pdp.discovery.DiscoveryConfig.UseStaticEndpoint {
		log.Fatalln("not impl")
	} else if pdp.discovery.DiscoveryConfig.UseSimpleEndpoint {
		pdp.EDP = edp.NewEDPSimple(pdp, pdp.rtpsParticipant)
		if !pdp.EDP.InitEDP(pdp.discovery) {
			log.Fatalln("Endpoint discovery configuration failed")
			return false
		}
	} else {
		log.Fatalln("No EndpointDiscoveryProtocol defined")
		return false
	}

	return true
}

// Force the sending of our local DPD to all remote RTPSParticipants and multicast Locators.
func (pdp *PDPSimple) AnnounceParticipantState(newChange, dispose bool, wparams *common.WriteParamsT) {
	pdp.pdpBase.AnnounceParticipantState(newChange, dispose, wparams)

	if !(dispose || newChange) {
		statelessWriter, ok := pdp.writer.(*writer.StatelessWriter)
		if ok {
			statelessWriter.UnsentChangesReset()
		} else {
			log.Fatalln("Using PDPSimple protocol with a reliable writer")
		}
	}
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
	pdp.readerPayloadPool = history.GetTopicPayloadPoolProxy("DCPSParticipant", readerPoolCfg)
	pdp.readerPayloadPool.ReserveHistory(readerPoolCfg, true)

	pdp.pdpReaderHistory = history.NewReaderHistory(&hatt)
	ratt := attributes.NewReaderAttributes()
	ratt.EndpointAtt.MulticastLocatorList = *pdp.builtin.GetMetatrafficMulticastLocatorList()
	ratt.EndpointAtt.UnicastLocatorList = *pdp.builtin.GetMetatrafficUnicastLocatorList()
	ratt.EndpointAtt.TopicKind = common.KNoKey
	ratt.EndpointAtt.DurabilityKind = common.KTransientLocal
	ratt.EndpointAtt.ReliabilityKind = common.KBestEffort
	ratt.MatchedWritersAllocation = *allocation.Participants

	pdp.listener = newPDPListener(pdp)
	var ok bool
	pdp.reader, ok = pdp.rtpsParticipant.CreateReader(ratt, pdp.readerPayloadPool,
		pdp.pdpReaderHistory, pdp.listener,
		common.KEntityIDSPDPReader, true, false)
	if ok {
		if os.Getenv("HAVE_SECURITY") != "" {
			log.Fatal("todo HAVE_SECURITY CreatePDPEndpoints")
		}
	} else {
		log.Fatal("SimplePDP Reader creation failed")
		pdp.pdpReaderHistory = nil
		pdp.listener = nil
		pdp.readerPayloadPool.ReleaseHistory(readerPoolCfg, true)
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
	pdp.writerPayloadPool = history.GetTopicPayloadPoolProxy("DCPSParticipant", writerPoolCfg)
	pdp.writerPayloadPool.ReserveHistory(writerPoolCfg, false)

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
	ok, wout := pdp.rtpsParticipant.CreateWriter(watt, pdp.writerPayloadPool, pdp.pdpWriterHistory, nil,
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
		pdp.writerPayloadPool.ReleaseHistory(writerPoolCfg, false)
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
