package rtps

import "os"

type GuidPrefix struct {
	value [12]octet
}

var (
	unknownGuid = Guid{}
)

type Guid struct {
	guid_prefix GuidPrefix
	entity_id   EntityId
}

func (guid *Guid) IsOnTheSameHost(id *GuidPrefix) bool {
	return false
}

func (guid *Guid) IsOnTheSameProcess(id *GuidPrefix) bool {
	return false
}

func CreateGuidPrefix(id uint32) (guid_prefix GuidPrefix) {
	pid := os.Getppid()
	guid_prefix.value[0] = c_VendorId_eProsima.m_vendor[0]
	guid_prefix.value[1] = c_VendorId_eProsima.m_vendor[1]
}
