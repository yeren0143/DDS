package attributes

import (
	. "github.com/yeren0143/DDS/common"
)

type PropertyPolicyT struct {
	Properties       PropertySeq
	BinaryProperties BinaryPropertySeq
}

func NewPropertyPolicy() *PropertyPolicyT {
	return &PropertyPolicyT{}
}
