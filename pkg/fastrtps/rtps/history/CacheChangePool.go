package history

import "github.com/yeren0143/DDS/common"

var _ IChangePool = (*CacheChangePool)(nil)

type CacheChangePool struct {
}

func (changePool *CacheChangePool) ReserveCache() (*common.CacheChangeT, bool) {
	return nil, false
}

func (changePool *CacheChangePool) ReleaseCache(*common.CacheChangeT) bool {
	return false
}

func NewCacheChangePool(cfg *PoolConfig) *CacheChangePool {
	return &CacheChangePool{}
}
