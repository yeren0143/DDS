package writer

import (
	"log"

	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
)

// ReaderProxy class that helps to keep the state of a specific Reader with respect to the RTPSWriter.
type ReaderProxy struct {
}

func NewReaderProxy(times *attributes.WriterTimes, locAlloc *attributes.RemoteLocatorsAllocationAttributes,
	awriter *StatefulWriter) *ReaderProxy {
	log.Fatalln("not impl")
	return &ReaderProxy{}
}
