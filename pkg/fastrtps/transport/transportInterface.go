package transport

import (
	common "github.com/yeren0143/DDS/common"
	//network "github.com/yeren0143/DDS/fastrtps/network"
)

// const to descripe transport layer
const (
	CMaximumMessageSize       = 65500
	CMaximumInitialPeersRange = 4
	CMinimumSocketBuffer      = 65536
	CIPv4AddressAny           = "0.0.0.0"
	CIPv6AddressAny           = "::"
)

// SenderSourceList is a slice of SenderResource
//type SenderSourceList []*network.SenderResource
type SenderSourceList []interface{}

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

	//Must report whether the given locator is allowed by this transport.
	IsLocatorAllowed(*common.Locator) bool

	//Returns the locator describing the main (most general) channel that can write to the provided remote locator.
	RemoteToMainLocal(remote *common.Locator) *common.Locator

	// Transforms a remote locator into a locator optimized for local communications.
	// If the remote locator corresponds to one of the local interfaces, it is converted
	// to the corresponding local address.
	//false if the input locator is not supported/allowed by this transport, true otherwise.
	TransformRemoteLocator(remote *common.Locator, locator *common.Locator) bool

	//Must open the channel that maps to/from the given locator. This method must allocate, reserve and mark
	//any resources that are needed for said channel.
	OpenOutputChannel(senderList SenderSourceList, locator *common.Locator) bool

	//return transport kind
	Kind() int32
}
