package status

// BaseStatus storing the base status
type BaseStatus struct {
	TotalCount       int32
	TotalCountChange int32 // Increment since the last time the status was read
}

// SampleLostStatus is a alias of BaseStatus
type SampleLostStatus = BaseStatus

// LivelinessLostStatus is a alias of BaseStatus
type LivelinessLostStatus = BaseStatus

// InconsistentTopicStatus is a alias of BaseStatus
type InconsistentTopicStatus = BaseStatus
