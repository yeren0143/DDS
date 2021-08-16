package data

import (
	"log"

	"dds/common"
	"dds/fastrtps/rtps/attributes"
	"dds/fastrtps/rtps/network"
	"dds/fastrtps/rtps/qos"
)

type WriterProxyData struct {
	Qos qos.WriterQos
}

func (proxy *WriterProxyData) GetSerializedSize(includeEncapsulation bool) uint32 {
	log.Panic("not impl")
	return 0
}

func (proxy *WriterProxyData) Clear() {
	log.Fatalln("not impl")
}

func (proxy *WriterProxyData) SetPersistenceGuid(id *common.GUIDT) {
	log.Fatalln("not impl")
}

func (proxy *WriterProxyData) SetPersistenceEntityID(nid *common.EntityIDT) {
	log.Fatalln("not impl")
}

func (proxy *WriterProxyData) Guid() *common.GUIDT {
	log.Fatalln("not impl")
	return nil
}

func (proxy *WriterProxyData) WriteToCDRMessage(msg *common.CDRMessage, writeEncapsulation bool) bool {
	log.Panic("not impl")
	return false
}

func (proxy *WriterProxyData) SetRemoteLocators(locators *common.RemoteLocatorList,
	network *network.NetFactory, useMulticastLocators bool) {

}

func NewWriterProxyData(maxUnicastLocators, maxMulticastLocators uint32,
	dataLimits *attributes.VariableLengthDataLimits) *WriterProxyData {
	return &WriterProxyData{}
}
