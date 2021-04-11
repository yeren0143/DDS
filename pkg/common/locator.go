package common

// LocatorEnum ...
type LocatorEnum = int32

// transport enum
const (
	KLocatorPortInvalid  LocatorEnum = 0
	KLocatorKindReserved LocatorEnum = 0
	KLocatorKindUDPv4    LocatorEnum = 1
	KLocatorKindUDPv6    LocatorEnum = 2
	KLocatorKindTCPv4    LocatorEnum = 4
	KLocatorKindTCPv6    LocatorEnum = 8
	KLocatorKindShm      LocatorEnum = 16
)

//Locator ...
type Locator struct {
	Kind    int32
	Port    uint32
	Address [16]Octet
}

//NewLocator create locator with default value
func NewLocator() *Locator {
	return &Locator{
		Kind: KLocatorKindUDPv4,
	}
}

func (locator *Locator) Valid() bool {
	return locator.Kind >= 0
}

//SetInvalidAddress memset locator's address to 0
func (locator *Locator) SetInvalidAddress() {
	locator.Address = [16]Octet{0}
}

//LocatorList ...
type LocatorList struct {
	Locators []Locator
}

func (list *LocatorList) Set(aList *LocatorList) {
	if len(aList.Locators) == 0 {
		return
	}
	list.Locators = make([]Locator, len(aList.Locators))
	copy(list.Locators, aList.Locators)
}

//Valid ...
func (list *LocatorList) Valid() bool {
	for _, locator := range list.Locators {
		if locator.Valid() == false {
			return false
		}
	}
	return true
}

//Length ...
func (list *LocatorList) Length() int {
	return len(list.Locators)
}

func (list *LocatorList) PushBack(local *Locator) {
	list.Locators = append(list.Locators, *local)
}

//Append add locators to the end of locatorlist
func (list *LocatorList) Append(newList *LocatorList) *LocatorList {
	list.Locators = append(list.Locators, newList.Locators...)
	return list
}

func (list *LocatorList) Empty() bool {
	if len(list.Locators) > 0 {
		return false
	}

	return true
}

func (list *LocatorList) Clear() {
	list.Locators = nil
}

func NewLocatorList() *LocatorList {
	return &LocatorList{}
}
