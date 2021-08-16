package writer

import (
	"log"

	"dds/common"
	"dds/core/policy"
	"dds/core/status"
)

// Class WriterListener with virtual method so the user can implement callbacks to certain events.
type IWriterListener interface {
	OnWriterMatched(awriter IRTPSWriter, info *common.MatchingInfo)

	/**
	* This method is called when a new Reader is discovered, with a Topic that
	* matches that of a local writer, but with a requested QoS that is incompatible
	* with the one offered by the local writer
	 */
	OnOfferedIncompatibleQos(awriter IRTPSWriter, qos policy.QosMask)

	/**
	    * This method is called when all the readers matched with this Writer acknowledge that a cache
		* change has been received.
	*/
	OnWriterChangeReceivedByAll(awriter IRTPSWriter, aChange *common.CacheChangeT)

	// Method called when the livelivess of a writer is lost
	OnLivelinessLost(awriter IRTPSWriter, status *status.LivelinessLostStatus)
}

type WriterListenerBase struct {
}

func (listener *WriterListenerBase) OnWriterMatched(awriter IRTPSWriter, info *common.MatchingInfo) {
	log.Println("nothing doing in base class")
}

func (listener *WriterListenerBase) OnOfferedIncompatibleQos(awriter IRTPSWriter, qos policy.QosMask) {
	log.Println("nothing doing in base class")
}

func (listener *WriterListenerBase) OnWriterChangeReceivedByAll(awriter IRTPSWriter, aChange *common.CacheChangeT) {
	log.Println("nothing doing in base class")
}

func (listener *WriterListenerBase) OnLivelinessLost(awriter IRTPSWriter, status *status.LivelinessLostStatus) {
	log.Println("nothing doing in base class")
}
