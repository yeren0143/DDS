package history

import (
	"dds/common"
	"dds/fastrtps/rtps/resources"
)

var _ iTopicPayloadPoolImpl = (*PreallocatedTopicPayloadPool)(nil)

type PreallocatedTopicPayloadPool struct {
	TopicPayloadPoolBase
	payloadSize     uint32
	minimumPoolSize uint32
}

func NewPreallocatedTopicPayloadPool(payloadSize, minimumPoolSize uint32) *PreallocatedTopicPayloadPool {
	// return &PreallocatedTopicPayloadPool{
	// 	payloadSize:     payloadSize,
	// 	minimumPoolSize: minimumPoolSize,
	// }
	payloadPool := PreallocatedTopicPayloadPool{
		payloadSize:     payloadSize,
		minimumPoolSize: minimumPoolSize,
	}
	payloadPool.TopicPayloadPoolBase = *NewTopicPayloadPoolBase(&payloadPool)

	return &payloadPool
}

func (payloadPool *PreallocatedTopicPayloadPool) memoryPolicy() resources.MemoryManagementPolicy {
	return resources.KPreallocatedMemoryMode
}

func (payloadPool *PreallocatedTopicPayloadPool) GetPayload(size uint32, cacheChange *common.CacheChangeT) bool {
	payloadSize := size
	if payloadSize < payloadPool.minimumPoolSize {
		payloadSize = payloadPool.minimumPoolSize
	}
	return payloadPool.getPayload(payloadSize, cacheChange, true)
}
