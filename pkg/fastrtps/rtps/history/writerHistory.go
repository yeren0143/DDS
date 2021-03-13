package history

import (
	"log"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
)

type writerCallback = func() uint32
type IWriterWithHistory interface {
	GetGUID() *common.GUIDT
	UnsentChangeAddedToHistory(aChange *common.CacheChangeT, maxBlockingTime common.Time)
	ReleaseChange(aChange *common.CacheChangeT) bool
	NewChange(dataCdrSerializedSize writerCallback, changeKind common.ChangeKindT, handle common.InstanceHandleT) *common.CacheChangeT
}

var _ IHistory = (*WriterHistory)(nil)
var _ historyImpl = (*WriterHistory)(nil)

// Class WriterHistory, container of the different CacheChanges of a writer
type WriterHistory struct {
	historyBase
	lastCacheChangeSeqNum common.SequenceNumberT
	Writer                IWriterWithHistory
}

func (wHistory *WriterHistory) AddChange(aChange *common.CacheChangeT, wparams *common.WriteParamsT) bool {
	return false
}

// Remove the CacheChange_t with the minimum sequenceNumber.
func (wHistory *WriterHistory) RemoveMinChange() bool {
	if wHistory.Writer == nil || wHistory.Mutex == nil {
		log.Fatalln("You need to create a Writer with this History before removing any changes")
		return false
	}

	wHistory.Mutex.Lock()
	defer wHistory.Mutex.Unlock()
	if len(wHistory.changes) > 0 && wHistory.RemoveChange(wHistory.changes[0]) {
		return true
	}
	return false
}

func (wHistory *WriterHistory) addChange(aChange *common.CacheChangeT, wparams *common.WriteParamsT,
	maxBlockingTime common.Time) bool {
	if wHistory.Writer == nil || wHistory.Mutex == nil {
		log.Fatalln("You need to create a Writer with this History before adding any changes")
		return false
	}

	wHistory.Mutex.Lock()
	defer wHistory.Mutex.Unlock()
	if aChange.WriterGUID != *wHistory.Writer.GetGUID() {
		log.Fatalln("Change writerGUID ", aChange.WriterGUID,
			" different than Writer GUID ", wHistory.Writer.GetGUID())
		return false
	}
	if wHistory.Att.MemoryPolicy == resources.KPreallocatedMemoryMode &&
		aChange.SerializedPayload.Length > wHistory.Att.PayloadMaxSize {
		log.Fatalln("Change payload size of '", aChange.SerializedPayload.Length,
			"' bytes is larger than the history payload size of '", wHistory.Att.PayloadMaxSize,
			"' bytes and cannot be resized.")
		return false
	}
	if wHistory.isHistoryFull {
		log.Fatalln("History full for writer ", aChange.WriterGUID)
		return false
	}

	wHistory.lastCacheChangeSeqNum.Value++
	aChange.SequenceNumber = wHistory.lastCacheChangeSeqNum
	aChange.SourceTimestamp = common.CurrentTime()

	aChange.WriteParams = *wparams
	// Updated sample and related sample identities on the user's write params
	wparams.SampleIdentity.WriterGUID = aChange.WriterGUID
	wparams.SampleIdentity.SequenceNumber = aChange.SequenceNumber
	wparams.ReleatedSampleIdentity = wparams.SampleIdentity

	wHistory.changes = append(wHistory.changes, aChange)
	if len(wHistory.changes) == int(wHistory.Att.MaximumReservedCaches) {
		wHistory.isHistoryFull = true
	}

	log.Println("Change ", aChange.SequenceNumber, " added with ",
		aChange.SerializedPayload.Length, " bytes")
	wHistory.Writer.UnsentChangeAddedToHistory(aChange, maxBlockingTime)
	return true
}

func (wHistory *WriterHistory) doReleaseCache(ch *common.CacheChangeT) {
	wHistory.Writer.ReleaseChange(ch)
}

func (wHistory *WriterHistory) doReserveCache(size uint32) (*common.CacheChangeT, bool) {
	callback := func() uint32 {
		return size
	}
	aChange := wHistory.Writer.NewChange(callback, common.KAlive, common.KInstanceHandleUnknown)
	return aChange, aChange != nil
}

func (wHistory *WriterHistory) RemoveChangeNts(removal uint32, release bool) {
	wHistory.isHistoryFull = false
	if release {
		wHistory.doReleaseCache(wHistory.changes[removal])
	}
	wHistory.changes = append(wHistory.changes[:removal], wHistory.changes[removal+1:]...)
}

func NewWriterHistory(att *attributes.HistoryAttributes) *WriterHistory {
	var hist WriterHistory
	hist.historyBase = *NewhistoryBase(att)
	hist.impl = &hist
	return &hist
}
