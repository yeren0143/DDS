package transport

//UDPv4TransportDescriptor ...
type UDPv4TransportDescriptor struct {
	UDPTransportDescriptor
}

//NewUDPv4TransportDescriptor create UDPv4TransportDescriptor with default value
func NewUDPv4TransportDescriptor() *UDPv4TransportDescriptor {
	return &UDPv4TransportDescriptor{
		UDPTransportDescriptor: *NewUDPTransportDescriptor(),
	}
}

//CreateTransport ...
func (descriptor *UDPv4TransportDescriptor) CreateTransport() ITransport {
	//var transport UDPv4Transport
	transport := NewUDPv4Transport(descriptor)
	transport.sendBufferSize = descriptor.SendBufferSize
	transport.rcvBufferSize = descriptor.RcvBufferSize

	return transport
}
