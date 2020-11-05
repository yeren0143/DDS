package participant

import (
	. "attributes"
	. "common"
	. "writer"
)

type RTPSParticipant interface {
	GetGuid() GUID_t

	/**
	 * Indicate the Participant that you have discovered a new Remote Writer.
	 * This method can be used by the user to implements its own Static Endpoint
	 * Discovery Protocol
	 * @param pguid GUID_t of the discovered Writer.
	 * @param userDefinedId ID of the discovered Writer.
	 * @return True if correctly added.
	 */
	NewRemoteWriterDiscovered(pguid GUID_t, userDefinedId int16) bool

	/**
	 * Indicate the Participant that you have discovered a new Remote Reader.
	 * This method can be used by the user to implements its own Static Endpoint
	 * Discovery Protocol
	 * @param pguid GUID_t of the discovered Reader.
	 * @param userDefinedId ID of the discovered Reader.
	 * @return True if correctly added.
	 */
	NewRemoteReaderDiscovered(pguid GUID_t, userDefinedId int16) bool

	GetRTPSParticipantID() uint32

	RegisterWriter(writerRTPSWriter *RTPSWriter, topicAtt *TopicAttributes, wqos *WriterQos) bool
}
