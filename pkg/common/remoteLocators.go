package common

// Holds information about the locators of a remote entity.
type RemoteLocatorList struct {
	Unicast   []Locator
	Multicast []Locator
}

func (locators *RemoteLocatorList) AddUnicastLocator(locator *Locator) {
	locators.Unicast = append(locators.Unicast, *locator)
}

func (locators *RemoteLocatorList) AddMulticastLocator(locator *Locator) {
	locators.Multicast = append(locators.Multicast, *locator)
}

func NewRemoteLocatorList(maxUnicastLocators, maxMulticastLocators uint32) *RemoteLocatorList {
	var remoteLocatorList RemoteLocatorList
	return &remoteLocatorList
}
