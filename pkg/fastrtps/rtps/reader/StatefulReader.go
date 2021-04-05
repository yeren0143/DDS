package reader

import (
	"log"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/data"
	"github.com/yeren0143/DDS/fastrtps/rtps/endpoint"
	"github.com/yeren0143/DDS/fastrtps/rtps/history"
	"github.com/yeren0143/DDS/fastrtps/utils"
)

var _ IRTPSReader = (*StatefulReader)(nil)

// Class StatefulReader, specialization of RTPSReader than stores the state of the matched writers.
type StatefulReader struct {
	RTPSReader
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

func (reader *StatefulReader) MatchedWriterAdd(wdata *data.WriterProxyData) bool {
	log.Fatalln("notimpl")
	return false
}

func (reader *StatefulReader) MatchedWriterRemove(writerGuid *common.GUIDT, removedByLease bool) bool {
	log.Fatalln("notimpl")
	return false
}

func (reader *StatefulReader) MatchedWriterIsMatched(writerGuid *common.GUIDT) bool {
	log.Fatalln("notimpl")
	return false
}

func (reader *StatefulReader) MatchedWriterLookup(writerGuid *common.GUIDT) (bool, *WriterProxy) {
	log.Fatalln("notimpl")
	return false, nil
}

func (reader *StatefulReader) ProcessDataMsg(aChange *common.CacheChangeT) bool {
	log.Fatalln("notimpl")
	return false
}

func (reader *StatefulReader) ProcessDataFragMsg(change *common.CacheChangeT,
	sampleSize uint32, fragmentStartingNum uint32, fragmentsInSubmessage uint16) bool {
	log.Fatalln("notimpl")
	return false
}

func (reader *StatefulReader) ProcessHeartbeatMsg(writerGUID *common.GUIDT,
	hbCount uint32, firstSN *common.SequenceNumberT, lastSN *common.SequenceNumberT,
	finalFlag bool, livelinessFlag bool) bool {
	log.Fatalln("notimpl")
	return false
}

func (reader *StatefulReader) ProcessGapMsg(writerGUID *common.GUIDT,
	gapStart *common.SequenceNumberT,
	gapList *common.SequenceNumberSet) bool {
	log.Fatalln("notimpl")
	return false
}

func (reader *StatefulReader) ChangeReceived(aChange *common.CacheChangeT, prox *WriterProxy) bool {
	log.Fatalln("notimpl")
	return false
}

func (reader *StatefulReader) initWithParticipant(pimpl endpoint.IEndpointParent, att *attributes.ReaderAttributes) {
	partAtt := pimpl.GetAttributes()
	for n := uint32(0); n < att.MatchedWritersAllocation.Initial; n++ {
		proxy := NewWriterProxy(reader, partAtt.Allocation.Locators, &reader.proxyChangesConfig)
		reader.matchedWritersPool = append(reader.matchedWritersPool, proxy)
	}
}

func NewStatefulReader(parent endpoint.IEndpointParent, guid *common.GUIDT, att *attributes.ReaderAttributes,
	payloadPool history.IPayloadPool, hist *history.ReaderHistory, rlisten IReaderListener) *StatefulReader {
	var retReader StatefulReader
	retReader.readerTimes = att.Times
	retReader.proxyChangesConfig = history.ResourceLimitsFromHistory(&hist.Att, 0)
	retReader.disablePositiveAcks = att.DisablePositiveAcks
	retReader.isAlive = true
	var aChangePool history.IChangePool
	retReader.RTPSReader = *NewRtpsReader(parent, guid, att, payloadPool, aChangePool, hist, rlisten)
	retReader.RTPSReader.impl = &retReader

	retReader.initWithParticipant(parent, att)

	return &retReader
}
