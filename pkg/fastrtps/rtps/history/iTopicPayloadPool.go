package history

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
)

var _ IPayloadPool = (ITopicPayloadPool)(nil)

type ITopicPayloadPool interface {
	// IPayloadPool

	/**
	* @brief Ensures the pool has capacity to fullfill the requirements of a new history.
	*
	* @param [in]  config              The new history's pool requirements.
	* @param [in]  is_reader_history   True if the new history is for a reader. False otherwise.
	*
	* @post
	*   - If @c config.maximum_size is not zero
	*     - The maximum size of the pool is increased by @c config.maximum_size.
	*   - else
	*     - The maximum size of the pool is set to the largest representable value.
	*   - If the pool is configured for PREALLOCATED or PREALLOCATED WITH REALLOC memory policy:
	*     - The pool has at least as many elements allocated (including elements already in use)
	*       as the sum of the @c config.initial_size for all reserved writer histories
	*       plus the maximum of the @c config.initial_size for all reserved reader histories.
	 */
	ReserveHistory(cfg *PoolConfig, isReader bool) bool

	/**
	    * @brief Informs the pool that some history requirements are not longer active.
	    *
	    * The pool can release some resources that are not needed any longer.
	    *
	    * @param [in]  config              The old history's pool requirements, which are no longer active.
	    * @param [in]  is_reader_history   True if the history was for a reader. False otherwise.
	    *
	    * @pre
	    *   - If all remaining histories were reserved with non zero @c config.maximum_size
	    *      - The number of elements in use is less than
	    *        the sum of the @c config.maximum_size for all remaining histories
	    *
	    * @post
	    *   - If all remaining histories were reserved with non zero @c config.maximum_size
	    *      - The maximum size of the pool is set to
	    *        the sum of the @c config.maximum_size for all remaining histories
		*   - else
		*     - The maximum size of the pool remains the largest representable value.
		*   - If the number of allocated elements is greater than the new maximum size,
		*     the excess of elements are freed until the number of allocated elemets is equal to the new maximum.
	*/
	ReleaseHistory(cfg *PoolConfig, isReader bool) bool

	// Get the number of allocated payloads (reserved and not reserved).
	PayloadPoolAllocatedSize() uint32

	// Get the number of available payloads (not reserved).
	PayloadPoolAvailableSize() uint32

	GetPayload(size uint32, cacheChange *common.CacheChangeT) bool
	ReleasePayload(cacheChange *common.CacheChangeT) bool
}

func GetTopicPayloadPool(config IBasicPoolConfig) ITopicPayloadPool {
	if config.GetPayloadInitialSize() == 0 {
		return nil
	}

	switch config.GetMemoryPolicy() {
	case resources.KPreallocatedMemoryMode:
		return NewPreallocatedTopicPayloadPool(config.GetPayloadInitialSize(), 0)
	case resources.KPreallocatedWithReallocMemoryMode:
		return NewPreallocatedReallocTopicPayloadPool(config.GetPayloadInitialSize(), 0)
	case resources.KDynamicReserveMemoryMode:
	case resources.KDynamicReusableMemoryMode:
	}

	return nil
}
