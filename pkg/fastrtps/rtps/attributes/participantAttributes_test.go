package attributes

import (
	"testing"

	common "github.com/yeren0143/DDS/common"
)

func TestNewParticipantAttributes(t *testing.T) {
	attr := NewParticipantAttributes()
	attr.RTPS.Builtin.DiscoveryConfig.DiscoveryProtocol = SIMPLE
	attr.RTPS.Builtin.DiscoveryConfig.UseSimpleEndpointDiscoveryProtocol = true
	attr.RTPS.Builtin.DiscoveryConfig.SimpleEDP.UsePublicationReaderAndSubscriptionWriter = true
	attr.RTPS.Builtin.DiscoveryConfig.SimpleEDP.UsePublicationWriterAndSubscriptionReader = true
	attr.RTPS.Builtin.DiscoveryConfig.LeaseDuration = common.CTimeInfinite
	attr.RTPS.Name = "Participant_sub_test"
}
