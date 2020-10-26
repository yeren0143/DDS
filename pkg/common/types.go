package rtps/common

import (
	"encoding/hex"
)

var (
	c_VendorId_Unknown  = NewVendorId("0x00", "0x00")
	c_VendorId_eProsima = NewVendorId("0x01", "0x0F")
)

type octet = uint8

type VendorId struct {
	m_vendor [2]octet
}

// func NewVendorId(v *VendorId) VendorId {
// 	vendor := VendorId{}
// 	vendor.m_vendor[0] = v.m_vendor[0]
// 	vendor.m_vendor[1] = v.m_vendor[1]
// 	return vendor
// }

func NewVendorId(str0 string, str1 string) VendorId {
	vendor := VendorId{}
	a, _ := hex.DecodeString(str0)
	b, _ := hex.DecodeString(str1)
	vendor.m_vendor[0] = a[0]
	vendor.m_vendor[1] = b[0]
	return vendor
}

func (vendor_id *VendorId) Equal(v *VendorId) bool {
	if vendor_id.m_vendor[0] == v.m_vendor[0] && vendor_id.m_vendor[1] == v.m_vendor[1] {
		return true
	} else {
		return false
	}
}
