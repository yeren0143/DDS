package message

import (
	"github.com/yeren0143/DDS/common"
)

func addHeader(msg *common.CDRMessage, guidPrefix *common.GUIDPrefixT,
	version *common.ProtocolVersionT, vendorID *common.VendorIDT) bool {
	msg.AddOctet('R')
	msg.AddOctet('T')
	msg.AddOctet('P')
	msg.AddOctet('S')

	msg.AddOctet(version.Major)
	msg.AddOctet(version.Minor)

	msg.AddOctet(vendorID.Vendor[0])
	msg.AddOctet(vendorID.Vendor[1])

	if msg.Pos+12 < msg.Length {
		for i := 0; i < 12; i++ {
			msg.AddOctet(guidPrefix.Value[i])
		}
	}

	return true
}
