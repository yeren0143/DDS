package attributes

import (
	. "github.com/yeren0143/DDS/common"
)

type PropertyPolicy struct {
	Properties       PropertySeq
	BinaryProperties BinaryPropertySeq
}

func NewPropertyPolicy() *PropertyPolicy {
	return &PropertyPolicy{}
}
