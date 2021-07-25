package common

type MatchingStatus = uint8

const (
	KMatchedMatching MatchingStatus = iota
	KRemovedMatching
)

type MatchingInfo struct {
	Status             MatchingStatus
	RemoteEndpointGuid GUIDT
}
