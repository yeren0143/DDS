package history

import (
	"log"
	"math"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
)

var _ IChangePool = (*CacheChangePool)(nil)

type CacheChangePool struct {
	currentPoolSize uint32
	maxPoolSize     uint32
	memoryMode      resources.MemoryManagementPolicy
	freeCaches      []*common.CacheChangeT
	allCaches       []*common.CacheChangeT
}

func (changePool *CacheChangePool) ReserveCache() (*common.CacheChangeT, bool) {
	var cacheChange *common.CacheChangeT
	if len(changePool.freeCaches) == 0 {
		switch changePool.memoryMode {
		case resources.KPreallocatedMemoryMode:
			if !changePool.allocateGroup(changePool.currentPoolSize/10 + 10) {
				return nil, false
			}
		case resources.KPreallocatedWithReallocMemoryMode:
			if !changePool.allocateGroup(changePool.currentPoolSize/10 + 10) {
				return nil, false
			}
		case resources.KDynamicReserveMemoryMode:
			cacheChange = changePool.allocateSingle()
			return cacheChange, cacheChange != nil
		case resources.KDynamicReusableMemoryMode:
			cacheChange = changePool.allocateSingle()
			return cacheChange, cacheChange != nil
		default:
			return nil, false
		}
	}

	cacheChange = changePool.freeCaches[0]
	changePool.freeCaches = changePool.freeCaches[1:]
	return cacheChange, true
}

func (changePool *CacheChangePool) ReleaseCache(cacheChange *common.CacheChangeT) bool {
	switch changePool.memoryMode {
	case resources.KPreallocatedMemoryMode:
		changePool.returnCacheToPool(cacheChange)
	case resources.KPreallocatedWithReallocMemoryMode:
		changePool.returnCacheToPool(cacheChange)
	case resources.KDynamicReusableMemoryMode:
		changePool.returnCacheToPool(cacheChange)
	case resources.KDynamicReserveMemoryMode:
		// Find pointer in CacheChange vector, remove element, then delete it
		target := -1
		for i := 0; i < len(changePool.allCaches); i++ {
			if cacheChange == changePool.allCaches[i] {
				target = i
			}
		}
		if target != -1 {
			changePool.allCaches = append(changePool.allCaches[:target], changePool.allCaches[:target+1]...)
		} else {
			log.Println("Tried to release a CacheChange that is not logged in the Pool")
			return false
		}

		changePool.currentPoolSize--
	}
	return true
}

func (changePool *CacheChangePool) allocateSingle() *common.CacheChangeT {
	/*
	 *   In Dynamic Memory Mode CacheChanges are only allocated when they are needed.
	 *   This means when the buffer of the message receiver is copied into this struct, the size is allocated.
	 *   When the change is released and comes back to the pool, it is deallocated correspondingly.
	 *
	 *   In Preallocated mode, changes are allocated with a static maximum size and then they are dealt as
	 *   they are needed. In Dynamic mode, they are only allocated when they are needed. In Dynamic mode only
	 *   the all_caches_ vector is used, in order to keep track of all the changes that are dealt for destruction
	 *   purposes.
	 *
	 */
	var added bool
	var ch *common.CacheChangeT

	if changePool.memoryMode != resources.KPreallocatedMemoryMode &&
		changePool.memoryMode != resources.KPreallocatedWithReallocMemoryMode {
		log.Fatalln("invalid memoryMode")
	}

	if changePool.currentPoolSize < changePool.maxPoolSize {
		changePool.currentPoolSize++
		ch = common.NewCacheChangeT()
		changePool.allCaches = append(changePool.allCaches, ch)
		added = true
	}

	if !added {
		log.Fatalln("Maximum number of allowed reserved caches reached")
		return nil
	}

	return ch
}

// Returns a CacheChange to the free caches pool
func (changePool *CacheChangePool) returnCacheToPool(ch *common.CacheChangeT) {
	ch.Kind = common.KAlive
	ch.SequenceNumber.Value = 0
	ch.WriterGUID = common.KGuidUnknown
	ch.InstanceHandle.Value = [16]common.Octet{0}
	ch.IsRead = false
	ch.SourceTimestamp.Seconds = 0
	ch.SourceTimestamp.Nanosec = 0
	ch.SetFragmentSize(0, false)
	changePool.freeCaches = append(changePool.freeCaches, ch)
}

func (changePool *CacheChangePool) allocateGroup(groupSize uint32) bool {
	// This method should only called from within PREALLOCATED_MEMORY_MODE or
	// PREALLOCATED_WITH_REALLOC_MEMORY_MODE
	if changePool.memoryMode != resources.KPreallocatedMemoryMode &&
		changePool.memoryMode != resources.KPreallocatedWithReallocMemoryMode {
		log.Fatalln("invalid memoryMode")
	}
	log.Println("Allocating group of cache changes of size: ", groupSize)
	desiredSize := changePool.currentPoolSize + groupSize
	if desiredSize > changePool.maxPoolSize {
		desiredSize = changePool.maxPoolSize
		groupSize = changePool.maxPoolSize - changePool.currentPoolSize
	}

	if groupSize <= 0 {
		log.Println("Maximum number of allowed reserved caches reached")
		return false
	}

	for changePool.currentPoolSize < desiredSize {
		ch := common.NewCacheChangeT()
		changePool.allCaches = append(changePool.allCaches, ch)
		changePool.freeCaches = append(changePool.freeCaches, ch)
		changePool.currentPoolSize++
	}

	return true
}

func NewCacheChangePool(cfg *PoolConfig) *CacheChangePool {
	var cachePool CacheChangePool
	cachePool.memoryMode = cfg.MemoryPolicy

	// Common for all modes: Set the pool size and size limit
	poolSize := cfg.InitialSize
	maxPoolSize := cfg.MaximumSize

	log.Println("Creating CacheChangePool of size:", poolSize)

	cachePool.currentPoolSize = 0
	if maxPoolSize > 0 {
		if poolSize > maxPoolSize {
			cachePool.maxPoolSize = poolSize
		} else {
			cachePool.maxPoolSize = maxPoolSize
		}
	} else {
		cachePool.maxPoolSize = math.MaxUint32
	}

	if poolSize == 0 {
		poolSize = 1
	}

	switch cachePool.memoryMode {
	case resources.KPreallocatedMemoryMode:
		log.Println("Static Mode is active, preallocating memory for pool_size elements")
		cachePool.allocateGroup(poolSize)
	case resources.KPreallocatedWithReallocMemoryMode:
		log.Println(`Semi-Static Mode is active,
            preallocating memory for pool_size. Size of the cachechanges can be increased`)
		cachePool.allocateGroup(poolSize)
	case resources.KDynamicReserveMemoryMode:
		log.Println("Dynamic Mode is active, CacheChanges are allocated on request")
	case resources.KDynamicReusableMemoryMode:
		log.Println(`Semi-Dynamic Mode is active, 
        no preallocation but dynamically allocated CacheChanges are reused for future cachechanges`)
	}

	return &cachePool
}
