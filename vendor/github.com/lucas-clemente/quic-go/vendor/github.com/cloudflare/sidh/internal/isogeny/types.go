package internal

const (
	FP_MAX_WORDS = 12 // Currently p751.NumWords
)

// Representation of an element of the base field F_p.
//
// No particular meaning is assigned to the representation -- it could represent
// an element in Montgomery form, or not.  Tracking the meaning of the field
// element is left to higher types.
type FpElement [FP_MAX_WORDS]uint64

// Represents an intermediate product of two elements of the base field F_p.
type FpElementX2 [2 * FP_MAX_WORDS]uint64

// Represents an element of the extended field Fp^2 = Fp(x+i)
type Fp2Element struct {
	A FpElement
	B FpElement
}

type DomainParams struct {
	// P, Q and R=P-Q base points
	Affine_P, Affine_Q, Affine_R Fp2Element
	// Size of a compuatation strategy for x-torsion group
	IsogenyStrategy []uint32
	// Max size of secret key for x-torsion group
	SecretBitLen uint
	// Max size of secret key for x-torsion group
	SecretByteLen uint
}

type SidhParams struct {
	Id uint8
	// Bytelen of P
	Bytelen int
	// The public key size, in bytes.
	PublicKeySize int
	// The shared secret size, in bytes.
	SharedSecretSize uint
	// 2- and 3-torsion group parameter definitions
	A, B DomainParams
	// Precomputed identity element in the Fp2 in Montgomery domain
	OneFp2 Fp2Element
	// Precomputed 1/2 in the Fp2 in Montgomery domain
	HalfFp2 Fp2Element
	// Length of SIKE secret message. Must be one of {24,32,40},
	// depending on size of prime field used (see [SIKE], 1.4 and 5.1)
	MsgLen uint
	// Length of SIKE ephemeral KEM key (see [SIKE], 1.4 and 5.1)
	KemSize uint
	// Access to field arithmetic
	Op FieldOps
}

// Interface for working with isogenies.
type Isogeny interface {
	// Given a torsion point on a curve computes isogenous curve.
	// Returns curve coefficients (A:C), so that E_(A/C) = E_(A/C)/<P>,
	// where P is a provided projective point. Sets also isogeny constants
	// that are needed for isogeny evaluation.
	GenerateCurve(*ProjectivePoint) CurveCoefficientsEquiv
	// Evaluates isogeny at caller provided point. Requires isogeny curve constants
	// to be earlier computed by GenerateCurve.
	EvaluatePoint(*ProjectivePoint) ProjectivePoint
}

// Stores curve projective parameters equivalent to A/C. Meaning of the
// values depends on the context. When working with isogenies over
// subgroup that are powers of:
// * three then  (A:C) ~ (A+2C:A-2C)
// * four then   (A:C) ~ (A+2C:  4C)
// See Appendix A of SIKE for more details
type CurveCoefficientsEquiv struct {
	A Fp2Element
	C Fp2Element
}

// A point on the projective line P^1(F_{p^2}).
//
// This represents a point on the Kummer line of a Montgomery curve.  The
// curve is specified by a ProjectiveCurveParameters struct.
type ProjectivePoint struct {
	X Fp2Element
	Z Fp2Element
}

// A point on the projective line P^1(F_{p^2}).
//
// This is used to work projectively with the curve coefficients.
type ProjectiveCurveParameters struct {
	A Fp2Element
	C Fp2Element
}

// Stores Isogeny 3 curve constants
type isogeny3 struct {
	Field FieldOps
	K1    Fp2Element
	K2    Fp2Element
}

// Stores Isogeny 4 curve constants
type isogeny4 struct {
	isogeny3
	K3 Fp2Element
}

type FieldOps interface {
	// Set res = lhs + rhs.
	//
	// Allowed to overlap lhs or rhs with res.
	Add(res, lhs, rhs *Fp2Element)

	// Set res = lhs - rhs.
	//
	// Allowed to overlap lhs or rhs with res.
	Sub(res, lhs, rhs *Fp2Element)

	// Set res = lhs * rhs.
	//
	// Allowed to overlap lhs or rhs with res.
	Mul(res, lhs, rhs *Fp2Element)
	// Set res = x * x
	//
	// Allowed to overlap res with x.
	Square(res, x *Fp2Element)
	// Set res = 1/x
	//
	// Allowed to overlap res with x.
	Inv(res, x *Fp2Element)
	// If choice = 1u8, set (x,y) = (y,x). If choice = 0u8, set (x,y) = (x,y).
	CondSwap(xPx, xPz, xQx, xQz *Fp2Element, choice uint8)
	// Converts Fp2Element to Montgomery domain (x*R mod p)
	ToMontgomery(x *Fp2Element)
	// Converts 'a' in montgomery domain to element from Fp2Element
	// and stores it in 'x'
	FromMontgomery(x *Fp2Element, a *Fp2Element)
}
