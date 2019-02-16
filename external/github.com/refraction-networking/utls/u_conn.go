// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tls

import (
	"bufio"
	"bytes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync/atomic"
)

type UConn struct {
	*Conn

	Extensions    []TLSExtension
	clientHelloID ClientHelloID

	ClientHelloBuilt bool
	HandshakeState   ClientHandshakeState

	// sessionID may or may not depend on ticket; nil => random
	GetSessionID func(ticket []byte) [32]byte

	greaseSeed [ssl_grease_last_index]uint16
}

// UClient returns a new uTLS client, with behavior depending on clientHelloID.
// Config CAN be nil, but make sure to eventually specify ServerName.
func UClient(conn net.Conn, config *Config, clientHelloID ClientHelloID) *UConn {
	if config == nil {
		config = &Config{}
	}
	tlsConn := Conn{conn: conn, config: config, isClient: true}
	handshakeState := ClientHandshakeState{C: &tlsConn, Hello: &ClientHelloMsg{}}
	uconn := UConn{Conn: &tlsConn, clientHelloID: clientHelloID, HandshakeState: handshakeState}
	return &uconn
}

// BuildHandshakeState behavior varies based on ClientHelloID and
// whether it was already called before.
// If HelloGolang:
//   [only once] make default ClientHello and overwrite existing state
// If any other mimicking ClientHelloID is used:
//   [only once] make ClientHello based on ID and overwrite existing state
//   [each call] apply uconn.Extensions config to internal crypto/tls structures
//   [each call] marshal ClientHello.
//
// BuildHandshakeState is automatically called before uTLS performs handshake,
// amd should only be called explicitly to inspect/change fields of
// default/mimicked ClientHello.
func (uconn *UConn) BuildHandshakeState() error {
	if uconn.clientHelloID == HelloGolang {
		if uconn.ClientHelloBuilt {
			return nil
		}

		// use default Golang ClientHello.
		hello, ecdheParams, err := uconn.makeClientHello()
		if err != nil {
			return err
		}

		uconn.HandshakeState.Hello = hello.getPublicPtr()
		uconn.HandshakeState.State13.EcdheParams = ecdheParams
		uconn.HandshakeState.C = uconn.Conn
	} else {
		if !uconn.ClientHelloBuilt {
			err := uconn.applyPresetByID(uconn.clientHelloID)
			if err != nil {
				return err
			}
		}

		err := uconn.ApplyConfig()
		if err != nil {
			return err
		}
		err = uconn.MarshalClientHello()
		if err != nil {
			return err
		}
	}
	uconn.ClientHelloBuilt = true
	return nil
}

// SetSessionState sets the session ticket, which may be preshared or fake.
// If session is nil, the body of session ticket extension will be unset,
// but the extension itself still MAY be present for mimicking purposes.
// Session tickets to be reused - use same cache on following connections.
func (uconn *UConn) SetSessionState(session *ClientSessionState) error {
	uconn.HandshakeState.Session = session
	var sessionTicket []uint8
	if session != nil {
		sessionTicket = session.sessionTicket
	}
	uconn.HandshakeState.Hello.TicketSupported = true
	uconn.HandshakeState.Hello.SessionTicket = sessionTicket

	for _, ext := range uconn.Extensions {
		st, ok := ext.(*SessionTicketExtension)
		if !ok {
			continue
		}
		st.Session = session
		if session != nil {
			if len(session.SessionTicket()) > 0 {
				if uconn.GetSessionID != nil {
					sid := uconn.GetSessionID(session.SessionTicket())
					uconn.HandshakeState.Hello.SessionId = sid[:]
					return nil
				}
			}
			var sessionID [32]byte
			_, err := io.ReadFull(uconn.config.rand(), uconn.HandshakeState.Hello.SessionId)
			if err != nil {
				return err
			}
			uconn.HandshakeState.Hello.SessionId = sessionID[:]
		}
		return nil
	}
	return nil
}

// If you want session tickets to be reused - use same cache on following connections
func (uconn *UConn) SetSessionCache(cache ClientSessionCache) {
	uconn.config.ClientSessionCache = cache
	uconn.HandshakeState.Hello.TicketSupported = true
}

// SetClientRandom sets client random explicitly.
// BuildHandshakeFirst() must be called before SetClientRandom.
// r must to be 32 bytes long.
func (uconn *UConn) SetClientRandom(r []byte) error {
	if len(r) != 32 {
		return errors.New("Incorrect client random length! Expected: 32, got: " + strconv.Itoa(len(r)))
	} else {
		uconn.HandshakeState.Hello.Random = make([]byte, 32)
		copy(uconn.HandshakeState.Hello.Random, r)
		return nil
	}
}

func (uconn *UConn) SetSNI(sni string) {
	hname := hostnameInSNI(sni)
	uconn.config.ServerName = hname
	for _, ext := range uconn.Extensions {
		sniExt, ok := ext.(*SNIExtension)
		if ok {
			sniExt.ServerName = hname
		}
	}
}

// Handshake runs the client handshake using given clientHandshakeState
// Requires hs.hello, and, optionally, hs.session to be set.
func (c *UConn) Handshake() error {
	c.handshakeMutex.Lock()
	defer c.handshakeMutex.Unlock()

	if err := c.handshakeErr; err != nil {
		return err
	}
	if c.handshakeComplete() {
		return nil
	}

	c.in.Lock()
	defer c.in.Unlock()

	if c.isClient {
		// [uTLS section begins]
		err := c.BuildHandshakeState()
		if err != nil {
			return err
		}
		// [uTLS section ends]

		c.handshakeErr = c.clientHandshake()
	} else {
		c.handshakeErr = c.serverHandshake()
	}
	if c.handshakeErr == nil {
		c.handshakes++
	} else {
		// If an error occurred during the hadshake try to flush the
		// alert that might be left in the buffer.
		c.flush()
	}

	if c.handshakeErr == nil && !c.handshakeComplete() {
		c.handshakeErr = errors.New("tls: internal error: handshake should have had a result")
	}

	return c.handshakeErr
}

// Copy-pasted from tls.Conn in its entirety. But c.Handshake() is now utls' one, not tls.
// Write writes data to the connection.
func (c *UConn) Write(b []byte) (int, error) {
	// interlock with Close below
	for {
		x := atomic.LoadInt32(&c.activeCall)
		if x&1 != 0 {
			return 0, errClosed
		}
		if atomic.CompareAndSwapInt32(&c.activeCall, x, x+2) {
			defer atomic.AddInt32(&c.activeCall, -2)
			break
		}
	}

	if err := c.Handshake(); err != nil {
		return 0, err
	}

	c.out.Lock()
	defer c.out.Unlock()

	if err := c.out.err; err != nil {
		return 0, err
	}

	if !c.handshakeComplete() {
		return 0, alertInternalError
	}

	if c.closeNotifySent {
		return 0, errShutdown
	}

	// SSL 3.0 and TLS 1.0 are susceptible to a chosen-plaintext
	// attack when using block mode ciphers due to predictable IVs.
	// This can be prevented by splitting each Application Data
	// record into two records, effectively randomizing the IV.
	//
	// https://www.openssl.org/~bodo/tls-cbc.txt
	// https://bugzilla.mozilla.org/show_bug.cgi?id=665814
	// https://www.imperialviolet.org/2012/01/15/beastfollowup.html

	var m int
	if len(b) > 1 && c.vers <= VersionTLS10 {
		if _, ok := c.out.cipher.(cipher.BlockMode); ok {
			n, err := c.writeRecordLocked(recordTypeApplicationData, b[:1])
			if err != nil {
				return n, c.out.setErrorLocked(err)
			}
			m, b = 1, b[1:]
		}
	}

	n, err := c.writeRecordLocked(recordTypeApplicationData, b)
	return n + m, c.out.setErrorLocked(err)
}

// clientHandshakeWithOneState checks that exactly one expected state is set (1.2 or 1.3)
// and performs client TLS handshake with that state
func (c *UConn) clientHandshake() (err error) {
	// [uTLS section begins]
	hello := c.HandshakeState.Hello.getPrivatePtr()
	defer func() { c.HandshakeState.Hello = hello.getPublicPtr() }()

	sessionIsAlreadySet := c.HandshakeState.Session != nil

	// after this point exactly 1 out of 2 HandshakeState pointers is non-nil,
	// useTLS13 variable tells which pointer
	// [uTLS section ends]

	if c.config == nil {
		c.config = defaultConfig()
	}

	// This may be a renegotiation handshake, in which case some fields
	// need to be reset.
	c.didResume = false

	// [uTLS section begins]
	// don't make new ClientHello, use hs.hello
	// preserve the checks from beginning and end of makeClientHello()
	if len(c.config.ServerName) == 0 && !c.config.InsecureSkipVerify {
		return errors.New("tls: either ServerName or InsecureSkipVerify must be specified in the tls.Config")
	}

	nextProtosLength := 0
	for _, proto := range c.config.NextProtos {
		if l := len(proto); l == 0 || l > 255 {
			return errors.New("tls: invalid NextProtos value")
		} else {
			nextProtosLength += 1 + l
		}
	}

	if nextProtosLength > 0xffff {
		return errors.New("tls: NextProtos values too large")
	}

	if c.handshakes > 0 {
		hello.secureRenegotiation = c.clientFinished[:]
	}
	// [uTLS section ends]

	cacheKey, session, earlySecret, binderKey := c.loadSession(hello)
	if cacheKey != "" && session != nil {
		defer func() {
			// If we got a handshake failure when resuming a session, throw away
			// the session ticket. See RFC 5077, Section 3.2.
			//
			// RFC 8446 makes no mention of dropping tickets on failure, but it
			// does require servers to abort on invalid binders, so we need to
			// delete tickets to recover from a corrupted PSK.
			if err != nil {
				c.config.ClientSessionCache.Put(cacheKey, nil)
			}
		}()
	}

	if !sessionIsAlreadySet { // uTLS: do not overwrite already set session
		err = c.SetSessionState(session)
		if err != nil {
			return
		}
	}

	if _, err := c.writeRecord(recordTypeHandshake, hello.marshal()); err != nil {
		return err
	}

	msg, err := c.readHandshake()
	if err != nil {
		return err
	}

	serverHello, ok := msg.(*serverHelloMsg)
	if !ok {
		c.sendAlert(alertUnexpectedMessage)
		return unexpectedMessageError(serverHello, msg)
	}

	if err := c.pickTLSVersion(serverHello); err != nil {
		return err
	}

	// uTLS: do not create new handshakeState, use existing one
	if c.vers == VersionTLS13 {
		hs13 := c.HandshakeState.toPrivate13()
		hs13.serverHello = serverHello
		hs13.hello = hello
		if !sessionIsAlreadySet {
			hs13.earlySecret = earlySecret
			hs13.binderKey = binderKey
		}
		// In TLS 1.3, session tickets are delivered after the handshake.
		err = hs13.handshake()
		c.HandshakeState = *hs13.toPublic13()
		return err
	}

	hs12 := c.HandshakeState.toPrivate12()
	hs12.serverHello = serverHello
	hs12.hello = hello
	err = hs12.handshake()
	c.HandshakeState = *hs12.toPublic13()
	if err != nil {
		return err
	}

	// If we had a successful handshake and hs.session is different from
	// the one already cached - cache a new one.
	if cacheKey != "" && hs12.session != nil && session != hs12.session {
		c.config.ClientSessionCache.Put(cacheKey, hs12.session)
	}
	return nil
}

func (uconn *UConn) ApplyConfig() error {
	for _, ext := range uconn.Extensions {
		err := ext.writeToUConn(uconn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (uconn *UConn) MarshalClientHello() error {
	hello := uconn.HandshakeState.Hello
	headerLength := 2 + 32 + 1 + len(hello.SessionId) +
		2 + len(hello.CipherSuites)*2 +
		1 + len(hello.CompressionMethods)

	extensionsLen := 0
	var paddingExt *UtlsPaddingExtension
	for _, ext := range uconn.Extensions {
		if pe, ok := ext.(*UtlsPaddingExtension); !ok {
			// If not padding - just add length of extension to total length
			extensionsLen += ext.Len()
		} else {
			// If padding - process it later
			if paddingExt == nil {
				paddingExt = pe
			} else {
				return errors.New("Multiple padding extensions!")
			}
		}
	}

	if paddingExt != nil {
		// determine padding extension presence and length
		paddingExt.Update(headerLength + 4 + extensionsLen + 2)
		extensionsLen += paddingExt.Len()
	}

	helloLen := headerLength
	if len(uconn.Extensions) > 0 {
		helloLen += 2 + extensionsLen // 2 bytes for extensions' length
	}

	helloBuffer := bytes.Buffer{}
	bufferedWriter := bufio.NewWriterSize(&helloBuffer, helloLen+4) // 1 byte for tls record type, 3 for length
	// We use buffered Writer to avoid checking write errors after every Write(): whenever first error happens
	// Write() will become noop, and error will be accessible via Flush(), which is called once in the end

	binary.Write(bufferedWriter, binary.BigEndian, typeClientHello)
	helloLenBytes := []byte{byte(helloLen >> 16), byte(helloLen >> 8), byte(helloLen)} // poor man's uint24
	binary.Write(bufferedWriter, binary.BigEndian, helloLenBytes)
	binary.Write(bufferedWriter, binary.BigEndian, hello.Vers)

	binary.Write(bufferedWriter, binary.BigEndian, hello.Random)

	binary.Write(bufferedWriter, binary.BigEndian, uint8(len(hello.SessionId)))
	binary.Write(bufferedWriter, binary.BigEndian, hello.SessionId)

	binary.Write(bufferedWriter, binary.BigEndian, uint16(len(hello.CipherSuites)<<1))
	for _, suite := range hello.CipherSuites {
		binary.Write(bufferedWriter, binary.BigEndian, suite)
	}

	binary.Write(bufferedWriter, binary.BigEndian, uint8(len(hello.CompressionMethods)))
	binary.Write(bufferedWriter, binary.BigEndian, hello.CompressionMethods)

	if len(uconn.Extensions) > 0 {
		binary.Write(bufferedWriter, binary.BigEndian, uint16(extensionsLen))
		for _, ext := range uconn.Extensions {
			bufferedWriter.ReadFrom(ext)
		}
	}

	err := bufferedWriter.Flush()
	if err != nil {
		return err
	}

	if helloBuffer.Len() != 4+helloLen {
		return errors.New("utls: unexpected ClientHello length. Expected: " + strconv.Itoa(4+helloLen) +
			". Got: " + strconv.Itoa(helloBuffer.Len()))
	}

	hello.Raw = helloBuffer.Bytes()
	return nil
}

// get current state of cipher and encrypt zeros to get keystream
func (uconn *UConn) GetOutKeystream(length int) ([]byte, error) {
	zeros := make([]byte, length)

	if outCipher, ok := uconn.out.cipher.(cipher.AEAD); ok {
		// AEAD.Seal() does not mutate internal state, other ciphers might
		return outCipher.Seal(nil, uconn.out.seq[:], zeros, nil), nil
	}
	return nil, errors.New("Could not convert OutCipher to cipher.AEAD")
}

// SetVersCreateState set min and max TLS version in all appropriate places.
func (uconn *UConn) SetTLSVers(minTLSVers, maxTLSVers uint16) error {
	if minTLSVers < VersionTLS10 || minTLSVers > VersionTLS12 {
		return fmt.Errorf("uTLS does not support 0x%X as min version", minTLSVers)
	}

	if maxTLSVers < VersionTLS10 || maxTLSVers > VersionTLS13 {
		return fmt.Errorf("uTLS does not support 0x%X as max version", maxTLSVers)
	}

	uconn.HandshakeState.Hello.SupportedVersions = makeSupportedVersions(minTLSVers, maxTLSVers)
	uconn.config.MinVersion = minTLSVers
	uconn.config.MaxVersion = maxTLSVers

	return nil
}

func (uconn *UConn) SetUnderlyingConn(c net.Conn) {
	uconn.Conn.conn = c
}

func (uconn *UConn) GetUnderlyingConn() net.Conn {
	return uconn.Conn.conn
}

// MakeConnWithCompleteHandshake allows to forge both server and client side TLS connections.
// Major Hack Alert.
func MakeConnWithCompleteHandshake(tcpConn net.Conn, version uint16, cipherSuite uint16, masterSecret []byte, clientRandom []byte, serverRandom []byte, isClient bool) *Conn {
	tlsConn := &Conn{conn: tcpConn, config: &Config{}, isClient: isClient}
	cs := cipherSuiteByID(cipherSuite)

	// This is mostly borrowed from establishKeys()
	clientMAC, serverMAC, clientKey, serverKey, clientIV, serverIV :=
		keysFromMasterSecret(version, cs, masterSecret, clientRandom, serverRandom,
			cs.macLen, cs.keyLen, cs.ivLen)

	var clientCipher, serverCipher interface{}
	var clientHash, serverHash macFunction
	if cs.cipher != nil {
		clientCipher = cs.cipher(clientKey, clientIV, true /* for reading */)
		clientHash = cs.mac(version, clientMAC)
		serverCipher = cs.cipher(serverKey, serverIV, false /* not for reading */)
		serverHash = cs.mac(version, serverMAC)
	} else {
		clientCipher = cs.aead(clientKey, clientIV)
		serverCipher = cs.aead(serverKey, serverIV)
	}

	if isClient {
		tlsConn.in.prepareCipherSpec(version, serverCipher, serverHash)
		tlsConn.out.prepareCipherSpec(version, clientCipher, clientHash)
	} else {
		tlsConn.in.prepareCipherSpec(version, clientCipher, clientHash)
		tlsConn.out.prepareCipherSpec(version, serverCipher, serverHash)
	}

	// skip the handshake states
	tlsConn.handshakeStatus = 1
	tlsConn.cipherSuite = cipherSuite
	tlsConn.haveVers = true
	tlsConn.vers = version

	// Update to the new cipher specs
	// and consume the finished messages
	tlsConn.in.changeCipherSpec()
	tlsConn.out.changeCipherSpec()

	tlsConn.in.incSeq()
	tlsConn.out.incSeq()

	return tlsConn
}

func makeSupportedVersions(minVers, maxVers uint16) []uint16 {
	a := make([]uint16, maxVers-minVers+1)
	for i := range a {
		a[i] = maxVers - uint16(i)
	}
	return a
}
