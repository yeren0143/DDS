package rtps

import (
	"fmt"
	"testing"
)

func TestNewParticipant(t *testing.T) {
	var att RTPSParticipantAttributes
	//var listen ParticipantListener
	participant := NewParticipant(&att, nil)
	fmt.Printf("%v", participant)
}
