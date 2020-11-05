package attributes

type ParticipantAttributes struct {
	DomainId uint32
	RTPS     *RTPSParticipantAttributes
}

func NewParticipantAttributes() *ParticipantAttributes {
	return &ParticipantAttributes{
		DomainId: 0,
		RTPS:     NewRTPSParticipantAttributes(),
	}
}
