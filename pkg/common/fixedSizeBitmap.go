package common

// BitmapRange is a class to hold a range of items using a custom bitmap.
type BitmapRange struct {
	base       interface{} // Holds base value of the range.
	rangeMax   interface{} // Holds maximum allowed value of the range.
	bitmap     []uint32    // Holds the bitmap values.
	numBits    uint32      // Holds the highest bit set in the bitmap.
	maxNumBits uint32
}

// Base set a new base for the range.
// This method resets the range and sets a new value for its base.
func (bitmapRange *BitmapRange) Base(base interface{}) {
	// bitmapRange.base = base
	// bitmapRange.rangeMax = bitmapRange.base + bitmapRange.maxNumBits - 1
}
