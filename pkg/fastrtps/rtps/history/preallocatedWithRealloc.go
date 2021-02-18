package history

import (
	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
)

var _ topicPayloadPoolImpl = (*PreallocatedTopicPayloadPool)(nil)

type PreallocatedReallocTopicPayloadPool struct {
	TopicPayloadPool
	minPayloadSize  uint32
	minimumPoolSize uint32
}

func (payloadPool *PreallocatedReallocTopicPayloadPool) memoryPolicy() resources.MemoryManagementPolicy {
	return resources.KPreallocatedWithReallocMemoryMode
}

func (pool *PreallocatedReallocTopicPayloadPool) ReleaseHistory(config *PoolConfig, isReader bool) bool {
	pool.mutex.Lock()
	pool.minimumPoolSize -= config.InitialSize
	pool.mutex.Unlock()

	return pool.TopicPayloadPool.ReleaseHistory(config, isReader)
}

func (pool *PreallocatedReallocTopicPayloadPool) ReserveHistory(config *PoolConfig, isReader bool) bool {
	if !pool.TopicPayloadPool.ReserveHistory(config, isReader) {
		return false
	}

	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	pool.minimumPoolSize += config.InitialSize
	pool.Reserve(pool.minimumPoolSize, pool.minPayloadSize)

	return true
}

func NewPreallocatedReallocTopicPayloadPool(payloadSize, poolSize uint32) *PreallocatedReallocTopicPayloadPool {
	var pool PreallocatedReallocTopicPayloadPool
	pool.minPayloadSize = payloadSize
	pool.minimumPoolSize = poolSize
	pool.impl = &pool

	return &pool
}
