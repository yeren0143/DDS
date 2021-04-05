package network

import (
	"sync"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/message"
	"github.com/yeren0143/DDS/fastrtps/rtps/transport"
)

type receiverResourceCleanupFunc func()
type locatorMapsToManagedChannelFunc = func(*common.Locator) bool

// ReceiverResource is RAII object that encapsulates the Receive operation over one channel in an unknown transport.
// A Receiver resource is always univocally associated to a transport channel; the
// act of constructing a Receiver Resource opens the channel and its destruction
// closes it.
type ReceiverResource struct {
	mutex                       sync.Mutex
	receiver                    *message.Receiver
	MaxMessageSize              uint32
	Valid                       bool
	Cleanup                     receiverResourceCleanupFunc
	LocatorMapsToManagedChannel *locatorMapsToManagedChannelFunc
}

var _ transport.ITransportReceiver = (*ReceiverResource)(nil)

func (resource *ReceiverResource) OnDataReceived(data []common.Octet, length uint32, localLocator *common.Locator,
	remoteLocator *common.Locator) {
	resource.mutex.Lock()
	defer resource.mutex.Unlock()

	if resource.receiver != nil {
		msg := common.NewCDRMessage(length)
		msg.Wraps = false
		copy(msg.Buffer, data)
		msg.Length = length
		msg.MaxSize = msg.Length
		msg.ReservedSize = msg.Length

		// TODO: Should we unlock in case UnregisterReceiver is called from callback ?
		resource.receiver.ProcessCDRMsg(remoteLocator, msg)
	}
}

func (resource *ReceiverResource) SupportsLocator(locator *common.Locator) bool {
	if resource.LocatorMapsToManagedChannel != nil {
		return (*resource.LocatorMapsToManagedChannel)(locator)
	}
	return false
}

func (resource *ReceiverResource) RegisterReceiver(rcv *message.Receiver) {
	resource.mutex.Lock()
	defer resource.mutex.Unlock()
	if rcv != nil {
		resource.receiver = rcv
	}
}

func NewReceiverResource(inTransport transport.ITransport, locator *common.Locator,
	maxRecvBufferSize uint32) *ReceiverResource {
	resource := &ReceiverResource{}
	transport := inTransport
	resource.MaxMessageSize = maxRecvBufferSize
	resource.Valid = transport.OpenInputChannel(locator, resource, maxRecvBufferSize)
	if !resource.Valid {
		return resource // Invalid resource to be discarded by the factory.
	}

	resource.Cleanup = func() {
		transport.CloseInputChannel(locator)
	}

	changeFunc := func(locatorToCheck *common.Locator) bool {
		return transport.DoInputLocatorsMatch(locator, locatorToCheck)
	}
	resource.LocatorMapsToManagedChannel = &changeFunc

	return resource
}
