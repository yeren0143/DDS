package history

import (
	"github.com/yeren0143/DDS/common"
)

// An interface for classes responsible of cache changes allocation management.
type IChangePool interface {
	// Get a new cache change from the pool
	ReserveCache() (bool, *common.CacheChangeT)

	// Return a cache change to the pool
	ReleaseCache(*common.CacheChangeT) bool
}
