package discovery

import (
	. "github.com/yeren0143/DDS/common"
	. "github.com/yeren0143/DDS/fastrtps/rtps/attributes"
)

type BuiltinProtocols struct {
	Att                             BuiltinAttributes
	ParticipantImpl                 interface{}
	PDP                             interface{}
	WLP                             interface{}
	TLM                             interface{} //TypeLookupManager
	MetatrafficMulticastLocatorList LocatorList
	MetatrafficUnicastLocatorList   LocatorList
	InitialPeersList                LocatorList
	DiscoveryServers                []interface{} //Known discovery and backup server container
}

func CreateBuiltinProtocols() BuiltinProtocols {
	var buildProtocols BuiltinProtocols
	return buildProtocols
}
