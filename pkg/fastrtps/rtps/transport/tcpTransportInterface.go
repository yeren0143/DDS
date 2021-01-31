package transport

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/utils"
	"sync"
)

type ITCPTransport interface {
	CompareLocatorIP(lh, rh *common.Locator) bool
	CompareLocatorIPAndPort(lh, rh *common.Locator) bool
	FillLocalIP(loc *common.Locator)

	SetReceiveBufferSize(size uint32)
	SetSendBufferSize(size uint32)
}

type tcpTransport struct {
	transportKind     int8
	socketsMapMutex   sync.Mutex
	receiverResources map[uint16]ITransportReceiver
	currentInterfaces []utils.InfoIP
}

func (transport *tcpTransport) isInputPortOpen(port uint16) bool {
	transport.socketsMapMutex.Lock()
	defer transport.socketsMapMutex.Unlock()
	_, ok := transport.receiverResources[port]
	return ok
}

func (transport *tcpTransport) Kind() int8 {
	return transport.transportKind
}

func (transport *tcpTransport) IsLocatorSupported(locator *common.Locator) bool {
	return locator.Kind == transport.transportKind
}

func (transport *tcpTransport) IsInputChannelOpen(locator *common.Locator) bool {
	result := transport.IsLocatorSupported(locator)
	result = result && transport.isInputPortOpen(utils.GetLogicalPort(locator))
	return result
}

func (transport *tcpTransport) OpenOutputChannel(senderList SenderSourceList, locator *common.Locator) bool {
	if transport.IsLocatorSupported(locator) == false {
		return false
	}

	success := false
	logicalPort := utils.GetLogicalPort(locator)

	if logicalPort != 0 {
		physicalLocator := utils.ToPhysicalLocator(locator)

		// We try to find a SenderResource that can be reuse to this locator.
		// Note: This is done in this level because if we do in NetworkFactory level, we have to mantain what transport
		// already reuses a SenderResource.
		for _, resource := range senderList {
			tcpResource := resource.(*tcpSenderResource)
			if tcpResource != nil && physicalLocator == tcpResource.channel.locator {
				return true
			}
		}
	}

	return success
}

func (transport *tcpTransport) RemoteToMainLocal(remote *common.Locator) *common.Locator {
	if transport.IsLocatorSupported(remote) == false {
		return nil
	}

	mainLocal := remote
	mainLocal.SetInvalidAddress()
	return mainLocal
}

func (transport *tcpTransport) CloseInputChannel(locator *common.Locator) bool {
	return false
}

func (transport *tcpTransport) DoInputLocatorsMatch(left, right *common.Locator) bool {
	return utils.GetPhysicalPort(left) == utils.GetPhysicalPort(right)
}

func (transport *tcpTransport) OpenInputChannel(locator *common.Locator, receiver ITransportReceiver, maxMsgSize uint32) bool {
	return false
}
