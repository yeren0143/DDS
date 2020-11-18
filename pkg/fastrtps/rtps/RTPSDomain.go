package rtps

import (
	"log"
	"os"
	"sync"

	common "github.com/yeren0143/DDS/common"
	rtpsAtt "github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	rtpsPant "github.com/yeren0143/DDS/fastrtps/rtps/participant"
	utils "github.com/yeren0143/DDS/fastrtps/utils"
)

var (
	gMaxRTPSParticipantID uint32 = 1
	gRTPSParticipantIDs   map[uint32]bool
	gLock                 sync.Mutex
)

type RTPSDomain struct {
	ParticipantList []*rtpsPant.RTPSParticipant
}

func (domain *RTPSDomain) getNewId() uint32 {
	ret := gMaxRTPSParticipantID
	gMaxRTPSParticipantID += 1
	return ret
}

func createGuidPrefix(ID uint32) *common.GuidPrefix_t {
	var guid common.GuidPrefix_t
	pid := os.Getppid()

	guid.Value[0] = common.CVendorIDTeProsima.Vendor[0]
	guid.Value[1] = common.CVendorIDTeProsima.Vendor[1]

	host_id := utils.GetHost().Id()
	guid.Value[2] = common.Octet(host_id)
	guid.Value[3] = common.Octet(host_id >> 8)

	guid.Value[4] = common.Octet(pid)
	guid.Value[5] = common.Octet(pid >> 8)
	guid.Value[6] = common.Octet(pid >> 16)
	guid.Value[7] = common.Octet(pid >> 24)
	guid.Value[8] = common.Octet(ID)
	guid.Value[9] = common.Octet(ID >> 8)
	guid.Value[10] = common.Octet(ID >> 16)
	guid.Value[11] = common.Octet(ID >> 24)

	return &guid
}

func (domain *RTPSDomain) NewRTPSParticipant(domainId uint32, useProtocol bool, attrs *rtpsAtt.RTPSParticipantAttributes, listen *rtpsPant.RTPSParticipantListener) *rtpsPant.RTPSParticipant {
	log.Println("cretae RTPS participant")

	if attrs.Builtin.DiscoveryConfig.LeaseDuration.Less(common.CTimeInfinite) &&
		!attrs.Builtin.DiscoveryConfig.LeaseDurationAnnouncementPeriod.Less(attrs.Builtin.DiscoveryConfig.LeaseDuration) {

		log.Fatal("RTPSParticipant Attributes: LeaseDuration should be >= leaseDuration announcement period")
		return nil
	}

	var id uint32
	{
		gLock.Lock()
		defer gLock.Unlock()

		if attrs.ParticipantID < 0 {
			id = domain.getNewId()
			_, ok := gRTPSParticipantIDs[id]
			for ok {
				id = domain.getNewId()
				_, ok = gRTPSParticipantIDs[id]
			}
			gRTPSParticipantIDs[id] = true
		} else {
			id = attrs.ParticipantID
			if _, ok := gRTPSParticipantIDs[id]; ok {
				log.Fatal("RTPSParticipant with the same ID already exists")
				return nil
			}
		}
	}

	if attrs.DefaultUnicastLocatorList.Valid() == false {
		log.Fatal("Default Unicast Locator List contains invalid Locator")
		return nil
	}

	if attrs.DefaultMulticastLocatorList.Valid() == false {
		log.Fatal("Default Multicast Locator List contains invalid Locator")
		return nil
	}

	attrs.ParticipantID = id
	loc := utils.GetIP4Address()

	if loc.Empty() && attrs.Builtin.InitialPeersList.Empty() {
		var local common.Locator
		utils.SetIPv4(&local, 127, 0, 0, 1)
		attrs.Builtin.InitialPeersList.PushBack(&local)
	}

	participant := rtpsPant.NewParticipant(domainId, useProtocol, attrs, listen)
	domain.ParticipantList = append(domain.ParticipantList, participant)

	return participant
}
