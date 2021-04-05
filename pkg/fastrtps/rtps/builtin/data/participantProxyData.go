package data

import (
	"log"
	"os"
	"time"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/core/policy"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
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
	ProtoVersion             common.ProtocolVersionT
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

func (proxy *ParticipantProxyData) WriteToCDRMessage(msg *common.CDRMessage, writeEncapsulation bool) bool {
	if writeEncapsulation {
		if !policy.WriteEncapsulationToCDRMsg(msg) {
			return false
		}
	}

	{
		p := policy.NewParameterProtocolVersion(policy.KPidProtocolVersion, 4)
		p.ProtocolVersion = proxy.ProtoVersion
		if !policy.AddProtocolVersionToMsg(p, msg) {
			return false
		}
	}
	{
		p := policy.NewParameterVendorIDT(policy.KPidVendorID, 4)
		p.VendorID.Vendor[0] = proxy.VendorID.Vendor[0]
		p.VendorID.Vendor[1] = proxy.VendorID.Vendor[1]
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
	log.Panic("not impl")
}

func NewParticipantProxyData(att *attributes.RTPSParticipantAllocationAttributes) *ParticipantProxyData {
	var proxyData ParticipantProxyData
	return &proxyData
}
