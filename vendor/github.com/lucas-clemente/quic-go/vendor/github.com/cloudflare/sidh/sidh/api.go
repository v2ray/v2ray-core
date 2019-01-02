package sidh

import (
	"errors"
	. "github.com/cloudflare/sidh/internal/isogeny"
	"io"
)

// I keep it bool in order to be able to apply logical NOT
type KeyVariant uint

// Id's correspond to bitlength of the prime field characteristic
// Currently FP_751 is the only one supported by this implementation
const (
	FP_503 uint8 = iota
	FP_751
	FP_964
	maxPrimeFieldId
)

const (
	// First 2 bits identify SIDH variant third bit indicates
	// wether key is a SIKE variant (set) or SIDH (not set)

	// 001 - SIDH: corresponds to 2-torsion group
	KeyVariant_SIDH_A KeyVariant = 1 << 0
	// 010 - SIDH: corresponds to 3-torsion group
	KeyVariant_SIDH_B = 1 << 1
	// 110 - SIKE
	KeyVariant_SIKE = 1<<2 | KeyVariant_SIDH_B
)

// Base type for public and private key. Used mainly to carry domain
// parameters.
type key struct {
	// Domain parameters of the algorithm to be used with a key
	params *SidhParams
	// Flag indicates wether corresponds to 2-, 3-torsion group or SIKE
	keyVariant KeyVariant
}

// Defines operations on public key
type PublicKey struct {
	key
	affine_xP   Fp2Element
	affine_xQ   Fp2Element
	affine_xQmP Fp2Element
}

// Defines operations on private key
type PrivateKey struct {
	key
	// Secret key
	Scalar []byte
	// Used only by KEM
	S []byte
}

// Accessor to the domain parameters
func (key *key) Params() *SidhParams {
	return key.params
}

// Accessor to key variant
func (key *key) Variant() KeyVariant {
	return key.keyVariant
}

// NewPrivateKey initializes private key.
// Usage of this function guarantees that the object is correctly initialized.
func NewPrivateKey(id uint8, v KeyVariant) *PrivateKey {
	prv := &PrivateKey{key: key{params: Params(id), keyVariant: v}}
	if (v & KeyVariant_SIDH_A) == KeyVariant_SIDH_A {
		prv.Scalar = make([]byte, prv.params.A.SecretByteLen)
	} else {
		prv.Scalar = make([]byte, prv.params.B.SecretByteLen)
	}
	if v == KeyVariant_SIKE {
		prv.S = make([]byte, prv.params.MsgLen)
	}
	return prv
}

// NewPublicKey initializes public key.
// Usage of this function guarantees that the object is correctly initialized.
func NewPublicKey(id uint8, v KeyVariant) *PublicKey {
	return &PublicKey{key: key{params: Params(id), keyVariant: v}}
}

// Import clears content of the public key currently stored in the structure
// and imports key stored in the byte string. Returns error in case byte string
// size is wrong. Doesn't perform any validation.
func (pub *PublicKey) Import(input []byte) error {
	if len(input) != pub.Size() {
		return errors.New("sidh: input to short")
	}
	op := CurveOperations{Params: pub.params}
	ssSz := pub.params.SharedSecretSize
	op.Fp2FromBytes(&pub.affine_xP, input[0:ssSz])
	op.Fp2FromBytes(&pub.affine_xQ, input[ssSz:2*ssSz])
	op.Fp2FromBytes(&pub.affine_xQmP, input[2*ssSz:3*ssSz])
	return nil
}

// Exports currently stored key. In case structure hasn't been filled with key data
// returned byte string is filled with zeros.
func (pub *PublicKey) Export() []byte {
	output := make([]byte, pub.params.PublicKeySize)
	op := CurveOperations{Params: pub.params}
	ssSz := pub.params.SharedSecretSize
	op.Fp2ToBytes(output[0:ssSz], &pub.affine_xP)
	op.Fp2ToBytes(output[ssSz:2*ssSz], &pub.affine_xQ)
	op.Fp2ToBytes(output[2*ssSz:3*ssSz], &pub.affine_xQmP)
	return output
}

// Size returns size of the public key in bytes
func (pub *PublicKey) Size() int {
	return pub.params.PublicKeySize
}

// Exports currently stored key. In case structure hasn't been filled with key data
// returned byte string is filled with zeros.
func (prv *PrivateKey) Export() []byte {
	ret := make([]byte, len(prv.Scalar)+len(prv.S))
	copy(ret, prv.S)
	copy(ret[len(prv.S):], prv.Scalar)
	return ret
}

// Size returns size of the private key in bytes
func (prv *PrivateKey) Size() int {
	tmp := len(prv.Scalar)
	if prv.Variant() == KeyVariant_SIKE {
		tmp += int(prv.params.MsgLen)
	}
	return tmp
}

// Import clears content of the private key currently stored in the structure
// and imports key from octet string. In case of SIKE, the random value 'S'
// must be prepended to the value of actual private key (see SIKE spec for details).
// Function doesn't import public key value to PrivateKey object.
func (prv *PrivateKey) Import(input []byte) error {
	if len(input) != prv.Size() {
		return errors.New("sidh: input to short")
	}
	copy(prv.S, input[:len(prv.S)])
	copy(prv.Scalar, input[len(prv.S):])
	return nil
}

// Generates random private key for SIDH or SIKE. Generated value is
// formed as little-endian integer from key-space <2^(e2-1)..2^e2 - 1>
// for KeyVariant_A or <2^(s-1)..2^s - 1>, where s = floor(log_2(3^e3)),
// for KeyVariant_B.
//
// Returns error in case user provided RNG fails.
func (prv *PrivateKey) Generate(rand io.Reader) error {
	var err error
	var dp *DomainParams

	if (prv.keyVariant & KeyVariant_SIDH_A) == KeyVariant_SIDH_A {
		dp = &prv.params.A
	} else {
		dp = &prv.params.B
	}

	if prv.keyVariant == KeyVariant_SIKE && err == nil {
		_, err = io.ReadFull(rand, prv.S)
	}

	// Private key generation takes advantage of the fact that keyspace for secret
	// key is (0, 2^x - 1), for some possitivite value of 'x' (see SIKE, 1.3.8).
	// It means that all bytes in the secret key, but the last one, can take any
	// value between <0x00,0xFF>. Similarily for the last byte, but generation
	// needs to chop off some bits, to make sure generated value is an element of
	// a key-space.
	_, err = io.ReadFull(rand, prv.Scalar)
	if err != nil {
		return err
	}
	prv.Scalar[len(prv.Scalar)-1] &= (1 << (dp.SecretBitLen % 8)) - 1
	// Make sure scalar is SecretBitLen long. SIKE spec says that key
	// space starts from 0, but I'm not confortable with having low
	// value scalars used for private keys. It is still secrure as per
	// table 5.1 in [SIKE].
	prv.Scalar[len(prv.Scalar)-1] |= 1 << ((dp.SecretBitLen % 8) - 1)
	return err
}

// Generates public key.
//
// Constant time.
func (prv *PrivateKey) GeneratePublicKey() *PublicKey {
	if (prv.keyVariant & KeyVariant_SIDH_A) == KeyVariant_SIDH_A {
		return publicKeyGenA(prv)
	}
	return publicKeyGenB(prv)
}

// Computes a shared secret which is a j-invariant. Function requires that pub has
// different KeyVariant than prv. Length of returned output is 2*ceil(log_2 P)/8),
// where P is a prime defining finite field.
//
// It's important to notice that each keypair must not be used more than once
// to calculate shared secret.
//
// Function may return error. This happens only in case provided input is invalid.
// Constant time for properly initialized private and public key.
func DeriveSecret(prv *PrivateKey, pub *PublicKey) ([]byte, error) {

	if (pub == nil) || (prv == nil) {
		return nil, errors.New("sidh: invalid arguments")
	}

	if (pub.keyVariant == prv.keyVariant) || (pub.params.Id != prv.params.Id) {
		return nil, errors.New("sidh: public and private are incompatbile")
	}

	if (prv.keyVariant & KeyVariant_SIDH_A) == KeyVariant_SIDH_A {
		return deriveSecretA(prv, pub), nil
	} else {
		return deriveSecretB(prv, pub), nil
	}
}
