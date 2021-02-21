package writer

import (
	"sync"

	"github.com/yeren0143/DDS/fastrtps/utils"
)

// This static class owns a thread that manages asynchronous writes.
// Asynchronous writes happen directly (when using an async writer) and
// indirectly (when responding to a NACK).
type AsyncWriterThread struct {
	conditionVariableMutex *sync.Mutex
	cv                     *utils.TimedConditionVariable
	interestTree           *AsyncInterestTree
	running                bool
	runScheduled           bool
	routeRunning           bool
}

func (thread *AsyncWriterThread) run() {
	thread.conditionVariableMutex.Lock()
	defer thread.conditionVariableMutex.Unlock()
	for thread.running {
		if thread.runScheduled {
			thread.runScheduled = false
			thread.conditionVariableMutex.Unlock()
			thread.interestTree.swap()

			thread.interestTree.mutexActive.Lock()
			curr := thread.interestTree.nextActiveNts()
			
			for ; curr != nil; curr = thread.interestTree.nextActiveNts() {
				curr.SendAnyUnsentChanges()
			}
			thread.interestTree.mutexActive.Unlock()

			thread.conditionVariableMutex.Lock()
		} else {
			thread.cv.Wait()
		}
	}
}

// Wakes the thread up and starts processing async writers.
func (thread *AsyncWriterThread) Wakeup(awriter IRTPSWriter) {
	if thread.interestTree.registerInterest(awriter) {
		thread.conditionVariableMutex.Lock()
		defer thread.conditionVariableMutex.Unlock()
		thread.runScheduled = true
		if !thread.routeRunning {
			thread.running = true
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				wg.Done()
				thread.run()
			}()
			wg.Wait()
		} else {
			thread.cv.Broadcast()
		}
	}
}

func NewAsyncWriterThread() *AsyncWriterThread {
	var asyThread AsyncWriterThread
	asyThread.conditionVariableMutex = new(sync.Mutex)
	asyThread.cv = utils.NewTimedCond(asyThread.conditionVariableMutex)
	asyThread.runScheduled = false
	asyThread.interestTree = &AsyncInterestTree{}
	asyThread.runScheduled = false
	return &asyThread
}
