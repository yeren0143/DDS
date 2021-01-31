package transport

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/utils"
	"os"
)

var _ ITransport = (*TCPv4Transport)(nil)
var _ ITCPTransport = (*TCPv4Transport)(nil)

type TCPv4Transport struct {
	tcpTransport
	configuration      *TCPv4TransportDescriptor
	interfaceWhiteList map[string]bool
}

func newTCPv4Transport(descriptor *TCPv4TransportDescriptor) *TCPv4Transport {
	return &TCPv4Transport{configuration: descriptor}
}

func (transport *TCPv4Transport) Init() bool {
	return true
}

func (transport *TCPv4Transport) Configuration() ITransportDescriptor {
	return transport.configuration
}

func (transport *TCPv4Transport) CompareLocatorIP(lh, rh *common.Locator) bool {
	return utils.CompareAddress(lh, rh, false)
}

func (transport *TCPv4Transport) CompareLocatorIPAndPort(lh, rh *common.Locator) bool {
	return utils.CompareAddressAndPhysicalPort(lh, rh)
}

func (transport *TCPv4Transport) FillLocalIP(loc *common.Locator) {
	utils.SetIPv4WithBytes(loc, []byte{127, 0, 0, 1})
	loc.Kind = common.KLocatorKindTCPv4
}

func (transport *TCPv4Transport) FillMetatrafficMulticastLocator(locator *common.Locator, wellKnownPort uint32) bool {
	return true
}

func (transport *TCPv4Transport) fillMetrafficUnicastLocator(locator *common.Locator, metatrafficUnicastPort uint32) bool {
	if utils.GetPhysicalPort(locator) == 0 {
		config := transport.configuration
		if config != nil {
			if len(config.ListeningPorts) > 0 {
				utils.SetPhysicalPort(locator, config.ListeningPorts[0])
			} else {
				utils.SetPhysicalPort(locator, uint16(os.Getpid()))
			}
		}
	}

	if utils.GetLogicalPort(locator) == 0 {
		utils.SetLogicalPort(locator, uint16(metatrafficUnicastPort))
	}

	return true
}

func (transport *TCPv4Transport) FillMetatrafficUnicastLocator(locator *common.Locator, metatrafficUnicastPort uint32) bool {
	result := transport.fillMetrafficUnicastLocator(locator, metatrafficUnicastPort)
	utils.SetWan(locator, transport.configuration.wanAddr[0], transport.configuration.wanAddr[1],
		transport.configuration.wanAddr[2], transport.configuration.wanAddr[3])
	return result
}

func (transport *TCPv4Transport) fillUnicastLocator(locator *common.Locator, wellKnownPort uint32) bool {
	if utils.GetPhysicalPort(locator) == 0 {
		config := transport.configuration
		if config != nil {
			if len(config.ListeningPorts) > 0 {
				utils.SetPhysicalPort(locator, config.ListeningPorts[0])
			} else {
				utils.SetPhysicalPort(locator, uint16(os.Getpid()))
			}
		}
	}

	if utils.GetLogicalPort(locator) == 0 {
		utils.SetLogicalPort(locator, uint16(wellKnownPort))
	}

	return true
}

func (transport *TCPv4Transport) FillUnicastLocator(locator *common.Locator, wellKnownPort uint32) bool {
	result := transport.fillUnicastLocator(locator, wellKnownPort)

	config := transport.configuration
	utils.SetWan(locator, config.wanAddr[0], config.wanAddr[1], config.wanAddr[2], config.wanAddr[3])
	return result
}

func (transport *TCPv4Transport) GetDefaultMetatrafficMulticastLocators(locators *common.LocatorList,
	multicastPort uint32) bool {
	return true
}

func (transport *TCPv4Transport) GetDefaultMetatrafficUnicastLocators(locators *common.LocatorList,
	unicastPort uint32) bool {
	locator := common.NewLocator()
	locator.Kind = transport.transportKind
	locator.SetInvalidAddress()
	transport.fillMetrafficUnicastLocator(locator, unicastPort)
	locators.PushBack(locator)
	return true
}

func (transport *TCPv4Transport) GetDefaultUnicastLocators(locators *common.LocatorList, unicastPort uint32) bool {
	locator := common.NewLocator()
	locator.Kind = transport.transportKind
	locator.SetInvalidAddress()
	transport.fillUnicastLocator(locator, unicastPort)
	locators.PushBack(locator)
	return true
}

func (transport *TCPv4Transport) SetReceiveBufferSize(size uint32) {
	transport.configuration.RcvBufferSize = size
}

func (transport *TCPv4Transport) SetSendBufferSize(size uint32) {
	transport.configuration.SendBufferSize = size
}

func (transport *TCPv4Transport) IsLocatorAllowed(locator *common.Locator) bool {
	if transport.IsLocatorSupported(locator) == false {
		return false
	}

	if len(transport.interfaceWhiteList) == 0 {
		return true
	}

	return transport.isInterfaceAllowed(utils.ToIPv4String(locator))
}

func (transport *TCPv4Transport) isInterfaceAllowed(inter string) bool {
	if len(transport.interfaceWhiteList) == 0 {
		return true
	}

	isAddressAny := false
	for i := 0; i < len(inter); i++ {
		if inter[i] != 0 {
			isAddressAny = true
			break
		}
	}
	if isAddressAny {
		return true
	}

	_, ok := transport.interfaceWhiteList[inter]
	return ok
}

func (transport *TCPv4Transport) MaxRecvBufferSize() uint32 {
	return transport.configuration.maxMessageSize
}

func (transport *TCPv4Transport) IsLocalLocator(locator *common.Locator) bool {
	if utils.HasWan(locator) {
		wan := utils.GetWan(locator)
		if wan != transport.configuration.wanAddr {
			return false
		}
	}

	if utils.IsLocal(locator) {
		return true
	}

	// TODO:
	//for _, current := range transport.current

	return false
}

func (transport *TCPv4Transport) getIPv4s(returnLoopBack bool) []*utils.InfoIP {
	ips, _ := utils.GetIPs(false)
	validIps := []*utils.InfoIP{}
	for _, ip := range ips {
		if ip.Type == utils.KIP4 || ip.Type == utils.KIP4Local {
			validIps = append(validIps, ip)
		}
	}

	for _, ip := range validIps {
		ip.Locator.Kind = common.KLocatorKindTCPv4
	}

	return validIps
}

func (transport *TCPv4Transport) NormalizeLocator(locator *common.Locator) *common.LocatorList {
	list := common.LocatorList{}
	if utils.IsAny(locator) {
		locNames := transport.getIPv4s(false)
		for _, infoIP := range locNames {
			if transport.isInterfaceAllowed(infoIP.Name) {
				newLoc := locator
				utils.SetIPv4(newLoc, &infoIP.Locator)
			}
		}

	}

	return &list
}

func (transport *TCPv4Transport) ConfigureInitialPeerLocator(locator *common.Locator, portParams *common.PortParameters,
	domainID uint32, list *common.LocatorList) bool {
	if utils.GetPhysicalPort(locator) == 0 {
		for i := uint32(0); i < transport.Configuration().MaxInitialPeersRange(); i++ {
			auxloc := *locator
			auxloc.Port = portParams.GetUnicastPort(domainID, i)

			if utils.GetLogicalPort(locator) == 0 {
				utils.SetLogicalPort(&auxloc, uint16(portParams.GetUnicastPort(domainID, i)))
			}

			list.PushBack(&auxloc)
		}
	} else {
		if utils.GetLogicalPort(locator) == 0 {
			for i := uint32(0); i < transport.Configuration().MaxInitialPeersRange(); i++ {
				auxloc := *locator
				utils.SetLogicalPort(&auxloc, uint16(portParams.GetUnicastPort(domainID, i)))
				list.PushBack(&auxloc)
			}

		} else {
			list.PushBack(locator)
		}
	}

	return true
}

func (transport *TCPv4Transport) fillLocalIP(local *common.Locator) {
	utils.SetIPv4WithBytes(local, []common.Octet{127, 0, 0, 1})
	local.Kind = common.KLocatorKindTCPv4
}

func (transport *TCPv4Transport) TransformRemoteLocator(remoteLocator *common.Locator) (*common.Locator, bool) {
	if transport.IsLocatorSupported(remoteLocator) == false {
		return nil, false
	}

	if transport.IsLocalLocator(remoteLocator) == false {
		return remoteLocator, true
	}

	if transport.IsLocatorAllowed(remoteLocator) == false {
		return nil, false
	}

	var resultLocator common.Locator
	transport.fillLocalIP(&resultLocator)
	if transport.IsLocatorAllowed(&resultLocator) {
		utils.SetPhysicalPort(&resultLocator, utils.GetPhysicalPort(remoteLocator))
		utils.SetLogicalPort(&resultLocator, utils.GetLogicalPort(remoteLocator))
		return &resultLocator, true
	}

	resultLocator = *remoteLocator
	return &resultLocator, true
}
