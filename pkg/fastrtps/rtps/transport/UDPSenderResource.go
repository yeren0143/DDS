package transport

type UDPSenderResource struct {
	SenderResource
	onlyMulticastPurpose bool
}
