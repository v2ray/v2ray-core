package handshake

import (
	"crypto/x509"
	"io"

	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/protocol"
	"v2ray.com/core/external/github.com/marten-seemann/qtls"
)

// Opener opens a packet
type Opener interface {
	Open(dst, src []byte, packetNumber protocol.PacketNumber, associatedData []byte) ([]byte, error)
	DecryptHeader(sample []byte, firstByte *byte, pnBytes []byte)
}

// Sealer seals a packet
type Sealer interface {
	Seal(dst, src []byte, packetNumber protocol.PacketNumber, associatedData []byte) []byte
	EncryptHeader(sample []byte, firstByte *byte, pnBytes []byte)
	Overhead() int
}

// A tlsExtensionHandler sends and received the QUIC TLS extension.
type tlsExtensionHandler interface {
	GetExtensions(msgType uint8) []qtls.Extension
	ReceivedExtensions(msgType uint8, exts []qtls.Extension) error
}

// CryptoSetup handles the handshake and protecting / unprotecting packets
type CryptoSetup interface {
	RunHandshake() error
	io.Closer

	HandleMessage([]byte, protocol.EncryptionLevel) bool
	ConnectionState() ConnectionState

	GetSealer() (protocol.EncryptionLevel, Sealer)
	GetSealerWithEncryptionLevel(protocol.EncryptionLevel) (Sealer, error)
	GetOpener(protocol.EncryptionLevel) (Opener, error)
}

// ConnectionState records basic details about the QUIC connection.
// Warning: This API should not be considered stable and might change soon.
type ConnectionState struct {
	HandshakeComplete bool                // handshake is complete
	ServerName        string              // server name requested by client, if any (server side only)
	PeerCertificates  []*x509.Certificate // certificate chain presented by remote peer
}
