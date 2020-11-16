package transport

import (
	"net"
	"syscall"

	common "github.com/yeren0143/DDS/common"
	utils "github.com/yeren0143/DDS/fastrtps/utils"
)

//UDPTransport ...
type UDPTransport struct {
	ITransport
	configure        *UDPv4TransportDescriptor
	sendBufferSize   uint32
	rcvBufferSize    uint32
	currentInterface []*utils.InfoIP
}

//CloseInputChannel Removes the listening socket for the specified port.
func (transport *UDPTransport) CloseInputChannel(*common.Locator) bool {
	return true
}

//DoInputLocatorsMatch Reports whether Locators correspond to the same port.
func (transport *UDPTransport) DoInputLocatorsMatch(*common.Locator, *common.Locator) bool {
	return true
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

	return uint32(sendBuffer), uint32(rcvBuffer)

	//return 106496, 106496
}

// Init impl
func (transport *UDPTransport) Init() bool {
	if transport.sendBufferSize == 0 || transport.rcvBufferSize == 0 {
		sendSize, rcvSize := getUDPBufferSize()

		if sendSize < CMinimumSocketBuffer {
			sendSize = CMinimumSocketBuffer
		}
		transport.sendBufferSize = sendSize
		transport.configure.sendBufferSize = sendSize

		if rcvSize < CMinimumSocketBuffer {
			rcvSize = CMinimumSocketBuffer
		}
		transport.rcvBufferSize = rcvSize
		transport.configure.rcvBufferSize = rcvSize
	}

	if transport.configure.maxMessageSize > CMaximumMessageSize {
		panic("maxMessageSize cannot be greatorthan 65000")
	}

	if transport.configure.maxMessageSize > transport.configure.sendBufferSize {
		panic("maxMessageSize cannot be greator than send buffer size")
	}

	if transport.configure.maxMessageSize > transport.configure.rcvBufferSize {
		panic("maxMessageSize cannot be greator than receive buffer size")
	}

	transport.getIPs()

	return true
}

func (transport *UDPTransport) getIPs() {
	ips, _ := utils.GetIPs(false)
	for _, ip := range ips {
		if ip.Type == utils.CIP4 || ip.Type == utils.CIP4Local {
			ip.Locator.Kind = common.LocatorKindUDPv4
			transport.currentInterface = append(transport.currentInterface, ip)
		}
	}
}

func (transport *UDPTransport) configuration() *UDPv4TransportDescriptor {
	return transport.configure
}
