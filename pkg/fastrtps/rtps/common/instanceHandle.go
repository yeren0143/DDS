package common

//Struct InstanceHandle_t, used to contain the key for WITH_KEY topics.
type InstanceHandle_t struct {
	Value [16]Octet
}

var (
	C_InstanceHandle_Unknown = InstanceHandle_t{}
)

func CreateInstanceHandle(guid GUID_t) InstanceHandle_t {
	var instance InstanceHandle_t
	for i := 0; i < 12; i += 1 {
		instance.Value[i] = guid.Prefix.Value[i]
	}
	for i := 12; i < 16; i += 1 {
		instance.Value[i] = guid.Entity_id.Value[i-12]
	}

	return instance
}

func IHandle2GUID(ihandle *InstanceHandle_t) GUID_t {
	var guid GUID_t
	for i := 0; i < 12; i += 1 {
		guid.Prefix.Value[i] = ihandle.Value[i]
	}
	for i := 12; i < 16; i += 1 {
		guid.Entity_id.Value[i-12] = ihandle.Value[i]
	}
	return guid
}
