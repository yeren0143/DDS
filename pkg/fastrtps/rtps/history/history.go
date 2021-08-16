package history

import (
	"log"
	"sync"

	"dds/common"
	"dds/fastrtps/rtps/attributes"
)

type IHistory interface {
	doReserveCache(size uint32) (*common.CacheChangeT, bool)
	doReleaseCache(ch *common.CacheChangeT)
	GetHistorySize() uint32
	RemoveChange(ch *common.CacheChangeT) bool
}

type historyImpl interface {
	RemoveChangeNts(removal uint32, release bool)
}

type historyBase struct {
	Att     attributes.HistoryAttributes
	Changes []*common.CacheChangeT
	Mutex   *sync.Mutex

	// Variable to know if the history is full without needing to block the History mutex.
	isHistoryFull bool

	impl historyImpl
}

func (hist *historyBase) GetHistorySize() uint32 {
	hist.Mutex.Lock()
	defer hist.Mutex.Unlock()
	return uint32(len(hist.Changes))
}

func (hist *historyBase) RemoveChange(ch *common.CacheChangeT) bool {
	// hist.mutex.Lock()
	// defer hist.mutex.Unlock()

	index := 0
	for ; index < len(hist.Changes); index++ {
		if hist.Changes[index] == ch {
			break
		}
	}
	if index == len(hist.Changes) {
		log.Fatalln("Trying to remove a change not in history")
		return false
	}

	// remove using the virtual method
	hist.impl.RemoveChangeNts(uint32(index), true)

	return true

}

func NewhistoryBase(att *attributes.HistoryAttributes) *historyBase {
	var his historyBase
	his.Att = *att
	// initialSize := att.InitialReservedCaches
	// if initialSize < 0 {
	// 	initialSize = 0
	// }
	// his.Changes = make([]*common.CacheChangeT, initialSize)
	return &his
}
