package reader

import (
	"log"

	"dds/common"
	"dds/core/policy"
	"dds/core/status"
)

// ReaderListener to be used by the user to override some of is virtual method to program actions to
// certain events.
type IReaderListener interface {
	// This method is invoked when a new reader matches
	OnReaderMatched(reader IRTPSReader)

	// This method is invoked when a new reader matches
	OnReaderMatchedWithStatus(reader IRTPSReader, info *status.SubscriptionMatchedStatus)

	/**
		 * This method is called when a new CacheChange_t is added to the ReaderHistory.
	     * CacheChange_t. This is a const pointer to const data
	     * to indicate that the user should not dispose of this data himself.
	     * To remove the data call the remove_change method of the ReaderHistory.
	     * reader->getHistory()->remove_change((CacheChange_t*)change).
	*/
	OnNewCacheChangeAdded(reader IRTPSReader, change *common.CacheChangeT)

	// Method called when the livelivess of a reader changes
	OnLivelinessChanged(reader IRTPSReader, liveness *status.LivelinessChangedStatus)

	// This method is called when a new Writer is discovered, with a Topic that
	OnRequestIncompatibleQos(reader IRTPSReader, qos policy.QosMask)
}

type ReaderListenerBase struct {
}

func (listener *ReaderListenerBase) OnReaderMatched(reader IRTPSReader) {
	log.Println("not impl OnReaderMatched")
}

func (listener *ReaderListenerBase) OnReaderMatchedWithStatus(reader IRTPSReader, info *status.SubscriptionMatchedStatus) {
	log.Println("not impl OnReaderMatchedWithStatus")
}

func (listener *ReaderListenerBase) OnNewCacheChangeAdded(reader IRTPSReader, change *common.CacheChangeT) {
	log.Println("not impl OnNewCacheChangeAdded")
	log.Fatalln("notImpl")
}

func (listener *ReaderListenerBase) OnLivelinessChanged(reader IRTPSReader, liveness *status.LivelinessChangedStatus) {
	log.Println("not impl OnLivelinessChanged")
}

func (listener *ReaderListenerBase) OnRequestIncompatibleQos(reader IRTPSReader, qos policy.QosMask) {
	log.Println("not impl OnRequestIncompatibleQos")
}
