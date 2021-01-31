package history

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"sync"
)

type IHistory interface {
	doReserveCache(changes []*common.CacheChangeT) bool
	doReleaseCache(ch *common.CacheChangeT)
}

type historyBase struct {
	Att     attributes.HistoryAttributes
	changes []*common.CacheChangeT
	mutex   sync.Mutex

	// Variable to know if the history is full without needing to block the History mutex.
	isHistoryFull bool
}

func NewhistoryBase(att *attributes.HistoryAttributes) *historyBase {
	var his historyBase
	his.Att = *att
	initialSize := att.InitialReservedCaches
	if initialSize < 0 {
		initialSize = 0
	}
	his.changes = make([]*common.CacheChangeT, initialSize)
	return &his
}
