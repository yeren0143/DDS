package utils

import (
	"math"
)

type ResourceLimitedContainerConfig struct {
	//! Number of elements to be preallocated in the collection.
	Initial uint32
	//! Maximum number of elements allowed in the collection.
	Maximum uint32
	//! Number of items to add when capacity limit is reached.
	Increment uint32
}

var KDefaultResourceLimitedContainerConfig = ResourceLimitedContainerConfig{
	Initial:   0,
	Maximum:   math.MaxUint32,
	Increment: 1,
}

func NewResourceLimitedContainerConfig() *ResourceLimitedContainerConfig {
	ret := KDefaultResourceLimitedContainerConfig
	return &ret
}
