package endpoint

import (
	"sync"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/core/policy"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/history"
	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
)

type IWlp interface {
	InitWL(p interface{}) bool
	AddWriter(guid *common.GUIDT, kind policy.LivelinessQosPolicyKind, leaseDuration *common.DurationT) bool
	RemoveWriter(guid *common.GUIDT, kind policy.LivelinessQosPolicyKind, leaseDuration *common.DurationT) bool
	AssertLiveliness(writerGuid *common.GUIDT, kind policy.LivelinessQosPolicyKind, leaseDuration *common.DurationT) bool
}

type IEndpointParent interface {
	Wlp() IWlp
	GetAttributes() *attributes.RTPSParticipantAttributes
	CreateSenderResources(locator *common.Locator)
	GetMinNetworkSendBufferSize() uint32
	GetEventResource() *resources.ResourceEvent
	SendSync(msg *common.CDRMessage, locators []common.Locator, maxBlockingTimePoint common.Time) bool
}

type IEndpoint interface {
	GetGUID() *common.GUIDT
	GetMutex() *sync.Mutex
	GetAttributes() *attributes.EndpointAttributes
	GetRtpsParticipant() IEndpointParent
}

/**
 * Class Endpoint, all entities of the RTPS network derive from this class.
 * Although the RTPSParticipant is also defined as an endpoint in the RTPS specification, in this implementation
 * the RTPSParticipant class **does not** inherit from the endpoint class. Each Endpoint object owns a pointer to the
 * RTPSParticipant it belongs to.
 * @ingroup COMMON_MODULE
 */
type EndpointBase struct {
	Mutex            sync.Mutex
	GUID             common.GUIDT
	Att              attributes.EndpointAttributes
	PayloadPool      history.IPayloadPool
	ChangePool       history.IChangePool
	FixedPayloadSize uint32
	RTPSParticipant  IEndpointParent
}

func (endpointBase *EndpointBase) GetAttributes() *attributes.EndpointAttributes {
	return &endpointBase.Att
}

func (endpointBase *EndpointBase) GetGUID() *common.GUIDT {
	return &endpointBase.GUID
}

func (endpointBase *EndpointBase) GetMutex() *sync.Mutex {
	return &endpointBase.Mutex
}

func (endpointBase *EndpointBase) GetRtpsParticipant() IEndpointParent {
	return endpointBase.RTPSParticipant
}

func NewEndPointBase(parent IEndpointParent, guid *common.GUIDT, att *attributes.EndpointAttributes) *EndpointBase {
	return &EndpointBase{
		RTPSParticipant: parent,
		Att:             *att,
		GUID:            *guid,
	}
}
