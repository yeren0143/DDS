package reader

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/endpoint"
)

// IRtpsReader manages the reception of data from its matched writers.
type IRTPSReader interface {
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
	ChangeRemovedByHistory() bool

	// Accept msg to unknwon readers (default=true)
	AcceptMessagesToUnknownReaders() bool

	GetAttributes() *attributes.EndpointAttributes
}

type rtpsReaderBase struct {
	endpoint.Endpoint
	acceptMessagesToUnknownReaders bool // Accept msg to unknwon readers (default=true)
}

func (reader *rtpsReaderBase) AcceptMessagesToUnknownReaders() bool {
	return reader.acceptMessagesToUnknownReaders
}

func (reader *rtpsReaderBase) ReleaseCache(change *common.CacheChangeT) {
	reader.Mutex.Lock()
	defer reader.Mutex.Unlock()

	if pool := change.PayloadOwner(); pool != nil {
		pool.ReleasePayload(change)
	}
	reader.ChangePool.ReleaseCache(change)
}
