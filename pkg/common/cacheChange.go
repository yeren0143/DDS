package common

import (
	"unsafe"
)

// ChangeKindT different types of CacheChange_t
type ChangeKindT = uint8

const (
	KAlive ChangeKindT = iota
	KNotAliveDisposed
	KNotAliveUnregistered
	KNotAliveDisposedUnregistered
)

// ICacheChangeOwner is a pool that created the payload of cache change
type ICacheChangeOwner interface {
	ReleasePayload(*CacheChangeT) bool
}

// CacheChangeT contains information on a specific CacheChange.
type CacheChangeT struct {
	Kind                 ChangeKindT
	WriterGUID           GUIDT
	InstanceHandle       InstanceHandleT
	SequenceNumber       SequenceNumberT
	SerializedPayload    SerializedPayloadT
	IsRead               bool
	SourceTimestamp      Time // Source TimeStamp (only used in Readers)
	ReceptionTimestamp   Time // Reception TimeStamp (only used in Readers)
	WriteParams          WriteParamsT
	IsUntyped            bool
	fragmentSize         uint16
	fragmentCount        uint32
	firstMissingFragment uint32
	payloadOwner         ICacheChangeOwner
}

func (cache *CacheChangeT) PayloadOwner() ICacheChangeOwner {
	return cache.payloadOwner
}

func (cache *CacheChangeT) SetPayloadOwner(owner ICacheChangeOwner) {
	cache.payloadOwner = owner
}

func (cache *CacheChangeT) nextFragmentPointer(fragmentIndex uint32) *uint32 {
	offset := uint32(cache.fragmentSize) * fragmentIndex
	offset = (offset + 3) & (^uint32(3))
	add := &(cache.SerializedPayload.Data[offset])
	return (*uint32)(unsafe.Pointer(add))
}

func (cache *CacheChangeT) setNextMissingFragment(fragmentIndex, nextFragmentIndex uint32) {
	ptr := cache.nextFragmentPointer(fragmentIndex)
	*ptr = nextFragmentIndex
}

// SetFragmentSize Set fragment size for this change.
func (cache *CacheChangeT) SetFragmentSize(fragmentSize uint16, createFragmentList bool) {
	cache.fragmentSize = fragmentSize
	cache.fragmentCount = 0
	cache.firstMissingFragment = 0

	if fragmentSize > 0 {
		// This follows RTPS 8.3.7.3.5
		cache.fragmentCount = (cache.SerializedPayload.Length + uint32(fragmentSize) - 1) / uint32(fragmentSize)

		if createFragmentList {
			// Keep index of next fragment on the payload portion at the beginning of each fragment. Last
			// fragment will have fragment_count_ as 'next fragment index'
			offset := uint32(0)
			for i := uint32(1); i < cache.fragmentCount; i++ {
				cache.setNextMissingFragment(i-1, i) // index to next fragment in missing list
				offset += uint32(cache.fragmentSize)
			}
		} else {
			// List not created. This means we are going to send this change fragmented, so it is already
			// assembled, and the missing list is empty (i.e. first missing points to fragment count)
			cache.firstMissingFragment = cache.fragmentCount
		}
	}
}

func NewCacheChangeT() *CacheChangeT {
	return &CacheChangeT{
		Kind:                 KAlive,
		IsRead:               false,
		IsUntyped:            true,
		fragmentSize:         0,
		fragmentCount:        0,
		firstMissingFragment: 0,
	}
}
