package writer

import (
	"log"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/core/policy"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/endpoint"
	"github.com/yeren0143/DDS/fastrtps/rtps/flowcontrol"
	"github.com/yeren0143/DDS/fastrtps/rtps/history"
)

var _ history.IWriterWithHistory = (IRTPSWriter)(nil)

// type IWriterParent interface {
// 	GetAttributes() *attributes.RTPSParticipantAttributes
// 	Wlp() endpoint.IWlp
// 	SendSync(msg *common.CDRMessage, locators []common.Locator, maxBlockingTimePoint common.Time) bool
// }

type writerCallback = func() uint32

// RTPSWriter manages the sending of data to the readers. Is always associated with a HistoryCache.
type IRTPSWriter interface {
	endpoint.IEndpoint
	IRTPSMessageSender

	// Process an incoming ACKNACK submessage.
	// result true if the writer could process the submessage. Only valid when returned value is true.
	// valid true when the submessage was destinated to this writer, false otherwise.
	ProcessAcknack(writerGUID, readerGUID *common.GUIDT, ackCount uint32,
		snSet *common.SequenceNumberSet, finalFlag bool) (result, valid bool)

	// result true if the writer could process the submessage. Only valid when returned value is true
	// true when the submessage was destinated to this writer, false otherwise.
	ProcessNackFrag(writerGUID, readerGUID *common.GUIDT, ackCount uint32,
		seqNum *common.SequenceNumberT, fragmentsState *common.FragmentNumberSet) (result, valid bool)

	GetAttributes() *attributes.EndpointAttributes

	GetGUID() *common.GUIDT

	NewChange(dataCdrSerializedSize writerCallback, changeKind common.ChangeKindT,
		handle common.InstanceHandleT) *common.CacheChangeT

	ReleaseChange(aChange *common.CacheChangeT) bool

	UnsentChangeAddedToHistory(aChange *common.CacheChangeT, maxBlockingTime common.Time)

	AddFlowController(controller flowcontrol.IFlowController)
}

type rtpsWriterBase struct {
	endpoint.EndpointBase
	// Is the data sent directly or announced by HB and THEN sent to the ones who ask for it?.
	pushMode bool
	hist     *history.WriterHistory
	listen   IWriterListener
	// Asynchronous publication activated
	isAsync bool
	// Separate sending activated
	separateSendingEnabled       bool
	allRemoteReaders             []common.GUIDT
	allRemoteParticipants        []common.GUIDPrefixT
	locatorSelector              common.LocatorSelector
	livelinessKind               policy.LivelinessQosPolicyKind
	livelinessLeaseDuration      common.DurationT
	livelinessAnnouncementPeriod common.DurationT
}

func (writerBase *rtpsWriterBase) DestinationGuidPrefix() common.GUIDPrefixT {
	if len(writerBase.allRemoteParticipants) == 1 {
		return writerBase.allRemoteParticipants[0]
	} else {
		return common.KGuidPrefixUnknown
	}
}

// Check if the destinations managed by this sender interface have changed.
func (writerBase *rtpsWriterBase) DestinationHaveChanged() bool {
	return false
}

func (writerBase *rtpsWriterBase) ProcessAcknack(writerGUID, readerGUID *common.GUIDT,
	ackCount uint32, snSet *common.SequenceNumberSet, finalFlag bool) (result, valid bool) {
	return result, *writerGUID == writerBase.GUID
}

func (writerBase *rtpsWriterBase) ProcessNackFrag(writerGUID, readerGUID *common.GUIDT,
	ackCount uint32, seqNum *common.SequenceNumberT,
	fragmentsState *common.FragmentNumberSet) (result, valid bool) {
	return false, *writerGUID == writerBase.GUID
}

func (writerBase *rtpsWriterBase) ReleaseChange(aChange *common.CacheChangeT) bool {
	writerBase.Mutex.Lock()
	defer writerBase.Mutex.Unlock()
	owner := aChange.PayloadOwner()
	if owner != nil {
		owner.ReleasePayload(aChange)
	}
	return writerBase.ChangePool.ReleaseCache(aChange)
}

func (writerBase *rtpsWriterBase) RemoteGUIDs() []common.GUIDT {
	return writerBase.allRemoteReaders
}

func (writerBase *rtpsWriterBase) Send(msg *common.CDRMessage, maxBlockingTimePoint common.Time) bool {
	// participant := writerBase.RTPSParticipant
	// return writerBase.locatorSelector.SelectedSize() == 0 ||
	// 	participant.SendSync(msg, maxBlockingTimePoint)
	log.Panic("not impl")
	return false
}

func (writerBase *rtpsWriterBase) RemoteParticipants() []common.GUIDPrefixT {
	return writerBase.allRemoteParticipants
}

func (writerBase *rtpsWriterBase) NewChange(dataCdrSerializedSize writerCallback,
	changeKind common.ChangeKindT, handle common.InstanceHandleT) *common.CacheChangeT {
	log.Println("Creating new change")
	writerBase.Mutex.Lock()
	writerBase.Mutex.Unlock()
	reservedChange, ok := writerBase.ChangePool.ReserveCache()
	if !ok {
		log.Println("Problem reserving cache from pool")
		return nil
	}

	payloadSize := writerBase.FixedPayloadSize
	if payloadSize == 0 {
		payloadSize = dataCdrSerializedSize()
	}

	if !writerBase.PayloadPool.GetPayload(payloadSize, reservedChange) {
		writerBase.ChangePool.ReleaseCache(reservedChange)
		log.Println("Problem reserving payload from pool")
		return nil
	}

	reservedChange.Kind = changeKind
	if writerBase.Att.TopicKind == common.KWithKey && !handle.IsDefined() {
		log.Println("Changes in KEYED Writers need a valid instanceHandle")
	}

	reservedChange.InstanceHandle = handle
	reservedChange.WriterGUID = writerBase.GUID
	return reservedChange
}

func (writerBase *rtpsWriterBase) init(payloadPool history.IPayloadPool, changePool history.IChangePool) {
	writerBase.PayloadPool = payloadPool
	log.Println("RTPSWriter created")
}

func newRtpsWriterBase(parent endpoint.IEndpointParent, guid *common.GUIDT,
	att *attributes.WriterAttributes, payloadPool history.IPayloadPool, changePool history.IChangePool,
	hist *history.WriterHistory, wlisten IWriterListener) *rtpsWriterBase {
	var retWriter rtpsWriterBase
	retWriter.EndpointBase = *endpoint.NewEndPointBase(parent, guid, &att.EndpointAtt)
	retWriter.hist = hist
	retWriter.listen = wlisten
	retWriter.isAsync = (att.PubMode == attributes.KAsynchronousWriter)
	// TODO:
	//writerBase.locatorSelector = att.MatchedReadersAllocation
	retWriter.allRemoteReaders = make([]common.GUIDT, att.MatchedReadersAllocation.Initial)
	retWriter.allRemoteParticipants = make([]common.GUIDPrefixT, att.MatchedReadersAllocation.Initial)
	retWriter.livelinessKind = att.LivelinessKind
	retWriter.livelinessLeaseDuration = att.LivelinessLeaseDuration
	retWriter.livelinessAnnouncementPeriod = att.LivelinessAnnouncementPeriod
	retWriter.init(payloadPool, changePool)
	return &retWriter
}
