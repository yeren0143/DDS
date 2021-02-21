package participant

import (
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/protocol"
	"github.com/yeren0143/DDS/fastrtps/rtps/reader"
)

var _ reader.IReaderListener = (*PDPListener)(nil)

/**
 * Class PDPListener, specification used by the PDP to perform the History check when a new message is received.
 * This class is implemented in order to use the same structure than with any other RTPSReader.
 * @ingroup DISCOVERY_MODULE
 */
type PDPListener struct {
	reader.ReaderListenerBase
}

func newPDPListener(parent protocol.IPDP) *PDPListener {
	return &PDPListener{}
}
