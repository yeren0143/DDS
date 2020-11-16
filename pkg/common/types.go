package common

// version info
var (
	CVendorIDUnknown   = NewVendorId(0x00, 0x00)
	CVendorIdeProsima  = NewVendorId(0x01, 0x0F)
	CProtocolVersion20 = ProtocolVersionT{2, 0}
	CProtocolVersion21 = ProtocolVersionT{2, 1}
	CProtocolVersion22 = ProtocolVersionT{2, 2}
	CProtocolVersion23 = ProtocolVersionT{2, 3}
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

type TopicKindT int8

const (
	CNoKey TopicKindT = iota
	CWithKey
)

type Octet = byte
type SubmessageFlag = byte
type BuiltinEndpointSet = uint32
type Count_t = uint32

type VendorId struct {
	Vendor [2]Octet
}

type ProtocolVersionT struct {
	Major Octet
	Minor Octet
}

func NewProtocolVersion(maj Octet, min Octet) ProtocolVersionT {
	return ProtocolVersionT{maj, min}
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
