package rtps

import (
	. "attributes"
	. "common"
	"log"
	"os"
	. "participant"
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

	guid.Value[0] = C_VendorId_eProsima.Vendor[0]
	guid.Value[1] = C_VendorId_eProsima.Vendor[1]

	host_id := GetHost().Id()
	guid.Value[2] = Octet(host_id)
	guid.Value[3] = Octet(host_id >> 8)

	guid.Value[4] = Octet(pid)
	guid.Value[5] = Octet(pid >> 8)
	guid.Value[6] = Octet(pid >> 16)
	guid.Value[7] = Octet(pid >> 24)
	guid.Value[8] = Octet(ID)
	guid.Value[9] = Octet(ID >> 8)
	guid.Value[10] = Octet(ID >> 16)
	guid.Value[11] = Octet(ID >> 24)

	return &guid
}

func NewParticipant(domain_id uint32, attrs *RTPSParticipantAttributes) *RTPSParticipant {
	log.Println("")
}
