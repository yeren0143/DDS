package history

import (
	"dds/fastrtps/rtps/resources"
)

var _ iTopicPayloadPoolImpl = (*PreallocatedTopicPayloadPool)(nil)

type PreallocatedReallocTopicPayloadPool struct {
	TopicPayloadPoolBase
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

	return pool.TopicPayloadPoolBase.ReleaseHistory(config, isReader)
}

func (pool *PreallocatedReallocTopicPayloadPool) ReserveHistory(config *PoolConfig, isReader bool) bool {
	if !pool.TopicPayloadPoolBase.ReserveHistory(config, isReader) {
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
