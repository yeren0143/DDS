package rtps

import (
	. "attributes"
	. "common"
	"log"
	"os"
	. "utils"
)

var (
	maxRTPSParticipantID uint32 = 1
	mRTPSParticipantIDs  map[uint32]bool
)

type RTPSDomain struct {
}

// type t_p_RTPSParticipant {
// 	var participant *RTPSParticipant
// 	var participant_impl *RTPSParticipantImpl
// }

func createGuidPrefix(ID uint32) *GuidPrefix_t {
	var guid GuidPrefix_t
	pid := os.Getppid()

	guid.value[0] = C_VendorId_eProsima.Vendor[0]
	guid.value[1] = C_VendorId_eProsima.Vendor[1]

	host_id := GetHost().Id()
	guid.value[2] = Octet(host_id)
	guid.value[3] = Octet(host_id >> 8)

	guid.value[4] = Octet(pid)
	guid.value[5] = Octet(pid >> 8)
	guid.value[6] = Octet(pid >> 16)
	guid.value[7] = Octet(pid >> 24)
	guid.value[8] = Octet(ID)
	guid.value[9] = Octet(ID >> 8)
	guid.value[10] = Octet(ID >> 16)
	guid.value[11] = Octet(ID >> 24)

	return &guid
}

func CreateParticipant(domain_id uint32, attrs *RTPSParticipantAttributes) *RTPSParticipant {
	log.Println("")
}
