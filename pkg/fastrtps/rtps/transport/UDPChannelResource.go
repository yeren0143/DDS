package transport

import (
	"github.com/yeren0143/DDS/common"
	//"golang.org/x/net/ipv4"
	"log"
	"sync/atomic"
	"syscall"
)

//UDPChannelResource is container of UDP transport
type UDPChannelResource struct {
	ChannelResource
	messageReceiver ITransportReceiver
	transport       IUDPTransport
	udpConn         int
	//Conn       *ipv4.PacketConn
	sinterface string
}

func (resource *UDPChannelResource) Receive(receiverMsg *common.CDRMessage, remoteLocator *common.Locator) bool {
	numBytes, sockaddr, err := syscall.Recvfrom(resource.udpConn, receiverMsg.Buffer, 0)

	if err != nil {
		log.Fatalf("receive udp channel resource falied: (%v)", err)
	}

	if numBytes > 0 {
		header := []rune("EPRORTPSCLOSE")
		if numBytes == 13 && string(receiverMsg.Buffer[:13]) == string(header) {
			return false
		}

		sockaddr4, ok := sockaddr.(*syscall.SockaddrInet4)
		log.Printf("receive msg from addr(%v), port(%v)", sockaddr4.Addr, sockaddr4.Port)
		if ok {
			receiverMsg.Length = uint32(numBytes)
			//resource.transport.endPointToLocator(remoteAddr, remoteLocator)
			remoteLocator.Address[12] = sockaddr4.Addr[0]
			remoteLocator.Address[13] = sockaddr4.Addr[1]
			remoteLocator.Address[14] = sockaddr4.Addr[2]
			remoteLocator.Address[15] = sockaddr4.Addr[3]
			remoteLocator.Port = uint32(sockaddr4.Port)
			remoteLocator.Kind = common.KLocatorKindUDPv4
			return true
		}
	}

	return false
}

func NewUDPChannelResource(transport IUDPTransport, conn int, maxMsgSize uint32, locator *common.Locator,
	sInterface string, receiver ITransportReceiver) *UDPChannelResource {
	udp := UDPChannelResource{}
	atomic.StoreUint32(&udp.alive, 1)
	udp.messageBuffer = common.NewCDRMessage(maxMsgSize)
	udp.messageReceiver = receiver
	udp.udpConn = conn
	udp.sinterface = sInterface
	udp.transport = transport

	go func(inputLocator *common.Locator) {
		var remoteLocator common.Locator
		for udp.Alive() {
			msg := udp.MessageBuffer()
			if !udp.Receive(msg, &remoteLocator) {
				continue
			}

			if udp.messageReceiver != nil {
				udp.messageReceiver.OnDataReceived(msg.Buffer, msg.Length, inputLocator, &remoteLocator)
			} else if udp.Alive() {
				log.Printf("Received Message, but no receiver attached")
			}
		}

		udp.messageReceiver = nil

	}(locator)

	return &udp
}
