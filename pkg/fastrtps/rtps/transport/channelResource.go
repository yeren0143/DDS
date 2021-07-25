package transport

import (
	"github.com/yeren0143/DDS/common"
	"sync/atomic"
)

type ChannelResource struct {
	messageBuffer *common.CDRMessage
	alive         uint32
}

func (resource *ChannelResource) Clear() {
	atomic.StoreUint32(&resource.alive, 0)
}

func (resource *ChannelResource) Disable() {
	atomic.StoreUint32(&resource.alive, 0)
}

func (resource *ChannelResource) Alive() bool {
	return atomic.LoadUint32(&resource.alive) > 0
}

func (resource *ChannelResource) MessageBuffer() *common.CDRMessage {
	return resource.messageBuffer
}
