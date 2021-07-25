package participant

type DiscoveryStatus int

const (
	DISCOVERED_PARTICIPANT DiscoveryStatus = iota
	CHANGED_QOS_PARTICIPANT
	REMOVED_PARTICIPANT
	DROPPED_PARTICIPANT
)

type ParticipantDiscoveryInfo struct {
	status DiscoveryStatus
}
