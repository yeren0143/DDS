package history

import (
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
)

type IBasicPoolConfig interface {
	GetMemoryPolicy() resources.MemoryManagementPolicy
	GetPayloadInitialSize() uint32
}

type basicPoolConfig struct {
	MemoryPolicy       resources.MemoryManagementPolicy
	PayloadInitialSize uint32
}

var _ IBasicPoolConfig = (*PoolConfig)(nil)
type PoolConfig struct {
	basicPoolConfig

	// Initial number of elements when preallocating data.
	InitialSize uint32

	// Maximum number of elements in the pool. Default value is 0,
	// indicating to make allocations until they fail.
	MaximumSize uint32
}

func (config *PoolConfig) GetMemoryPolicy() resources.MemoryManagementPolicy {
	return config.MemoryPolicy
}

func (config *PoolConfig) GetPayloadInitialSize() uint32 {
	return config.PayloadInitialSize
}

func FromHistoryAttributes(history *attributes.HistoryAttributes) PoolConfig {
	var config PoolConfig
	config.basicPoolConfig = basicPoolConfig{
		MemoryPolicy:       history.MemoryPolicy,
		PayloadInitialSize: history.PayloadMaxSize,
	}

	// Negative or 0 means no preallocation.
	// Otherwise, to avoid dynamic allocations after the pools are created, we need to
	// reserve an extra slot, as a call to `reserve` needs to succeed even if the history is full.
	if history.InitialReservedCaches <= 0 {
		config.InitialSize = 0
	} else {
		config.InitialSize = history.InitialReservedCaches + 1
	}

	// Negative or 0 means infinite maximum.
	// Otherwise, we need to allow one extra slot, as a call to `reserve` needs to succeed
	// even if the history is full, as old changes will only be removed when a change is added
	// to the history.
	if history.MaximumReservedCaches <= 0 {
		config.MaximumSize = 0
	} else {
		config.MaximumSize = history.MaximumReservedCaches + 1
	}

	return config
}
