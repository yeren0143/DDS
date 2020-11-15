package transport

import (
	"fmt"
	"testing"
)

func TestUDPv4Transport(t *testing.T) {
	var descriptor UDPv4TransportDescriptor
	descriptor.sendBufferSize = 0
	descriptor.receiveBufferSize = 0

	transport := descriptor.CreateTransport()
	transport.Init()
	fmt.Println(transport)
}
