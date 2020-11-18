package flowcontrol

type ThroghputControllerDescriptor struct {
	BytesPerPeriod  uint32
	PeriodMillisecs uint32
}

func NewThroghputControllerDescriptor() *ThroghputControllerDescriptor {
	return &ThroghputControllerDescriptor{}
}
