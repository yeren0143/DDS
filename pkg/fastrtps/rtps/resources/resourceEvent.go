package resources

import (
	"log"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"dds/fastrtps/utils"
)

// ResourceEvent centralizes all operations over timed events in the same thread.
type ResourceEvent struct {
	stop                    int32 //Warns the internal thread can stop.
	allowVectorManipulation bool  //Flag used to allow a thread to manipulate the timer collections when the execution thread is not using them.
	mutex                   sync.Mutex
	cvManipulation          utils.TimedConditionVariable
	cv                      utils.TimedConditionVariable
	timersCount             int //The total number of created timers.
	PendingTimers           TimeEventVector
	ActiveTimers            TimeEventVector
	currentTime             time.Time
}

//NewResourceEvent create resource event with default value
func NewResourceEvent() *ResourceEvent {
	var event ResourceEvent
	// event.mutex = new(sync.Mutex)
	event.stop = 0
	event.allowVectorManipulation = true
	event.timersCount = 0
	event.cvManipulation = *utils.NewTimedCond(&event.mutex)
	event.cv = *utils.NewTimedCond(&event.mutex)
	return &event
}

// This method informs that a TimedEventImpl has been created.
// This method has to be called when creating a TimedEventImpl object.
func (resource *ResourceEvent) RegisterTimer(event *TimedEventImpl) {
	resource.mutex.Lock()
	defer resource.mutex.Unlock()
	resource.timersCount++
	resource.cv.Signal()
}

// Registers a new TimedEventImpl object in the internal queue to be processed.
// Non thread safe.
// return True value if the insertion was successful. In other case, it return False.
func (resource *ResourceEvent) registerTimerNts(event *TimedEventImpl) bool {
	if !resource.PendingTimers.Has(event) {
		resource.PendingTimers.Push(event)
		return true
	}
	return false
}

// This method notifies to ResourceEvent that the TimedEventImpl object has operations to be
// scheduled. These operations can be the cancellation of the timer or starting another async_wait.
//
func (resource *ResourceEvent) Notify(event *TimedEventImpl) {
	resource.mutex.Lock()
	defer resource.mutex.Unlock()
	if resource.registerTimerNts(event) {
		resource.cv.Signal()
	}
}

//InitThread to initialize the internal thread.
func (resource *ResourceEvent) InitThread() {
	resource.mutex.Lock()
	defer resource.mutex.Unlock()

	resource.allowVectorManipulation = false
	resource.ResizeCollections()

	go func(resource *ResourceEvent) {
		resource.eventService()
	}(resource)

	log.Println("InitThread finished")
}

//ResizeCollections Ensures internal collections can accommodate current total number of timers.
func (resource *ResourceEvent) ResizeCollections() {
}

func (resource *ResourceEvent) eventService() {
	for atomic.LoadInt32(&resource.stop) <= 0 {
		// Perform update and execution of timers
		resource.updateCurrentTime()
		resource.doTimerActions()

		resource.mutex.Lock()

		// If the thread has already been instructed to stop, do it.
		if atomic.LoadInt32(&resource.stop) > 0 {
			break
		}

		// If pending timers exist, there is some work to be done, so no need to wait.
		if len(resource.PendingTimers.Events) > 0 {
			continue
		}

		resource.allowVectorManipulation = true
		resource.cvManipulation.Broadcast()

		// Wait for the first timer to be triggered
		nextTrigger := time.Time{}
		if len(resource.ActiveTimers.Events) == 0 {
			nextTrigger = resource.currentTime.Add(time.Second)
		} else {
			nextTrigger = resource.ActiveTimers.Events[0].nextTriggerTime
		}

		if nextTrigger.Before(resource.currentTime) {
			log.Fatal("next trigger time can not brefore current time")
		}

		// test test
		//resource.cv.WaitOrTimeout(nextTrigger.Sub(resource.currentTime))

		// Don't allow other threads to manipulate the timer collections
		resource.allowVectorManipulation = false
		resource.ResizeCollections()

		resource.mutex.Unlock()
	}
}

func (resource *ResourceEvent) updateCurrentTime() {
	resource.currentTime = time.Now()
}

func (resource *ResourceEvent) sortTimers() {
	sort.Sort(&resource.ActiveTimers)
}

func (resource *ResourceEvent) doTimerActions() {
	d, _ := time.ParseDuration("24h")
	cancelTime := resource.currentTime
	cancelTime.Add(d)

	didSomeThing := false

	//Process pending orders
	// TODO: 可重入锁
	{
		resource.mutex.Lock()

		sort.Sort(&resource.ActiveTimers)
		for _, tp := range resource.PendingTimers.Events {
			// Remove item from active timers
			for i, activeTp := range resource.ActiveTimers.Events {
				log.Println("looping resource.ActiveTimers.Events")
				if activeTp == tp {
					events := resource.ActiveTimers.Events
					if i == len(events)-1 {
						resource.ActiveTimers.Events = events[:len(events)-1]
					} else {
						resource.ActiveTimers.Events = append(events[:i], events[i+1:]...)
					}
					break
				}
			}

			// Update timer info
			if tp.Update(resource.currentTime, cancelTime) {
				index := len(resource.ActiveTimers.Events)
				for i, activeTp := range resource.ActiveTimers.Events {
					if tp.nextTriggerTime.Before(activeTp.nextTriggerTime) {
						index = i
						break
					}
				}

				// Insert on correct position
				events := append([]*TimedEventImpl{}, resource.ActiveTimers.Events[index:]...)
				resource.ActiveTimers.Events = append(resource.ActiveTimers.Events[:index], tp)
				resource.ActiveTimers.Events = append(resource.ActiveTimers.Events, events...)
			}
		}
		resource.PendingTimers = TimeEventVector{}

		resource.mutex.Unlock()
	}

	// Trigger active timers
	for _, tp := range resource.ActiveTimers.Events {
		if tp.nextTriggerTime.Before(resource.currentTime) {
			didSomeThing = true
			tp.Trigger(resource.currentTime, cancelTime)
		} else {
			break
		}
	}

	// If an action was made, keep active_timers_ sorted
	if didSomeThing {
		resource.sortTimers()

		index := len(resource.ActiveTimers.Events)
		for i, activeTp := range resource.ActiveTimers.Events {
			if cancelTime.Before(activeTp.nextTriggerTime) {
				index = i
				break
			}
		}

		if index < len(resource.ActiveTimers.Events) {
			resource.ActiveTimers.Events = append([]*TimedEventImpl{}, resource.ActiveTimers.Events[index:]...)
		}
	}
}
