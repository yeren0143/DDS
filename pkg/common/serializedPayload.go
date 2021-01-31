package common

const (
	CDR_BE    = 0x0000
	CDR_LE    = 0x0001
	PL_CDR_BE = 0x0002
	PL_CDR_LE = 0x0003
)

type SerializedPayloadT struct {
	Encapsulation uint16
	Length        uint32
	Data          []Octet
	MaxSize       uint32
	Pos           uint32 //!Position when reading
}

func CreateSerializedPayload() SerializedPayloadT {
	return SerializedPayloadT{
		Encapsulation: CDR_BE,
		Length:        0,
		MaxSize:       0,
		Pos:           0,
	}
}
