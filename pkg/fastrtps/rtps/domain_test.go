package rtps

import (
	"fmt"
	. "github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	. "github.com/yeren0143/DDS/fastrtps/rtps/participant"
	"testing"
)

func TestNewParticipant(t *testing.T) {
	var att RTPSParticipantAttributes
	participant := NewParticipant(0, true, &att, nil)
	fmt.Printf("%v", participant)
}
