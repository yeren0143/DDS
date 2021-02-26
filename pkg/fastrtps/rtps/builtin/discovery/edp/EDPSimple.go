package edp

import (
	"log"
	"sync"

	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/data"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/protocol"
	"github.com/yeren0143/DDS/fastrtps/rtps/history"
	"github.com/yeren0143/DDS/fastrtps/rtps/reader"
	"github.com/yeren0143/DDS/fastrtps/rtps/writer"
)

var _ IEDP = (*EDPSimple)(nil)
var _ edpImpl = (*EDPSimple)(nil)

// Class EDPSimple, implements the Simple Endpoint Discovery Protocol defined in the RTPS specification.
// Inherits from EDP class.
type EDPSimple struct {
	edpBase
	Discovery             attributes.BuiltinAttributes
	PublicationsListener  IEDPListener
	SubscriptionsListener IEDPListener
	pubWriterPayloadPool  history.ITopicPayloadPool
	pubReaderPayloadPool  history.ITopicPayloadPool
	subWriterPayloadPool  history.ITopicPayloadPool
	subReaderPayloadPool  history.ITopicPayloadPool
	tempDataMutex         sync.Mutex
	tempReaderProxyData   data.ReaderProxyData
	tempWriterProxyData   data.WriterProxyData
}

func (edp *EDPSimple) AssignRemoteEndpoints(pdata *data.ParticipantProxyData) {
	log.Fatalln("not Impl")
}

//Initialization of history attributes for EDP built-in readers
func (edp *EDPSimple) setBuiltinReaderHistoryAttributes(att *attributes.HistoryAttributes) {
	log.Fatalln("not Impl")
}

//Initialization of history attributes for EDP built-in writers
func (edp *EDPSimple) setBuiltinWriterHistoryAttributes(att *attributes.HistoryAttributes) {
	log.Fatalln("not Impl")
}

//Initialization of reader attributes for EDP built-in readers
func (edp *EDPSimple) setBuiltinReaderAttributes(att *attributes.ReaderAttributes) {
	log.Fatalln("not Impl")
}

//Initialization of writer attributes for EDP built-in writers
func (edp *EDPSimple) setBuiltinWriterAttributes(att *attributes.WriterAttributes) {
	log.Fatalln("not Impl")
}

func (edp *EDPSimple) createSEDPEndpoints() bool {
	var watt attributes.WriterAttributes
	var ratt attributes.ReaderAttributes
	var readerHistoryAtt attributes.HistoryAttributes
	var writerHistoryAtt attributes.HistoryAttributes

	edp.setBuiltinReaderHistoryAttributes(&readerHistoryAtt)
	edp.setBuiltinWriterHistoryAttributes(&writerHistoryAtt)
	edp.setBuiltinReaderAttributes(&ratt)
	edp.setBuiltinWriterAttributes(&watt)

	edp.PublicationsListener = newEDPSimplePubListener(edp)
	edp.SubscriptionsListener = newEDPSimpleSubListener(edp)

	if edp.Discovery.DiscoveryConfig.SimpleEDP.UsePublicationWriterAndSubscriptionReader {
		log.Fatalln("notImpl")
	}

	if edp.Discovery.DiscoveryConfig.SimpleEDP.UsePublicationReaderAndSubscriptionWriter {
		log.Fatalln("notImpl")
	}

	log.Println("Creation finished")
	return true
}

func (edp *EDPSimple) InitEDP(att *attributes.BuiltinAttributes) bool {
	log.Println("Beginning Simple Endpoint Discovery Protocol")
	edp.Discovery = *att

	if !edp.createSEDPEndpoints() {
		log.Fatalln("Problem creation SimpleEDP endpoints")
		return false
	}
	return true
}

func (edp *EDPSimple) ProcessLocalWriterProxyData(awriter *writer.IRTPSWriter, wdata *data.WriterProxyData) bool {
	log.Fatalln("not Impl")
	return false
}

func (edp *EDPSimple) ProcessLocalReaderProxyData(areader *reader.IRTPSReader, rdata *data.ReaderProxyData) bool {
	log.Fatalln("not Impl")
	return false
}

func (edp *EDPSimple) RemoveLocalWriter(awriter *writer.IRTPSWriter) bool {
	log.Fatalln("not Impl")
	return false
}

func (edp *EDPSimple) RemoveLocalReader(areader *reader.IRTPSReader) bool {
	log.Fatalln("not Impl")
	return false
}

func NewEDPSimple(p protocol.IPDP, part protocol.IParticipant) *EDPSimple {
	var edpSimple EDPSimple
	edpSimple.edpBase = *NewEDPBase(p, part)
	edpSimple.edpImpl = &edpSimple

	return &edpSimple
}
