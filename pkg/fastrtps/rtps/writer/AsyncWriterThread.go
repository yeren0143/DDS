package writer

import (
	"log"
	"sync"

	"dds/fastrtps/utils"
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
		if !thread.running {
			thread.running = true
			go func() {
				thread.run()
			}()
		} else {
			thread.cv.Broadcast()
		}
	}

	log.Println("Wakeup finished")
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
