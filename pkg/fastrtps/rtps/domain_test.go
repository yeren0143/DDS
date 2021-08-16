package rtps

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"testing"

	"dds/common"
	"dds/fastrtps/rtps/attributes"
)

func TestNewParticipant(t *testing.T) {
	go StartHTTPDebuger()

	log.SetFlags(log.Ldate | log.Llongfile)

	os.Setenv("SHM_TRANSPORT_BUILTIN", "1")
	os.Setenv("FASTDDS_SHM_TRANSPORT_DISABLED", "true")

	att := attributes.NewRTPSParticipantAttributes()
	att.Builtin.DiscoveryConfig.Protocol = attributes.KDisPSimple
	att.Builtin.DiscoveryConfig.UseSimpleEndpoint = true
	att.Builtin.DiscoveryConfig.SimpleEDP.UsePublicationReaderAndSubscriptionWriter = true
	att.Builtin.DiscoveryConfig.SimpleEDP.UsePublicationWriterAndSubscriptionReader = true
	att.Builtin.DiscoveryConfig.LeaseDuration = common.KTimeInfinite
	att.Name = "Participant_sub_test"

	pparam := attributes.NewParticipantAttributes()
	pparam.RTPS = att

	participant := NewRTPSParticipant(pparam, nil)
	fmt.Printf("%v", participant)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	select {
	case sig := <-sigChan:
		fmt.Printf("接受到来自系统的信号(%v)：", sig)
	}
}
