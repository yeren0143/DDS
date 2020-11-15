package rtps

import (
	. "github.com/yeren0143/DDS/common"
	. "github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	. "github.com/yeren0143/DDS/fastrtps/rtps/participant"
	// . "github.com/yeren0143/DDS/fastrtps/rtps/resources"
	. "github.com/yeren0143/DDS/fastrtps/utils"
	"os"
)

var (
	maxRTPSParticipantID uint32 = 1
	mRTPSParticipantIDs  map[uint32]bool
)

type RTPSDomain struct {
	Participant_list []*RTPSParticipant
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

func (domain *RTPSDomain) NewParticipant(domainId uint32, useProtocol bool, attrs *RTPSParticipantAttributes, listen *RTPSParticipantListener) *RTPSParticipant {

	participant := NewParticipant(domainId, useProtocol, attrs, listen)
	domain.Participant_list = append(domain.Participant_list, participant)

	return participant
}
