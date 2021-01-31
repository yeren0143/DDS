package data

import (
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
)

type WriterProxyData struct {
}

func NewWriterProxyData(maxUnicastLocators, maxMulticastLocators uint32,
	dataLimits *attributes.VariableLengthDataLimits) *WriterProxyData {
	return &WriterProxyData{}
}
