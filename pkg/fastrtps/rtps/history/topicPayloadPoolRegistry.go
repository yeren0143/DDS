package history

import (
	"sync"

	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
)

type TopicPayloadPoolRegistry struct {
	poolMap map[string]topicPayloadPoolRegistryEntry
	mutex   sync.Mutex
}

var TopicPayloadPoolRegistryInstance = TopicPayloadPoolRegistry{
	poolMap: make(map[string]topicPayloadPoolRegistryEntry),
}

func (registry *TopicPayloadPoolRegistry) getTopicPayloadPoolProxy(proxy *TopicPayloadPoolProxy, topic string, config IBasicPoolConfig) *TopicPayloadPoolProxy {
	if proxy == nil {
		proxy = NewTopicPayloadPoolProxy(topic, config)
	}

	return proxy
}

func (registry *TopicPayloadPoolRegistry) Release(pool *TopicPayloadPoolProxy) {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	// A reference count of 2 means the only ones referencing the pointer are the caller and the registry.
	// This means we can release the pointer on the registry also.
	// TODO:

}

func GetTopicPayloadPoolProxy(topic string, config IBasicPoolConfig) *TopicPayloadPoolProxy {
	TopicPayloadPoolRegistryInstance.mutex.Lock()
	defer TopicPayloadPoolRegistryInstance.mutex.Unlock()

	entry, ok := TopicPayloadPoolRegistryInstance.poolMap[topic]
	if !ok {
		entry = *NewTopicPayloadPoolRegistryEntry()
		TopicPayloadPoolRegistryInstance.poolMap[topic] = entry
	}

	switch config.GetMemoryPolicy() {
	case resources.KPreallocatedMemoryMode:
		return TopicPayloadPoolRegistryInstance.getTopicPayloadPoolProxy(entry.PoolForPreAllocated, topic, config)
	case resources.KPreallocatedWithReallocMemoryMode:
		return TopicPayloadPoolRegistryInstance.getTopicPayloadPoolProxy(entry.PoolForPreallocatedRealloc, topic, config)
	case resources.KDynamicReserveMemoryMode:
		return TopicPayloadPoolRegistryInstance.getTopicPayloadPoolProxy(entry.PoolForDynamic, topic, config)
	case resources.KDynamicReusableMemoryMode:
		return TopicPayloadPoolRegistryInstance.getTopicPayloadPoolProxy(entry.PoolForDynamicResuable, topic, config)

	}

	return nil
}

func ReleaseTopicPayloadPool(pool ITopicPayloadPool) {
	topicPool := pool.(*TopicPayloadPoolProxy)

	if topicPool != nil {
		TopicPayloadPoolRegistryInstance.Release(topicPool)
	}
}
