package writer

import (
	"dds/common"
)

type ChangeForReaderStatus uint8

const (
	KUnset ChangeForReaderStatus = iota
	KRequested
	KUnacknowledged
	KAcknowledged
	KUnderway
)

// Struct ChangeForReader used to represent the state of a specific change with respect to a specific reader,
// as well as its relevance.
type ChangeForReader struct {
	status          ChangeForReaderStatus
	isRelevant      bool
	seqNum          common.SequenceNumberT
	change          *common.CacheChangeT
	unsentFragments common.FragmentNumberSet
}

func (change *ChangeForReader) Less(otherChange *ChangeForReader) bool {
	return change.seqNum.Less(&otherChange.seqNum)
}
