package participant

import "github.com/yeren0143/DDS/fastrtps/rtps/builtin/data"

type DiscoveryStatus uint8

const (
	KDiscoveryParticipant DiscoveryStatus = iota
	KChangedQosParticipant
	KRemovedParticipant
	KDroppedParticipant
)

type ParticipantDiscoveryInfo struct {
	status DiscoveryStatus
	info   *data.ParticipantProxyData
}
