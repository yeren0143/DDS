package participant

import (
	. "github.com/yeren0143/DDS/fastrtps/rtps/reader"
	. "github.com/yeren0143/DDS/fastrtps/rtps/writer"
)

type ParticipantListener struct {
}

func (pl *ParticipantListener) OnParticipantDiscovery(participant *Participant, info *ParticipantDiscoveryInfo) {

}

func (pl *ParticipantListener) OnSubscriberDiscovery(participant *Participant, info *ReaderDiscoveryInfo) {

}

func (pl *ParticipantListener) OnPublisherDiscovery(participant *Participant, info *WriterDiscoveryInfo) {

}
