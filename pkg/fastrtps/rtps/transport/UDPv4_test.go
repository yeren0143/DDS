package transport

import (
	"fmt"
	"testing"
)

func TestUDPv4Transport(t *testing.T) {
	var descriptor UDPv4TransportDescriptor
	descriptor.SendBufferSize = 0
	descriptor.RcvBufferSize = 0

	transport := descriptor.CreateTransport()
	transport.Init()
	fmt.Println("ip4:")
}

// func TestUDPv4Description(t *testing.T) {
// 	udpv4 := UDPv4TransportDescriptor{}
// 	var transportDescription ITransportDescriptor
// 	transportDescription = &udpv4
// }

// func TestUDPv4Transport(t *testing.T) {

// 	udpv4 := UDPv4Transport{}
// 	var transportInterface ITransport
// 	transportInterface = &udpv4
// }
