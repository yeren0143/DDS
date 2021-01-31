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
	KInfiniteSeconds     = 0x7fffffff
	KInfiniteNanoSeconds = 0xffffffff
)

var KTimeInfinite = Time{Seconds: KInfiniteSeconds, Nanosec: KInfiniteNanoSeconds}
var KTimeZero = Time{Seconds: 0, Nanosec: 0}
var KTimeInvalid = Time{Seconds: -1, Nanosec: KInfiniteNanoSeconds}
