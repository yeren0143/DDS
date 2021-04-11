package pdp

import (
	"log"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/core/policy"
	"github.com/yeren0143/DDS/fastrtps/message"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/protocol"
	"github.com/yeren0143/DDS/fastrtps/rtps/reader"
)

var _ reader.IReaderListener = (*PDPListener)(nil)

/**
 * Class PDPListener, specification used by the PDP to perform the History check when a new message is received.
 * This class is implemented in order to use the same structure than with any other RTPSReader.
 * @ingroup DISCOVERY_MODULE
 */
type PDPListener struct {
	reader.ReaderListenerBase

	parentPDP protocol.IPDP
}

func (listener *PDPListener) OnNewCacheChangeAdded(reader reader.IRTPSReader, change *common.CacheChangeT) {
	log.Println("SPDP Message received from: ", change.WriterGUID)

	// Make sure we have an instance handle (i.e GUID)
	if change.InstanceHandle == common.KInstanceHandleUnknown {
		log.Fatalln("not impl")
	}

	guid := change.InstanceHandle.Convert2GUID()
	if change.Kind == common.KAlive {

	} else {
		reader.GetMutex().Unlock()
		if listener.parentPDP.RemoveRemoteParticipant(&guid, KRemovedParticipant) {
			reader.GetMutex().Lock()
			// All changes related with this participant have been removed from history by remove_remote_participant
			return
		}
		reader.GetMutex().Lock()
	}
	log.Println("not impl OnNewCacheChangeAdded")
}

func (listener *PDPListener) getKey(achange *common.CacheChangeT) bool {
	return message.ReadInstanceHandleFromCdrMsg(achange, policy.KPidParticipantGUID)
}

func newPDPListener(parent protocol.IPDP) *PDPListener {
	return &PDPListener{}
}
