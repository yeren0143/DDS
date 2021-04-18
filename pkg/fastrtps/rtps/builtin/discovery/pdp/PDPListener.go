package pdp

import (
	"log"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/core/policy"
	"github.com/yeren0143/DDS/fastrtps/message"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/data"
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

	// Temporary data to avoid reallocations.
	// This should be always accessed with the pdp_reader lock taken
	tempParticipantData data.ParticipantProxyData
}

func (listener *PDPListener) OnNewCacheChangeAdded(reader reader.IRTPSReader, achange *common.CacheChangeT) {
	log.Println("SPDP Message received from: ", achange.WriterGUID)

	// Make sure we have an instance handle (i.e GUID)
	if achange.InstanceHandle == common.KInstanceHandleUnknown {
		if !listener.getKey(achange) {
			log.Println("Problem getting the key of the change, removing")
			log.Fatalln("not Impl")
			// listener.parentPDP.GetPDPReaderHistory().RemoveChange(change)
			return
		}
	}

	guid := achange.InstanceHandle.Convert2GUID()
	if achange.Kind == common.KAlive {
		// Ignore announcement from own RTPSParticipant
		if guid == *listener.parentPDP.GetRTPSParticipant().GetGuid() {
			log.Println("Message from own RTPSParticipant, removing")
			listener.parentPDP.GetPDPReaderHistory().RemoveChange(achange)
			return
		}

		writerGuid := achange.WriterGUID

		// Release reader lock to avoid ABBA lock. PDP mutex should always be first.
		// Keep change information on local variables to check consistency later
		seqNum := achange.SequenceNumber
		reader.GetMutex().Unlock()
		listener.parentPDP.GetMutex().Lock()
		reader.GetMutex().Lock()

		// If change is not consistent, it will be processed on the thread that has overriten it
		if achange.Kind != common.KAlive || seqNum != achange.SequenceNumber || achange.WriterGUID != writerGuid {
			return
		}

		// Access to temp_participant_data_ is protected by reader lock

		// Load information on temp_participant_data_
		msg := common.NewCDRMessageWithPayload(&achange.SerializedPayload)
		listener.tempParticipantData.Clear()

		participant := listener.parentPDP.GetRTPSParticipant()
		if listener.tempParticipantData.ReadFromCDRMessage(msg, true, participant.NetworkFactory(),
			participant.HasShmTransport()) {
			log.Fatalln("not impl")
		}

		listener.parentPDP.GetMutex().Unlock()

	} else {
		reader.GetMutex().Unlock()
		if listener.parentPDP.RemoveRemoteParticipant(&guid, KRemovedParticipant) {
			reader.GetMutex().Lock()
			// All changes related with this participant have been removed from history by remove_remote_participant
			return
		}
		reader.GetMutex().Lock()
	}

	//Remove change form history.
	listener.parentPDP.GetPDPReaderHistory().RemoveChange(achange)
}

func (listener *PDPListener) getKey(achange *common.CacheChangeT) bool {
	return message.ReadInstanceHandleFromCdrMsg(achange, policy.KPidParticipantGUID)
}

func newPDPListener(parent protocol.IPDP) *PDPListener {
	allocAtt := parent.GetRTPSParticipant().GetAttributes().Allocation
	participantData := data.NewParticipantProxyData(allocAtt)
	return &PDPListener{
		parentPDP:           parent,
		tempParticipantData: *participantData,
	}
}
