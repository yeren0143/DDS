package data

import (
	"log"
	"time"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
)

var BUILTIN_PARTICIPANT_DATA_MAX_SIZE = uint64(100)
var TYPELOOKUP_DATA_MAX_SIZE = uint64(5000)
var DISC_BUILTIN_ENDPOINT_PARTICIPANT_ANNOUNCER = uint64(0x00000001 << 0)
var DISC_BUILTIN_ENDPOINT_PARTICIPANT_DETECTOR = uint64(0x00000001 << 1)
var DISC_BUILTIN_ENDPOINT_PUBLICATION_ANNOUNCER = uint64(0x00000001 << 2)
var DISC_BUILTIN_ENDPOINT_PUBLICATION_DETECTOR = uint64(0x00000001 << 3)
var DISC_BUILTIN_ENDPOINT_SUBSCRIPTION_ANNOUNCER = uint64(0x00000001 << 4)
var DISC_BUILTIN_ENDPOINT_SUBSCRIPTION_DETECTOR = uint64(0x00000001 << 5)
var DISC_BUILTIN_ENDPOINT_PARTICIPANT_PROXY_ANNOUNCER = uint64(0x00000001 << 6)
var DISC_BUILTIN_ENDPOINT_PARTICIPANT_PROXY_DETECTOR = uint64(0x00000001 << 7)
var DISC_BUILTIN_ENDPOINT_PARTICIPANT_STATE_ANNOUNCER = uint64(0x00000001 << 8)
var DISC_BUILTIN_ENDPOINT_PARTICIPANT_STATE_DETECTOR = uint64(0x00000001 << 9)
var BUILTIN_ENDPOINT_PARTICIPANT_MESSAGE_DATA_WRITER = uint64(0x00000001 << 10)
var BUILTIN_ENDPOINT_PARTICIPANT_MESSAGE_DATA_READER = uint64(0x00000001 << 11)
var BUILTIN_ENDPOINT_TYPELOOKUP_SERVICE_REQUEST_DATA_WRITER = uint64(0x00000001 << 12)
var BUILTIN_ENDPOINT_TYPELOOKUP_SERVICE_REQUEST_DATA_READER = uint64(0x00000001 << 13)
var BUILTIN_ENDPOINT_TYPELOOKUP_SERVICE_REPLY_DATA_WRITER = uint64(0x00000001 << 14)
var BUILTIN_ENDPOINT_TYPELOOKUP_SERVICE_REPLY_DATA_READER = uint64(0x00000001 << 15)
var DISC_BUILTIN_ENDPOINT_PUBLICATION_SECURE_ANNOUNCER = uint64(0x00000001 << 16)
var DISC_BUILTIN_ENDPOINT_PUBLICATION_SECURE_DETECTOR = uint64(0x00000001 << 17)
var DISC_BUILTIN_ENDPOINT_SUBSCRIPTION_SECURE_ANNOUNCER = uint64(0x00000001 << 18)
var DISC_BUILTIN_ENDPOINT_SUBSCRIPTION_SECURE_DETECTOR = uint64(0x00000001 << 19)
var BUILTIN_ENDPOINT_PARTICIPANT_MESSAGE_SECURE_DATA_WRITER = uint64(0x00000001 << 20)
var BUILTIN_ENDPOINT_PARTICIPANT_MESSAGE_SECURE_DATA_READER = uint64(0x00000001 << 21)
var DISC_BUILTIN_ENDPOINT_PARTICIPANT_SECURE_ANNOUNCER = uint64(0x00000001 << 26)
var DISC_BUILTIN_ENDPOINT_PARTICIPANT_SECURE_DETECTOR = uint64(0x00000001 << 27)

type ParticipantProxyData struct {
	ProtoVersion             common.ProtocolVersionT
	Guid                     common.GUIDT
	VendorIDT                common.VendorIDT
	ExpectsInlineQos         bool
	AviableBuiltinEndpoints  common.BuiltinEndpointSet
	MetatrafficLocators      common.RemoteLocatorList
	DefaultLocators          common.RemoteLocatorList
	ManualLivelinessCount    common.CountT
	ParticipantName          string
	Key                      common.InstanceHandleT
	ShouldCheckLeaseDuration bool
	LeaseDurationEvent       *resources.TimedEvent
	// Store the last timestamp it was received a RTPS message from the remote participant.
	LastReceivedMessageTm time.Time
	LeaseDuration         time.Duration
}

func (proxy *ParticipantProxyData) GetSerializedSize(includeEncapsulation bool) uint32 {
	log.Panic("not impl")
	return 0
}

func (proxy *ParticipantProxyData) WriteToCDRMessage(msg *common.CDRMessage, writeEncapsulation bool) bool {
	log.Panic("not impl")
	return false
}

func NewParticipantProxyData(att *attributes.RTPSParticipantAllocationAttributes) *ParticipantProxyData {
	return &ParticipantProxyData{}
}
