package types

import (
	"dds/common"
)

// OctetSeq is a array of bytes
type OctetSeq = []common.Octet

// predefined value
const (
	KTrue                       string = "true"
	KFalse                      string = "false"
	KAnnotationKeyID            string = "key"
	KAnnotationEPKeyID          string = "Key"
	KAnnotationTopicID          string = "Topic"
	KAnnotationExtensibilityID  string = "extensibility"
	KAnnotationFinalID          string = "final"
	KAnnotationAppendableID     string = "appendable"
	KAnnotationMutableID        string = "mutable"
	KAnnotationNestedID         string = "nested"
	KAnnotationOptionalID       string = "optional"
	KAnnotationMustUnderstandID string = "must_understand"
	KAnnotationNonSerializedID  string = "non_serialized"
	KAnnotationBitBoundID       string = "bit_bound"
	KAnnotationDefaultID        string = "default"
	KAnnotationDefaultLiteralID string = "default_literal"
	KAnnotationValueID          string = "value"
	KAnnotationPositionID       string = "position"

	KExtensibilityFinal      string = "FINAL"
	KExtensibilityAppendable string = "APPENDABLE"
	KExtensibilityMutable    string = "MUTABLE"

	KTKNameBoolean  string = "bool"
	KTKNameInt16    string = "int16_t"
	KTKNameUInt16   string = "uint16_t"
	KTKNameInt32    string = "int32_t"
	KTKNameUInt32   string = "uint32_t"
	KTKNameInt64    string = "int64_t"
	KTKNameUInt64   string = "uint64_t"
	KTKNameChar8    string = "char"
	KTKNameByte     string = "octet"
	KTKNameInt8     string = "int8_t"
	KTKNameUInt8    string = "uint8_t"
	KTKNameChar16   string = "wchar"
	KTKNameChar16T  string = "wchar_t"
	KTKNameFloat32  string = "float"
	KTKNameFloat64  string = "double"
	KTKNameFloat128 string = "longdouble"

	KTKNameString8    string = "string"
	KTKNameString16   string = "wstring"
	KTKNameAlias      string = "alias"
	KTKNameEnum       string = "enum"
	KTKNameBitMask    string = "bitmask"
	KTKNameAnnotation string = "annotation"
	KTKNameStruct     string = "struct"
	KTKNameUnion      string = "union"
	KTKNameBitSet     string = "bitset"
	KTKNameSequence   string = "sequence"
	KTKNameArray      string = "array"
	KTKNameMap        string = "map"
)

// EquivalenceKind is alias of octet
type EquivalenceKind common.Octet

// default values of EquivalenceKind
const (
	KEKMinimal  EquivalenceKind = 0xF1 // 0x1111 0001
	KEKComplete EquivalenceKind = 0xF2 // 0x1111 0010
	KEKBoth     EquivalenceKind = 0xF3 // 0x1111 0011
)

// TypeKind is alias of common.Octet
type TypeKind common.Octet

const (
	KTKNone     TypeKind = 0x00
	KTKBoolean  TypeKind = 0x01
	KTKByte     TypeKind = 0x02
	KTKInt16    TypeKind = 0x03
	KTKInt32    TypeKind = 0x04
	KTKInt64    TypeKind = 0x05
	KTKUInt16   TypeKind = 0x06
	KTKUInt32   TypeKind = 0x07
	KTKUInt64   TypeKind = 0x08
	KTKFloat32  TypeKind = 0x09
	KTKFloat64  TypeKind = 0x0A
	KTKFloat128 TypeKind = 0x0B
	KTKChar8    TypeKind = 0x10
	KTKChar16   TypeKind = 0x11

	KTKString8 TypeKind = 0x20
	KTString16 TypeKind = 0x21

	KTKAlias TypeKind = 0x30

	KTKEnum    TypeKind = 0x40
	KTKBitmask TypeKind = 0x41

	KTKAnnotation TypeKind = 0x50
	KTKStructure  TypeKind = 0x51
	KTKUnion      TypeKind = 0x52
	KTKBitSet     TypeKind = 0x53

	KTKSequence TypeKind = 0x60
	KTKArray    TypeKind = 0x61
	KTKMap      TypeKind = 0x62
)

// MemberName is The name of some element (e.g. type, type member, module)
type MemberName = string

// Valid characters are alphanumeric plus the "_" cannot start with digit
const (
	KMemberNameMaxLength int32 = 256
)

// QualifiedTypeName name includes the name of containing modules
type QualifiedTypeName string

// using "::" as separator. No leading "::". E.g. "MyModule::MyType"
const (
	KTypeNameMaxLength int32 = 256
)

// PrimitiveTypeID is a alias of bytes, Every type has an ID. Those of the primitive types are pre-defined.
type PrimitiveTypeID = common.Octet

// NameHash is First 4 bytes of MD5 of of a member name converted to bytes
// using UTF-8 encoding and without a 'nul' terminator.
// Example: the member name "color" has NameHash {0x70, 0xDD, 0xA5, 0xDF}
type NameHash [4]uint8

// Mask used to remove the flags that do no affect assignability
// Selects  T1, T2, O, M, K, D
const (
	KMemberFlagMinimalMask uint16 = 0x003f
)

// the enumeration ReturnCode_t.
const (
	KRetCodeOK                 uint32 = 0
	KRetCodeError              uint32 = 1
	KRetCodeUnsupported        uint32 = 2
	KRetCodeBadParameter       uint32 = 3
	KRetCodePreConditionNotMet uint32 = 4
	KRetCodeOutOfResources     uint32 = 5
	KRetCodeNotEnabled         uint32 = 6
	KRetCodeImmutablePolicy    uint32 = 7
	KRetCodeInconsistentPolicy uint32 = 8
	KRetCodeAlreadyDeleted     uint32 = 9
	KRetCodeTimeOut            uint32 = 10
	KRetCodeNoData             uint32 = 11
	KRetCodeIllegalOperation   uint32 = 12
)

type ReturnCodeT struct {
	Value uint32
}

func CreateReturnCode() ReturnCodeT {
	return ReturnCodeT{
		Value: KRetCodeOK,
	}
}

type ResponseCode = ReturnCodeT

type MemberID = uint32

const (
	KMemberIDInvalid  uint32 = 0x0FFFFFFF
	KIndexInvalid     uint32 = ^uint32(0)
	KMaxBitMaskLength int32  = 64
	KMaxElementsCount int32  = 100
	KMaxStringLength  int32  = 255
)

// Long Bound of a collection type
type LBound = uint32
type LBoundSeq = []LBound

// Short Bound of a collection type
type SBound = uint8
type SBoundSeq = []SBound

const (
	KInvalidLBound LBound = 0
	KInvalidSVound SBound = 0
)

// Flags that apply to struct/union/collection/enum/bitmask/bitset
// members/elements and DO affect type assignability
// Depending on the flag it may not apply to members of all types
type MemberFlag struct {
	memberFlag uint16
}

type CollectionElementFlag = MemberFlag
type StructMemberFlag = MemberFlag
type UnionMemberFlag = MemberFlag
type UnionDiscriminatorFlag = MemberFlag
type EnumeratedLiteralFlag = MemberFlag
type AnnotationParameterFlag = MemberFlag
type AliasMemberFlag = MemberFlag
type BitflagFlag = MemberFlag
type BitsetMemberFlag = MemberFlag

// Flags that apply to type declarationa and DO affect assignability
// Depending on the flag it may not apply to all types
// When not all, the applicable  types are listed
type TypeFlag struct {
	typeFlag uint16
}

type StructTypeFlag = TypeFlag
type UnionTypeFlag = TypeFlag
type CollectionTypeFlag = TypeFlag
type AnnotationTypeFlag = TypeFlag
type AliasTypeFlag = TypeFlag
type EnumTypeFlag = TypeFlag
type BitmaskTypeFlag = TypeFlag
type BitsetTypeFlag = TypeFlag

// Mask used to remove the flags that do no affect assignability
const (
	KTypeFlagMinimalMask uint16 = 0x0007
)

// ID of a type member
const (
	KAnnotationStrValueMaxLen      uint32 = 128
	KAnnotationOctetsecValueMaxLen uint32 = 128
)
