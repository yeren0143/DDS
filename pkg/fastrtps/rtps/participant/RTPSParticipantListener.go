package participant

type RTPSParticipantListener interface {
	OnParticipantDiscovery(particiant *RTPSParticipant, info *ParticipantDiscoveryInfo)
}
