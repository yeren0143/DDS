package transport

import "log"

var _ ITransportDescriptor = (*SharedMemTransportDescriptor)(nil)

//SharedMemTransportDescriptor configure Shared memory transport
type SharedMemTransportDescriptor struct {
	transportDescriptor
	SegmentSize          uint32
	portQueueCapacity    uint32
	healthyCheckTimeoutn uint32
	rtpsDumpFile         string
}

func (descriptor *SharedMemTransportDescriptor) CreateTransport() ITransport {
	// TODO::
	log.Println("SharedMemTransportDescriptor CreateTransport not Impl !!!!")
	return (ITransport)(nil)
}

func (descriptor *SharedMemTransportDescriptor) MinSendBufferSize() uint32 {
	return uint32(0)
}

func NewSharedMemTransportDescriptor() *SharedMemTransportDescriptor {
	return &SharedMemTransportDescriptor{}
}
