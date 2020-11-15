package common

// Time used to describe times
type Time struct {
	Seconds int32
	Nanosec uint32
}

// DurationT used to describe Duration
type DurationT = Time

// time const
const (
	InfiniteSeconds     = 0x7fffffff
	InfiniteNanoSeconds = 0xffffffff
)

var (
	// TimeInfinite representing an infinite time
	TimeInfinite Time
	// TimeZero representing a zero time
	TimeZero Time
	// TimeInvalid representing an invalid time
	TimeInvalid Time
)

func init() {
	TimeInfinite = Time{Seconds: 0x7fffffff, Nanosec: 0xffffffff}
	TimeZero = Time{Seconds: 0, Nanosec: 0}
	TimeInvalid = Time{Seconds: -1, Nanosec: InfiniteNanoSeconds}
}
