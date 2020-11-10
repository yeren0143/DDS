package attributes

import (
	. "common"
	. "fastrtps/rtps/flowcontrol"
	. "fastrtps/rtps/resources"
	. "fastrtps/transport"
)

type DiscoveryProtocol_t int8

const (
	// NONE: NO discovery whatsoever would be used.
	// Publisher and Subscriber defined with the same topic name would NOT be linked.
	// All matching must be done manually through the addReaderLocator, addReaderProxy, addWriterProxy methods.
	NONE DiscoveryProtocol_t = iota

	//Discovery works according to 'The Real-time Publish-Subscribe Protocol(RTPS) DDS
	//Interoperability Wire Protocol Specification'
	SIMPLE

	//A user defined PDP subclass object must be provided in the attributes that deals with the discovery.
	//Framework is not responsible of this object lifetime.
	EXTERNAL

	//The participant will behave as a client concerning discovery operation.
	//Server locators should be specified as attributes.
	CLIENT

	//The participant will behave as a server concerning discovery operation.
	//Discovery operation is volatile (discovery handshake must take place if shutdown
	SERVER

	//The participant will behave as a server concerning discovery operation.
	//Discovery operation persist on a file (discovery handshake wouldn't repeat if shutdown
	BACKUP
)

type ParticipantFilteringFlags int8

const (
	NO_FILTER                ParticipantFilteringFlags = 0
	FILTER_DIFFERENT_HOST    ParticipantFilteringFlags = 0x1
	FILTER_DIFFERENT_PROCESS ParticipantFilteringFlags = 0x2
	FILTER_SAME_PROCESS      ParticipantFilteringFlags = 0x4
)

const (
	BUILTIN_DATA_MAX_SIZE uint32 = 512
)

// type PDPFactory interface {
// 	CreatePDPInstance(BuiltinProtocols)
// }

// Class SimpleEDPAttributes, to define the attributes of the Simple Endpoint Discovery Protocol.
type SimpleEDPAttributes struct {
	Use_PublicationWriterANDSubscriptionReader                         bool
	Use_PublicationReaderANDSubscriptionWriter                         bool
	Enable_builtin_secure_publications_writer_and_subscriptions_reader bool
	Enable_builtin_secure_subscriptions_writer_and_publications_reader bool
}

func NewSimpleEDPAttributes() *SimpleEDPAttributes {
	return &SimpleEDPAttributes{
		Use_PublicationWriterANDSubscriptionReader:                         true,
		Use_PublicationReaderANDSubscriptionWriter:                         true,
		Enable_builtin_secure_publications_writer_and_subscriptions_reader: true,
		Enable_builtin_secure_subscriptions_writer_and_publications_reader: true,
	}
}

//defines the behavior of the RTPSParticipant initial announcements.
type InitialAnnouncementConfig struct {
	Count  uint32
	Period Duration_t
}

func NewDefaultInitialAnnouncementConfig() InitialAnnouncementConfig {
	return InitialAnnouncementConfig{
		Count:  5,
		Period: CreateDuration(0, 100000000),
	}
}

type DiscoverySettings struct {
	DiscoveryProtocol DiscoveryProtocol_t

	//If set to true, SimpleEDP would be used.
	Use_SIMPLE_EndpointDiscoveryProtocol bool

	//If set to true, StaticEDP based on an XML file would be implemented.
	Use_STATIC_EndpointDiscoveryProtocol bool

	//indicating how much time remote RTPSParticipants should consider this RTPSParticipant alive.
	LeaseDuration Duration_t

	//The period for the RTPSParticipant to send its Discovery Message to all other discovered
	//RTPSParticipants as well as to all Multicast ports.
	LeaseDuration_AnnouncementPeriod Duration_t

	Initial_Announcements InitialAnnouncementConfig
	SimpleEDP             *SimpleEDPAttributes

	//function that returns a PDP object (only if EXTERNAL selected)
	PDPFactory interface{}

	DiscoveryServer_client_syncperiod Duration_t
}

func NewDiscoverySettings() *DiscoverySettings {
	var discoverySettings DiscoverySettings
	discoverySettings.DiscoveryProtocol = SIMPLE
	discoverySettings.Use_SIMPLE_EndpointDiscoveryProtocol = true
	discoverySettings.Use_STATIC_EndpointDiscoveryProtocol = false
	discoverySettings.LeaseDuration = CreateDuration(20, 0)
	discoverySettings.LeaseDuration_AnnouncementPeriod = CreateDuration(3, 0)
	discoverySettings.Initial_Announcements = NewDefaultInitialAnnouncementConfig()
	discoverySettings.SimpleEDP = NewSimpleEDPAttributes()
	discoverySettings.DiscoveryServer_client_syncperiod = CreateDuration(0, 450*1000000)

	return &discoverySettings
}

type TypeLookupSettings struct {
	Use_Client bool
	Use_Server bool
}

func NewTypeLookupSettings() TypeLookupSettings {
	return TypeLookupSettings{
		Use_Client: false,
		Use_Server: false,
	}
}

type BuiltinAttributes struct {
	DiscoveryConfig                 *DiscoverySettings
	Use_WriterLivelinessProtocol    bool
	TypeLookup_Config               TypeLookupSettings
	MetatrafficUnicastLocatorList   LocatorList
	MetatrafficMulticastLocatorList LocatorList
	InitialPeersList                LocatorList
	ReaderHostoryMemoryPolicy       MemoryManagementPolicy
	ReaderPayloadSize               uint32 //Maximum payload size for builtin readers
	WriterHistoryMemoryPolicy       MemoryManagementPolicy
	WriterPayloadSize               uint32 //Maximum payload size for builtin writers
	Mutation_Tries                  uint32 //Mutation tries if the port is being used.
	Avoid_Builtin_Multicast         bool   //Set to true to avoid multicast traffic on builtin endpoints
}

func NewBuiltinAttributes() *BuiltinAttributes {
	// var att BuiltinAttributes
	// att.DiscoveryConfig = NewDiscoveryConfig()
	// att.Use_WriterLivelinessProtocol = true
	// att.TypeLookup_Config = NewTypeLookupSettings()
	// att.MetatrafficUnicastLocatorList = NewLocatorList()
	// att.MetatrafficMulticastLocatorList = NewLocatorList()
	// att.InitialPeersList = NewLocatorList()
	// att.ReaderHostoryMemoryPolicy = PREALLOCATED_WITH_REALLOC_MEMORY_MODE
	// att.ReaderPayloadSize = BUILTIN_DATA_MAX_SIZE
	// att.WriterHistoryMemoryPolicy = PREALLOCATED_WITH_REALLOC_MEMORY_MODE
	// att.WriterPayloadSize = BUILTIN_DATA_MAX_SIZE
	// att.Mutation_Tries = 100
	// att.Avoid_Builtin_Multicast = true

	att := BuiltinAttributes{
		DiscoveryConfig:                 NewDiscoverySettings(),
		Use_WriterLivelinessProtocol:    true,
		TypeLookup_Config:               NewTypeLookupSettings(),
		MetatrafficUnicastLocatorList:   NewLocatorList(),
		MetatrafficMulticastLocatorList: NewLocatorList(),
		InitialPeersList:                NewLocatorList(),
		ReaderHostoryMemoryPolicy:       PREALLOCATED_WITH_REALLOC_MEMORY_MODE,
		ReaderPayloadSize:               BUILTIN_DATA_MAX_SIZE,
		WriterHistoryMemoryPolicy:       PREALLOCATED_WITH_REALLOC_MEMORY_MODE,
		WriterPayloadSize:               BUILTIN_DATA_MAX_SIZE,
		Mutation_Tries:                  100,
		Avoid_Builtin_Multicast:         true,
	}

	return &att
}

type RTPSParticipantAttributes struct {
	name                        string
	DefaultUnicastLocatorList   LocatorList
	DefaultMulticastLocatorList LocatorList
	SendSocketBufferSize        uint32
	ListenSocketBufferSize      uint32
	Prefix                      *GuidPrefix_t
	Builtin                     BuiltinAttributes
	Port                        PortParameters
	UserData                    []Octet
	ParticipantID               int32

	//!Throughput controller parameters. Leave default for uncontrolled flow.
	ThroghputController ThroghputControllerDescriptor

	//!User defined transports to use alongside or in place of builtins.
	UserTransports []*TransportDescriptorInterface

	//!Set as false to disable the default UDPv4 implementation.
	UseBuiltinTransports bool

	//!Holds allocation limits affecting collections managed by a participant.
	Allocation *RTPSParticipantAllocationAttributes

	Properties *PropertyPolicy
}

func NewRTPSParticipantAttributes() *RTPSParticipantAttributes {
	participantAttributes := RTPSParticipantAttributes{
		name:                        "RTPSParticipant",
		ParticipantID:               -1,
		UseBuiltinTransports:        true,
		DefaultUnicastLocatorList:   NewLocatorList(),
		DefaultMulticastLocatorList: NewLocatorList(),
		Prefix:                      NewGuiPrefix(),
		Port:                        NewDefaultPortParameters(),
		ThroghputController:         ThroghputControllerDescriptor{},
		Allocation:                  NewRTPSParticipantAllocationAttributes(),
		Properties:                  NewPropertyPolicy(),
	}

	return &participantAttributes
}
