package network

import (
	transport "dds/fastrtps/rtps/transport"
	"testing"
)

func TestNetWorkRegister(t *testing.T) {
	var descriptor transport.UDPv4TransportDescriptor
	descriptor.SendBufferSize = 0
	descriptor.RcvBufferSize = 0

	factory := NewNetworkFactory()
	factory.RegisterTransport(&descriptor)
}
