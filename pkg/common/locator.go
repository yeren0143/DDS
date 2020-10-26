package common

type Locator_Enum = int8

const (
	LOCATOR_PORT_INVALID  Locator_Enum = 0
	LOCATOR_KIND_RESERVED Locator_Enum = 0
	LOCATOR_KIND_UDPv4    Locator_Enum = 1
	LOCATOR_KIND_UDPv6    Locator_Enum = 2
	LOCATOR_KIND_TCPv4    Locator_Enum = 4
	LOCATOR_KIND_TCPv6    Locator_Enum = 8
	LOCATOR_KIND_SHM      Locator_Enum = 16
)

type Locator struct {
	kind    int32
	port    uint32
	address [16]octet
}

type LocatorList = []Locator