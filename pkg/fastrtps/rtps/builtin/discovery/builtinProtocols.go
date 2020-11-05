package discovery

import (
	. "attributes"
	. "common"
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
