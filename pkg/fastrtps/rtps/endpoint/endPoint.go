package endpoint

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/history"
	"sync"
)

/**
 * Class Endpoint, all entities of the RTPS network derive from this class.
 * Although the RTPSParticipant is also defined as an endpoint in the RTPS specification, in this implementation
 * the RTPSParticipant class **does not** inherit from the endpoint class. Each Endpoint object owns a pointer to the
 * RTPSParticipant it belongs to.
 * @ingroup COMMON_MODULE
 */
type Endpoint struct {
	Mutex            sync.Mutex
	guid             common.GUIDT
	att              attributes.EndpointAttributes
	payloadPool      history.IPayloadPool
	ChangePool       history.IChangePool
	fixedPayloadSize uint32
}
