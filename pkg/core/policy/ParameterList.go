package policy

import (
	"github.com/yeren0143/DDS/common"
)

// ParameterList class has static methods to update or read a list of Parameter_t
func WriteEncapsulationToCDRMsg(msg *common.CDRMessage) bool {
	valid := msg.AddOctet(0)
	valid = valid && msg.AddOctet((common.Octet)(common.PL_CDR_LE-msg.MsgEndian))
	valid = valid && msg.AddUInt16(0)
	return valid
}
