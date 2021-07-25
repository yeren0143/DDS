package transport

import (
	common "github.com/yeren0143/DDS/common"
)

type ISenderResource interface {
	Send(data []common.Octet, locators common.LocatorList, maxBlockingTimePoint common.Time) bool
}

type cleanFunc func()
type sendFunc func([]common.Octet, common.LocatorList, common.Time) bool

// SenderResource is RAII object that encapsulates the Send operation over one chanel in an unknown transport.
// A Sender resource is always univocally associated to a transport channel; the
// act of constructing a Sender Resource opens the channel and its destruction
// closes it.
var _ ISenderResource = (*SenderResource)(nil)

type SenderResource struct {
	transportKind int32
	cleanUp       *cleanFunc
	sendImpl      *sendFunc
}

type SenderResourceList = []ISenderResource

// Send to a destination locator, through the channel managed by this resource.
func (sendResource *SenderResource) Send(data []common.Octet, locators common.LocatorList, maxBlockingTimePoint common.Time) bool {
	ret := false
	if sendResource.sendImpl != nil {
		ret = (*sendResource.sendImpl)(data, locators, maxBlockingTimePoint)
	}
	return ret
}

// Kind return transport kind
func (sendResource *SenderResource) Kind() int32 {
	return sendResource.transportKind
}
