package policy

import (
	"log"

	"dds/common"
)

// ParameterList class has static methods to update or read a list of Parameter_t
func WriteEncapsulationToCDRMsg(msg *common.CDRMessage) bool {
	valid := msg.AddOctet(0)
	valid = valid && msg.AddOctet((common.Octet)(common.PL_CDR_LE-msg.MsgEndian))
	valid = valid && msg.AddUInt16(0)
	return valid
}

/**
 * Read change instanceHandle from the KEY_HASH or another specific PID parameter of a CDRMessage
 * @param[in,out] change Pointer to the cache change.
 * @param[in] search_pid Specific PID to search
 * @return True when instanceHandle is updated.
 */
func ReadInstanceHandleFromCDRMsg(achange *common.CacheChangeT, searchPID uint16) bool {
	// Only process data when change does not already have a handle
	if achange.InstanceHandle.IsDefined() {
		return true
	}

	// Use a temporary wraping message
	msg := common.NewCDRMessageWithPayload(&achange.SerializedPayload)

	// Read encapsulation
	msg.Pos++
	var encapsulation common.Octet
	if !msg.ReadOctet(&encapsulation) {
		return false
	}

	switch encapsulation {
	case common.PL_CDR_BE:
		msg.MsgEndian = common.BIGEND
	case common.PL_CDR_LE:
		msg.MsgEndian = common.LITTLEEND
	default:
		return false
	}

	achange.SerializedPayload.Encapsulation = uint16(encapsulation)

	// Skip encapsulation options
	msg.Pos += 2

	valid := false
	var pid, plength uint16
	for msg.Pos < msg.Length {
		valid = true
		valid = valid && msg.ReadUInt16(&pid)
		valid = valid && msg.ReadUInt16(&plength)
		log.Fatalln("not Impl")
	}

	return false
}

type parameterProcessor = func(msg *common.CDRMessage, pid ParameterIDT, plength uint16) bool

// Read a parameterList from a CDRMessage
func ReadParameterListFromCDRMsg(msg *common.CDRMessage, processor parameterProcessor,
	useEncapsulation bool, qosSize *uint32) bool {
	*qosSize = 0
	if useEncapsulation {
		msg.Pos += 1
		var encapsulation common.Octet
		msg.ReadOctet(&encapsulation)
		if encapsulation == common.PL_CDR_BE {
			msg.MsgEndian = common.BIGEND
		} else if encapsulation == common.PL_CDR_LE {
			msg.MsgEndian = common.LITTLEEND
		} else {
			log.Println("read encapsulation error")
			return false
		}
		msg.Pos += 2
	}

	originalPos := msg.Pos
	isSentinel := false
	for !isSentinel {
		msg.Pos = originalPos + *qosSize

		var pid ParameterIDT
		var plength uint16
		valid := msg.ReadUInt16(&pid)
		valid = valid && msg.ReadUInt16(&plength)

		if pid == KPidSentinel {
			// PID_SENTINEL is always considered of length 0
			plength = 0
			isSentinel = true
		}

		*qosSize += uint32(4 + plength)

		// Align to 4 byte boundary and prepare for next iteration
		*qosSize = uint32(int32(*qosSize+3) & (^int32(3)))

		if (!valid) || ((msg.Pos + uint32(plength)) > msg.Length) {
			return false
		} else if !isSentinel {
			if !processor(msg, pid, plength) {
				return false
			}
		}
	}
	return true
}
