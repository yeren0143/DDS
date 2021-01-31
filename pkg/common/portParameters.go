package common

type PortParameters struct {
	PortBase          uint16
	DomainIDGain      uint16
	ParticipantIDGain uint16
	Offsetd0          uint16
	Offsetd1          uint16
	Offsetd2          uint16
	Offsetd3          uint16
}

func NewDefaultPortParameters() *PortParameters {
	return &PortParameters{
		PortBase:          7400,
		DomainIDGain:      250,
		ParticipantIDGain: 2,
		Offsetd0:          0,
		Offsetd1:          10,
		Offsetd2:          1,
		Offsetd3:          11,
	}
}

func (portParams *PortParameters) GetMulticastPort(domainId uint32) uint32 {
	port := uint32(portParams.PortBase) +
		uint32(portParams.DomainIDGain)*domainId + uint32(portParams.Offsetd0)

	if port > 65535 {
		panic("Calculated port number is too high. Probably the domainId is over 232," +
			"or portBase is too high.")
	}

	return port
}

//GetUnicastPort return unicast port
func (portParams *PortParameters) GetUnicastPort(domaindID uint32, RTPSParticipantID uint32) uint32 {
	port := uint32(portParams.PortBase) + uint32(portParams.DomainIDGain)*domaindID +
		uint32(portParams.Offsetd1) + uint32(portParams.ParticipantIDGain)*uint32(RTPSParticipantID)

	if port > 65535 {
		panic("Calculated port number is too high. Probably the domainId is over 232, there are " +
			"too much participants created or portBase is too high.")
	}
	return port
}
