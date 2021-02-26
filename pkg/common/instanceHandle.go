package common

// InstanceHandleT used to contain the key for WITH_KEY topics.
type InstanceHandleT struct {
	Value [16]Octet
}

var (
	KInstanceHandleUnknown = InstanceHandleT{}
)

func (instance *InstanceHandleT) IsDefined() bool {
	for i := 0; i < 16; i++ {
		if instance.Value[i] != 0 {
			return true
		}
	}
	return false
}

func (instance *InstanceHandleT) InitWithGUID(guid *GUIDT) {
	for i:= 0; i < 16; i++ {
		if i < 12 {
			instance.Value[i] = guid.Prefix.Value[i]
		} else {
			instance.Value[i] = guid.EntityID.Value[i - 12]
		}
	}
}

func CreateInstanceHandle(guid GUIDT) InstanceHandleT {
	var instance InstanceHandleT
	for i := 0; i < 12; i++ {
		instance.Value[i] = guid.Prefix.Value[i]
	}
	for i := 12; i < 16; i++ {
		instance.Value[i] = guid.EntityID.Value[i-12]
	}

	return instance
}

func IHandle2GUID(ihandle *InstanceHandleT) GUIDT {
	var guid GUIDT
	for i := 0; i < 12; i++ {
		guid.Prefix.Value[i] = ihandle.Value[i]
	}
	for i := 12; i < 16; i++ {
		guid.EntityID.Value[i-12] = ihandle.Value[i]
	}
	return guid
}
