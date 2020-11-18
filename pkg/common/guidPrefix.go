package common

//GUIDPrefixT guid prefix of GUID
type GUIDPrefixT struct {
	Value [12]Octet
}

//CUnknownGuidPrefix guid
var (
	CUnknownGUIDPrefix GUIDPrefixT = GUIDPrefixT{}
)

//NewGUIDPrefix ...
func NewGUIDPrefix() *GUIDPrefixT {
	return &GUIDPrefixT{}
}
