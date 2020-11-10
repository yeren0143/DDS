package rtps

import (
	"fmt"
	"testing"
)

func TestNewParticipant(t *testing.T) {
	var att RTPSParticipantAttributes
	participant := NewParticipant(0, &att)
	fmt.Printf("%v", participant)
}
