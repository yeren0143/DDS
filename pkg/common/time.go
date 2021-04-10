package common

import "time"

// Time used to describe times
type Time struct {
	Seconds int32
	// TODO:
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

func (t *Time) MilliSeconds() int64 {
	return int64(t.Seconds*1000) + int64(t.Nanosec)/1000
}

func CurrentTime() Time {
	timeStamp := time.Now()
	return Time{
		Seconds: int32(timeStamp.Unix()),
		Nanosec: uint32(timeStamp.UnixNano() - (timeStamp.Unix())*1e9),
	}
}

func NewTime(duration *time.Duration) *Time {
	return &Time{
		Seconds: int32(duration.Milliseconds()),
		Nanosec: uint32(duration.Nanoseconds() - duration.Milliseconds()),
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
