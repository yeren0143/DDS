package participant

type ParticipantListener struct {
}

func (pl *ParticipantListener) OnParticipantDiscovery(participant *Participant, info *ParticipantDiscoveryInfo) {

}

func (pl *ParticipantListener) OnSubscriberDiscovery(participant *Participant, info *ReaderDiscoveryInfo) {

}

func (pl *ParticipantListener) OnPublisherDiscovery(participant *Participant, info *WriterDiscoveryInfo) {

}
