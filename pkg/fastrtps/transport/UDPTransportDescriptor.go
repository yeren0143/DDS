package transport

//UDPTransportDescriptor define UDP Transport configuration
type UDPTransportDescriptor struct {
	SocketTransportDescriptor
	outputUDPSocket uint16
	//Whether to use non-blocking calls to send_to().
	nonBlockingSend bool
}

//NewUDPTransportDescriptor create UDPTransportDescriptor with default value
func NewUDPTransportDescriptor() *UDPTransportDescriptor {
	return &UDPTransportDescriptor{
		nonBlockingSend:           false,
		SocketTransportDescriptor: *NewSocketTransportDescriptor(),
	}
}
