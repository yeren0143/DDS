package transport

import (
	common "github.com/yeren0143/DDS/common"
)

// const to descripe transport layer
const (
	KMaximumMessageSize       = 65500
	KMaximumInitialPeersRange = 4
	KMinimumSocketBuffer      = 65536
	kIPv4AddressAny           = "0.0.0.0"
	KIPv6AddressAny           = "::"
)

// ITransport against which to implement a transport layer, decoupled from FastRTPS internals.
// TransportInterface expects the user to implement a logical equivalence between Locators and protocol-specific "channels".
// This equivalence can be narrowing: For instance in UDP/IP, a port would take the role of channel, and several different
// locators can map to the same port, and hence the same channel.
type ITransport interface {
	//Initialize this transport. This method will prepare all the internals of the transport.
	Init() bool

	// Must report whether the input channel associated to this locator is open. Channels must either be
	// fully closed or fully open, so that "open" and "close" operations are whole and definitive.
	IsInputChannelOpen(*common.Locator) bool

	//Must report whether the given locator is supported by this transport (typically inspecting its "kind" value).
	IsLocatorSupported(*common.Locator) bool

	GetDefaultUnicastLocators(locators *common.LocatorList, unicastPost uint32) bool

	GetDefaultMetatrafficMulticastLocators(locators *common.LocatorList, multicastPort uint32) bool
	GetDefaultMetatrafficUnicastLocators(locators *common.LocatorList, multicastPort uint32) bool

	//Must report whether the given locator is allowed by this transport.
	IsLocatorAllowed(*common.Locator) bool

	//Returns the locator describing the main (most general) channel that can write to the provided remote locator.
	RemoteToMainLocal(remote *common.Locator) *common.Locator

	// Transforms a remote locator into a locator optimized for local communications.
	// If the remote locator corresponds to one of the local interfaces, it is converted
	// to the corresponding local address.
	// false if the input locator is not supported/allowed by this transport, true otherwise.
	TransformRemoteLocator(remote *common.Locator) (*common.Locator, bool)

	//Must open the channel that maps to/from the given locator. This method must allocate, reserve and mark
	//any resources that are needed for said channel.
	OpenOutputChannel(senderList SenderResourceList, locator *common.Locator) bool

	OpenInputChannel(locator *common.Locator, receiver ITransportReceiver, maxMsgSize uint32) bool

	/**
	 * Must close the channel that maps to/from the given locator.
	 * IMPORTANT: It MUST be safe to call this method even during a Receive operation on another thread. You must implement
	 * any necessary mutual exclusion and timeout mechanisms to make sure the channel can be closed without damage.
	 */
	CloseInputChannel(locator *common.Locator) bool

	//! Must report whether two locators map to the same internal channel.
	DoInputLocatorsMatch(left *common.Locator, right *common.Locator) bool

	Configuration() ITransportDescriptor

	// The maximum datagram size for reception supported by the transport
	MaxRecvBufferSize() uint32

	FillUnicastLocator(Locator *common.Locator, wellKnownPort uint32) bool

	FillMetatrafficMulticastLocator(locator *common.Locator, wellKnownPort uint32) bool
	FillMetatrafficUnicastLocator(locator *common.Locator, unicastPort uint32) bool

	NormalizeLocator(locator *common.Locator) *common.LocatorList

	ConfigureInitialPeerLocator(locator *common.Locator, portParams *common.PortParameters,
		domainID uint32, list *common.LocatorList) bool

	//return transport kind
	Kind() int32

	AddDefaultOutputLocator(defaultList *common.LocatorList)
}
