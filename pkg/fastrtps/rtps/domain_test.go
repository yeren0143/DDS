package rtps

import (
	"fmt"
	"testing"

	common "github.com/yeren0143/DDS/common"
	rtpsAtt "github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	rtpsPant "github.com/yeren0143/DDS/fastrtps/rtps/participant"
)

func TestNewParticipant(t *testing.T) {
	att := rtpsAtt.NewRTPSParticipantAttributes()
	att.Builtin.DiscoveryConfig.DiscoveryProtocol = rtpsAtt.SIMPLE
	att.Builtin.DiscoveryConfig.UseSimpleEndpointDiscoveryProtocol = true
	att.Builtin.DiscoveryConfig.SimpleEDP.UsePublicationReaderAndSubscriptionWriter = true
	att.Builtin.DiscoveryConfig.SimpleEDP.UsePublicationWriterAndSubscriptionReader = true
	att.Builtin.DiscoveryConfig.LeaseDuration = common.CTimeInfinite
	att.Name = "Participant_sub_test"

	participant := rtpsPant.NewParticipant(0, true, att, nil)
	fmt.Printf("%v", participant)
}
