package transport

import (
	common "github.com/yeren0143/DDS/common"
)

//ITransportReceiver against which to implement a data receiver, decoupled from transport internals.
type ITransportReceiver interface {
	OnDataReceived(data []common.Octet, local *common.Locator, remote *common.Locator)
}
