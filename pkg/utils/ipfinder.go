package utils

import (
	"common"
	"net"
)

type IPTYPE = int8

const (
	IP4       IPTYPE = 0
	IP6       IPTYPE = 1
	IP4_LOCAL IPTYPE = 2
	IP6_LOCAL IPTYPE = 3
)

type Info_IP struct {
	ip_type IPTYPE
	name    string
	dev     string
	locator common.Locator
}

func parseIP4(info *Info_IP) {
	info.locator.Kind = 1
	info.locator.Port = 0
	setIPv4WithIP(&info.locator, info.name)
	if IsLocal(&info.locator) {
		info.ip_type = IP4_LOCAL
	}
}

func parseIP6(info *Info_IP) {
	info.locator.Kind = common.LOCATOR_KIND_UDPv6
	info.locator.Port = 0
	setIP6WithString(&info.locator, info.name)
	if IsLocal(&info.locator) {
		info.ip_type = IP6_LOCAL
	}
}

func getIPsNoLoopBack() []Info_IP {
	var info_ips []Info_IP

	interfaces, _ := net.Interfaces()

	for _, inter := range interfaces {
		addrs, _ := inter.Addrs()
		if inter.Flags&net.FlagUp == 0 {
			continue
		}
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				var info Info_IP
				if ipnet.IP.To4() != nil {
					info.ip_type = IP4
					info.name = ipnet.IP.String()
					info.dev = inter.Name

					parseIP4(&info)
					if info.ip_type != IP4_LOCAL {
						info_ips = append(info_ips, info)
					}

				} else {
					info.ip_type = IP6
					info.name = ipnet.IP.String()
					info.dev = inter.Name

					parseIP6(&info)
					if info.ip_type != IP6_LOCAL {
						info_ips = append(info_ips, info)
					}
				}
			}
		}

	}

	return info_ips
}

func GetIP4Address() common.LocatorList {

	locators := common.NewLocatorList()

	ip_names := getIPsNoLoopBack()
	for _, ip := range ip_names {
		if ip.ip_type == IP4 {
			locators = append(locators, ip.locator)
		}
	}
	return locators
}
