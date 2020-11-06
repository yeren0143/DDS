package types

import (
	. "common"
)

type TypeIdentifierKind = Octet

const (
	TI_STRING8_SMALL                Octet = 0x70
	TI_STRING8_LARGE                Octet = 0x71
	TI_STRING16_SMALL               Octet = 0x72
	TI_STRING16_LARGE               Octet = 0x73
	TI_PLAIN_SEQUENCE_SMALL         Octet = 0x80
	TI_PLAIN_SEQUENCE_LARGE         Octet = 0x81
	TI_PLAIN_ARRAY_SMALL            Octet = 0x90
	TI_PLAIN_ARRAY_LARGE            Octet = 0x91
	TI_PLAIN_MAP_SMALL              Octet = 0xA0
	TI_PLAIN_MAP_LARGE              Octet = 0xA1
	TI_STRONGLY_CONNECTED_COMPONENT Octet = 0xB0
)

//TODO:
type TypeIdentifier_t struct {
	Md Octet
}
