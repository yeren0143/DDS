package status

// Mask is a bitmap or bitset field.
type Mask uint32

// enum of StatusMask
const (
	KInconsistentTopic        Mask = (0x00000001 << 0)
	KOfferedDeadlineMissed    Mask = (0x00000001 << 1)
	KRequestDeadlineMissed    Mask = (0x00000001 << 2)
	KOfferedIncompatibleQos   Mask = (0x00000001 << 5)
	KRequestedIncompatibleQos Mask = (0x00000001 << 6)
	KSampleLost               Mask = (0x00000001 << 7)
	KSampleRejected           Mask = (0x00000001 << 8)
	KDataOnReaders            Mask = (0x00000001 << 9)
	KDataAvailable            Mask = (0x00000001 << 10)
	KLivelinessLost           Mask = (0x00000001 << 11)
	KLivelinessChanged        Mask = (0x00000001 << 12)
	KPublicationMatched       Mask = (0x00000001 << 13)
	KSubscriptionMatched      Mask = (0x00000001 << 14)
)
