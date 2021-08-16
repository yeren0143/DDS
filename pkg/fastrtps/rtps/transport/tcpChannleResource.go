package transport

import (
	"dds/common"
)

type SocketErrorCodes uint8

const (
	KNoError SocketErrorCodes = iota
	KBrokenPipe
	KSystemError
	KException
	KConnectionAborted = 125
)

type TCPConnectionType uint8

const (
	KTcpAcceptType  TCPConnectionType = 0
	KTcpConnectType TCPConnectionType = 1
)

type ConnectionStatus uint8

const (
	KDisconnected ConnectionStatus = iota
	KConnecting
	KConnected
	KWaitingForBind
	KWaitingForBindResponse
	KEstablished
	KUnbinding
)

type TCPChannelResource struct {
	ChannelResource
	locator             *common.Locator
	waitingForKeepAlive bool
}
