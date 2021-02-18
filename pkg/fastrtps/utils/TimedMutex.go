package utils

import (
	"sync"
	"time"
)

type TimeMutex struct {
	locker sync.Mutex
	ch     chan bool
}

func NewTimeMutex() *TimeMutex {
	return &TimeMutex{
		ch: make(chan bool),
	}
}

func (tMutex *TimeMutex) TryLockUntil(d time.Duration) {
	tmo := time.NewTimer(d)
	var wg sync.WaitGroup
	wg.Add(1)
	tMutex.locker.Lock()
	go func() {
		wg.Done()
		select {
		case <-tmo.C:
			tMutex.locker.Unlock()
		}
	}()
	wg.Wait()
}
