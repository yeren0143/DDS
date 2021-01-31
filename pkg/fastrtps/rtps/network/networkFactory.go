package network

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/transport"
	"github.com/yeren0143/DDS/fastrtps/utils"
	"log"
)

// NetFactory Provides the FastRTPS library with abstract resources, which
// in turn manage the SEND and RECEIVE operations over some transport.
// Once a transport is registered, it becomes invisible to the library
// and is abstracted away for good.
type NetFactory struct {
	maxMessageSizeBetweenTransports uint32
	minSendBufferSize               uint32
	registeredTransports            []transport.ITransport
}

//RegisterTransport Allows registration of a transport statically, by specifying the transport type and
// its associated descriptor type. This is particularly useful for user-defined transports.
func (factory *NetFactory) RegisterTransport(descriptor transport.ITransportDescriptor) bool {
	wasRegistered := false
	minSendBufferSize := ^uint32(0)

	trans := descriptor.CreateTransport()

	if trans != nil {
		if trans.Init() == true {
			minSendBufferSize = trans.Configuration().MinSendBufferSize()
			factory.registeredTransports = append(factory.registeredTransports, trans)
			wasRegistered = true
		}

		if wasRegistered {
			if descriptor.MaxMessageSize() < factory.maxMessageSizeBetweenTransports {
				factory.maxMessageSizeBetweenTransports = descriptor.MaxMessageSize()
			}

			if minSendBufferSize < factory.minSendBufferSize {
				factory.minSendBufferSize = minSendBufferSize
			}
		}
	}

	return wasRegistered
}

//NumberOfRegisteredTransports return registered transport number
func (factory *NetFactory) NumberOfRegisteredTransports() int {
	return len(factory.registeredTransports)
}

func (factory *NetFactory) GetMaxMessageSizeBetweenTransports() uint32 {
	return factory.maxMessageSizeBetweenTransports
}

//GetDefaultMetatrafficUnicastLocators adds locators to the metatraffic unicast list.
func (factory *NetFactory) GetDefaultMetatrafficUnicastLocators(locators *common.LocatorList, port uint32) bool {
	result := false
	for _, transport := range factory.registeredTransports {
		result = transport.GetDefaultMetatrafficUnicastLocators(locators, port)
	}
	return result
}

//GetDefaultMetatrafficMulticastLocators adds locators to the metatraffic multicast list.
func (factory *NetFactory) GetDefaultMetatrafficMulticastLocators(locators *common.LocatorList, multicastPort uint32) bool {
	result := false

	var shmTransport transport.ITransport

	for _, transport := range factory.registeredTransports {
		// For better fault-tolerance reasons, SHM multicast metatraffic is avoided if it is already provided
		// by another transport
		if transport.Kind() != common.KLocatorKindShm {
			result = result || transport.GetDefaultMetatrafficMulticastLocators(locators, multicastPort)
		} else {
			shmTransport = transport
		}
	}

	if locators.Length() == 0 && shmTransport != nil {
		result = result || shmTransport.GetDefaultMetatrafficMulticastLocators(locators, multicastPort)
	}

	return result
}

func calculateWellKnownPort(domainID uint32, att *attributes.RTPSParticipantAttributes, multicast bool) uint16 {
	port := uint32(att.Port.PortBase) + uint32(att.Port.DomainIDGain)*domainID
	if multicast {
		port += uint32(att.Port.Offsetd2)
	} else {
		port += uint32(att.Port.Offsetd3) + uint32(att.Port.ParticipantIDGain)*uint32(att.ParticipantID)
	}

	if port > 65535 {
		msg := `Calculated port number is too high. Probably the domainId is over 232, there are 
		too much participants created or portBase is too high.`
		log.Fatal(msg)
	}

	return uint16(port)
}

// GetDefaultUnicastLocators adds locators to the default unicast configuration.
func (factory *NetFactory) GetDefaultUnicastLocators(domainID uint32, locators *common.LocatorList,
	att *attributes.RTPSParticipantAttributes) bool {
	result := false
	wellKnownPort := calculateWellKnownPort(domainID, att, false)
	for _, transport := range factory.registeredTransports {
		result = result || transport.GetDefaultUnicastLocators(locators, uint32(wellKnownPort))
	}
	return result
}

// NormalizedLocators ...
func (factory *NetFactory) NormalizedLocators(locators *common.LocatorList) *common.LocatorList {
	var normalizedLocators common.LocatorList

	for _, loc := range locators.Locators {
		normalized := false
		for _, transport := range factory.registeredTransports {
			// Check if the locator is supported and filter unicast locators.
			if transport.IsLocatorSupported(&loc) &&
				(utils.IsMulticast(&loc) || transport.IsLocatorAllowed(&loc)) {
				normalizedLocators.Append(transport.NormalizeLocator(&loc))
				normalized = true
			}
		}

		if normalized == false {
			newloc := loc
			normalizedLocators.PushBack(&newloc)
		}
	}

	*locators = normalizedLocators
	return &normalizedLocators
}

// FillDefaultLocatorPort fills the locator with the default unicast configuration.
func (factory *NetFactory) FillDefaultLocatorPort(domainID uint32, locator *common.Locator,
	att *attributes.RTPSParticipantAttributes, isMulticast bool) bool {
	result := false
	wellKnownPort := calculateWellKnownPort(domainID, att, isMulticast)
	for _, transport := range factory.registeredTransports {
		if transport.IsLocatorSupported(locator) {
			result = result || transport.FillUnicastLocator(locator, uint32(wellKnownPort))
		}
	}
	return result
}

//FillMetatrafficMulticastLocator fills the locator with the metatraffic multicast configuration.
func (factory *NetFactory) FillMetatrafficMulticastLocator(locator *common.Locator, multicastPort uint32) bool {
	result := false
	for _, transport := range factory.registeredTransports {
		if transport.IsLocatorSupported(locator) {
			result = result || transport.FillMetatrafficMulticastLocator(locator, multicastPort)
		}
	}

	return result
}

//FillMetatrafficUnicastLocator fills the locator with the metatraffic unicast configuration.
func (factory *NetFactory) FillMetatrafficUnicastLocator(locator *common.Locator, unicastPort uint32) bool {
	result := false
	for _, transport := range factory.registeredTransports {
		if transport.IsLocatorSupported(locator) {
			result = result || transport.FillMetatrafficUnicastLocator(locator, unicastPort)
		}
	}
	return result
}

//TransformRemoteLocator transforms a remote locator into a locator optimized for local communications.
//If the remote locator corresponds to one of the local interfaces, it is converted
//to the corresponding local address.
//return false if the input locator is not supported/allowed by any of the registered transports,
//true otherwise.
func (factory *NetFactory) TransformRemoteLocator(remote, result *common.Locator) bool {
	for _, trans := range factory.registeredTransports {
		if _, ok := trans.TransformRemoteLocator(remote); ok {
			return true
		}
	}

	return false
}

// ConfigureInitialPeerLocator the locator with the initial peer configuration
func (factory *NetFactory) ConfigureInitialPeerLocator(domainID uint32, locator *common.Locator,
	att *attributes.RTPSParticipantAttributes) bool {
	result := false
	for _, transport := range factory.registeredTransports {
		if transport.IsLocatorSupported(locator) {
			result = result || transport.ConfigureInitialPeerLocator(locator, att.Port, domainID,
				att.Builtin.InitialPeersList)
		}
	}

	return result
}

// Walks over the list of transports, opening every possible channel that we can listen to
// from the given locator, and returns a vector of Receiver Resources for this goal.
// @param local Locator from which to listen.
// @param returnResourcesList that will be filled with the created ReceiverResources.
// @param receiver_max_message_size Max message size allowed by the message receiver.
func (factory *NetFactory) BuildReceiverResources(local *common.Locator, receiverMaxMsgSize uint32) (bool, []*ReceiverResource) {
	returnedValue := false
	var returnResourcesList []*ReceiverResource
	for _, transport := range factory.registeredTransports {
		if transport.IsLocatorSupported(local) {
			if !transport.IsInputChannelOpen(local) {
				var maxRecvBufferSize uint32
				if transport.MaxRecvBufferSize() > receiverMaxMsgSize {
					maxRecvBufferSize = receiverMaxMsgSize
				} else {
					maxRecvBufferSize = transport.MaxRecvBufferSize()
				}

				newReceiverResource := NewReceiverResource(transport, local, maxRecvBufferSize)

				if newReceiverResource.Valid {
					returnResourcesList = append(returnResourcesList, newReceiverResource)
					returnedValue = true
				}
			} else {
				returnedValue = true
			}
		}
	}

	return returnedValue, returnResourcesList
}

//NewNetworkFactory create network factory
func NewNetworkFactory() *NetFactory {
	return &NetFactory{
		maxMessageSizeBetweenTransports: ^uint32(0),
		minSendBufferSize:               ^uint32(0),
	}
}
