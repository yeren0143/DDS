package history

import (
	"log"
	"math"
	"sync"
	"unsafe"

	"dds/common"
	"dds/fastrtps/rtps/resources"
)

var _ ITopicPayloadPool = (*TopicPayloadPoolBase)(nil)
var _ IPayloadPool = (*TopicPayloadPoolBase)(nil)

type payloadNode struct {
	refCounter uint32
	dataSize   uint32
	dataIndex  uint32
	data       []common.Octet
}

func (node *payloadNode) resize(size uint32) bool {
	if size < node.dataSize {
		log.Fatalln("assert size < node.dataSize failed")
		return false
	}

	oldBuffer := node.data[:]
	node.data = make([]common.Octet, size)
	copy(node.data[:node.dataSize], oldBuffer[:node.dataSize])
	node.dataSize = size
	return true
}

func (node *payloadNode) reference() {
	node.refCounter++
}

func NewpayloadNode(size uint32) *payloadNode {
	return &payloadNode{
		refCounter: 0,
		dataSize:   size,
		dataIndex:  0,
		data:       make([]common.Octet, size),
	}
}

func dereference(data []common.Octet) bool {
	return false
}

type iTopicPayloadPoolImpl interface {
	memoryPolicy() resources.MemoryManagementPolicy
	// updateMaximumSize(config *PoolConfig, reserve bool)
	//GetPayload(size uint32, cache *common.CacheChangeT, resizeable bool) bool
}

type TopicPayloadPoolBase struct {
	maxPoolSize            uint32 // Maximum size of the pool
	infiniteHistoriesCount uint32 // Number of infinite histories reserved
	finiteMaxPoolSize      uint32 //Maximum size of the pool if no infinite histories were reserved
	freePayloads           []*payloadNode
	allPayloads            []*payloadNode
	impl                   iTopicPayloadPoolImpl
	mutex                  sync.Mutex
}

func (pool *TopicPayloadPoolBase) Reserve(minNumPayloads, size uint32) {
	if minNumPayloads > pool.maxPoolSize {
		log.Fatalln("assert (min_num_payloads <= max_pool_size_) failed")
		return
	}

	for i := len(pool.allPayloads); i < int(minNumPayloads); i++ {
		payload := pool.allocate(size)
		pool.freePayloads = append(pool.freePayloads, payload)
	}
}

func (pool *TopicPayloadPoolBase) ReserveHistory(cfg *PoolConfig, isReader bool) bool {
	if cfg.MemoryPolicy != pool.impl.memoryPolicy() {
		log.Fatalln("cfg.MemoryPolicy != pool.impl.memoryPolicy()")
	}

	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	pool.updateMaximumSize(cfg, true)

	return true
}

func (pool *TopicPayloadPoolBase) updateMaximumSize(config *PoolConfig, reserve bool) {
	if reserve {
		if config.MaximumSize == 0 {
			pool.maxPoolSize = math.MaxUint32
			pool.infiniteHistoriesCount++
		} else {
			if config.InitialSize > config.MaximumSize {
				pool.finiteMaxPoolSize += config.MaximumSize
			} else {
				pool.finiteMaxPoolSize += config.InitialSize
			}

			if pool.infiniteHistoriesCount == 0 {
				pool.maxPoolSize = pool.finiteMaxPoolSize
			}
		}
	} else {
		if config.MaximumSize == 0 {
			pool.infiniteHistoriesCount--
		} else {
			if config.InitialSize > config.MaximumSize {
				pool.finiteMaxPoolSize -= config.InitialSize
			} else {
				pool.finiteMaxPoolSize -= config.MaximumSize
			}
		}

		if pool.infiniteHistoriesCount == 0 {
			pool.maxPoolSize = pool.finiteMaxPoolSize
		}

	}

}

func (pool *TopicPayloadPoolBase) shrink(maxNumPayloads uint32) bool {
	for maxNumPayloads < uint32(len(pool.allPayloads)) {
		length := len(pool.freePayloads)
		payload := pool.freePayloads[length-1]
		pool.freePayloads = pool.freePayloads[:length-1]

		// Find data in allPayloads, remove element, then delete it
		length = len(pool.allPayloads)
		pool.allPayloads[payload.dataIndex] = pool.allPayloads[length-1]
		pool.allPayloads[length-1].dataIndex = payload.dataIndex
		pool.allPayloads = pool.allPayloads[:length-1]
	}
	return true
}

func (pool *TopicPayloadPoolBase) GetPayload(size uint32, cacheChange *common.CacheChangeT) bool {
	return pool.getPayload(size, cacheChange, false)
}

func (pool *TopicPayloadPoolBase) GetPayloadWithOwner(data *common.SerializedPayloadT, dataOwner *common.ICacheChangeParent,
	aChange *common.CacheChangeT) bool {
	if aChange.WriterGUID == common.KGuidUnknown {
		log.Panic("aChange.WriterGUID == common.KGuidUnknown")
	}

	if aChange.SequenceNumber == common.KSequenceNumberUnknown {
		log.Panic("aChange.SequenceNumber == common.KSequenceNumberUnknown")
	}

	if *dataOwner == pool {
		// TODO:
		log.Panic("not impl")
	} else {
		if pool.GetPayload(data.Length, aChange) {
			if !aChange.SerializedPayload.Copy(data, true) {
				pool.ReleasePayload(aChange)
				return false
			}

			if dataOwner == nil {
				*dataOwner = pool
				data.Data = aChange.SerializedPayload.Data
			}
			return true
		}
	}

	return false
}

func (pool *TopicPayloadPoolBase) ReleasePayload(cacheChange *common.CacheChangeT) bool {
	if cacheChange.PayloadOwner() != pool {
		log.Fatalln("cacheChange.PayloadOwner() != pool")
	}

	if dereference(cacheChange.SerializedPayload.Data) {
		pool.mutex.Lock()
		defer pool.mutex.Unlock()
		node := (*payloadNode)(unsafe.Pointer(&cacheChange.SerializedPayload.Data))
		payload := pool.allPayloads[node.dataIndex]
		pool.freePayloads = append(pool.freePayloads, payload)
	}

	cacheChange.SerializedPayload.Length = 0
	cacheChange.SerializedPayload.Pos = 0
	cacheChange.SerializedPayload.MaxSize = 0
	cacheChange.SerializedPayload.Data = nil
	cacheChange.SetPayloadOwner(nil)

	return true
}

func (pool *TopicPayloadPoolBase) ReleaseHistory(cfg *PoolConfig, isReader bool) bool {
	if cfg.MemoryPolicy != pool.impl.memoryPolicy() {
		return false
	}

	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	pool.updateMaximumSize(cfg, false)

	return pool.shrink(pool.maxPoolSize)
}

func (pool *TopicPayloadPoolBase) PayloadPoolAllocatedSize() uint32 {
	return uint32(len(pool.allPayloads))
}

func (pool *TopicPayloadPoolBase) PayloadPoolAvailableSize() uint32 {
	return uint32(len(pool.freePayloads))
}

func (pool *TopicPayloadPoolBase) allocate(size uint32) *payloadNode {
	if len(pool.allPayloads) >= int(pool.maxPoolSize) {
		log.Fatalln("Maximum number of allowed reserved payloads reached")
		return nil
	}

	payload := NewpayloadNode(size)
	payload.dataIndex = uint32(len(pool.allPayloads))
	pool.allPayloads = append(pool.allPayloads, payload)
	return payload
}

func (pool *TopicPayloadPoolBase) getPayload(size uint32, cacheChange *common.CacheChangeT, resizeable bool) bool {
	var payload *payloadNode
	pool.mutex.Lock()
	if len(pool.freePayloads) == 0 {
		payload = pool.allocate(size)
		if payload == nil {
			pool.mutex.Unlock()
			cacheChange.SerializedPayload.Data = nil
			cacheChange.SerializedPayload.MaxSize = 0
			cacheChange.SetPayloadOwner(nil)
			return false
		}
	} else {
		length := len(pool.freePayloads)
		payload = pool.freePayloads[length-1]
		pool.freePayloads = pool.freePayloads[:length-1]
	}

	// Resize if needed
	if resizeable && size > payload.dataSize {
		if !payload.resize(size) {
			// Failed to resize, but we can still keep it for later.
			pool.freePayloads = append(pool.freePayloads, payload)
			pool.mutex.Unlock()
			log.Fatalln("Failed to resize the payload")
			cacheChange.SerializedPayload.Data = nil
			cacheChange.SerializedPayload.MaxSize = 0
			cacheChange.SetPayloadOwner(nil)
			return false
		}
	}

	pool.mutex.Unlock()
	payload.reference()
	cacheChange.SerializedPayload.Data = payload.data
	cacheChange.SerializedPayload.MaxSize = payload.dataSize
	cacheChange.SetPayloadOwner(pool)

	return true
}

func NewTopicPayloadPoolBase(payloadPoolImpl iTopicPayloadPoolImpl) *TopicPayloadPoolBase {
	return &TopicPayloadPoolBase{impl: payloadPoolImpl}
}
