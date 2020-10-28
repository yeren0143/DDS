package common

var (
	c_VendorId_Unknown  = NewVendorId(0x00, 0x00)
	c_VendorId_eProsima = NewVendorId(0x01, 0x0F)
)

type Octet = byte

type VendorId struct {
	m_vendor [2]Octet
}

// func NewVendorId(v *VendorId) VendorId {
// 	vendor := VendorId{}
// 	vendor.m_vendor[0] = v.m_vendor[0]
// 	vendor.m_vendor[1] = v.m_vendor[1]
// 	return vendor
// }

func NewVendorId(a uint8, b uint8) VendorId {
	var vendor VendorId
	vendor.m_vendor[0] = a
	vendor.m_vendor[1] = b
	return vendor
}

func (vendor_id *VendorId) Equal(v *VendorId) bool {
	if vendor_id.m_vendor[0] == v.m_vendor[0] && vendor_id.m_vendor[1] == v.m_vendor[1] {
		return true
	} else {
		return false
	}
}
