package attributes

import (
	"dds/common"
)

type PropertyPolicyT struct {
	Properties       common.PropertySeq
	BinaryProperties common.BinaryPropertySeq
}

func FindProperty(policy *PropertyPolicyT, name string) string {
	var returnValue string
	for _, property := range policy.Properties {
		if property.Name == name {
			returnValue = property.Value
			break
		}
	}

	return returnValue
}

func NewPropertyPolicy() *PropertyPolicyT {
	return &PropertyPolicyT{}
}
