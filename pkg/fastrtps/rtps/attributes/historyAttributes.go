package attributes

import (
	"dds/fastrtps/rtps/resources"
)

/**
 * Class HistoryAttributes, to specify the attributes of a WriterHistory or a ReaderHistory.
 * This class is only intended to be used with the RTPS API.
 * The Publsiher-Subscriber API has other fields to define this values (HistoryQosPolicy and ResourceLimitsQosPolicy).
 */
type HistoryAttributes struct {
	MemoryPolicy resources.MemoryManagementPolicy

	// Maximum payload size of the history, default value 500.
	PayloadMaxSize uint32

	// Number of the initial Reserved Caches, default value 500.
	InitialReservedCaches uint32

	/**
	* Maximum number of reserved caches. Default value is 0 that indicates to keep reserving until something
	* breaks.
	 */
	MaximumReservedCaches uint32
}

var KDefaultHistoryAttributes = HistoryAttributes{}
