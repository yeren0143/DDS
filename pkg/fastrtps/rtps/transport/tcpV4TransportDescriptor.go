package transport

import (
	common "dds/common"
)

type TCPv4TransportDescriptor struct {
	TCPTransportDescriptor
	wanAddr [4]common.Octet
}

func (descriptor *TCPv4TransportDescriptor) CreateTransport() ITransport {
	return newTCPv4Transport(descriptor)
}
