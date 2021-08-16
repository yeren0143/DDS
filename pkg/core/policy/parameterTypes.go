package policy

import (
	"dds/common"
)

type ParameterIDT = uint16

const (
	KParameterKeyHashLength = 16
)

const (
	KPidPad                              ParameterIDT = 0x0000
	KPidSentinel                         ParameterIDT = 0x0001
	KPidUserData                         ParameterIDT = 0x002c
	KPidTopicName                        ParameterIDT = 0x0005
	KPidTypeName                         ParameterIDT = 0x0007
	KPidGroupData                        ParameterIDT = 0x002d
	KPidTopicData                        ParameterIDT = 0x002e
	KPidDurability                       ParameterIDT = 0x001d
	KPidDurabilityService                ParameterIDT = 0x001e
	KPidDeadline                         ParameterIDT = 0x0023
	KPidLatencyBudget                    ParameterIDT = 0x0027
	KPidLiveliness                       ParameterIDT = 0x001b
	KPidReliability                      ParameterIDT = 0x001A
	KPidLifeSpan                         ParameterIDT = 0x002b
	KPidDestinationOrder                 ParameterIDT = 0x0025
	KPidHistory                          ParameterIDT = 0x0040
	KPidResourceLimits                   ParameterIDT = 0x0041
	KPidOwnership                        ParameterIDT = 0x001f
	KPidOwnershipStrength                ParameterIDT = 0x0006
	KPidPresentation                     ParameterIDT = 0x0021
	KPidPartition                        ParameterIDT = 0x0029
	KPidTimeBasedFilter                  ParameterIDT = 0x0004
	KPidTransportPriority                ParameterIDT = 0x0049
	KPidProtocolVersion                  ParameterIDT = 0x0015
	KPidVendorID                         ParameterIDT = 0x0016
	KPidUnicastLocator                   ParameterIDT = 0x002f
	KPidMulticastLocator                 ParameterIDT = 0x0030
	KPidMulticastIPAddress               ParameterIDT = 0x0011
	KPidDefaultUnicastLocator            ParameterIDT = 0x0031
	KPidDefaultMulticastLocator          ParameterIDT = 0x0048
	KPidMetatrafficUnicastLocator        ParameterIDT = 0x0032
	KPidMetatrafficMulticastLocator      ParameterIDT = 0x0033
	KPidDefaultUnicastIPAddress          ParameterIDT = 0x000c
	KPidDefaultUnicastPort               ParameterIDT = 0x000e
	KPidMetatrafficUnicastIPAddress      ParameterIDT = 0x0045
	KPidMetatrafficUnicastPort           ParameterIDT = 0x000d
	KPidMetatrafficMulticastIPAddress    ParameterIDT = 0x000b
	KPidMetatrafficMulticastPort         ParameterIDT = 0x0046
	KPidExpectsInlineQos                 ParameterIDT = 0x0043
	KPidParticipantManualLivelinessCount ParameterIDT = 0x0034
	KPidParticipantBuiltinEndpoints      ParameterIDT = 0x0044
	KPidParticipantLeaseDuration         ParameterIDT = 0x0002
	KPidContentFilterProperty            ParameterIDT = 0x0035
	KPidParticipantGUID                  ParameterIDT = 0x0050
	KPidParticipantEntityID              ParameterIDT = 0x0051
	KPidGroupGUID                        ParameterIDT = 0x0052
	KPidGroupEntityID                    ParameterIDT = 0x0053
	KPidBuiltinEndpointSet               ParameterIDT = 0x0058
	KPidPropertyList                     ParameterIDT = 0x0059
	KPidTypeMaxSizeSerialized            ParameterIDT = 0x0060
	KPidEntityName                       ParameterIDT = 0x0062
	KPidTypeIDV1                         ParameterIDT = 0x0069
	KPidKeyHash                          ParameterIDT = 0x0070
	KPidStatusInfo                       ParameterIDT = 0x0071
	KPidTypeObjectV1                     ParameterIDT = 0x0072
	KPidEndpointGUID                     ParameterIDT = 0x005a
	KPidIdentityToken                    ParameterIDT = 0x1001
	KPidPremissionsToken                 ParameterIDT = 0x1002
	KPidDataTags                         ParameterIDT = 0x1003
	KPidEndpointSecurityInfo             ParameterIDT = 0x1004
	KPidParticipantSecurityInfo          ParameterIDT = 0x1005
	KPidIdentityStatusToken              ParameterIDT = 0x1006
	KPidPersistenceGUID                  ParameterIDT = 0x8002
	KPidRelatedSampleIdentity            ParameterIDT = 0x800f
	KPidDataRepresentation               ParameterIDT = 0x0073
	KPidTypeConsistencyEnforcement       ParameterIDT = 0x0074
	KPidTypeInfoRmation                  ParameterIDT = 0x0075
	KPidDisablePositiveAcks              ParameterIDT = 0x8005
	//PID_RELATED_SAMPLE_IDENTITY = 0x0083
)

type ParameterT struct {
	Pid    ParameterIDT
	Length uint16
}

func NewParameterT(pid ParameterIDT, length uint16) *ParameterT {
	return &ParameterT{
		Pid:    pid,
		Length: length,
	}
}

var (
	KDefaultParameterT ParameterT
)

func init() {
	KDefaultParameterT = *NewParameterT(KPidPad, 0)
}

type ParameterKeyT struct {
	ParameterT
	Key common.InstanceHandleT
}

func NewParameterKey(pid ParameterIDT, length uint16) *ParameterKeyT {
	return &ParameterKeyT{
		ParameterT: ParameterT{
			Pid:    pid,
			Length: length,
		},
	}
}

type ParamaterLocatorT struct {
	ParameterT
	Locator common.Locator
}

func NewParamaterLocator(pid ParameterIDT, len uint16, loc *common.Locator) *ParamaterLocatorT {
	var parameter ParamaterLocatorT
	parameter.ParameterT = *NewParameterT(pid, len)
	parameter.Locator = *loc
	return &parameter
}

const (
	KParameterLocatorLength uint16 = 24
)

type ParameterStringT struct {
	ParameterT
	// Name. <br> By default, empty string.
	Name string
}

func NewParameterString(pid ParameterIDT, len uint16, strin string) *ParameterStringT {
	var parameter ParameterStringT
	parameter.ParameterT = *NewParameterT(pid, len)
	parameter.Name = strin
	return &parameter
}

type ParameterPortT struct {
	ParameterT
	Port uint32
}

const (
	KParameterPortLength uint32 = 4
)

type ParameterGuidT struct {
	ParameterT
	Guid common.GUIDT
}

func NewParameterGuid(pid ParameterIDT, len uint16, guidIn *common.GUIDT) *ParameterGuidT {
	var parameter ParameterGuidT
	parameter.ParameterT = *NewParameterT(pid, len)
	parameter.Guid = *guidIn
	return &parameter
}

const (
	KParameterGuidLength uint16 = 16
)

type ParameterProtocolVersionT struct {
	ParameterT
	ProtocolVersion common.ProtocolVersionT
}

func NewParameterProtocolVersion(pid ParameterIDT, len uint16) *ParameterProtocolVersionT {
	var protocol ParameterProtocolVersionT
	protocol.ParameterT = *NewParameterT(pid, len)
	protocol.ProtocolVersion = common.KProtocolVersion
	return &protocol
}

const (
	KParameterProtocolLength uint16 = 4
)

type ParameterVendorIDT struct {
	ParameterT
	VendorID common.VendorIDT
}

func NewParameterVendorIDT(pid ParameterIDT, len uint16) *ParameterVendorIDT {
	var parameter ParameterVendorIDT
	parameter.ParameterT = *NewParameterT(pid, len)
	parameter.VendorID = common.KVendorIDTeProsima
	return &parameter
}

const (
	KParameterVendorLength uint16 = 4
)

type ParameterIP4AddressT struct {
	ParameterT

	// Address <br> By default [0,0,0,0].
	Address [4]common.Octet
}

const (
	KParameterIP4Length uint32 = 4
)

type ParameterBoolT struct {
	ParameterT
	// By default, false.
	Value bool
}

const (
	KParameterBoolLength uint32 = 4
)

type ParameterStatusInfoT struct {
	ParameterT
	// By default, 0.
	Status uint8
}

const (
	KParameterStatusInfoLength uint32 = 4
)

type ParameterCountT struct {
	ParameterT
	Count common.CountT
}

const (
	KParameterCountLength uint32 = 4
)

type ParameterEntityIDT struct {
	ParameterT
	EntityID common.EntityIDT
}

const (
	KParameterEntityIDLength uint32 = 4
)

type ParameterTimeT struct {
	ParameterT
	Time common.Time
}

const (
	KParameterTimeLength uint16 = 8
)

func NewParameterTimeT(pid ParameterIDT, len uint16) *ParameterTimeT {
	var param ParameterTimeT
	param.ParameterT = *NewParameterT(pid, len)
	return &param
}

type ParameterBuiltinEndpointSetT struct {
	ParameterT
	EndpointSet common.BuiltinEndpointSet
}

const (
	KParameterBuiltinEndpointsetLength uint16 = 4
)

func NewParameterBuiltinEndpointSetT(pid ParameterIDT, len uint16) *ParameterBuiltinEndpointSetT {
	var param ParameterBuiltinEndpointSetT
	param.ParameterT = *NewParameterT(pid, len)
	return &param
}

type ParameterPropertyT struct {
	data []common.Octet
}

const (
	KParameterPropertyPersistenceGuid  = "PID_PERSISTENCE_GUID"
	KParameterPropertyParticipantType  = "PARTICIPANT_TYPE"
	KParameterPropertyDsVersion        = "DS_VERSION"
	KPArameterPropertyCurrentDSVersion = "2.0"
)

type ParameterPropertyListT struct {
	ParameterT
	properties  common.SerializedPayloadT
	nproperties uint32
	limitSize   bool
	ptr         []common.Octet
	value       ParameterPropertyT
}

func (paramList *ParameterPropertyListT) Size() uint32 {
	return paramList.nproperties
}

func (paramList *ParameterPropertyListT) Clear() {
	paramList.Length = 0
	paramList.nproperties = 0
}

func NewParameterPropertyListT(size uint32) *ParameterPropertyListT {
	var property ParameterPropertyListT
	property.properties = common.CreateSerializedPayload()
	property.properties.Reserve(size)
	property.nproperties = 0
	if size == 0 {
		property.limitSize = false
	} else {
		property.limitSize = true
	}
	property.ParameterT = *NewParameterT(KPidPropertyList, 0)
	return &property
}

type ParameterSampleIdentityT struct {
	ParameterT
	SampleID common.SampleIdentityT
}

const (
	KParameterSampleIdentidyLength uint32 = 24
)
