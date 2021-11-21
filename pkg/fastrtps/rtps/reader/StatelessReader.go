package reader

import (
	"log"
	"math"

	"dds/common"
	"dds/core/policy"
	"dds/fastrtps/rtps/attributes"
	"dds/fastrtps/rtps/builtin/data"
	"dds/fastrtps/rtps/endpoint"
	"dds/fastrtps/rtps/history"

	"github.com/golang/glog"
)

//var _ IRTPSReader = (*StatelessReader)(nil)
var _ ireaderImpl = (*StatelessReader)(nil)

// Class StatelessReader, specialization of the RTPSReader for Best Effort Readers.
type StatelessReader struct {
	RTPSReader
	matchedWriters    []remoteWriterInfoT
	maxMatchedWriters int
}

type remoteWriterInfoT struct {
	GUID                     common.GUIDT
	PersistenceGUID          common.GUIDT
	HasManualTopicLiveliness bool
	FragmentedChange         *common.CacheChangeT
}

func (statelessReader *StatelessReader) ChangeRemovedByHistory(change *common.CacheChangeT, prox history.IWriterProxyWithHistory) bool {
	if !change.IsRead {
		if statelessReader.totalUnread > 0 {
			statelessReader.totalUnread--
		}
	}
	return true
}

func (statelessReader *StatelessReader) acceptMsgFrom(writerID *common.GUIDT, changeKind common.ChangeKindT) bool {
	if changeKind == common.KAlive {
		if statelessReader.acceptMessageFromUnKnowWriters {
			return true
		} else if writerID.EntityID == statelessReader.trustedWriterEntityID {
			return true
		}
	}
	for i := 0; i < len(statelessReader.matchedWriters); i++ {
		if statelessReader.matchedWriters[i].GUID == *writerID {
			return true
		}
	}
	return false
}

func (statelessReader *StatelessReader) MatchedWriterAdd(wdata *data.WriterProxyData) bool {
	statelessReader.Mutex.Lock()
	defer statelessReader.Mutex.Unlock()
	for i := 0; i < len(statelessReader.matchedWriters); i++ {
		if statelessReader.matchedWriters[i].GUID == *wdata.Guid() {
			glog.Info("Attempting to add existing writer")
			return false
		}
	}

	info := remoteWriterInfoT{
		GUID:            *wdata.Guid(),
		PersistenceGUID: *wdata.PersistentGuid(),
	}

	if wdata.Qos.Liveliness.Kind == policy.MANUAL_BY_TOPIC_LIVELINESS_QOS {
		info.HasManualTopicLiveliness = true
	}

	if len(statelessReader.matchedWriters) < statelessReader.maxMatchedWriters {
		glog.Error("Finite liveliness lease duration but WLP not enabled")
		log.Println("No space to add writer ", *wdata.Guid(), " to reader ", statelessReader.GUID)
		return false
	}
	statelessReader.matchedWriters = append(statelessReader.matchedWriters, info)

	//statelessReader.AddPersistenceGuid(info.GUID, info.PersistenceGUID)

	persistence_guid_to_store := info.PersistenceGUID
	if persistence_guid_to_store == common.KGuidUnknown {
		persistence_guid_to_store = info.GUID
	}
	statelessReader.historyState.PersistenceGUIDMap[info.GUID] = persistence_guid_to_store
	statelessReader.historyState.PersistenceGUIDCount[persistence_guid_to_store]++

	statelessReader.acceptMessageFromUnKnowWriters = false
	log.Println("Writer ", info.GUID, " add to reader ", statelessReader.GUID)

	if statelessReader.livelinessLeaseDuration.Less(common.KTimeInfinite) {
		wlp := statelessReader.RTPSParticipant.Wlp()
		if wlp == nil {
			glog.Fatalln("Finite liveliness lease duration but WLP not enabled")
		}

		wlp.SubLivelinessManager().AddWriter(*wdata.Guid(), statelessReader.livelinessKind,
			statelessReader.livelinessLeaseDuration)
	}

	return true
}

func (statelessReader *StatelessReader) ProcessDataMsg(change *common.CacheChangeT) bool {

	statelessReader.Mutex.Lock()
	defer statelessReader.Mutex.Unlock()

	if !statelessReader.acceptMsgFrom(&change.WriterGUID, change.Kind) {
		log.Println("refuse to accept msg")
		return true
	}

	log.Println("Trying to add change ", change.SequenceNumber, " TO reader: ", statelessReader.GUID)
	statelessReader.assertWriterLiveliness(&change.WriterGUID)

	// Ask the pool for a cache change
	changeToAdd, ok := statelessReader.ChangePool.ReserveCache()
	if !ok {
		log.Fatalln("Problem reserving CacheChange in reader: ", statelessReader.GUID)
		return false
	}

	// Copy metadata to reserved change
	changeToAdd.CopyNotMemcpy(change)

	// Ask payload pool to copy the payload
	payloadOwner := change.PayloadOwner()
	if statelessReader.PayloadPool.GetPayloadWithOwner(&change.SerializedPayload, &payloadOwner, changeToAdd) {
		change.SetPayloadOwner(payloadOwner)
	} else {
		dataSize := statelessReader.FixedPayloadSize
		if dataSize == 0 {
			dataSize = math.MaxUint32
		}
		log.Println("Problem copying CacheChange, received data is: ",
			change.SerializedPayload.Length, " bytes and max size in reader ",
			statelessReader.GUID, " is ", dataSize)
		statelessReader.ReleaseCache(changeToAdd)
		return false
	}

	// Perform reception of cache change
	if !statelessReader.ChangeReceived(changeToAdd) {
		log.Println("MessageReceiver not add change ", changeToAdd.SequenceNumber)
		statelessReader.PayloadPool.ReleasePayload(changeToAdd)
		statelessReader.ChangePool.ReleaseCache(changeToAdd)
	}

	return true
}

func (statelessReader *StatelessReader) ProcessGapMsg(writerGUID *common.GUIDT,
	gapStart *common.SequenceNumberT, gapList *common.SequenceNumberSet) bool {
	return true
}

func (statelessReader *StatelessReader) ProcessHeartbeatMsg(writerGUID *common.GUIDT,
	hbCount uint32, firstSN *common.SequenceNumberT,
	lastSN *common.SequenceNumberT, finalFlag bool, livelinessFlag bool) bool {
	return true
}

// * This method is called when a new change is received. This method calls the received_change of the History
// * and depending on the implementation performs different actions.
func (statelessReader *StatelessReader) ChangeReceived(aChange *common.CacheChangeT) bool {
	// Only make visible the change if there is not other with bigger sequence number.
	if !statelessReader.thereIsUpperRecordOf(&aChange.WriterGUID, &aChange.SequenceNumber) {
		if statelessReader.readerHistory.ReceivedChange(aChange) {
			aChange.ReceptionTimestamp = common.CurrentTime()
			statelessReader.updateLastNotified(&aChange.WriterGUID, &aChange.SequenceNumber)
			statelessReader.totalUnread++

			if statelessReader.listener != nil {
				statelessReader.listener.OnNewCacheChangeAdded(statelessReader, aChange)
			}
			statelessReader.newNotificationCV.Broadcast()
			return true
		}
	}
	return false
}

func (statelessReader *StatelessReader) getLastNotified(guid *common.GUIDT) common.SequenceNumberT {
	var retVal common.SequenceNumberT
	// TODO::
	// statelessReader.Mutex.Lock()
	// defer statelessReader.Mutex.Unlock()
	guidToLook := *guid
	pGUID, ok := statelessReader.historyState.PersistenceGUIDMap[*guid]
	if ok {
		guidToLook = pGUID
	}

	pSeq, ok := statelessReader.historyState.HistoryRecord[guidToLook]
	if ok {
		retVal = pSeq
	}

	return retVal
}

func (statelessReader *StatelessReader) thereIsUpperRecordOf(guid *common.GUIDT, seq *common.SequenceNumberT) bool {
	lastNotifiedSeq := statelessReader.getLastNotified(guid)
	return !lastNotifiedSeq.Less(seq)
}

func (statelessReader *StatelessReader) assertWriterLiveliness(guid *common.GUIDT) {
	if statelessReader.livelinessLeaseDuration.Less(common.KTimeInfinite) {
		wlp := statelessReader.RTPSParticipant.Wlp()
		if wlp != nil {
			wlp.AssertLiveliness(guid, statelessReader.livelinessKind, &statelessReader.livelinessLeaseDuration)
		} else {
			log.Fatal("Finite liveliness lease duration but WLP not enabled")
		}
	}
}

func (statelessReader *StatelessReader) ProcessDataFragMsg(aChange *common.CacheChangeT, sampleSize uint32,
	fragmentStartingNum uint32, fragmentsInSubmessage uint16) bool {
	writerGUID := aChange.WriterGUID
	statelessReader.Mutex.Lock()
	defer statelessReader.Mutex.Unlock()

	for _, loopWriter := range statelessReader.matchedWriters {
		if loopWriter.GUID != writerGUID {
			continue
		}

		statelessReader.assertWriterLiveliness(&writerGUID)

		// Check if CacheChange was received.
		if !statelessReader.thereIsUpperRecordOf(&writerGUID, &aChange.SequenceNumber) {
			log.Println("Trying to add fragment ", aChange.SequenceNumber.Value, " TO reader: ", statelessReader.GUID)

			// Early return if we already know abount a greater sequence number
			workChange := loopWriter.FragmentedChange
			if workChange != nil && aChange.SequenceNumber.Less(&workChange.SequenceNumber) {
				return true
			}

			changeToAdd := aChange
			// Check if pending fragmented change should be dropped
			if workChange != nil {
				if workChange.SequenceNumber.Less(&changeToAdd.SequenceNumber) {
					// Pending change should be dropped. Check if it can be reused
					if sampleSize <= workChange.SerializedPayload.MaxSize {
						// Sample fits inside pending change. Reuse it.
						workChange.CopyNotMemcpy(changeToAdd)
						workChange.SerializedPayload.Length = sampleSize
						workChange.SetFragmentSize(changeToAdd.GetFragmentSize(), true)
					} else {
						// Release change, and let it be reserved later
						statelessReader.ReleaseCache(workChange)
						workChange = nil
					}
				}
			}

			// Check if a new change should be reserved
			if workChange == nil {
				if workChange, ok := statelessReader.ReserveCache(sampleSize); ok {
					if workChange.SerializedPayload.MaxSize < sampleSize {
						statelessReader.ReleaseCache(workChange)
						workChange = nil
					} else {
						workChange.CopyNotMemcpy(changeToAdd)
						workChange.SerializedPayload.Length = sampleSize
						workChange.SetFragmentSize(changeToAdd.GetFragmentSize(), true)
					}
				}
			}

			// Process fragment and set change_completed if it is fully reassembled
			var changeCompleted *common.CacheChangeT
			if workChange != nil {
				if workChange.AddFragments(&changeToAdd.SerializedPayload, fragmentStartingNum, fragmentsInSubmessage) {
					changeCompleted = workChange
					workChange = nil
				}
			}
			loopWriter.FragmentedChange = workChange

			// If the change was completed, process it.
			if changeCompleted != nil {
				if !statelessReader.ChangeReceived(changeCompleted) {
					log.Println("MessageReceiver not add change ", changeCompleted.SequenceNumber.Value)

					// Release CacheChange_t
					statelessReader.ReleaseCache(changeCompleted)
				}
			}
		}
		return true
	}

	log.Println("Reader ", statelessReader.GUID, " received DATA_FRAG from unknown writer", writerGUID)
	return true
}

// func (statelessReader *StatelessReader) init(payloadPool history.IPayloadPool, changePool history.IChangePool) {
// 	statelessReader.PayloadPool = payloadPool
// 	statelessReader.ChangePool = changePool
// 	statelessReader.FixedPayloadSize = 0
// 	if statelessReader.readerHistory.Att.MemoryPolicy == resources.KPreallocatedMemoryMode {
// 		statelessReader.FixedPayloadSize = statelessReader.readerHistory.Att.PayloadMaxSize
// 	}

// 	statelessReader.readerHistory.Reader = statelessReader
// 	statelessReader.readerHistory.Mutex = &statelessReader.Mutex

// 	log.Println("RTPSReader created correctly")
// }

func NewStatelessReader(parent endpoint.IEndpointParent, guid *common.GUIDT, att *attributes.ReaderAttributes,
	payloadPool history.IPayloadPool, hist *history.ReaderHistory, rlisten IReaderListener) *StatelessReader {
	var retReader StatelessReader
	poolCfg := history.FromHistoryAttributes(&hist.Att)
	aChangePool := history.NewCacheChangePool(poolCfg)
	retReader.RTPSReader = *NewRtpsReader(parent, guid, att, payloadPool, aChangePool, hist, rlisten)
	retReader.RTPSReader.impl = &retReader
	retReader.matchedWriters = make([]remoteWriterInfoT, att.MatchedWritersAllocation.Initial)
	retReader.init(&retReader, payloadPool, aChangePool)
	return &retReader
}
