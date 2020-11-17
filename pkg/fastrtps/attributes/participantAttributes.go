package attributes

import (
	. "github.com/yeren0143/DDS/fastrtps/rtps/attributes"
)

type ParticipantAttributes struct {
	DomainID uint32
	RTPS     *RTPSParticipantAttributes
}

func NewParticipantAttributes() *ParticipantAttributes {
	return &ParticipantAttributes{
		DomainID: 0,
		RTPS:     NewRTPSParticipantAttributes(),
	}
}
