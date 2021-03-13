package common

var (
	KVendorIDTUnknown  VendorIDT
	KVendorIDTeProsima VendorIDT
	KProtocolVersion20 ProtocolVersionT
	KProtocolVersion21 ProtocolVersionT
	KProtocolVersion22 ProtocolVersionT
	KProtocolVersion23 ProtocolVersionT
	KProtocolVersion   ProtocolVersionT
)

// version info
func init() {
	KVendorIDTUnknown = CreateVendorID(0x00, 0x00)
	KVendorIDTeProsima = CreateVendorID(0x01, 0x0F)
	KProtocolVersion20 = ProtocolVersionT{2, 0}
	KProtocolVersion21 = ProtocolVersionT{2, 1}
	KProtocolVersion22 = ProtocolVersionT{2, 2}
	KProtocolVersion23 = ProtocolVersionT{2, 3}
	KProtocolVersion = ProtocolVersionT{2, 2}
}

type Endianness int8

const (
	BIGEND    Endianness = 0x1
	LITTLEEND Endianness = 0x0
)

var KDefaultEndian Endianness = LITTLEEND

type ReliabilityKindT int8

const (
	KReliable ReliabilityKindT = iota
	KBestEffort
)

type DurabilityKindT int8
const (
	KVolatile DurabilityKindT = iota
	KTransientLocal
	KTransient
	KPersistent //!< NOT IMPLEMENTED.
)

type EndpointKindT int8
const (
	KReader EndpointKindT = iota
	KWriter
)

type TopicKindT int8
const (
	KNoKey TopicKindT = iota
	KWithKey
)

type Octet = byte
type SubmessageFlag = byte
type BuiltinEndpointSet = uint32
type CountT = uint32

type VendorIDT struct {
	Vendor [2]Octet
}

type ProtocolVersionT struct {
	Major Octet
	Minor Octet
}

func NewProtocolVersion(maj Octet, min Octet) ProtocolVersionT {
	return ProtocolVersionT{maj, min}
}

func CreateVendorID(a uint8, b uint8) VendorIDT {
	var vendor VendorIDT
	vendor.Vendor[0] = a
	vendor.Vendor[1] = b
	return vendor
}

func (vendor_id *VendorIDT) Equal(v *VendorIDT) bool {
	if vendor_id.Vendor[0] == v.Vendor[0] && vendor_id.Vendor[1] == v.Vendor[1] {
		return true
	} else {
		return false
	}
}

func Bit(i uint32) byte {
	return 1 << i
}
