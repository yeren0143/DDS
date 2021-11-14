package data

import (
	"log"

	"dds/common"
	"dds/core/policy"
	"dds/fastrtps/rtps/attributes"
	"dds/fastrtps/rtps/network"
	"dds/fastrtps/rtps/qos"
)

type WriterProxyData struct {
	Qos                qos.WriterQos
	guid               common.GUIDT
	remoteLocators     common.RemoteLocatorList
	key                common.InstanceHandleT
	rtpsParticipantKey common.InstanceHandleT
	typeName           string
	topicName          string
	userDefinedId      uint16
	typeMaxSerialized  uint32
	topicKind          common.TopicKindT
	persistenceGuid    common.GUIDT
	typeId             policy.TypeIDV1
	mType              policy.TypeObjectV1
	typeInformation    policy.TypeInformation
	mProperties        policy.ParameterPropertyListT
}

func (proxy *WriterProxyData) GetSerializedSize(includeEncapsulation bool) uint32 {
	log.Panic("not impl")
	return 0
}

func (proxy *WriterProxyData) Clear() {
	proxy.guid = common.KGuidUnknown
	proxy.remoteLocators.Unicast = []common.Locator{}
	proxy.remoteLocators.Multicast = []common.Locator{}
	proxy.key = common.InstanceHandleT{}
	proxy.rtpsParticipantKey = common.InstanceHandleT{}
	proxy.typeName = ""
	proxy.topicName = ""
	proxy.userDefinedId = 0
	proxy.Qos.Clear()
	proxy.typeMaxSerialized = 0
	proxy.topicKind = common.KNoKey
	proxy.persistenceGuid = common.KGuidUnknown
	proxy.mProperties.Clear()
	proxy.mProperties.Length = 0

	proxy.typeId = policy.TypeIDV1{}
	proxy.mType = policy.TypeObjectV1{}
	proxy.typeInformation = policy.TypeInformation{}
}

func (proxy *WriterProxyData) SetPersistenceGuid(guid *common.GUIDT) {
	log.Fatalln("not impl")
	if guid == &common.KGuidUnknown {
		return
	}
}

func (proxy *WriterProxyData) SetPersistenceEntityID(nid *common.EntityIDT) {
	log.Fatalln("not impl")
}

func (proxy *WriterProxyData) Guid() *common.GUIDT {
	return &proxy.guid
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
