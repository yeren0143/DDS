package message

import (
	"dds/common"
	"os"
	"sync"
)

type SendBuffersManager struct {
	mutex        sync.Mutex
	availableCV  *sync.Cond
	commonBuffer []common.Octet
	pools        []*RTPSMessageGroupT
	createdCount uint64 //Creation counter
	allowGrowing bool   // Whether we allow n_created_ to grow beyond the pool_ capacity.
}

func (manager *SendBuffersManager) Init(guidPrefix *common.GUIDPrefixT, maxMsgSize uint32) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if manager.createdCount >= uint64(len(manager.pools)) {
		return
	}

	// Single allocation for the data of all the buffers.
	// We align the payload size to the size of a pointer, so all buffers will
	// be aligned as if directly allocated.
	alignSize := uint32(8) - 1
	payloadSize := (maxMsgSize + alignSize) & (^alignSize)
	advance := payloadSize
	if os.Getenv("HAVE_SECURITY") != "true" {
		advance *= 2
	}

	dataSize := advance * (uint32(len(manager.pools)) - uint32(manager.createdCount))
	manager.commonBuffer = make([]common.Octet, dataSize)

	offset := 0
	for i := uint64(0); manager.createdCount < uint64(len(manager.pools)); i++ {
		manager.pools[i] = NewRTPSMessageGroup(manager.commonBuffer[offset:], uint32(payloadSize), guidPrefix)
		offset += int(advance)
		manager.createdCount++
	}
}

func NewSendBuffersManager(reservedSize uint64, allowGrowingBuffers bool) *SendBuffersManager {
	sendBufferMgr := SendBuffersManager{
		allowGrowing: allowGrowingBuffers,
		pools:        make([]*RTPSMessageGroupT, reservedSize),
	}
	return &sendBufferMgr
}
