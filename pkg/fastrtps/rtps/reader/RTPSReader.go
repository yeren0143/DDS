package reader

import (
	"log"
	"sync"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/core/policy"
	"github.com/yeren0143/DDS/fastrtps/message"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/endpoint"
	"github.com/yeren0143/DDS/fastrtps/rtps/history"
	"github.com/yeren0143/DDS/fastrtps/utils"
)

type IReaderHistory interface {
}

var _ history.IReaderWithHistory = (IRTPSReader)(nil)

//var _ endpoint.IEndpoint = (IRTPSReader)(nil)

// IRtpsReader manages the reception of data from its matched writers.
type IRTPSReader interface {
	endpoint.IEndpoint
	message.IRtpsMsgReader

	// Processes a new DATA message. Previously the message must have been accepted by
	// function acceptMsgDirectedTo.
	ProcessDataMsg(change *common.CacheChangeT) bool

	ProcessDataFragMsg(change *common.CacheChangeT,
		sampleSize uint32,
		fragmentStartingNum uint32,
		fragmentsInSubmessage uint16) bool

	ProcessHeartbeatMsg(writerGUID *common.GUIDT,
		hbCount uint32,
		firstSN *common.SequenceNumberT,
		lastSN *common.SequenceNumberT,
		finalFlag bool,
		livelinessFlag bool) bool

	ProcessGapMsg(writerGUID *common.GUIDT,
		gapStart *common.SequenceNumberT,
		gapList *common.SequenceNumberSet) bool

	// Method to indicate the reader that some change has been removed due to
	// HistoryQos requirements.
	ChangeRemovedByHistory(change *common.CacheChangeT, prox history.IWriterProxyWithHistory) bool

	// Accept msg to unknwon readers (default=true)
	AcceptMessagesToUnknownReaders() bool

	GetAttributes() *attributes.EndpointAttributes

	SetTrustedWriter(writerEnt *common.EntityIDT)

	GetGUID() *common.GUIDT

	GetMutex() *sync.Mutex

	ReleaseCache(change *common.CacheChangeT)
	ReserveCache(size uint32) (*common.CacheChangeT, bool)

	ExpectsInlineQos() bool
}

// reader who devired from IRtpsReader must implement ireaderImpl
type ireaderImpl interface {
	mayRemoveHistoryRecord(removedByLease bool) bool
	setLastNotified(persistenceGUID *common.GUIDT, seq *common.SequenceNumberT)
	init(payloadPool history.IPayloadPool, changePool history.IChangePool)
}

type rtpsReaderBase struct {
	endpoint.EndpointBase
	impl                           ireaderImpl
	trustedWriterEntityID          common.EntityIDT
	acceptMessagesToUnknownReaders bool // Accept msg to unknwon readers (default=true)
	acceptMessageFromUnKnowWriters bool
	expectsInlineQos               bool
	totalUnread                    uint64
	listener                       IReaderListener
	newNotificationCV              *utils.TimedConditionVariable
	readerHistory                  *history.ReaderHistory
	historyState                   *ReaderHistoryState
	livelinessKind                 policy.LivelinessQosPolicyKind
	livelinessLeaseDuration        common.DurationT
}

func (reader *rtpsReaderBase) AcceptMessagesToUnknownReaders() bool {
	return reader.acceptMessagesToUnknownReaders
}

func (reader *rtpsReaderBase) ExpectsInlineQos() bool {
	return reader.expectsInlineQos
}

func (reader *rtpsReaderBase) mayRemoveHistoryRecord(removedByLease bool) bool {
	return !removedByLease
}

func (reader *rtpsReaderBase) setLastNotified(persistenceGUID *common.GUIDT, seq *common.SequenceNumberT) {
	reader.historyState.HistoryRecord[*persistenceGUID] = *seq
}

// func (reader *rtpsReaderBase) GetAttributes() *attributes.EndpointAttributes {
// 	return &reader.Att
// }

// func (reader *rtpsReaderBase) GetMutex() *sync.Mutex {
// 	return reader.GetMutex()
// }

// func (reader *rtpsReaderBase) GetGUID() *common.GUIDT {
// 	return &reader.GUID
// }

func (reader *rtpsReaderBase) SetTrustedWriter(writerEnt *common.EntityIDT) {
	reader.acceptMessagesToUnknownReaders = false
	reader.trustedWriterEntityID = *writerEnt
}

func (reader *rtpsReaderBase) ReleaseCache(change *common.CacheChangeT) {
	reader.Mutex.Lock()
	defer reader.Mutex.Unlock()

	if pool := change.PayloadOwner(); pool != nil {
		pool.ReleasePayload(change)
	}
	reader.ChangePool.ReleaseCache(change)
}

// Update the last notified sequence for a RTPS guid
func (reader *rtpsReaderBase) updateLastNotified(guid *common.GUIDT, seq *common.SequenceNumberT) common.SequenceNumberT {
	var retVal common.SequenceNumberT
	reader.Mutex.Lock()
	defer reader.Mutex.Unlock()
	guidToLook := *guid
	pguid, ok := reader.historyState.PersistenceGUIDMap[*guid]
	if ok {
		guidToLook = pguid
	}

	pseq, ok := reader.historyState.HistoryRecord[guidToLook]
	if ok {
		retVal = pseq
	}

	if retVal.Less(seq) {
		reader.impl.setLastNotified(&guidToLook, seq)
		reader.newNotificationCV.Broadcast()
	}

	return retVal
}

func (reader *rtpsReaderBase) ReserveCache(dataCdrSerializedSize uint32) (*common.CacheChangeT, bool) {
	reader.Mutex.Lock()
	reader.Mutex.Unlock()

	reservedChange, ok := reader.ChangePool.ReserveCache()
	if !ok {
		log.Fatalln("Problem reserving cache from pool")
		return nil, false
	}

	payloadSize := reader.FixedPayloadSize
	if payloadSize < 0 {
		payloadSize = dataCdrSerializedSize
	}
	if !reader.PayloadPool.GetPayload(payloadSize, reservedChange) {
		reader.ChangePool.ReleaseCache(reservedChange)
		log.Fatalln("Problem reserving payload from pool")
		return nil, false
	}

	return reservedChange, true
}

func NewRtpsReaderBase(parent endpoint.IEndpointParent, guid *common.GUIDT, att *attributes.ReaderAttributes,
	payloadPool history.IPayloadPool, changePool history.IChangePool, hist *history.ReaderHistory,
	rlisten IReaderListener) *rtpsReaderBase {

	var retReader rtpsReaderBase
	retReader.EndpointBase = *endpoint.NewEndPointBase(parent, guid, &att.EndpointAtt)
	retReader.readerHistory = hist
	retReader.listener = rlisten
	retReader.acceptMessagesToUnknownReaders = true
	retReader.acceptMessageFromUnKnowWriters = false
	retReader.expectsInlineQos = att.ExpectsInlineQos
	retReader.historyState = NewReaderHistoryState(att.MatchedWritersAllocation.Initial)
	retReader.livelinessKind = att.LivelinessKind
	retReader.livelinessLeaseDuration = att.LivelinessLeaseDuration
	//retReader.impl.init(payloadPool, changePool)
	return &retReader
}
