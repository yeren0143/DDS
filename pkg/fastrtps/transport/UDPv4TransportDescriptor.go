package transport

import (
	common "github.com/yeren0143/DDS/common"
)

//UDPv4TransportDescriptor ...
type UDPv4TransportDescriptor struct {
	UDPTransportDescriptor
}

//CreateTransport ...
func (descriptor *UDPv4TransportDescriptor) CreateTransport() ITransport {
	var transport UDPv4Transport
	transport.sendBufferSize = descriptor.sendBufferSize
	transport.receiveBufferSize = descriptor.receiveBufferSize
	transport.configure = descriptor
	transport.kind = common.LocatorKindUDPv4
	return &transport
}
