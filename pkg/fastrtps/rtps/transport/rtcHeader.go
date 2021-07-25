package transport

import (
	common "github.com/yeren0143/DDS/common"
)

type TCPHeader struct {
	rtpc   [4]common.Octet
	length uint32
	crc    uint32
	uint16 uint16
}

type TCPTransactionID struct {
	octets [12]common.Octet
}

type TCPCPMKind uint8

const (
	KBindConnectionRequest      TCPCPMKind = 0xD1
	KBindConnectionResponse     TCPCPMKind = 0xE1
	KOpenLogicalPortRequest     TCPCPMKind = 0xD2
	KOpenLogicalPortResponse    TCPCPMKind = 0xE2
	KCheckLogicalPortRequest    TCPCPMKind = 0xD3
	KCheckLogicalPortResponse   TCPCPMKind = 0xE3
	KKeepAliveRequest           TCPCPMKind = 0xD4
	KKeepAliveResponse          TCPCPMKind = 0xE4
	KLogicalPortIsClosedRequest TCPCPMKind = 0xD5
	KUnbindConnectionRequest    TCPCPMKind = 0xD6
)

type TCPControlMsgHeader struct {
	kind          TCPCPMKind
	flags         common.Octet
	length        uint16
	transactionID TCPTransactionID
}
