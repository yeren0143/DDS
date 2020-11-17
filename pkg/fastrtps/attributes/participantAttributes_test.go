package attributes

import (
	"testing"

	. "github.com/yeren0143/DDS/fastrtps/rtps/attributes"
)

func TestNewParticipantAttributes(t testing.T) {
	attr := NewParticipantAttributes()
	attr.RTPS.Builtin.DiscoveryConfig.DiscoveryProtocol = SIMPLE
}
