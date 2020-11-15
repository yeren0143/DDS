package transport

import (
	common "github.com/yeren0143/DDS/common"
)

//UDPv4Transport is a default UDPv4 implementation.
//  - Opening an output channel by passing a locator will open a socket per interface on the given port.
//     This collection of sockets constitute the "outbound channel". In other words, a channel corresponds
//     to a port + a direction.
//  - It is possible to provide a white list at construction, which limits the interfaces the transport
//     will ever be able to interact with. If left empty, all interfaces are allowed.
//  - Opening an input channel by passing a locator will open a socket listening on the given port on every
//     whitelisted interface, and join the multicast channel specified by the locator address. Hence, any locator
//     that does not correspond to the multicast range will simply open the port without a subsequent join. Joining
//     multicast groups late is supported by attempting to open the channel again with the same port + a
//     multicast address (the OpenInputChannel function will fail, however, because no new channel has been
//     opened in a strict sense).
type UDPv4Transport struct {
	UDPTransport
	kind common.LocatorEnum
}

/*OpenInputChannel Starts listening on the specified port, and if the specified address is in the
* multicast range, it joins the specified multicast group,
 */
func (udp *UDPv4Transport) OpenInputChannel(locator *common.Locator, receiver *ITransportReceiver, maxMsgSize uint32) bool {
	success := false
	return success
}
