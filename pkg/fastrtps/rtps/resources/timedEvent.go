package resources

import "time"

type TimedEvent_StateCode uint8

type TimedEvent_Callback func() bool

const (
	INACTIVE TimedEvent_StateCode = iota // The event is inactive. The event service is not waiting for it.
	READY                                // The event is ready for being processed by ResourceEvent and added to the event service.
	WAITING                              // The event is waiting for the event service to be triggered.
)

type TimedEvent struct {
	Interval_microsec time.Duration // Expiration time in microseconds of the event
	Next_trigger_time time.Time
	state             TimedEvent_StateCode
	Callback          TimedEvent_Callback
}
