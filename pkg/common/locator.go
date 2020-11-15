package common

// LocatorEnum ...
type LocatorEnum = int8

// transport enum
const (
	LocatorPortInvalid LocatorEnum = 0
	LocatorKindReserved LocatorEnum = 0
	LocatorKindUDPv4 LocatorEnum = 1
	LocatorKindUDPv6 LocatorEnum = 2
	LocatorKindTCPv4 LocatorEnum = 4
	LocatorKindTCPv6 LocatorEnum = 8
	LocatorKindShm LocatorEnum = 16
)

//Locator ...
type Locator struct {
	Kind    int8
	Port    uint32
	Address [16]Octet
}

type LocatorList = []Locator

func NewLocatorList() LocatorList {
	locator_list := []Locator{}
	return locator_list
}
