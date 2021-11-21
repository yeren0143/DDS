package endpoint

import (
	"dds/common"
	"dds/core/policy"
)

type IWlp interface {
	InitWL(p interface{}) bool
	AddWriter(guid *common.GUIDT, kind policy.LivelinessQosPolicyKind, leaseDuration *common.DurationT) bool
	RemoveWriter(guid *common.GUIDT, kind policy.LivelinessQosPolicyKind, leaseDuration *common.DurationT) bool
	AssertLiveliness(writerGuid *common.GUIDT, kind policy.LivelinessQosPolicyKind, leaseDuration *common.DurationT) bool
	SubLivelinessManager() ILivelinessManager
}

type ILivelinessManager interface {
	AddWriter(guid common.GUIDT, kind policy.LivelinessQosPolicyKind, leaseDuration common.DurationT) bool
}
