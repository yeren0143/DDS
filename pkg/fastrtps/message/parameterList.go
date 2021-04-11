package message

import (
	"log"

	"github.com/yeren0143/DDS/common"
)

// Write parameterList encapsulation to the CDRMessage.
func WriteEncapsulationToCdrMsg(msg *common.CDRMessage) bool {
	log.Fatalln("not impl")
	return true
}

// Update the information of a cache change parsing the inline qos from a CDRMessage
// change Reference to the cache change to be updated.
func UpdateCacheChangeFromInlineQos(change *common.CacheChangeT, msg *common.CDRMessage,
	qosSize *uint32) bool {
	log.Fatalln("not impl")
	return true
}

// Read a parameterList from a CDRMessage
func ReadParameterListfromCdrMsg(msg *common.CDRMessage) bool {
	log.Fatalln("not impl")
	return true
}

// Read guid from the KEY_HASH or another specific PID parameter of a CDRMessage
func ReadGUIDfromCdrMsg(msg *common.CDRMessage, searchPid uint16, guid *common.GUIDT) bool {
	log.Fatalln("not impl")
	return true
}

func ParameterProcess(change *common.CacheChangeT, msg *common.CDRMessage, qosSize uint32) bool {
	log.Fatalln("not impl")
	return true
}

// Read change instanceHandle from the KEY_HASH or another specific PID parameter of a CDRMessage
func ReadInstanceHandleFromCdrMsg(change *common.CacheChangeT, searchPid uint16) bool {
	log.Fatalln("not impl")
	return true
}
