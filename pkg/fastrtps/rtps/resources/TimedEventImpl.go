package resources

import (
	"sync"
	"sync/atomic"
	"time"
)

type TimedEventState = int32

const (
	KInactive TimedEventState = iota // The event is inactive. The event service is not waiting for it
	KReady                           // The event is ready for being processed by ResourceEvent and added to the event service.
	KWaiting                         // The event is waiting for the event service to be triggered.
)

type TimedEventImpl struct {
	intervalMicrosec time.Duration
	nextTriggerTime  time.Time
	callback         *TimedEventCallback
	stateCode        TimedEventState
	mutex            sync.Mutex //  Protects interval_microsec_ and next_trigger_time_
}

//NextTriggerTime Returns next trigger time as a time point
func (timeEvent *TimedEventImpl) NextTriggerTime() time.Time {
	timeEvent.mutex.Lock()
	defer timeEvent.mutex.Unlock()
	return timeEvent.nextTriggerTime
}

func (timeEvent *TimedEventImpl) UpdateInterval(interval time.Duration) bool {
	timeEvent.mutex.Lock()
	defer timeEvent.mutex.Unlock()
	timeEvent.intervalMicrosec = interval
	return true
}

// Update It updates the timer depending on the state of the TimedEventImpl object.
// This method has to be called from ResourceEvent's internal thread.
// false if the event was canceled, true otherwise.
func (timeEvent *TimedEventImpl) Update(currentTime, cancelTime time.Time) bool {
	expected := KReady
	setTime := atomic.CompareAndSwapInt32(&timeEvent.stateCode, expected, KWaiting)

	if setTime {
		timeEvent.mutex.Lock()
		timeEvent.nextTriggerTime = currentTime.Add(timeEvent.intervalMicrosec * time.Microsecond)
		timeEvent.mutex.Unlock()
	} else if expected == KInactive {
		timeEvent.mutex.Lock()
		timeEvent.nextTriggerTime = cancelTime
		timeEvent.mutex.Unlock()
	}

	return expected != KInactive
}

//Trigger triggers the callback action
func (timeEvent *TimedEventImpl) Trigger(currentTime, cancelTime time.Time) {
	if timeEvent.callback != nil {
		expected := KWaiting
		atomic.CompareAndSwapInt32(&timeEvent.stateCode, expected, KInactive)

		restart := (*timeEvent.callback)()
		if restart {
			expected = KInactive
			if atomic.CompareAndSwapInt32(&timeEvent.stateCode, expected, KWaiting) {
				timeEvent.mutex.Lock()
				defer timeEvent.mutex.Unlock()
				timeEvent.nextTriggerTime = currentTime.Add(timeEvent.intervalMicrosec * time.Microsecond)
				return
			}
		}

		timeEvent.mutex.Lock()
		defer timeEvent.mutex.Unlock()
		timeEvent.nextTriggerTime = cancelTime
	}
}

func NewTimedEventImpl(callback *TimedEventCallback, interval int64) *TimedEventImpl {
	var eventImpl TimedEventImpl
	eventImpl.stateCode = KInactive
	eventImpl.callback = callback
	eventImpl.intervalMicrosec = time.Duration(interval * int64(time.Microsecond))
	return &eventImpl
}
