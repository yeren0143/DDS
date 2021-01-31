package transport

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/utils"
	"net"
	"runtime"
	"sync"
	"syscall"
)

type IUDPTransport interface {
	Configuration() ITransportDescriptor
	// endPointToLocator(addr *net.UDPAddr, locator *common.Locator)
}

//UDPTransport ...
type udpTransport struct {
	kind             common.LocatorEnum
	configure        UDPv4TransportDescriptor
	sendBufferSize   uint32
	rcvBufferSize    uint32
	inputMapMutex    sync.Mutex
	inputSockets     map[uint16][]*UDPChannelResource
	currentInterface []*utils.InfoIP
}

//CloseInputChannel Removes the listening socket for the specified port.
func (transport *udpTransport) CloseInputChannel(*common.Locator) bool {
	return true
}

func (transport *udpTransport) MaxRecvBufferSize() uint32 {
	return transport.configure.MaxMessageSize()
}

//DoInputLocatorsMatch Reports whether Locators correspond to the same port.
func (transport *udpTransport) DoInputLocatorsMatch(left, right *common.Locator) bool {
	return utils.GetPhysicalPort(left) == utils.GetPhysicalPort(right)
}

func getUDPBufferSize() (uint32, uint32) {
	//check system buffer sizes
	addr, _ := net.ResolveUDPAddr("udp", ":8089")
	conn, _ := net.DialUDP("udp", nil, addr)
	fd, _ := conn.File()
	sendBuffer, _ := syscall.GetsockoptInt(int(fd.Fd()), syscall.SOL_SOCKET, syscall.SO_SNDBUF)
	rcvBuffer, _ := syscall.GetsockoptInt(int(fd.Fd()), syscall.SOL_SOCKET, syscall.SO_RCVBUF)
	fd.Close()
	conn.Close()

	// On Linux, setting SO_SNDBUF or SO_RCVBUF to N actually causes the kernel
	// to set the buffer size to N*2. Linux puts additional stuff into the
	// buffers so that only about half is actually available to the application.
	// The retrieved value is divided by 2 here to make it appear as though the
	// correct value has been set.
	if runtime.GOOS == "linux" {
		sendBuffer /= 2
		rcvBuffer /= 2
	}
	return uint32(sendBuffer), uint32(rcvBuffer)
}

//IsLocatorAllowed here for occupy
func (transport *udpTransport) IsLocatorAllowed(locator *common.Locator) bool {
	panic("IsLocatorAllowed must be implement by subclass of transport")
}

// TransformRemoteLocator transforms a remote locator into a locator optimized for local communications.
// If the remote locator corresponds to one of the local interfaces, it is converted
// to the corresponding local address.
// false if the input locator is not supported/allowed by this transport, true otherwise.
func (transport *udpTransport) TransformRemoteLocator(remoteLocator *common.Locator) (*common.Locator, bool) {
	var resultLocator *common.Locator
	if transport.IsLocatorSupported(remoteLocator) {
		*resultLocator = *remoteLocator
		if transport.IsLocatorSupported(resultLocator) {
			// is_local_locator will return false for multicast addresses as well as
			// remote unicast ones.
			return resultLocator, true
		}

		// If we get here, the locator is a local unicast address
		if transport.IsLocatorAllowed(resultLocator) {
			return nil, false
		}
	}

	return nil, false
}

func (transport *udpTransport) ConfigureInitialPeerLocator(locator *common.Locator, portParams *common.PortParameters,
	domainID uint32, list *common.LocatorList) bool {

	if locator.Port == 0 {
		for i := uint32(0); i < transport.configure.maxInitialPeersRange; i++ {
			auxloc := *locator
			auxloc.Port = portParams.GetUnicastPort(domainID, i)
			list.PushBack(&auxloc)
		}
	} else {
		list.PushBack(locator)
	}
	return true
}

// RemoteToMainLocal converts a given remote locator (that is, a locator referring to a remote
// destination) to the main local locator whose channel can write to that
// destination. In this case it will return a 0.0.0.0 address on that port.
func (transport *udpTransport) RemoteToMainLocal(locator *common.Locator) *common.Locator {
	return locator
}

//IsInputChannelOpen return true if locator is in inputsockets
func (transport *udpTransport) IsInputChannelOpen(locator *common.Locator) bool {
	// transport.inputMapMutex.Lock()
	// defer transport.inputMapMutex.Unlock()

	result := transport.IsLocatorSupported(locator)
	if result {
		physicalPort := utils.GetPhysicalPort(locator)
		_, found := transport.inputSockets[physicalPort]
		result = result && found
	}
	return result
}

//IsLocatorSupported return true if locator's kind equal transport's kind
func (transport *udpTransport) IsLocatorSupported(locator *common.Locator) bool {
	return locator.Kind == transport.kind
}

func (transport *udpTransport) FillMetatrafficMulticastLocator(locator *common.Locator, wellKnownPort uint32) bool {
	if locator.Port == 0 {
		locator.Port = wellKnownPort
	}
	return true
}

func (transport *udpTransport) FillMetatrafficUnicastLocator(locator *common.Locator, unicastPort uint32) bool {
	if locator.Port == 0 {
		locator.Port = unicastPort
	}
	return true
}

func (transport *udpTransport) getIPv4s(returnLoopBack bool) []*utils.InfoIP {
	locNames, _ := utils.GetIPs(false)
	var ip4Names []*utils.InfoIP
	for _, loc := range locNames {
		if loc.Type == utils.KIP4 || loc.Type == utils.KIP4Local {
			loc.Locator.Kind = common.KLocatorKindUDPv4
			ip4Names = append(ip4Names, loc)
		}
	}

	return ip4Names
}

//OpenOutputChannel Opens a socket on the given address and port (as long as they are white listed).
func (transport *udpTransport) OpenOutputChannel(senderList SenderSourceList, locator *common.Locator) bool {
	if transport.IsLocatorSupported(locator) == false {
		return false
	}

	// We try to find a SenderResource that can be reuse to this locator.
	// Note: This is done in this level because if we do in NetworkFactory level, we have to mantain what transport
	// already reuses a SenderResource.
	for _, sendRes := range senderList {
		if _, ok := sendRes.(*udpTransport); ok {
			return true
		}
	}

	//TODO
	// port := transport.configure.outputUDPSocket
	// locNames := transport.getIPs(false)

	// if transport.isInterfaceWhitelistEmpty() {

	// }

	return true
}

//Configuration return transport's config
func (transport *udpTransport) Configuration() ITransportDescriptor {
	return &transport.configure
}

//Kind return udp locator enum kind
func (transport *udpTransport) Kind() common.LocatorEnum {
	return transport.kind
}

func (transport *udpTransport) FillUnicastLocator(locator *common.Locator, wellKnownPort uint32) bool {
	if locator.Port == 0 {
		locator.Port = wellKnownPort
	}
	return true
}
