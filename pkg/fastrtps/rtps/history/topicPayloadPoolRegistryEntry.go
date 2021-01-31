package history

type topicPayloadPoolRegistryEntry struct {
	PoolForPreAllocated        *TopicPayloadPoolProxy
	PoolForPreallocatedRealloc *TopicPayloadPoolProxy
	PoolForDynamic             *TopicPayloadPoolProxy
	PoolForDynamicResuable     *TopicPayloadPoolProxy
}

func NewTopicPayloadPoolRegistryEntry() *topicPayloadPoolRegistryEntry {
	return &topicPayloadPoolRegistryEntry{}
}
