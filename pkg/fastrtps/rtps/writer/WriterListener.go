package writer

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/core/policy"
	"github.com/yeren0143/DDS/core/status"
)

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
