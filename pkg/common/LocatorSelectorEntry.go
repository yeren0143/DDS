package common

/**
* Construct an EntryState object.
*
* @param max_unicast_locators    Maximum number of unicast locators to held by parent LocatorSelectorEntry.
* @param max_multicast_locators  Maximum number of multicast locators to held by parent LocatorSelectorEntry.
 */
type EntryState struct {
	Unicast   []uint32
	Multicast []uint32
}

func NewEntryState(maxUnicastLocators, maxMulticastLocators uint32) *EntryState {
	return &EntryState{
		Unicast:   make([]uint32, maxUnicastLocators),
		Multicast: make([]uint32, maxMulticastLocators),
	}
}

/**
* An entry for the @ref LocatorSelector.
*
* This class holds the locators of a remote endpoint along with data required for the locator selection algorithm.
* Can be easyly integrated inside other classes, such as @ref ReaderProxyData and @ref WriterProxyData.
 */
type LocatorSelectorEntry struct {
	RemoteGUID GUIDT
	Unicast    []Locator
	Multicast  []Locator
	State      EntryState
	// Indicates whether this entry should be taken into consideration.
	Enable bool
	// A temporary value for each transport to help optimizing some use cases.
	TransportShouldProcess bool
}

func NewLocatorSelectorEntry(maxUnicastLocators, maxMulticastLocators uint32) *LocatorSelectorEntry {
	return &LocatorSelectorEntry{
		RemoteGUID:             KGuidUnknown,
		Unicast:                make([]Locator, maxUnicastLocators),
		Multicast:              make([]Locator, maxMulticastLocators),
		State:                  *NewEntryState(maxUnicastLocators, maxMulticastLocators),
		Enable:                 false,
		TransportShouldProcess: false,
	}
}
