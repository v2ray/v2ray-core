package internet

import (
	"crypto/tls"
	"net"
)

type ConnectionHandler func(Connection)

type Reusable interface {
	Reusable() bool
	SetReusable(reuse bool)
}

type StreamConnectionType int

const (
	StreamConnectionTypeRawTCP StreamConnectionType = 1
	StreamConnectionTypeTCP    StreamConnectionType = 2
	StreamConnectionTypeKCP    StreamConnectionType = 4
)

type StreamSecurityType int

const (
	StreamSecurityTypeNone StreamSecurityType = 0
	StreamSecurityTypeTLS  StreamSecurityType = 1
)

type TLSSettings struct {
	AllowInsecure bool
	Certs         []tls.Certificate
}

func (this *TLSSettings) GetTLSConfig() *tls.Config {
	config := &tls.Config{
		InsecureSkipVerify: this.AllowInsecure,
	}

	config.Certificates = this.Certs
	config.BuildNameToCertificate()

	return config
}

type StreamSettings struct {
	Type        StreamConnectionType
	Security    StreamSecurityType
	TLSSettings *TLSSettings
}

func (this *StreamSettings) IsCapableOf(streamType StreamConnectionType) bool {
	return (this.Type & streamType) == streamType
}

type Connection interface {
	net.Conn
	Reusable
}

type SysFd interface {
	SysFd() (int, error)
}
