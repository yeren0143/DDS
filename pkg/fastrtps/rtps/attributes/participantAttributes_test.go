package attributes

import (
	"testing"

	common "dds/common"
)

func TestNewParticipantAttributes(t *testing.T) {
	attr := NewParticipantAttributes()
	attr.RTPS.Builtin.DiscoveryConfig.Protocol = KDisPSimple
	attr.RTPS.Builtin.DiscoveryConfig.UseSimpleEndpoint = true
	attr.RTPS.Builtin.DiscoveryConfig.SimpleEDP.UsePublicationReaderAndSubscriptionWriter = true
	attr.RTPS.Builtin.DiscoveryConfig.SimpleEDP.UsePublicationWriterAndSubscriptionReader = true
	attr.RTPS.Builtin.DiscoveryConfig.LeaseDuration = common.KTimeInfinite
	attr.RTPS.Name = "Participant_sub_test"
}
