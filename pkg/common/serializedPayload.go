package common

const (
	CDR_BE    = 0x0000
	CDR_LE    = 0x0001
	PL_CDR_BE = 0x0002
	PL_CDR_LE = 0x0003
)

type SerializedPayload_t struct {
	Encapsulation uint16
	Length        uint32
	data          []Octet
	MaxSize       uint32
	Pos           uint32 //!Position when reading
}

func CreateSerializedPayload() SerializedPayload_t {
	return SerializedPayload_t{
		Encapsulation: CDR_BE,
		Length:        0,
		MaxSize:       0,
		Pos:           0,
	}
}
