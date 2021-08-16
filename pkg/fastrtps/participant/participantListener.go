package participant

import (
	"log"

	"dds/fastrtps/rtps/reader"
	"dds/fastrtps/rtps/writer"
)

type ParticipantListener struct {
}

func (pl *ParticipantListener) OnParticipantDiscovery(participant *Participant, info *ParticipantDiscoveryInfo) {
	log.Fatalln("not impl")
}

func (pl *ParticipantListener) OnSubscriberDiscovery(participant *Participant, info *reader.ReaderDiscoveryInfo) {
	log.Fatalln("not impl")
}

func (pl *ParticipantListener) OnPublisherDiscovery(participant *Participant, info *writer.WriterDiscoveryInfo) {
	log.Fatalln("not impl")
}
