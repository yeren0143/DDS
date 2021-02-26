package writer

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/endpoint"
	"github.com/yeren0143/DDS/fastrtps/rtps/reader"
)

var _ IRTPSMessageSender = (*ReaderLocator)(nil)

// type readerLocatorParent interface {
// 	SendSync(msg *common.CDRMessage, locators []common.Locator, maxBlockingTimePoint common.Time) bool
// }

/**
 * Class ReaderLocator, contains information about a remote reader, without saving its state.
 * It also implements IRTPSMessageSender, so it can be used when separate sending is enabled.
 */
type ReaderLocator struct {
	owner            IRTPSWriter
	participantOwner endpoint.IEndpointParent
	expectsInlineQos bool
	isLocalReader    bool
	localReader      reader.IRTPSReader
	localInfo        common.LocatorSelectorEntry
	guidPrefixs      []common.GUIDPrefixT
	guids            []common.GUIDT
}

func (readerLocator *ReaderLocator) DestinationHaveChanged() bool {
	return false
}

func (readerLocator *ReaderLocator) DestinationGuidPrefix() common.GUIDPrefixT {
	return readerLocator.localInfo.RemoteGUID.Prefix
}

func (readerLocator *ReaderLocator) RemoteParticipants() []common.GUIDPrefixT {
	return readerLocator.guidPrefixs
}

func (readerLocator *ReaderLocator) RemoteGUIDs() []common.GUIDT {
	return readerLocator.guids
}

func (readerLocator *ReaderLocator) Send(msg *common.CDRMessage, maxBlockingTimePoint common.Time) bool {
	if readerLocator.localInfo.RemoteGUID != common.KGuidUnknown && !readerLocator.isLocalReader {
		if len(readerLocator.localInfo.Unicast) > 0 {
			return readerLocator.participantOwner.SendSync(msg, readerLocator.localInfo.Unicast, maxBlockingTimePoint)
		} else {
			return readerLocator.participantOwner.SendSync(msg, readerLocator.localInfo.Multicast, maxBlockingTimePoint)
		}
	}
	return true
}

func NewReaderLocator(owner IRTPSWriter, maxUnicastLocatos, maxMulticastLocators uint32) *ReaderLocator {
	var readerLocator ReaderLocator
	readerLocator.owner = owner
	readerLocator.participantOwner = owner.GetRtpsParticipant()
	readerLocator.localInfo = *common.NewLocatorSelectorEntry(uint32(maxUnicastLocatos), uint32(maxMulticastLocators))
	return &readerLocator
}
