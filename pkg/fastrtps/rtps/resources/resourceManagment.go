package resources

//MemoryManagementPolicy memory policy
type MemoryManagementPolicy uint8

//
const (
	KPreallocatedMemoryMode MemoryManagementPolicy = iota
	KPreallocatedWithReallocMemoryMode
	KDynamicReserveMemoryMode
	KDynamicReusableMemoryMode
)
