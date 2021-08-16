package transport

import (
	"log"

	"dds/common"
	"dds/fastrtps/utils"

	//"net"
	"runtime"
	"sort"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

// var _ ITransport = (*UDPv4Transport)(nil)
var _ IUDPTransport = (*UDPv4Transport)(nil)
var _ udpTransportImpl = (*UDPv4Transport)(nil)

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
	udpTransport
	interfaceWhiteList map[string]bool
}

//NewUDPv4Transport create UDPv4 transport with UDPv4TransportDescriptor
func NewUDPv4Transport(descriptor *UDPv4TransportDescriptor) *UDPv4Transport {
	transport := &UDPv4Transport{
		udpTransport: udpTransport{
			kind:         common.KLocatorKindUDPv4,
			inputSockets: make(map[uint16][]*UDPChannelResource),
			configure:    *descriptor},
		interfaceWhiteList: make(map[string]bool),
	}
	transport.sendBufferSize = descriptor.SendBufferSize
	transport.rcvBufferSize = descriptor.RcvBufferSize
	if len(transport.interfaceWhiteList) > 0 {
		ips := getIPv4s(true)
		for _, ip := range ips {
			if _, ok := transport.interfaceWhiteList[ip.Name]; ok {
				transport.interfaceWhiteList[ip.Name] = true
			}
		}

		if len(transport.interfaceWhiteList) == 0 {
			log.Fatalln("All whitelist interfaces where filtered out")
			transport.interfaceWhiteList["192.0.2.0"] = true
		}
	}

	return transport
}

// Init impl
func (udp *UDPv4Transport) Init() bool {
	if udp.configure.SendBufferSize == 0 || udp.configure.RcvBufferSize == 0 {
		sendSize, rcvSize := getUDPBufferSize()

		if udp.configure.SendBufferSize == 0 {
			udp.configure.SendBufferSize = sendSize
			if udp.configure.SendBufferSize < KMinimumSocketBuffer {
				udp.sendBufferSize = KMinimumSocketBuffer
				udp.configure.SendBufferSize = KMinimumSocketBuffer
			}
		}

		if udp.configure.RcvBufferSize == 0 {
			udp.configure.RcvBufferSize = rcvSize
			if udp.configure.RcvBufferSize < KMinimumSocketBuffer {
				udp.rcvBufferSize = KMinimumSocketBuffer
				udp.configure.RcvBufferSize = KMinimumSocketBuffer
			}
		}
	}

	if udp.configure.maxMessageSize > KMaximumMessageSize {
		panic("maxMessageSize cannot be greatorthan 65000")
	}

	if udp.configure.maxMessageSize > udp.configure.SendBufferSize {
		panic("maxMessageSize cannot be greator than send buffer size")
	}

	if udp.configure.maxMessageSize > udp.configure.RcvBufferSize {
		panic("maxMessageSize cannot be greator than receive buffer size")
	}

	udp.GetIPs(false)
	return true
}

func (udp *UDPv4Transport) AddDefaultOutputLocator(defaultList *common.LocatorList) {
	loc := utils.CreateLocator(common.KLocatorKindUDPv4, "239.255.0.1", udp.configure.outputUDPSocket)
	defaultList.PushBack(&loc)
}

/**
* Structure info_IP with information about a specific IP obtained from a NIC.
 */
type infoIPT struct {
	ipType  utils.IPTYPE
	name    string
	dev     string
	locator *common.Locator
}

type locatorsWrap struct {
	locators []*utils.InfoIP
}

func (wrap *locatorsWrap) Len() int {
	return len(wrap.locators)
}

func (wrap *locatorsWrap) Swap(i, j int) {
	wrap.locators[i], wrap.locators[j] = wrap.locators[j], wrap.locators[i]
}

func (wrap *locatorsWrap) Less(i, j int) bool {
	return wrap.locators[i].Dev < wrap.locators[j].Dev
}

func (udp *UDPv4Transport) getIPv4UniqueInterfaces(returnLoopBack bool) []*utils.InfoIP {
	locNames := udp.getIPv4s(returnLoopBack)
	wraps := &locatorsWrap{locNames}
	sort.Sort(wraps)
	return wraps.locators
}

func (udp *UDPv4Transport) fillLocalIP(loc *common.Locator) {
	utils.SetIPv4WithIP(loc, "127.0.0.1")
	loc.Kind = common.KLocatorKindUDPv4
}

//OpenInputChannel Starts listening on the specified port, and if the specified address is in the
//multicast range, it joins the specified multicast group,
func (udp *UDPv4Transport) OpenInputChannel(locator *common.Locator, receiver ITransportReceiver, maxMsgSize uint32) bool {
	udp.inputMapMutex.Lock()
	defer udp.inputMapMutex.Unlock()

	if !udp.IsLocatorAllowed(locator) {
		return false
	}

	success := false

	if !udp.IsInputChannelOpen(locator) {
		isMulticast := utils.IsMulticast(locator)
		success = udp.openAndBindInputSockets(locator, receiver, isMulticast, maxMsgSize)
	}

	if utils.IsMulticast(locator) && udp.IsInputChannelOpen(locator) {
		locatorAddrStr := utils.ToIPv4String(locator)

		if runtime.GOOS != "windows" {
			if len(udp.interfaceWhiteList) > 0 {
				// Either wildcard address or the multicast address needs to be bound on non-windows systems
				found := false

				// First check if the multicast address is already bound
				channelResources := udp.inputSockets[utils.GetPhysicalPort(locator)]
				for _, resource := range channelResources {
					if resource.sinterface == locatorAddrStr {
						found = true
						break
					}
				}

				// Create a new resource if no one is found
				if !found {
					channelResource := udp.CreateInputChannelResource(locatorAddrStr, locator, true, maxMsgSize, receiver)
					resources := udp.inputSockets[utils.GetPhysicalPort(locator)]
					resources = append(resources, channelResource)
					udp.inputSockets[utils.GetPhysicalPort(locator)] = resources
					// TODO:
					// syscall.socket.Set
				}
			} else {
				// The multicast group will be joined silently, because we do not
				// want to return another resource.
				if channelResources, ok := udp.inputSockets[utils.GetPhysicalPort(locator)]; ok {
					addr := utils.ToIP4(locatorAddrStr)
					for _, channelResource := range channelResources {
						if channelResource.sinterface == kIPv4AddressAny {
							locNames := udp.getIPv4UniqueInterfaces(true)
							for _, infoIP := range locNames {
								ip := utils.ToIP4(infoIP.Name)
								var mreq = &syscall.IPMreq{Multiaddr: addr, Interface: ip}
								err := syscall.SetsockoptIPMreq(channelResource.udpConn, syscall.IPPROTO_IP, syscall.IP_ADD_MEMBERSHIP, mreq)
								if err != nil {
									log.Fatalf("SetsockoptIPMreq error: %v", err)
								}
							}
						} else {
							ip := utils.ToIP4(channelResource.sinterface)
							var mreq = &syscall.IPMreq{Multiaddr: addr, Interface: ip}
							err := syscall.SetsockoptIPMreq(channelResource.udpConn, syscall.IPPROTO_IP, syscall.IP_ADD_MEMBERSHIP, mreq)
							if err != nil {
								log.Fatalf("Error joining multicast group on (%v)", ip)
							}
						}
					}
				}
			}
		}
	}

	return success
}

func (udp *UDPv4Transport) getBindingInterfacesList() []string {
	var outputInterfaces []string
	if udp.isInterfaceWhitelistEmpty() {
		outputInterfaces = append(outputInterfaces, kIPv4AddressAny)
	} else {
		for ip := range udp.interfaceWhiteList {
			outputInterfaces = append(outputInterfaces, ip)
		}
	}
	return outputInterfaces
}

func (udp *UDPv4Transport) CreateInputChannelResource(sInterface string, locator *common.Locator, isMulticast bool,
	maxMsgSize uint32, receiver ITransportReceiver) *UDPChannelResource {
	unicastSocket := udp.OpenAndBindInputSocket(sInterface, utils.GetPhysicalPort(locator), isMulticast)
	channelResource := NewUDPChannelResource(udp, unicastSocket, maxMsgSize, locator, sInterface, receiver)

	return channelResource
}

func (udp *UDPv4Transport) OpenAndBindInputSocket(sIP string, port uint16, isMulticast bool) int {
	socketMC, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		log.Fatalf("create socket fd error: %v", err)
	}

	addr4 := utils.ToIP4(sIP)

	if isMulticast {
		// if err := syscall.SetsockoptInt(socketMC, syscall.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
		// 	log.Fatalf("reuse port error: %v\n", err)
		// }
		if err := syscall.SetsockoptInt(socketMC, syscall.SOL_SOCKET, unix.SO_REUSEADDR, 1); err != nil {
			log.Fatalln(err)
		}
	}

	err = syscall.Bind(socketMC, &syscall.SockaddrInet4{Port: int(port), Addr: addr4})
	if err != nil {
		log.Printf("Bind socket error: %v, ip(%v), port(%v), isMulticast(%v).", err, sIP, port, isMulticast)
	}

	return socketMC
}

func (udp *UDPv4Transport) OpenAndBindInputSockets(locator *common.Locator, receiver ITransportReceiver,
	isMulticast bool, maxMsgSize uint32) bool {
	udp.inputMapMutex.Lock()
	defer udp.inputMapMutex.Unlock()
	return udp.openAndBindInputSockets(locator, receiver, isMulticast, maxMsgSize)
}

func (udp *UDPv4Transport) openAndBindInputSockets(locator *common.Locator, receiver ITransportReceiver,
	isMulticast bool, maxMsgSize uint32) bool {
	interfaceList := udp.getBindingInterfacesList()
	for _, sInterface := range interfaceList {
		resource := udp.CreateInputChannelResource(sInterface, locator, isMulticast, maxMsgSize, receiver)
		channelResources := udp.inputSockets[utils.GetPhysicalPort(locator)]
		channelResources = append(channelResources, resource)
		udp.inputSockets[utils.GetPhysicalPort(locator)] = channelResources
	}

	return true
}

func (udp *UDPv4Transport) isInterfaceWhitelistEmpty() bool {
	return len(udp.interfaceWhiteList) == 0
}

func getIPv4s(returnLoopBack bool) []*utils.InfoIP {
	locNames, _ := utils.GetIPs(returnLoopBack)
	var results []*utils.InfoIP
	for _, loc := range locNames {
		if loc.Type == utils.KIP4 || loc.Type == utils.KIP4Local {
			results = append(results, loc)
		}
	}

	for _, loc := range results {
		loc.Locator.Kind = common.KLocatorKindUDPv4
	}

	return results
}

func (udp *UDPv4Transport) GetIPs(returnLoopBack bool) []*utils.InfoIP {
	return getIPv4s(returnLoopBack)
}

func (udp *UDPv4Transport) GetDefaultMetatrafficMulticastLocators(locators *common.LocatorList, multicastPort uint32) bool {
	loc := common.NewLocator()
	loc.Kind = common.KLocatorKindUDPv4
	loc.Port = multicastPort
	utils.SetIPv4WithBytes(loc, []common.Octet{239, 255, 0, 1})
	locators.PushBack(loc)
	return true
}

//IsLocatorAllowed Checks for whether locator is allowed.
func (udp *UDPv4Transport) IsLocatorAllowed(locator *common.Locator) bool {
	if udp.IsLocatorSupported(locator) == false {
		return false
	}
	if len(udp.interfaceWhiteList) == 0 || utils.IsMulticast(locator) {
		return true
	}

	var ip string
	for i := 0; i < len(locator.Address); i++ {
		ip = ip + string(locator.Address[i])
	}
	return udp.isInetrfaceAllowed(ip)
}

func (udp *UDPv4Transport) GetDefaultUnicastLocators(locators *common.LocatorList, unicastPort uint32) bool {
	locator := common.NewLocator()
	locator.Kind = common.KLocatorKindUDPv4
	locator.SetInvalidAddress()
	udp.FillUnicastLocator(locator, unicastPort)
	locators.PushBack(locator)

	return true
}

func (udp *UDPv4Transport) GetDefaultMetatrafficUnicastLocators(locators *common.LocatorList, unicastPort uint32) bool {
	loc := common.NewLocator()
	loc.Kind = common.KLocatorKindUDPv4
	loc.Port = unicastPort
	loc.SetInvalidAddress()
	locators.PushBack(loc)
	return true
}

func (udp *UDPv4Transport) isInetrfaceAllowed(ip string) bool {
	if len(udp.interfaceWhiteList) == 0 {
		return true
	}

	if ip[12] == 0 && ip[13] == 0 && ip[14] == 0 && ip[15] == 0 {
		return true
	}

	for allowed := range udp.interfaceWhiteList {
		if strings.Compare(allowed, ip) == 0 {
			return true
		}
	}

	return false
}

//NormalizeLocator ...
func (udp *UDPv4Transport) NormalizeLocator(locator *common.Locator) *common.LocatorList {
	var list common.LocatorList
	if utils.IsAny(locator) {
		locNames := getIPv4s(false)

		for _, infoIP := range locNames {
			if udp.isInetrfaceAllowed(infoIP.Name) {
				newloc := *locator
				utils.SetIPv4(&newloc, &infoIP.Locator)
				list.PushBack(&newloc)
			}
		}
		if list.Empty() {
			newloc := *locator
			utils.SetIPv4WithBytes(&newloc, []common.Octet{127, 0, 0, 1})
			list.PushBack(&newloc)
		}
	} else {
		list.PushBack(locator)
	}

	return &list
}
