package attributes

import (
	"fmt"
	"testing"
)

func TestRTPSParticipantAttributes(t *testing.T) {
	att := NewRTPSParticipantAttributes()
	fmt.Printf("%v", att)
}
