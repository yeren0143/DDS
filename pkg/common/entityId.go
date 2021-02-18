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

var (
	KEIDUnknown                                          = NewEntityID(0x00000000)
	KEntityIDDSServerVirtualWriter                       = NewEntityID(0x00030073)
	KDsServerVirtualWriter                               = KEntityIDDSServerVirtualWriter
	KEntityIDDSServerVirtualReader                       = NewEntityID(0x00030074)
	KDsServerVirtualReader                               = KEntityIDDSServerVirtualReader
	KEidRTPSParticipant                                  = NewEntityID(0x000001c1)
	KEntityIDRTPSParticipant                             = KEidRTPSParticipant
	KEidSEDPBuiltinTopicWriter                           = NewEntityID(0x000002c2)
	KEidSEDPBuiltinTopicReader                           = NewEntityID(0x000002c7)
	KEidSEDPBuiltinPublicationsWriter                    = NewEntityID(0x000003c2)
	KEntityIDSEDPPubWriter                               = KEidSEDPBuiltinPublicationsWriter
	KEidSEDPBuiltinPublicationsReader                    = NewEntityID(0x000003c7)
	KEntityIDSEDPPubReader                               = KEidSEDPBuiltinPublicationsReader
	KEidSEDPBuiltinSubscriptionsWriter                   = NewEntityID(0x000004c2)
	KEntityIDSEDPSubWriter                               = KEidSEDPBuiltinSubscriptionsWriter
	KEidSEDPBuiltinSubscriptionsReader                   = NewEntityID(0x000004c7)
	KEntityIDSEDPSubReader                               = KEidSEDPBuiltinSubscriptionsReader
	KEidSPDPBuiltinRTPSParticipantWriter                 = NewEntityID(0x000100c2)
	KEntityIDSPDPWriter                                  = KEidSPDPBuiltinRTPSParticipantWriter
	KEidSPDPBuiltinRTPSParticipantReader                 = NewEntityID(0x000100c7)
	KEntityIDSPDPReader                                  = KEidSPDPBuiltinRTPSParticipantReader
	KEidP2PBuiltinRTPSParticipantMessageWriter           = NewEntityID(0x000200C2)
	KEntityIDWriterLiveliness                            = KEidP2PBuiltinRTPSParticipantMessageWriter
	KEidP2PBuiltinRTPSParticipantMessageReader           = NewEntityID(0x000200C7)
	KEntityIDReaderLiveliness                            = KEidP2PBuiltinRTPSParticipantMessageReader
	KEidP2PBuiltinParticipantStalessWriter               = NewEntityID(0x000201C3)
	KParticipantStatelessMessageWriterEntityID           = KEidP2PBuiltinParticipantStalessWriter
	KEidP2PBuiltinParticipantStalessReader               = NewEntityID(0x000201C4)
	KParticipantStatelessMessageEeaderEntityID           = KEidP2PBuiltinParticipantStalessReader
	KEidTLSvcReqWriter                                   = NewEntityID(0x000300C3)
	KEntityIDTypeLookupRequestWriter                     = KEidTLSvcReqWriter
	KEidTLSvcReqReader                                   = NewEntityID(0x000300C4)
	KEntityIDTypeLookupRequestReader                     = KEidTLSvcReqReader
	KEidTLSvcReplyWriter                                 = NewEntityID(0x000301C3)
	KEntityIDTypeLookupReplyWriter                       = KEidTLSvcReplyWriter
	KEidTLSvcReplyReader                                 = NewEntityID(0x000301C4)
	KEntityIDTypeLookupReplyReader                       = KEidTLSvcReplyReader
	KEidSEDPBuiltinPubSecureWriter                       = NewEntityID(0xff0003c2)
	KEidEDPBuiltinPubSecureReader                        = NewEntityID(0xff0003c7)
	KEidSEDPBuiltinSubSecureWriter                       = NewEntityID(0xff0004c2)
	KEidSEDPBuiltinSubSecureReader                       = NewEntityID(0xff0004c7)
	KEidP2PBuiltinParticipantMessageSecureWriter         = NewEntityID(0xff0200c2)
	KEidP2PBuiltinParticipantMessageSecureReader         = NewEntityID(0xff0200c7)
	KEidP2PBuiltinParticipantVolatileMessageSecureWriter = NewEntityID(0xff0202C3)
	KEidP2PBuiltinParticipantVolatileMessageSecureReader = NewEntityID(0xff0202C4)
	KEidSPDPReliableBuiltinParticipantSecureWriter       = NewEntityID(0xff0101c2)
	KEidSPDPReliableBuiltinParticipantSecureReader       = NewEntityID(0xff0101c7)
)
