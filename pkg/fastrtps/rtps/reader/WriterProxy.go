package reader

import (
	"log"

	"dds/fastrtps/rtps/attributes"
	"dds/fastrtps/rtps/history"
	"dds/fastrtps/utils"
)

var _ history.IWriterProxyWithHistory = (*WriterProxy)(nil)

type WriterProxy struct {
}

func NewWriterProxy(areader *StatefulReader, locAlloc *attributes.RemoteLocatorsAllocationAttributes,
	changesAllocation *utils.ResourceLimitedContainerConfig) *WriterProxy {
	log.Fatalln("not impl")
	return nil
}
