package reader

import "dds/common"

//  Class RTPSReader, manages the reception of data from its matched writers.
type ReaderHistoryState struct {
	// Physical GUID to persistence GUID map
	PersistenceGUIDMap map[common.GUIDT]common.GUIDT
	// Persistence GUID count map
	PersistenceGUIDCount map[common.GUIDT]uint16
	// Information about max notified change
	HistoryRecord map[common.GUIDT]common.SequenceNumberT
}

func NewReaderHistoryState(initialWritersAllocation uint32) *ReaderHistoryState {
	var rState ReaderHistoryState
	rState.PersistenceGUIDMap = make(map[common.GUIDT]common.GUIDT, initialWritersAllocation)
	rState.PersistenceGUIDCount = make(map[common.GUIDT]uint16, initialWritersAllocation)
	rState.HistoryRecord = make(map[common.GUIDT]common.SequenceNumberT, initialWritersAllocation)
	return &rState
}
