package utils

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/yeren0143/DDS/common"
)

func setIPv4WithBytes(locator *common.Locator, ipv4 []common.Octet) {
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
	last_index := len(hexdigits) - 1
	last := hexdigits[last_index]
	if strings.Contains(last, ".") {
		return false
	}

	pos := strings.Index(hexdigits[last_index], "%")
	if pos > 0 {
		hexdigits[last_index] = hexdigits[last_index][:pos]
	}

	index := 15
	for i := last_index - 1; i >= 0; i -= 1 {
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
	for i := 0; i < len(hexdigits); i += 1 {
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

func setIPv4(dest *common.Locator, origin *common.Locator) {
	setIPv4WithBytes(dest, GetIPv4(origin))
}

func GetIPv4(locator *common.Locator) []common.Octet {
	return locator.Address[12:]
}

func IsLocal(locator *common.Locator) bool {
	if locator.Kind == common.LOCATOR_KIND_UDPv4 ||
		locator.Kind == common.LOCATOR_KIND_TCPv4 {
		local := []byte{127, 0, 0, 1}
		return bytes.Equal(locator.Address[12:], local)
	} else {
		local := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
		return bytes.Equal(locator.Address[0:], local)
	}
}
