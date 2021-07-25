package policy

import (
	"github.com/yeren0143/DDS/common"
)

func AddCommonToMsg(parameter *ParameterT, msg *common.CDRMessage) bool {
	valid := msg.AddUInt16(parameter.Pid)
	valid = valid && msg.AddUInt16(parameter.Length)
	return valid
}

func AddGuidToMsg(parameter *ParameterGuidT, msg *common.CDRMessage) bool {
	valid := AddCommonToMsg(&parameter.ParameterT, msg)
	valid = valid && msg.AddData(parameter.Guid.Prefix.Value[:12])
	valid = valid && msg.AddData(parameter.Guid.EntityID.Value[:4])
	return valid
}

func AddLocatorToMsg(parameter *ParamaterLocatorT, msg *common.CDRMessage) bool {
	valid := AddCommonToMsg(&parameter.ParameterT, msg)
	valid = valid && msg.AddLocator(&parameter.Locator)
	return valid
}

func AddTimeToMsg(parameter *ParameterTimeT, msg *common.CDRMessage) bool {
	valid := AddCommonToMsg(&parameter.ParameterT, msg)
	valid = valid && msg.AddInt32(parameter.Time.Seconds)
	valid = valid && msg.AddInt32(int32(parameter.Time.Nanosec))
	return valid
}

func AddBuiltinEndpointSetToMsg(parameter *ParameterBuiltinEndpointSetT, msg *common.CDRMessage) bool {
	valid := AddCommonToMsg(&parameter.ParameterT, msg)
	valid = valid && msg.AddUInt32(parameter.EndpointSet)
	return valid
}

func AddParameterStringToMsg(parameter *ParameterStringT, msg *common.CDRMessage) bool {
	if len(parameter.Name) == 0 {
		return false
	}

	valid := msg.AddUInt16(parameter.Pid)
	strSize := len(parameter.Name) + 1
	len := (strSize + 4 + 3) & (^3)
	valid = valid && msg.AddUInt16(uint16(len))
	valid = valid && msg.AddString(parameter.Name)
	return valid
}

func AddProtocolVersionToMsg(parameter *ParameterProtocolVersionT, msg *common.CDRMessage) bool {
	valid := AddCommonToMsg(&parameter.ParameterT, msg)
	valid = valid && msg.AddOctet(parameter.ProtocolVersion.Major)
	valid = valid && msg.AddOctet(parameter.ProtocolVersion.Minor)
	valid = valid && msg.AddUInt16(0)
	return valid
}

func AddVendorIDToMsg(parameter *ParameterVendorIDT, msg *common.CDRMessage) bool {
	valid := AddCommonToMsg(&parameter.ParameterT, msg)
	valid = valid && msg.AddOctet(parameter.VendorID.Value[0])
	valid = valid && msg.AddOctet(parameter.VendorID.Value[1])
	valid = valid && msg.AddUInt16(0)
	return valid
}

func AddParameterSentinelToMsg(msg *common.CDRMessage) bool {
	if msg.Pos+4 > msg.MaxSize {
		return false
	}

	msg.AddUInt16(KPidSentinel)
	msg.AddUInt16(0)

	return true
}

func ReadVendorIdFromCDRMessage(parameter *ParameterVendorIDT, msg *common.CDRMessage, parameterLength uint16) bool {
	if parameterLength != KParameterVendorLength {
		return false
	}
	parameter.Length = parameterLength
	valid := msg.ReadOctet(&parameter.VendorID.Value[0])
	valid = valid && msg.ReadOctet(&parameter.VendorID.Value[1])
	msg.Pos += 2
	return valid
}

func ReadParameterProtocolFromCDRMessage(parameter *ParameterProtocolVersionT, msg *common.CDRMessage, parameterLength uint16) bool {
	if parameterLength != KParameterProtocolLength {
		return false
	}
	parameter.Length = parameterLength
	valid := msg.ReadOctet(&parameter.ProtocolVersion.Major)
	valid = valid && msg.ReadOctet(&parameter.ProtocolVersion.Minor)
	msg.Pos += 2
	return valid
}

func ReadGuidFromCDRMessage(parameter *ParameterGuidT, msg *common.CDRMessage, parameterLength uint16) bool {
	if parameterLength != KParameterGuidLength {
		return false
	}
	parameter.Length = parameterLength
	prefix, ok1 := msg.ReadData(12)
	copy(parameter.Guid.Prefix.Value[:12], prefix[:])
	entityID, ok2 := msg.ReadData(4)
	copy(parameter.Guid.EntityID.Value[:4], entityID[:])
	return ok1 && ok2
}

func ReadLocatorFromCDRMessage(parameter *ParamaterLocatorT, msg *common.CDRMessage, parameterLength uint16) bool {
	if parameterLength != KParameterLocatorLength {
		return false
	}
	parameter.Length = parameterLength
	return msg.ReadLocator(&parameter.Locator)
}

func ReadTimeFromCDRMessage(parameter *ParameterTimeT, msg *common.CDRMessage, parameterLength uint16) bool {
	if parameterLength != KParameterTimeLength {
		return false
	}
	parameter.Length = parameterLength
	var sec int32
	valid := msg.ReadInt32(&sec)
	parameter.Time.Seconds = sec
	var frac uint32
	valid = valid && msg.ReadUInt32(&frac)
	parameter.Time.Nanosec = frac
	return valid
}

func ReadBuiltinEndpointSetFromCDRMessage(parameter *ParameterBuiltinEndpointSetT, msg *common.CDRMessage, parameterLength uint16) bool {

	if parameterLength != KParameterBuiltinEndpointsetLength {
		return false
	}
	parameter.Length = parameterLength
	return msg.ReadUInt32(&parameter.EndpointSet)
}

func ReadEntityNameFromCDRMessage(parameter *ParameterStringT, msg *common.CDRMessage, parameterLength uint16) bool {
	if parameterLength > 256 {
		return false
	}

	parameter.Length = parameterLength
	stri, valid := msg.ReadString()
	parameter.Name = stri
	return valid
}
