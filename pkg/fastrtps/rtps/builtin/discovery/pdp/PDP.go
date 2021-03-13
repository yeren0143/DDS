package pdp

import (
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/data"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/edp"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/protocol"
	"github.com/yeren0143/DDS/fastrtps/rtps/history"
	"github.com/yeren0143/DDS/fastrtps/rtps/reader"
	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
	"github.com/yeren0143/DDS/fastrtps/rtps/writer"
	"github.com/yeren0143/DDS/fastrtps/utils"
)

type IPDPParent interface {
	UpdateMetatrafficLocators(loclist *common.LocatorList) bool
	GetBuiltinAttributes() *attributes.BuiltinAttributes
	GetMetatrafficMulticastLocatorList() *common.LocatorList
	GetMetatrafficUnicastLocatorList() *common.LocatorList
	GetInitialPeers() *common.LocatorList
}

type IpdpBaseImpl interface {
	CreatePDPEndpoints() bool
}

type ParticipantProxyDataVector struct {
	Proxies []*data.ParticipantProxyData
	Config  utils.ResourceLimitedContainerConfig
}

type ReaderProxyDataVector struct {
	Proxies []*data.ReaderProxyData
	Config  utils.ResourceLimitedContainerConfig
}

type WriterProxyDataVector struct {
	Proxies []*data.WriterProxyData
	Config  utils.ResourceLimitedContainerConfig
}

type pdpBase struct {
	impl            IpdpBaseImpl
	builtin         IPDPParent
	rtpsParticipant protocol.IParticipant
	discovery       *attributes.BuiltinAttributes
	writer          writer.IRTPSWriter
	reader          reader.IRTPSReader
	EDP             edp.IEDP
	// Number of participant proxy data objects created
	participantProxiesNumber uint32
	// Registered RTPSParticipants (including the local one, that is the first one.)
	participantProxies ParticipantProxyDataVector
	// Pool of participant proxy data objects ready for reuse
	participantProxiesPool ParticipantProxyDataVector
	// Number of reader proxy data objects created
	readerProxiesNumber uint32
	// Pool of reader proxy data objects ready for reuse
	readerProxiesPool ReaderProxyDataVector
	// Number of writer proxy data objects created
	writerProxiesNumber uint32
	// Pool of writer proxy data objects ready for reuse
	writerProxiesPool WriterProxyDataVector
	// Variable to indicate if any parameter has changed.
	hasChangedLocalPDP int32
	// Listener for the SPDP messages.
	listener          reader.IReaderListener
	pdpWriterHistory  *history.WriterHistory
	pdpReaderHistory  *history.ReaderHistory
	writerPayloadPool history.ITopicPayloadPool
	readerPayloadPool history.ITopicPayloadPool
	// ReaderProxyData to allow preallocation of remote locators
	tempReaderData *data.ReaderProxyData
	// WriterProxyData to allow preallocation of remote locators
	tempWriterData *data.WriterProxyData
	tempDataMutex  sync.Mutex
	// Participant data atomic access assurance
	mutex *sync.Mutex
	// To protect callbacks (ParticipantProxyData&)
	callbackMutex sync.Mutex

	// TimedEvent to periodically resend the local RTPSParticipant information.
	resendParticipantInfoEvent *resources.TimedEvent
	initialAnnouncements       attributes.InitialAnnouncementConfig
}

func (pdp *pdpBase) BuiltinAttributes() *attributes.BuiltinAttributes {
	return pdp.builtin.GetBuiltinAttributes()
}

func (pdp *pdpBase) GetRTPSParticipant() protocol.IParticipant {
	return pdp.rtpsParticipant
}

// Force the sending of our local DPD to all remote RTPSParticipants and multicast Locators.
func (pdp *pdpBase) AnnounceParticipantState(newChange, dispose bool, wparams *common.WriteParamsT) {
	var aChange *common.CacheChangeT
	if !dispose {
		if atomic.CompareAndSwapInt32(&pdp.hasChangedLocalPDP, 1, -1) || newChange {
			pdp.mutex.Lock()
			localParticipantData := pdp.GetLocalParticipantProxyData()
			key := localParticipantData.Key
			proxyDataCopy := *localParticipantData
			pdp.mutex.Unlock()

			if pdp.pdpWriterHistory.GetHistorySize() > 0 {
				pdp.pdpWriterHistory.RemoveMinChange()
			}
			cdrSize := proxyDataCopy.GetSerializedSize(true)
			cdrCallback := func() uint32 {
				return cdrSize
			}
			aChange = pdp.writer.NewChange(cdrCallback, common.KAlive, key)
			if aChange != nil {
				auxMsg := common.NewCDRMessageWithPayload(&aChange.SerializedPayload)
				if os.Getenv("BIG_ENDIAN") != "" {
					aChange.SerializedPayload.Encapsulation = common.PL_CDR_BE
					auxMsg.MsgEndian = common.BIGEND
				} else {
					aChange.SerializedPayload.Encapsulation = common.PL_CDR_LE
					auxMsg.MsgEndian = common.LITTLEEND
				}

				if proxyDataCopy.WriteToCDRMessage(auxMsg, true) {
					aChange.SerializedPayload.Length = auxMsg.Length
					pdp.pdpWriterHistory.AddChange(aChange, wparams)
				} else {
					log.Panic("Cannot serialize ParticipantProxyData.")
				}
			}
		} else {
			pdp.mutex.Lock()
			proxyDataCopy := *pdp.GetLocalParticipantProxyData()
			pdp.mutex.Unlock()

			if pdp.pdpWriterHistory.GetHistorySize() > 0 {
				pdp.pdpWriterHistory.RemoveMinChange()
			}
			cdrSize := proxyDataCopy.GetSerializedSize(true)
			cdrCallback := func() uint32 {
				return cdrSize
			}
			aChange := pdp.writer.NewChange(cdrCallback, common.KNotAliveDisposedUnregistered, proxyDataCopy.Key)
			if aChange != nil {
				auxMsg := common.NewCDRMessageWithPayload(&aChange.SerializedPayload)
				if os.Getenv("BIG_ENDIAN") != "" {
					aChange.SerializedPayload.Encapsulation = common.PL_CDR_BE
					auxMsg.MsgEndian = common.BIGEND
				} else {
					aChange.SerializedPayload.Encapsulation = common.PL_CDR_LE
					auxMsg.MsgEndian = common.LITTLEEND
				}
				if proxyDataCopy.WriteToCDRMessage(auxMsg, true) {
					aChange.SerializedPayload.Length = auxMsg.Length
					pdp.pdpWriterHistory.AddChange(aChange, wparams)
				} else {
					log.Panic("Cannot serialize ParticipantProxyData.")
				}
			}
		}
	}
}

func (pdp *pdpBase) GetLocalParticipantProxyData() *data.ParticipantProxyData {
	return pdp.participantProxies.Proxies[0]
}

func (pdp *pdpBase) initPDP(participant protocol.IParticipant) bool {
	log.Println("Beginning")
	pdp.rtpsParticipant = participant
	pdp.discovery = participant.GetAttributes().Builtin
	pdp.initialAnnouncements = *pdp.discovery.DiscoveryConfig.InitialAnnouncements
	//CREATE ENDPOINTS
	if !pdp.impl.CreatePDPEndpoints() {
		return false
	}
	//UPDATE METATRAFFIC.
	pdp.builtin.UpdateMetatrafficLocators(&pdp.reader.GetAttributes().UnicastLocatorList)

	pdp.mutex.Lock()
	pdata := pdp.addParticipantProxyData(participant.GetGuid(), false)
	pdp.mutex.Unlock()

	if pdata == nil {
		return false
	}
	pdp.initializeParticipantProxyData(pdata)

	// Create lease events on already created proxy data objects
	for _, item := range pdp.participantProxiesPool.Proxies {
		callback := func() bool {
			pdp.checkRemoteParticipantLiveliness(item)
			return false
		}
		item.LeaseDurationEvent = resources.NewTimedEvent(pdp.rtpsParticipant.GetEventResource(), &callback, 0)
	}

	callback := func() bool {
		pdp.AnnounceParticipantState(false, false, &common.KWriteParamDefault)
		pdp.setNextAnnouncementInterval()
		return true
	}
	pdp.resendParticipantInfoEvent = resources.NewTimedEvent(pdp.rtpsParticipant.GetEventResource(), &callback, 0)

	pdp.setInitialAnnouncementInterval()
	return true
}

func (pdp *pdpBase) setInitialAnnouncementInterval() {
	if pdp.initialAnnouncements.Count > 0 && pdp.initialAnnouncements.Period.Less(common.KTimeZero) {
		// Force a small interval (1ms) between initial announcements
		log.Println("Initial announcement period is not strictly positive. Changing to 1ms.")
		pdp.initialAnnouncements.Period = common.DurationT{
			Seconds: 0,
			Nanosec: 1000000,
		}
	}
	pdp.setNextAnnouncementInterval()
}

func (pdp *pdpBase) setNextAnnouncementInterval() {
	if pdp.initialAnnouncements.Count > 0 {
		pdp.initialAnnouncements.Count--
		pdp.resendParticipantInfoEvent.UpdateInterval(pdp.initialAnnouncements.Period)
	} else {
		inter := pdp.discovery.DiscoveryConfig.LeaseDurationAnnouncementPeriod
		pdp.resendParticipantInfoEvent.UpdateInterval(inter)
	}
}

func (pdp *pdpBase) removeRemoteParticipant(partGUID *common.GUIDT, reason DiscoveryStatus) bool {
	log.Fatalln("not impl")
	return false
}

func (pdp *pdpBase) checkRemoteParticipantLiveliness(remoteParticipant *data.ParticipantProxyData) {
	pdp.mutex.Lock()
	defer pdp.mutex.Unlock()
	if remoteParticipant.ShouldCheckLeaseDuration {
		now := time.Now()
		realLeaseTm := remoteParticipant.LastReceivedMessageTm.Add(remoteParticipant.LeaseDuration)
		if now.After(realLeaseTm) {
			pdp.mutex.Unlock()
			pdp.removeRemoteParticipant(&remoteParticipant.Guid, KDroppedParticipant)
			return
		}
	}
}

func (pdp *pdpBase) AddReaderProxyData(readerGUID, participantGUID *common.GUIDT,
	initializer protocol.ReaderProxyDataInitFunc) *data.ReaderProxyData {
	log.Println("Adding reader proxy data ", *readerGUID)
	log.Fatalln("not impl")
	return nil
}

func (pdp *pdpBase) addParticipantProxyData(participantGUID *common.GUIDT, withLeaseDuration bool) *data.ParticipantProxyData {
	var retVal *data.ParticipantProxyData

	// Try to take one entry from the pool
	if len(pdp.participantProxiesPool.Proxies) == 0 {
		maxProxies := pdp.participantProxies.Config.Maximum
		if pdp.participantProxiesNumber < maxProxies {
			// Pool is empty but limit has not been reached, so we create a new entry.
			pdp.participantProxiesNumber++
			retVal = data.NewParticipantProxyData(pdp.rtpsParticipant.GetAttributes().Allocation)
			if participantGUID != pdp.rtpsParticipant.GetGuid() {
				eventCallBack := func() bool {
					pdp.checkRemoteParticipantLiveliness(retVal)
					return false
				}
				pevent := pdp.rtpsParticipant.GetEventResource()
				retVal.LeaseDurationEvent = resources.NewTimedEvent(pevent, &eventCallBack, 0)
			}
		} else {
			log.Fatalln("Maximum number of participant proxies (", maxProxies, ") reached for participant ",
				pdp.rtpsParticipant.GetGuid())
			return nil
		}
	} else {
		// Pool is not empty, use entry from pool
		length := len(pdp.participantProxiesPool.Proxies)
		retVal = pdp.participantProxiesPool.Proxies[length-1]
		pdp.participantProxiesPool.Proxies = pdp.participantProxiesPool.Proxies[:length-1]
	}

	// Add returned entry to the collection
	retVal.ShouldCheckLeaseDuration = withLeaseDuration
	retVal.Guid = *participantGUID
	pdp.participantProxies.Proxies = append(pdp.participantProxies.Proxies, retVal)

	return retVal
}

func (pdp *pdpBase) initializeParticipantProxyData(*data.ParticipantProxyData) {

}

func newPDP(protocol IPDPParent, att *attributes.RTPSParticipantAllocationAttributes, impl IpdpBaseImpl) *pdpBase {
	var pdp pdpBase
	pdp.mutex = new(sync.Mutex)
	pdp.builtin = protocol
	// pdp.participantProxies = att.Participants
	pdp.participantProxiesNumber = att.Participants.Initial
	pdp.participantProxies.Config = *att.Participants
	pdp.participantProxiesPool.Config = *att.Participants
	pdp.readerProxiesPool.Config = *att.TotalReaders()
	pdp.readerProxiesNumber = pdp.readerProxiesPool.Config.Initial
	pdp.writerProxiesPool.Config = *att.TotalWriters()
	pdp.writerProxiesNumber = pdp.writerProxiesPool.Config.Initial
	pdp.hasChangedLocalPDP = 1

	maxUnicastLocators := uint32(att.Locators.MaxUnicastLocators)
	maxMulticastLocators := uint32(att.Locators.MaxMulticastLocators)

	pdp.tempReaderData = data.NewReaderProxyData(maxUnicastLocators, maxMulticastLocators, att.DataLimits)
	pdp.tempWriterData = data.NewWriterProxyData(maxUnicastLocators, maxMulticastLocators, att.DataLimits)

	for i := uint32(0); i < att.Participants.Initial; i++ {
		proxyData := data.NewParticipantProxyData(att)
		pdp.participantProxiesPool.Proxies = append(pdp.participantProxiesPool.Proxies, proxyData)
	}

	for i := uint32(0); i < att.TotalReaders().Initial; i++ {
		proxyData := data.NewReaderProxyData(maxUnicastLocators, maxMulticastLocators, att.DataLimits)
		pdp.readerProxiesPool.Proxies = append(pdp.readerProxiesPool.Proxies, proxyData)
	}

	for i := uint32(0); i < att.TotalWriters().Initial; i++ {
		proxyData := data.NewWriterProxyData(maxUnicastLocators, maxMulticastLocators, att.DataLimits)
		pdp.writerProxiesPool.Proxies = append(pdp.writerProxiesPool.Proxies, proxyData)
	}

	pdp.impl = impl

	return &pdp
}
