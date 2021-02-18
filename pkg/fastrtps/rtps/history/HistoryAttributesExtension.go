package history

import (
	"math"

	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/utils"
)

func ResourceLimitsFromHistory(att *attributes.HistoryAttributes, increment uint32) utils.ResourceLimitedContainerConfig {
	if att.MaximumReservedCaches > 0 && att.InitialReservedCaches == att.MaximumReservedCaches {
		return utils.ResourceLimitedContainerConfig{
			Initial:   att.MaximumReservedCaches,
			Maximum:   att.MaximumReservedCaches,
			Increment: 0,
		}
	}

	var config utils.ResourceLimitedContainerConfig
	if att.InitialReservedCaches <= 0 {
		config.Initial = 0
	}
	if att.MaximumReservedCaches <= 0 {
		config.Maximum = math.MaxUint32
	}
	if increment == 0 {
		config.Increment = 1
	}

	return config
}
