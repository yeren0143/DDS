package common

// Holds information about the locators of a remote entity.
type RemoteLocatorList struct {
	Unicast   []Locator
	Multicast []Locator
}

func NewRemoteLocatorList(maxUnicastLocators, maxMulticastLocators uint32) *RemoteLocatorList {
	var remoteLocatorList RemoteLocatorList
	return &remoteLocatorList
}
