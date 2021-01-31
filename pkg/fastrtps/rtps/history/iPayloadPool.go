package history

import (
	"github.com/yeren0143/DDS/common"
)

var _ common.ICacheChangeOwner = (IPayloadPool)(nil)

// IPayloadPool is a interface for classes responsible of serialized payload management.
type IPayloadPool interface {
	/**
	 * @brief Get a serialized payload for a new sample.
	 *
	 * This method will usually be called in one of the following situations:
	 *     @li When a writer creates a new cache change
	 *     @li When a reader receives the first fragment of a cache change
	 *
	 * In both cases, the received @c size will be for the whole serialized payload.
	 *
	 * @param [in]     size          Number of bytes required for the serialized payload.
	 *                               Should be greater than 0.
	 * @param [in,out] cache_change  Cache change to assign the payload to
	 *
	 * @returns whether the operation succeeded or not
	 *
	 * @pre Fields @c writerGUID and @c sequenceNumber of @c cache_change are either:
	 *     @li Both equal to @c unknown (meaning a writer is creating a new change)
	 *     @li Both different from @c unknown (meaning a reader has received the first fragment of a cache change)
	 *
	 * @post
	 *     @li Field @c cache_change.payload_owner equals this
	 *     @li Field @c serializedPayload.data points to a buffer of at least @c size bytes
	 *     @li Field @c serializedPayload.max_size is greater than or equal to @c size
	 */
	GetPayload(size uint32, cacheChange *common.CacheChangeT) bool

	// Release a serialized payload from a sample.
	// This method will be called when a cache change is removed from a history.
	ReleasePayload(cacheChange *common.CacheChangeT) bool
}
