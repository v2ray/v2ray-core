// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qtls

// Delegated credentials for TLS
// (https://tools.ietf.org/html/draft-ietf-tls-subcerts-02) is an IETF Internet
// draft and proposed TLS extension. This allows a backend server to delegate
// TLS termination to a trusted frontend. If the client supports this extension,
// then the frontend may use a "delegated credential" as the signing key in the
// handshake. A delegated credential is a short lived key pair delegated to the
// server by an entity trusted by the client. Once issued, credentials can't be
// revoked; in order to mitigate risk in case the frontend is compromised, the
// credential is only valid for a short time (days, hours, or even minutes).
//
// This implements draft 02. This draft doesn't specify an object identifier for
// the X.509 extension; we use one assigned by Cloudflare. In addition, IANA has
// not assigned an extension ID for this extension; we picked up one that's not
// yet taken.
//
// TODO(cjpatton) Only ECDSA is supported with delegated credentials for now;
// we'd like to suppoort for EcDSA signatures once these have better support
// upstream.

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/asn1"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

const (
	// length of the public key field
	dcPubKeyFieldLen  = 3
	dcMaxTTLSeconds   = 60 * 60 * 24 * 7 // 7 days
	dcMaxTTL          = time.Duration(dcMaxTTLSeconds * time.Second)
	dcMaxPublicKeyLen = 1 << 24 // Bytes
	dcMaxSignatureLen = 1 << 16 // Bytes
)

var errNoDelegationUsage = errors.New("certificate not authorized for delegation")

// delegationUsageId is the DelegationUsage X.509 extension OID
//
// NOTE(cjpatton) This OID is a child of Cloudflare's IANA-assigned OID.
var delegationUsageId = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 44363, 44}

// canDelegate returns true if a certificate can be used for delegated
// credentials.
func canDelegate(cert *x509.Certificate) bool {
	// Check that the digitalSignature key usage is set.
	if (cert.KeyUsage & x509.KeyUsageDigitalSignature) == 0 {
		return false
	}

	// Check that the certificate has the DelegationUsage extension and that
	// it's non-critical (per the spec).
	for _, extension := range cert.Extensions {
		if extension.Id.Equal(delegationUsageId) {
			return true
		}
	}
	return false
}

// credential stores the public components of a credential.
type credential struct {
	// The serialized form of the credential.
	raw []byte

	// The amount of time for which the credential is valid. Specifically, the
	// the credential expires `ValidTime` seconds after the `notBefore` of the
	// delegation certificate. The delegator shall not issue delegated
	// credentials that are valid for more than 7 days from the current time.
	//
	// When this data structure is serialized, this value is converted to a
	// uint32 representing the duration in seconds.
	validTime time.Duration

	// The signature scheme associated with the delegated credential public key.
	expectedCertVerifyAlgorithm SignatureScheme

	// The version of TLS in which the credential will be used.
	expectedVersion uint16

	// The credential public key.
	publicKey crypto.PublicKey
}

// isExpired returns true if the credential has expired. The end of the validity
// interval is defined as the delegator certificate's notBefore field (`start`)
// plus ValidTime seconds. This function simply checks that the current time
// (`now`) is before the end of the valdity interval.
func (cred *credential) isExpired(start, now time.Time) bool {
	end := start.Add(cred.validTime)
	return !now.Before(end)
}

// invalidTTL returns true if the credential's validity period is longer than the
// maximum permitted. This is defined by the certificate's notBefore field
// (`start`) plus the ValidTime, minus the current time (`now`).
func (cred *credential) invalidTTL(start, now time.Time) bool {
	return cred.validTime > (now.Sub(start) + dcMaxTTL).Round(time.Second)
}

// marshalSubjectPublicKeyInfo returns a DER encoded SubjectPublicKeyInfo structure
// (as defined in the X.509 standard) for the credential.
func (cred *credential) marshalSubjectPublicKeyInfo() ([]byte, error) {
	switch cred.expectedCertVerifyAlgorithm {
	case ECDSAWithP256AndSHA256,
		ECDSAWithP384AndSHA384,
		ECDSAWithP521AndSHA512:
		serializedPublicKey, err := x509.MarshalPKIXPublicKey(cred.publicKey)
		if err != nil {
			return nil, err
		}
		return serializedPublicKey, nil

	default:
		return nil, fmt.Errorf("unsupported signature scheme: 0x%04x", cred.expectedCertVerifyAlgorithm)
	}
}

// marshal encodes a credential in the wire format specified in
// https://tools.ietf.org/html/draft-ietf-tls-subcerts-02.
func (cred *credential) marshal() ([]byte, error) {
	// The number of bytes comprising the DC parameters, which includes the
	// validity time (4 bytes), the signature scheme of the public key (2 bytes), and
	// the protocol version (2 bytes).
	paramsLen := 8

	// The first 4 bytes are the valid_time, scheme, and version fields.
	serialized := make([]byte, paramsLen+dcPubKeyFieldLen)
	binary.BigEndian.PutUint32(serialized, uint32(cred.validTime/time.Second))
	binary.BigEndian.PutUint16(serialized[4:], uint16(cred.expectedCertVerifyAlgorithm))
	binary.BigEndian.PutUint16(serialized[6:], cred.expectedVersion)

	// Encode the public key and assert that the encoding is no longer than 2^16
	// bytes (per the spec).
	serializedPublicKey, err := cred.marshalSubjectPublicKeyInfo()
	if err != nil {
		return nil, err
	}
	if len(serializedPublicKey) > dcMaxPublicKeyLen {
		return nil, errors.New("public key is too long")
	}

	// The next 3 bytes are the length of the public key field, which may be up
	// to 2^24 bytes long.
	putUint24(serialized[paramsLen:], len(serializedPublicKey))

	// The remaining bytes are the public key itself.
	serialized = append(serialized, serializedPublicKey...)
	cred.raw = serialized
	return serialized, nil
}

// unmarshalCredential decodes a credential and returns it.
func unmarshalCredential(serialized []byte) (*credential, error) {
	// The number of bytes comprising the DC parameters.
	paramsLen := 8

	if len(serialized) < paramsLen+dcPubKeyFieldLen {
		return nil, errors.New("credential is too short")
	}

	// Parse the valid_time, scheme, and version fields.
	validTime := time.Duration(binary.BigEndian.Uint32(serialized)) * time.Second
	scheme := SignatureScheme(binary.BigEndian.Uint16(serialized[4:]))
	version := binary.BigEndian.Uint16(serialized[6:])

	// Parse the SubjectPublicKeyInfo.
	pk, err := x509.ParsePKIXPublicKey(serialized[paramsLen+dcPubKeyFieldLen:])
	if err != nil {
		return nil, err
	}

	if _, ok := pk.(*ecdsa.PublicKey); !ok {
		return nil, fmt.Errorf("unsupported delegation key type: %T", pk)
	}

	return &credential{
		raw:                         serialized,
		validTime:                   validTime,
		expectedCertVerifyAlgorithm: scheme,
		expectedVersion:             version,
		publicKey:                   pk,
	}, nil
}

// getCredentialLen returns the number of bytes comprising the serialized
// credential that starts at the beginning of the input slice. It returns an
// error if the input is too short to contain a credential.
func getCredentialLen(serialized []byte) (int, error) {
	paramsLen := 8
	if len(serialized) < paramsLen+dcPubKeyFieldLen {
		return 0, errors.New("credential is too short")
	}
	// First several bytes are the valid_time, scheme, and version fields.
	serialized = serialized[paramsLen:]

	// The next 3 bytes are the length of the serialized public key, which may
	// be up to 2^24 bytes in length.
	serializedPublicKeyLen := getUint24(serialized)
	serialized = serialized[dcPubKeyFieldLen:]

	if len(serialized) < serializedPublicKeyLen {
		return 0, errors.New("public key of credential is too short")
	}

	return paramsLen + dcPubKeyFieldLen + serializedPublicKeyLen, nil
}

// delegatedCredential stores a credential and its delegation.
type delegatedCredential struct {
	raw []byte

	// The credential, which contains a public and its validity time.
	cred *credential

	// The signature scheme used to sign the credential.
	algorithm SignatureScheme

	// The credential's delegation.
	signature []byte
}

// ensureCertificateHasLeaf parses the leaf certificate if needed.
func ensureCertificateHasLeaf(cert *Certificate) error {
	var err error
	if cert.Leaf == nil {
		if len(cert.Certificate[0]) == 0 {
			return errors.New("missing leaf certificate")
		}
		cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
		if err != nil {
			return err
		}
	}
	return nil
}

// validate checks that that the signature is valid, that the credential hasn't
// expired, and that the TTL is valid. It also checks that certificate can be
// used for delegation.
func (dc *delegatedCredential) validate(cert *x509.Certificate, now time.Time) (bool, error) {
	// Check that the cert can delegate.
	if !canDelegate(cert) {
		return false, errNoDelegationUsage
	}

	if dc.cred.isExpired(cert.NotBefore, now) {
		return false, errors.New("credential has expired")
	}

	if dc.cred.invalidTTL(cert.NotBefore, now) {
		return false, errors.New("credential TTL is invalid")
	}

	// Prepare the credential for verification.
	rawCred, err := dc.cred.marshal()
	if err != nil {
		return false, err
	}
	hash := getHash(dc.algorithm)
	in := prepareDelegation(hash, rawCred, cert.Raw, dc.algorithm)

	// TODO(any) This code overlaps significantly with verifyHandshakeSignature()
	// in ../auth.go. This should be refactored.
	switch dc.algorithm {
	case ECDSAWithP256AndSHA256,
		ECDSAWithP384AndSHA384,
		ECDSAWithP521AndSHA512:
		pk, ok := cert.PublicKey.(*ecdsa.PublicKey)
		if !ok {
			return false, errors.New("expected ECDSA public key")
		}
		sig := new(ecdsaSignature)
		if _, err = asn1.Unmarshal(dc.signature, sig); err != nil {
			return false, err
		}
		return ecdsa.Verify(pk, in, sig.R, sig.S), nil

	default:
		return false, fmt.Errorf(
			"unsupported signature scheme: 0x%04x", dc.algorithm)
	}
}

// unmarshalDelegatedCredential decodes a DelegatedCredential structure.
func unmarshalDelegatedCredential(serialized []byte) (*delegatedCredential, error) {
	// Get the length of the serialized credential that begins at the start of
	// the input slice.
	serializedCredentialLen, err := getCredentialLen(serialized)
	if err != nil {
		return nil, err
	}

	// Parse the credential.
	cred, err := unmarshalCredential(serialized[:serializedCredentialLen])
	if err != nil {
		return nil, err
	}

	// Parse the signature scheme.
	serialized = serialized[serializedCredentialLen:]
	if len(serialized) < 4 {
		return nil, errors.New("delegated credential is too short")
	}
	scheme := SignatureScheme(binary.BigEndian.Uint16(serialized))

	// Parse the signature length.
	serialized = serialized[2:]
	serializedSignatureLen := binary.BigEndian.Uint16(serialized)

	// Prase the signature.
	serialized = serialized[2:]
	if len(serialized) < int(serializedSignatureLen) {
		return nil, errors.New("signature of delegated credential is too short")
	}
	sig := serialized[:serializedSignatureLen]

	return &delegatedCredential{
		raw:       serialized,
		cred:      cred,
		algorithm: scheme,
		signature: sig,
	}, nil
}

// getCurve maps the SignatureScheme to its corresponding elliptic.Curve.
func getCurve(scheme SignatureScheme) elliptic.Curve {
	switch scheme {
	case ECDSAWithP256AndSHA256:
		return elliptic.P256()
	case ECDSAWithP384AndSHA384:
		return elliptic.P384()
	case ECDSAWithP521AndSHA512:
		return elliptic.P521()
	default:
		return nil
	}
}

// getHash maps the SignatureScheme to its corresponding hash function.
//
// TODO(any) This function overlaps with hashForSignatureScheme in 13.go.
func getHash(scheme SignatureScheme) crypto.Hash {
	switch scheme {
	case ECDSAWithP256AndSHA256:
		return crypto.SHA256
	case ECDSAWithP384AndSHA384:
		return crypto.SHA384
	case ECDSAWithP521AndSHA512:
		return crypto.SHA512
	default:
		return 0 // Unknown hash function
	}
}

// prepareDelegation returns a hash of the message that the delegator is to
// sign. The inputs are the credential (`cred`), the DER-encoded delegator
// certificate (`delegatorCert`) and the signature scheme of the delegator
// (`delegatorAlgorithm`).
func prepareDelegation(hash crypto.Hash, cred, delegatorCert []byte, delegatorAlgorithm SignatureScheme) []byte {
	h := hash.New()

	// The header.
	h.Write(bytes.Repeat([]byte{0x20}, 64))
	h.Write([]byte("TLS, server delegated credentials"))
	h.Write([]byte{0x00})

	// The delegation certificate.
	h.Write(delegatorCert)

	// The credential.
	h.Write(cred)

	// The delegator signature scheme.
	var serializedScheme [2]byte
	binary.BigEndian.PutUint16(serializedScheme[:], uint16(delegatorAlgorithm))
	h.Write(serializedScheme[:])

	return h.Sum(nil)
}
