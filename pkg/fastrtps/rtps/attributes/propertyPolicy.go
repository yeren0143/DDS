package attributes

import (
	. "common"
)

type PropertyPolicy struct {
	Properties       PropertySeq
	BinaryProperties BinaryPropertySeq
}

func NewPropertyPolicy() *PropertyPolicy {
	return &PropertyPolicy{}
}
