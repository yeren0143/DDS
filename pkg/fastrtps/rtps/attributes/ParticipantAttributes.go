package attributes

//ParticipantAttributes ...
type ParticipantAttributes struct {
	DomainID uint32
	RTPS     *RTPSParticipantAttributes
}

//NewParticipantAttributes ...
func NewParticipantAttributes() *ParticipantAttributes {
	return &ParticipantAttributes{
		DomainID: 0,
		RTPS:     NewRTPSParticipantAttributes(),
	}
}
