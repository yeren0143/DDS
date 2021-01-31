package transport

import (
	common "github.com/yeren0143/DDS/common"
)

//ITransportReceiver against which to implement a data receiver, decoupled from transport internals.
type ITransportReceiver interface {
	OnDataReceived(data []common.Octet, length uint32, local *common.Locator, remote *common.Locator)
}
