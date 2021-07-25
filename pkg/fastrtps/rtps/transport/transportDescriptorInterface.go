package transport

//ITransportDescriptor base class for the data type used to define transport configuration.
type ITransportDescriptor interface {

	// Factory method pattern. It will create and return a TransportInterface
	// corresponding to this descriptor. This provides an interface to the NetworkFactory
	// to create the transports without the need to know about their type
	CreateTransport() ITransport

	//Returns the minimum size required for a send operation
	MinSendBufferSize() uint32

	//Returns the maximum size expected for received messages.
	MaxMessageSize() uint32

	MaxInitialPeersRange() uint32

	// GetConfiguration() ITransportDescriptor
}

type transportDescriptor struct {
	maxMessageSize       uint32
	maxInitialPeersRange uint32
}

func (descriptor *transportDescriptor) MaxMessageSize() uint32 {
	return descriptor.maxMessageSize
}

func (descriptor *transportDescriptor) MaxInitialPeersRange() uint32 {
	return descriptor.maxInitialPeersRange
}

func NewTransportDescriptor(maxMsgSize, maxInitPeersRange uint32) *transportDescriptor {
	return &transportDescriptor{
		maxMessageSize:       maxMsgSize,
		maxInitialPeersRange: maxInitPeersRange,
	}
}
