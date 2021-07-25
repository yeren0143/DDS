package transport

//UDPTransportDescriptor define UDP Transport configuration
type UDPTransportDescriptor struct {
	socketTransportDescriptor
	outputUDPSocket uint32
	//Whether to use non-blocking calls to send_to().
	nonBlockingSend bool
}

//NewUDPTransportDescriptor create UDPTransportDescriptor with default value
func NewUDPTransportDescriptor() *UDPTransportDescriptor {
	return &UDPTransportDescriptor{
		nonBlockingSend:           false,
		socketTransportDescriptor: *NewSocketTransportDescriptor(KMaximumMessageSize, KMaximumInitialPeersRange),
	}
}
