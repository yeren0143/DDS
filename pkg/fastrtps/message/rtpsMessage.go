package message

import (
	"github.com/yeren0143/DDS/common"
)

// SubmessageID is Enumeration of the different Submessages types
type SubmessageID = uint8

const (
	KPad           SubmessageID = 0x01
	KAcknack       SubmessageID = 0x06
	KHeartbeat     SubmessageID = 0x07
	KGap           SubmessageID = 0x08
	KInfoTs        SubmessageID = 0x09
	KInfoSrc       SubmessageID = 0x0c
	KInfoReplyIP4  SubmessageID = 0x0d
	KInfoDst       SubmessageID = 0x0e
	KInfoReply     SubmessageID = 0x0f
	KNackFrag      SubmessageID = 0x12
	KHeartbeatFrag SubmessageID = 0x13
	KData          SubmessageID = 0x15
	KDataFrag      SubmessageID = 0x16
)

type HeaderT struct {
	Version    common.ProtocolVersionT
	VendorID   common.VendorIDT
	GUIDPrefix common.GUIDPrefixT
}

// SubmessageHeaderT used to contain the header information of a submessage.
type SubmessageHeaderT struct {
	SubMessageID     common.Octet
	SubmessageLength uint32
	Flags            common.SubmessageFlag
	IsLast           bool
}

type IRtpsMsgReader interface {
	AcceptMessagesToUnknownReaders() bool
	ProcessHeartbeatMsg(writerGUID *common.GUIDT,
		hbCount uint32,
		firstSN *common.SequenceNumberT,
		lastSN *common.SequenceNumberT,
		finalFlag bool,
		livelinessFlag bool) bool

	ProcessGapMsg(writerGUID *common.GUIDT,
		gapStart *common.SequenceNumberT,
		gapList *common.SequenceNumberSet) bool

	ProcessDataMsg(change *common.CacheChangeT) bool

	ProcessDataFragMsg(change *common.CacheChangeT,
		sampleSize uint32,
		fragmentStartingNum uint32,
		fragmentsInSubmessage uint16) bool
}

type IRtpsMsgWriter interface {
	// Process an incoming ACKNACK submessage.
	// result true if the writer could process the submessage. Only valid when returned value is true.
	// valid true when the submessage was destinated to this writer, false otherwise.
	ProcessAcknack(writerGUID, readerGUID *common.GUIDT, ackCount uint32,
		snSet *common.SequenceNumberSet, finalFlag bool) (result, valid bool)

	// result true if the writer could process the submessage. Only valid when returned value is true
	// true when the submessage was destinated to this writer, false otherwise.
	ProcessNackFrag(writerGUID, readerGUID *common.GUIDT, ackCount uint32,
		seqNum *common.SequenceNumberT, fragmentsState *common.FragmentNumberSet) (result, valid bool)
}
