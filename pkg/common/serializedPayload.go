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

// Copy another structure (including allocating new space for the data.)
func (payload *SerializedPayloadT) Copy(serData *SerializedPayloadT, withLimit bool) bool {
	payload.Length = serData.Length
	if serData.Length > payload.MaxSize {
		if withLimit {
			return false
		}

		payload.Reserve(serData.Length)
	}
	payload.Encapsulation = serData.Encapsulation
	copy(payload.Data, serData.Data)
	return true
}

func (payload *SerializedPayloadT) Reserve(newSize uint32) {
	if newSize <= payload.MaxSize {
		return
	}

	if len(payload.Data) == 0 {
		payload.Data = make([]Octet, newSize)
	} else {
		newData := make([]Octet, newSize)
		copy(newData, payload.Data)
		payload.Data = newData
	}
	payload.MaxSize = newSize
}

func CreateSerializedPayload() SerializedPayloadT {
	return SerializedPayloadT{
		Encapsulation: CDR_BE,
		Length:        0,
		MaxSize:       0,
		Pos:           0,
	}
}
