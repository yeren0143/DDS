package common

// SequenceNumberT different for each change in the same writer.
type SequenceNumberT struct {
	Value uint64
}

var KSequenceNumberUnknown SequenceNumberT

func init() {
	KSequenceNumberUnknown.Value = ^uint64(0)
}

// Increment SequenceNumber.
func (seqNum *SequenceNumberT) Increment() {
	seqNum.Value++
}

func (seqNum *SequenceNumberT) Less(that *SequenceNumberT) bool {
	return seqNum.Value < that.Value
}

func (seqNum *SequenceNumberT) Equal(that *SequenceNumberT) bool {
	return seqNum.Value == that.Value
}

type SequenceNumberSet struct {
	base     SequenceNumberT
	rangeMax SequenceNumberT
	bitmap   [8]uint32
	numBits  uint32
}

// Base set a new base for the range.
// This method resets the range and sets a new value for its base.
func (sequences *SequenceNumberSet) Base(base SequenceNumberT) {
	sequences.base = base
	sequences.rangeMax.Value = base.Value + 255
	sequences.numBits = 0
	for i := 0; i < 8; i++ {
		sequences.bitmap[i] = 0
	}
}

// BitmapSet sets the current value of the bitmap.
// This method is designed to be used when performing deserialization of a bitmap range.
func (sequences *SequenceNumberSet) BitmapSet(numBits uint32, bitmap *[8]uint32) {
	if numBits < 256 {
		sequences.numBits = numBits
	} else {
		sequences.numBits = 256
	}

	sequences.bitmap = *bitmap
	sequences.bitmap[7] = sequences.bitmap[7] & ^(^uint32(0) >> (numBits & 31))
	sequences.calMaximumBitSet(8, 0)
}

func (sequences *SequenceNumberSet) clz(bits uint32) uint32 {
	firstNotZero := uint32(0)
	for i := uint32(0); i < 32; i++ {
		if (bits & (1 << i)) != 0 {
			firstNotZero = i
			break
		}
	}

	return firstNotZero
}

func (sequences *SequenceNumberSet) calMaximumBitSet(startingIndex, minIndex uint32) {
	sequences.numBits = 0
	for i := startingIndex; i > minIndex; {
		i--
		bits := sequences.bitmap[i]
		if bits != 0 {
			bits = bits & ^(bits - 1)
			offset := sequences.clz(bits) + 1
			sequences.numBits = (i << 5) + offset
			break
		}
	}
}

// NewSequenceNumberSet return SequenceNumberSet
func NewSequenceNumberSet(seqNum SequenceNumberT) *SequenceNumberSet {
	return &SequenceNumberSet{
		base:     seqNum,
		rangeMax: SequenceNumberT{Value: seqNum.Value + 255},
	}
}
