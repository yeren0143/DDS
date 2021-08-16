package history

import (
	"dds/common"
	"dds/fastrtps/rtps/resources"
)

var _ ITopicPayloadPool = (*TopicPayloadPoolProxy)(nil)

/**
 * Proxy class that adds the topic name to a ITopicPayloadPool, so we can look-up
 * the corresponding entry in the registry when releasing the pool.
 */
type TopicPayloadPoolProxy struct {
	topicName string
	policy    resources.MemoryManagementPolicy
	innerPool ITopicPayloadPool
}

func (proxy *TopicPayloadPoolProxy) GetPayload(size uint32, cacheChange *common.CacheChangeT) bool {
	return proxy.innerPool.GetPayload(size, cacheChange)
}

func (proxy *TopicPayloadPoolProxy) GetPayloadWithOwner(data *common.SerializedPayloadT, dataOwner *common.ICacheChangeParent,
	aChange *common.CacheChangeT) bool {
	return proxy.innerPool.GetPayloadWithOwner(data, dataOwner, aChange)
}

func (proxy *TopicPayloadPoolProxy) ReleasePayload(cacheChange *common.CacheChangeT) bool {
	return proxy.innerPool.ReleasePayload(cacheChange)
}

func (proxy *TopicPayloadPoolProxy) ReserveHistory(cfg *PoolConfig, isReader bool) bool {
	return proxy.innerPool.ReserveHistory(cfg, isReader)
}

func (proxy *TopicPayloadPoolProxy) ReleaseHistory(cfg *PoolConfig, isReader bool) bool {
	return proxy.innerPool.ReleaseHistory(cfg, isReader)
}

func (proxy *TopicPayloadPoolProxy) PayloadPoolAllocatedSize() uint32 {
	return proxy.innerPool.PayloadPoolAllocatedSize()
}

func (proxy *TopicPayloadPoolProxy) PayloadPoolAvailableSize() uint32 {
	return proxy.innerPool.PayloadPoolAvailableSize()
}

func NewTopicPayloadPoolProxy(topic string, config IBasicPoolConfig) *TopicPayloadPoolProxy {
	var proxy TopicPayloadPoolProxy
	proxy.topicName = topic
	proxy.policy = config.GetMemoryPolicy()
	proxy.innerPool = GetTopicPayloadPool(config)
	return &proxy
}
