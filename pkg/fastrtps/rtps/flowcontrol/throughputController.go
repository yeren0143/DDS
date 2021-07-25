package flowcontrol

var _ IFlowController = (*ThroughputController)(nil)

// ThroughputController is aSimple filter that only clears changes up to a certain accumulated payload size.
// It refreshes after a given time in MS, in a staggered way (e.g. if it clears
// 500kb at t=0 and 800 kb at t=10, it will refresh 500kb at t = 0 + period, and
// then fully refresh at t = 10 + period).
type ThroughputController struct {
	flowControllerImpl
}

//NewThroughputController create ThroughputController
func NewThroughputController(descriptor *ThroughputControllerDescriptor, subject IFlowControllerParent) *ThroughputController {
	return &ThroughputController{}
}
