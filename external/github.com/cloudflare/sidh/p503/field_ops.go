package p503

import (
	. "v2ray.com/core/external/github.com/cloudflare/sidh/internal/isogeny"
)

type fp503Ops struct{}

func FieldOperations() FieldOps {
	return &fp503Ops{}
}

func (fp503Ops) Add(dest, lhs, rhs *Fp2Element) {
	fp503AddReduced(&dest.A, &lhs.A, &rhs.A)
	fp503AddReduced(&dest.B, &lhs.B, &rhs.B)
}

func (fp503Ops) Sub(dest, lhs, rhs *Fp2Element) {
	fp503SubReduced(&dest.A, &lhs.A, &rhs.A)
	fp503SubReduced(&dest.B, &lhs.B, &rhs.B)
}

func (fp503Ops) Mul(dest, lhs, rhs *Fp2Element) {
	// Let (a,b,c,d) = (lhs.a,lhs.b,rhs.a,rhs.b).
	a := &lhs.A
	b := &lhs.B
	c := &rhs.A
	d := &rhs.B

	// We want to compute
	//
	// (a + bi)*(c + di) = (a*c - b*d) + (a*d + b*c)i
	//
	// Use Karatsuba's trick: note that
	//
	// (b - a)*(c - d) = (b*c + a*d) - a*c - b*d
	//
	// so (a*d + b*c) = (b-a)*(c-d) + a*c + b*d.

	var ac, bd FpElementX2
	fp503Mul(&ac, a, c) // = a*c*R*R
	fp503Mul(&bd, b, d) // = b*d*R*R

	var b_minus_a, c_minus_d FpElement
	fp503SubReduced(&b_minus_a, b, a) // = (b-a)*R
	fp503SubReduced(&c_minus_d, c, d) // = (c-d)*R

	var ad_plus_bc FpElementX2
	fp503Mul(&ad_plus_bc, &b_minus_a, &c_minus_d) // = (b-a)*(c-d)*R*R
	fp503X2AddLazy(&ad_plus_bc, &ad_plus_bc, &ac) // = ((b-a)*(c-d) + a*c)*R*R
	fp503X2AddLazy(&ad_plus_bc, &ad_plus_bc, &bd) // = ((b-a)*(c-d) + a*c + b*d)*R*R

	fp503MontgomeryReduce(&dest.B, &ad_plus_bc) // = (a*d + b*c)*R mod p

	var ac_minus_bd FpElementX2
	fp503X2SubLazy(&ac_minus_bd, &ac, &bd)       // = (a*c - b*d)*R*R
	fp503MontgomeryReduce(&dest.A, &ac_minus_bd) // = (a*c - b*d)*R mod p
}

// Set dest = 1/x
//
// Allowed to overlap dest with x.
//
// Returns dest to allow chaining operations.
func (fp503Ops) Inv(dest, x *Fp2Element) {
	a := &x.A
	b := &x.B

	// We want to compute
	//
	//    1          1     (a - bi)	    (a - bi)
	// -------- = -------- -------- = -----------
	// (a + bi)   (a + bi) (a - bi)   (a^2 + b^2)
	//
	// Letting c = 1/(a^2 + b^2), this is
	//
	// 1/(a+bi) = a*c - b*ci.

	var asq_plus_bsq primeFieldElement
	var asq, bsq FpElementX2
	fp503Mul(&asq, a, a)                         // = a*a*R*R
	fp503Mul(&bsq, b, b)                         // = b*b*R*R
	fp503X2AddLazy(&asq, &asq, &bsq)             // = (a^2 + b^2)*R*R
	fp503MontgomeryReduce(&asq_plus_bsq.A, &asq) // = (a^2 + b^2)*R mod p
	// Now asq_plus_bsq = a^2 + b^2

	inv := asq_plus_bsq
	inv.Mul(&asq_plus_bsq, &asq_plus_bsq)
	inv.P34(&inv)
	inv.Mul(&inv, &inv)
	inv.Mul(&inv, &asq_plus_bsq)

	var ac FpElementX2
	fp503Mul(&ac, a, &inv.A)
	fp503MontgomeryReduce(&dest.A, &ac)

	var minus_b FpElement
	fp503SubReduced(&minus_b, &minus_b, b)
	var minus_bc FpElementX2
	fp503Mul(&minus_bc, &minus_b, &inv.A)
	fp503MontgomeryReduce(&dest.B, &minus_bc)
}

func (fp503Ops) Square(dest, x *Fp2Element) {
	a := &x.A
	b := &x.B

	// We want to compute
	//
	// (a + bi)*(a + bi) = (a^2 - b^2) + 2abi.

	var a2, a_plus_b, a_minus_b FpElement
	fp503AddReduced(&a2, a, a)        // = a*R + a*R = 2*a*R
	fp503AddReduced(&a_plus_b, a, b)  // = a*R + b*R = (a+b)*R
	fp503SubReduced(&a_minus_b, a, b) // = a*R - b*R = (a-b)*R

	var asq_minus_bsq, ab2 FpElementX2
	fp503Mul(&asq_minus_bsq, &a_plus_b, &a_minus_b) // = (a+b)*(a-b)*R*R = (a^2 - b^2)*R*R
	fp503Mul(&ab2, &a2, b)                          // = 2*a*b*R*R

	fp503MontgomeryReduce(&dest.A, &asq_minus_bsq) // = (a^2 - b^2)*R mod p
	fp503MontgomeryReduce(&dest.B, &ab2)           // = 2*a*b*R mod p
}

// In case choice == 1, performs following swap in constant time:
// 	xPx <-> xQx
//	xPz <-> xQz
// Otherwise returns xPx, xPz, xQx, xQz unchanged
func (fp503Ops) CondSwap(xPx, xPz, xQx, xQz *Fp2Element, choice uint8) {
	fp503ConditionalSwap(&xPx.A, &xQx.A, choice)
	fp503ConditionalSwap(&xPx.B, &xQx.B, choice)
	fp503ConditionalSwap(&xPz.A, &xQz.A, choice)
	fp503ConditionalSwap(&xPz.B, &xQz.B, choice)
}

// Converts values in x.A and x.B to Montgomery domain
// x.A = x.A * R mod p
// x.B = x.B * R mod p
// Performs v = v*R^2*R^(-1) mod p, for both x.A and x.B
func (fp503Ops) ToMontgomery(x *Fp2Element) {
	var aRR FpElementX2

	// convert to montgomery domain
	fp503Mul(&aRR, &x.A, &p503R2)     // = a*R*R
	fp503MontgomeryReduce(&x.A, &aRR) // = a*R mod p
	fp503Mul(&aRR, &x.B, &p503R2)
	fp503MontgomeryReduce(&x.B, &aRR)
}

// Converts values in x.A and x.B from Montgomery domain
// a = x.A mod p
// b = x.B mod p
//
// After returning from the call x is not modified.
func (fp503Ops) FromMontgomery(x *Fp2Element, out *Fp2Element) {
	var aR FpElementX2

	// convert from montgomery domain
	// TODO: make fpXXXMontgomeryReduce use stack instead of reusing aR
	//       so that we don't have do this copy here
	copy(aR[:], x.A[:])
	fp503MontgomeryReduce(&out.A, &aR) // = a mod p in [0, 2p)
	fp503StrongReduce(&out.A)          // = a mod p in [0, p)
	for i := range aR {
		aR[i] = 0
	}
	copy(aR[:], x.B[:])
	fp503MontgomeryReduce(&out.B, &aR)
	fp503StrongReduce(&out.B)
}

//------------------------------------------------------------------------------
// Prime Field
//------------------------------------------------------------------------------

// Represents an element of the prime field F_p.
type primeFieldElement struct {
	// This field element is in Montgomery form, so that the value `A` is
	// represented by `aR mod p`.
	A FpElement
}

// Set dest = lhs * rhs.
//
// Allowed to overlap lhs or rhs with dest.
//
// Returns dest to allow chaining operations.
func (dest *primeFieldElement) Mul(lhs, rhs *primeFieldElement) *primeFieldElement {
	a := &lhs.A // = a*R
	b := &rhs.A // = b*R

	var ab FpElementX2
	fp503Mul(&ab, a, b)                 // = a*b*R*R
	fp503MontgomeryReduce(&dest.A, &ab) // = a*b*R mod p

	return dest
}

// Set dest = x^(2^k), for k >= 1, by repeated squarings.
//
// Allowed to overlap x with dest.
//
// Returns dest to allow chaining operations.
func (dest *primeFieldElement) Pow2k(x *primeFieldElement, k uint8) *primeFieldElement {
	dest.Mul(x, x)
	for i := uint8(1); i < k; i++ {
		dest.Mul(dest, dest)
	}

	return dest
}

// Set dest = x^((p-3)/4).  If x is square, this is 1/sqrt(x).
// Uses variation of sliding-window algorithm from with window size
// of 5 and least to most significant bit sliding (left-to-right)
// See HAC 14.85 for general description.
//
// Allowed to overlap x with dest.
//
// Returns dest to allow chaining operations.
func (dest *primeFieldElement) P34(x *primeFieldElement) *primeFieldElement {
	// Sliding-window strategy computed with etc/scripts/sliding_window_strat_calc.py
	//
	// This performs sum(powStrategy) + 1 squarings and len(lookup) + len(mulStrategy)
	// multiplications.
	powStrategy := []uint8{1, 12, 5, 5, 2, 7, 11, 3, 8, 4, 11, 4, 7, 5, 6, 3, 7, 5, 7, 2, 12, 5, 6, 4, 6, 8, 6, 4, 7, 5, 5, 8, 5, 8, 5, 5, 8, 9, 3, 6, 2, 10, 6, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 3}
	mulStrategy := []uint8{0, 12, 11, 10, 0, 1, 8, 3, 7, 1, 8, 3, 6, 7, 14, 2, 14, 14, 9, 0, 13, 9, 15, 5, 12, 7, 13, 7, 15, 6, 7, 9, 0, 5, 7, 6, 8, 8, 3, 7, 0, 10, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 3}

	// Precompute lookup table of odd multiples of x for window
	// size k=5.
	lookup := [16]primeFieldElement{}
	xx := &primeFieldElement{}
	xx.Mul(x, x)
	lookup[0] = *x
	for i := 1; i < 16; i++ {
		lookup[i].Mul(&lookup[i-1], xx)
	}

	// Now lookup = {x, x^3, x^5, ... }
	// so that lookup[i] = x^{2*i + 1}
	// so that lookup[k/2] = x^k, for odd k
	*dest = lookup[mulStrategy[0]]
	for i := uint8(1); i < uint8(len(powStrategy)); i++ {
		dest.Pow2k(dest, powStrategy[i])
		dest.Mul(dest, &lookup[mulStrategy[i]])
	}

	return dest
}
