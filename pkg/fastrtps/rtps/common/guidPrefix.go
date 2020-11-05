package common

type GuidPrefix_t struct {
	Value [12]Octet
}

var (
	UnknownGuid GuidPrefix_t = GuidPrefix_t{}
)

func NewGuiPrefix() *GuidPrefix_t {
	return &GuidPrefix_t{}
}
