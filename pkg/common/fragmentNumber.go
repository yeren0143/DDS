package common

type FragmentNumberT = uint32

const (
	defaultNItems = 8
)

//
type FragmentNumberSet struct {
	base     FragmentNumberT
	rangeMax FragmentNumberT
	bitmap   [8]uint32
	numBits  uint32
}

func (frags *FragmentNumberSet) Base(base FragmentNumberT) {
	frags.base = base
	frags.rangeMax = base + 255
	frags.numBits = 0
	for i := 0; i < 8; i++ {
		frags.bitmap[i] = 0
	}
}

func (frags *FragmentNumberSet) BitmapSet(numBits uint32, bitmap *[8]uint32) {
	if numBits < 256 {
		frags.numBits = numBits
	} else {
		frags.numBits = 256
	}

	frags.bitmap = *bitmap
	frags.bitmap[7] = frags.bitmap[7] & ^(^uint32(0) >> (numBits & 31))
	frags.calMaximumBitSet(8, 0)
}

func (frags *FragmentNumberSet) clz(bits uint32) uint32 {
	firstNotZero := uint32(0)
	for i := uint32(0); i < 32; i++ {
		if (bits & (1 << i)) != 0 {
			firstNotZero = i
			break
		}
	}

	return firstNotZero
}

func (frags *FragmentNumberSet) calMaximumBitSet(startingIndex, minIndex uint32) {
	frags.numBits = 0
	for i := startingIndex; i > minIndex; {
		i--
		bits := frags.bitmap[i]
		if bits != 0 {
			bits = bits & ^(bits - 1)
			offset := frags.clz(bits) + 1
			frags.numBits = (i << 5) + offset
			break
		}
	}
}

func NewFragmentNumberSet() *FragmentNumberSet {
	return &FragmentNumberSet{
		base:     256,
		rangeMax: 256 + 255,
	}
}
