package attributes

import (
	"fmt"
	"testing"
	//common "github.com/yeren0143/DDS/common"
)

func TestRTPSParticipantAttributes(t *testing.T) {
	att := NewRTPSParticipantAttributes()
	fmt.Printf("%v", att)
}

func TestNewParticipantAttributes(t *testing.T) {
	attr := NewParticipantAttributes()
	fmt.Println(attr)
	attr.RTPS.Builtin.DiscoveryConfig.DiscoveryProtocol = SIMPLE
}
