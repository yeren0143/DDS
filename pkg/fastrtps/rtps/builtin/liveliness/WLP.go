package liveliness

type IWlpParent interface {
}

// Class WLP that implements the Writer Liveliness Protocol described in the RTPS specification.
type WLP struct {
	// Minimum time among liveliness periods of automatic writers, in milliseconds
	minAutomaticMs float64
	// Minimum time among liveliness periods of manual by participant writers, in milliseconds
	minManualByParticipantMs float64
	participant              IWlpParent
}
