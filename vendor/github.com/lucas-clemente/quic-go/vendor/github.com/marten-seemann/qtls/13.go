package qtls

import (
	"bytes"
	"crypto"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cloudflare/sidh/sidh"
	"golang.org/x/crypto/curve25519"
)

// numSessionTickets is the number of different session tickets the
// server sends to a TLS 1.3 client, who will use each only once.
const numSessionTickets = 2

type secretLabel int

const (
	x25519SharedSecretSz = 32

	P503PubKeySz                  = 378
	P503PrvKeySz                  = 32
	P503SharedSecretSz            = 126
	SIDHp503Curve25519PubKeySz    = x25519SharedSecretSz + P503PubKeySz
	SIDHp503Curve25519PrvKeySz    = x25519SharedSecretSz + P503PrvKeySz
	SIDHp503Curve25519SharedKeySz = x25519SharedSecretSz + P503SharedSecretSz
)

const (
	secretResumptionPskBinder secretLabel = iota
	secretEarlyClient
	secretHandshakeClient
	secretHandshakeServer
	secretApplicationClient
	secretApplicationServer
	secretResumption
)

type keySchedule13 struct {
	suite          *cipherSuite
	transcriptHash hash.Hash // uses the cipher suite hash algo
	secret         []byte    // Current secret as used for Derive-Secret
	handshakeCtx   []byte    // cached handshake context, invalidated on updates.
	clientRandom   []byte    // Used for keylogging, nil if keylogging is disabled.
	config         *Config   // Used for KeyLogWriter callback, nil if keylogging is disabled.
}

// Interface implemented by DH key exchange strategies
type dhKex interface {
	// c - context of current TLS handshake, groupId - ID of an algorithm
	// (curve/field) being chosen for key agreement. Methods implmenting an
	// interface always assume that provided groupId is correct.
	//
	// In case of success, function returns secret key and ephemeral key. Otherwise
	// error is set.
	generate(c *Conn, groupId CurveID) ([]byte, keyShare, error)
	// c - context of current TLS handshake, ks - public key received
	// from the other side of the connection, secretKey - is a private key
	// used for DH key agreement. Function returns shared secret in case
	// of success or empty slice otherwise.
	derive(c *Conn, ks keyShare, secretKey []byte) []byte
}

// Key Exchange strategies per curve type
type kexNist struct{}     // Used by NIST curves; P-256, P-384, P-512
type kexX25519 struct{}   // Used by X25519
type kexSIDHp503 struct{} // Used by SIDH/P503
type kexHybridSIDHp503X25519 struct {
	classicKEX kexX25519
	pqKEX      kexSIDHp503
} // Used by SIDH-ECDH hybrid scheme

// Routing map for key exchange strategies
var dhKexStrat = map[CurveID]dhKex{
	CurveP256:                &kexNist{},
	CurveP384:                &kexNist{},
	CurveP521:                &kexNist{},
	X25519:                   &kexX25519{},
	HybridSIDHp503Curve25519: &kexHybridSIDHp503X25519{},
}

func newKeySchedule13(suite *cipherSuite, config *Config, clientRandom []byte) *keySchedule13 {
	if config.KeyLogWriter == nil {
		clientRandom = nil
		config = nil
	}
	return &keySchedule13{
		suite:          suite,
		transcriptHash: hashForSuite(suite).New(),
		clientRandom:   clientRandom,
		config:         config,
	}
}

// setSecret sets the early/handshake/master secret based on the given secret
// (IKM). The salt is based on previous secrets (nil for the early secret).
func (ks *keySchedule13) setSecret(secret []byte) {
	hash := hashForSuite(ks.suite)
	salt := ks.secret
	if salt != nil {
		h0 := hash.New().Sum(nil)
		salt = hkdfExpandLabel(hash, salt, h0, "derived", hash.Size())
	}
	ks.secret = hkdfExtract(hash, secret, salt)
}

// Depending on role returns pair of key variant to be used by
// local and remote process.
func getSidhKeyVariant(isClient bool) (sidh.KeyVariant, sidh.KeyVariant) {
	if isClient {
		return sidh.KeyVariant_SIDH_A, sidh.KeyVariant_SIDH_B
	}
	return sidh.KeyVariant_SIDH_B, sidh.KeyVariant_SIDH_A
}

// write appends the data to the transcript hash context.
func (ks *keySchedule13) write(data []byte) {
	ks.handshakeCtx = nil
	ks.transcriptHash.Write(data)
}

func (ks *keySchedule13) getLabel(secretLabel secretLabel) (label, keylogType string) {
	switch secretLabel {
	case secretResumptionPskBinder:
		label = "res binder"
	case secretEarlyClient:
		label = "c e traffic"
		keylogType = "CLIENT_EARLY_TRAFFIC_SECRET"
	case secretHandshakeClient:
		label = "c hs traffic"
		keylogType = "CLIENT_HANDSHAKE_TRAFFIC_SECRET"
	case secretHandshakeServer:
		label = "s hs traffic"
		keylogType = "SERVER_HANDSHAKE_TRAFFIC_SECRET"
	case secretApplicationClient:
		label = "c ap traffic"
		keylogType = "CLIENT_TRAFFIC_SECRET_0"
	case secretApplicationServer:
		label = "s ap traffic"
		keylogType = "SERVER_TRAFFIC_SECRET_0"
	case secretResumption:
		label = "res master"
	}
	return
}

// deriveSecret returns the secret derived from the handshake context and label.
func (ks *keySchedule13) deriveSecret(secretLabel secretLabel) []byte {
	label, keylogType := ks.getLabel(secretLabel)
	if ks.handshakeCtx == nil {
		ks.handshakeCtx = ks.transcriptHash.Sum(nil)
	}
	hash := hashForSuite(ks.suite)
	secret := hkdfExpandLabel(hash, ks.secret, ks.handshakeCtx, label, hash.Size())
	if keylogType != "" && ks.config != nil {
		ks.config.writeKeyLog(keylogType, ks.clientRandom, secret)
	}
	return secret
}

func (ks *keySchedule13) prepareCipher(trafficSecret []byte) cipher.AEAD {
	hash := hashForSuite(ks.suite)
	key := hkdfExpandLabel(hash, trafficSecret, nil, "key", ks.suite.keyLen)
	iv := hkdfExpandLabel(hash, trafficSecret, nil, "iv", ks.suite.ivLen)
	return ks.suite.aead(key, iv)
}

func (hs *serverHandshakeState) doTLS13Handshake() error {
	config := hs.c.config
	c := hs.c

	hs.c.cipherSuite, hs.hello.cipherSuite = hs.suite.id, hs.suite.id
	hs.c.clientHello = hs.clientHello.marshal()

	// When picking the group for the handshake, priority is given to groups
	// that the client provided a keyShare for, so to avoid a round-trip.
	// After that the order of CurvePreferences is respected.
	var ks keyShare
CurvePreferenceLoop:
	for _, curveID := range config.curvePreferences() {
		for _, keyShare := range hs.clientHello.keyShares {
			if curveID == keyShare.group {
				ks = keyShare
				break CurvePreferenceLoop
			}
		}
	}
	if ks.group == 0 {
		c.sendAlert(alertInternalError)
		return errors.New("tls: HelloRetryRequest not implemented") // TODO(filippo)
	}

	privateKey, serverKS, err := c.generateKeyShare(ks.group)
	if err != nil {
		c.sendAlert(alertInternalError)
		return err
	}
	hs.hello.keyShare = serverKS

	hash := hashForSuite(hs.suite)
	hashSize := hash.Size()
	hs.keySchedule = newKeySchedule13(hs.suite, config, hs.clientHello.random)

	// Check for PSK and update key schedule with new early secret key
	isResumed, pskAlert := hs.checkPSK()
	switch {
	case pskAlert != alertSuccess:
		c.sendAlert(pskAlert)
		return errors.New("tls: invalid client PSK")
	case !isResumed:
		// apply an empty PSK if not resumed.
		hs.keySchedule.setSecret(nil)
	case isResumed:
		c.didResume = true
	}

	hs.keySchedule.write(hs.clientHello.marshal())

	earlyClientTrafficSecret := hs.keySchedule.deriveSecret(secretEarlyClient)

	ecdheSecret := c.deriveDHESecret(ks, privateKey)
	if ecdheSecret == nil {
		c.sendAlert(alertIllegalParameter)
		return errors.New("tls: bad ECDHE client share")
	}

	hs.keySchedule.write(hs.hello.marshal())
	if _, err := c.writeRecord(recordTypeHandshake, hs.hello.marshal()); err != nil {
		return err
	}

	// middlebox compatibility mode: send CCS after first handshake message
	if _, err := c.writeRecord(recordTypeChangeCipherSpec, []byte{1}); err != nil {
		return err
	}

	hs.keySchedule.setSecret(ecdheSecret)
	hs.hsClientTrafficSecret = hs.keySchedule.deriveSecret(secretHandshakeClient)
	hsServerTrafficSecret := hs.keySchedule.deriveSecret(secretHandshakeServer)
	c.out.exportKey(hs.keySchedule.suite, hsServerTrafficSecret)
	c.out.setKey(c.vers, hs.keySchedule.suite, hsServerTrafficSecret)

	serverFinishedKey := hkdfExpandLabel(hash, hsServerTrafficSecret, nil, "finished", hashSize)
	hs.clientFinishedKey = hkdfExpandLabel(hash, hs.hsClientTrafficSecret, nil, "finished", hashSize)

	// EncryptedExtensions
	hs.keySchedule.write(hs.hello13Enc.marshal())
	if _, err := c.writeRecord(recordTypeHandshake, hs.hello13Enc.marshal()); err != nil {
		return err
	}

	// TODO: we should have 2 separated methods - one for full-handshake and the other for PSK-handshake
	if !c.didResume {
		// Server MUST NOT send CertificateRequest if authenticating with PSK
		if c.config.ClientAuth >= RequestClientCert {

			certReq := new(certificateRequestMsg13)
			// extension 'signature_algorithms' MUST be specified
			certReq.supportedSignatureAlgorithms = supportedSignatureAlgorithms13
			certReq.supportedSignatureAlgorithmsCert = supportedSigAlgorithmsCert(supportedSignatureAlgorithms13)
			hs.keySchedule.write(certReq.marshal())
			if _, err := hs.c.writeRecord(recordTypeHandshake, certReq.marshal()); err != nil {
				return err
			}
		}

		if err := hs.sendCertificate13(); err != nil {
			return err
		}
	}

	verifyData := hmacOfSum(hash, hs.keySchedule.transcriptHash, serverFinishedKey)
	serverFinished := &finishedMsg{
		verifyData: verifyData,
	}
	hs.keySchedule.write(serverFinished.marshal())
	if _, err := c.writeRecord(recordTypeHandshake, serverFinished.marshal()); err != nil {
		return err
	}

	hs.keySchedule.setSecret(nil) // derive master secret
	serverAppTrafficSecret := hs.keySchedule.deriveSecret(secretApplicationServer)
	c.out.exportKey(hs.keySchedule.suite, serverAppTrafficSecret)
	c.out.setKey(c.vers, hs.keySchedule.suite, serverAppTrafficSecret)

	if c.hand.Len() > 0 {
		return c.sendAlert(alertUnexpectedMessage)
	}
	hs.appClientTrafficSecret = hs.keySchedule.deriveSecret(secretApplicationClient)
	if hs.hello13Enc.earlyData {
		c.in.exportKey(hs.keySchedule.suite, earlyClientTrafficSecret)
		c.in.setKey(c.vers, hs.keySchedule.suite, earlyClientTrafficSecret)
		c.phase = readingEarlyData
	} else {
		c.in.exportKey(hs.keySchedule.suite, hs.hsClientTrafficSecret)
		c.in.setKey(c.vers, hs.keySchedule.suite, hs.hsClientTrafficSecret)
		if hs.clientHello.earlyData {
			c.phase = discardingEarlyData
		} else {
			c.phase = waitingClientFinished
		}
	}

	return nil
}

// readClientFinished13 is called during the server handshake (when no early
// data it available) or after reading all early data. It discards early data if
// the server did not accept it and then verifies the Finished message. Once
// done it sends the session tickets. Under c.in lock.
func (hs *serverHandshakeState) readClientFinished13(hasConfirmLock bool) error {
	c := hs.c

	// If the client advertised and sends early data while the server does
	// not accept it, it must be fully skipped until the Finished message.
	for c.phase == discardingEarlyData {
		if err := c.readRecord(recordTypeApplicationData); err != nil {
			return err
		}
		// Assume receipt of Finished message (will be checked below).
		if c.hand.Len() > 0 {
			c.phase = waitingClientFinished
			break
		}
	}

	// If the client sends early data followed by a Finished message (but
	// no end_of_early_data), the server MUST terminate the connection.
	if c.phase != waitingClientFinished {
		c.sendAlert(alertUnexpectedMessage)
		return errors.New("tls: did not expect Client Finished yet")
	}

	c.phase = readingClientFinished
	msg, err := c.readHandshake()
	if err != nil {
		return err
	}

	// client authentication
	// (4.4.2) Client MUST send certificate msg if requested by server
	if c.config.ClientAuth >= RequestClientCert && !c.didResume {
		certMsg, ok := msg.(*certificateMsg13)
		if !ok {
			c.sendAlert(alertCertificateRequired)
			return unexpectedMessageError(certMsg, msg)
		}

		hs.keySchedule.write(certMsg.marshal())
		certs := getCertsFromEntries(certMsg.certificates)
		pubKey, err := hs.processCertsFromClient(certs)
		if err != nil {
			return err
		}

		if len(certs) > 0 {
			// 4.4.3: CertificateVerify MUST appear immediately after Certificate msg
			msg, err = c.readHandshake()
			if err != nil {
				return err
			}

			certVerify, ok := msg.(*certificateVerifyMsg)
			if !ok {
				c.sendAlert(alertUnexpectedMessage)
				return unexpectedMessageError(certVerify, msg)
			}

			err, alertCode := verifyPeerHandshakeSignature(
				certVerify,
				pubKey,
				supportedSignatureAlgorithms13,
				hs.keySchedule.transcriptHash.Sum(nil),
				"TLS 1.3, client CertificateVerify")
			if err != nil {
				c.sendAlert(alertCode)
				return err
			}
			hs.keySchedule.write(certVerify.marshal())
		}

		// Read next chunk
		msg, err = c.readHandshake()
		if err != nil {
			return err
		}
	}

	clientFinished, ok := msg.(*finishedMsg)
	if !ok {
		c.sendAlert(alertUnexpectedMessage)
		return unexpectedMessageError(clientFinished, msg)
	}

	hash := hashForSuite(hs.suite)
	expectedVerifyData := hmacOfSum(hash, hs.keySchedule.transcriptHash, hs.clientFinishedKey)
	if len(expectedVerifyData) != len(clientFinished.verifyData) ||
		subtle.ConstantTimeCompare(expectedVerifyData, clientFinished.verifyData) != 1 {
		c.sendAlert(alertDecryptError)
		return errors.New("tls: client's Finished message is incorrect")
	}
	hs.keySchedule.write(clientFinished.marshal())

	c.hs = nil // Discard the server handshake state
	if c.hand.Len() > 0 {
		return c.sendAlert(alertUnexpectedMessage)
	}
	c.in.exportKey(hs.keySchedule.suite, hs.appClientTrafficSecret)
	c.in.setKey(c.vers, hs.keySchedule.suite, hs.appClientTrafficSecret)
	c.in.traceErr, c.out.traceErr = nil, nil
	c.phase = handshakeConfirmed
	atomic.StoreInt32(&c.handshakeConfirmed, 1)

	// Any read operation after handshakeRunning and before handshakeConfirmed
	// will be holding this lock, which we release as soon as the confirmation
	// happens, even if the Read call might do more work.
	// If a Handshake is pending, c.confirmMutex will never be locked as
	// ConfirmHandshake will wait for the handshake to complete. If a
	// handshake was complete, and this was a confirmation, unlock
	// c.confirmMutex now to allow readers to proceed.
	if hasConfirmLock {
		c.confirmMutex.Unlock()
	}

	return hs.sendSessionTicket13() // TODO: do in a goroutine
}

func (hs *serverHandshakeState) sendCertificate13() error {
	c := hs.c

	certEntries := []certificateEntry{}
	for _, cert := range hs.cert.Certificate {
		certEntries = append(certEntries, certificateEntry{data: cert})
	}
	if len(certEntries) > 0 && hs.clientHello.ocspStapling {
		certEntries[0].ocspStaple = hs.cert.OCSPStaple
	}
	if len(certEntries) > 0 && hs.clientHello.scts {
		certEntries[0].sctList = hs.cert.SignedCertificateTimestamps
	}

	// If hs.delegatedCredential is set (see hs.readClientHello()) then the
	// server is using the delegated credential extension. The DC is added as an
	// extension to the end-entity certificate, i.e., the last CertificateEntry
	// of Certificate.certficate_list. (For details, see
	// https://tools.ietf.org/html/draft-ietf-tls-subcerts-02.)
	if len(certEntries) > 0 && hs.clientHello.delegatedCredential && hs.delegatedCredential != nil {
		certEntries[0].delegatedCredential = hs.delegatedCredential
	}

	certMsg := &certificateMsg13{certificates: certEntries}

	hs.keySchedule.write(certMsg.marshal())
	if _, err := c.writeRecord(recordTypeHandshake, certMsg.marshal()); err != nil {
		return err
	}

	sigScheme, err := hs.selectTLS13SignatureScheme()
	if err != nil {
		c.sendAlert(alertInternalError)
		return err
	}

	sigHash := hashForSignatureScheme(sigScheme)
	opts := crypto.SignerOpts(sigHash)
	if signatureSchemeIsPSS(sigScheme) {
		opts = &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: sigHash}
	}

	toSign := prepareDigitallySigned(sigHash, "TLS 1.3, server CertificateVerify", hs.keySchedule.transcriptHash.Sum(nil))
	signature, err := hs.privateKey.(crypto.Signer).Sign(c.config.rand(), toSign[:], opts)
	if err != nil {
		c.sendAlert(alertInternalError)
		return err
	}

	verifyMsg := &certificateVerifyMsg{
		hasSignatureAndHash: true,
		signatureAlgorithm:  sigScheme,
		signature:           signature,
	}
	hs.keySchedule.write(verifyMsg.marshal())
	if _, err := c.writeRecord(recordTypeHandshake, verifyMsg.marshal()); err != nil {
		return err
	}

	return nil
}

func (c *Conn) handleEndOfEarlyData() error {
	if c.phase != readingEarlyData || c.vers < VersionTLS13 {
		return c.in.setErrorLocked(c.sendAlert(alertUnexpectedMessage))
	}
	msg, err := c.readHandshake()
	if err != nil {
		return err
	}
	endOfEarlyData, ok := msg.(*endOfEarlyDataMsg)
	// No handshake messages are allowed after EOD.
	if !ok || c.hand.Len() > 0 {
		return c.in.setErrorLocked(c.sendAlert(alertUnexpectedMessage))
	}
	c.hs.keySchedule.write(endOfEarlyData.marshal())
	c.phase = waitingClientFinished
	c.in.exportKey(c.hs.keySchedule.suite, c.hs.hsClientTrafficSecret)
	c.in.setKey(c.vers, c.hs.keySchedule.suite, c.hs.hsClientTrafficSecret)
	return nil
}

// selectTLS13SignatureScheme chooses the SignatureScheme for the CertificateVerify
// based on the certificate type and client supported schemes. If no overlap is found,
// a fallback is selected.
//
// See https://tools.ietf.org/html/draft-ietf-tls-tls13-18#section-4.4.1.2
func (hs *serverHandshakeState) selectTLS13SignatureScheme() (sigScheme SignatureScheme, err error) {
	var supportedSchemes []SignatureScheme
	signer, ok := hs.privateKey.(crypto.Signer)
	if !ok {
		return 0, errors.New("tls: private key does not implement crypto.Signer")
	}
	pk := signer.Public()
	if _, ok := pk.(*rsa.PublicKey); ok {
		sigScheme = PSSWithSHA256
		supportedSchemes = []SignatureScheme{PSSWithSHA256, PSSWithSHA384, PSSWithSHA512}
	} else if pk, ok := pk.(*ecdsa.PublicKey); ok {
		switch pk.Curve {
		case elliptic.P256():
			sigScheme = ECDSAWithP256AndSHA256
			supportedSchemes = []SignatureScheme{ECDSAWithP256AndSHA256}
		case elliptic.P384():
			sigScheme = ECDSAWithP384AndSHA384
			supportedSchemes = []SignatureScheme{ECDSAWithP384AndSHA384}
		case elliptic.P521():
			sigScheme = ECDSAWithP521AndSHA512
			supportedSchemes = []SignatureScheme{ECDSAWithP521AndSHA512}
		default:
			return 0, errors.New("tls: unknown ECDSA certificate curve")
		}
	} else {
		return 0, errors.New("tls: unknown certificate key type")
	}

	for _, ss := range supportedSchemes {
		for _, cs := range hs.clientHello.supportedSignatureAlgorithms {
			if ss == cs {
				return ss, nil
			}
		}
	}

	return sigScheme, nil
}

func signatureSchemeIsPSS(s SignatureScheme) bool {
	return s == PSSWithSHA256 || s == PSSWithSHA384 || s == PSSWithSHA512
}

// hashForSignatureScheme returns the Hash used by a SignatureScheme which is
// supported by selectTLS13SignatureScheme.
func hashForSignatureScheme(ss SignatureScheme) crypto.Hash {
	switch ss {
	case PSSWithSHA256, ECDSAWithP256AndSHA256:
		return crypto.SHA256
	case PSSWithSHA384, ECDSAWithP384AndSHA384:
		return crypto.SHA384
	case PSSWithSHA512, ECDSAWithP521AndSHA512:
		return crypto.SHA512
	default:
		panic("unsupported SignatureScheme passed to hashForSignatureScheme")
	}
}

func hashForSuite(suite *cipherSuite) crypto.Hash {
	if suite.flags&suiteSHA384 != 0 {
		return crypto.SHA384
	}
	return crypto.SHA256
}

func prepareDigitallySigned(hash crypto.Hash, context string, data []byte) []byte {
	message := bytes.Repeat([]byte{32}, 64)
	message = append(message, context...)
	message = append(message, 0)
	message = append(message, data...)
	h := hash.New()
	h.Write(message)
	return h.Sum(nil)
}

// generateKeyShare generates keypair. Private key is returned as first argument, public key
// is returned in keyShare.data. keyshare.curveID stores ID of the scheme used.
func (c *Conn) generateKeyShare(curveID CurveID) ([]byte, keyShare, error) {
	if val, ok := dhKexStrat[curveID]; ok {
		return val.generate(c, curveID)
	}
	return nil, keyShare{}, errors.New("tls: preferredCurves includes unsupported curve")
}

// DH key agreement. ks stores public key, secretKey stores private key used for ephemeral
// key agreement. Function returns shared secret in case of success or empty slice otherwise.
func (c *Conn) deriveDHESecret(ks keyShare, secretKey []byte) []byte {
	if val, ok := dhKexStrat[ks.group]; ok {
		return val.derive(c, ks, secretKey)
	}
	return nil
}

// HkdfExpandLabel HKDF expands a label
func HkdfExpandLabel(hash crypto.Hash, secret, hashValue []byte, label string, L int) []byte {
	return hkdfExpandLabel(hash, secret, hashValue, label, L)
}

func hkdfExpandLabel(hash crypto.Hash, secret, hashValue []byte, label string, L int) []byte {
	prefix := "tls13 "
	hkdfLabel := make([]byte, 4+len(prefix)+len(label)+len(hashValue))
	hkdfLabel[0] = byte(L >> 8)
	hkdfLabel[1] = byte(L)
	hkdfLabel[2] = byte(len(prefix) + len(label))
	copy(hkdfLabel[3:], prefix)
	z := hkdfLabel[3+len(prefix):]
	copy(z, label)
	z = z[len(label):]
	z[0] = byte(len(hashValue))
	copy(z[1:], hashValue)

	return hkdfExpand(hash, secret, hkdfLabel, L)
}

func hmacOfSum(f crypto.Hash, hash hash.Hash, key []byte) []byte {
	h := hmac.New(f.New, key)
	h.Write(hash.Sum(nil))
	return h.Sum(nil)
}

// Maximum allowed mismatch between the stated age of a ticket
// and the server-observed one. See
// https://tools.ietf.org/html/draft-ietf-tls-tls13-18#section-4.2.8.2.
const ticketAgeSkewAllowance = 10 * time.Second

// checkPSK tries to resume using a PSK, returning true (and updating the
// early secret in the key schedule) if the PSK was used and false otherwise.
func (hs *serverHandshakeState) checkPSK() (isResumed bool, alert alert) {
	if hs.c.config.SessionTicketsDisabled {
		return false, alertSuccess
	}

	foundDHE := false
	for _, mode := range hs.clientHello.pskKeyExchangeModes {
		if mode == pskDHEKeyExchange {
			foundDHE = true
			break
		}
	}
	if !foundDHE {
		return false, alertSuccess
	}

	hash := hashForSuite(hs.suite)
	hashSize := hash.Size()
	for i := range hs.clientHello.psks {
		sessionTicket := append([]uint8{}, hs.clientHello.psks[i].identity...)
		if hs.c.config.SessionTicketSealer != nil {
			var ok bool
			sessionTicket, ok = hs.c.config.SessionTicketSealer.Unseal(hs.clientHelloInfo(), sessionTicket)
			if !ok {
				continue
			}
		} else {
			sessionTicket, _ = hs.c.decryptTicket(sessionTicket)
			if sessionTicket == nil {
				continue
			}
		}
		s := &sessionState13{}
		if s.unmarshal(sessionTicket) != alertSuccess {
			continue
		}
		if s.vers != hs.c.vers {
			continue
		}
		clientAge := time.Duration(hs.clientHello.psks[i].obfTicketAge-s.ageAdd) * time.Millisecond
		serverAge := time.Since(time.Unix(int64(s.createdAt), 0))
		if clientAge-serverAge > ticketAgeSkewAllowance || clientAge-serverAge < -ticketAgeSkewAllowance {
			// XXX: NSS is off spec and sends obfuscated_ticket_age as seconds
			clientAge = time.Duration(hs.clientHello.psks[i].obfTicketAge-s.ageAdd) * time.Second
			if clientAge-serverAge > ticketAgeSkewAllowance || clientAge-serverAge < -ticketAgeSkewAllowance {
				continue
			}
		}

		// This enforces the stricter 0-RTT requirements on all ticket uses.
		// The benefit of using PSK+ECDHE without 0-RTT are small enough that
		// we can give them up in the edge case of changed suite or ALPN or SNI.
		if s.suite != hs.suite.id {
			continue
		}
		if s.alpnProtocol != hs.c.clientProtocol {
			continue
		}
		if s.SNI != hs.c.serverName {
			continue
		}

		hs.keySchedule.setSecret(s.pskSecret)
		binderKey := hs.keySchedule.deriveSecret(secretResumptionPskBinder)
		binderFinishedKey := hkdfExpandLabel(hash, binderKey, nil, "finished", hashSize)
		chHash := hash.New()
		chHash.Write(hs.clientHello.rawTruncated)
		expectedBinder := hmacOfSum(hash, chHash, binderFinishedKey)

		if subtle.ConstantTimeCompare(expectedBinder, hs.clientHello.psks[i].binder) != 1 {
			return false, alertDecryptError
		}

		if i == 0 && hs.clientHello.earlyData {
			// This is a ticket intended to be used for 0-RTT
			if s.maxEarlyDataLen == 0 {
				// But we had not tagged it as such.
				return false, alertIllegalParameter
			}
			if hs.c.config.Accept0RTTData {
				hs.c.binder = expectedBinder
				hs.c.ticketMaxEarlyData = int64(s.maxEarlyDataLen)
				hs.hello13Enc.earlyData = true
			}
		}
		hs.hello.psk = true
		hs.hello.pskIdentity = uint16(i)
		return true, alertSuccess
	}

	return false, alertSuccess
}

func (hs *serverHandshakeState) sendSessionTicket13() error {
	c := hs.c
	if c.config.SessionTicketsDisabled {
		return nil
	}

	foundDHE := false
	for _, mode := range hs.clientHello.pskKeyExchangeModes {
		if mode == pskDHEKeyExchange {
			foundDHE = true
			break
		}
	}
	if !foundDHE {
		return nil
	}

	resumptionMasterSecret := hs.keySchedule.deriveSecret(secretResumption)

	ageAddBuf := make([]byte, 4)
	sessionState := &sessionState13{
		vers:            c.vers,
		suite:           hs.suite.id,
		createdAt:       uint64(time.Now().Unix()),
		alpnProtocol:    c.clientProtocol,
		SNI:             c.serverName,
		maxEarlyDataLen: c.config.Max0RTTDataSize,
	}
	hash := hashForSuite(hs.suite)

	for i := 0; i < numSessionTickets; i++ {
		if _, err := io.ReadFull(c.config.rand(), ageAddBuf); err != nil {
			c.sendAlert(alertInternalError)
			return err
		}
		sessionState.ageAdd = uint32(ageAddBuf[0])<<24 | uint32(ageAddBuf[1])<<16 |
			uint32(ageAddBuf[2])<<8 | uint32(ageAddBuf[3])
		// ticketNonce must be a unique value for this connection.
		// Assume there are no more than 255 tickets, otherwise two
		// tickets might have the same PSK which could be a problem if
		// one of them is compromised.
		ticketNonce := []byte{byte(i)}
		sessionState.pskSecret = hkdfExpandLabel(hash, resumptionMasterSecret, ticketNonce, "resumption", hash.Size())
		ticket := sessionState.marshal()
		var err error
		if c.config.SessionTicketSealer != nil {
			cs := c.ConnectionState()
			ticket, err = c.config.SessionTicketSealer.Seal(&cs, ticket)
		} else {
			ticket, err = c.encryptTicket(ticket)
		}
		if err != nil {
			c.sendAlert(alertInternalError)
			return err
		}
		if ticket == nil {
			continue
		}
		ticketMsg := &newSessionTicketMsg13{
			lifetime:           24 * 3600, // TODO(filippo)
			maxEarlyDataLength: c.config.Max0RTTDataSize,
			withEarlyDataInfo:  c.config.Max0RTTDataSize > 0,
			ageAdd:             sessionState.ageAdd,
			nonce:              ticketNonce,
			ticket:             ticket,
		}
		if _, err := c.writeRecord(recordTypeHandshake, ticketMsg.marshal()); err != nil {
			return err
		}
	}

	return nil
}

func (hs *serverHandshakeState) traceErr(err error) {
	if err == nil {
		return
	}
	if os.Getenv("TLSDEBUG") == "error" {
		if hs != nil && hs.clientHello != nil {
			os.Stderr.WriteString(hex.Dump(hs.clientHello.marshal()))
		} else if err == io.EOF {
			return // don't stack trace on EOF before CH
		}
		fmt.Fprintf(os.Stderr, "\n%s\n", debug.Stack())
	}
	if os.Getenv("TLSDEBUG") == "short" {
		var pcs [4]uintptr
		frames := runtime.CallersFrames(pcs[0:runtime.Callers(3, pcs[:])])
		for {
			frame, more := frames.Next()
			if frame.Function != "crypto/tls.(*halfConn).setErrorLocked" &&
				frame.Function != "crypto/tls.(*Conn).sendAlertLocked" &&
				frame.Function != "crypto/tls.(*Conn).sendAlert" {
				file := frame.File[strings.LastIndex(frame.File, "/")+1:]
				log.Printf("%s:%d (%s): %v", file, frame.Line, frame.Function, err)
				return
			}
			if !more {
				break
			}
		}
	}
}

func getCertsFromEntries(certEntries []certificateEntry) [][]byte {
	certs := make([][]byte, len(certEntries))
	for i, cert := range certEntries {
		certs[i] = cert.data
	}
	return certs
}

func (hs *clientHandshakeState) processEncryptedExtensions(ee *encryptedExtensionsMsg) error {
	c := hs.c
	if ee.alpnProtocol != "" {
		c.clientProtocol = ee.alpnProtocol
		c.clientProtocolFallback = false
	}
	if hs.c.config.ReceivedExtensions != nil {
		return hs.c.config.ReceivedExtensions(typeEncryptedExtensions, ee.additionalExtensions)
	}
	return nil
}

func verifyPeerHandshakeSignature(
	certVerify *certificateVerifyMsg,
	pubKey crypto.PublicKey,
	signAlgosKnown []SignatureScheme,
	transHash []byte,
	contextString string) (error, alert) {

	_, sigType, hashFunc, err := pickSignatureAlgorithm(
		pubKey,
		[]SignatureScheme{certVerify.signatureAlgorithm},
		signAlgosKnown,
		VersionTLS13)
	if err != nil {
		return err, alertHandshakeFailure
	}

	digest := prepareDigitallySigned(hashFunc, contextString, transHash)
	err = verifyHandshakeSignature(sigType, pubKey, hashFunc, digest, certVerify.signature)

	if err != nil {
		return err, alertDecryptError
	}

	return nil, alertSuccess
}

func (hs *clientHandshakeState) getCertificate13(certReq *certificateRequestMsg13) (*Certificate, error) {
	certReq12 := &certificateRequestMsg{
		hasSignatureAndHash:          true,
		supportedSignatureAlgorithms: certReq.supportedSignatureAlgorithms,
		certificateAuthorities:       certReq.certificateAuthorities,
	}

	var rsaAvail, ecdsaAvail bool
	for _, sigAlg := range certReq.supportedSignatureAlgorithms {
		switch signatureFromSignatureScheme(sigAlg) {
		case signaturePKCS1v15, signatureRSAPSS:
			rsaAvail = true
		case signatureECDSA:
			ecdsaAvail = true
		}
	}
	if rsaAvail {
		certReq12.certificateTypes = append(certReq12.certificateTypes, certTypeRSASign)
	}
	if ecdsaAvail {
		certReq12.certificateTypes = append(certReq12.certificateTypes, certTypeECDSASign)
	}

	return hs.getCertificate(certReq12)
}

func (hs *clientHandshakeState) sendCertificate13(chainToSend *Certificate, certReq *certificateRequestMsg13) error {
	c := hs.c

	certEntries := []certificateEntry{}
	for _, cert := range chainToSend.Certificate {
		certEntries = append(certEntries, certificateEntry{data: cert})
	}
	certMsg := &certificateMsg13{certificates: certEntries}

	hs.keySchedule.write(certMsg.marshal())
	if _, err := c.writeRecord(recordTypeHandshake, certMsg.marshal()); err != nil {
		return err
	}

	if len(certEntries) == 0 {
		// No client cert available, nothing to sign.
		return nil
	}

	key, ok := chainToSend.PrivateKey.(crypto.Signer)
	if !ok {
		c.sendAlert(alertInternalError)
		return fmt.Errorf("tls: client certificate private key of type %T does not implement crypto.Signer", chainToSend.PrivateKey)
	}

	signatureAlgorithm, sigType, hashFunc, err := pickSignatureAlgorithm(key.Public(), certReq.supportedSignatureAlgorithms, hs.hello.supportedSignatureAlgorithms, c.vers)
	if err != nil {
		hs.c.sendAlert(alertHandshakeFailure)
		return err
	}

	digest := prepareDigitallySigned(hashFunc, "TLS 1.3, client CertificateVerify", hs.keySchedule.transcriptHash.Sum(nil))
	signOpts := crypto.SignerOpts(hashFunc)
	if sigType == signatureRSAPSS {
		signOpts = &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: hashFunc}
	}
	signature, err := key.Sign(c.config.rand(), digest, signOpts)
	if err != nil {
		c.sendAlert(alertInternalError)
		return err
	}

	verifyMsg := &certificateVerifyMsg{
		hasSignatureAndHash: true,
		signatureAlgorithm:  signatureAlgorithm,
		signature:           signature,
	}
	hs.keySchedule.write(verifyMsg.marshal())
	if _, err := c.writeRecord(recordTypeHandshake, verifyMsg.marshal()); err != nil {
		return err
	}

	return nil
}

func (hs *clientHandshakeState) doTLS13Handshake() error {
	c := hs.c
	hash := hashForSuite(hs.suite)
	hashSize := hash.Size()
	serverHello := hs.serverHello
	c.scts = serverHello.scts

	// middlebox compatibility mode, send CCS before second flight.
	if _, err := c.writeRecord(recordTypeChangeCipherSpec, []byte{1}); err != nil {
		return err
	}

	// TODO check if keyshare is unacceptable, raise HRR.

	clientKS := hs.hello.keyShares[0]
	if serverHello.keyShare.group != clientKS.group {
		c.sendAlert(alertIllegalParameter)
		return errors.New("bad or missing key share from server")
	}

	// 0-RTT is not supported yet, so use an empty PSK.
	hs.keySchedule.setSecret(nil)
	ecdheSecret := c.deriveDHESecret(serverHello.keyShare, hs.privateKey)
	if ecdheSecret == nil {
		c.sendAlert(alertIllegalParameter)
		return errors.New("tls: bad ECDHE server share")
	}

	// Calculate handshake secrets.
	hs.keySchedule.setSecret(ecdheSecret)
	clientHandshakeSecret := hs.keySchedule.deriveSecret(secretHandshakeClient)
	if c.hand.Len() > 0 {
		c.sendAlert(alertUnexpectedMessage)
		return errors.New("tls: unexpected data after Server Hello")
	}
	serverHandshakeSecret := hs.keySchedule.deriveSecret(secretHandshakeServer)
	c.in.exportKey(hs.keySchedule.suite, serverHandshakeSecret)
	// Already the sender key yet, when using an alternative record layer.
	// QUIC needs the handshake write key in order to acknowlege Handshake packets.
	c.out.exportKey(hs.keySchedule.suite, clientHandshakeSecret)
	// Do not change the sender key yet, the server must authenticate first.
	c.in.setKey(c.vers, hs.keySchedule.suite, serverHandshakeSecret)

	// Calculate MAC key for Finished messages.
	serverFinishedKey := hkdfExpandLabel(hash, serverHandshakeSecret, nil, "finished", hashSize)
	clientFinishedKey := hkdfExpandLabel(hash, clientHandshakeSecret, nil, "finished", hashSize)

	msg, err := c.readHandshake()
	if err != nil {
		return err
	}
	encryptedExtensions, ok := msg.(*encryptedExtensionsMsg)
	if !ok {
		c.sendAlert(alertUnexpectedMessage)
		return unexpectedMessageError(encryptedExtensions, msg)
	}
	if err := hs.processEncryptedExtensions(encryptedExtensions); err != nil {
		return err
	}
	hs.keySchedule.write(encryptedExtensions.marshal())

	// PSKs are not supported, so receive Certificate message.
	msg, err = c.readHandshake()
	if err != nil {
		return err
	}

	var chainToSend *Certificate
	certReq, isCertRequested := msg.(*certificateRequestMsg13)
	if isCertRequested {
		hs.keySchedule.write(certReq.marshal())

		if chainToSend, err = hs.getCertificate13(certReq); err != nil {
			c.sendAlert(alertInternalError)
			return err
		}

		msg, err = c.readHandshake()
		if err != nil {
			return err
		}
	}

	certMsg, ok := msg.(*certificateMsg13)
	if !ok {
		c.sendAlert(alertUnexpectedMessage)
		return unexpectedMessageError(certMsg, msg)
	}
	hs.keySchedule.write(certMsg.marshal())

	// Validate certificates.
	certs := getCertsFromEntries(certMsg.certificates)
	if err := hs.processCertsFromServer(certs); err != nil {
		return err
	}

	// Receive CertificateVerify message.
	msg, err = c.readHandshake()
	if err != nil {
		return err
	}
	certVerifyMsg, ok := msg.(*certificateVerifyMsg)
	if !ok {
		c.sendAlert(alertUnexpectedMessage)
		return unexpectedMessageError(certVerifyMsg, msg)
	}

	// Validate the DC if present. The DC is only processed if the extension was
	// indicated by the ClientHello; otherwise this call will result in an
	// "illegal_parameter" alert.
	if len(certMsg.certificates) > 0 {
		if err := hs.processDelegatedCredentialFromServer(
			certMsg.certificates[0].delegatedCredential,
			certVerifyMsg.signatureAlgorithm); err != nil {
			return err
		}
	}

	// Set the public key used to verify the handshake.
	pk := hs.c.peerCertificates[0].PublicKey

	// If the delegated credential extension has successfully been negotiated,
	// then the  CertificateVerify signature will have been produced with the
	// DelegatedCredential's private key.
	if hs.c.verifiedDc != nil {
		pk = hs.c.verifiedDc.cred.publicKey
	}

	// Verify the handshake signature.
	err, alertCode := verifyPeerHandshakeSignature(
		certVerifyMsg,
		pk,
		hs.hello.supportedSignatureAlgorithms,
		hs.keySchedule.transcriptHash.Sum(nil),
		"TLS 1.3, server CertificateVerify")
	if err != nil {
		c.sendAlert(alertCode)
		return err
	}
	hs.keySchedule.write(certVerifyMsg.marshal())

	// Receive Finished message.
	msg, err = c.readHandshake()
	if err != nil {
		return err
	}
	serverFinished, ok := msg.(*finishedMsg)
	if !ok {
		c.sendAlert(alertUnexpectedMessage)
		return unexpectedMessageError(serverFinished, msg)
	}
	// Validate server Finished hash.
	expectedVerifyData := hmacOfSum(hash, hs.keySchedule.transcriptHash, serverFinishedKey)
	if subtle.ConstantTimeCompare(expectedVerifyData, serverFinished.verifyData) != 1 {
		c.sendAlert(alertDecryptError)
		return errors.New("tls: server's Finished message is incorrect")
	}
	hs.keySchedule.write(serverFinished.marshal())

	// Server has authenticated itself. Calculate application traffic secrets.
	hs.keySchedule.setSecret(nil) // derive master secret

	// Change outbound handshake cipher for final step
	c.out.setKey(c.vers, hs.keySchedule.suite, clientHandshakeSecret)

	clientAppTrafficSecret := hs.keySchedule.deriveSecret(secretApplicationClient)
	serverAppTrafficSecret := hs.keySchedule.deriveSecret(secretApplicationServer)
	// TODO store initial traffic secret key for KeyUpdate GH #85

	// Client auth requires sending a (possibly empty) Certificate followed
	// by a CertificateVerify message (if there was an actual certificate).
	if isCertRequested {
		if err := hs.sendCertificate13(chainToSend, certReq); err != nil {
			return err
		}
	}

	// Send Finished
	verifyData := hmacOfSum(hash, hs.keySchedule.transcriptHash, clientFinishedKey)
	clientFinished := &finishedMsg{
		verifyData: verifyData,
	}
	if _, err := c.writeRecord(recordTypeHandshake, clientFinished.marshal()); err != nil {
		return err
	}

	// Handshake done, set application traffic secret
	// TODO store initial traffic secret key for KeyUpdate GH #85
	c.out.exportKey(hs.keySchedule.suite, clientAppTrafficSecret)
	c.out.setKey(c.vers, hs.keySchedule.suite, clientAppTrafficSecret)
	if c.hand.Len() > 0 {
		c.sendAlert(alertUnexpectedMessage)
		return errors.New("tls: unexpected data after handshake")
	}
	c.in.exportKey(hs.keySchedule.suite, serverAppTrafficSecret)
	c.in.setKey(c.vers, hs.keySchedule.suite, serverAppTrafficSecret)
	return nil
}

// supportedSigAlgorithmsCert iterates over schemes and filters out those algorithms
// which are not supported for certificate verification.
func supportedSigAlgorithmsCert(schemes []SignatureScheme) (ret []SignatureScheme) {
	for _, sig := range schemes {
		// X509 doesn't support PSS signatures
		if !signatureSchemeIsPSS(sig) {
			ret = append(ret, sig)
		}
	}
	return
}

// Functions below implement dhKex interface for different DH shared secret agreements

// KEX: P-256, P-384, P-512 KEX
func (kexNist) generate(c *Conn, groupId CurveID) (private []byte, ks keyShare, err error) {
	// never fails
	curve, _ := curveForCurveID(groupId)
	private, x, y, err := elliptic.GenerateKey(curve, c.config.rand())
	if err != nil {
		return nil, keyShare{}, err
	}
	ks.group = groupId
	ks.data = elliptic.Marshal(curve, x, y)
	return
}
func (kexNist) derive(c *Conn, ks keyShare, secretKey []byte) []byte {
	// never fails
	curve, _ := curveForCurveID(ks.group)
	x, y := elliptic.Unmarshal(curve, ks.data)
	if x == nil {
		return nil
	}
	x, _ = curve.ScalarMult(x, y, secretKey)
	xBytes := x.Bytes()
	curveSize := (curve.Params().BitSize + 8 - 1) >> 3
	if len(xBytes) == curveSize {
		return xBytes
	}
	buf := make([]byte, curveSize)
	copy(buf[len(buf)-len(xBytes):], xBytes)
	return buf
}

// KEX: X25519
func (kexX25519) generate(c *Conn, groupId CurveID) ([]byte, keyShare, error) {
	var scalar, public [x25519SharedSecretSz]byte
	if _, err := io.ReadFull(c.config.rand(), scalar[:]); err != nil {
		return nil, keyShare{}, err
	}
	curve25519.ScalarBaseMult(&public, &scalar)
	return scalar[:], keyShare{group: X25519, data: public[:]}, nil
}

func (kexX25519) derive(c *Conn, ks keyShare, secretKey []byte) []byte {
	var theirPublic, sharedKey, scalar [x25519SharedSecretSz]byte
	if len(ks.data) != x25519SharedSecretSz {
		return nil
	}
	copy(theirPublic[:], ks.data)
	copy(scalar[:], secretKey)
	curve25519.ScalarMult(&sharedKey, &scalar, &theirPublic)
	return sharedKey[:]
}

// KEX: SIDH/503
func (kexSIDHp503) generate(c *Conn, groupId CurveID) ([]byte, keyShare, error) {
	var variant, _ = getSidhKeyVariant(c.isClient)
	var prvKey = sidh.NewPrivateKey(sidh.FP_503, variant)
	if prvKey.Generate(c.config.rand()) != nil {
		return nil, keyShare{}, errors.New("tls: private SIDH key generation failed")
	}
	pubKey := prvKey.GeneratePublicKey()
	return prvKey.Export(), keyShare{group: 0 /*UNUSED*/, data: pubKey.Export()}, nil
}

func (kexSIDHp503) derive(c *Conn, ks keyShare, key []byte) []byte {
	var prvVariant, pubVariant = getSidhKeyVariant(c.isClient)
	var prvKeySize = P503PrvKeySz

	if len(ks.data) != P503PubKeySz || len(key) != prvKeySize {
		return nil
	}

	prvKey := sidh.NewPrivateKey(sidh.FP_503, prvVariant)
	pubKey := sidh.NewPublicKey(sidh.FP_503, pubVariant)

	if err := prvKey.Import(key); err != nil {
		return nil
	}
	if err := pubKey.Import(ks.data); err != nil {
		return nil
	}

	// Never fails
	sharedKey, _ := sidh.DeriveSecret(prvKey, pubKey)
	return sharedKey
}

// KEX Hybrid SIDH/503-X25519
func (kex *kexHybridSIDHp503X25519) generate(c *Conn, groupId CurveID) (private []byte, ks keyShare, err error) {
	var pubHybrid [SIDHp503Curve25519PubKeySz]byte
	var prvHybrid [SIDHp503Curve25519PrvKeySz]byte

	// Generate ephemeral key for classic x25519
	private, ks, err = kex.classicKEX.generate(c, groupId)
	if err != nil {
		return
	}
	copy(prvHybrid[:], private)
	copy(pubHybrid[:], ks.data)

	// Generate PQ ephemeral key for SIDH
	private, ks, err = kex.pqKEX.generate(c, groupId)
	if err != nil {
		return
	}
	copy(prvHybrid[x25519SharedSecretSz:], private)
	copy(pubHybrid[x25519SharedSecretSz:], ks.data)
	return prvHybrid[:], keyShare{group: HybridSIDHp503Curve25519, data: pubHybrid[:]}, nil
}

func (kex *kexHybridSIDHp503X25519) derive(c *Conn, ks keyShare, key []byte) []byte {
	var sharedKey [SIDHp503Curve25519SharedKeySz]byte
	var ret []byte
	var tmpKs keyShare

	// Key agreement for classic
	tmpKs.group = X25519
	tmpKs.data = ks.data[:x25519SharedSecretSz]
	ret = kex.classicKEX.derive(c, tmpKs, key[:x25519SharedSecretSz])
	if ret == nil {
		return nil
	}
	copy(sharedKey[:], ret)

	// Key agreement for PQ
	tmpKs.group = 0 /*UNUSED*/
	tmpKs.data = ks.data[x25519SharedSecretSz:]
	ret = kex.pqKEX.derive(c, tmpKs, key[x25519SharedSecretSz:])
	if ret == nil {
		return nil
	}
	copy(sharedKey[x25519SharedSecretSz:], ret)
	return sharedKey[:]
}
