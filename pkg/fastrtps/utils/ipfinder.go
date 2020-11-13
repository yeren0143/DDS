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
	"github.com/yeren0143/DDS/common"
	"errors"
	"fmt"
	"syscall"
	"unsafe"
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

func parseIP4(info *Info_IP) {
	info.locator.Kind = 1
	info.locator.Port = 0
	setIPv4WithIP(&info.locator, info.name)
	if IsLocal(&info.locator) {
		info.ip_type = IP4_LOCAL
	}
}

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

type addrinfoErrno int

func parseIP6(info *Info_IP) {
	info.locator.Kind = common.LOCATOR_KIND_UDPv6
	info.locator.Port = 0
	setIP6WithString(&info.locator, info.name)
	if IsLocal(&info.locator) {
		info.ip_type = IP6_LOCAL
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

func getIPs(return_loopback bool) ([]Info_IP, error) {
	var info_ips []Info_IP

	var ifap *C.struct_ifaddrs
	if C.getifaddrs(&ifap) == -1 {
		return nil, errors.New("getifaddes() failed")
	}
	defer C.freeifaddrs(ifap)

	for ifa := ifap; ifa != nil; ifa = ifa.ifa_next {
		fmt.Println(ifa)
		if ifa.ifa_addr == nil || (ifa.ifa_flags&C.IFF_RUNNING) == 0 {
			continue
		}

		family := ifa.ifa_addr.sa_family
		ifa_name := C.GoString(ifa.ifa_name)

		if family == C.AF_INET {

			var info Info_IP
			info.ip_type = IP4
			salen := C.socklen_t(syscall.SizeofSockaddrInet4)
			info.name = lookupAddr(ifa.ifa_addr, salen)
			info.dev = ifa_name
			parseIP4(&info)

			if return_loopback || info.ip_type != IP4_LOCAL {
				info_ips = append(info_ips, info)
			}

		} else if family == C.AF_INET6 {
			var info Info_IP
			info.ip_type = IP6
			salen := C.socklen_t(syscall.SizeofSockaddrInet6)
			info.name = lookupAddr(ifa.ifa_addr, salen)
			info.dev = ifa_name
			parseIP6(&info)

			if return_loopback || info.ip_type != IP6_LOCAL {
				info_ips = append(info_ips, info)
			}
		}
	}

	return info_ips, nil
}

func GetIP4Address() common.LocatorList {

	locators := common.NewLocatorList()

	ip_names, _ := getIPs(false)
	for _, ip := range ip_names {
		if ip.ip_type == IP4 {
			locators = append(locators, ip.locator)
		}
	}
	return locators
}
