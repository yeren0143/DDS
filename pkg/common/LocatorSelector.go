package common

/**
 * A class used for the efficient selection of locators when sending data to multiple entities.
 *
 * Algorithm:
 *   - Entries are added/removed with add_entry/remove_entry when matched/unmatched.
 *   - When data is to be sent:
 *     - A reference to this object is passed to the message group
 *     - For each submessage:
 *       - A call to reset is performed
 *       - A call to enable is performed per desired destination
 *       - If state_has_changed() returns true:
 *         - the message group is flushed
 *         - selection_start is called
 *         - for each transport:
 *           - transport_starts is called
 *           - transport handles the selection state of each entry
 *           - select may be called
 *       - Submessage is added to the message group
 */
type LocatorSelector struct {
	Entries    []*LocatorSelectorEntry
	Selections []uint32
	LastState  []int
}

func (selector *LocatorSelector) SelectedSize() uint32 {
	result := 0
	for _, index := range selector.Selections {
		entry := selector.Entries[index]
		result += len(entry.State.Multicast)
		result += len(entry.State.Unicast)
	}
	return uint32(result)
}
