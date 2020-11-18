package common

// LocatorEnum ...
type LocatorEnum = int8

// transport enum
const (
	LocatorPortInvalid  LocatorEnum = 0
	LocatorKindReserved LocatorEnum = 0
	LocatorKindUDPv4    LocatorEnum = 1
	LocatorKindUDPv6    LocatorEnum = 2
	LocatorKindTCPv4    LocatorEnum = 4
	LocatorKindTCPv6    LocatorEnum = 8
	LocatorKindShm      LocatorEnum = 16
)

//Locator ...
type Locator struct {
	Kind    int8
	Port    uint32
	Address [16]Octet
}

func (locator *Locator) Valid() bool {
	return locator.Kind >= 0
}

//type LocatorList = []Locator

type LocatorList struct {
	Locators []Locator
}

func (list *LocatorList) Valid() bool {
	for _, locator := range list.Locators {
		if locator.Valid() == false {
			return false
		}
	}
	return true
}

func (list *LocatorList) PushBack(local *Locator) {
	list.Locators = append(list.Locators, *local)
}

func (list *LocatorList) Empty() bool {
	if len(list.Locators) > 0 {
		return false
	}

	return true
}

func NewLocatorList() *LocatorList {
	return &LocatorList{}
}
