package transport

import (
	common "github.com/yeren0143/DDS/common"
)

type TCPv4TransportDescriptor struct {
	TCPTransportDescriptor
	wanAddr [4]common.Octet
}

func (descriptor *TCPv4TransportDescriptor) CreateTransport() ITransport {
	return newTCPv4Transport(descriptor)
}
