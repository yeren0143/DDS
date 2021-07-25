package transport

//
const (
	CDefaultTTL uint8 = 1
)

//socketTransportDescriptor define configuration of transports using sockets.
type socketTransportDescriptor struct {
	transportDescriptor

	SendBufferSize uint32
	RcvBufferSize  uint32
	//Allowed interfaces in an IP string format.
	InterfaceWhiteList []string
	//Specified time to live (8bit - 255 max TTL)
	TTL uint8
}

//MinSendBufferSize implement
func (socketDescribe *socketTransportDescriptor) MinSendBufferSize() uint32 {
	return socketDescribe.SendBufferSize
}

//NewSocketTransportDescriptor create socketTransportDescriptor with default value
func NewSocketTransportDescriptor(maxMsgSize, maxInitPeersRange uint32) *socketTransportDescriptor {
	return &socketTransportDescriptor{
		transportDescriptor: *NewTransportDescriptor(maxMsgSize, maxInitPeersRange),
		RcvBufferSize:       0,
		SendBufferSize:      0,
		TTL:                 CDefaultTTL,
	}
}
