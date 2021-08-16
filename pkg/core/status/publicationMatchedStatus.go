package status

import (
	"dds/common"
)

// PublicationMatchedStatus storing the publication status
type PublicationMatchedStatus struct {
	MatchedStatus

	// Handle to the last reader that matched the writer causing the status to change
	LastSubscriptionHandle common.InstanceHandleT
}
