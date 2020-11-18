package common

//EntityID ...
type EntityID struct {
	Value [4]Octet
}

//NewEntityID ...
func NewEntityID(id uint32) *EntityID {
	var entityID = EntityID{}
	entityID.Value[0] = Octet(id)
	entityID.Value[1] = Octet(id >> 8)
	entityID.Value[2] = Octet(id >> 16)
	entityID.Value[3] = Octet(id >> 24)
	return &entityID
}

// EntityIdUnknown default entietyID
var (
	CEidUnknown                                          *EntityID
	CEidRTPSParticipant                                  *EntityID
	CEidSEDPBuiltinTopicWriter                           *EntityID
	CEidSEDPBuiltinTopicReader                           *EntityID
	CEidSEDPBuiltinPublicationsWriter                    *EntityID
	CEidSEDPBuiltinPublicationsReader                    *EntityID
	CEidSEDPBuiltinSubscriptionsWriter                   *EntityID
	CEidSEDPBuiltinSubscriptionsReader                   *EntityID
	CEidSPDPBuiltinRTPSParticipantWriter                 *EntityID
	CEidSPDPBuiltinRTPSParticipantReader                 *EntityID
	CEidP2PBuiltinRTPSParticipantMessageWriter           *EntityID
	CEidP2PBuiltinRTPSParticipantMessageReader           *EntityID
	CEidP2PBuiltinParticipantStalessWriter               *EntityID
	CEidP2PBuiltinParticipantStalessReader               *EntityID
	CEidTLSvcReqWriter                                   *EntityID
	CEidTLSvcReqReader                                   *EntityID
	CEidTLSvcReplyWriter                                 *EntityID
	CEidTLSvcReplyReader                                 *EntityID
	CEidSEDPBuiltinPubSecureWriter                       *EntityID
	CEidEDPBuiltinPubSecureReader                        *EntityID
	CEidSEDPBuiltinSubSecureWriter                       *EntityID
	CEidSEDPBuiltinSubSecureReader                       *EntityID
	CEidP2PBuiltinParticipantMessageSecureWriter         *EntityID
	CEidP2PBuiltinParticipantMessageSecureReader         *EntityID
	CEidP2PBuiltinParticipantVolatileMessageSecureWriter *EntityID
	CEidP2PBuiltinParticipantVolatileMessageSecureReader *EntityID
	CEidSPDPReliableBuiltinParticipantSecureWriter       *EntityID
	CEidSPDPReliableBuiltinParticipantSecureReader       *EntityID
)

// ...
func init() {
	CEidUnknown = NewEntityID(0x00000000)
	CEidRTPSParticipant = NewEntityID(0x000001c1)
	CEidSEDPBuiltinTopicWriter = NewEntityID(0x000002c2)
	CEidSEDPBuiltinTopicReader = NewEntityID(0x000002c7)
	CEidSEDPBuiltinPublicationsWriter = NewEntityID(0x000003c2)
	CEidSEDPBuiltinPublicationsReader = NewEntityID(0x000003c7)
	CEidSEDPBuiltinSubscriptionsWriter = NewEntityID(0x000004c2)
	CEidSEDPBuiltinSubscriptionsReader = NewEntityID(0x000004c7)
	CEidSPDPBuiltinRTPSParticipantWriter = NewEntityID(0x000100c2)
	CEidSPDPBuiltinRTPSParticipantReader = NewEntityID(0x000100c7)
	CEidP2PBuiltinRTPSParticipantMessageWriter = NewEntityID(0x000200C2)
	CEidP2PBuiltinRTPSParticipantMessageReader = NewEntityID(0x000200C7)
	CEidP2PBuiltinParticipantStalessWriter = NewEntityID(0x000201C3)
	CEidP2PBuiltinParticipantStalessReader = NewEntityID(0x000201C4)
	CEidTLSvcReqWriter = NewEntityID(0x000300C3)
	CEidTLSvcReqReader = NewEntityID(0x000300C4)
	CEidTLSvcReplyWriter = NewEntityID(0x000301C3)
	CEidTLSvcReplyReader = NewEntityID(0x000301C4)
	CEidSEDPBuiltinPubSecureWriter = NewEntityID(0xff0003c2)
	CEidEDPBuiltinPubSecureReader = NewEntityID(0xff0003c7)
	CEidSEDPBuiltinSubSecureWriter = NewEntityID(0xff0004c2)
	CEidSEDPBuiltinSubSecureReader = NewEntityID(0xff0004c7)
	CEidP2PBuiltinParticipantMessageSecureWriter = NewEntityID(0xff0200c2)
	CEidP2PBuiltinParticipantMessageSecureReader = NewEntityID(0xff0200c7)
	CEidP2PBuiltinParticipantVolatileMessageSecureWriter = NewEntityID(0xff0202C3)
	CEidP2PBuiltinParticipantVolatileMessageSecureReader = NewEntityID(0xff0202C4)
	CEidSPDPReliableBuiltinParticipantSecureWriter = NewEntityID(0xff0101c2)
	CEidSPDPReliableBuiltinParticipantSecureReader = NewEntityID(0xff0101c7)
}
