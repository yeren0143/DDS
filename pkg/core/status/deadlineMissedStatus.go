package status

import (
	"dds/common"
)

// DeadlineMissedStatus storing the deadline status
type DeadlineMissedStatus struct {
	TotalCount         uint32
	TotalCountChange   uint32
	LastInstanceHandle common.InstanceHandleT
}

// OfferedDeadlineMissedStatus is alias of DeadlineMissedStatus
type OfferedDeadlineMissedStatus = DeadlineMissedStatus

// RequestedDeadlineMissedStatus is alias of DeadlineMissedStatus
type RequestedDeadlineMissedStatus = DeadlineMissedStatus
