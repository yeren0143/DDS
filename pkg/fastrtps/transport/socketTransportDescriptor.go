package transport

//
const (
	CDefaultTTL uint8 = 1
)

//SocketTransportDescriptor define configuration of transports using sockets.
type SocketTransportDescriptor struct {
	transportDescriptor

	sendBufferSize uint32
	rcvBufferSize  uint32
	//Allowed interfaces in an IP string format.
	interfaceWhiteList []string
	//Specified time to live (8bit - 255 max TTL)
	ttl uint8
}

//NewSocketTransportDescriptor create SocketTransportDescriptor with default value
func NewSocketTransportDescriptor() *SocketTransportDescriptor {
	return &SocketTransportDescriptor{
		ttl: CDefaultTTL,
	}
}
