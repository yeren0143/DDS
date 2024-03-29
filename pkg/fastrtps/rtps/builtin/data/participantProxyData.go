package data

import (
	"log"
	"os"
	"time"

	"dds/common"
	"dds/core/policy"
	"dds/fastrtps/rtps/attributes"
	"dds/fastrtps/rtps/network"
	"dds/fastrtps/rtps/resources"
)

var BUILTIN_PARTICIPANT_DATA_MAX_SIZE = uint32(100)
var TYPELOOKUP_DATA_MAX_SIZE = uint32(5000)
var DISC_BUILTIN_ENDPOINT_PARTICIPANT_ANNOUNCER = uint32(0x00000001 << 0)
var DISC_BUILTIN_ENDPOINT_PARTICIPANT_DETECTOR = uint32(0x00000001 << 1)
var DISC_BUILTIN_ENDPOINT_PUBLICATION_ANNOUNCER = uint32(0x00000001 << 2)
var DISC_BUILTIN_ENDPOINT_PUBLICATION_DETECTOR = uint32(0x00000001 << 3)
var DISC_BUILTIN_ENDPOINT_SUBSCRIPTION_ANNOUNCER = uint32(0x00000001 << 4)
var DISC_BUILTIN_ENDPOINT_SUBSCRIPTION_DETECTOR = uint32(0x00000001 << 5)
var DISC_BUILTIN_ENDPOINT_PARTICIPANT_PROXY_ANNOUNCER = uint32(0x00000001 << 6)
var DISC_BUILTIN_ENDPOINT_PARTICIPANT_PROXY_DETECTOR = uint32(0x00000001 << 7)
var DISC_BUILTIN_ENDPOINT_PARTICIPANT_STATE_ANNOUNCER = uint32(0x00000001 << 8)
var DISC_BUILTIN_ENDPOINT_PARTICIPANT_STATE_DETECTOR = uint32(0x00000001 << 9)
var BUILTIN_ENDPOINT_PARTICIPANT_MESSAGE_DATA_WRITER = uint32(0x00000001 << 10)
var BUILTIN_ENDPOINT_PARTICIPANT_MESSAGE_DATA_READER = uint32(0x00000001 << 11)
var BUILTIN_ENDPOINT_TYPELOOKUP_SERVICE_REQUEST_DATA_WRITER = uint32(0x00000001 << 12)
var BUILTIN_ENDPOINT_TYPELOOKUP_SERVICE_REQUEST_DATA_READER = uint32(0x00000001 << 13)
var BUILTIN_ENDPOINT_TYPELOOKUP_SERVICE_REPLY_DATA_WRITER = uint32(0x00000001 << 14)
var BUILTIN_ENDPOINT_TYPELOOKUP_SERVICE_REPLY_DATA_READER = uint32(0x00000001 << 15)
var DISC_BUILTIN_ENDPOINT_PUBLICATION_SECURE_ANNOUNCER = uint32(0x00000001 << 16)
var DISC_BUILTIN_ENDPOINT_PUBLICATION_SECURE_DETECTOR = uint32(0x00000001 << 17)
var DISC_BUILTIN_ENDPOINT_SUBSCRIPTION_SECURE_ANNOUNCER = uint32(0x00000001 << 18)
var DISC_BUILTIN_ENDPOINT_SUBSCRIPTION_SECURE_DETECTOR = uint32(0x00000001 << 19)
var BUILTIN_ENDPOINT_PARTICIPANT_MESSAGE_SECURE_DATA_WRITER = uint32(0x00000001 << 20)
var BUILTIN_ENDPOINT_PARTICIPANT_MESSAGE_SECURE_DATA_READER = uint32(0x00000001 << 21)
var DISC_BUILTIN_ENDPOINT_PARTICIPANT_SECURE_ANNOUNCER = uint32(0x00000001 << 26)
var DISC_BUILTIN_ENDPOINT_PARTICIPANT_SECURE_DETECTOR = uint32(0x00000001 << 27)

type ParticipantProxyData struct {
	ProtocolVersion          common.ProtocolVersionT
	Guid                     common.GUIDT
	VendorID                 common.VendorIDT
	ExpectsInlineQos         bool
	AviableBuiltinEndpoints  common.BuiltinEndpointSet
	MetatrafficLocators      common.RemoteLocatorList
	DefaultLocators          common.RemoteLocatorList
	ManualLivelinessCount    common.CountT
	ParticipantName          string
	Key                      common.InstanceHandleT
	Properties               policy.ParameterPropertyListT
	UserData                 *policy.UserDataQosPolicy
	ShouldCheckLeaseDuration bool
	LeaseDurationEvent       *resources.TimedEvent
	// Store the last timestamp it was received a RTPS message from the remote participant.
	LastReceivedMessageTm time.Time
	LeaseDuration         time.Duration
	leaseDurationMill     int64
	IsAlive               bool
	Readers               map[*ReaderProxyData]bool
	Writers               map[*WriterProxyData]bool
}

func (proxy *ParticipantProxyData) GetSerializedSize(includeEncapsulation bool) uint32 {
	retVal := uint32(0)
	if includeEncapsulation {
		retVal = 4
	}
	// PID_PROTOCOL_VERSION
	retVal += 4 + 4

	// PID_VENDORID
	retVal += 4 + 4

	if proxy.ExpectsInlineQos {
		// PID_EXPECTS_INLINE_QOS
		retVal += 4 + policy.PARAMETER_BOOL_LENGTH
	}

	// PID_PARTICIPANT_GUID
	retVal += 4 + uint32(policy.KParameterGuidLength)

	// PID_METATRAFFIC_MULTICAST_LOCATOR
	retVal += uint32(4+policy.KParameterLocatorLength) * uint32(len(proxy.MetatrafficLocators.Multicast))

	// PID_METATRAFFIC_UNICAST_LOCATOR
	retVal += uint32(4+policy.KParameterLocatorLength) * uint32(len(proxy.MetatrafficLocators.Unicast))

	// PID_DEFAULT_UNICAST_LOCATOR
	retVal += uint32(4+policy.KParameterLocatorLength) * uint32(len(proxy.DefaultLocators.Multicast))

	// PID_PARTICIPANT_LEASE_DURATION
	retVal += 4 + uint32(policy.KParameterTimeLength)

	// PID_BUILTIN_ENDPOINT_SET
	retVal += 4 + uint32(policy.KParameterBuiltinEndpointsetLength)

	if len(proxy.ParticipantName) > 0 {
		// PID_ENTITY_NAME
		retVal += uint32(len(proxy.ParticipantName))
	}

	if proxy.UserData.Size() > 0 {
		log.Fatalln("not Impl")
	}

	if proxy.Properties.Size() > 0 {
		log.Fatalln("not Impl")
	}

	return retVal + 4
}

func (proxy *ParticipantProxyData) Copy(pdata *ParticipantProxyData) {
	proxy.ProtocolVersion = pdata.ProtocolVersion
	proxy.Guid = pdata.Guid
	proxy.VendorID.Value[0] = pdata.VendorID.Value[0]
	proxy.VendorID.Value[1] = pdata.VendorID.Value[1]
	proxy.AviableBuiltinEndpoints = pdata.AviableBuiltinEndpoints
	proxy.MetatrafficLocators = pdata.MetatrafficLocators
	proxy.DefaultLocators = pdata.DefaultLocators
	proxy.ParticipantName = pdata.ParticipantName
	proxy.LeaseDuration = pdata.LeaseDuration
	proxy.leaseDurationMill = pdata.LeaseDuration.Milliseconds()
	proxy.Key = pdata.Key
	proxy.IsAlive = pdata.IsAlive
	proxy.UserData = pdata.UserData
	proxy.Properties = pdata.Properties
}

func (proxy *ParticipantProxyData) ReadFromCDRMessage(msg *common.CDRMessage, useEncapsulation bool,
	network *network.NetFactory, isShmTransportAvailable bool) bool {
	areShmMetrafficLocatorsPresent := false
	areShmDefaultLocatorsPresent := false
	isShmTransportPossible := false

	paramProcess := func(msg *common.CDRMessage, pid policy.ParameterIDT, plength uint16) bool {
		switch pid {
		case policy.KPidKeyHash:
			//p := policy.NewParameterKey(pid, plength)
			log.Fatalln("not Impl")
		case policy.KPidProtocolVersion:
			p := policy.NewParameterProtocolVersion(pid, plength)
			if !policy.ReadParameterProtocolFromCDRMessage(p, msg, plength) {
				return false
			}

			if p.ProtocolVersion.Major < common.KProtocolVersion.Major {
				return false
			}
			proxy.ProtocolVersion = p.ProtocolVersion
		case policy.KPidVendorID:
			p := policy.NewParameterVendorIDT(pid, plength)
			if !policy.ReadVendorIdFromCDRMessage(p, msg, plength) {
				return false
			}

			proxy.VendorID.Value[0] = p.VendorID.Value[0]
			proxy.VendorID.Value[1] = p.VendorID.Value[1]
			valid := (proxy.VendorID == common.KVendorIDTeProsima)
			isShmTransportAvailable = isShmTransportAvailable && valid
		case policy.KPidExpectsInlineQos:
			log.Fatalln("not impl")
		case policy.KPidParticipantGUID:
			p := policy.NewParameterGuid(pid, plength, &common.KGuidUnknown)
			if !policy.ReadGuidFromCDRMessage(p, msg, plength) {
				return false
			}
			proxy.Guid = p.Guid
			proxy.Key = common.CreateInstanceHandle(&p.Guid)
		case policy.KPidMetatrafficMulticastLocator:
			p := policy.NewParamaterLocator(pid, plength, common.NewLocator())
			if !policy.ReadLocatorFromCDRMessage(p, msg, plength) {
				return false
			}
			var tempLocator common.Locator
			if network.TransformRemoteLocator(&p.Locator, &tempLocator) {
				FilterLocators(isShmTransportAvailable, &isShmTransportPossible,
					&areShmMetrafficLocatorsPresent, &proxy.MetatrafficLocators, &tempLocator, false)
			}
		case policy.KPidMetatrafficUnicastLocator:
			p := policy.NewParamaterLocator(pid, plength, common.NewLocator())
			if !policy.ReadLocatorFromCDRMessage(p, msg, plength) {
				return false
			}
			var tempLocator common.Locator
			if network.TransformRemoteLocator(&p.Locator, &tempLocator) {
				FilterLocators(isShmTransportAvailable, &isShmTransportPossible,
					&areShmMetrafficLocatorsPresent, &proxy.MetatrafficLocators, &tempLocator, true)
			}
		case policy.KPidDefaultUnicastLocator:
			p := policy.NewParamaterLocator(pid, plength, common.NewLocator())
			if !policy.ReadLocatorFromCDRMessage(p, msg, plength) {
				return false
			}
			var tempLocator common.Locator
			if network.TransformRemoteLocator(&p.Locator, &tempLocator) {
				FilterLocators(isShmTransportAvailable, &isShmTransportPossible,
					&areShmDefaultLocatorsPresent, &proxy.DefaultLocators, &tempLocator, true)
			}
		case policy.KPidDefaultMulticastLocator:
			p := policy.NewParamaterLocator(pid, plength, common.NewLocator())
			if !policy.ReadLocatorFromCDRMessage(p, msg, plength) {
				return false
			}
			var tempLocator common.Locator
			if network.TransformRemoteLocator(&p.Locator, &tempLocator) {
				FilterLocators(isShmTransportAvailable, &isShmTransportPossible,
					&areShmDefaultLocatorsPresent, &proxy.DefaultLocators, &tempLocator, false)
			}
		case policy.KPidParticipantLeaseDuration:
			p := policy.NewParameterTimeT(pid, plength)
			if !policy.ReadTimeFromCDRMessage(p, msg, plength) {
				return false
			}
			proxy.LeaseDuration = time.Duration(p.Time.Seconds)*time.Second + time.Duration(p.Time.Nanosec)*time.Nanosecond
			proxy.leaseDurationMill = int64(p.Time.Seconds*1000000) + int64(p.Time.Nanosec/1000)

		case policy.KPidBuiltinEndpointSet:
			p := policy.NewParameterBuiltinEndpointSetT(pid, plength)
			if !policy.ReadBuiltinEndpointSetFromCDRMessage(p, msg, plength) {
				return false
			}
			proxy.AviableBuiltinEndpoints = p.EndpointSet
		case policy.KPidEntityName:
			p := policy.NewParameterString(pid, plength, "")
			if !policy.ReadEntityNameFromCDRMessage(p, msg, plength) {
				return false
			}
			proxy.ParticipantName = p.Name
		case policy.KPidPropertyList:
			log.Fatalln("not impl")
		case policy.KPidUserData:
			log.Fatalln("not impl")
		case policy.KPidIdentityToken:
			log.Fatalln("not impl")
		case policy.KPidPremissionsToken:
			log.Fatalln("not impl")
		case policy.KPidParticipantSecurityInfo:
			log.Fatalln("not impl")
		default:
			log.Println("policy kind not impl")
		}
		return true
	}

	proxy.Clear()

	var qosSize uint32
	return policy.ReadParameterListFromCDRMsg(msg, paramProcess, useEncapsulation, &qosSize)
}

func (proxy *ParticipantProxyData) WriteToCDRMessage(msg *common.CDRMessage, writeEncapsulation bool) bool {
	if writeEncapsulation {
		if !policy.WriteEncapsulationToCDRMsg(msg) {
			return false
		}
	}

	{
		p := policy.NewParameterProtocolVersion(policy.KPidProtocolVersion, 4)
		p.ProtocolVersion = proxy.ProtocolVersion
		if !policy.AddProtocolVersionToMsg(p, msg) {
			return false
		}
	}
	{
		p := policy.NewParameterVendorIDT(policy.KPidVendorID, 4)
		p.VendorID.Value[0] = proxy.VendorID.Value[0]
		p.VendorID.Value[1] = proxy.VendorID.Value[1]
		if !policy.AddVendorIDToMsg(p, msg) {
			return false
		}
	}

	if proxy.ExpectsInlineQos {
		log.Fatalln("not Impl")
	}

	{
		p := policy.NewParameterGuid(policy.KPidParticipantGUID, policy.KParameterGuidLength, &proxy.Guid)
		if !policy.AddGuidToMsg(p, msg) {
			return false
		}
	}

	for i := 0; i < len(proxy.MetatrafficLocators.Multicast); i++ {
		log.Fatalln("not Impl")
	}

	for i := 0; i < len(proxy.MetatrafficLocators.Unicast); i++ {
		p := policy.NewParamaterLocator(policy.KPidMetatrafficUnicastPort,
			policy.KParameterLocatorLength, &proxy.MetatrafficLocators.Unicast[i])
		if !policy.AddLocatorToMsg(p, msg) {
			return false
		}
	}

	for i := 0; i < len(proxy.DefaultLocators.Unicast); i++ {
		p := policy.NewParamaterLocator(policy.KPidDefaultUnicastLocator,
			policy.KParameterLocatorLength, &proxy.DefaultLocators.Unicast[i])
		if !policy.AddLocatorToMsg(p, msg) {
			return false
		}
	}

	for i := 0; i < len(proxy.DefaultLocators.Multicast); i++ {
		p := policy.NewParamaterLocator(policy.KPidDefaultMulticastLocator,
			policy.KParameterLocatorLength, &proxy.DefaultLocators.Multicast[i])
		if !policy.AddLocatorToMsg(p, msg) {
			return false
		}
	}

	{
		p := policy.NewParameterTimeT(policy.KPidParticipantLeaseDuration, policy.KParameterTimeLength)
		p.Time = *common.NewTime(&proxy.LeaseDuration)
		if !policy.AddTimeToMsg(p, msg) {
			return false
		}
	}

	{
		p := policy.NewParameterBuiltinEndpointSetT(policy.KPidBuiltinEndpointSet, policy.KParameterBuiltinEndpointsetLength)
		p.EndpointSet = proxy.AviableBuiltinEndpoints
		if !policy.AddBuiltinEndpointSetToMsg(p, msg) {
			return false
		}
	}

	if len(proxy.ParticipantName) > 0 {
		p := policy.NewParameterString(policy.KPidEntityName, 0, proxy.ParticipantName)
		if !policy.AddParameterStringToMsg(p, msg) {
			return false
		}
	}

	if proxy.UserData.Size() > 0 {
		log.Fatalln("not Impl")
	}

	if proxy.Properties.Size() > 0 {
		log.Fatalln("not Impl")
	}

	if os.Getenv("HAVE_SECURITY") != "" {
		log.Fatalln("not Impl")
	}

	return policy.AddParameterSentinelToMsg(msg)
}

func (proxy *ParticipantProxyData) SetPersistenceGuid(guid *common.GUIDT) {
	log.Fatalln("not impl")
}

func (proxy *ParticipantProxyData) GetPersistenceGuid() *common.GUIDT {
	log.Fatalln("not impl")
	return nil
}

func (proxy *ParticipantProxyData) SetPersistenceEntityID(id *common.EntityIDT) {
	log.Fatalln("not impl")
}

func (proxy *ParticipantProxyData) Clear() {
	proxy.ProtocolVersion = common.KDefaultProtocolVersion
	proxy.VendorID = common.KVendorIDTUnknown
	proxy.ExpectsInlineQos = false
	proxy.AviableBuiltinEndpoints = 0
	proxy.MetatrafficLocators.Unicast = []common.Locator{}
	proxy.MetatrafficLocators.Multicast = []common.Locator{}
	proxy.DefaultLocators.Unicast = []common.Locator{}
	proxy.DefaultLocators.Multicast = []common.Locator{}
	proxy.ParticipantName = ""
	proxy.Key = common.KInstanceHandleUnknown
	proxy.LeaseDuration = 0
	proxy.leaseDurationMill = 0
	proxy.IsAlive = true
	proxy.Properties.Clear()
	proxy.Properties.Length = 0
	proxy.UserData.Clear()
	proxy.UserData.Length = 0

}

func NewParticipantProxyData(allocation *attributes.RTPSParticipantAllocationAttributes) *ParticipantProxyData {
	var proxyData ParticipantProxyData
	proxyData.ProtocolVersion = common.KProtocolVersion
	proxyData.VendorID = common.KVendorIDTUnknown
	proxyData.ExpectsInlineQos = false
	proxyData.AviableBuiltinEndpoints = 0
	proxyData.MetatrafficLocators = *common.NewRemoteLocatorList(allocation.Locators.MaxUnicastLocators,
		allocation.Locators.MaxMulticastLocators)
	proxyData.DefaultLocators = *common.NewRemoteLocatorList(allocation.Locators.MaxUnicastLocators,
		allocation.Locators.MaxMulticastLocators)
	proxyData.IsAlive = false
	proxyData.Properties = *policy.NewParameterPropertyListT(uint32(allocation.DataLimits.MaxProperties))
	proxyData.ShouldCheckLeaseDuration = false
	proxyData.Readers = make(map[*ReaderProxyData]bool)
	proxyData.Writers = make(map[*WriterProxyData]bool)
	proxyData.UserData = &policy.UserDataQosPolicy{}

	// TODO:
	//m_userData.set_max_size(static_cast<uint32_t>(allocation.data_limits.max_user_data));

	return &proxyData
}
