// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qtls

import (
	"container/list"
	"crypto"
	"crypto/rand"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	VersionSSL30 = 0x0300
	VersionTLS10 = 0x0301
	VersionTLS11 = 0x0302
	VersionTLS12 = 0x0303
	VersionTLS13 = 0x0304
)

const (
	maxPlaintext      = 16384        // maximum plaintext payload length
	maxCiphertext     = 16384 + 2048 // maximum ciphertext payload length
	recordHeaderLen   = 5            // record header length
	maxHandshake      = 65536        // maximum handshake we support (protocol max is 16 MB)
	maxWarnAlertCount = 5            // maximum number of consecutive warning alerts

	minVersion = VersionTLS12
	maxVersion = VersionTLS13
)

// TLS record types.
type recordType uint8

const (
	recordTypeChangeCipherSpec recordType = 20
	recordTypeAlert            recordType = 21
	recordTypeHandshake        recordType = 22
	recordTypeApplicationData  recordType = 23
)

// TLS handshake message types.
const (
	typeHelloRequest        uint8 = 0
	typeClientHello         uint8 = 1
	typeServerHello         uint8 = 2
	typeNewSessionTicket    uint8 = 4
	typeEndOfEarlyData      uint8 = 5
	typeEncryptedExtensions uint8 = 8
	typeCertificate         uint8 = 11
	typeServerKeyExchange   uint8 = 12
	typeCertificateRequest  uint8 = 13
	typeServerHelloDone     uint8 = 14
	typeCertificateVerify   uint8 = 15
	typeClientKeyExchange   uint8 = 16
	typeFinished            uint8 = 20
	typeCertificateStatus   uint8 = 22
	typeNextProtocol        uint8 = 67 // Not IANA assigned
)

// TLS compression types.
const (
	compressionNone uint8 = 0
)

type Extension struct {
	Type uint16
	Data []byte
}

// TLS extension numbers
const (
	extensionServerName              uint16 = 0
	extensionStatusRequest           uint16 = 5
	extensionSupportedCurves         uint16 = 10 // Supported Groups in 1.3 nomenclature
	extensionSupportedPoints         uint16 = 11
	extensionSignatureAlgorithms     uint16 = 13
	extensionALPN                    uint16 = 16
	extensionSCT                     uint16 = 18 // https://tools.ietf.org/html/rfc6962#section-6
	extensionEMS                     uint16 = 23
	extensionSessionTicket           uint16 = 35
	extensionPreSharedKey            uint16 = 41
	extensionEarlyData               uint16 = 42
	extensionSupportedVersions       uint16 = 43
	extensionPSKKeyExchangeModes     uint16 = 45
	extensionCAs                     uint16 = 47
	extensionSignatureAlgorithmsCert uint16 = 50
	extensionKeyShare                uint16 = 51
	extensionNextProtoNeg            uint16 = 13172 // not IANA assigned
	extensionRenegotiationInfo       uint16 = 0xff01
	extensionDelegatedCredential     uint16 = 0xff02 // TODO(any) Get IANA assignment
)

// TLS signaling cipher suite values
const (
	scsvRenegotiation uint16 = 0x00ff
)

// PSK Key Exchange Modes
// https://tools.ietf.org/html/draft-ietf-tls-tls13-18#section-4.2.7
const (
	pskDHEKeyExchange uint8 = 1
)

// CurveID is tls.CurveID
// TLS 1.3 refers to these as Groups, but this library implements only
// curve-based ones anyway. See https://tools.ietf.org/html/draft-ietf-tls-tls13-18#section-4.2.4.
type CurveID = tls.CurveID

const (
	// Exported IDs
	CurveP256 = tls.CurveP256
	CurveP384 = tls.CurveP384
	CurveP521 = tls.CurveP521
	X25519    = tls.X25519
)

// TLS 1.3 Key Share
// See https://tools.ietf.org/html/draft-ietf-tls-tls13-18#section-4.2.5
type keyShare struct {
	group CurveID
	data  []byte
}

// TLS 1.3 PSK Identity and Binder, as sent by the client
// https://tools.ietf.org/html/draft-ietf-tls-tls13-18#section-4.2.6

type psk struct {
	identity     []byte
	obfTicketAge uint32
	binder       []byte
}

// TLS Elliptic Curve Point Formats
// http://www.iana.org/assignments/tls-parameters/tls-parameters.xml#tls-parameters-9
const (
	pointFormatUncompressed uint8 = 0
)

// TLS CertificateStatusType (RFC 3546)
const (
	statusTypeOCSP uint8 = 1
)

// Certificate types (for certificateRequestMsg)
const (
	certTypeRSASign    = 1 // A certificate containing an RSA key
	certTypeDSSSign    = 2 // A certificate containing a DSA key
	certTypeRSAFixedDH = 3 // A certificate containing a static DH key
	certTypeDSSFixedDH = 4 // A certificate containing a static DH key

	// See RFC 4492 sections 3 and 5.5.
	certTypeECDSASign      = 64 // A certificate containing an ECDSA-capable public key, signed with ECDSA.
	certTypeRSAFixedECDH   = 65 // A certificate containing an ECDH-capable public key, signed with RSA.
	certTypeECDSAFixedECDH = 66 // A certificate containing an ECDH-capable public key, signed with ECDSA.

	// Rest of these are reserved by the TLS spec
)

// Signature algorithms for TLS 1.2 (See RFC 5246, section A.4.1)
const (
	signaturePKCS1v15 uint8 = iota + 1
	signatureECDSA
	signatureRSAPSS
)

// supportedSignatureAlgorithms contains the signature and hash algorithms that
// the code advertises as supported in a TLS 1.2 ClientHello and in a TLS 1.2
// CertificateRequest. The two fields are merged to match with TLS 1.3.
// Note that in TLS 1.2, the ECDSA algorithms are not constrained to P-256, etc.
var supportedSignatureAlgorithms = []SignatureScheme{
	PKCS1WithSHA256,
	ECDSAWithP256AndSHA256,
	PKCS1WithSHA384,
	ECDSAWithP384AndSHA384,
	PKCS1WithSHA512,
	ECDSAWithP521AndSHA512,
	PKCS1WithSHA1,
	ECDSAWithSHA1,
}

// supportedSignatureAlgorithms13 lists the advertised signature algorithms
// allowed for digital signatures. It includes TLS 1.2 + PSS.
var supportedSignatureAlgorithms13 = []SignatureScheme{
	PSSWithSHA256,
	PKCS1WithSHA256,
	ECDSAWithP256AndSHA256,
	PSSWithSHA384,
	PKCS1WithSHA384,
	ECDSAWithP384AndSHA384,
	PSSWithSHA512,
	PKCS1WithSHA512,
	ECDSAWithP521AndSHA512,
	PKCS1WithSHA1,
	ECDSAWithSHA1,
}

// ConnectionState records basic TLS details about the connection.
type ConnectionState struct {
	ConnectionID                []byte                // Random unique connection id
	Version                     uint16                // TLS version used by the connection (e.g. VersionTLS12)
	HandshakeComplete           bool                  // TLS handshake is complete
	DidResume                   bool                  // connection resumes a previous TLS connection
	CipherSuite                 uint16                // cipher suite in use (TLS_RSA_WITH_RC4_128_SHA, ...)
	NegotiatedProtocol          string                // negotiated next protocol (not guaranteed to be from Config.NextProtos)
	NegotiatedProtocolIsMutual  bool                  // negotiated protocol was advertised by server (client side only)
	ServerName                  string                // server name requested by client, if any (server side only)
	PeerCertificates            []*x509.Certificate   // certificate chain presented by remote peer
	VerifiedChains              [][]*x509.Certificate // verified chains built from PeerCertificates
	SignedCertificateTimestamps [][]byte              // SCTs from the server, if any
	OCSPResponse                []byte                // stapled OCSP response from server, if any
	DelegatedCredential         []byte                // Delegated credential sent by the server, if any

	// TLSUnique contains the "tls-unique" channel binding value (see RFC
	// 5929, section 3). For resumed sessions this value will be nil
	// because resumption does not include enough context (see
	// https://mitls.org/pages/attacks/3SHAKE#channelbindings). This will
	// change in future versions of Go once the TLS master-secret fix has
	// been standardized and implemented.
	TLSUnique []byte

	// HandshakeConfirmed is true once all data returned by Read
	// (past and future) is guaranteed not to be replayed.
	HandshakeConfirmed bool

	// Unique0RTTToken is a value that never repeats, and can be used
	// to detect replay attacks against 0-RTT connections.
	// Unique0RTTToken is only present if HandshakeConfirmed is false.
	Unique0RTTToken []byte

	ClientHello []byte // ClientHello packet
}

// The ClientAuthType is the tls.ClientAuthType
type ClientAuthType = tls.ClientAuthType

const (
	NoClientCert               = tls.NoClientCert
	RequestClientCert          = tls.RequestClientCert
	RequireAnyClientCert       = tls.RequireAnyClientCert
	VerifyClientCertIfGiven    = tls.VerifyClientCertIfGiven
	RequireAndVerifyClientCert = tls.RequireAndVerifyClientCert
)

// ClientSessionState contains the state needed by clients to resume TLS
// sessions.
type ClientSessionState struct {
	sessionTicket      []uint8               // Encrypted ticket used for session resumption with server
	vers               uint16                // SSL/TLS version negotiated for the session
	cipherSuite        uint16                // Ciphersuite negotiated for the session
	masterSecret       []byte                // MasterSecret generated by client on a full handshake
	serverCertificates []*x509.Certificate   // Certificate chain presented by the server
	verifiedChains     [][]*x509.Certificate // Certificate chains we built for verification
	useEMS             bool                  // State of extended master secret
}

// ClientSessionCache is a cache of ClientSessionState objects that can be used
// by a client to resume a TLS session with a given server. ClientSessionCache
// implementations should expect to be called concurrently from different
// goroutines. Only ticket-based resumption is supported, not SessionID-based
// resumption.
type ClientSessionCache interface {
	// Get searches for a ClientSessionState associated with the given key.
	// On return, ok is true if one was found.
	Get(sessionKey string) (session *ClientSessionState, ok bool)

	// Put adds the ClientSessionState to the cache with the given key.
	Put(sessionKey string, cs *ClientSessionState)
}

// SignatureScheme is a tls.SignatureScheme
type SignatureScheme = tls.SignatureScheme

const (
	PKCS1WithSHA1   = tls.PKCS1WithSHA1
	PKCS1WithSHA256 = tls.PKCS1WithSHA256
	PKCS1WithSHA384 = tls.PKCS1WithSHA384
	PKCS1WithSHA512 = tls.PKCS1WithSHA512

	PSSWithSHA256 = tls.PSSWithSHA256
	PSSWithSHA384 = tls.PSSWithSHA384
	PSSWithSHA512 = tls.PSSWithSHA512

	ECDSAWithP256AndSHA256 = tls.ECDSAWithP256AndSHA256
	ECDSAWithP384AndSHA384 = tls.ECDSAWithP384AndSHA384
	ECDSAWithP521AndSHA512 = tls.ECDSAWithP521AndSHA512

	// Legacy signature and hash algorithms for TLS 1.2.
	ECDSAWithSHA1 = tls.ECDSAWithSHA1
)

// ClientHelloInfo contains information from a ClientHello message in order to
// guide certificate selection in the GetCertificate callback.
type ClientHelloInfo struct {
	// CipherSuites lists the CipherSuites supported by the client (e.g.
	// TLS_RSA_WITH_RC4_128_SHA).
	CipherSuites []uint16

	// ServerName indicates the name of the server requested by the client
	// in order to support virtual hosting. ServerName is only set if the
	// client is using SNI (see
	// http://tools.ietf.org/html/rfc4366#section-3.1).
	ServerName string

	// SupportedCurves lists the elliptic curves supported by the client.
	// SupportedCurves is set only if the Supported Elliptic Curves
	// Extension is being used (see
	// http://tools.ietf.org/html/rfc4492#section-5.1.1).
	SupportedCurves []CurveID

	// SupportedPoints lists the point formats supported by the client.
	// SupportedPoints is set only if the Supported Point Formats Extension
	// is being used (see
	// http://tools.ietf.org/html/rfc4492#section-5.1.2).
	SupportedPoints []uint8

	// SignatureSchemes lists the signature and hash schemes that the client
	// is willing to verify. SignatureSchemes is set only if the Signature
	// Algorithms Extension is being used (see
	// https://tools.ietf.org/html/rfc5246#section-7.4.1.4.1).
	SignatureSchemes []SignatureScheme

	// SupportedProtos lists the application protocols supported by the client.
	// SupportedProtos is set only if the Application-Layer Protocol
	// Negotiation Extension is being used (see
	// https://tools.ietf.org/html/rfc7301#section-3.1).
	//
	// Servers can select a protocol by setting Config.NextProtos in a
	// GetConfigForClient return value.
	SupportedProtos []string

	// SupportedVersions lists the TLS versions supported by the client.
	// For TLS versions less than 1.3, this is extrapolated from the max
	// version advertised by the client, so values other than the greatest
	// might be rejected if used.
	SupportedVersions []uint16

	// Conn is the underlying net.Conn for the connection. Do not read
	// from, or write to, this connection; that will cause the TLS
	// connection to fail.
	Conn net.Conn

	// Offered0RTTData is true if the client announced that it will send
	// 0-RTT data. If the server Config.Accept0RTTData is true, and the
	// client offered a session ticket valid for that purpose, it will
	// be notified that the 0-RTT data is accepted and it will be made
	// immediately available for Read.
	Offered0RTTData bool

	// AcceptsDelegatedCredential is true if the client indicated willingness
	// to negotiate the delegated credential extension.
	AcceptsDelegatedCredential bool

	// The Fingerprint is an sequence of bytes unique to this Client Hello.
	// It can be used to prevent or mitigate 0-RTT data replays as it's
	// guaranteed that a replayed connection will have the same Fingerprint.
	Fingerprint []byte
}

// The CertificateRequestInfo is a tls.CertificateRequestInfo
type CertificateRequestInfo = tls.CertificateRequestInfo

// RenegotiationSupport is a tls.RenegotiationSupport
type RenegotiationSupport = tls.RenegotiationSupport

const (
	// RenegotiateNever disables renegotiation.
	RenegotiateNever = tls.RenegotiateNever

	// RenegotiateOnceAsClient allows a remote server to request
	// renegotiation once per connection.
	RenegotiateOnceAsClient = tls.RenegotiateOnceAsClient

	// RenegotiateFreelyAsClient allows a remote server to repeatedly
	// request renegotiation.
	RenegotiateFreelyAsClient = tls.RenegotiateFreelyAsClient
)

// A Config structure is used to configure a TLS client or server.
// After one has been passed to a TLS function it must not be
// modified. A Config may be reused; the tls package will also not
// modify it.
type Config struct {
	// Rand provides the source of entropy for nonces and RSA blinding.
	// If Rand is nil, TLS uses the cryptographic random reader in package
	// crypto/rand.
	// The Reader must be safe for use by multiple goroutines.
	Rand io.Reader

	// Time returns the current time as the number of seconds since the epoch.
	// If Time is nil, TLS uses time.Now.
	Time func() time.Time

	// Certificates contains one or more certificate chains to present to
	// the other side of the connection. Server configurations must include
	// at least one certificate or else set GetCertificate. Clients doing
	// client-authentication may set either Certificates or
	// GetClientCertificate.
	Certificates []Certificate

	// NameToCertificate maps from a certificate name to an element of
	// Certificates. Note that a certificate name can be of the form
	// '*.example.com' and so doesn't have to be a domain name as such.
	// See Config.BuildNameToCertificate
	// The nil value causes the first element of Certificates to be used
	// for all connections.
	NameToCertificate map[string]*Certificate

	// GetCertificate returns a Certificate based on the given
	// ClientHelloInfo. It will only be called if the client supplies SNI
	// information or if Certificates is empty.
	//
	// If GetCertificate is nil or returns nil, then the certificate is
	// retrieved from NameToCertificate. If NameToCertificate is nil, the
	// first element of Certificates will be used.
	GetCertificate func(*ClientHelloInfo) (*Certificate, error)

	// GetClientCertificate, if not nil, is called when a server requests a
	// certificate from a client. If set, the contents of Certificates will
	// be ignored.
	//
	// If GetClientCertificate returns an error, the handshake will be
	// aborted and that error will be returned. Otherwise
	// GetClientCertificate must return a non-nil Certificate. If
	// Certificate.Certificate is empty then no certificate will be sent to
	// the server. If this is unacceptable to the server then it may abort
	// the handshake.
	//
	// GetClientCertificate may be called multiple times for the same
	// connection if renegotiation occurs or if TLS 1.3 is in use.
	GetClientCertificate func(*CertificateRequestInfo) (*Certificate, error)

	// GetConfigForClient, if not nil, is called after a ClientHello is
	// received from a client. It may return a non-nil Config in order to
	// change the Config that will be used to handle this connection. If
	// the returned Config is nil, the original Config will be used. The
	// Config returned by this callback may not be subsequently modified.
	//
	// If GetConfigForClient is nil, the Config passed to Server() will be
	// used for all connections.
	//
	// Uniquely for the fields in the returned Config, session ticket keys
	// will be duplicated from the original Config if not set.
	// Specifically, if SetSessionTicketKeys was called on the original
	// config but not on the returned config then the ticket keys from the
	// original config will be copied into the new config before use.
	// Otherwise, if SessionTicketKey was set in the original config but
	// not in the returned config then it will be copied into the returned
	// config before use. If neither of those cases applies then the key
	// material from the returned config will be used for session tickets.
	GetConfigForClient func(*ClientHelloInfo) (*Config, error)

	// VerifyPeerCertificate, if not nil, is called after normal
	// certificate verification by either a TLS client or server. It
	// receives the raw ASN.1 certificates provided by the peer and also
	// any verified chains that normal processing found. If it returns a
	// non-nil error, the handshake is aborted and that error results.
	//
	// If normal verification fails then the handshake will abort before
	// considering this callback. If normal verification is disabled by
	// setting InsecureSkipVerify, or (for a server) when ClientAuth is
	// RequestClientCert or RequireAnyClientCert, then this callback will
	// be considered but the verifiedChains argument will always be nil.
	VerifyPeerCertificate func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error

	// RootCAs defines the set of root certificate authorities
	// that clients use when verifying server certificates.
	// If RootCAs is nil, TLS uses the host's root CA set.
	RootCAs *x509.CertPool

	// NextProtos is a list of supported, application level protocols.
	NextProtos []string

	// ServerName is used to verify the hostname on the returned
	// certificates unless InsecureSkipVerify is given. It is also included
	// in the client's handshake to support virtual hosting unless it is
	// an IP address.
	ServerName string

	// ClientAuth determines the server's policy for
	// TLS Client Authentication. The default is NoClientCert.
	ClientAuth ClientAuthType

	// ClientCAs defines the set of root certificate authorities
	// that servers use if required to verify a client certificate
	// by the policy in ClientAuth.
	ClientCAs *x509.CertPool

	// InsecureSkipVerify controls whether a client verifies the
	// server's certificate chain and host name.
	// If InsecureSkipVerify is true, TLS accepts any certificate
	// presented by the server and any host name in that certificate.
	// In this mode, TLS is susceptible to man-in-the-middle attacks.
	// This should be used only for testing.
	InsecureSkipVerify bool

	// CipherSuites is a list of supported cipher suites to be used in
	// TLS 1.0-1.2. If CipherSuites is nil, TLS uses a list of suites
	// supported by the implementation.
	CipherSuites []uint16

	// PreferServerCipherSuites controls whether the server selects the
	// client's most preferred ciphersuite, or the server's most preferred
	// ciphersuite. If true then the server's preference, as expressed in
	// the order of elements in CipherSuites, is used.
	PreferServerCipherSuites bool

	// SessionTicketsDisabled may be set to true to disable session ticket
	// (resumption) support.
	SessionTicketsDisabled bool

	// SessionTicketKey is used by TLS servers to provide session
	// resumption. See RFC 5077. If zero, it will be filled with
	// random data before the first server handshake.
	//
	// If multiple servers are terminating connections for the same host
	// they should all have the same SessionTicketKey. If the
	// SessionTicketKey leaks, previously recorded and future TLS
	// connections using that key are compromised.
	SessionTicketKey [32]byte

	// ClientSessionCache is a cache of ClientSessionState entries for TLS
	// session resumption.
	ClientSessionCache ClientSessionCache

	// MinVersion contains the minimum SSL/TLS version that is acceptable.
	// If zero, then TLS 1.0 is taken as the minimum.
	MinVersion uint16

	// MaxVersion contains the maximum SSL/TLS version that is acceptable.
	// If zero, then the maximum version supported by this package is used,
	// which is currently TLS 1.2.
	MaxVersion uint16

	// CurvePreferences contains the elliptic curves that will be used in
	// an ECDHE handshake, in preference order. If empty, the default will
	// be used.
	CurvePreferences []CurveID

	// DynamicRecordSizingDisabled disables adaptive sizing of TLS records.
	// When true, the largest possible TLS record size is always used. When
	// false, the size of TLS records may be adjusted in an attempt to
	// improve latency.
	DynamicRecordSizingDisabled bool

	// Renegotiation controls what types of renegotiation are supported.
	// The default, none, is correct for the vast majority of applications.
	Renegotiation RenegotiationSupport

	// KeyLogWriter optionally specifies a destination for TLS master secrets
	// in NSS key log format that can be used to allow external programs
	// such as Wireshark to decrypt TLS connections.
	// See https://developer.mozilla.org/en-US/docs/Mozilla/Projects/NSS/Key_Log_Format.
	// Use of KeyLogWriter compromises security and should only be
	// used for debugging.
	KeyLogWriter io.Writer

	// If Max0RTTDataSize is not zero, the client will be allowed to use
	// session tickets to send at most this number of bytes of 0-RTT data.
	// 0-RTT data is subject to replay and has memory DoS implications.
	// The server will later be able to refuse the 0-RTT data with
	// Accept0RTTData, or wait for the client to prove that it's not
	// replayed with Conn.ConfirmHandshake.
	//
	// It has no meaning on the client.
	//
	// See https://tools.ietf.org/html/draft-ietf-tls-tls13-18#section-2.3.
	Max0RTTDataSize uint32

	// Accept0RTTData makes the 0-RTT data received from the client
	// immediately available to Read. 0-RTT data is subject to replay.
	// Use Conn.ConfirmHandshake to wait until the data is known not
	// to be replayed after reading it.
	//
	// It has no meaning on the client.
	//
	// See https://tools.ietf.org/html/draft-ietf-tls-tls13-18#section-2.3.
	Accept0RTTData bool

	// SessionTicketSealer, if not nil, is used to wrap and unwrap
	// session tickets, instead of SessionTicketKey.
	SessionTicketSealer SessionTicketSealer

	// AcceptDelegatedCredential is true if the client is willing to negotiate
	// the delegated credential extension.
	//
	// This value has no meaning for the server.
	//
	// See https://tools.ietf.org/html/draft-ietf-tls-subcerts-02.
	AcceptDelegatedCredential bool

	// GetDelegatedCredential returns a DC and its private key for use in the
	// delegated credential extension. The inputs to the callback are some
	// information parsed from the ClientHello, as well as the protocol version
	// selected by the server. This is necessary because the DC is bound to the
	// protocol version in which it's used. The return value is the raw DC
	// encoded in the wire format specified in
	// https://tools.ietf.org/html/draft-ietf-tls-subcerts-02. If the return
	// value is nil, then the server will not offer negotiate the extension.
	//
	// This value has no meaning for the client.
	GetDelegatedCredential func(*ClientHelloInfo, uint16) ([]byte, crypto.PrivateKey, error)

	// GetExtensions, if not nil, is called before a message that allows
	// sending of extensions is sent.
	// Currently only implemented for the ClientHello message (for the client)
	// and for the EncryptedExtensions message (for the server).
	// Only valid for TLS 1.3.
	GetExtensions func(handshakeMessageType uint8) []Extension

	// ReceivedExtensions, if not nil, is called when a message that allows the
	// inclusion of extensions is received.
	// It is called with an empty slice of extensions, if the message didn't
	// contain any extensions.
	// Currently only implemented for the ClientHello message (sent by the
	// client) and for the EncryptedExtensions message (sent by the server).
	// Only valid for TLS 1.3.
	ReceivedExtensions func(handshakeMessageType uint8, exts []Extension) error

	serverInitOnce sync.Once // guards calling (*Config).serverInit

	// mutex protects sessionTicketKeys.
	mutex sync.RWMutex
	// sessionTicketKeys contains zero or more ticket keys. If the length
	// is zero, SessionTicketsDisabled must be true. The first key is used
	// for new tickets and any subsequent keys can be used to decrypt old
	// tickets.
	sessionTicketKeys []ticketKey

	// UseExtendedMasterSecret indicates whether or not the connection
	// should use the extended master secret computation if available
	UseExtendedMasterSecret bool

	// AlternativeRecordLayer is used by QUIC
	AlternativeRecordLayer RecordLayer
}

type RecordLayer interface {
	SetReadKey(suite *CipherSuite, trafficSecret []byte)
	SetWriteKey(suite *CipherSuite, trafficSecret []byte)
	ReadHandshakeMessage() ([]byte, error)
	WriteRecord([]byte) (int, error)
}

// ticketKeyNameLen is the number of bytes of identifier that is prepended to
// an encrypted session ticket in order to identify the key used to encrypt it.
const ticketKeyNameLen = 16

// ticketKey is the internal representation of a session ticket key.
type ticketKey struct {
	// keyName is an opaque byte string that serves to identify the session
	// ticket key. It's exposed as plaintext in every session ticket.
	keyName [ticketKeyNameLen]byte
	aesKey  [16]byte
	hmacKey [16]byte
}

// ticketKeyFromBytes converts from the external representation of a session
// ticket key to a ticketKey. Externally, session ticket keys are 32 random
// bytes and this function expands that into sufficient name and key material.
func ticketKeyFromBytes(b [32]byte) (key ticketKey) {
	hashed := sha512.Sum512(b[:])
	copy(key.keyName[:], hashed[:ticketKeyNameLen])
	copy(key.aesKey[:], hashed[ticketKeyNameLen:ticketKeyNameLen+16])
	copy(key.hmacKey[:], hashed[ticketKeyNameLen+16:ticketKeyNameLen+32])
	return key
}

// Clone returns a shallow clone of c. It is safe to clone a Config that is
// being used concurrently by a TLS client or server.
func (c *Config) Clone() *Config {
	// Running serverInit ensures that it's safe to read
	// SessionTicketsDisabled.
	c.serverInitOnce.Do(func() { c.serverInit(nil) })

	var sessionTicketKeys []ticketKey
	c.mutex.RLock()
	sessionTicketKeys = c.sessionTicketKeys
	c.mutex.RUnlock()

	return &Config{
		Rand:                        c.Rand,
		Time:                        c.Time,
		Certificates:                c.Certificates,
		NameToCertificate:           c.NameToCertificate,
		GetCertificate:              c.GetCertificate,
		GetClientCertificate:        c.GetClientCertificate,
		GetConfigForClient:          c.GetConfigForClient,
		VerifyPeerCertificate:       c.VerifyPeerCertificate,
		RootCAs:                     c.RootCAs,
		NextProtos:                  c.NextProtos,
		ServerName:                  c.ServerName,
		ClientAuth:                  c.ClientAuth,
		ClientCAs:                   c.ClientCAs,
		InsecureSkipVerify:          c.InsecureSkipVerify,
		CipherSuites:                c.CipherSuites,
		PreferServerCipherSuites:    c.PreferServerCipherSuites,
		SessionTicketsDisabled:      c.SessionTicketsDisabled,
		SessionTicketKey:            c.SessionTicketKey,
		ClientSessionCache:          c.ClientSessionCache,
		MinVersion:                  c.MinVersion,
		MaxVersion:                  c.MaxVersion,
		CurvePreferences:            c.CurvePreferences,
		DynamicRecordSizingDisabled: c.DynamicRecordSizingDisabled,
		Renegotiation:               c.Renegotiation,
		KeyLogWriter:                c.KeyLogWriter,
		Accept0RTTData:              c.Accept0RTTData,
		Max0RTTDataSize:             c.Max0RTTDataSize,
		SessionTicketSealer:         c.SessionTicketSealer,
		AcceptDelegatedCredential:   c.AcceptDelegatedCredential,
		GetDelegatedCredential:      c.GetDelegatedCredential,
		GetExtensions:               c.GetExtensions,
		ReceivedExtensions:          c.ReceivedExtensions,
		sessionTicketKeys:           sessionTicketKeys,
		UseExtendedMasterSecret:     c.UseExtendedMasterSecret,
	}
}

// serverInit is run under c.serverInitOnce to do initialization of c. If c was
// returned by a GetConfigForClient callback then the argument should be the
// Config that was passed to Server, otherwise it should be nil.
func (c *Config) serverInit(originalConfig *Config) {
	if c.SessionTicketsDisabled || len(c.ticketKeys()) != 0 || c.SessionTicketSealer != nil {
		return
	}

	alreadySet := false
	for _, b := range c.SessionTicketKey {
		if b != 0 {
			alreadySet = true
			break
		}
	}

	if !alreadySet {
		if originalConfig != nil {
			copy(c.SessionTicketKey[:], originalConfig.SessionTicketKey[:])
		} else if _, err := io.ReadFull(c.rand(), c.SessionTicketKey[:]); err != nil {
			c.SessionTicketsDisabled = true
			return
		}
	}

	if originalConfig != nil {
		originalConfig.mutex.RLock()
		c.sessionTicketKeys = originalConfig.sessionTicketKeys
		originalConfig.mutex.RUnlock()
	} else {
		c.sessionTicketKeys = []ticketKey{ticketKeyFromBytes(c.SessionTicketKey)}
	}
}

func (c *Config) ticketKeys() []ticketKey {
	c.mutex.RLock()
	// c.sessionTicketKeys is constant once created. SetSessionTicketKeys
	// will only update it by replacing it with a new value.
	ret := c.sessionTicketKeys
	c.mutex.RUnlock()
	return ret
}

// SetSessionTicketKeys updates the session ticket keys for a server. The first
// key will be used when creating new tickets, while all keys can be used for
// decrypting tickets. It is safe to call this function while the server is
// running in order to rotate the session ticket keys. The function will panic
// if keys is empty.
func (c *Config) SetSessionTicketKeys(keys [][32]byte) {
	if len(keys) == 0 {
		panic("tls: keys must have at least one key")
	}

	newKeys := make([]ticketKey, len(keys))
	for i, bytes := range keys {
		newKeys[i] = ticketKeyFromBytes(bytes)
	}

	c.mutex.Lock()
	c.sessionTicketKeys = newKeys
	c.mutex.Unlock()
}

func (c *Config) rand() io.Reader {
	r := c.Rand
	if r == nil {
		return rand.Reader
	}
	return r
}

func (c *Config) time() time.Time {
	t := c.Time
	if t == nil {
		t = time.Now
	}
	return t()
}

func hasOverlappingCipherSuites(cs1, cs2 []uint16) bool {
	for _, c1 := range cs1 {
		for _, c2 := range cs2 {
			if c1 == c2 {
				return true
			}
		}
	}
	return false
}

func (c *Config) cipherSuites() []uint16 {
	s := c.CipherSuites
	if s == nil {
		s = defaultCipherSuites()
	} else if c.maxVersion() >= VersionTLS13 {
		// Ensure that TLS 1.3 suites are always present, but respect
		// the application cipher suite preferences.
		s13 := defaultTLS13CipherSuites()
		if !hasOverlappingCipherSuites(s, s13) {
			allSuites := make([]uint16, len(s13)+len(s))
			allSuites = append(allSuites, s13...)
			s = append(allSuites, s...)
		}
	}
	return s
}

func (c *Config) minVersion() uint16 {
	if c == nil || c.MinVersion == 0 {
		return minVersion
	}
	return c.MinVersion
}

func (c *Config) maxVersion() uint16 {
	if c == nil || c.MaxVersion == 0 {
		return maxVersion
	}
	return c.MaxVersion
}

var defaultCurvePreferences = []CurveID{X25519, CurveP256, CurveP384, CurveP521}

func (c *Config) curvePreferences() []CurveID {
	if c == nil || len(c.CurvePreferences) == 0 {
		return defaultCurvePreferences
	}
	return c.CurvePreferences
}

// mutualVersion returns the protocol version to use given the advertised
// version of the peer using the legacy non-extension methods.
func (c *Config) mutualVersion(vers uint16) (uint16, bool) {
	minVersion := c.minVersion()
	maxVersion := c.maxVersion()

	// Version 1.3 and higher are not negotiated via this mechanism.
	if maxVersion > VersionTLS12 {
		maxVersion = VersionTLS12
	}

	if vers < minVersion {
		return 0, false
	}
	if vers > maxVersion {
		vers = maxVersion
	}
	return vers, true
}

// pickVersion returns the protocol version to use given the advertised
// versions of the peer using the Supported Versions extension.
func (c *Config) pickVersion(peerSupportedVersions []uint16) (uint16, bool) {
	supportedVersions := c.getSupportedVersions()
	for _, supportedVersion := range supportedVersions {
		for _, version := range peerSupportedVersions {
			if version == supportedVersion {
				return version, true
			}
		}
	}
	return 0, false
}

// configSuppVersArray is the backing array of Config.getSupportedVersions
var configSuppVersArray = [...]uint16{VersionTLS13, VersionTLS12, VersionTLS11, VersionTLS10, VersionSSL30}

// getSupportedVersions returns the protocol versions that are supported by the
// current configuration.
func (c *Config) getSupportedVersions() []uint16 {
	minVersion := c.minVersion()
	maxVersion := c.maxVersion()
	// Sanity check to avoid advertising unsupported versions.
	if minVersion < VersionSSL30 {
		minVersion = VersionSSL30
	}
	if maxVersion > VersionTLS13 {
		maxVersion = VersionTLS13
	}
	if maxVersion < minVersion {
		return nil
	}
	return configSuppVersArray[VersionTLS13-maxVersion : VersionTLS13-minVersion+1]
}

// getCertificate returns the best certificate for the given ClientHelloInfo,
// defaulting to the first element of c.Certificates.
func (c *Config) getCertificate(clientHello *ClientHelloInfo) (*Certificate, error) {
	if c.GetCertificate != nil &&
		(len(c.Certificates) == 0 || len(clientHello.ServerName) > 0) {
		cert, err := c.GetCertificate(clientHello)
		if cert != nil || err != nil {
			return cert, err
		}
	}

	if len(c.Certificates) == 0 {
		return nil, errors.New("tls: no certificates configured")
	}

	if len(c.Certificates) == 1 || c.NameToCertificate == nil {
		// There's only one choice, so no point doing any work.
		return &c.Certificates[0], nil
	}

	name := strings.ToLower(clientHello.ServerName)
	for len(name) > 0 && name[len(name)-1] == '.' {
		name = name[:len(name)-1]
	}

	if cert, ok := c.NameToCertificate[name]; ok {
		return cert, nil
	}

	// try replacing labels in the name with wildcards until we get a
	// match.
	labels := strings.Split(name, ".")
	for i := range labels {
		labels[i] = "*"
		candidate := strings.Join(labels, ".")
		if cert, ok := c.NameToCertificate[candidate]; ok {
			return cert, nil
		}
	}

	// If nothing matches, return the first certificate.
	return &c.Certificates[0], nil
}

// BuildNameToCertificate parses c.Certificates and builds c.NameToCertificate
// from the CommonName and SubjectAlternateName fields of each of the leaf
// certificates.
func (c *Config) BuildNameToCertificate() {
	c.NameToCertificate = make(map[string]*Certificate)
	for i := range c.Certificates {
		cert := &c.Certificates[i]
		x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
		if err != nil {
			continue
		}
		if len(x509Cert.Subject.CommonName) > 0 {
			c.NameToCertificate[x509Cert.Subject.CommonName] = cert
		}
		for _, san := range x509Cert.DNSNames {
			c.NameToCertificate[san] = cert
		}
	}
}

// writeKeyLog logs client random and master secret if logging was enabled by
// setting c.KeyLogWriter.
func (c *Config) writeKeyLog(what string, clientRandom, masterSecret []byte) error {
	if c.KeyLogWriter == nil {
		return nil
	}

	logLine := []byte(fmt.Sprintf("%s %x %x\n", what, clientRandom, masterSecret))

	writerMutex.Lock()
	_, err := c.KeyLogWriter.Write(logLine)
	writerMutex.Unlock()

	return err
}

// writerMutex protects all KeyLogWriters globally. It is rarely enabled,
// and is only for debugging, so a global mutex saves space.
var writerMutex sync.Mutex

// A Certificate is a tls.Certificate
type Certificate = tls.Certificate

type handshakeMessage interface {
	marshal() []byte
	unmarshal([]byte) alert
}

// lruSessionCache is a ClientSessionCache implementation that uses an LRU
// caching strategy.
type lruSessionCache struct {
	sync.Mutex

	m        map[string]*list.Element
	q        *list.List
	capacity int
}

type lruSessionCacheEntry struct {
	sessionKey string
	state      *ClientSessionState
}

// NewLRUClientSessionCache returns a ClientSessionCache with the given
// capacity that uses an LRU strategy. If capacity is < 1, a default capacity
// is used instead.
func NewLRUClientSessionCache(capacity int) ClientSessionCache {
	const defaultSessionCacheCapacity = 64

	if capacity < 1 {
		capacity = defaultSessionCacheCapacity
	}
	return &lruSessionCache{
		m:        make(map[string]*list.Element),
		q:        list.New(),
		capacity: capacity,
	}
}

// Put adds the provided (sessionKey, cs) pair to the cache.
func (c *lruSessionCache) Put(sessionKey string, cs *ClientSessionState) {
	c.Lock()
	defer c.Unlock()

	if elem, ok := c.m[sessionKey]; ok {
		entry := elem.Value.(*lruSessionCacheEntry)
		entry.state = cs
		c.q.MoveToFront(elem)
		return
	}

	if c.q.Len() < c.capacity {
		entry := &lruSessionCacheEntry{sessionKey, cs}
		c.m[sessionKey] = c.q.PushFront(entry)
		return
	}

	elem := c.q.Back()
	entry := elem.Value.(*lruSessionCacheEntry)
	delete(c.m, entry.sessionKey)
	entry.sessionKey = sessionKey
	entry.state = cs
	c.q.MoveToFront(elem)
	c.m[sessionKey] = elem
}

// Get returns the ClientSessionState value associated with a given key. It
// returns (nil, false) if no value is found.
func (c *lruSessionCache) Get(sessionKey string) (*ClientSessionState, bool) {
	c.Lock()
	defer c.Unlock()

	if elem, ok := c.m[sessionKey]; ok {
		c.q.MoveToFront(elem)
		return elem.Value.(*lruSessionCacheEntry).state, true
	}
	return nil, false
}

// TODO(jsing): Make these available to both crypto/x509 and crypto/tls.
type dsaSignature struct {
	R, S *big.Int
}

type ecdsaSignature dsaSignature

var emptyConfig Config

func defaultConfig() *Config {
	return &emptyConfig
}

var (
	once                        sync.Once
	varDefaultCipherSuites      []uint16
	varDefaultTLS13CipherSuites []uint16
)

func defaultCipherSuites() []uint16 {
	once.Do(initDefaultCipherSuites)
	return varDefaultCipherSuites
}

func defaultTLS13CipherSuites() []uint16 {
	once.Do(initDefaultCipherSuites)
	return varDefaultTLS13CipherSuites
}

func initDefaultCipherSuites() {
	var topCipherSuites, topTLS13CipherSuites []uint16
	// TODO: check for hardware support
	// This used to be: if cipherhw.AESGCMSupport() {
	// However, cipherhw is an internal package
	if true {
		// If AES-GCM hardware is provided then prioritise AES-GCM
		// cipher suites.
		topTLS13CipherSuites = []uint16{
			TLS_AES_128_GCM_SHA256,
			TLS_AES_256_GCM_SHA384,
			TLS_CHACHA20_POLY1305_SHA256,
		}
		topCipherSuites = []uint16{
			TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		}
	} else {
		// Without AES-GCM hardware, we put the ChaCha20-Poly1305
		// cipher suites first.
		topTLS13CipherSuites = []uint16{
			TLS_CHACHA20_POLY1305_SHA256,
			TLS_AES_128_GCM_SHA256,
			TLS_AES_256_GCM_SHA384,
		}
		topCipherSuites = []uint16{
			TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		}
	}

	varDefaultTLS13CipherSuites = make([]uint16, 0, len(cipherSuites))
	varDefaultTLS13CipherSuites = append(varDefaultTLS13CipherSuites, topTLS13CipherSuites...)
	varDefaultCipherSuites = make([]uint16, 0, len(cipherSuites))
	varDefaultCipherSuites = append(varDefaultCipherSuites, topCipherSuites...)

NextCipherSuite:
	for _, suite := range cipherSuites {
		if suite.flags&suiteDefaultOff != 0 {
			continue
		}
		if suite.flags&suiteTLS13 != 0 {
			for _, existing := range varDefaultTLS13CipherSuites {
				if existing == suite.id {
					continue NextCipherSuite
				}
			}
			varDefaultTLS13CipherSuites = append(varDefaultTLS13CipherSuites, suite.id)
		} else {
			for _, existing := range varDefaultCipherSuites {
				if existing == suite.id {
					continue NextCipherSuite
				}
			}
			varDefaultCipherSuites = append(varDefaultCipherSuites, suite.id)
		}
	}
	varDefaultCipherSuites = append(varDefaultTLS13CipherSuites, varDefaultCipherSuites...)
}

func unexpectedMessageError(wanted, got interface{}) error {
	return fmt.Errorf("tls: received unexpected handshake message of type %T when waiting for %T", got, wanted)
}

func isSupportedSignatureAlgorithm(sigAlg SignatureScheme, supportedSignatureAlgorithms []SignatureScheme) bool {
	for _, s := range supportedSignatureAlgorithms {
		if s == sigAlg {
			return true
		}
	}
	return false
}

// signatureFromSignatureScheme maps a signature algorithm to the underlying
// signature method (without hash function).
func signatureFromSignatureScheme(signatureAlgorithm SignatureScheme) uint8 {
	switch signatureAlgorithm {
	case PKCS1WithSHA1, PKCS1WithSHA256, PKCS1WithSHA384, PKCS1WithSHA512:
		return signaturePKCS1v15
	case PSSWithSHA256, PSSWithSHA384, PSSWithSHA512:
		return signatureRSAPSS
	case ECDSAWithSHA1, ECDSAWithP256AndSHA256, ECDSAWithP384AndSHA384, ECDSAWithP521AndSHA512:
		return signatureECDSA
	default:
		return 0
	}
}

// TODO(kk): Use variable length encoding?
func getUint24(b []byte) int {
	n := int(b[2])
	n += int(b[1] << 8)
	n += int(b[0] << 16)
	return n
}

func putUint24(b []byte, n int) {
	b[0] = byte(n >> 16)
	b[1] = byte(n >> 8)
	b[2] = byte(n & 0xff)
}
