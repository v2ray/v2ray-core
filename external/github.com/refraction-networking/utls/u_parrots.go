// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tls

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sort"
	"strconv"
	"time"
)

func utlsIdToSpec(id ClientHelloID) (ClientHelloSpec, error) {
	switch id {
	case HelloChrome_58, HelloChrome_62:
		return ClientHelloSpec{
			TLSVersMax: VersionTLS12,
			TLSVersMin: VersionTLS10,
			CipherSuites: []uint16{
				GREASE_PLACEHOLDER,
				TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				TLS_RSA_WITH_AES_128_GCM_SHA256,
				TLS_RSA_WITH_AES_256_GCM_SHA384,
				TLS_RSA_WITH_AES_128_CBC_SHA,
				TLS_RSA_WITH_AES_256_CBC_SHA,
				TLS_RSA_WITH_3DES_EDE_CBC_SHA,
			},
			CompressionMethods: []byte{compressionNone},
			Extensions: []TLSExtension{
				&UtlsGREASEExtension{},
				&RenegotiationInfoExtension{renegotiation: RenegotiateOnceAsClient},
				&SNIExtension{},
				&UtlsExtendedMasterSecretExtension{},
				&SessionTicketExtension{},
				&SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []SignatureScheme{
					ECDSAWithP256AndSHA256,
					PSSWithSHA256,
					PKCS1WithSHA256,
					ECDSAWithP384AndSHA384,
					PSSWithSHA384,
					PKCS1WithSHA384,
					PSSWithSHA512,
					PKCS1WithSHA512,
					PKCS1WithSHA1},
				},
				&StatusRequestExtension{},
				&SCTExtension{},
				&ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&FakeChannelIDExtension{},
				&SupportedPointsExtension{SupportedPoints: []byte{pointFormatUncompressed}},
				&SupportedCurvesExtension{[]CurveID{CurveID(GREASE_PLACEHOLDER),
					X25519, CurveP256, CurveP384}},
				&UtlsGREASEExtension{},
				&UtlsPaddingExtension{GetPaddingLen: BoringPaddingStyle},
			},
			GetSessionID: sha256.Sum256,
		}, nil
	case HelloChrome_70:
		return ClientHelloSpec{
			TLSVersMin: VersionTLS10,
			TLSVersMax: VersionTLS13,
			CipherSuites: []uint16{
				GREASE_PLACEHOLDER,
				TLS_AES_128_GCM_SHA256,
				TLS_AES_256_GCM_SHA384,
				TLS_CHACHA20_POLY1305_SHA256,
				TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				TLS_RSA_WITH_AES_128_GCM_SHA256,
				TLS_RSA_WITH_AES_256_GCM_SHA384,
				TLS_RSA_WITH_AES_128_CBC_SHA,
				TLS_RSA_WITH_AES_256_CBC_SHA,
				TLS_RSA_WITH_3DES_EDE_CBC_SHA,
			},
			CompressionMethods: []byte{
				compressionNone,
			},
			Extensions: []TLSExtension{
				&UtlsGREASEExtension{},
				&RenegotiationInfoExtension{renegotiation: RenegotiateOnceAsClient},
				&SNIExtension{},
				&UtlsExtendedMasterSecretExtension{},
				&SessionTicketExtension{},
				&SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []SignatureScheme{
					ECDSAWithP256AndSHA256,
					PSSWithSHA256,
					PKCS1WithSHA256,
					ECDSAWithP384AndSHA384,
					PSSWithSHA384,
					PKCS1WithSHA384,
					PSSWithSHA512,
					PKCS1WithSHA512,
					PKCS1WithSHA1,
				}},
				&StatusRequestExtension{},
				&SCTExtension{},
				&ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&FakeChannelIDExtension{},
				&SupportedPointsExtension{SupportedPoints: []byte{
					pointFormatUncompressed,
				}},
				&KeyShareExtension{[]KeyShare{
					{Group: CurveID(GREASE_PLACEHOLDER), Data: []byte{0}},
					{Group: X25519},
				}},
				&PSKKeyExchangeModesExtension{[]uint8{pskModeDHE}},
				&SupportedVersionsExtension{[]uint16{
					GREASE_PLACEHOLDER,
					VersionTLS13,
					VersionTLS12,
					VersionTLS11,
					VersionTLS10}},
				&SupportedCurvesExtension{[]CurveID{
					CurveID(GREASE_PLACEHOLDER),
					X25519,
					CurveP256,
					CurveP384,
				}},
				&GenericExtension{id: fakeCertCompressionAlgs, data: []byte{02, 00, 02}},
				&UtlsGREASEExtension{},
				&UtlsPaddingExtension{GetPaddingLen: BoringPaddingStyle},
			},
		}, nil
	case HelloFirefox_55, HelloFirefox_56:
		return ClientHelloSpec{
			TLSVersMax: VersionTLS12,
			TLSVersMin: VersionTLS10,
			CipherSuites: []uint16{
				TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				FAKE_TLS_DHE_RSA_WITH_AES_128_CBC_SHA,
				FAKE_TLS_DHE_RSA_WITH_AES_256_CBC_SHA,
				TLS_RSA_WITH_AES_128_CBC_SHA,
				TLS_RSA_WITH_AES_256_CBC_SHA,
				TLS_RSA_WITH_3DES_EDE_CBC_SHA,
			},
			CompressionMethods: []byte{compressionNone},
			Extensions: []TLSExtension{
				&SNIExtension{},
				&UtlsExtendedMasterSecretExtension{},
				&RenegotiationInfoExtension{renegotiation: RenegotiateOnceAsClient},
				&SupportedCurvesExtension{[]CurveID{X25519, CurveP256, CurveP384, CurveP521}},
				&SupportedPointsExtension{SupportedPoints: []byte{pointFormatUncompressed}},
				&SessionTicketExtension{},
				&ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&StatusRequestExtension{},
				&SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []SignatureScheme{
					ECDSAWithP256AndSHA256,
					ECDSAWithP384AndSHA384,
					ECDSAWithP521AndSHA512,
					PSSWithSHA256,
					PSSWithSHA384,
					PSSWithSHA512,
					PKCS1WithSHA256,
					PKCS1WithSHA384,
					PKCS1WithSHA512,
					ECDSAWithSHA1,
					PKCS1WithSHA1},
				},
				&UtlsPaddingExtension{GetPaddingLen: BoringPaddingStyle},
			},
			GetSessionID: nil,
		}, nil
	case HelloFirefox_63:
		return ClientHelloSpec{
			TLSVersMin: VersionTLS10,
			TLSVersMax: VersionTLS13,
			CipherSuites: []uint16{
				TLS_AES_128_GCM_SHA256,
				TLS_CHACHA20_POLY1305_SHA256,
				TLS_AES_256_GCM_SHA384,
				TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				FAKE_TLS_DHE_RSA_WITH_AES_128_CBC_SHA,
				FAKE_TLS_DHE_RSA_WITH_AES_256_CBC_SHA,
				TLS_RSA_WITH_AES_128_CBC_SHA,
				TLS_RSA_WITH_AES_256_CBC_SHA,
				TLS_RSA_WITH_3DES_EDE_CBC_SHA,
			},
			CompressionMethods: []byte{
				compressionNone,
			},
			Extensions: []TLSExtension{
				&SNIExtension{},
				&UtlsExtendedMasterSecretExtension{},
				&RenegotiationInfoExtension{renegotiation: RenegotiateOnceAsClient},
				&SupportedCurvesExtension{[]CurveID{
					X25519,
					CurveP256,
					CurveP384,
					CurveP521,
					CurveID(FakeFFDHE2048),
					CurveID(FakeFFDHE3072),
				}},
				&SupportedPointsExtension{SupportedPoints: []byte{
					pointFormatUncompressed,
				}},
				&SessionTicketExtension{},
				&ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&StatusRequestExtension{},
				&KeyShareExtension{[]KeyShare{
					{Group: X25519},
					{Group: CurveP256},
				}},
				&SupportedVersionsExtension{[]uint16{
					VersionTLS13,
					VersionTLS12,
					VersionTLS11,
					VersionTLS10}},
				&SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []SignatureScheme{
					ECDSAWithP256AndSHA256,
					ECDSAWithP384AndSHA384,
					ECDSAWithP521AndSHA512,
					PSSWithSHA256,
					PSSWithSHA384,
					PSSWithSHA512,
					PKCS1WithSHA256,
					PKCS1WithSHA384,
					PKCS1WithSHA512,
					ECDSAWithSHA1,
					PKCS1WithSHA1,
				}},
				&PSKKeyExchangeModesExtension{[]uint8{pskModeDHE}},
				&GenericExtension{id: fakeRecordSizeLimit, data: []byte{0x40, 0x01}},
				&UtlsPaddingExtension{GetPaddingLen: BoringPaddingStyle},
			}}, nil
	case HelloIOS_11_1:
		return ClientHelloSpec{
			TLSVersMax: VersionTLS12,
			TLSVersMin: VersionTLS10,
			CipherSuites: []uint16{
				TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				DISABLED_TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384,
				TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
				TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				DISABLED_TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384,
				TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
				TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				TLS_RSA_WITH_AES_256_GCM_SHA384,
				TLS_RSA_WITH_AES_128_GCM_SHA256,
				DISABLED_TLS_RSA_WITH_AES_256_CBC_SHA256,
				TLS_RSA_WITH_AES_128_CBC_SHA256,
				TLS_RSA_WITH_AES_256_CBC_SHA,
				TLS_RSA_WITH_AES_128_CBC_SHA,
			},
			CompressionMethods: []byte{
				compressionNone,
			},
			Extensions: []TLSExtension{
				&RenegotiationInfoExtension{renegotiation: RenegotiateOnceAsClient},
				&SNIExtension{},
				&UtlsExtendedMasterSecretExtension{},
				&SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []SignatureScheme{
					ECDSAWithP256AndSHA256,
					PSSWithSHA256,
					PKCS1WithSHA256,
					ECDSAWithP384AndSHA384,
					PSSWithSHA384,
					PKCS1WithSHA384,
					PSSWithSHA512,
					PKCS1WithSHA512,
					PKCS1WithSHA1,
				}},
				&StatusRequestExtension{},
				&NPNExtension{},
				&SCTExtension{},
				&ALPNExtension{AlpnProtocols: []string{"h2", "h2-16", "h2-15", "h2-14", "spdy/3.1", "spdy/3", "http/1.1"}},
				&SupportedPointsExtension{SupportedPoints: []byte{
					pointFormatUncompressed,
				}},
				&SupportedCurvesExtension{Curves: []CurveID{
					X25519,
					CurveP256,
					CurveP384,
					CurveP521,
				}},
			},
		}, nil
	default:
		return ClientHelloSpec{}, errors.New("ClientHello ID " + id.Str() + " is unknown")
	}
}

func (uconn *UConn) applyPresetByID(id ClientHelloID) (err error) {
	var spec ClientHelloSpec
	// choose/generate the spec
	switch id {
	case HelloRandomized:
		if tossBiasedCoin(0.5) {
			return uconn.applyPresetByID(HelloRandomizedALPN)
		} else {
			return uconn.applyPresetByID(HelloRandomizedNoALPN)
		}
	case HelloRandomizedALPN:
		spec, err = uconn.generateRandomizedSpec(true)
		if err != nil {
			return err
		}
	case HelloRandomizedNoALPN:
		spec, err = uconn.generateRandomizedSpec(false)
		if err != nil {
			return err
		}
	case HelloCustom:
		return nil

	default:
		spec, err = utlsIdToSpec(id)
		if err != nil {
			return err
		}
	}

	uconn.clientHelloID = id
	return uconn.ApplyPreset(&spec)
}

// ApplyPreset should only be used in conjunction with HelloCustom to apply custom specs.
// Fields of TLSExtensions that are slices/pointers are shared across different connections with
// same ClientHelloSpec. It is advised to use different specs and avoid any shared state.
func (uconn *UConn) ApplyPreset(p *ClientHelloSpec) error {
	var err error
	err = uconn.SetTLSVers(p.TLSVersMin, p.TLSVersMax)
	if err != nil {
		return err
	}

	privateHello, ecdheParams, err := uconn.makeClientHello()
	if err != nil {
		return err
	}
	uconn.HandshakeState.Hello = privateHello.getPublicPtr()
	uconn.HandshakeState.State13.EcdheParams = ecdheParams
	hello := uconn.HandshakeState.Hello
	session := uconn.HandshakeState.Session

	switch len(hello.Random) {
	case 0:
		hello.Random = make([]byte, 32)
		_, err := io.ReadFull(uconn.config.rand(), hello.Random)
		if err != nil {
			return errors.New("tls: short read from Rand: " + err.Error())
		}
	case 32:
	// carry on
	default:
		return errors.New("ClientHello expected length: 32 bytes. Got: " +
			strconv.Itoa(len(hello.Random)) + " bytes")
	}
	if len(hello.CipherSuites) == 0 {
		hello.CipherSuites = defaultCipherSuites()
	}
	if len(hello.CompressionMethods) == 0 {
		hello.CompressionMethods = []uint8{compressionNone}
	}

	// Currently, GREASE is assumed to come from BoringSSL
	grease_bytes := make([]byte, 2*ssl_grease_last_index)
	grease_extensions_seen := 0
	_, err = io.ReadFull(uconn.config.rand(), grease_bytes)
	if err != nil {
		return errors.New("tls: short read from Rand: " + err.Error())
	}
	for i := range uconn.greaseSeed {
		uconn.greaseSeed[i] = binary.LittleEndian.Uint16(grease_bytes[2*i : 2*i+2])
	}
	if uconn.greaseSeed[ssl_grease_extension1] == uconn.greaseSeed[ssl_grease_extension2] {
		uconn.greaseSeed[ssl_grease_extension2] ^= 0x1010
	}

	hello.CipherSuites = make([]uint16, len(p.CipherSuites))
	copy(hello.CipherSuites, p.CipherSuites)
	for i := range hello.CipherSuites {
		if hello.CipherSuites[i] == GREASE_PLACEHOLDER {
			hello.CipherSuites[i] = GetBoringGREASEValue(uconn.greaseSeed, ssl_grease_cipher)
		}
	}
	uconn.GetSessionID = p.GetSessionID
	uconn.Extensions = make([]TLSExtension, len(p.Extensions))
	copy(uconn.Extensions, p.Extensions)

	// reGrease, and point things to each other
	for _, e := range uconn.Extensions {
		switch ext := e.(type) {
		case *SNIExtension:
			if ext.ServerName == "" {
				ext.ServerName = uconn.config.ServerName
			}
		case *UtlsGREASEExtension:
			switch grease_extensions_seen {
			case 0:
				ext.Value = GetBoringGREASEValue(uconn.greaseSeed, ssl_grease_extension1)
			case 1:
				ext.Value = GetBoringGREASEValue(uconn.greaseSeed, ssl_grease_extension2)
				ext.Body = []byte{0}
			default:
				return errors.New("at most 2 grease extensions are supported")
			}
			grease_extensions_seen += 1
		case *SessionTicketExtension:
			err := uconn.SetSessionState(session)
			if err != nil {
				return err
			}
		case *SupportedCurvesExtension:
			for i := range ext.Curves {
				if ext.Curves[i] == GREASE_PLACEHOLDER {
					ext.Curves[i] = CurveID(GetBoringGREASEValue(uconn.greaseSeed, ssl_grease_group))
				}
			}
		case *KeyShareExtension:
			preferredCurveIsSet := false
			for i := range ext.KeyShares {
				curveID := ext.KeyShares[i].Group
				if curveID == GREASE_PLACEHOLDER {
					ext.KeyShares[i].Group = CurveID(GetBoringGREASEValue(uconn.greaseSeed, ssl_grease_group))
					continue
				}
				if len(ext.KeyShares[i].Data) > 1 {
					continue
				}

				ecdheParams, err := generateECDHEParameters(uconn.config.rand(), curveID)
				if err != nil {
					return fmt.Errorf("unsupported Curve in KeyShareExtension: %v."+
						"To mimic it, fill the Data(key) field manually.", curveID)
				}
				ext.KeyShares[i].Data = ecdheParams.PublicKey()
				if !preferredCurveIsSet {
					// only do this once for the first non-grease curve
					uconn.HandshakeState.State13.EcdheParams = ecdheParams
					preferredCurveIsSet = true
				}
			}
		case *SupportedVersionsExtension:
			for i := range ext.Versions {
				if ext.Versions[i] == GREASE_PLACEHOLDER {
					ext.Versions[i] = GetBoringGREASEValue(uconn.greaseSeed, ssl_grease_version)
				}
			}
		}
	}
	return nil
}

func (uconn *UConn) generateRandomizedSpec(WithALPN bool) (ClientHelloSpec, error) {
	p := ClientHelloSpec{}

	p.CipherSuites = make([]uint16, len(defaultCipherSuites()))
	copy(p.CipherSuites, defaultCipherSuites())
	shuffledSuites, err := shuffledCiphers()
	if err != nil {
		return p, err
	}

	if tossBiasedCoin(0.4) {
		p.TLSVersMin = VersionTLS10
		p.TLSVersMax = VersionTLS13
		tls13ciphers := defaultCipherSuitesTLS13()
		err = shuffleUInts16(tls13ciphers)
		if err != nil {
			return p, err
		}
		// appending TLS 1.3 ciphers before TLS 1.2, since that's what popular implementations do
		shuffledSuites = append(tls13ciphers, shuffledSuites...)

		// TLS 1.3 forbids RC4 in any configurations
		shuffledSuites = removeRC4Ciphers(shuffledSuites)
	} else {
		p.TLSVersMin = VersionTLS10
		p.TLSVersMax = VersionTLS12
	}

	p.CipherSuites = removeRandomCiphers(shuffledSuites, 0.4)

	sni := SNIExtension{uconn.config.ServerName}
	sessionTicket := SessionTicketExtension{Session: uconn.HandshakeState.Session}

	sigAndHashAlgos := []SignatureScheme{
		ECDSAWithP256AndSHA256,
		PKCS1WithSHA256,
		ECDSAWithP384AndSHA384,
		PKCS1WithSHA384,
		PKCS1WithSHA1,
		PKCS1WithSHA512,
	}

	if tossBiasedCoin(0.63) {
		sigAndHashAlgos = append(sigAndHashAlgos, ECDSAWithSHA1)
	}
	if tossBiasedCoin(0.59) {
		sigAndHashAlgos = append(sigAndHashAlgos, ECDSAWithP521AndSHA512)
	}
	if tossBiasedCoin(0.51) || p.TLSVersMax == VersionTLS13 {
		// https://tools.ietf.org/html/rfc8446 says "...RSASSA-PSS (which is mandatory in TLS 1.3)..."
		sigAndHashAlgos = append(sigAndHashAlgos, PSSWithSHA256)
		if tossBiasedCoin(0.9) {
			// these usually go together
			sigAndHashAlgos = append(sigAndHashAlgos, PSSWithSHA384)
			sigAndHashAlgos = append(sigAndHashAlgos, PSSWithSHA512)
		}
	}

	err = shuffleSignatures(sigAndHashAlgos)
	if err != nil {
		return p, err
	}
	sigAndHash := SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: sigAndHashAlgos}

	status := StatusRequestExtension{}
	sct := SCTExtension{}
	ems := UtlsExtendedMasterSecretExtension{}
	points := SupportedPointsExtension{SupportedPoints: []byte{pointFormatUncompressed}}

	curveIDs := []CurveID{}
	if tossBiasedCoin(0.71) || p.TLSVersMax == VersionTLS13 {
		curveIDs = append(curveIDs, X25519)
	}
	curveIDs = append(curveIDs, CurveP256, CurveP384)
	if tossBiasedCoin(0.46) {
		curveIDs = append(curveIDs, CurveP521)
	}

	curves := SupportedCurvesExtension{curveIDs}

	padding := UtlsPaddingExtension{GetPaddingLen: BoringPaddingStyle}
	reneg := RenegotiationInfoExtension{renegotiation: RenegotiateOnceAsClient}

	p.Extensions = []TLSExtension{
		&sni,
		&sessionTicket,
		&sigAndHash,
		&points,
		&curves,
	}

	if WithALPN {
		if len(uconn.config.NextProtos) == 0 {
			// if user didn't specify alpn yet, choose something popular
			uconn.config.NextProtos = []string{"h2", "http/1.1"}
		}
		alpn := ALPNExtension{AlpnProtocols: uconn.config.NextProtos}
		p.Extensions = append(p.Extensions, &alpn)
	}

	if tossBiasedCoin(0.62) || p.TLSVersMax == VersionTLS13 {
		// always include for TLS 1.3, since TLS 1.3 ClientHellos are often over 256 bytes
		// and that's when padding is required to work around buggy middleboxes
		p.Extensions = append(p.Extensions, &padding)
	}
	if tossBiasedCoin(0.74) {
		p.Extensions = append(p.Extensions, &status)
	}
	if tossBiasedCoin(0.46) {
		p.Extensions = append(p.Extensions, &sct)
	}
	if tossBiasedCoin(0.75) {
		p.Extensions = append(p.Extensions, &reneg)
	}
	if tossBiasedCoin(0.77) {
		p.Extensions = append(p.Extensions, &ems)
	}
	if p.TLSVersMax == VersionTLS13 {
		ks := KeyShareExtension{[]KeyShare{
			{Group: X25519}, // the key for the group will be generated later
		}}
		if tossBiasedCoin(0.25) {
			// do not ADD second keyShare because crypto/tls does not support multiple ecdheParams
			// TODO: add it back when they implement multiple keyShares, or implement it oursevles
			// ks.KeyShares = append(ks.KeyShares, KeyShare{Group: CurveP256})
			ks.KeyShares[0].Group = CurveP256
		}
		pskExchangeModes := PSKKeyExchangeModesExtension{[]uint8{pskModeDHE}}
		supportedVersionsExt := SupportedVersionsExtension{
			Versions: makeSupportedVersions(p.TLSVersMin, p.TLSVersMax),
		}
		p.Extensions = append(p.Extensions, &ks, &pskExchangeModes, &supportedVersionsExt)
	}
	err = shuffleTLSExtensions(p.Extensions)
	if err != nil {
		return p, err
	}
	err = uconn.SetTLSVers(p.TLSVersMin, p.TLSVersMax)
	if err != nil {
		return p, err
	}

	return p, nil
}

func tossBiasedCoin(probability float32) bool {
	// probability is expected to be in [0,1]
	// this function never returns errors for ease of use
	const precision = 0xffff
	threshold := float32(precision) * probability
	value, err := getRandInt(precision)
	if err != nil {
		// I doubt that this code will ever actually be used, as other functions are expected to complain
		// about used source of entropy. Nonetheless, this is more than enough for given purpose
		return ((time.Now().Unix() & 1) == 0)
	}

	if float32(value) <= threshold {
		return true
	} else {
		return false
	}
}

func removeRandomCiphers(s []uint16, maxRemovalProbability float32) []uint16 {
	// removes elements in place
	// probability to remove increases for further elements
	// never remove first cipher
	if len(s) <= 1 {
		return s
	}

	// remove random elements
	floatLen := float32(len(s))
	sliceLen := len(s)
	for i := 1; i < sliceLen; i++ {
		if tossBiasedCoin(maxRemovalProbability * float32(i) / floatLen) {
			s = append(s[:i], s[i+1:]...)
			sliceLen--
			i--
		}
	}
	return s[:sliceLen]
}

func removeRC4Ciphers(s []uint16) []uint16 {
	// removes elements in place
	sliceLen := len(s)
	for i := 0; i < sliceLen; i++ {
		cipher := s[i]
		if cipher == TLS_ECDHE_ECDSA_WITH_RC4_128_SHA ||
			cipher == TLS_ECDHE_RSA_WITH_RC4_128_SHA ||
			cipher == TLS_RSA_WITH_RC4_128_SHA {
			s = append(s[:i], s[i+1:]...)
			sliceLen--
			i--
		}
	}
	return s[:sliceLen]
}

func getRandInt(max int) (int, error) {
	bigInt, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(bigInt.Int64()), err
}

func getRandPerm(n int) ([]int, error) {
	permArray := make([]int, n)
	for i := 1; i < n; i++ {
		j, err := getRandInt(i + 1)
		if err != nil {
			return permArray, err
		}
		permArray[i] = permArray[j]
		permArray[j] = i
	}
	return permArray, nil
}

func shuffledCiphers() ([]uint16, error) {
	ciphers := make(sortableCiphers, len(cipherSuites))
	perm, err := getRandPerm(len(cipherSuites))
	if err != nil {
		return nil, err
	}
	for i, suite := range cipherSuites {
		ciphers[i] = sortableCipher{suite: suite.id,
			isObsolete: ((suite.flags & suiteTLS12) == 0),
			randomTag:  perm[i]}
	}
	sort.Sort(ciphers)
	return ciphers.GetCiphers(), nil
}

type sortableCipher struct {
	isObsolete bool
	randomTag  int
	suite      uint16
}

type sortableCiphers []sortableCipher

func (ciphers sortableCiphers) Len() int {
	return len(ciphers)
}

func (ciphers sortableCiphers) Less(i, j int) bool {
	if ciphers[i].isObsolete && !ciphers[j].isObsolete {
		return false
	}
	if ciphers[j].isObsolete && !ciphers[i].isObsolete {
		return true
	}
	return ciphers[i].randomTag < ciphers[j].randomTag
}

func (ciphers sortableCiphers) Swap(i, j int) {
	ciphers[i], ciphers[j] = ciphers[j], ciphers[i]
}

func (ciphers sortableCiphers) GetCiphers() []uint16 {
	cipherIDs := make([]uint16, len(ciphers))
	for i := range ciphers {
		cipherIDs[i] = ciphers[i].suite
	}
	return cipherIDs
}

// so much for generics
func shuffleTLSExtensions(s []TLSExtension) error {
	// shuffles array in place
	perm, err := getRandPerm(len(s))
	if err != nil {
		return err
	}
	for i := range s {
		s[i], s[perm[i]] = s[perm[i]], s[i]
	}
	return nil
}

// so much for generics
func shuffleSignatures(s []SignatureScheme) error {
	// shuffles array in place
	perm, err := getRandPerm(len(s))
	if err != nil {
		return err
	}
	for i := range s {
		s[i], s[perm[i]] = s[perm[i]], s[i]
	}
	return nil
}

// so much for generics
func shuffleUInts16(s []uint16) error {
	// shuffles array in place
	perm, err := getRandPerm(len(s))
	if err != nil {
		return err
	}
	for i := range s {
		s[i], s[perm[i]] = s[perm[i]], s[i]
	}
	return nil
}
