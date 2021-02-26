package builtin

import (
	"log"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/pdp"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/protocol"
	"github.com/yeren0143/DDS/fastrtps/rtps/endpoint"
	"github.com/yeren0143/DDS/fastrtps/rtps/network"
)

//var _ participant.IPDPParent = (*Protocols)(nil)

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

func (protocol *Protocols) UpdateMetatrafficLocators(loclist *common.LocatorList) bool {
	*protocol.MetatrafficUnicastLocatorList = *loclist
	return true
}

func (protocol *Protocols) GetMetatrafficMulticastLocatorList() *common.LocatorList {
	return protocol.MetatrafficMulticastLocatorList
}

func (protocol *Protocols) AnnounceParticipantState() {
	if protocol.PDP != nil {
		protocol.PDP.AnnounceParticipantState(false, false, &common.KWriteParamDefault)
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
	log.Fatalln("WLP not impl")

	return true
}

func NewBuiltinProtocols() *Protocols {
	var buildinProtocols Protocols
	return &buildinProtocols
}
