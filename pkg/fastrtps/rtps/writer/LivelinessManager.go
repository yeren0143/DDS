package writer

import (
	"log"
	"math"
	"sync"
	"time"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/core/policy"
	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
)

type LivelinessCallback = func(guid *common.GUIDT, kind policy.LivelinessQosPolicyKind,
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
	timer *resources.TimedEvent
}

func (manager *LivelinessManager) calculateNext() bool {
	manager.timerOwner = nil
	minTime := time.Now().Add(math.MaxInt32 * time.Nanosecond)

	anyAlive := false
	for _, awriter := range manager.writers {
		if awriter.Status == KAlive {
			if awriter.TimePoint.Before(minTime) {
				minTime = awriter.TimePoint
				manager.timerOwner = &awriter
			}
			anyAlive = true
		}
	}

	return anyAlive
}

func (manager *LivelinessManager) timerExpired() bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if manager.timerOwner == nil {
		log.Fatalln("Liveliness timer expired but there is no writer")
		return false
	}

	if manager.callback != nil {
		manager.callback(&manager.timerOwner.GUID, manager.timerOwner.Kind,
			&manager.timerOwner.LeaseDuration, -1, 1)
	}
	manager.timerOwner.Status = KNotAlive

	if manager.calculateNext() {
		// Some times the interval could be negative if a writer expired during the call to this function
		// Once in this situation there is not much we can do but let asio timers expire inmediately
		interval := manager.timerOwner.TimePoint.Sub(time.Now())
		manager.timer.UpdateInterval(common.Time{Nanosec: uint32(interval.Nanoseconds())})
		return true
	}

	return false
}

type livelinessCallback = func(guid *common.GUIDT, akind policy.LivelinessQosPolicyKind,
	leaseDuration *common.DurationT, aliveCount int32, notAliveCount int32)

func NewLivelinessManager(callback *livelinessCallback, event *resources.ResourceEvent, manageAutomatic bool) *LivelinessManager {
	var livelinessManager LivelinessManager
	livelinessManager.callback = *callback
	livelinessManager.manageAutomatic = manageAutomatic
	interval := func() bool {
		return livelinessManager.timerExpired()
	}
	livelinessManager.timer = resources.NewTimedEvent(event, &interval, 0)
	return &livelinessManager
}
