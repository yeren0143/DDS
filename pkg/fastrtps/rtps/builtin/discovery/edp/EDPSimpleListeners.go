package edp

import (
	"log"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/reader"
	"github.com/yeren0143/DDS/fastrtps/rtps/writer"
)

var _ reader.IReaderListener = (IEDPListener)(nil)
var _ writer.IWriterListener = (IEDPListener)(nil)

type IEDPListener interface {
	reader.IReaderListener
	writer.IWriterListener
}

type EDPListenerBase struct {
	reader.ReaderListenerBase
	writer.WriterListenerBase
}

func (edpListener *EDPListenerBase) ComputeKey(aChange *common.CacheChangeT) bool {
	log.Fatalln("notImpl")
	return false
}

// returns true if loading info from persistency database
func ongoingDeserialization(edp IEDP) bool {
	log.Fatalln("notImpl")
	return false
}

// Class EDPSimplePUBReaderListener, used to define the behavior when a new WriterProxyData is received.
type EDPSimplePubListener struct {
	EDPListenerBase
}

func newEDPSimplePubListener(edp *EDPSimple) *EDPSimplePubListener {
	return &EDPSimplePubListener{}
}

//Class EDPSimpleSUBReaderListener, used to define the behavior when a new ReaderProxyData is received.
type EDPSimpleSubListener struct {
	EDPListenerBase
}

func newEDPSimpleSubListener(edp *EDPSimple) *EDPSimpleSubListener {
	return &EDPSimpleSubListener{}
}
