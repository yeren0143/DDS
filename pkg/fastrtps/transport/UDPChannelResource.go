package transport

//UDPChannelResource is container of UDP transport
type UDPChannelResource struct {
	IChannelResource
	messageReceiver *ITransportReceiver
	transport       *UDPTransport
}
