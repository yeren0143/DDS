package transport

import (
	common "github.com/yeren0143/DDS/common"
)

//UDPTransport ...
type UDPTransport struct {
	ITransport
	configure         ITransportDescriptor
	sendBufferSize    uint32
	receiveBufferSize uint32
}

//CloseInputChannel Removes the listening socket for the specified port.
func (transport *UDPTransport) CloseInputChannel(*common.Locator) bool {
	return true
}

//DoInputLocatorsMatch Reports whether Locators correspond to the same port.
func (transport *UDPTransport) DoInputLocatorsMatch(*common.Locator, *common.Locator) bool {
	return true
}

// Init impl
func (transport *UDPTransport) Init() bool {
	return true
}

func (transport *UDPTransport) configuration() ITransportDescriptor {
	return transport.configure
}
