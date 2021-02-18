package reader

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/history"
	"github.com/yeren0143/DDS/fastrtps/utils"
)

//var _ IRTPSReader = (*StatefulReader)(nil)

// Class StatefulReader, specialization of RTPSReader than stores the state of the matched writers.
type StatefulReader struct {
	rtpsReaderBase
	acknackCount  uint32
	nackfragCount uint32
	readerTimes   attributes.ReaderTimes
	// Vector containing pointers to all the active WriterProxies.
	matchedWriters []*WriterProxy
	// Vector containing pointers to all the inactive, ready for reuse, WriterProxies.
	matchedWritersPool  []*WriterProxy
	proxyChangesConfig  utils.ResourceLimitedContainerConfig
	disablePositiveAcks bool
	isAlive             bool
}

func (reader *StatefulReader) ChangeRemovedByHistory(change *common.CacheChangeT, prox history.IWriterProxyWithHistory) bool {
	if !change.IsRead {
		if reader.totalUnread > 0 {
			reader.totalUnread--
		}
	}
	return true
}
