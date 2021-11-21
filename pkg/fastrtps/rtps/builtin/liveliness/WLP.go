package liveliness

import (
	"log"
	"math"
	"sync"

	"dds/common"
	"dds/core/policy"
	"dds/fastrtps/rtps/attributes"
	"dds/fastrtps/rtps/builtin/data"
	"dds/fastrtps/rtps/builtin/discovery/protocol"
	"dds/fastrtps/rtps/endpoint"
	"dds/fastrtps/rtps/history"
	"dds/fastrtps/rtps/reader"
	"dds/fastrtps/rtps/resources"
	"dds/fastrtps/rtps/writer"
	"dds/fastrtps/utils"
)

type IProtocolWithWlp interface {
	GetParticipant() protocol.IParticipant
	GetInitialPeers() *common.LocatorList
	GetMetatrafficUnicastLocators() *common.LocatorList
}

var _ endpoint.IWlp = (*WLP)(nil)

// Class WLP that implements the Writer Liveliness Protocol described in the RTPS specification.
type WLP struct {
	// Minimum time among liveliness periods of automatic writers, in milliseconds
	minAutomaticMs float64
	// Minimum time among liveliness periods of manual by participant writers, in milliseconds
	minManualByParticipantMs     float64
	builtinProtocols             IProtocolWithWlp
	pParticipant                 protocol.IParticipant
	builtinWriter                *writer.StatefulWriter
	builtinReader                *reader.StatefulReader
	builtinWriterHistory         *history.WriterHistory
	builtinReaderHistory         *history.ReaderHistory
	listener                     *WLPListener
	automaticLivelinessAssertion *resources.TimedEvent
	manualLivelinessAssertion    *resources.TimedEvent
	automaticWriters             []writer.IRTPSWriter
	// List of the writers using manual by participant liveliness.
	manualByParticipantWriters []writer.IRTPSWriter
	manualByTopicWriters       []writer.IRTPSWriter
	readers                    []reader.IRTPSReader
	// A boolean indicating that there is at least one reader requesting automatic liveliness
	automaticReaders bool
	// A class used by writers in this participant to keep track of their liveliness
	pubLivelinessManager              *writer.LivelinessManager
	subLivelinessManager              *writer.LivelinessManager
	automaticInstanceHandle           common.InstanceHandleT
	manualByParticipantInstanceHandle common.InstanceHandleT
	tempDataLock                      sync.Mutex
	tempReaderProxyData               *data.ReaderProxyData
	tempWriterProxyData               *data.WriterProxyData
	payloadPool                       history.ITopicPayloadPool
}

func (wlp *WLP) AddWriter(guid *common.GUIDT, kind policy.LivelinessQosPolicyKind, leaseDuration *common.DurationT) bool {
	log.Fatalln("notImpl")
	return false
}

func (wlp *WLP) RemoveWriter(guid *common.GUIDT, kind policy.LivelinessQosPolicyKind, leaseDuration *common.DurationT) bool {
	log.Fatalln("notImpl")
	return false
}

func (wlp *WLP) AssertLiveliness(writerGuid *common.GUIDT, kind policy.LivelinessQosPolicyKind, leaseDuration *common.DurationT) bool {
	log.Fatalln("notImpl")
	return false
}

func (wlp *WLP) setBuiltinWriterHistoryAttributes(hatt *attributes.HistoryAttributes, isSecure bool) {
	if isSecure {
		hatt.PayloadMaxSize = 128
	} else {
		hatt.PayloadMaxSize = 28
	}
	hatt.InitialReservedCaches = 2
	hatt.MaximumReservedCaches = 2
}

func (wlp *WLP) SubLivelinessManager() endpoint.ILivelinessManager {
	return wlp.subLivelinessManager
}

func (wlp *WLP) setBuiltinReaderHistoryAttributes(hatt *attributes.HistoryAttributes,
	cfg *utils.ResourceLimitedContainerConfig, isSecure bool) {
	cUpperLimit := uint32(math.MaxUint32 / 2)
	if isSecure {
		hatt.PayloadMaxSize = 128
	} else {
		hatt.PayloadMaxSize = 28
	}

	if (uint32)(cfg.Maximum) < cUpperLimit && (uint32)(cfg.Initial) < cUpperLimit {
		hatt.InitialReservedCaches = (uint32)(cfg.Initial * 2)
		hatt.MaximumReservedCaches = (uint32)(cfg.Maximum * 2)
	} else {
		hatt.InitialReservedCaches = (uint32)(cfg.Initial * 2)
		hatt.MaximumReservedCaches = 0
	}
}

func (wlp *WLP) createEndpoints() bool {
	participantAlloc := wlp.pParticipant.GetAttributes().Allocation.Participants
	// Builtin writer history
	hatt := attributes.KDefaultHistoryAttributes
	wlp.setBuiltinWriterHistoryAttributes(&hatt, false)
	wlp.builtinWriterHistory = history.NewWriterHistory(&hatt)
	writerPoolCfg := history.FromHistoryAttributes(&hatt)
	wlp.payloadPool = history.GetTopicPayloadPoolProxy("DCPSParticipantMessage", writerPoolCfg)
	wlp.payloadPool.ReserveHistory(writerPoolCfg, false)
	// Built-in writer
	watt := attributes.NewWriterAttributes()
	watt.EndpointAtt.UnicastLocatorList.Append(wlp.builtinProtocols.GetInitialPeers())
	watt.MatchedReadersAllocation = *participantAlloc
	watt.EndpointAtt.TopicKind = common.KWithKey
	watt.EndpointAtt.DurabilityKind = common.KTransientLocal
	watt.EndpointAtt.ReliabilityKind = common.KReliable
	throughputController := wlp.pParticipant.GetAttributes().ThroughputController
	if throughputController.BytesPerPeriod != math.MaxUint32 && throughputController.PeriodMillisecs != 0 {
		watt.PubMode = attributes.KAsynchronousWriter
	}

	ok, wout := wlp.pParticipant.CreateWriter(watt, wlp.payloadPool, wlp.builtinWriterHistory, nil,
		common.KEntityIDWriterLiveliness, true)
	if ok {
		if wlp.builtinWriter, ok = wout.(*writer.StatefulWriter); ok {
			log.Println("Builtin Liveliness Writer created")
		} else {
			log.Fatalln("Builtin Liveliness Writer created failed !")
		}
	} else {
		wlp.payloadPool.ReleaseHistory(writerPoolCfg, false)
		log.Fatalln("Liveliness Writer Creation failed ")
		return false
	}
	// Built-in reader history
	wlp.setBuiltinReaderHistoryAttributes(&hatt, participantAlloc, false)
	wlp.builtinReaderHistory = history.NewReaderHistory(&hatt)
	readerPoolCfg := history.FromHistoryAttributes(&hatt)
	wlp.payloadPool.ReserveHistory(readerPoolCfg, true)

	// WLP listener
	wlp.listener = NewWlpListener(wlp)

	// Built-in reader
	ratt := attributes.NewReaderAttributes()
	ratt.EndpointAtt.TopicKind = common.KWithKey
	ratt.EndpointAtt.DurabilityKind = common.KTransientLocal
	ratt.EndpointAtt.ReliabilityKind = common.KReliable
	ratt.ExpectsInlineQos = true
	ratt.EndpointAtt.UnicastLocatorList.Append(wlp.builtinProtocols.GetMetatrafficUnicastLocators())
	ratt.EndpointAtt.MulticastLocatorList.Append(wlp.builtinProtocols.GetMetatrafficUnicastLocators())
	ratt.EndpointAtt.RemoteLocatorList.Append(wlp.builtinProtocols.GetInitialPeers())
	ratt.MatchedWritersAllocation = *participantAlloc
	ratt.EndpointAtt.TopicKind = common.KWithKey

	rout, ok := wlp.pParticipant.CreateReader(ratt, wlp.payloadPool, wlp.builtinReaderHistory,
		wlp.listener, common.KEntityIDReaderLiveliness, true, true)
	if ok {
		wlp.builtinReader, ok = rout.(*reader.StatefulReader)
		if ok {
			log.Println("Builtin Liveliness Reader created")
		} else {
			log.Fatalln("Builtin Liveliness Reader created failed")
		}
	} else {
		log.Fatalln("Liveliness Reader Creation failed.")
		wlp.payloadPool.ReleaseHistory(readerPoolCfg, true)
		return false
	}

	return true
}

func (wlp *WLP) InitWL(p interface{}) bool {
	aprotocol, ok := p.(protocol.IParticipant)
	if !ok {
		log.Fatalln("initWL failed with invalid protocol")
	}
	log.Println("Initializing Liveliness Protocol")
	wlp.pParticipant = aprotocol
	pubCallback := func(guid *common.GUIDT, akind policy.LivelinessQosPolicyKind, leaseDuration *common.DurationT, aliveCount int32, notAliveCount int32) {
		log.Fatalln("notImpl")
	}
	wlp.pubLivelinessManager = writer.NewLivelinessManager(&pubCallback, wlp.pParticipant.GetEventResource(), false)
	subCallback := func(guid *common.GUIDT, akind policy.LivelinessQosPolicyKind, leaseDuration *common.DurationT, aliveCount int32, notAliveCount int32) {
		log.Fatalln("notImpl")
	}
	wlp.subLivelinessManager = writer.NewLivelinessManager(&subCallback, wlp.pParticipant.GetEventResource(), true)
	return wlp.createEndpoints()
}

func NewWLP(parent IProtocolWithWlp) *WLP {
	var wlp WLP
	wlp.minAutomaticMs = math.MaxFloat64
	wlp.minManualByParticipantMs = math.MaxFloat64
	wlp.builtinProtocols = parent
	wlp.automaticReaders = false
	allocAtt := wlp.builtinProtocols.GetParticipant().GetAttributes().Allocation
	wlp.tempReaderProxyData = data.NewReaderProxyData(allocAtt.Locators.MaxUnicastLocators,
		allocAtt.Locators.MaxMulticastLocators, allocAtt.DataLimits)
	wlp.tempWriterProxyData = data.NewWriterProxyData(allocAtt.Locators.MaxUnicastLocators, allocAtt.Locators.MaxMulticastLocators, allocAtt.DataLimits)
	wlp.automaticInstanceHandle = common.CreateInstanceHandle(wlp.builtinProtocols.GetParticipant().GetGuid())
	copy(wlp.automaticInstanceHandle.Value[:3], []common.Octet{0, 0, 0})
	wlp.manualByParticipantInstanceHandle = wlp.automaticInstanceHandle
	wlp.automaticInstanceHandle.Value[15] = policy.AUTOMATIC_LIVELINESS_QOS + 0x01
	wlp.manualByParticipantInstanceHandle.Value[15] = policy.MANUAL_BY_PARTICIPANT_LIVELINESS_QOS + 0x01
	return &wlp
}
