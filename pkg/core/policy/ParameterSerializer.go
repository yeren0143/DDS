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
	valid = valid && msg.AddOctet(parameter.VendorID.Vendor[0])
	valid = valid && msg.AddOctet(parameter.VendorID.Vendor[1])
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
