package resources

import "time"

/**
 * This class centralizes all operations over timed events in the same thread.
 * @ingroup MANAGEMENT_MODULE
 */
type ResourceEvent struct {
	Pending_timers []*TimedEvent
	Active_timers  []*TimedEvent
	current_time   time.Time
}
