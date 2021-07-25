package status

import (
	"github.com/yeren0143/DDS/common"
)

// SampleRejectedStatusKind is kind of possible values for the sample rejected reason
type SampleRejectedStatusKind = uint8

// enum of SampleRejectedStatusKind
const (
	KNotRejected SampleRejectedStatusKind = iota
	KRejectedByInstancesLimit
	KRejectedBySamplesLimit
	KRejectedBySamplesPerInstanceLimit
)

// SampleRejectedStatus storing the sample lost status
type SampleRejectedStatus struct {
	// Total cumulative count of samples rejected by the DataReader.
	TotalCount uint32

	// The incremental number of samples rejected since the last time the listener
	// was called or the status was read.
	TotalCountChange uint32

	// Reason for rejecting the last sample rejected.
	// If no samples have been rejected, the reason is the special value NOT_REJECTED.
	LastReason SampleRejectedStatusKind

	// Handle to the instance being updated by the last sample that was rejected.
	LastSubscriptionHandle common.InstanceHandleT
}

// CreateSampleRejectedStatus return SampleRejectedStatus with default value, KNotRejected
func CreateSampleRejectedStatus() SampleRejectedStatus {
	return SampleRejectedStatus{
		LastReason: KNotRejected,
	}
}
