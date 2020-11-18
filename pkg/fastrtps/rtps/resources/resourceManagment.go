package resources

type MemoryManagementPolicy uint8

// const (
// 	PREALLOCATED_MEMORY_MODE MemoryManagementPolicy = iota
// 	PREALLOCATED_WITH_REALLOC_MEMORY_MODE
// 	DYNAMIC_RESERVE_MEMORY_MODE
// 	DYNAMIC_REUSABLE_MEMORY_MODE
// )

const (
	CPreallocatedMemoryMode MemoryManagementPolicy = iota
	CPreallocatedWithReallocMemoryMode
	CDynamicReserveMemoryMode
	CDynamicReusableMemoryMode
)
