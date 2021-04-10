package writer

import (
	"log"
	"sync"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/endpoint"
	"github.com/yeren0143/DDS/fastrtps/rtps/flowcontrol"
	"github.com/yeren0143/DDS/fastrtps/rtps/history"
	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
)

var _ IRTPSWriter = (*StatefulWriter)(nil)

// Class StatefulWriter, specialization of RTPSWriter that maintains information of each matched Reader.
type StatefulWriter struct {
	rtpsWriterBase

	// Timed Event to manage the periodic HB to the Reader.
	periodHbEvent *resources.TimedEvent
	// Timed Event to manage the Acknack response delay.
	nackResponseEvent *resources.TimedEvent
	// A timed event to mark samples as acknowledget (used only if disable positive ACKs QoS is enabled)
	ackEvent *resources.TimedEvent

	heartbeatCount uint32
	times          attributes.WriterTimes

	// Vector containing all the active ReaderProxies.
	matchedReaders []*ReaderProxy
	// Vector containing all the inactive, ready for reuse, ReaderProxies.
	matchedReadersPool []*ReaderProxy

	allAckedMutex             *sync.Mutex
	allAckedCond              *sync.Cond
	allAcked                  bool
	mayRemoveChangeCond       *sync.Cond
	mayRemoveChange           uint32
	disbaleHeartbeatPiggyback bool
	disablePositiveAcks       bool
	// Keep duration for disable positive ACKs QoS, in microseconds
	keepDuration common.DurationT
	// To avoid notifying twice of the same sequence number
	nextAllAckedNotifySequence common.SequenceNumberT
	minReadersLowMark          common.SequenceNumberT

	sendBufferSize             uint32
	currentUsageSendBufferSize uint32
	flowControllers            []flowcontrol.IFlowController

	threreAreRemoteReaders bool
	thereAreLocalReaders   bool
	readerDataFilter       IReaderDataFilter
}

func (statefulWriter *StatefulWriter) AddFlowController(controller flowcontrol.IFlowController) {
	statefulWriter.flowControllers = append(statefulWriter.flowControllers, controller)
}

func (statefulWriter *StatefulWriter) SendAnyUnsentChanges() {
	log.Fatalln("notImpl")
}

func (statefulWriter *StatefulWriter) SendPeriodHeartbeat(final, liveliness bool) bool {
	log.Fatalln("notImpl")
	return false
}

func (statefulWriter *StatefulWriter) UnsentChangeAddedToHistory(aChange *common.CacheChangeT,
	maxBlockingTime common.Time) {
	log.Fatalln("notImpl")
}

func (statefulWriter *StatefulWriter) PerformNackResponse() {
	log.Fatalln("notImpl")
}

func (statefulWriter *StatefulWriter) AckTimerExpired() bool {
	log.Fatalln("notImpl")
	return false
}

func (statefulWriter *StatefulWriter) init(parent endpoint.IEndpointParent, att *attributes.WriterAttributes) {
	partAtt := parent.GetAttributes()
	periodHbEventCallback := func() bool {
		return statefulWriter.SendPeriodHeartbeat(false, false)
	}
	statefulWriter.periodHbEvent = resources.NewTimedEvent(parent.GetEventResource(),
		&periodHbEventCallback, statefulWriter.times.HeartbeatPeriod.MilliSeconds())

	nackResponseEventCallback := func() bool {
		statefulWriter.PerformNackResponse()
		return false
	}
	statefulWriter.nackResponseEvent = resources.NewTimedEvent(parent.GetEventResource(),
		&nackResponseEventCallback, statefulWriter.times.NackResponseDelay.MilliSeconds())

	if statefulWriter.disablePositiveAcks {
		ackEventCallback := func() bool {
			return statefulWriter.AckTimerExpired()
		}
		statefulWriter.ackEvent = resources.NewTimedEvent(parent.GetEventResource(),
			&ackEventCallback, att.KeepDuration.MilliSeconds())
	}

	for n := uint32(0); n < att.MatchedReadersAllocation.Initial; n++ {
		readerProxy := NewReaderProxy(&statefulWriter.times, partAtt.Allocation.Locators, statefulWriter)
		statefulWriter.matchedReadersPool = append(statefulWriter.matchedReadersPool, readerProxy)
	}
}

func NewStatefulWriter(parent endpoint.IEndpointParent, guid *common.GUIDT, att *attributes.WriterAttributes,
	payloadPool history.IPayloadPool, hist *history.WriterHistory, listener IWriterListener) *StatefulWriter {
	var awriter StatefulWriter
	poolCfg := history.FromHistoryAttributes(&hist.Att)
	cacheChangePool := history.NewCacheChangePool(poolCfg)
	awriter.rtpsWriterBase = *newRtpsWriterBase(parent, guid, att, payloadPool, cacheChangePool, hist, listener)
	awriter.heartbeatCount = 0
	awriter.times = att.Times
	awriter.matchedReaders = make([]*ReaderProxy, att.MatchedReadersAllocation.Initial)
	awriter.matchedReadersPool = make([]*ReaderProxy, att.MatchedReadersAllocation.Initial)
	awriter.nextAllAckedNotifySequence.Value = 1
	awriter.allAcked = false
	awriter.mayRemoveChange = 0
	awriter.allAckedMutex = new(sync.Mutex)
	awriter.disbaleHeartbeatPiggyback = att.DisableHeartbeatPiggyback
	awriter.disablePositiveAcks = att.DisablePositiveAcks
	awriter.keepDuration = att.KeepDuration
	awriter.sendBufferSize = parent.GetMinNetworkSendBufferSize()
	awriter.currentUsageSendBufferSize = parent.GetMinNetworkSendBufferSize()

	awriter.hist.Writer = &awriter
	awriter.hist.Mutex = awriter.Mutex
	awriter.init(parent, att)
	return &awriter
}
