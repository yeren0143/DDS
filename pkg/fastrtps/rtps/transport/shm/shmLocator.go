package shm

import (
	"log"

	"github.com/yeren0143/DDS/common"
)

type ShmType uint8

const (
	KUnicast ShmType = iota
	KMulticast
)

// Check whether a given locator is shared-memory kind and belongs to this host
func IsShmAndFromThisHost(locator *common.Locator) bool {
	if locator.Kind == common.KLocatorKindShm {
		log.Fatalln("not impl")
	}
	return false
}
