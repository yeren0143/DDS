package message

import (
	common "github.com/yeren0143/DDS/common"
	reader "github.com/yeren0143/DDS/fastrtps/rtps/reader"
	writer "github.com/yeren0143/DDS/fastrtps/rtps/writer"
)

// Receiver process the received messages.
type Receiver struct {
	sourceVersion     common.ProtocolVertionT
	sourceVendorIDT   common.VendorIDT
	sourceGUIDPrefix  common.GuidPrefix
	destGUIDPrefix    common.GuidPrefix
	haveTimeStamp     bool
	timeStamp         common.Time
	associatedWriters []*writer.RTPSWriter
	associatedReaders map[common.EntityId]*reader.RTPSReader
}
