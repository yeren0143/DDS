package history

import (
	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
	"sync"
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

// func (register *TopicPayloadPoolRegistry) Get(topic string, config *BasicPoolConfig) ITopicPayloadPool {
// 	register.mutex.Lock()
// 	defer register.mutex.Unlock()

// 	entry, ok := register.poolMap[topic]
// 	if !ok {
// 		entry = *NewTopicPayloadPoolRegistryEntry()
// 		register.poolMap[topic] = entry
// 	}

// 	switch config.MemoryPolicy {
// 	case resources.KPreallocatedMemoryMode:
// 	case resources.KPreallocatedWithReallocMemoryMode:
// 	case resources.KDynamicReserveMemoryMode:
// 	case resources.KDynamicReusableMemoryMode:
// 	}

// 	return nil
// }
