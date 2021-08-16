package writer

import "dds/common"

// Abstract class IReaderDataFilter that acts as virtual interface for data filters in ReaderProxy.
type IReaderDataFilter interface {
	// This method checks whether a CacheChange_t is relevant for the remote reader
	// This callback should return always the same result given the same arguments
	IsRelevant(aChange *common.CacheChangeT, readerGuid *common.GUIDT) bool
}
