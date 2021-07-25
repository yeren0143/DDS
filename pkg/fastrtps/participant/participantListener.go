package participant

import (
	"log"

	"github.com/yeren0143/DDS/fastrtps/rtps/reader"
	"github.com/yeren0143/DDS/fastrtps/rtps/writer"
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
