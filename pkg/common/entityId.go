package common

type EntityIDT struct {
	Value [4]Octet
}

//NewEntityID ...
func NewEntityID(id uint32) *EntityIDT {
	var entityID = EntityIDT{}
	entityID.Value[0] = Octet(id)
	entityID.Value[1] = Octet(id >> 8)
	entityID.Value[2] = Octet(id >> 16)
	entityID.Value[3] = Octet(id >> 24)
	return &entityID
}

// EntityIdUnknown default entietyID
var (
	KEidUnknown                                          *EntityIDT
	KEidRTPSParticipant                                  *EntityIDT
	KEidSEDPBuiltinTopicWriter                           *EntityIDT
	KEidSEDPBuiltinTopicReader                           *EntityIDT
	KEidSEDPBuiltinPublicationsWriter                    *EntityIDT
	KEidSEDPBuiltinPublicationsReader                    *EntityIDT
	KEidSEDPBuiltinSubscriptionsWriter                   *EntityIDT
	KEidSEDPBuiltinSubscriptionsReader                   *EntityIDT
	KEidSPDPBuiltinRTPSParticipantWriter                 *EntityIDT
	KEidSPDPBuiltinRTPSParticipantReader                 *EntityIDT
	KEidP2PBuiltinRTPSParticipantMessageWriter           *EntityIDT
	KEidP2PBuiltinRTPSParticipantMessageReader           *EntityIDT
	KEidP2PBuiltinParticipantStalessWriter               *EntityIDT
	KEidP2PBuiltinParticipantStalessReader               *EntityIDT
	KEidTLSvcReqWriter                                   *EntityIDT
	KEidTLSvcReqReader                                   *EntityIDT
	KEidTLSvcReplyWriter                                 *EntityIDT
	KEidTLSvcReplyReader                                 *EntityIDT
	KEidSEDPBuiltinPubSecureWriter                       *EntityIDT
	KEidEDPBuiltinPubSecureReader                        *EntityIDT
	KEidSEDPBuiltinSubSecureWriter                       *EntityIDT
	KEidSEDPBuiltinSubSecureReader                       *EntityIDT
	KEidP2PBuiltinParticipantMessageSecureWriter         *EntityIDT
	KEidP2PBuiltinParticipantMessageSecureReader         *EntityIDT
	KEidP2PBuiltinParticipantVolatileMessageSecureWriter *EntityIDT
	KEidP2PBuiltinParticipantVolatileMessageSecureReader *EntityIDT
	KEidSPDPReliableBuiltinParticipantSecureWriter       *EntityIDT
	KEidSPDPReliableBuiltinParticipantSecureReader       *EntityIDT
)

// ...
func init() {
	KEidUnknown = NewEntityID(0x00000000)
	KEidRTPSParticipant = NewEntityID(0x000001c1)
	KEidSEDPBuiltinTopicWriter = NewEntityID(0x000002c2)
	KEidSEDPBuiltinTopicReader = NewEntityID(0x000002c7)
	KEidSEDPBuiltinPublicationsWriter = NewEntityID(0x000003c2)
	KEidSEDPBuiltinPublicationsReader = NewEntityID(0x000003c7)
	KEidSEDPBuiltinSubscriptionsWriter = NewEntityID(0x000004c2)
	KEidSEDPBuiltinSubscriptionsReader = NewEntityID(0x000004c7)
	KEidSPDPBuiltinRTPSParticipantWriter = NewEntityID(0x000100c2)
	KEidSPDPBuiltinRTPSParticipantReader = NewEntityID(0x000100c7)
	KEidP2PBuiltinRTPSParticipantMessageWriter = NewEntityID(0x000200C2)
	KEidP2PBuiltinRTPSParticipantMessageReader = NewEntityID(0x000200C7)
	KEidP2PBuiltinParticipantStalessWriter = NewEntityID(0x000201C3)
	KEidP2PBuiltinParticipantStalessReader = NewEntityID(0x000201C4)
	KEidTLSvcReqWriter = NewEntityID(0x000300C3)
	KEidTLSvcReqReader = NewEntityID(0x000300C4)
	KEidTLSvcReplyWriter = NewEntityID(0x000301C3)
	KEidTLSvcReplyReader = NewEntityID(0x000301C4)
	KEidSEDPBuiltinPubSecureWriter = NewEntityID(0xff0003c2)
	KEidEDPBuiltinPubSecureReader = NewEntityID(0xff0003c7)
	KEidSEDPBuiltinSubSecureWriter = NewEntityID(0xff0004c2)
	KEidSEDPBuiltinSubSecureReader = NewEntityID(0xff0004c7)
	KEidP2PBuiltinParticipantMessageSecureWriter = NewEntityID(0xff0200c2)
	KEidP2PBuiltinParticipantMessageSecureReader = NewEntityID(0xff0200c7)
	KEidP2PBuiltinParticipantVolatileMessageSecureWriter = NewEntityID(0xff0202C3)
	KEidP2PBuiltinParticipantVolatileMessageSecureReader = NewEntityID(0xff0202C4)
	KEidSPDPReliableBuiltinParticipantSecureWriter = NewEntityID(0xff0101c2)
	KEidSPDPReliableBuiltinParticipantSecureReader = NewEntityID(0xff0101c7)
}
