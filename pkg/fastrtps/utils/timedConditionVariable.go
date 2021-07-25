package utils

import (
	"sync"
	"time"
)

//TimedConditionVariable implementation
type TimedConditionVariable struct {
	locker sync.Locker
	cond   *sync.Cond
	ch     chan bool
}

//NewTimedCond ...
func NewTimedCond(l sync.Locker) *TimedConditionVariable {
	return &TimedConditionVariable{
		locker: l,
		cond:   sync.NewCond(l),
		ch:     make(chan bool),
	}
}

//Wait ...
func (timedCond *TimedConditionVariable) Wait() {
	// cond.locker.Unlock()
	// <-cond.ch
	// cond.locker.Lock()
	timedCond.cond.Wait()
}

// // test
// func (cond *TimedConditionVariable) Lock() {
// 	cond.locker.Lock()
// }

// func (cond *TimedConditionVariable) Unlock() {
// 	cond.locker.Unlock()
// 	cond.ch <- true
// }

//WaitOrTimeout ...
func (timedCond *TimedConditionVariable) WaitOrTimeout(d time.Duration) bool {
	// tmo := time.NewTimer(d)
	// cond.locker.Unlock()
	// var ret bool
	// select {
	// case <-tmo.C:
	// 	ret = false
	// case <-cond.ch:
	// 	ret = true
	// }

	// if !tmo.Stop() {
	// 	select {
	// 	case <-tmo.C:
	// 	default:
	// 	}
	// }
	// cond.locker.Lock()

	tmo := time.NewTimer(d)
	var ret bool
	select {
	case <-tmo.C:
		ret = false
		timedCond.locker.Unlock()
		timedCond.locker.Lock()
	case <-timedCond.ch:
		ret = true
	}

	return ret
}

//Signal ...
func (timedCond *TimedConditionVariable) Signal() {
	timedCond.cond.Signal()
	timedCond.signal()
}

func (timedCond *TimedConditionVariable) signal() bool {
	select {
	case timedCond.ch <- true:
		return true
	default:
		return false
	}
}

//Broadcast ...
func (timedCond *TimedConditionVariable) Broadcast() {
	timedCond.cond.Broadcast()
	for {
		// Stop when we run out of waiters
		if !timedCond.signal() {
			return
		}
	}
}
