package history

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
)

var _ topicPayloadPoolImpl = (*PreallocatedTopicPayloadPool)(nil)

type PreallocatedTopicPayloadPool struct {
	TopicPayloadPool
	payloadSize     uint32
	minimumPoolSize uint32
}

func NewPreallocatedTopicPayloadPool(payloadSize, minimumPoolSize uint32) *PreallocatedTopicPayloadPool {
	return &PreallocatedTopicPayloadPool{
		payloadSize:     payloadSize,
		minimumPoolSize: minimumPoolSize,
	}
}

func (payloadPool *PreallocatedTopicPayloadPool) memoryPolicy() resources.MemoryManagementPolicy {
	return resources.KPreallocatedWithReallocMemoryMode
}

func (payloadPool *PreallocatedTopicPayloadPool) GetPayload(size uint32, cacheChange *common.CacheChangeT) bool {
	payloadSize := size
	if payloadSize < payloadPool.minimumPoolSize {
		payloadSize = payloadPool.minimumPoolSize
	}
	return payloadPool.getPayload(payloadSize, cacheChange, true)
}
