package message

import (
	"log"

	"dds/common"
	"dds/core/policy"
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
func ReadInstanceHandleFromCdrMsg(achange *common.CacheChangeT, searchPid uint16) bool {
	// Only process data when change does not already have a handle
	if achange.InstanceHandle.IsDefined() {
		return true
	}

	// Use a temporary wraping message
	msg := common.NewCDRMessageWithPayload(&achange.SerializedPayload)

	// Read encapsulation
	msg.Pos++
	var encapsulation common.Octet
	msg.ReadOctet(&encapsulation)
	if encapsulation == common.PL_CDR_BE {
		msg.MsgEndian = common.BIGEND
	} else if encapsulation == common.PL_CDR_LE {
		msg.MsgEndian = common.LITTLEEND
	} else {
		return false
	}

	achange.SerializedPayload.Encapsulation = uint16(encapsulation)

	// Skip encapsulation options
	msg.Pos += 2

	valid := true
	var pid, plength uint16
	for msg.Pos < msg.Length {
		valid = valid && msg.ReadUInt16(&pid)
		valid = valid && msg.ReadUInt16(&plength)
		if pid == policy.KPidSentinel || !valid {
			break
		}
		if pid == policy.KPidKeyHash || pid == searchPid {
			handle, ok := msg.ReadData(16)
			copy(achange.InstanceHandle.Value[:16], handle[:])
			return valid && ok
		}

		msg.Pos += uint32(int(plength+3) & int(^3))
	}

	return false
}
