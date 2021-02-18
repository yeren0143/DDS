package writer

import (
	"sync"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/core/policy"
	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
)

type LivelinessCallback = func(guid *common.GUIDT, kind *policy.LivelinessQosPolicyKind,
	duration *common.DurationT, aliveChange, notAliveChange int32)

// A class managing the liveliness of a set of writers.
// Writers are represented by their LivelinessData
// Uses a shared timed event and informs outside classes on liveliness changes
type LivelinessManager struct {
	// A callback to inform outside classes that a writer changed its liveliness status
	callback LivelinessCallback
	// A boolean indicating whether we are managing writers with automatic liveliness
	manageAutomatic bool
	writers         []LivelinessData
	mutex           sync.Mutex
	// The timer owner, i.e. the writer which is next due to lose its liveliness
	timerOwner *LivelinessData

	// A timed callback expiring when a writer (the timer owner) loses its liveliness
	timer resources.TimedEvent
}
