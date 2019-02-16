// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tls

import (
	"crypto"
	"crypto/x509"
	"hash"
)

// ClientHandshakeState includes both TLS 1.3-only and TLS 1.2-only states,
// only one of them will be used, depending on negotiated version.
//
// ClientHandshakeState will be converted into and from either
//   - clientHandshakeState      (TLS 1.2)
//   - clientHandshakeStateTLS13 (TLS 1.3)
// uTLS will call .handshake() on one of these private internal states,
// to perform TLS handshake using standard crypto/tls implementation.
type ClientHandshakeState struct {
	C            *Conn
	ServerHello  *ServerHelloMsg
	Hello        *ClientHelloMsg
	MasterSecret []byte
	Session      *ClientSessionState

	State12 TLS12OnlyState
	State13 TLS13OnlyState
}

// TLS 1.3 only
type TLS13OnlyState struct {
	Suite         *CipherSuiteTLS13
	EcdheParams   EcdheParameters
	EarlySecret   []byte
	BinderKey     []byte
	CertReq       *CertificateRequestMsgTLS13
	UsingPSK      bool
	SentDummyCCS  bool
	Transcript    hash.Hash
	TrafficSecret []byte // client_application_traffic_secret_0
}

// TLS 1.2 and before only
type TLS12OnlyState struct {
	FinishedHash FinishedHash
	Suite        CipherSuite
}

func (chs *ClientHandshakeState) toPrivate13() *clientHandshakeStateTLS13 {
	if chs == nil {
		return nil
	} else {
		return &clientHandshakeStateTLS13{
			c:           chs.C,
			serverHello: chs.ServerHello.getPrivatePtr(),
			hello:       chs.Hello.getPrivatePtr(),
			ecdheParams: chs.State13.EcdheParams,

			session:     chs.Session,
			earlySecret: chs.State13.EarlySecret,
			binderKey:   chs.State13.BinderKey,

			certReq:       chs.State13.CertReq.toPrivate(),
			usingPSK:      chs.State13.UsingPSK,
			sentDummyCCS:  chs.State13.SentDummyCCS,
			suite:         chs.State13.Suite.toPrivate(),
			transcript:    chs.State13.Transcript,
			masterSecret:  chs.MasterSecret,
			trafficSecret: chs.State13.TrafficSecret,
		}
	}
}

func (chs13 *clientHandshakeStateTLS13) toPublic13() *ClientHandshakeState {
	if chs13 == nil {
		return nil
	} else {
		tls13State := TLS13OnlyState{
			EcdheParams:   chs13.ecdheParams,
			EarlySecret:   chs13.earlySecret,
			BinderKey:     chs13.binderKey,
			CertReq:       chs13.certReq.toPublic(),
			UsingPSK:      chs13.usingPSK,
			SentDummyCCS:  chs13.sentDummyCCS,
			Suite:         chs13.suite.toPublic(),
			TrafficSecret: chs13.trafficSecret,
			Transcript:    chs13.transcript,
		}
		return &ClientHandshakeState{
			C:           chs13.c,
			ServerHello: chs13.serverHello.getPublicPtr(),
			Hello:       chs13.hello.getPublicPtr(),

			Session: chs13.session,

			MasterSecret: chs13.masterSecret,

			State13: tls13State,
		}
	}
}

func (chs *ClientHandshakeState) toPrivate12() *clientHandshakeState {
	if chs == nil {
		return nil
	} else {
		return &clientHandshakeState{
			c:           chs.C,
			serverHello: chs.ServerHello.getPrivatePtr(),
			hello:       chs.Hello.getPrivatePtr(),
			suite:       chs.State12.Suite.getPrivatePtr(),
			session:     chs.Session,

			masterSecret: chs.MasterSecret,

			finishedHash: *chs.State12.FinishedHash.getPrivatePtr(),
		}
	}
}

func (chs12 *clientHandshakeState) toPublic13() *ClientHandshakeState {
	if chs12 == nil {
		return nil
	} else {
		tls12State := TLS12OnlyState{
			Suite:        *chs12.suite.getPublicPtr(),
			FinishedHash: *chs12.finishedHash.getPublicPtr(),
		}
		return &ClientHandshakeState{
			C:           chs12.c,
			ServerHello: chs12.serverHello.getPublicPtr(),
			Hello:       chs12.hello.getPublicPtr(),

			Session: chs12.session,

			MasterSecret: chs12.masterSecret,

			State12: tls12State,
		}
	}
}

type EcdheParameters interface {
	ecdheParameters
}

type CertificateRequestMsgTLS13 struct {
	Raw                              []byte
	OcspStapling                     bool
	Scts                             bool
	SupportedSignatureAlgorithms     []SignatureScheme
	SupportedSignatureAlgorithmsCert []SignatureScheme
	CertificateAuthorities           [][]byte
}

func (crm *certificateRequestMsgTLS13) toPublic() *CertificateRequestMsgTLS13 {
	if crm == nil {
		return nil
	} else {
		return &CertificateRequestMsgTLS13{
			Raw:                              crm.raw,
			OcspStapling:                     crm.ocspStapling,
			Scts:                             crm.scts,
			SupportedSignatureAlgorithms:     crm.supportedSignatureAlgorithms,
			SupportedSignatureAlgorithmsCert: crm.supportedSignatureAlgorithmsCert,
			CertificateAuthorities:           crm.certificateAuthorities,
		}
	}
}

func (crm *CertificateRequestMsgTLS13) toPrivate() *certificateRequestMsgTLS13 {
	if crm == nil {
		return nil
	} else {
		return &certificateRequestMsgTLS13{
			raw:                              crm.Raw,
			ocspStapling:                     crm.OcspStapling,
			scts:                             crm.Scts,
			supportedSignatureAlgorithms:     crm.SupportedSignatureAlgorithms,
			supportedSignatureAlgorithmsCert: crm.SupportedSignatureAlgorithmsCert,
			certificateAuthorities:           crm.CertificateAuthorities,
		}
	}
}

type CipherSuiteTLS13 struct {
	Id     uint16
	KeyLen int
	Aead   func(key, fixedNonce []byte) aead
	Hash   crypto.Hash
}

func (c *cipherSuiteTLS13) toPublic() *CipherSuiteTLS13 {
	if c == nil {
		return nil
	} else {
		return &CipherSuiteTLS13{
			Id:     c.id,
			KeyLen: c.keyLen,
			Aead:   c.aead,
			Hash:   c.hash,
		}
	}
}

func (c *CipherSuiteTLS13) toPrivate() *cipherSuiteTLS13 {
	if c == nil {
		return nil
	} else {
		return &cipherSuiteTLS13{
			id:     c.Id,
			keyLen: c.KeyLen,
			aead:   c.Aead,
			hash:   c.Hash,
		}
	}
}

type ServerHelloMsg struct {
	Raw                          []byte
	Vers                         uint16
	Random                       []byte
	SessionId                    []byte
	CipherSuite                  uint16
	CompressionMethod            uint8
	NextProtoNeg                 bool
	NextProtos                   []string
	OcspStapling                 bool
	Scts                         [][]byte
	Ems                          bool
	TicketSupported              bool
	SecureRenegotiation          []byte
	SecureRenegotiationSupported bool
	AlpnProtocol                 string

	// 1.3
	SupportedVersion        uint16
	ServerShare             keyShare
	SelectedIdentityPresent bool
	SelectedIdentity        uint16
	Cookie                  []byte  // HelloRetryRequest extension
	SelectedGroup           CurveID // HelloRetryRequest extension

}

func (shm *ServerHelloMsg) getPrivatePtr() *serverHelloMsg {
	if shm == nil {
		return nil
	} else {
		return &serverHelloMsg{
			raw:                          shm.Raw,
			vers:                         shm.Vers,
			random:                       shm.Random,
			sessionId:                    shm.SessionId,
			cipherSuite:                  shm.CipherSuite,
			compressionMethod:            shm.CompressionMethod,
			nextProtoNeg:                 shm.NextProtoNeg,
			nextProtos:                   shm.NextProtos,
			ocspStapling:                 shm.OcspStapling,
			scts:                         shm.Scts,
			ems:                          shm.Ems,
			ticketSupported:              shm.TicketSupported,
			secureRenegotiation:          shm.SecureRenegotiation,
			secureRenegotiationSupported: shm.SecureRenegotiationSupported,
			alpnProtocol:                 shm.AlpnProtocol,
			supportedVersion:             shm.SupportedVersion,
			serverShare:                  shm.ServerShare,
			selectedIdentityPresent:      shm.SelectedIdentityPresent,
			selectedIdentity:             shm.SelectedIdentity,
			cookie:                       shm.Cookie,
			selectedGroup:                shm.SelectedGroup,
		}
	}
}

func (shm *serverHelloMsg) getPublicPtr() *ServerHelloMsg {
	if shm == nil {
		return nil
	} else {
		return &ServerHelloMsg{
			Raw:                          shm.raw,
			Vers:                         shm.vers,
			Random:                       shm.random,
			SessionId:                    shm.sessionId,
			CipherSuite:                  shm.cipherSuite,
			CompressionMethod:            shm.compressionMethod,
			NextProtoNeg:                 shm.nextProtoNeg,
			NextProtos:                   shm.nextProtos,
			OcspStapling:                 shm.ocspStapling,
			Scts:                         shm.scts,
			Ems:                          shm.ems,
			TicketSupported:              shm.ticketSupported,
			SecureRenegotiation:          shm.secureRenegotiation,
			SecureRenegotiationSupported: shm.secureRenegotiationSupported,
			AlpnProtocol:                 shm.alpnProtocol,
			SupportedVersion:             shm.supportedVersion,
			ServerShare:                  shm.serverShare,
			SelectedIdentityPresent:      shm.selectedIdentityPresent,
			SelectedIdentity:             shm.selectedIdentity,
			Cookie:                       shm.cookie,
			SelectedGroup:                shm.selectedGroup,
		}
	}
}

type ClientHelloMsg struct {
	Raw                          []byte
	Vers                         uint16
	Random                       []byte
	SessionId                    []byte
	CipherSuites                 []uint16
	CompressionMethods           []uint8
	NextProtoNeg                 bool
	ServerName                   string
	OcspStapling                 bool
	Scts                         bool
	Ems                          bool // [UTLS] actually implemented due to its prevalence
	SupportedCurves              []CurveID
	SupportedPoints              []uint8
	TicketSupported              bool
	SessionTicket                []uint8
	SupportedSignatureAlgorithms []SignatureScheme
	SecureRenegotiation          []byte
	SecureRenegotiationSupported bool
	AlpnProtocols                []string

	// 1.3
	SupportedSignatureAlgorithmsCert []SignatureScheme
	SupportedVersions                []uint16
	Cookie                           []byte
	KeyShares                        []KeyShare
	EarlyData                        bool
	PskModes                         []uint8
	PskIdentities                    []pskIdentity
	PskBinders                       [][]byte
}

func (chm *ClientHelloMsg) getPrivatePtr() *clientHelloMsg {
	if chm == nil {
		return nil
	} else {
		return &clientHelloMsg{
			raw:                          chm.Raw,
			vers:                         chm.Vers,
			random:                       chm.Random,
			sessionId:                    chm.SessionId,
			cipherSuites:                 chm.CipherSuites,
			compressionMethods:           chm.CompressionMethods,
			nextProtoNeg:                 chm.NextProtoNeg,
			serverName:                   chm.ServerName,
			ocspStapling:                 chm.OcspStapling,
			scts:                         chm.Scts,
			ems:                          chm.Ems,
			supportedCurves:              chm.SupportedCurves,
			supportedPoints:              chm.SupportedPoints,
			ticketSupported:              chm.TicketSupported,
			sessionTicket:                chm.SessionTicket,
			supportedSignatureAlgorithms: chm.SupportedSignatureAlgorithms,
			secureRenegotiation:          chm.SecureRenegotiation,
			secureRenegotiationSupported: chm.SecureRenegotiationSupported,
			alpnProtocols:                chm.AlpnProtocols,

			supportedSignatureAlgorithmsCert: chm.SupportedSignatureAlgorithmsCert,
			supportedVersions:                chm.SupportedVersions,
			cookie:                           chm.Cookie,
			keyShares:                        KeyShares(chm.KeyShares).ToPrivate(),
			earlyData:                        chm.EarlyData,
			pskModes:                         chm.PskModes,
			pskIdentities:                    chm.PskIdentities,
			pskBinders:                       chm.PskBinders,
		}
	}
}

func (chm *clientHelloMsg) getPublicPtr() *ClientHelloMsg {
	if chm == nil {
		return nil
	} else {
		return &ClientHelloMsg{
			Raw:                          chm.raw,
			Vers:                         chm.vers,
			Random:                       chm.random,
			SessionId:                    chm.sessionId,
			CipherSuites:                 chm.cipherSuites,
			CompressionMethods:           chm.compressionMethods,
			NextProtoNeg:                 chm.nextProtoNeg,
			ServerName:                   chm.serverName,
			OcspStapling:                 chm.ocspStapling,
			Scts:                         chm.scts,
			Ems:                          chm.ems,
			SupportedCurves:              chm.supportedCurves,
			SupportedPoints:              chm.supportedPoints,
			TicketSupported:              chm.ticketSupported,
			SessionTicket:                chm.sessionTicket,
			SupportedSignatureAlgorithms: chm.supportedSignatureAlgorithms,
			SecureRenegotiation:          chm.secureRenegotiation,
			SecureRenegotiationSupported: chm.secureRenegotiationSupported,
			AlpnProtocols:                chm.alpnProtocols,

			SupportedSignatureAlgorithmsCert: chm.supportedSignatureAlgorithmsCert,
			SupportedVersions:                chm.supportedVersions,
			Cookie:                           chm.cookie,
			KeyShares:                        keyShares(chm.keyShares).ToPublic(),
			EarlyData:                        chm.earlyData,
			PskModes:                         chm.pskModes,
			PskIdentities:                    chm.pskIdentities,
			PskBinders:                       chm.pskBinders,
		}
	}
}

// A CipherSuite is a specific combination of key agreement, cipher and MAC
// function. All cipher suites currently assume RSA key agreement.
type CipherSuite struct {
	Id uint16
	// the lengths, in bytes, of the key material needed for each component.
	KeyLen int
	MacLen int
	IvLen  int
	Ka     func(version uint16) keyAgreement
	// flags is a bitmask of the suite* values, above.
	Flags  int
	Cipher func(key, iv []byte, isRead bool) interface{}
	Mac    func(version uint16, macKey []byte) macFunction
	Aead   func(key, fixedNonce []byte) aead
}

func (cs *CipherSuite) getPrivatePtr() *cipherSuite {
	if cs == nil {
		return nil
	} else {
		return &cipherSuite{
			id:     cs.Id,
			keyLen: cs.KeyLen,
			macLen: cs.MacLen,
			ivLen:  cs.IvLen,
			ka:     cs.Ka,
			flags:  cs.Flags,
			cipher: cs.Cipher,
			mac:    cs.Mac,
			aead:   cs.Aead,
		}
	}
}

func (cs *cipherSuite) getPublicPtr() *CipherSuite {
	if cs == nil {
		return nil
	} else {
		return &CipherSuite{
			Id:     cs.id,
			KeyLen: cs.keyLen,
			MacLen: cs.macLen,
			IvLen:  cs.ivLen,
			Ka:     cs.ka,
			Flags:  cs.flags,
			Cipher: cs.cipher,
			Mac:    cs.mac,
			Aead:   cs.aead,
		}
	}
}

// A FinishedHash calculates the hash of a set of handshake messages suitable
// for including in a Finished message.
type FinishedHash struct {
	Client hash.Hash
	Server hash.Hash

	// Prior to TLS 1.2, an additional MD5 hash is required.
	ClientMD5 hash.Hash
	ServerMD5 hash.Hash

	// In TLS 1.2, a full buffer is sadly required.
	Buffer []byte

	Version uint16
	Prf     func(result, secret, label, seed []byte)
}

func (fh *FinishedHash) getPrivatePtr() *finishedHash {
	if fh == nil {
		return nil
	} else {
		return &finishedHash{
			client:    fh.Client,
			server:    fh.Server,
			clientMD5: fh.ClientMD5,
			serverMD5: fh.ServerMD5,
			buffer:    fh.Buffer,
			version:   fh.Version,
			prf:       fh.Prf,
		}
	}
}

func (fh *finishedHash) getPublicPtr() *FinishedHash {
	if fh == nil {
		return nil
	} else {
		return &FinishedHash{
			Client:    fh.client,
			Server:    fh.server,
			ClientMD5: fh.clientMD5,
			ServerMD5: fh.serverMD5,
			Buffer:    fh.buffer,
			Version:   fh.version,
			Prf:       fh.prf}
	}
}

// TLS 1.3 Key Share. See RFC 8446, Section 4.2.8.
type KeyShare struct {
	Group CurveID
	Data  []byte
}

type KeyShares []KeyShare
type keyShares []keyShare

func (kss keyShares) ToPublic() []KeyShare {
	var KSS []KeyShare
	for _, ks := range kss {
		KSS = append(KSS, KeyShare{Data: ks.data, Group: ks.group})
	}
	return KSS
}
func (KSS KeyShares) ToPrivate() []keyShare {
	var kss []keyShare
	for _, KS := range KSS {
		kss = append(kss, keyShare{data: KS.Data, group: KS.Group})
	}
	return kss
}

// ClientSessionState is public, but all its fields are private. Let's add setters, getters and constructor

// ClientSessionState contains the state needed by clients to resume TLS sessions.
func MakeClientSessionState(
	SessionTicket []uint8,
	Vers uint16,
	CipherSuite uint16,
	MasterSecret []byte,
	ServerCertificates []*x509.Certificate,
	VerifiedChains [][]*x509.Certificate) *ClientSessionState {
	css := ClientSessionState{sessionTicket: SessionTicket,
		vers:               Vers,
		cipherSuite:        CipherSuite,
		masterSecret:       MasterSecret,
		serverCertificates: ServerCertificates,
		verifiedChains:     VerifiedChains}
	return &css
}

// Encrypted ticket used for session resumption with server
func (css *ClientSessionState) SessionTicket() []uint8 {
	return css.sessionTicket
}

// SSL/TLS version negotiated for the session
func (css *ClientSessionState) Vers() uint16 {
	return css.vers
}

// Ciphersuite negotiated for the session
func (css *ClientSessionState) CipherSuite() uint16 {
	return css.cipherSuite
}

// MasterSecret generated by client on a full handshake
func (css *ClientSessionState) MasterSecret() []byte {
	return css.masterSecret
}

// Certificate chain presented by the server
func (css *ClientSessionState) ServerCertificates() []*x509.Certificate {
	return css.serverCertificates
}

// Certificate chains we built for verification
func (css *ClientSessionState) VerifiedChains() [][]*x509.Certificate {
	return css.verifiedChains
}

func (css *ClientSessionState) SetSessionTicket(SessionTicket []uint8) {
	css.sessionTicket = SessionTicket
}
func (css *ClientSessionState) SetVers(Vers uint16) {
	css.vers = Vers
}
func (css *ClientSessionState) SetCipherSuite(CipherSuite uint16) {
	css.cipherSuite = CipherSuite
}
func (css *ClientSessionState) SetMasterSecret(MasterSecret []byte) {
	css.masterSecret = MasterSecret
}
func (css *ClientSessionState) SetServerCertificates(ServerCertificates []*x509.Certificate) {
	css.serverCertificates = ServerCertificates
}
func (css *ClientSessionState) SetVerifiedChains(VerifiedChains [][]*x509.Certificate) {
	css.verifiedChains = VerifiedChains
}
