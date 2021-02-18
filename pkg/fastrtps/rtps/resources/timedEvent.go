package resources

import (
	"time"

	"github.com/yeren0143/DDS/common"
)

//TimedEventCallback called when the expiration time expires.
type TimedEventCallback = func() bool

//TimedEvent can be used to launch an event through ResourceEvent's internal thread.
type TimedEvent struct {
	service *ResourceEvent
	impl    *TimedEventImpl
}

// Update event interval.
// When updating the interval, the timer is not restarted and the new interval will only
// be used the next time you call restart_timer().
func (timeEvent *TimedEvent) UpdateInterval(inter common.DurationT) bool {
	duration := time.Duration(inter.Seconds)*time.Second + time.Duration(inter.Nanosec)*time.Nanosecond
	return timeEvent.impl.UpdateInterval(duration)
}

func NewTimedEvent(service *ResourceEvent, callback *TimedEventCallback, milliseconds int64) *TimedEvent {
	var event TimedEvent
	event.impl = NewTimedEventImpl(callback, milliseconds)
	service.RegisterTimer(event.impl)
	return &event
}

//TimeEventVector implement sort
type TimeEventVector struct {
	Events []*TimedEventImpl
}

//Len get vector length
func (timeEvents *TimeEventVector) Len() int {
	return len(timeEvents.Events)
}

//Less ...
func (timeEvents *TimeEventVector) Less(i, j int) bool {
	lhs := timeEvents.Events[i].NextTriggerTime()
	rhs := timeEvents.Events[j].NextTriggerTime()
	return lhs.Before(rhs)
}

func (timeEvents *TimeEventVector) Swap(i, j int) {
	tmp := timeEvents.Events[i]
	timeEvents.Events[i] = timeEvents.Events[j]
	timeEvents.Events[j] = tmp
}
