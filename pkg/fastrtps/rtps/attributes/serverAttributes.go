package attributes

import (
	common "github.com/yeren0143/DDS/common"
)

// port use if the ros environment variable doesn't specified one
// default server base guidPrefix
const (
	CDefaultRos2ServerPort       = uint16(11811)
	CDefaultRos2ServerGUIDPrefix = string("44.49.53.43.53.45.52.56.45.52.5F.30")
)

//RemoteServerAttributes define the attributes of the Discovery Server Protocol.
type RemoteServerAttributes struct {
	MetrafficUnicastLocatorList   *common.LocatorList
	MetrafficMulticastLocatorList *common.LocatorList
	GUIDPrefix                    *common.GUIDPrefixT
	ParticipantProxy              interface{} // Live participant proxy reference
}

//RemoteServerList define list of RemoteServerAttributes
type RemoteServerList = []*RemoteServerAttributes
