package attributes

import (
	"github.com/yeren0143/DDS/common"
)

// Structure EndpointAttributes, describing the attributes associated with an RTPS Endpoint.
type EndpointAttributes struct {
	EndpointKind         common.EndpointKindT
	TopicKind            common.TopicKindT
	ReliabilityKind      common.ReliabilityKindT
	DurabilityKind       common.DurabilityKindT
	PersistenceGUID      common.GUIDT
	Properties           PropertyPolicyT
	UnicastLocatorList   common.LocatorList
	MulticastLocatorList common.LocatorList
	RemoteLocatorList    common.LocatorList

	// User Defined ID, used for StaticEndpointDiscovery, default value -1.
	userDefinedID int16

	// Entity ID, if the user want to specify the EntityID of the enpoint, default value -1.
	entityID int16
}

func (att *EndpointAttributes) GetEntityID() int16 {
	return att.entityID
}

var KDefaultEndpointAttributes = EndpointAttributes{
	EndpointKind:    common.KWriter,
	TopicKind:       common.KNoKey,
	ReliabilityKind: common.KBestEffort,
	DurabilityKind:  common.KVolatile,
	userDefinedID:   -1,
	entityID:        -1,
}
