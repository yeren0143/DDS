package resources

import (
	"sync"
	"sync/atomic"
	"time"
)

//TimedEventStateCode state code
type TimedEventStateCode = int32

//TimedEventCallback called when the expiration time expires.
type TimedEventCallback func() bool

//
const (
	CInactiveTimedEventStateCode TimedEventStateCode = iota // The event is inactive. The event service is not waiting for it.
	CReadyTimedEventStateCode                               // The event is ready for being processed by ResourceEvent and added to the event service.
	CWaitingTimedEventStateCode                             // The event is waiting for the event service to be triggered.
)

//TimedEvent can be used to launch an event through ResourceEvent's internal thread.
type TimedEvent struct {
	Interval        time.Duration // Expiration time in microseconds of the event
	nextTriggerTime time.Time
	state           TimedEventStateCode
	Callback        *TimedEventCallback
	mutex           sync.Mutex
}

//NextTriggerTime Returns next trigger time as a time point
func (timeEvent *TimedEvent) NextTriggerTime() time.Time {
	timeEvent.mutex.Lock()
	defer timeEvent.mutex.Unlock()
	return timeEvent.nextTriggerTime
}

// Update It updates the timer depending on the state of the TimedEventImpl object.
// This method has to be called from ResourceEvent's internal thread.
// false if the event was canceled, true otherwise.
func (timeEvent *TimedEvent) Update(currentTime, cancelTime time.Time) bool {
	expected := CReadyTimedEventStateCode
	setTime := atomic.CompareAndSwapInt32(&timeEvent.state, expected, CWaitingTimedEventStateCode)

	if setTime {
		timeEvent.mutex.Lock()
		timeEvent.nextTriggerTime = currentTime.Add(timeEvent.Interval)
		timeEvent.mutex.Unlock()
	} else if expected == CInactiveTimedEventStateCode {
		timeEvent.mutex.Lock()
		timeEvent.nextTriggerTime = cancelTime
		timeEvent.mutex.Unlock()
	}

	return expected != CInactiveTimedEventStateCode
}

//Trigger triggers the callback action
func (timeEvent *TimedEvent) Trigger(currentTime, cancelTime time.Time) {
	if timeEvent.Callback != nil {
		expected := CWaitingTimedEventStateCode
		atomic.CompareAndSwapInt32(&timeEvent.state, expected, CInactiveTimedEventStateCode)

		restart := (*timeEvent.Callback)()
		if restart {
			expected = CInactiveTimedEventStateCode
			if atomic.CompareAndSwapInt32(&timeEvent.state, expected, CWaitingTimedEventStateCode) {
				timeEvent.mutex.Lock()
				defer timeEvent.mutex.Unlock()

				timeEvent.nextTriggerTime = currentTime.Add(timeEvent.Interval)
				return
			}
		}

		timeEvent.mutex.Lock()
		defer timeEvent.mutex.Unlock()
		timeEvent.nextTriggerTime = cancelTime
	}
}

//TimeEventVector implement sort
type TimeEventVector struct {
	Events []*TimedEvent
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
