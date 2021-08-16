package pdp

import "dds/fastrtps/rtps/builtin/data"

type DiscoveryStatus = uint8

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
