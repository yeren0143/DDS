package data

import (
	"dds/common"
	"dds/fastrtps/rtps/transport/shm"
)

func FilterLocators(isShmTransportAviable bool, isShmTransportPossible *bool,
	areShmLocatorsPresent *bool, targetLocatorsList *common.RemoteLocatorList,
	tempLocator *common.Locator, isUnicast bool) {
	if isShmTransportAviable && !(*isShmTransportPossible) {
		*isShmTransportPossible = shm.IsShmAndFromThisHost(tempLocator)
	}

	if *isShmTransportPossible {
		if tempLocator.Kind == common.KLocatorKindShm {
			// First SHM locator
			if !(*areShmLocatorsPresent) {
				// Remove previously added locators from other transports
				targetLocatorsList.Unicast = []common.Locator{}
				targetLocatorsList.Multicast = []common.Locator{}
				*areShmLocatorsPresent = true
			}

			if isUnicast {
				targetLocatorsList.AddUnicastLocator(tempLocator)
			} else {
				targetLocatorsList.AddMulticastLocator(tempLocator)
			}
		} else if !(*areShmLocatorsPresent) {
			if isUnicast {
				targetLocatorsList.AddUnicastLocator(tempLocator)
			} else {
				targetLocatorsList.AddMulticastLocator(tempLocator)
			}
		}
	} else {
		if tempLocator.Kind != common.KLocatorKindShm {
			if isUnicast {
				targetLocatorsList.AddUnicastLocator(tempLocator)
			} else {
				targetLocatorsList.AddMulticastLocator(tempLocator)
			}
		}
	}
}
