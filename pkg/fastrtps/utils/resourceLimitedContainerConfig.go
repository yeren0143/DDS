package utils

type ResourceLimitedContainerConfig struct {
	//! Number of elements to be preallocated in the collection.
	Initial uint
	//! Maximum number of elements allowed in the collection.
	Maximum uint
	//! Number of items to add when capacity limit is reached.
	Increment uint
}

func NewResourceLimitedContainerConfig() *ResourceLimitedContainerConfig {
	return &ResourceLimitedContainerConfig{
		Initial:   0,
		Maximum:   ^uint(0),
		Increment: 1,
	}
}
