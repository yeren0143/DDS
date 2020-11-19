package attributes

import (
	. "github.com/yeren0143/DDS/common"
	. "github.com/yeren0143/DDS/fastrtps/rtps/flowcontrol"
	. "github.com/yeren0143/DDS/fastrtps/rtps/resources"
)

// DiscoveryProtocolT ...
type DiscoveryProtocolT int8

const (
	// CDisNone NO discovery whatsoever would be used.
	// Publisher and Subscriber defined with the same topic name would NOT be linked.
	// All matching must be done manually through the addReaderLocator, addReaderProxy, addWriterProxy methods.
	CDisNone DiscoveryProtocolT = iota

	// CDisSimple Discovery works according to 'The Real-time Publish-Subscribe Protocol(RTPS) DDS
	// Interoperability Wire Protocol Specification'
	CDisSimple

	// CDisExternal A user defined PDP subclass object must be provided in the attributes that deals with the discovery.
	// Framework is not responsible of this object lifetime.
	CDisExternal

	// CDisClient The participant will behave as a client concerning discovery operation.
	// Server locators should be specified as attributes.
	CDisClient

	// CDisServer participant will behave as a server concerning discovery operation.
	//Discovery operation is volatile (discovery handshake must take place if shutdown
	CDisServer

	// CDisBackup participant will behave as a server concerning discovery operation.
	//Discovery operation persist on a file (discovery handshake wouldn't repeat if shutdown
	CDisBackup
)

// ParticipantFilteringFlags ...
type ParticipantFilteringFlags int8

// ...
const (
	CNoFilter               ParticipantFilteringFlags = 0
	CFilterDifferentHost    ParticipantFilteringFlags = 0x1
	CFilterDifferentProcess ParticipantFilteringFlags = 0x2
	CFilterSameProcess      ParticipantFilteringFlags = 0x4
	CBuiltinDataMaxSize     uint32                    = 512
)

// type PDPFactory interface {
// 	CreatePDPInstance(BuiltinProtocols)
// }

// SimpleEDPAttributes define the attributes of the Simple Endpoint Discovery Protocol.
type SimpleEDPAttributes struct {
	UsePublicationWriterAndSubscriptionReader                   bool
	UsePublicationReaderAndSubscriptionWriter                   bool
	EnableBuiltinSecurePublicationsWriterAndSubscriptionsReader bool
	EnableBuiltinSecureSubscriptionsWriterAndPublicationsReader bool
}

// NewSimpleEDPAttributes create EDPAttributes
func NewSimpleEDPAttributes() *SimpleEDPAttributes {
	return &SimpleEDPAttributes{
		UsePublicationWriterAndSubscriptionReader:                   true,
		UsePublicationReaderAndSubscriptionWriter:                   true,
		EnableBuiltinSecurePublicationsWriterAndSubscriptionsReader: true,
		EnableBuiltinSecureSubscriptionsWriterAndPublicationsReader: true,
	}
}

// InitialAnnouncementConfig defines the behavior of the RTPSParticipant initial announcements.
type InitialAnnouncementConfig struct {
	Count  uint32
	Period DurationT
}

// NewDefaultInitialAnnouncementConfig create AnnouncementConfig with default config
func NewDefaultInitialAnnouncementConfig() *InitialAnnouncementConfig {
	return &InitialAnnouncementConfig{
		Count:  5,
		Period: DurationT{0, 100000000},
	}
}

// DiscoverySettings define discovery config
type DiscoverySettings struct {
	DiscoveryProtocol DiscoveryProtocolT

	//If set to true, SimpleEDP would be used.
	UseSimpleEndpointDiscoveryProtocol bool

	//If set to true, StaticEDP based on an XML file would be implemented.
	UseStaticEndpointDiscoveryProtocol bool

	//indicating how much time remote RTPSParticipants should consider this RTPSParticipant alive.
	LeaseDuration DurationT

	//The period for the RTPSParticipant to send its Discovery Message to all other discovered
	//RTPSParticipants as well as to all Multicast ports.
	LeaseDurationAnnouncementPeriod DurationT

	InitialAnnouncements *InitialAnnouncementConfig
	SimpleEDP            *SimpleEDPAttributes

	//function that returns a PDP object (only if EXTERNAL selected)
	PDPFactory interface{}

	DiscoveryServerClientSyncPeriod DurationT
}

// NewDiscoverySettings create DiscoverySetting with default value
func NewDiscoverySettings() *DiscoverySettings {
	var discoverySettings DiscoverySettings
	discoverySettings.DiscoveryProtocol = CDisSimple
	discoverySettings.UseSimpleEndpointDiscoveryProtocol = true
	discoverySettings.UseStaticEndpointDiscoveryProtocol = false
	discoverySettings.LeaseDuration = DurationT{20, 0}
	discoverySettings.LeaseDurationAnnouncementPeriod = DurationT{3, 0}
	discoverySettings.InitialAnnouncements = NewDefaultInitialAnnouncementConfig()
	discoverySettings.SimpleEDP = NewSimpleEDPAttributes()
	discoverySettings.DiscoveryServerClientSyncPeriod = DurationT{0, 450 * 1000000}

	return &discoverySettings
}

// TypeLookupSettings ...
type TypeLookupSettings struct {
	UseClient bool
	UseServer bool
}

// NewTypeLookupSettings ...
func NewTypeLookupSettings() *TypeLookupSettings {
	return &TypeLookupSettings{
		UseClient: false,
		UseServer: false,
	}
}

// BuiltinAttributes ...
type BuiltinAttributes struct {
	DiscoveryConfig                 *DiscoverySettings
	UseWriterLivelinessProtocol     bool
	TypeLookupConfig                *TypeLookupSettings
	MetatrafficUnicastLocatorList   *LocatorList
	MetatrafficMulticastLocatorList *LocatorList
	InitialPeersList                *LocatorList
	ReaderHostoryMemoryPolicy       MemoryManagementPolicy
	ReaderPayloadSize               uint32 //Maximum payload size for builtin readers
	WriterHistoryMemoryPolicy       MemoryManagementPolicy
	WriterPayloadSize               uint32 //Maximum payload size for builtin writers
	MutationTries                   uint32 //Mutation tries if the port is being used.
	AvoidBuiltinMulticast           bool   //Set to true to avoid multicast traffic on builtin endpoints
}

// NewBuiltinAttributes ...
func NewBuiltinAttributes() *BuiltinAttributes {
	return &BuiltinAttributes{
		DiscoveryConfig:                 NewDiscoverySettings(),
		UseWriterLivelinessProtocol:     true,
		TypeLookupConfig:                NewTypeLookupSettings(),
		MetatrafficUnicastLocatorList:   NewLocatorList(),
		MetatrafficMulticastLocatorList: NewLocatorList(),
		InitialPeersList:                NewLocatorList(),
		ReaderHostoryMemoryPolicy:       CPreallocatedWithReallocMemoryMode,
		ReaderPayloadSize:               CBuiltinDataMaxSize,
		WriterHistoryMemoryPolicy:       CPreallocatedWithReallocMemoryMode,
		WriterPayloadSize:               CBuiltinDataMaxSize,
		MutationTries:                   100,
		AvoidBuiltinMulticast:           true,
	}
}

// RTPSParticipantAttributes ...
type RTPSParticipantAttributes struct {
	Name                        string
	DefaultUnicastLocatorList   *LocatorList
	DefaultMulticastLocatorList *LocatorList
	SendSocketBufferSize        uint32
	ListenSocketBufferSize      uint32
	Prefix                      *GUIDPrefixT
	Builtin                     *BuiltinAttributes
	Port                        *PortParameters
	UserData                    []Octet
	ParticipantID               uint32

	//!Throughput controller parameters. Leave default for uncontrolled flow.
	ThroghputController *ThroghputControllerDescriptor

	//!User defined transports to use alongside or in place of builtins.
	//UserTransports []*TransportDescriptorInterface

	//!Set as false to disable the default UDPv4 implementation.
	UseBuiltinTransports bool

	//!Holds allocation limits affecting collections managed by a participant.
	Allocation *RTPSParticipantAllocationAttributes

	Properties *PropertyPolicy
}

// NewRTPSParticipantAttributes ...
func NewRTPSParticipantAttributes() *RTPSParticipantAttributes {
	return &RTPSParticipantAttributes{
		Name:                        "RTPSParticipant",
		ParticipantID:               ^uint32(0),
		UseBuiltinTransports:        true,
		DefaultUnicastLocatorList:   NewLocatorList(),
		DefaultMulticastLocatorList: NewLocatorList(),
		Prefix:                      NewGUIDPrefix(),
		Builtin:                     NewBuiltinAttributes(),
		Port:                        NewDefaultPortParameters(),
		ThroghputController:         NewThroghputControllerDescriptor(),
		Allocation:                  NewRTPSParticipantAllocationAttributes(),
		Properties:                  NewPropertyPolicy(),
	}
}
