package network

import (
	common "github.com/yeren0143/DDS/common"
	"time"
)

type cleanFunc func()
type sendFunc func([]common.Octet, uint32, common.LocatorList, time.Time) bool

// SenderResource is RAII object that encapsulates the Send operation over one chanel in an unknown transport.
// A Sender resource is always univocally associated to a transport channel; the
// act of constructing a Sender Resource opens the channel and its destruction
// closes it.
type SenderResource struct {
	transportKind int32
	cleanUp       *cleanFunc
	sendImpl      *sendFunc
}

// Send to a destination locator, through the channel managed by this resource.
func (sendResource *SenderResource) Send(data []common.Octet, length uint32, locators common.LocatorList, maxBlockingTimePoint time.Time) bool {
	ret := false
	if sendResource.sendImpl != nil {
		ret = (*sendResource.sendImpl)(data, length, locators, maxBlockingTimePoint)
	}
	return ret
}

// Kind return transport kind
func (sendResource *SenderResource) Kind() int32 {
	return sendResource.transportKind
}