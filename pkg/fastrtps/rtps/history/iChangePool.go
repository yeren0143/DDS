package history

import (
	"dds/common"
)

// An interface for classes responsible of cache changes allocation management.
type IChangePool interface {
	// Get a new cache change from the pool
	ReserveCache() (*common.CacheChangeT, bool)

	// Return a cache change to the pool
	ReleaseCache(*common.CacheChangeT) bool
}
