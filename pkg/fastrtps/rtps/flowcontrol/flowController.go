package flowcontrol

//IFlowController take a vector of cache changes (by reference) and return a filtered
//vector, with a collection of changes this filter considers valid for sending,
//ordered by its subjective priority.
type IFlowController interface {
	RegisterAsListeningController()
	DeRegisterAsListeningController()
}

//IFlowControllerSubject is a subject who own FlowControllerIs
type IFlowControllerSubject interface {
}

type flowControllerImpl struct {
}

func (flowControl *flowControllerImpl) RegisterAsListeningController() {

}

func (flowControl *flowControllerImpl) DeRegisterAsListeningController() {

}
