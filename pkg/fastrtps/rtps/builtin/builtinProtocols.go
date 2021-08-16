package builtin

import (
	"log"

	"dds/common"
	"dds/fastrtps/rtps/attributes"
	"dds/fastrtps/rtps/builtin/discovery/pdp"
	"dds/fastrtps/rtps/builtin/discovery/protocol"
	"dds/fastrtps/rtps/builtin/liveliness"
	"dds/fastrtps/rtps/endpoint"
	"dds/fastrtps/rtps/network"
)

//Protocols that contains builtin endpoints implementing the discovery and liveliness protocols.
type Protocols struct {
	Att                             *attributes.BuiltinAttributes
	participantImpl                 protocol.IParticipant
	PDP                             protocol.IPDP
	WLP                             endpoint.IWlp
	TLM                             interface{} //TypeLookupManager
	MetatrafficMulticastLocatorList *common.LocatorList
	MetatrafficUnicastLocatorList   *common.LocatorList
	InitialPeersList                *common.LocatorList
	DiscoveryServers                []*attributes.RemoteServerAttributes //Known discovery and backup server container
}

func (protocol *Protocols) GetMetatrafficUnicastLocators() *common.LocatorList {
	return protocol.MetatrafficUnicastLocatorList
}

func (protocol *Protocols) transformServerRemoteLocators(nf *network.NetFactory) {
	for _, rs := range protocol.DiscoveryServers {
		for i, loc := range rs.MetrafficUnicastLocatorList.Locators {
			var localized common.Locator
			if nf.TransformRemoteLocator(&loc, &localized) {
				rs.MetrafficUnicastLocatorList.Locators[i] = localized
			}
		}
	}
}

func (protocol *Protocols) GetBuiltinAttributes() *attributes.BuiltinAttributes {
	return protocol.Att
}

func (protocol *Protocols) GetParticipant() protocol.IParticipant {
	return protocol.participantImpl
}

func (protocol *Protocols) UpdateMetatrafficLocators(loclist *common.LocatorList) bool {
	*protocol.MetatrafficUnicastLocatorList = *loclist
	return true
}

func (protocol *Protocols) GetMetatrafficMulticastLocatorList() *common.LocatorList {
	return protocol.MetatrafficMulticastLocatorList
}

func (protocol *Protocols) AnnounceParticipantState() {
	if protocol.PDP != nil {
		writeParam := common.KWriteParamDefault
		protocol.PDP.AnnounceParticipantState(false, false, &writeParam)
	} else {
		log.Fatalln("Trying to use BuiltinProtocols interfaces before initBuiltinProtocols call")
	}
}

func (protocol *Protocols) GetMetatrafficUnicastLocatorList() *common.LocatorList {
	return protocol.MetatrafficUnicastLocatorList
}

func (protocol *Protocols) GetInitialPeers() *common.LocatorList {
	return protocol.InitialPeersList
}

//InitBuiltinProtocol Initialize the builtin protocols.
func (protocol *Protocols) InitBuiltinProtocol(ppart protocol.IParticipant, att *attributes.BuiltinAttributes) bool {
	protocol.participantImpl = ppart
	protocol.Att = att
	protocol.MetatrafficMulticastLocatorList = att.MetatrafficMulticastLocatorList
	protocol.MetatrafficUnicastLocatorList = att.MetatrafficUnicastLocatorList
	protocol.InitialPeersList = att.InitialPeersList
	protocol.DiscoveryServers = att.DiscoveryConfig.DiscoveryServers

	protocol.transformServerRemoteLocators(ppart.NetworkFactory())
	allocation := ppart.GetAttributes().Allocation

	//PDP
	switch att.DiscoveryConfig.Protocol {
	case attributes.KDisPNone:
		log.Fatalln("No participant discovery protocol specified")
		return true
	case attributes.KDisPSimple:
		protocol.PDP = pdp.NewPDPSimple(protocol, allocation)
	case attributes.KDisPExternal:
		log.Fatalln("Flag only present for debugging purposes")
		return false
	case attributes.KDisPServer:
		log.Fatalln("Flag only present for debugging purposes")
	//TODO:: #if HAVE_SQLITE3
	default:
		log.Fatal("Unknown DiscoveryProtocol_t specified.")
	}

	if !protocol.PDP.Init(protocol.participantImpl) {
		log.Fatalln("Participant discovery configuration failed")
		return false
	}

	// WLP
	if protocol.Att.UseWriterLivelinessProtocol {
		protocol.WLP = liveliness.NewWLP(protocol)
		protocol.WLP.InitWL(protocol.participantImpl)
	}

	// TypeLookupManager
	if protocol.Att.TypeLookupConfig.UseClient || protocol.Att.TypeLookupConfig.UseServer {
		log.Fatalln("not impl")
	}

	writeParam := common.KWriteParamDefault
	protocol.PDP.AnnounceParticipantState(true, false, &writeParam)
	protocol.PDP.ResetParticipantAnnouncement()
	protocol.PDP.Enable()

	return true
}

func NewBuiltinProtocols() *Protocols {
	var buildinProtocols Protocols
	return &buildinProtocols
}
