package common

var (
	C_VendorId_Unknown    = NewVendorId(0x00, 0x00)
	C_VendorId_eProsima   = NewVendorId(0x01, 0x0F)
	c_ProtocolVersion_2_0 = ProtocolVersion{2, 0}
	c_ProtocolVersion_2_1 = ProtocolVersion{2, 1}
	c_ProtocolVersion_2_2 = ProtocolVersion{2, 2}
	c_ProtocolVersion_2_3 = ProtocolVersion{2, 3}
)

type Endianness int8

const (
	BIGEND    Endianness = 0x1
	LITTLEEND Endianness = 0x0
)

type ReliabilityKind int8

const (
	RELIABLE ReliabilityKind = iota
	BEST_EFFORT
)

type DurabilityKind int8

const (
	VOLATILE DurabilityKind = iota
	TRANSIENT_LOCAL
	TRANSIENT
	PERSISTENT //!< NOT IMPLEMENTED.
)

type EndpointKind_t int8

const (
	READER EndpointKind_t = iota
	WRITER
)

type TopicKind_t int8

const (
	NO_KEY TopicKind_t = iota
	WITH_KEY
)

type Octet = byte
type SubmessageFlag = byte
type BuiltinEndpointSet = uint32
type Count_t = uint32

type VendorId struct {
	Vendor [2]Octet
}

type ProtocolVersion struct {
	Major Octet
	Minor Octet
}

func NewProtocolVersion(maj Octet, min Octet) ProtocolVersion {
	return ProtocolVersion{maj, min}
}

// func NewVendorId(v *VendorId) VendorId {
// 	vendor := VendorId{}
// 	vendor.m_vendor[0] = v.m_vendor[0]
// 	vendor.m_vendor[1] = v.m_vendor[1]
// 	return vendor
// }

func NewVendorId(a uint8, b uint8) VendorId {
	var vendor VendorId
	vendor.Vendor[0] = a
	vendor.Vendor[1] = b
	return vendor
}

func (vendor_id *VendorId) Equal(v *VendorId) bool {
	if vendor_id.Vendor[0] == v.Vendor[0] && vendor_id.Vendor[1] == v.Vendor[1] {
		return true
	} else {
		return false
	}
}
