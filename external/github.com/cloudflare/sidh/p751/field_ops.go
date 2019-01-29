package p751

import . "v2ray.com/core/external/github.com/cloudflare/sidh/internal/isogeny"

// 2*p751
var ()

//------------------------------------------------------------------------------
// Implementtaion of FieldOperations
//------------------------------------------------------------------------------

// Implements FieldOps
type fp751Ops struct{}

func FieldOperations() FieldOps {
	return &fp751Ops{}
}

func (fp751Ops) Add(dest, lhs, rhs *Fp2Element) {
	fp751AddReduced(&dest.A, &lhs.A, &rhs.A)
	fp751AddReduced(&dest.B, &lhs.B, &rhs.B)
}

func (fp751Ops) Sub(dest, lhs, rhs *Fp2Element) {
	fp751SubReduced(&dest.A, &lhs.A, &rhs.A)
	fp751SubReduced(&dest.B, &lhs.B, &rhs.B)
}

func (fp751Ops) Mul(dest, lhs, rhs *Fp2Element) {
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
	fp751Mul(&ac, a, c) // = a*c*R*R
	fp751Mul(&bd, b, d) // = b*d*R*R

	var b_minus_a, c_minus_d FpElement
	fp751SubReduced(&b_minus_a, b, a) // = (b-a)*R
	fp751SubReduced(&c_minus_d, c, d) // = (c-d)*R

	var ad_plus_bc FpElementX2
	fp751Mul(&ad_plus_bc, &b_minus_a, &c_minus_d) // = (b-a)*(c-d)*R*R
	fp751X2AddLazy(&ad_plus_bc, &ad_plus_bc, &ac) // = ((b-a)*(c-d) + a*c)*R*R
	fp751X2AddLazy(&ad_plus_bc, &ad_plus_bc, &bd) // = ((b-a)*(c-d) + a*c + b*d)*R*R

	fp751MontgomeryReduce(&dest.B, &ad_plus_bc) // = (a*d + b*c)*R mod p

	var ac_minus_bd FpElementX2
	fp751X2SubLazy(&ac_minus_bd, &ac, &bd)       // = (a*c - b*d)*R*R
	fp751MontgomeryReduce(&dest.A, &ac_minus_bd) // = (a*c - b*d)*R mod p
}

func (fp751Ops) Square(dest, x *Fp2Element) {
	a := &x.A
	b := &x.B

	// We want to compute
	//
	// (a + bi)*(a + bi) = (a^2 - b^2) + 2abi.

	var a2, a_plus_b, a_minus_b FpElement
	fp751AddReduced(&a2, a, a)        // = a*R + a*R = 2*a*R
	fp751AddReduced(&a_plus_b, a, b)  // = a*R + b*R = (a+b)*R
	fp751SubReduced(&a_minus_b, a, b) // = a*R - b*R = (a-b)*R

	var asq_minus_bsq, ab2 FpElementX2
	fp751Mul(&asq_minus_bsq, &a_plus_b, &a_minus_b) // = (a+b)*(a-b)*R*R = (a^2 - b^2)*R*R
	fp751Mul(&ab2, &a2, b)                          // = 2*a*b*R*R

	fp751MontgomeryReduce(&dest.A, &asq_minus_bsq) // = (a^2 - b^2)*R mod p
	fp751MontgomeryReduce(&dest.B, &ab2)           // = 2*a*b*R mod p
}

// Set dest = 1/x
//
// Allowed to overlap dest with x.
//
// Returns dest to allow chaining operations.
func (fp751Ops) Inv(dest, x *Fp2Element) {
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
	fp751Mul(&asq, a, a)                         // = a*a*R*R
	fp751Mul(&bsq, b, b)                         // = b*b*R*R
	fp751X2AddLazy(&asq, &asq, &bsq)             // = (a^2 + b^2)*R*R
	fp751MontgomeryReduce(&asq_plus_bsq.A, &asq) // = (a^2 + b^2)*R mod p
	// Now asq_plus_bsq = a^2 + b^2

	// Invert asq_plus_bsq
	inv := asq_plus_bsq
	inv.Mul(&asq_plus_bsq, &asq_plus_bsq)
	inv.P34(&inv)
	inv.Mul(&inv, &inv)
	inv.Mul(&inv, &asq_plus_bsq)

	var ac FpElementX2
	fp751Mul(&ac, a, &inv.A)
	fp751MontgomeryReduce(&dest.A, &ac)

	var minus_b FpElement
	fp751SubReduced(&minus_b, &minus_b, b)
	var minus_bc FpElementX2
	fp751Mul(&minus_bc, &minus_b, &inv.A)
	fp751MontgomeryReduce(&dest.B, &minus_bc)
}

// In case choice == 1, performs following swap in constant time:
// 	xPx <-> xQx
//	xPz <-> xQz
// Otherwise returns xPx, xPz, xQx, xQz unchanged
func (fp751Ops) CondSwap(xPx, xPz, xQx, xQz *Fp2Element, choice uint8) {
	fp751ConditionalSwap(&xPx.A, &xQx.A, choice)
	fp751ConditionalSwap(&xPx.B, &xQx.B, choice)
	fp751ConditionalSwap(&xPz.A, &xQz.A, choice)
	fp751ConditionalSwap(&xPz.B, &xQz.B, choice)
}

// Converts values in x.A and x.B to Montgomery domain
// x.A = x.A * R mod p
// x.B = x.B * R mod p
func (fp751Ops) ToMontgomery(x *Fp2Element) {
	var aRR FpElementX2

	// convert to montgomery domain
	fp751Mul(&aRR, &x.A, &p751R2)     // = a*R*R
	fp751MontgomeryReduce(&x.A, &aRR) // = a*R mod p
	fp751Mul(&aRR, &x.B, &p751R2)
	fp751MontgomeryReduce(&x.B, &aRR)
}

// Converts values in x.A and x.B from Montgomery domain
// a = x.A mod p
// b = x.B mod p
//
// After returning from the call x is not modified.
func (fp751Ops) FromMontgomery(x *Fp2Element, out *Fp2Element) {
	var aR FpElementX2

	// convert from montgomery domain
	copy(aR[:], x.A[:])
	fp751MontgomeryReduce(&out.A, &aR) // = a mod p in [0, 2p)
	fp751StrongReduce(&out.A)          // = a mod p in [0, p)
	for i := range aR {
		aR[i] = 0
	}
	copy(aR[:], x.B[:])
	fp751MontgomeryReduce(&out.B, &aR)
	fp751StrongReduce(&out.B)
}

//------------------------------------------------------------------------------
// Prime Field
//------------------------------------------------------------------------------

// Represents an element of the prime field F_p in Montgomery domain
type primeFieldElement struct {
	// The value `A`is represented by `aR mod p`.
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
	fp751Mul(&ab, a, b)                 // = a*b*R*R
	fp751MontgomeryReduce(&dest.A, &ab) // = a*b*R mod p

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
//
// Allowed to overlap x with dest.
//
// Returns dest to allow chaining operations.
func (dest *primeFieldElement) P34(x *primeFieldElement) *primeFieldElement {
	// Sliding-window strategy computed with Sage, awk, sed, and tr.
	//
	// This performs sum(powStrategy) = 744 squarings and len(mulStrategy)
	// = 137 multiplications, in addition to 1 squaring and 15
	// multiplications to build a lookup table.
	//
	// In total this is 745 squarings, 152 multiplications.  Since squaring
	// is not implemented for the prime field, this is 897 multiplications
	// in total.
	powStrategy := [137]uint8{5, 7, 6, 2, 10, 4, 6, 9, 8, 5, 9, 4, 7, 5, 5, 4, 8, 3, 9, 5, 5, 4, 10, 4, 6, 6, 6, 5, 8, 9, 3, 4, 9, 4, 5, 6, 6, 2, 9, 4, 5, 5, 5, 7, 7, 9, 4, 6, 4, 8, 5, 8, 6, 6, 2, 9, 7, 4, 8, 8, 8, 4, 6, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 2}
	mulStrategy := [137]uint8{31, 23, 21, 1, 31, 7, 7, 7, 9, 9, 19, 15, 23, 23, 11, 7, 25, 5, 21, 17, 11, 5, 17, 7, 11, 9, 23, 9, 1, 19, 5, 3, 25, 15, 11, 29, 31, 1, 29, 11, 13, 9, 11, 27, 13, 19, 15, 31, 3, 29, 23, 31, 25, 11, 1, 21, 19, 15, 15, 21, 29, 13, 23, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 3}
	initialMul := uint8(27)

	// Build a lookup table of odd multiples of x.
	lookup := [16]primeFieldElement{}
	xx := &primeFieldElement{}
	xx.Mul(x, x) // Set xx = x^2
	lookup[0] = *x
	for i := 1; i < 16; i++ {
		lookup[i].Mul(&lookup[i-1], xx)
	}
	// Now lookup = {x, x^3, x^5, ... }
	// so that lookup[i] = x^{2*i + 1}
	// so that lookup[k/2] = x^k, for odd k

	*dest = lookup[initialMul/2]
	for i := uint8(0); i < 137; i++ {
		dest.Pow2k(dest, powStrategy[i])
		dest.Mul(dest, &lookup[mulStrategy[i]/2])
	}

	return dest
}
