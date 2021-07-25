package edp

import (
	"log"
	"math"
	"sync"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/data"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/protocol"
	"github.com/yeren0143/DDS/fastrtps/rtps/history"
	"github.com/yeren0143/DDS/fastrtps/rtps/reader"
	"github.com/yeren0143/DDS/fastrtps/rtps/writer"
)

var _ IEDP = (*EDPSimple)(nil)
var _ edpImpl = (*EDPSimple)(nil)

const (
	KEDPReaderInitialReservedCaches uint32 = 1
	KEDPWriterInitialReservedCaches uint32 = 20
)

var (
	KEDPHeartbeatPeriod        common.DurationT = common.DurationT{1, 0}
	KEDPNackResponseDelay      common.DurationT = common.DurationT{0, 100 * 1000}
	KEDPNackSupressionDuration common.DurationT = common.DurationT{0, 10 * 1000}
	KEDPHeartbeatResponseDelay common.DurationT = common.DurationT{0, 10 * 1000}
)

// Class EDPSimple, implements the Simple Endpoint Discovery Protocol defined in the RTPS specification.
// Inherits from EDP class.
type EDPSimple struct {
	edpBase
	Discovery             attributes.BuiltinAttributes
	PublicationsListener  IEDPListener
	SubscriptionsListener IEDPListener
	pubWriterPayloadPool  history.ITopicPayloadPool
	pubReaderPayloadPool  history.ITopicPayloadPool
	subWriterPayloadPool  history.ITopicPayloadPool
	subReaderPayloadPool  history.ITopicPayloadPool
	tempDataMutex         sync.Mutex
	tempReaderProxyData   data.ReaderProxyData
	tempWriterProxyData   data.WriterProxyData
	// Pointer to the Publications Writer (only created if indicated in the DiscoveryAtributes).
	publicationsWriter WriterHistoryPair
	// Pointer to the Subscriptions Writer (only created if indicated in the DiscoveryAtributes).
	subscriptionsWriter WriterHistoryPair
	// Pointer to the Publications Reader (only created if indicated in the DiscoveryAtributes).
	publicationsReader ReaderHistoryPair
	// Pointer to the Subscriptions Reader (only created if indicated in the DiscoveryAtributes)
	subscriptionsReader ReaderHistoryPair
}

func (edp *EDPSimple) AssignRemoteEndpoints(pdata *data.ParticipantProxyData) {
	log.Fatalln("not Impl")
}

//Initialization of history attributes for EDP built-in readers
func (edp *EDPSimple) setBuiltinReaderHistoryAttributes(att *attributes.HistoryAttributes) {
	att.InitialReservedCaches = KEDPReaderInitialReservedCaches
	att.PayloadMaxSize = edp.pPDP.BuiltinAttributes().ReaderPayloadSize
	att.MemoryPolicy = edp.pPDP.BuiltinAttributes().ReaderHostoryMemoryPolicy
}

//Initialization of history attributes for EDP built-in writers
func (edp *EDPSimple) setBuiltinWriterHistoryAttributes(att *attributes.HistoryAttributes) {
	att.InitialReservedCaches = KEDPWriterInitialReservedCaches
	att.PayloadMaxSize = edp.pPDP.BuiltinAttributes().WriterPayloadSize
	att.MemoryPolicy = edp.pPDP.BuiltinAttributes().WriterHistoryMemoryPolicy
}

//Initialization of reader attributes for EDP built-in readers
func (edp *EDPSimple) setBuiltinReaderAttributes(att *attributes.ReaderAttributes) {
	// Matched writers will depend on total number of participants
	att.MatchedWritersAllocation = *edp.pPDP.GetRTPSParticipant().GetAttributes().Allocation.Participants

	// As participants allocation policy includes the local participant, one has to be substracted
	if att.MatchedWritersAllocation.Initial > 1 {
		att.MatchedWritersAllocation.Initial--
	}
	if att.MatchedWritersAllocation.Maximum > 1 &&
		att.MatchedWritersAllocation.Maximum < math.MaxUint32 {
		att.MatchedWritersAllocation.Maximum--
	}

	// Locators are copied from the local participant metatraffic locators
	att.EndpointAtt.UnicastLocatorList.Clear()
	metatrafficLocators := edp.pPDP.GetLocalParticipantProxyData().MetatrafficLocators
	for _, loc := range metatrafficLocators.Unicast {
		att.EndpointAtt.UnicastLocatorList.PushBack(&loc)
	}
	att.EndpointAtt.MulticastLocatorList.Clear()
	for _, loc := range metatrafficLocators.Multicast {
		att.EndpointAtt.MulticastLocatorList.PushBack(&loc)
	}

	// Timings are configured using EDP default values
	att.Times.HeartBeatResponseDelay = KEDPHeartbeatResponseDelay

	// EDP endpoints are always reliable, transsient local, keyed topics
	att.EndpointAtt.ReliabilityKind = common.KReliable
	att.EndpointAtt.DurabilityKind = common.KTransientLocal
	att.EndpointAtt.TopicKind = common.KWithKey

	// Built-in EDP readers never expect inline qos
	att.ExpectsInlineQos = false
}

//Initialization of writer attributes for EDP built-in writers
func (edp *EDPSimple) setBuiltinWriterAttributes(att *attributes.WriterAttributes) {
	// Matched readers will depend on total number of participants
	att.MatchedReadersAllocation = *edp.pPDP.GetRTPSParticipant().GetAttributes().Allocation.Participants

	// As participants allocation policy includes the local participant, one has to be substracted
	if att.MatchedReadersAllocation.Initial > 1 {
		att.MatchedReadersAllocation.Initial--
	}
	if att.MatchedReadersAllocation.Maximum > 1 &&
		att.MatchedReadersAllocation.Maximum < math.MaxUint32 {
		att.MatchedReadersAllocation.Maximum--
	}

	metatrafficLocators := edp.pPDP.GetLocalParticipantProxyData().MetatrafficLocators
	// Locators are copied from the local participant metatraffic locators
	att.EndpointAtt.UnicastLocatorList.Clear()
	for _, loc := range metatrafficLocators.Unicast {
		att.EndpointAtt.UnicastLocatorList.PushBack(&loc)
	}
	att.EndpointAtt.MulticastLocatorList.Clear()
	for _, loc := range edp.pPDP.GetLocalParticipantProxyData().MetatrafficLocators.Multicast {
		att.EndpointAtt.MulticastLocatorList.PushBack(&loc)
	}

	// Timings are configured using EDP default values
	att.Times.HeartbeatPeriod = KEDPHeartbeatPeriod
	att.Times.NackResponseDelay = KEDPNackResponseDelay
	att.Times.NackSupressionDuration = KEDPNackSupressionDuration

	// EDP endpoints are always reliable, transsient local, keyed topics
	att.EndpointAtt.ReliabilityKind = common.KReliable
	att.EndpointAtt.DurabilityKind = common.KTransientLocal
	att.EndpointAtt.TopicKind = common.KWithKey

	// Set as asynchronous if there is a throughput controller installed
	throughputController := edp.rtpsParticipant.GetAttributes().ThroughputController
	if throughputController.BytesPerPeriod != math.MaxUint32 &&
		throughputController.PeriodMillisecs != 0 {
		att.PubMode = attributes.KAsynchronousWriter
	}
}

func (edp *EDPSimple) createSEDPEndpoints() bool {
	watt := attributes.NewWriterAttributes()
	ratt := attributes.NewReaderAttributes()
	readerHistoryAtt := attributes.KDefaultHistoryAttributes
	writerHistoryAtt := attributes.KDefaultHistoryAttributes

	edp.setBuiltinReaderHistoryAttributes(&readerHistoryAtt)
	edp.setBuiltinWriterHistoryAttributes(&writerHistoryAtt)
	edp.setBuiltinReaderAttributes(ratt)
	edp.setBuiltinWriterAttributes(watt)

	edp.PublicationsListener = newEDPSimplePubListener(edp)
	edp.SubscriptionsListener = newEDPSimpleSubListener(edp)

	if edp.Discovery.DiscoveryConfig.SimpleEDP.UsePublicationWriterAndSubscriptionReader {
		if !createEDPWriter(edp.rtpsParticipant, "DCPSPublications", common.KEntityIDSEDPPubWriter,
			&writerHistoryAtt, watt, edp.PublicationsListener, edp.pubWriterPayloadPool, &edp.publicationsWriter) {
			log.Println("createEDPWriter failed")
			return false
		}
		log.Println("SEDP Publication Writer created")

		if !createEDPReader(edp.rtpsParticipant, "DCPSSubscriptions", common.KEntityIDSEDPSubReader,
			&readerHistoryAtt, ratt, edp.SubscriptionsListener, edp.subReaderPayloadPool, &edp.subscriptionsReader) {
			log.Println("createEDPReader failed")
			return false
		}
		log.Println("SEDP Subscription Reader created")
	}

	if edp.Discovery.DiscoveryConfig.SimpleEDP.UsePublicationReaderAndSubscriptionWriter {
		if !createEDPReader(edp.rtpsParticipant, "DCPSPublications", common.KEntityIDSEDPPubReader,
			&readerHistoryAtt, ratt, edp.PublicationsListener, edp.pubReaderPayloadPool, &edp.publicationsReader) {
			log.Println("createEDPReader failed")
			return false
		}
		log.Println("SEDP Publication Reader created")

		if !createEDPWriter(edp.rtpsParticipant, "DCPSSubscriptions", common.KEntityIDSEDPSubWriter,
			&writerHistoryAtt, watt, edp.SubscriptionsListener, edp.subWriterPayloadPool, &edp.subscriptionsWriter) {
			log.Println("createEDPWriter failed")
			return false
		}

		log.Println("SEDP Subscription Writer created")
	}

	log.Println("Creation finished")
	return true
}

func (edp *EDPSimple) InitEDP(att *attributes.BuiltinAttributes) bool {
	log.Println("Beginning Simple Endpoint Discovery Protocol")
	edp.Discovery = *att

	if !edp.createSEDPEndpoints() {
		log.Fatalln("Problem creation SimpleEDP endpoints")
		return false
	}
	return true
}

func (edp *EDPSimple) ProcessLocalWriterProxyData(awriter *writer.IRTPSWriter, wdata *data.WriterProxyData) bool {
	log.Fatalln("not Impl")
	return false
}

func (edp *EDPSimple) ProcessLocalReaderProxyData(areader *reader.IRTPSReader, rdata *data.ReaderProxyData) bool {
	log.Fatalln("not Impl")
	return false
}

func (edp *EDPSimple) RemoveLocalWriter(awriter *writer.IRTPSWriter) bool {
	log.Fatalln("not Impl")
	return false
}

func (edp *EDPSimple) RemoveLocalReader(areader *reader.IRTPSReader) bool {
	log.Fatalln("not Impl")
	return false
}

func NewEDPSimple(p protocol.IPDP, part protocol.IParticipant) *EDPSimple {
	var edpSimple EDPSimple
	edpSimple.edpBase = *NewEDPBase(p, part)
	edpSimple.edpImpl = &edpSimple

	return &edpSimple
}
