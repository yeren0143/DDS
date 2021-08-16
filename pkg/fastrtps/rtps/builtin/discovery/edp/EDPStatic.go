package edp

import (
	"log"

	"dds/common"
	"dds/fastrtps/rtps/attributes"
	"dds/fastrtps/rtps/builtin/data"
	"dds/fastrtps/rtps/reader"
	"dds/fastrtps/rtps/writer"
)

// Class EDPStaticProperty,
// used to read and write the strings from the properties used to transmit the EntityId_t.
type EDPStaticProperty struct {
	EndpointType string
	Status       string
	UserIdStr    string
	UserId       uint16
	EntityID     common.EntityIDT
}

var _ IEDP = (*EDPStatic)(nil)

// Class EDPStatic, implements a static endpoint discovery module.
type EDPStatic struct {
	edpBase
	Att *attributes.BuiltinAttributes
}

func (edp *EDPStatic) AssignRemoteEndpoints(pdata *data.ParticipantProxyData) {
	log.Fatalln("not Impl")
}

func (edp *EDPStatic) InitEDP(att *attributes.BuiltinAttributes) bool {
	log.Fatalln("not Impl")
	return false
}

func (edp *EDPStatic) ProcessLocalReaderProxyData(areader *reader.IRTPSReader, rdata *data.ReaderProxyData) bool {
	log.Fatalln("not Impl")
	return false
}

func (edp *EDPStatic) RemoveLocalReader(areader *reader.IRTPSReader) bool {
	log.Fatalln("not Impl")
	return false
}

func (edp *EDPStatic) ProcessLocalWriterProxyData(awriter *writer.IRTPSWriter, wdata *data.WriterProxyData) bool {
	log.Fatalln("not Impl")
	return false
}

func (edp *EDPStatic) RemoveLocalWriter(awriter *writer.IRTPSWriter) bool {
	log.Fatalln("not Impl")
	return false
}

func NewEDPStatic() *EDPStatic {
	return &EDPStatic{}
}
