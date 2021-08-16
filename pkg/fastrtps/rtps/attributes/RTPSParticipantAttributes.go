package attributes

import (
	"dds/common"
	"dds/fastrtps/rtps/flowcontrol"
	"dds/fastrtps/rtps/resources"
	"dds/fastrtps/rtps/transport"
)

// DiscoveryProtocolT ...
type DiscoveryProtocolT int8

const (
	// KDisPNone NO discovery whatsoever would be used.
	// Publisher and Subscriber defined with the same topic name would NOT be linked.
	// All matching must be done manually through the addReaderLocator, addReaderProxy, addWriterProxy methods.
	KDisPNone DiscoveryProtocolT = iota

	// KDisSimple Discovery works according to 'The Real-time Publish-Subscribe Protocol(RTPS) DDS
	// Interoperability Wire Protocol Specification'
	KDisPSimple

	// KDisPExternal A user defined PDP subclass object must be provided in the attributes that deals with the discovery.
	// Framework is not responsible of this object lifetime.
	KDisPExternal

	// KDisPClient The participant will behave as a client concerning discovery operation.
	// Server locators should be specified as attributes.
	KDisPClient

	// KDisPServer participant will behave as a server concerning discovery operation.
	//Discovery operation is volatile (discovery handshake must take place if shutdown
	KDisPServer

	// KDisPBackup participant will behave as a server concerning discovery operation.
	//Discovery operation persist on a file (discovery handshake wouldn't repeat if shutdown
	KDisPBackup
)

// ParticipantFilteringFlags ...
type ParticipantFilteringFlags int8

// ...
const (
	KNoFilter               ParticipantFilteringFlags = 0
	KFilterDifferentHost    ParticipantFilteringFlags = 0x1
	KFilterDifferentProcess ParticipantFilteringFlags = 0x2
	KFilterSameProcess      ParticipantFilteringFlags = 0x4
	KBuiltinDataMaxSize     uint32                    = 512
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
	Period common.DurationT
}

// NewDefaultInitialAnnouncementConfig create AnnouncementConfig with default config
func NewDefaultInitialAnnouncementConfig() *InitialAnnouncementConfig {
	return &InitialAnnouncementConfig{
		Count:  5,
		Period: common.DurationT{Seconds: 0, Nanosec: 100000000},
	}
}

// DiscoverySettings define discovery config
type DiscoverySettings struct {
	Protocol DiscoveryProtocolT

	//If set to true, SimpleEDP would be used.
	UseSimpleEndpoint bool

	//If set to true, StaticEDP based on an XML file would be implemented.
	UseStaticEndpoint bool

	//indicating how much time remote RTPSParticipants should consider this RTPSParticipant alive.
	LeaseDuration common.DurationT

	//The period for the RTPSParticipant to send its Discovery Message to all other discovered
	//RTPSParticipants as well as to all Multicast ports.
	LeaseDurationAnnouncementPeriod common.DurationT

	InitialAnnouncements *InitialAnnouncementConfig
	SimpleEDP            *SimpleEDPAttributes

	//function that returns a PDP object (only if EXTERNAL selected)
	PDPFactory interface{}

	DiscoveryServerClientSyncPeriod common.DurationT

	//Discovery Server settings, only needed if use_CLIENT_DiscoveryProtocol=true
	DiscoveryServers RemoteServerList

	//! Filtering participants out depending on location
	IgnoreParticipantFlags ParticipantFilteringFlags
}

// NewDiscoverySettings create DiscoverySetting with default value
func NewDiscoverySettings() *DiscoverySettings {
	var discoverySettings DiscoverySettings
	discoverySettings.Protocol = KDisPSimple
	discoverySettings.UseSimpleEndpoint = true
	discoverySettings.UseStaticEndpoint = false
	discoverySettings.LeaseDuration = common.DurationT{Seconds: 20, Nanosec: 0}
	discoverySettings.LeaseDurationAnnouncementPeriod = common.DurationT{Seconds: 3, Nanosec: 0}
	discoverySettings.InitialAnnouncements = NewDefaultInitialAnnouncementConfig()
	discoverySettings.SimpleEDP = NewSimpleEDPAttributes()
	discoverySettings.DiscoveryServerClientSyncPeriod = common.DurationT{Seconds: 0, Nanosec: 450 * 1000000}
	discoverySettings.IgnoreParticipantFlags = KNoFilter

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
	MetatrafficUnicastLocatorList   *common.LocatorList
	MetatrafficMulticastLocatorList *common.LocatorList
	InitialPeersList                *common.LocatorList
	ReaderHostoryMemoryPolicy       resources.MemoryManagementPolicy
	ReaderPayloadSize               uint32 //Maximum payload size for builtin readers
	WriterHistoryMemoryPolicy       resources.MemoryManagementPolicy
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
		MetatrafficUnicastLocatorList:   common.NewLocatorList(),
		MetatrafficMulticastLocatorList: common.NewLocatorList(),
		InitialPeersList:                common.NewLocatorList(),
		ReaderHostoryMemoryPolicy:       resources.KPreallocatedWithReallocMemoryMode,
		ReaderPayloadSize:               KBuiltinDataMaxSize,
		WriterHistoryMemoryPolicy:       resources.KPreallocatedWithReallocMemoryMode,
		WriterPayloadSize:               KBuiltinDataMaxSize,
		MutationTries:                   100,
		AvoidBuiltinMulticast:           true,
	}
}

// RTPSParticipantAttributes ...
type RTPSParticipantAttributes struct {
	Name                        string
	DefaultUnicastLocatorList   *common.LocatorList
	DefaultMulticastLocatorList *common.LocatorList
	SendSocketBufferSize        uint32
	ListenSocketBufferSize      uint32
	Prefix                      *common.GUIDPrefixT
	Builtin                     *BuiltinAttributes
	Port                        *common.PortParameters
	UserData                    []common.Octet
	ParticipantID               int32

	//!Throughput controller parameters. Leave default for uncontrolled flow.
	ThroughputController *flowcontrol.ThroughputControllerDescriptor

	//!User defined transports to use alongside or in place of builtins.
	UserTransports []transport.ITransportDescriptor

	//!Set as false to disable the default UDPv4 implementation.
	UseBuiltinTransports bool

	//!Holds allocation limits affecting collections managed by a participant.
	Allocation *RTPSParticipantAllocationAttributes

	Properties *PropertyPolicyT
}

// NewRTPSParticipantAttributes ...
func NewRTPSParticipantAttributes() *RTPSParticipantAttributes {
	return &RTPSParticipantAttributes{
		Name:                        "RTPSParticipant",
		ParticipantID:               -1,
		UseBuiltinTransports:        true,
		DefaultUnicastLocatorList:   common.NewLocatorList(),
		DefaultMulticastLocatorList: common.NewLocatorList(),
		Prefix:                      common.NewGUIDPrefix(),
		Builtin:                     NewBuiltinAttributes(),
		Port:                        common.NewDefaultPortParameters(),
		ThroughputController:        flowcontrol.NewThroughputControllerDescriptor(),
		Allocation:                  NewRTPSParticipantAllocationAttributes(),
		Properties:                  NewPropertyPolicy(),
	}
}
