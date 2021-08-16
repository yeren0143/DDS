package transport

import (
	common "dds/common"
)

//ITransportReceiver against which to implement a data receiver, decoupled from transport internals.
type ITransportReceiver interface {
	OnDataReceived(data []common.Octet, length uint32, local *common.Locator, remote *common.Locator)
}
