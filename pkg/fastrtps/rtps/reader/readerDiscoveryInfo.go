package reader

type DISCOVERY_STATUS uint8

const (
	DISCOVERED_READER DISCOVERY_STATUS = iota
	CHANGED_QOS_READER
	REMOVED_READER
)

type ReaderDiscoveryInfo struct {
	Status DISCOVERY_STATUS
}
