package rtps

import (
	. "common"
	. "fastrtps/rtps/attributes"
	. "fastrtps/rtps/participant"
	. "fastrtps/rtps/resources"
	. "fastrtps/utils"
	"os"
)

var (
	maxRTPSParticipantID uint32 = 1
	mRTPSParticipantIDs  map[uint32]bool
)

type RTPSDomain struct {
	Participant_list []*Participant
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

func (domain *RTPSDomain) CreateParticipant(attrs *RTPSParticipantAttributes, listen *ParticipantListener) *RTPSParticipant {

	participant := NewRTPSParticipant(attrs, listen)
	domain.Participant_list = append(Participant_list, participant)

	return participant
}
