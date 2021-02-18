package utils

import (
	"sync"
	"time"
)

//TimedConditionVariable implementation
type TimedConditionVariable struct {
	locker sync.Locker
	ch     chan bool
}

//NewTimedCond ...
func NewTimedCond(l sync.Locker) *TimedConditionVariable {
	return &TimedConditionVariable{
		ch:     make(chan bool),
		locker: l,
	}
}

//Wait ...
func (cond *TimedConditionVariable) Wait() {
	cond.locker.Unlock()
	<-cond.ch
	cond.locker.Lock()
}

// test 
func (cond *TimedConditionVariable) Lock() {
	cond.locker.Lock()
}

func (cond *TimedConditionVariable) Unlock() {
	cond.locker.Unlock()
}


//WaitOrTimeout ...
func (cond *TimedConditionVariable) WaitOrTimeout(d time.Duration) bool {
	tmo := time.NewTimer(d)
	cond.locker.Unlock()
	var ret bool
	select {
	case <-tmo.C:
		ret = false
	case <-cond.ch:
		ret = true
	}

	if !tmo.Stop() {
		select {
		case <-tmo.C:
		default:
		}
	}
	cond.locker.Lock()

	return ret
}

//Signal ...
func (cond *TimedConditionVariable) Signal() {
	cond.signal()
}

func (cond *TimedConditionVariable) signal() bool {
	select {
	case cond.ch <- true:
		return true
	default:
		return false
	}
}

//Broadcast ...
func (cond *TimedConditionVariable) Broadcast() {
	for {
		// Stop when we run out of waiters
		if !cond.signal() {
			return
		}
	}
}
