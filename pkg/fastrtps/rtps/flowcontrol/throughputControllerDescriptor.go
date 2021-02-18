package flowcontrol

type ThroughputControllerDescriptor struct {
	BytesPerPeriod  uint32
	PeriodMillisecs uint32
}

func NewThroughputControllerDescriptor() *ThroughputControllerDescriptor {
	return &ThroughputControllerDescriptor{}
}
