package core

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/core/status"
)

// The Entity class is the abstract base class for all the objects that support QoS policies,
// a listener and a status condition.
type Entity struct {
	// StatusMask with relevant statuses set to 1
	statusMask status.Mask
	// StatusMask with triggered statuses set to 1
	statusChanges status.Mask
	// InstanceHandle associated to the Entity
	instanceHandle common.InstanceHandleT

	// Boolean that states if the Entity is enabled or disabled
	enable bool
}
