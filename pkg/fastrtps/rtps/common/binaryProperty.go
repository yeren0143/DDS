package common

type BinaryProperty struct {
	Name      string
	Value     []uint8
	Propagate bool
}

type BinaryPropertySeq = []BinaryProperty

func NewBinaryProperty() BinaryProperty {
	return BinaryProperty{}
}
