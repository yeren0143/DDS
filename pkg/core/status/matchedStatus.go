package status

// MatchedStatus storing the subscription status
type MatchedStatus struct {
	// total cumulative count the concerned reader discovered a match with a writer
	// it found a writer for the same topic with a requested Qos that is compatible
	// with that offered by the reader
	TotalCount int32

	// The change in total_count since the last time the listener was called or
	// the status was read
	TotalCountChange int32

	// The number of writers currently matched to the concerned reader
	CurrentCount int32

	// The change in current_count since the last time the listener was called or the status was read
	CurrentCountChange int32
}
