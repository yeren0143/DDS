package history

import (
	"github.com/yeren0143/DDS/common"
)

var _ common.ICacheChangeParent = (IPayloadPool)(nil)

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

	/**
	 * @brief Assign a serialized payload to a new sample.
	 *
	 * This method will usually be called when a reader receives a whole cache change.
	 *
	 * @param [in,out] data          Serialized payload received
	 * @param [in,out] data_owner    Payload pool owning incoming data
	 * @param [in,out] cache_change  Cache change to assign the payload to
	 *
	 * @returns whether the operation succeeded or not
	 *
	 * @note @c data and @c data_owner are received as references to accommodate the case where several readers
	 * receive the same payload. If the payload has no owner, it means it is allocated on the stack of a
	 * reception thread, and a copy should be performed. The pool may decide in that case to point @c data.data
	 * to the new copy and take ownership of the payload. In that case, when the reception thread is done with
	 * the payload (after all readers have been informed of the received data), method @c release_payload will be
	 * called to indicate that the reception thread is not using the payload anymore.
	 *
	 * @warning @c data_owner can only be changed from @c nullptr to @c this. If a value different from
	 * @c nullptr is received it should be left unchanged.
	 *
	 * @warning @c data fields can only be changed when @c data_owner is @c nullptr. If a value different from
	 * @c nullptr is received all fields in @c data should be left unchanged.
	 *
	 * @pre
	 *     @li Field @c cache_change.writerGUID is not @c unknown
	 *     @li Field @c cache_change.sequenceNumber is not @c unknown
	 *
	 * @post
	 *     @li Field @c cache_change.payload_owner equals this
	 *     @li Field @c cache_change.serializedPayload.data points to a buffer of at least @c data.length bytes
	 *     @li Field @c cache_change.serializedPayload.length is equal to @c data.length
	 *     @li Field @c cache_change.serializedPayload.max_size is greater than or equal to @c data.length
	 *     @li Content of @c cache_change.serializedPayload.data is the same as @c data.data
	 */
	GetPayloadWithOwner(data *common.SerializedPayloadT, dataOwner *common.ICacheChangeParent, aChange *common.CacheChangeT) bool

	// Release a serialized payload from a sample.
	// This method will be called when a cache change is removed from a history.
	ReleasePayload(cacheChange *common.CacheChangeT) bool
}
