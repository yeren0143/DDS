package utils

// /*
// #include <sys/socket.h>
// #include <net/if.h>
// #include <arpa/inet.h>
// #include <arpa/inet.h>
// #include <ifaddrs.h>
// #ifndef AF_LINK
// #define AF_LINK AF_PACKET
// #endif
// #ifndef __linux__ // NOT LINUX
// u_int32_t Ibytes(void *data) { return ((struct if_data *)data)->ifi_ibytes; }
// u_int32_t Obytes(void *data) { return ((struct if_data *)data)->ifi_obytes; }
// u_int32_t Ipackets(void *data) { return ((struct if_data *)data)->ifi_ipackets; }
// u_int32_t Opackets(void *data) { return ((struct if_data *)data)->ifi_opackets; }
// u_int32_t Ierrors(void *data) { return ((struct if_data *)data)->ifi_ierrors; }
// u_int32_t Oerrors(void *data) { return ((struct if_data *)data)->ifi_oerrors; }
// #else
// #include <linux/if_link.h>
// u_int32_t Ibytes(void *data) { return ((struct rtnl_link_stats *)data)->rx_bytes; }
// u_int32_t Obytes(void *data) { return ((struct rtnl_link_stats *)data)->tx_bytes; }
// u_int32_t Ipackets(void *data) { return ((struct rtnl_link_stats *)data)->rx_packets; }
// u_int32_t Opackets(void *data) { return ((struct rtnl_link_stats *)data)->tx_packets; }
// u_int32_t Ierrors(void *data) { return ((struct rtnl_link_stats *)data)->rx_errors; }
// u_int32_t Oerrors(void *data) { return ((struct rtnl_link_stats *)data)->tx_errors; }
// #endif
// char ADDR[INET_ADDRSTRLEN];
// */

/*
#if defined(_WIN32)
#pragma comment(lib, "Iphlpapi.lib")
#include <stdio.h>
#include <winsock2.h>
#include <iphlpapi.h>
#include <ws2tcpip.h>
#include <assert.h>
#else
#include <arpa/inet.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <netdb.h>
#include <ifaddrs.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <net/if.h>
#endif

#if defined(__FreeBSD__)
#include <netinet/in.h>
#endif
*/
import "C"

import (
	"errors"
	"syscall"
	"unsafe"

	"dds/common"
)

// These are roughly enough for the following:
//
// Source		Encoding			Maximum length of single name entry
// Unicast DNS		ASCII or			<=253 + a NUL terminator
//			Unicode in RFC 5892		252 * total number of labels + delimiters + a NUL terminator
// Multicast DNS	UTF-8 in RFC 5198 or		<=253 + a NUL terminator
//			the same as unicast DNS ASCII	<=253 + a NUL terminator
// Local database	various				depends on implementation
const (
	nameinfoLen    = 64
	maxNameinfoLen = 4096
)

func (eai addrinfoErrno) Error() string   { return C.GoString(C.gai_strerror(C.int(eai))) }
func (eai addrinfoErrno) Temporary() bool { return eai == C.EAI_AGAIN }
func (eai addrinfoErrno) Timeout() bool   { return false }

func parseIP4(info *InfoIP) {
	info.Locator.Kind = 1
	info.Locator.Port = 0
	SetIPv4WithIP(&info.Locator, info.Name)
	if IsLocal(&info.Locator) {
		info.Type = KIP4Local
	}
}

type IPTYPE = int8

const (
	KIP4      IPTYPE = 0
	KIP6      IPTYPE = 1
	KIP4Local IPTYPE = 2
	KIP6Local IPTYPE = 3
)

type InfoIP struct {
	Type    IPTYPE
	Name    string
	Dev     string
	Locator common.Locator
}

type addrinfoErrno int

func parseIP6(info *InfoIP) {
	info.Locator.Kind = common.KLocatorKindUDPv6
	info.Locator.Port = 0
	setIP6WithString(&info.Locator, info.Name)
	if IsLocal(&info.Locator) {
		info.Type = KIP6Local
	}
}

func cgoNameinfoPTR(b []byte, sa *C.struct_sockaddr, salen C.socklen_t) (int, error) {
	gerrno, err := C.getnameinfo(sa, salen, (*C.char)(unsafe.Pointer(&b[0])), C.socklen_t(len(b)), nil, 0, C.NI_NUMERICHOST)
	return int(gerrno), err
}

func lookupAddr(sa *C.struct_sockaddr, salen C.socklen_t) string {
	var b []byte
	var gerrno int
	var err error
	for l := nameinfoLen; l <= maxNameinfoLen; l *= 2 {
		b = make([]byte, l)
		gerrno, err = cgoNameinfoPTR(b, sa, salen)

		if gerrno == 0 || gerrno != C.EAI_OVERFLOW {
			break
		}
	}
	if gerrno != 0 {
		switch gerrno {
		case C.EAI_SYSTEM:
			if err == nil { // see golang.org/issue/6232
				err = syscall.EMFILE
			}
		default:
			err = addrinfoErrno(gerrno)
		}
		panic(err.Error())
	}

	for i := 0; i < len(b); i++ {
		if b[i] == 0 {
			b = b[:i]
			break
		}
	}

	return string(b)
}

func GetIPs(return_loopback bool) ([]*InfoIP, error) {
	var InfoIPs []*InfoIP

	var ifap *C.struct_ifaddrs
	if C.getifaddrs(&ifap) == -1 {
		return nil, errors.New("getifaddes() failed")
	}
	defer C.freeifaddrs(ifap)

	for ifa := ifap; ifa != nil; ifa = ifa.ifa_next {
		if ifa.ifa_addr == nil || (ifa.ifa_flags&C.IFF_RUNNING) == 0 {
			continue
		}

		family := ifa.ifa_addr.sa_family
		ifa_name := C.GoString(ifa.ifa_name)

		if family == C.AF_INET {

			var info InfoIP
			info.Type = KIP4
			salen := C.socklen_t(syscall.SizeofSockaddrInet4)
			info.Name = lookupAddr(ifa.ifa_addr, salen)
			info.Dev = ifa_name
			parseIP4(&info)

			if return_loopback || info.Type != KIP4Local {
				InfoIPs = append(InfoIPs, &info)
			}

		} else if family == C.AF_INET6 {
			var info InfoIP
			info.Type = KIP6
			salen := C.socklen_t(syscall.SizeofSockaddrInet6)
			info.Name = lookupAddr(ifa.ifa_addr, salen)
			info.Dev = ifa_name
			parseIP6(&info)

			if return_loopback || info.Type != KIP6Local {
				InfoIPs = append(InfoIPs, &info)
			}
		}
	}

	return InfoIPs, nil
}

func GetIP4Address() *common.LocatorList {

	locators := common.NewLocatorList()

	ip_names, _ := GetIPs(false)
	for _, ip := range ip_names {
		if ip.Type == KIP4 {
			locators.Locators = append(locators.Locators, ip.Locator)
		}
	}
	return locators
}
