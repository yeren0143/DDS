package writer

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
)

// RTPSWriter manages the sending of data to the readers. Is always associated with a HistoryCache.
type IRTPSWriter interface {
	// Process an incoming ACKNACK submessage.
	// result true if the writer could process the submessage. Only valid when returned value is true.
	// valid true when the submessage was destinated to this writer, false otherwise.
	ProcessAcknack(writerGUID, readerGUID *common.GUIDT, ackCount uint32,
		snSet *common.SequenceNumberSet, finalFlag bool) (result, valid bool)

	// result true if the writer could process the submessage. Only valid when returned value is true
	// true when the submessage was destinated to this writer, false otherwise.
	ProcessNackFrag(writerGUID, readerGUID *common.GUIDT, ackCount uint32,
		seqNum *common.SequenceNumberT, fragmentsState *common.FragmentNumberSet) (result, valid bool)

	GetAttributes() *attributes.EndpointAttributes
}
