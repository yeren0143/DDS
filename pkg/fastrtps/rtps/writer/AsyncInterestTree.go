package writer

import "sync"

/*
* Used by AsyncWriterThread to manage a double queue.
* One queue is being processed by AsyncWriterThread's internal thread while in the other one other threads can register
* RTPSWriter pointers that need to send samples asynchronously.
 */
type AsyncInterestTree struct {
	mutexActive   sync.Mutex
	mutexHidden   sync.Mutex
	activeWriters []IRTPSWriter
	hiddenWriters []IRTPSWriter
	activePos     uint32
	hiddenPos     uint32
}

func (tree *AsyncInterestTree) registerInterest(awriter IRTPSWriter) bool {
	tree.mutexHidden.Lock()
	defer tree.mutexHidden.Unlock()
	return tree.registerInterestNts(awriter)
}

func (tree *AsyncInterestTree) registerInterestNts(awriter IRTPSWriter) bool {
	for _, w := range tree.activeWriters {
		if w == awriter {
			return false
		}
	}
	tree.activeWriters = append(tree.activeWriters, awriter)
	return true
}

func (tree *AsyncInterestTree) swap() {

}

func (tree *AsyncInterestTree) nextActiveNts() IRTPSWriter {
	return nil
}
