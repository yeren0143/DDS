package common

type Time struct {
	Seconds int32
	Nanosec uint32
}

type Duration = Time

func CreateDuration(seconds int32, nanosec uint32) Time {
	return Time{
		Seconds: seconds,
		Nanosec: nanosec,
	}
}
