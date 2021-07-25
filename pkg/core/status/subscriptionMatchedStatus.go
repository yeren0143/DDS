package status

import (
	"github.com/yeren0143/DDS/common"
)

// SubscriptionMatchedStatus storing the subscription status
type SubscriptionMatchedStatus struct {
	MatchedStatus

	// Handle to the last writer that matched the reader causing the status change
	LastSubscriptionHandle common.InstanceHandleT
}
