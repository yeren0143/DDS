package common

type GUID_t struct {
	Prefix    GuidPrefix_t
	Entity_id EntityId
}

func (guid *GUID_t) IsOnTheSameHost(id *GuidPrefix_t) bool {
	return false
}

func (guid *GUID_t) IsOnTheSameProcess(id *GuidPrefix_t) bool {
	return false
}

// func CreateGuidPrefix(id uint32) *GuidPrefix {
// 	guid_prefix := GuidPrefix{}

// 	guid_prefix.value[0] = C_VendorId_eProsima.Vendor[0]
// 	guid_prefix.value[1] = C_VendorId_eProsima.Vendor[1]

// 	host_id := utils.GetHost().Id()
// 	guid_prefix.value[2] = Octet(host_id)
// 	guid_prefix.value[3] = Octet(host_id >> 8)

// 	pid := os.Getppid()
// 	guid_prefix.value[4] = Octet(pid)
// 	guid_prefix.value[5] = Octet(pid >> 8)
// 	guid_prefix.value[6] = Octet(pid >> 16)
// 	guid_prefix.value[7] = Octet(pid >> 24)
// 	guid_prefix.value[8] = Octet(id)
// 	guid_prefix.value[9] = Octet(id >> 8)
// 	guid_prefix.value[10] = Octet(id >> 16)
// 	guid_prefix.value[11] = Octet(id >> 24)

// 	return &guid_prefix
// }
