package reader

import (
	"log"

	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/history"
	"github.com/yeren0143/DDS/fastrtps/utils"
)

var _ history.IWriterProxyWithHistory = (*WriterProxy)(nil)

type WriterProxy struct {
}

func NewWriterProxy(areader *StatefulReader, locAlloc *attributes.RemoteLocatorsAllocationAttributes,
	changesAllocation *utils.ResourceLimitedContainerConfig) *WriterProxy {
	log.Fatalln("not impl")
	return nil

}
