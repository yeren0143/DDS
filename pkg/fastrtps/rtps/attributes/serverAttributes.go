package attributes

import (
	. "github.com/yeren0143/DDS/common"
)

const (
	DEFAULT_ROS2_SERVER_PORT       = uint16(11811)
	DEFAULT_ROS2_SERVER_GUIDPREFIX = string("44.49.53.43.53.45.52.56.45.52.5F.30")
)

type RemoteServerAttributes struct {
	MetrafficUnicastLocatorList   LocatorList
	MetrafficMulticastLocatorList LocatorList
	GuidPrefix                    GuidPrefix_t
	ParticipantProxy              interface{} // Live participant proxy reference
}

type RemoteServerList = []RemoteServerAttributes
