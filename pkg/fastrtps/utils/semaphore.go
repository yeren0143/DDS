package utils

import (
	"sync"
)

//Semaphore impl semaphore with go
type Semaphore struct {
	count   uint32
	mutex   sync.Mutex
	cv      *sync.Cond
	disable bool
}

// Post simulation semaphore post
func (sem *Semaphore) Post() {
	sem.mutex.Lock()
	defer sem.mutex.Unlock()
	if !sem.disable {
		sem.count++
		sem.cv.Signal()
	}
}

// PostN signal n times
func (sem *Semaphore) PostN(n uint32) {
	sem.mutex.Lock()
	defer sem.mutex.Unlock()

	if !sem.disable {
		sem.count += n
		sem.cv.Signal()
	}
}

// Disable unable semaphore
func (sem *Semaphore) Disable() {
	sem.mutex.Lock()
	defer sem.mutex.Unlock()
	if !sem.disable {
		sem.count--
		sem.cv.Broadcast()
		sem.disable = true
	}
}

// Enable semaphore
func (sem *Semaphore) Enable() {
	sem.mutex.Lock()
	defer sem.mutex.Unlock()
	if sem.disable {
		sem.count = 0
		sem.disable = false
	}
}

//Wait signal resource
func (sem *Semaphore) Wait() {
	sem.mutex.Lock()
	defer sem.mutex.Unlock()

	if !sem.disable {
		for !sem.disable && sem.count <= 0 {
			sem.cv.Wait()
		}
		sem.count--
	}
}

//NewSemaphore create semaphore with n resources
func NewSemaphore(count uint32) *Semaphore {
	var sem Semaphore
	sem.count = 0
	sem.disable = false
	sem.cv = sync.NewCond(&sem.mutex)
	return &sem
}
