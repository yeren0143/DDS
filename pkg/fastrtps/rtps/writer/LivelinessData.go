package writer

import (
	"time"

	"dds/common"
	"dds/core/policy"
)

type WriterStatus uint8

const (
	// Writer is matched but liveliness has not been asserted yet
	KNotAsserted WriterStatus = iota
	// Writer is alive
	KAlive
	// Writer is not alive
	KNotAlive
)

// A struct keeping relevant liveliness information of a writer
type LivelinessData struct {
	GUID          common.GUIDT
	Kind          policy.LivelinessQosPolicyKind
	LeaseDuration common.DurationT
	Count         uint
	Status        WriterStatus
	// The time when the writer will lose liveliness
	TimePoint time.Time
}
