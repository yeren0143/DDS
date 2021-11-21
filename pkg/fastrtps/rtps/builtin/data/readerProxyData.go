package data

import (
	"log"

	"dds/common"
	"dds/core/policy"
	"dds/fastrtps/rtps/attributes"
	"dds/fastrtps/rtps/network"
	"dds/fastrtps/rtps/qos"
)

/**
 * Class ReaderProxyData, used to represent all the information on a Reader
 * (both local and remote) with the purpose of implementing the discovery.
 */

type ReaderProxyData struct {
	GUID           common.GUIDT
	remoteLocators *common.RemoteLocatorList
	// GUID_t of the Reader converted to InstanceHandle_t
	Key common.InstanceHandleT
	// GUID_t of the participant converted to InstanceHandle
	RTPSParticipantKey common.InstanceHandleT
	typeName           string
	topicName          string
	userDefinedId      uint16
	// Field to indicate if the Reader is Alive.
	isAlive          bool
	topicKind        common.TopicKindT
	TypeID           policy.TypeIDV1
	TypeObj          policy.TypeObjectV1
	TypeInformation  *policy.TypeInformation
	Properties       *policy.ParameterPropertyListT
	ExpectsInlineQos bool
	Qos              *qos.ReaderQos
}

func (proxy *ReaderProxyData) GetSerializedSize(includeEncapsulation bool) uint32 {
	log.Panic("not impl")
	return 0
}

func (proxy *ReaderProxyData) SetAlive(isAlive bool) {
	proxy.isAlive = isAlive
}

func (proxy *ReaderProxyData) WriteToCDRMessage(msg *common.CDRMessage, writeEncapsulation bool) bool {
	log.Panic("not impl")
	return false
}

func (proxy *ReaderProxyData) SetLocators(locators *common.RemoteLocatorList) {
	copy(proxy.remoteLocators.Unicast, locators.Unicast)
	copy(proxy.remoteLocators.Multicast, locators.Multicast)
}

func (proxy *ReaderProxyData) RemoteLocators() *common.RemoteLocatorList {
	return proxy.remoteLocators
}

func (proxy *ReaderProxyData) SetRemoteLocators(locators *common.RemoteLocatorList,
	network *network.NetFactory, useMulticastLocators bool) {
	log.Fatalln("not impl")
}

func (proxy *ReaderProxyData) SetMulticastLocators(locators *common.LocatorList, networkFactory *network.NetFactory) {

}

func (proxy *ReaderProxyData) SetAnnouncedUnicastLocators(locators *common.LocatorList) {

}

func (proxy *ReaderProxyData) Clear() {

}

func NewReaderProxyData(maxUnicastLocators, maxMulticastLocators uint32,
	dataLimits *attributes.VariableLengthDataLimits) *ReaderProxyData {
	var proxyData = ReaderProxyData{
		ExpectsInlineQos: false,
		remoteLocators:   common.NewRemoteLocatorList(maxUnicastLocators, maxMulticastLocators),
		userDefinedId:    0,
		isAlive:          true,
		topicKind:        common.KNoKey,
	}

	// As DDS-XTypes, v1.2 (page 182) document stablishes, local default is ALLOW_TYPE_COERCION,
	// but when remotes doesn't send TypeConsistencyQos, we must assume DISALLOW.
	proxyData.Qos = qos.NewReaderQos()
	proxyData.Qos.TypeConsistency.Kind = policy.KAllowTypeCoercion

	return &proxyData
}
