package attributes

//ParticipantAttributes ...
type ParticipantAttributes struct {
	DomainID   uint32
	RTPS       *RTPSParticipantAttributes
	Allocation RTPSParticipantAllocationAttributes
}

//NewParticipantAttributes ...
func NewParticipantAttributes() *ParticipantAttributes {
	return &ParticipantAttributes{
		DomainID: 0,
		RTPS:     NewRTPSParticipantAttributes(),
	}
}
