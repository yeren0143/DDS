package rtps

import (
	"log"
	"os"
	"sync"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/participant"
	"github.com/yeren0143/DDS/fastrtps/utils"
)

var (
	gMaxRTPSParticipantID uint32          = 1
	gRTPSParticipantIDs   map[uint32]bool = make(map[uint32]bool)
	gLock                 sync.Mutex
	gParticipantList      []*participant.RTPSParticipant
)

func getNewID() uint32 {
	ret := gMaxRTPSParticipantID
	gMaxRTPSParticipantID++
	return ret
}

func createGUIDPrefix(ID uint32) *common.GUIDPrefixT {
	var guid common.GUIDPrefixT
	pid := os.Getppid()

	guid.Value[0] = common.KVendorIDTeProsima.Vendor[0]
	guid.Value[1] = common.KVendorIDTeProsima.Vendor[1]

	hostID := utils.GetHost().Id()
	guid.Value[2] = common.Octet(hostID)
	guid.Value[3] = common.Octet(hostID >> 8)

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

func ClientServerEnvironmentCreationOverride(domainID uint32, enable bool, attrs *attributes.RTPSParticipantAttributes,
	listen *participant.RTPSParticipantListener) *participant.RTPSParticipant {
	data := os.Getenv("DEFAULT_ROS2_MASTER_URI")
	if len(data) == 0 {
		return nil
	}

	log.Fatalln("not Impl")
	return nil
}

func createParticipant(domainID uint32, useProtocol bool, att *attributes.RTPSParticipantAttributes,
	listen *participant.RTPSParticipantListener) *participant.RTPSParticipant {

	if att.Builtin.DiscoveryConfig.LeaseDuration.Less(common.KTimeInfinite) &&
		!att.Builtin.DiscoveryConfig.LeaseDurationAnnouncementPeriod.Less(att.Builtin.DiscoveryConfig.LeaseDuration) {

		log.Fatal("RTPSParticipant Attributes: LeaseDuration should be >= leaseDuration announcement period")
		return nil
	}

	var id uint32
	gLock.Lock()
	{
		if att.ParticipantID < 0 {
			id = getNewID()
			_, found := gRTPSParticipantIDs[id]
			for found {
				id = getNewID()
				_, found = gRTPSParticipantIDs[id]
			}
			gRTPSParticipantIDs[id] = true
		} else {
			id = uint32(att.ParticipantID)
			if _, found := gRTPSParticipantIDs[id]; found {
				log.Fatal("RTPSParticipant with the same ID already exists")
				gLock.Unlock()
				return nil
			}
		}
	}
	gLock.Unlock()

	if att.DefaultUnicastLocatorList.Valid() == false {
		log.Fatal("Default Unicast Locator List contains invalid Locator")
		return nil
	}

	if att.DefaultMulticastLocatorList.Valid() == false {
		log.Fatal("Default Multicast Locator List contains invalid Locator")
		return nil
	}

	att.ParticipantID = int32(id)
	loc := utils.GetIP4Address()

	if loc.Empty() && att.Builtin.InitialPeersList.Empty() {
		var local common.Locator
		utils.SetIPv4WithBytes(&local, []common.Octet{127, 0, 0, 1})
		att.Builtin.InitialPeersList.PushBack(&local)
	}

	guidP := createGUIDPrefix(id)

	// If we force the participant to have a specific prefix we must define a different persistence GuidPrefix_t that
	// would ensure builtin endpoints are able to differentiate between a communication loss and a participant recovery
	var rtpsParticipant *participant.RTPSParticipant
	if *att.Prefix != common.KUnknownGUIDPrefix {
		rtpsParticipant = participant.NewParticipant(domainID, att, att.Prefix, guidP, listen)
	} else {
		rtpsParticipant = participant.NewParticipant(domainID, att, guidP, &common.KUnknownGUIDPrefix, listen)
	}

	if rtpsParticipant == nil {
		log.Fatal("create rtpsParticipant failed")
	}

	// Above constructors create the sender resources. If a given listening port cannot be allocated an iterative
	// mechanism will allocate another by default. Change the default listening port is unacceptable for server
	// discovery.
	if (att.Builtin.DiscoveryConfig.Protocol == attributes.KDisPServer ||
		att.Builtin.DiscoveryConfig.Protocol == attributes.KDisPBackup) &&
		rtpsParticipant.DidMutationTookPlaceOnMeta(att.Builtin.MetatrafficMulticastLocatorList,
			att.Builtin.MetatrafficUnicastLocatorList) {
		log.Fatal("Server wasn't able to allocate the specified listening port")
		return nil
	}

	// Check there is at least one transport registered.
	if rtpsParticipant.NetworkFactoryHasRegisteredTransports() == false {
		log.Fatal("Cannot create rtpsParticipant, because there is any transport")
		return nil
	}

	gLock.Lock()
	gParticipantList = append(gParticipantList, rtpsParticipant)
	gLock.Unlock()

	if useProtocol {
		rtpsParticipant.Enable()
	}

	return rtpsParticipant
}

//NewRTPSParticipant create new rtps participant
func NewRTPSParticipant(att *attributes.ParticipantAttributes,
	listen *participant.RTPSParticipantListener) *participant.RTPSParticipant {

	log.Println("cretae RTPS participant")

	// If DEFAULT_ROS2_MASTER_URI is specified then try to create default client if
	// that already exists.
	part := ClientServerEnvironmentCreationOverride(att.DomainID, false, att.RTPS, listen)

	if part == nil {
		part = createParticipant(att.DomainID, false, att.RTPS, listen)
	}

	if part == nil {
		log.Fatalln("Problem creating RTPSParticipant")
		return nil
	}

	// Enable participant
	part.Enable()

	return part
}
