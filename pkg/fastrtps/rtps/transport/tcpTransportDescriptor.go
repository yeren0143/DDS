package transport

import (
	"log"
)

//TLSOptions ...
type TLSOptions = uint32

//TLSOptions Options
const (
	CTlsNone              TLSOptions = 0      // 0000 0000 0000
	CTlsDefaultWorkaround TLSOptions = 1 << 0 // 0000 0000 0001
	CTlsNoCompression     TLSOptions = 1 << 1 // 0000 0000 0010
	CTlsNoSSLv2           TLSOptions = 1 << 2 // 0000 0000 0100
	CTlsNoSSLv3           TLSOptions = 1 << 3 // 0000 0000 1000
	CTlsv1                TLSOptions = 1 << 4 // 0000 0001 0000
	CTlsv1_1              TLSOptions = 1 << 5 // 0000 0010 0000
	CTlsv1_2              TLSOptions = 1 << 6 // 0000 0100 0000
	Ctlsv1_3              TLSOptions = 1 << 7 // 0000 1000 0000
	CtlsSignleDhUse       TLSOptions = 1 << 8 // 0001 0000 0000
)

//TLSVerifyMode configuration
type TLSVerifyMode = uint8

//TLSVerifyMode
const (
	CTlsVerUnused          TLSVerifyMode = 0      // 0000 0000
	CTlsVerNone            TLSVerifyMode = 1 << 0 // 0000 0001
	CTlsVerPeer            TLSVerifyMode = 1 << 1 // 0000 0010
	CTlsFailedIfNoPeerCert               = 1 << 2 // 0000 0100
	CtlsClientOnce                       = 1 << 3 // 0000 1000
)

//TLSHandShakeRole configuration
type TLSHandShakeRole = uint8

//TLSHandShakeRole
const (
	CTlsRoleDefault TLSHandShakeRole = 0      // 0000 0000
	CTlsRoleClient  TLSHandShakeRole = 1 << 0 // 0000 0001
	CTlsRoleServer  TLSHandShakeRole = 1 << 1 // 0000 0010
)

//TLSConfigT configuration
type TLSConfigT struct {
	PassWord          string
	options           uint32
	CertChainFile     string
	PrivateKeyFile    string
	TmpDhFile         string
	VerifyFile        string
	VerifyMode        uint8
	VerifyPaths       []string
	DefaultVerifyPath bool
	VerifyDepth       int32
	RsaPrivateKeyFile string
	HandShakeRole     TLSHandShakeRole
}

//NewTLSConfig create default TLSConfig
func NewTLSConfig() *TLSConfigT {
	return &TLSConfigT{
		DefaultVerifyPath: false,
		VerifyDepth:       -1,
	}
}

var _ ITransportDescriptor = (*TCPTransportDescriptor)(nil)

type TCPTransportDescriptor struct {
	socketTransportDescriptor

	ListeningPorts        []uint16
	KeepAliveFrequency    uint32 //ms
	KeepAliveTimeOut      uint32
	MaxLogicalPort        uint16
	LogicalPortRange      uint16
	LogicalPortIncrement  uint16
	TCPNegotiationTimeOut uint32
	EnableTCPNodelay      bool
	WaitTCPNegotiation    bool
	CalculateCrc          bool
	CheckCrc              bool
	ApplySecurity         bool
	TLS                   *TLSConfigT
}

const (
	defaultKeepAliveFrequency     = 5000
	defaultKeepAliveTimeout       = 15000
	defaultTCPNegotitationTimeout = 5000
)

// MaxMessageSize test
func (descriptor *TCPTransportDescriptor) MaxMessageSize() uint32 {
	return descriptor.maxMessageSize
}

func (descriptor *TCPTransportDescriptor) CreateTransport() ITransport {
	log.Fatal("CreateTransport not impl")
	return nil
}

//NewTCPTransportDescriptor create default TCPTransportDescriptor
func NewTCPTransportDescriptor() *TCPTransportDescriptor {
	var socketDescriptor socketTransportDescriptor
	socketDescriptor.SendBufferSize = 0
	socketDescriptor.RcvBufferSize = 0
	socketDescriptor.maxMessageSize = KMaximumMessageSize
	socketDescriptor.maxInitialPeersRange = KMaximumInitialPeersRange
	socketDescriptor.TTL = CDefaultTTL

	return &TCPTransportDescriptor{
		socketTransportDescriptor: socketDescriptor,
		KeepAliveFrequency:        defaultKeepAliveFrequency,
		KeepAliveTimeOut:          defaultKeepAliveTimeout,
		MaxLogicalPort:            100,
		LogicalPortRange:          20,
		LogicalPortIncrement:      2,
		TCPNegotiationTimeOut:     defaultTCPNegotitationTimeout,
		EnableTCPNodelay:          false,
		WaitTCPNegotiation:        false,
		CalculateCrc:              true,
		CheckCrc:                  true,
		ApplySecurity:             false,
		TLS:                       NewTLSConfig(),
	}
}
