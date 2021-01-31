package common

const (
	KGUIDPrefixSize = 12
)

// GUIDPrefixT guid prefix of GUID
type GUIDPrefixT struct {
	Value [12]Octet
}

//CUnknownGuidPrefix guid
var (
	KUnknownGUIDPrefix GUIDPrefixT
	KGuidPrefixUnknown GUIDPrefixT
)

func init() {
	KUnknownGUIDPrefix = *NewGUIDPrefix()
	KUnknownGUIDPrefix = *NewGUIDPrefix()
}

//NewGUIDPrefix ...
func NewGUIDPrefix() *GUIDPrefixT {
	return &GUIDPrefixT{}
}

func (guid *GUIDPrefixT) Equal(that *GUIDPrefixT) bool {
	return guid.Value == that.Value
}
