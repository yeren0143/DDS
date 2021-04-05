package history

import (
	"log"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
)

type IWriterProxyWithHistory interface {
}

type IReaderWithHistory interface {
	ChangeRemovedByHistory(change *common.CacheChangeT, proxy IWriterProxyWithHistory) bool
	ReleaseCache(change *common.CacheChangeT)
	ReserveCache(size uint32) (*common.CacheChangeT, bool)
}

var _ IHistory = (*ReaderHistory)(nil)

type ReaderHistory struct {
	historyBase
	Reader IReaderWithHistory
}

func NewReaderHistory(att *attributes.HistoryAttributes) *ReaderHistory {
	var rhist ReaderHistory
	rhist.historyBase = *NewhistoryBase(att)
	rhist.impl = &rhist
	return &rhist
}

func (history *ReaderHistory) ReceivedChange(change *common.CacheChangeT) bool {
	return history.addChange(change)
}

func (history *ReaderHistory) addChange(change *common.CacheChangeT) bool {
	if history.Reader == nil {
		log.Fatalln("You need to create a Reader with this History before adding any changes")
	}

	history.Mutex.Lock()
	defer history.Mutex.Unlock()

	if history.Att.MemoryPolicy == resources.KPreallocatedMemoryMode &&
		change.SerializedPayload.Length > history.Att.PayloadMaxSize {
		log.Fatalf("Change payload size of '%v' bytes is larger than the history payload size of '%v' bytes and cannot be resized",
			change.SerializedPayload.Length, history.Att.PayloadMaxSize)
		return false
	}

	if change.WriterGUID == common.KGuidUnknown {
		log.Fatalln("The Writer GUID_t must be defined")
	}

	if len(history.Changes) > 0 {
		if change.SourceTimestamp.Less(history.Changes[len(history.Changes)-1].SourceTimestamp) {
			i := 0
			for ; history.Changes[i].SourceTimestamp.Less(change.SourceTimestamp); i++ {
			}
			newChanges := history.Changes[:i]
			newChanges = append(newChanges, change)
			newChanges = append(newChanges)
			history.Changes = newChanges
		}
	} else {
		history.Changes = append(history.Changes, change)
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
func (history *ReaderHistory) RemoveChangeNts(removal uint32, release bool) {
	if history.Reader == nil {
		log.Fatalln("You need to create a Reader with this History before adding any changes")
	}

	if int(removal) >= len(history.Changes) {
		log.Println("Trying to remove without a proper CacheChange_t referenced")
		return
	}

	change := history.Changes[removal]
	newChanges := history.Changes[:removal-1]
	newChanges = append(newChanges, history.Changes[removal+1])
	history.Changes = newChanges
	history.isHistoryFull = false

	history.Reader.ChangeRemovedByHistory(change, nil)
	if release {
		history.Reader.ReleaseCache(change)
	}
}

func (history *ReaderHistory) doReleaseCache(ch *common.CacheChangeT) {
	history.Reader.ReleaseCache(ch)
}

func (history *ReaderHistory) doReserveCache(size uint32) (*common.CacheChangeT, bool) {
	return history.Reader.ReserveCache(size)
}
