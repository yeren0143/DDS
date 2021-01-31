package history

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
	"log"
)

var _ IHistory = (*ReaderHistory)(nil)

type ihistoryWriterProxy interface {
}

type IHistoryReader interface {
	ChangeRemovedByHistory(change *common.CacheChangeT, proxy ihistoryWriterProxy) bool
	ReleaseCache(change *common.CacheChangeT)
	ReserveCache(changes []*common.CacheChangeT) bool
}

type ReaderHistory struct {
	historyBase
	reader IHistoryReader
}

func NewReaderHistory(att *attributes.HistoryAttributes) *ReaderHistory {
	return &ReaderHistory{
		historyBase: *NewhistoryBase(att),
	}
}

func (history *ReaderHistory) ReceivedChange(change *common.CacheChangeT) bool {
	return history.addChange(change)
}

func (history *ReaderHistory) addChange(change *common.CacheChangeT) bool {
	if history.reader == nil {
		log.Fatalln("You need to create a Reader with this History before adding any changes")
	}

	history.mutex.Lock()
	defer history.mutex.Unlock()

	if history.Att.MemoryPolicy == resources.KPreallocatedMemoryMode &&
		change.SerializedPayload.Length > history.Att.PayloadMaxSize {
		log.Fatalf("Change payload size of '%v' bytes is larger than the history payload size of '%v' bytes and cannot be resized",
			change.SerializedPayload.Length, history.Att.PayloadMaxSize)
		return false
	}

	if change.WriterGUID == common.KGuidUnknown {
		log.Fatalln("The Writer GUID_t must be defined")
	}

	if len(history.changes) > 0 {
		if change.SourceTimestamp.Less(history.changes[len(history.changes)-1].SourceTimestamp) {
			i := 0
			for ; history.changes[i].SourceTimestamp.Less(change.SourceTimestamp); i++ {
			}
			newChanges := history.changes[:i]
			newChanges = append(newChanges, change)
			newChanges = append(newChanges)
			history.changes = newChanges
		}
	} else {
		history.changes = append(history.changes, change)
	}

	log.Printf("Change (%v) added with (%v) bytes", change.SequenceNumber, change.SerializedPayload)

	return true
}

func (history *ReaderHistory) MatchesChange(innerChange, outerChange *common.CacheChangeT) bool {
	if innerChange == nil || outerChange == nil {
		log.Fatalln("Pointer is not valid")
		return false
	}

	return innerChange.SequenceNumber == outerChange.SequenceNumber &&
		innerChange.WriterGUID == outerChange.WriterGUID
}

// Remove a specific change from the history.
// No Thread Safe
// @param release specifies if the change must be returned to the pool
func (history *ReaderHistory) RemoveChangeNts(removal int, release bool) {
	if history.reader == nil {
		log.Fatalln("You need to create a Reader with this History before adding any changes")
	}

	if removal >= len(history.changes) {
		log.Println("Trying to remove without a proper CacheChange_t referenced")
		return
	}

	change := history.changes[removal]
	newChanges := history.changes[:removal-1]
	newChanges = append(newChanges, history.changes[removal+1])
	history.changes = newChanges
	history.isHistoryFull = false

	history.reader.ChangeRemovedByHistory(change, nil)
	if release {
		history.reader.ReleaseCache(change)
	}
}

func (history *ReaderHistory) doReleaseCache(ch *common.CacheChangeT) {
	history.reader.ReleaseCache(ch)
}

func (history *ReaderHistory) doReserveCache(changes []*common.CacheChangeT) bool {
	return history.reader.ReserveCache(changes)
}
