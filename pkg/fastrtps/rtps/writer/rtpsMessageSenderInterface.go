package writer

import (
	"github.com/yeren0143/DDS/common"
)

// IRTPSMessageSender is an interface used in RTPSMessageGroup to handle destinations management
// and message sending
type IRTPSMessageSender interface {
	// Check if the destinations managed by this sender interface have changed.
	DestinationHaveChanged() bool

	// Get a GUID prefix representing all destinations.
	// When all the destinations share the same prefix (i.e. belong to the same participant)
	// that prefix is returned. When there are no destinations, or they belong to different
	// participants, c_GuidPrefix_Unknown is returned.
	DestinationGuidPrefix() common.GUIDPrefixT

	// Get the GUID prefix of all the destination participants.
	// a const reference to a vector with the GUID prefix of all destination participants.
	RemoteParticipants() []common.GUIDPrefixT

	// Get the GUID of all destinations.
	RemoteGUIDs() []common.GUIDT

	Send(msg *common.CDRMessage, maxBlockingTimePoint common.Time) bool
}
