package utils

import (
	"bytes"
	"github.com/yeren0143/DDS/common"
	"log"
	"os"
	"strconv"
	"strings"
)

//SetIPv4WithBytes ...
func SetIPv4WithBytes(locator *common.Locator, ipv4 []common.Octet) {
	copy(locator.Address[12:], ipv4)
}

func setIPv4WithIP(locator *common.Locator, ip string) bool {
	s := strings.Split(ip, ".")
	ip0, _ := strconv.Atoi(s[0])
	ip1, _ := strconv.Atoi(s[1])
	ip2, _ := strconv.Atoi(s[2])
	ip3, _ := strconv.Atoi(s[3])
	locator.Address[12] = uint8(ip0)
	locator.Address[13] = uint8(ip1)
	locator.Address[14] = uint8(ip2)
	locator.Address[15] = uint8(ip3)

	return true
}

func setIP6WithString(locator *common.Locator, ip string) bool {
	hexdigits := strings.Split(ip, ":")

	//FOUND a . in the last element (MAP TO IP4 address)
	lastIndex := len(hexdigits) - 1
	last := hexdigits[lastIndex]
	if strings.Contains(last, ".") {
		return false
	}

	pos := strings.Index(hexdigits[lastIndex], "%")
	if pos > 0 {
		hexdigits[lastIndex] = hexdigits[lastIndex][:pos]
	}

	index := 15
	for i := lastIndex - 1; i >= 0; i-- {
		s := hexdigits[i]
		if len(s) <= 2 {
			locator.Address[index-1] = 0
			auxnumber, _ := strconv.ParseUint(s, 16, 0)
			locator.Address[index] = uint8(auxnumber)
		} else {
			part1, _ := strconv.ParseUint(s[len(s)-2:], 16, 0)
			locator.Address[index] = byte(part1)
			part2, _ := strconv.ParseUint(s[:len(s)-2], 16, 0)
			locator.Address[index-1] = byte(part2)
		}
		index -= 2
	}

	index = 0
	for i := 0; i < len(hexdigits); i++ {
		s := hexdigits[i]
		if len(s) <= 2 {
			locator.Address[index] = 0
			auxnumber, _ := strconv.ParseUint(s, 16, 0)
			locator.Address[index+1] = byte(auxnumber)
		} else {
			part1, _ := strconv.ParseUint(s[len(s)-2:], 16, 0)
			locator.Address[index+1] = byte(part1)
			part2, _ := strconv.ParseUint(s[:len(s)-2], 16, 0)
			locator.Address[index] = byte(part2)
		}
		index += 2
	}

	return true
}

//SetIPv4 Copies locator's IPv4.
func SetIPv4(dest *common.Locator, origin *common.Locator) {
	SetIPv4WithBytes(dest, GetIPv4(origin))
}

func ToPhysicalLocator(locator *common.Locator) *common.Locator {
	result := *locator
	SetLogicalPort(&result, 0)
	return &result
}

func ToIPv4String(locator *common.Locator) string {
	var ret string
	v0 := int(locator.Address[12])
	v1 := int(locator.Address[13])
	v2 := int(locator.Address[14])
	v3 := int(locator.Address[15])
	ret = ret + string(strconv.Itoa(v0))
	ret = ret + string(".")
	ret = ret + string(strconv.Itoa(v1))
	ret = ret + string(".")
	ret = ret + string(strconv.Itoa(v2))
	ret = ret + string(".")
	ret = ret + string(strconv.Itoa(v3))
	return ret
}

func ToIP4(sIP string) [4]byte {
	s := strings.Split(sIP, ".")
	if len(s) != 4 {
		log.Fatalf("invalid sIP when toIP4:%v", sIP)
	}
	v0, _ := strconv.Atoi(s[0])
	v1, _ := strconv.Atoi(s[1])
	v2, _ := strconv.Atoi(s[2])
	v3, _ := strconv.Atoi(s[3])

	return [4]byte{byte(v0), byte(v1), byte(v2), byte(v3)}
}

//GetIPv4 return ipv4 of locator
func GetIPv4(locator *common.Locator) []common.Octet {
	return locator.Address[12:]
}

//IsMulticast return true if locator is multicast
func IsMulticast(locator *common.Locator) bool {
	if locator.Kind == common.KLocatorKindTCPv4 || locator.Kind == common.KLocatorKindTCPv6 {
		return false
	}

	if locator.Kind == common.KLocatorKindUDPv4 {
		return locator.Address[12] >= 224 && locator.Address[12] <= 239
	} else {
		return locator.Address[0] == 0xFF
	}
}

//IsLocal return true if locator is local address
func IsLocal(locator *common.Locator) bool {
	if locator.Kind == common.KLocatorKindUDPv4 ||
		locator.Kind == common.KLocatorKindTCPv4 {
		local := []byte{127, 0, 0, 1}
		return bytes.Equal(locator.Address[12:], local)
	} else {
		local := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
		return bytes.Equal(locator.Address[0:], local)
	}
}

//IsAny Checks if a locator has any IP address.
func IsAny(locator *common.Locator) bool {
	if locator.Kind == common.KLocatorKindUDPv4 || locator.Kind == common.KLocatorKindTCPv4 {
		return locator.Address[12] == 0 &&
			locator.Address[13] == 0 &&
			locator.Address[14] == 0 &&
			locator.Address[15] == 0
	} else {
		return locator.Address[0] == 0 &&
			locator.Address[1] == 0 &&
			locator.Address[2] == 0 &&
			locator.Address[3] == 0 &&
			locator.Address[4] == 0 &&
			locator.Address[5] == 0 &&
			locator.Address[6] == 0 &&
			locator.Address[7] == 0 &&
			locator.Address[8] == 0 &&
			locator.Address[9] == 0 &&
			locator.Address[10] == 0 &&
			locator.Address[11] == 0 &&
			locator.Address[12] == 0 &&
			locator.Address[13] == 0 &&
			locator.Address[14] == 0 &&
			locator.Address[15] == 0
	}
}

// Checks if a both locators has the same IP address.
func CompareAddress(loc1, loc2 *common.Locator, fullAddress bool) bool {
	if loc1.Kind != loc2.Kind {
		return false
	}

	if !fullAddress && (loc1.Kind == common.KLocatorKindUDPv4 || loc1.Kind == common.KLocatorKindTCPv4) {
		return (loc1.Address[12] == loc2.Address[12]) && (loc1.Address[13] == loc2.Address[13]) &&
			(loc1.Address[14] == loc2.Address[14]) && (loc1.Address[15] == loc2.Address[15])
	} else {
		ret := true
		for i := 0; i < 16; i++ {
			if loc1.Address[i] != loc2.Address[i] {
				ret = false
				break
			}
		}
		return ret
	}
}

func SetPhysicalPort(locator *common.Locator, port uint16) bool {
	if os.Getenv("FASTDDS_IS_BIG_ENDIAN_TARGET") == "true" {
		locator.Port = uint32(port << 16)
	} else {
		locator.Port = uint32(port)
	}

	return locator.Port != 0
}

func GetPhysicalPort(locator *common.Locator) uint16 {
	locPhysical := locator.Port
	if os.Getenv("FASTDDS_IS_BIG_ENDIAN_TARGET") == "true" {
		return uint16(locPhysical & 0xFFFF0000)
	} else {
		return uint16(locPhysical & 0x0000FFFF)
	}
}

func GetLogicalPort(locator *common.Locator) uint16 {
	locLogical := locator.Port
	if os.Getenv("FASTDDS_IS_BIG_ENDIAN_TARGET") == "true" {
		return uint16(locLogical & 0xFF00)
	} else {
		return uint16(locLogical & 0x00FF)
	}
}

func SetLogicalPort(locator *common.Locator, port uint16) bool {
	if os.Getenv("FASTDDS_IS_BIG_ENDIAN_TARGET") == "true" {
		locator.Port = uint32(port << 16)
	} else {
		locator.Port = uint32(port)
	}

	return locator.Port != 0
}

func SetWan(locator *common.Locator, o1, o2, o3, o4 common.Octet) bool {
	locator.Address[8] = o1
	locator.Address[9] = o2
	locator.Address[10] = o3
	locator.Address[11] = o4
	return true
}

func HasWan(locator *common.Locator) bool {
	ret := locator.Kind == common.KLocatorKindTCPv4 // TCPv6 doesn't use WAN
	for i := 8; i < 12; i++ {
		if locator.Address[i] == 0 {
			return false
		}
	}
	return ret
}

func GetWan(locator *common.Locator) [4]common.Octet {
	var ret [4]common.Octet
	ret[0] = locator.Address[8]
	ret[1] = locator.Address[9]
	ret[2] = locator.Address[10]
	ret[3] = locator.Address[11]
	return ret
}

// Checks if a both locators has the same IP address and physical port  (as in RTCP protocol).
func CompareAddressAndPhysicalPort(loc1, loc2 *common.Locator) bool {
	return CompareAddress(loc1, loc2, true) && GetPhysicalPort(loc1) == GetPhysicalPort(loc2)
}
