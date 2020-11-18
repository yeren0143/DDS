package common

// Time used to describe times
type Time struct {
	Seconds int32
	Nanosec uint32
}

// DurationT used to describe Duration
type DurationT = Time

func (t *Time) Less(value Time) bool {
	if t.Seconds < value.Seconds {
		return true
	} else if t.Seconds > value.Seconds {
		return false
	} else if t.Nanosec < value.Nanosec {
		return true
	} else {
		return false
	}
}

// time const
const (
	InfiniteSeconds     = 0x7fffffff
	InfiniteNanoSeconds = 0xffffffff
)

var (
	// TimeInfinite representing an infinite time
	CTimeInfinite Time
	// TimeZero representing a zero time
	CTimeZero Time
	// TimeInvalid representing an invalid time
	CTimeInvalid Time
)

func init() {
	CTimeInfinite = Time{Seconds: InfiniteSeconds, Nanosec: InfiniteNanoSeconds}
	CTimeZero = Time{Seconds: 0, Nanosec: 0}
	CTimeInvalid = Time{Seconds: -1, Nanosec: InfiniteNanoSeconds}
}
