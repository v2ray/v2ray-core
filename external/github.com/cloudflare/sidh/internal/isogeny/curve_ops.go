package internal

type CurveOperations struct {
	Params *SidhParams
}

// Computes j-invariant for a curve y2=x3+A/Cx+x with A,C in F_(p^2). Result
// is returned in jBytes buffer, encoded in little-endian format. Caller
// provided jBytes buffer has to be big enough to j-invariant value. In case
// of SIDH, buffer size must be at least size of shared secret.
// Implementation corresponds to Algorithm 9 from SIKE.
func (c *CurveOperations) Jinvariant(cparams *ProjectiveCurveParameters, jBytes []byte) {
	var j, t0, t1 Fp2Element

	op := c.Params.Op
	op.Square(&j, &cparams.A)  // j  = A^2
	op.Square(&t1, &cparams.C) // t1 = C^2
	op.Add(&t0, &t1, &t1)      // t0 = t1 + t1
	op.Sub(&t0, &j, &t0)       // t0 = j - t0
	op.Sub(&t0, &t0, &t1)      // t0 = t0 - t1
	op.Sub(&j, &t0, &t1)       // t0 = t0 - t1
	op.Square(&t1, &t1)        // t1 = t1^2
	op.Mul(&j, &j, &t1)        // j = j * t1
	op.Add(&t0, &t0, &t0)      // t0 = t0 + t0
	op.Add(&t0, &t0, &t0)      // t0 = t0 + t0
	op.Square(&t1, &t0)        // t1 = t0^2
	op.Mul(&t0, &t0, &t1)      // t0 = t0 * t1
	op.Add(&t0, &t0, &t0)      // t0 = t0 + t0
	op.Add(&t0, &t0, &t0)      // t0 = t0 + t0
	op.Inv(&j, &j)             // j  = 1/j
	op.Mul(&j, &t0, &j)        // j  = t0 * j

	c.Fp2ToBytes(jBytes, &j)
}

// Given affine points x(P), x(Q) and x(Q-P) in a extension field F_{p^2}, function
// recorvers projective coordinate A of a curve. This is Algorithm 10 from SIKE.
func (c *CurveOperations) RecoverCoordinateA(curve *ProjectiveCurveParameters, xp, xq, xr *Fp2Element) {
	var t0, t1 Fp2Element

	op := c.Params.Op
	op.Add(&t1, xp, xq)                          // t1 = Xp + Xq
	op.Mul(&t0, xp, xq)                          // t0 = Xp * Xq
	op.Mul(&curve.A, xr, &t1)                    // A  = X(q-p) * t1
	op.Add(&curve.A, &curve.A, &t0)              // A  = A + t0
	op.Mul(&t0, &t0, xr)                         // t0 = t0 * X(q-p)
	op.Sub(&curve.A, &curve.A, &c.Params.OneFp2) // A  = A - 1
	op.Add(&t0, &t0, &t0)                        // t0 = t0 + t0
	op.Add(&t1, &t1, xr)                         // t1 = t1 + X(q-p)
	op.Add(&t0, &t0, &t0)                        // t0 = t0 + t0
	op.Square(&curve.A, &curve.A)                // A  = A^2
	op.Inv(&t0, &t0)                             // t0 = 1/t0
	op.Mul(&curve.A, &curve.A, &t0)              // A  = A * t0
	op.Sub(&curve.A, &curve.A, &t1)              // A  = A - t1
}

// Computes equivalence (A:C) ~ (A+2C : A-2C)
func (c *CurveOperations) CalcCurveParamsEquiv3(cparams *ProjectiveCurveParameters) CurveCoefficientsEquiv {
	var coef CurveCoefficientsEquiv
	var c2 Fp2Element
	var op = c.Params.Op

	op.Add(&c2, &cparams.C, &cparams.C)
	// A24p = A+2*C
	op.Add(&coef.A, &cparams.A, &c2)
	// A24m = A-2*C
	op.Sub(&coef.C, &cparams.A, &c2)
	return coef
}

// Computes equivalence (A:C) ~ (A+2C : 4C)
func (c *CurveOperations) CalcCurveParamsEquiv4(cparams *ProjectiveCurveParameters) CurveCoefficientsEquiv {
	var coefEq CurveCoefficientsEquiv
	var op = c.Params.Op

	op.Add(&coefEq.C, &cparams.C, &cparams.C)
	// A24p = A+2C
	op.Add(&coefEq.A, &cparams.A, &coefEq.C)
	// C24 = 4*C
	op.Add(&coefEq.C, &coefEq.C, &coefEq.C)
	return coefEq
}

// Helper function for RightToLeftLadder(). Returns A+2C / 4.
func (c *CurveOperations) CalcAplus2Over4(cparams *ProjectiveCurveParameters) (ret Fp2Element) {
	var tmp Fp2Element
	var op = c.Params.Op

	// 2C
	op.Add(&tmp, &cparams.C, &cparams.C)
	// A+2C
	op.Add(&ret, &cparams.A, &tmp)
	// 1/4C
	op.Add(&tmp, &tmp, &tmp)
	op.Inv(&tmp, &tmp)
	// A+2C/4C
	op.Mul(&ret, &ret, &tmp)
	return
}

// Recovers (A:C) curve parameters from projectively equivalent (A+2C:A-2C).
func (c *CurveOperations) RecoverCurveCoefficients3(cparams *ProjectiveCurveParameters, coefEq *CurveCoefficientsEquiv) {
	var op = c.Params.Op

	op.Add(&cparams.A, &coefEq.A, &coefEq.C)
	// cparams.A = 2*(A+2C+A-2C) = 4A
	op.Add(&cparams.A, &cparams.A, &cparams.A)
	// cparams.C = (A+2C-A+2C) = 4C
	op.Sub(&cparams.C, &coefEq.A, &coefEq.C)
	return
}

// Recovers (A:C) curve parameters from projectively equivalent (A+2C:4C).
func (c *CurveOperations) RecoverCurveCoefficients4(cparams *ProjectiveCurveParameters, coefEq *CurveCoefficientsEquiv) {
	var op = c.Params.Op
	// cparams.C = (4C)*1/2=2C
	op.Mul(&cparams.C, &coefEq.C, &c.Params.HalfFp2)
	// cparams.A = A+2C - 2C = A
	op.Sub(&cparams.A, &coefEq.A, &cparams.C)
	// cparams.C = 2C * 1/2 = C
	op.Mul(&cparams.C, &cparams.C, &c.Params.HalfFp2)
	return
}

// Combined coordinate doubling and differential addition. Takes projective points
// P,Q,Q-P and (A+2C)/4C curve E coefficient. Returns 2*P and P+Q calculated on E.
// Function is used only by RightToLeftLadder. Corresponds to Algorithm 5 of SIKE
func (c *CurveOperations) xDblAdd(P, Q, QmP *ProjectivePoint, a24 *Fp2Element) (dblP, PaQ ProjectivePoint) {
	var t0, t1, t2 Fp2Element
	var op = c.Params.Op

	xQmP, zQmP := &QmP.X, &QmP.Z
	xPaQ, zPaQ := &PaQ.X, &PaQ.Z
	x2P, z2P := &dblP.X, &dblP.Z
	xP, zP := &P.X, &P.Z
	xQ, zQ := &Q.X, &Q.Z

	op.Add(&t0, xP, zP)      // t0   = Xp+Zp
	op.Sub(&t1, xP, zP)      // t1   = Xp-Zp
	op.Square(x2P, &t0)      // 2P.X = t0^2
	op.Sub(&t2, xQ, zQ)      // t2   = Xq-Zq
	op.Add(xPaQ, xQ, zQ)     // Xp+q = Xq+Zq
	op.Mul(&t0, &t0, &t2)    // t0   = t0 * t2
	op.Mul(z2P, &t1, &t1)    // 2P.Z = t1 * t1
	op.Mul(&t1, &t1, xPaQ)   // t1   = t1 * Xp+q
	op.Sub(&t2, x2P, z2P)    // t2   = 2P.X - 2P.Z
	op.Mul(x2P, x2P, z2P)    // 2P.X = 2P.X * 2P.Z
	op.Mul(xPaQ, a24, &t2)   // Xp+q = A24 * t2
	op.Sub(zPaQ, &t0, &t1)   // Zp+q = t0 - t1
	op.Add(z2P, xPaQ, z2P)   // 2P.Z = Xp+q + 2P.Z
	op.Add(xPaQ, &t0, &t1)   // Xp+q = t0 + t1
	op.Mul(z2P, z2P, &t2)    // 2P.Z = 2P.Z * t2
	op.Square(zPaQ, zPaQ)    // Zp+q = Zp+q ^ 2
	op.Square(xPaQ, xPaQ)    // Xp+q = Xp+q ^ 2
	op.Mul(zPaQ, xQmP, zPaQ) // Zp+q = Xq-p * Zp+q
	op.Mul(xPaQ, zQmP, xPaQ) // Xp+q = Zq-p * Xp+q
	return
}

// Given the curve parameters, xP = x(P), computes xP = x([2^k]P)
// Safe to overlap xP, x2P.
func (c *CurveOperations) Pow2k(xP *ProjectivePoint, params *CurveCoefficientsEquiv, k uint32) {
	var t0, t1 Fp2Element
	var op = c.Params.Op

	x, z := &xP.X, &xP.Z
	for i := uint32(0); i < k; i++ {
		op.Sub(&t0, x, z)           // t0  = Xp - Zp
		op.Add(&t1, x, z)           // t1  = Xp + Zp
		op.Square(&t0, &t0)         // t0  = t0 ^ 2
		op.Square(&t1, &t1)         // t1  = t1 ^ 2
		op.Mul(z, &params.C, &t0)   // Z2p = C24 * t0
		op.Mul(x, z, &t1)           // X2p = Z2p * t1
		op.Sub(&t1, &t1, &t0)       // t1  = t1 - t0
		op.Mul(&t0, &params.A, &t1) // t0  = A24+ * t1
		op.Add(z, z, &t0)           // Z2p = Z2p + t0
		op.Mul(z, z, &t1)           // Zp  = Z2p * t1
	}
}

// Given the curve parameters, xP = x(P), and k >= 0, compute xP = x([3^k]P).
//
// Safe to overlap xP, xR.
func (c *CurveOperations) Pow3k(xP *ProjectivePoint, params *CurveCoefficientsEquiv, k uint32) {
	var t0, t1, t2, t3, t4, t5, t6 Fp2Element
	var op = c.Params.Op

	x, z := &xP.X, &xP.Z
	for i := uint32(0); i < k; i++ {
		op.Sub(&t0, x, z)           // t0  = Xp - Zp
		op.Square(&t2, &t0)         // t2  = t0^2
		op.Add(&t1, x, z)           // t1  = Xp + Zp
		op.Square(&t3, &t1)         // t3  = t1^2
		op.Add(&t4, &t1, &t0)       // t4  = t1 + t0
		op.Sub(&t0, &t1, &t0)       // t0  = t1 - t0
		op.Square(&t1, &t4)         // t1  = t4^2
		op.Sub(&t1, &t1, &t3)       // t1  = t1 - t3
		op.Sub(&t1, &t1, &t2)       // t1  = t1 - t2
		op.Mul(&t5, &t3, &params.A) // t5  = t3 * A24+
		op.Mul(&t3, &t3, &t5)       // t3  = t5 * t3
		op.Mul(&t6, &t2, &params.C) // t6  = t2 * A24-
		op.Mul(&t2, &t2, &t6)       // t2  = t2 * t6
		op.Sub(&t3, &t2, &t3)       // t3  = t2 - t3
		op.Sub(&t2, &t5, &t6)       // t2  = t5 - t6
		op.Mul(&t1, &t2, &t1)       // t1  = t2 * t1
		op.Add(&t2, &t3, &t1)       // t2  = t3 + t1
		op.Square(&t2, &t2)         // t2  = t2^2
		op.Mul(x, &t2, &t4)         // X3p = t2 * t4
		op.Sub(&t1, &t3, &t1)       // t1  = t3 - t1
		op.Square(&t1, &t1)         // t1  = t1^2
		op.Mul(z, &t1, &t0)         // Z3p = t1 * t0
	}
}

// Set (y1, y2, y3)  = (1/x1, 1/x2, 1/x3).
//
// All xi, yi must be distinct.
func (c *CurveOperations) Fp2Batch3Inv(x1, x2, x3, y1, y2, y3 *Fp2Element) {
	var x1x2, t Fp2Element
	var op = c.Params.Op

	op.Mul(&x1x2, x1, x2) // x1*x2
	op.Mul(&t, &x1x2, x3) // 1/(x1*x2*x3)
	op.Inv(&t, &t)
	op.Mul(y1, &t, x2) // 1/x1
	op.Mul(y1, y1, x3)
	op.Mul(y2, &t, x1) // 1/x2
	op.Mul(y2, y2, x3)
	op.Mul(y3, &t, &x1x2) // 1/x3
}

// ScalarMul3Pt is a right-to-left point multiplication that given the
// x-coordinate of P, Q and P-Q calculates the x-coordinate of R=Q+[scalar]P.
// nbits must be smaller or equal to len(scalar).
func (c *CurveOperations) ScalarMul3Pt(cparams *ProjectiveCurveParameters, P, Q, PmQ *ProjectivePoint, nbits uint, scalar []uint8) ProjectivePoint {
	var R0, R2, R1 ProjectivePoint
	var op = c.Params.Op
	aPlus2Over4 := c.CalcAplus2Over4(cparams)
	R1 = *P
	R2 = *PmQ
	R0 = *Q

	// Iterate over the bits of the scalar, bottom to top
	prevBit := uint8(0)
	for i := uint(0); i < nbits; i++ {
		bit := (scalar[i>>3] >> (i & 7) & 1)
		swap := prevBit ^ bit
		prevBit = bit
		op.CondSwap(&R1.X, &R1.Z, &R2.X, &R2.Z, swap)
		R0, R2 = c.xDblAdd(&R0, &R2, &R1, &aPlus2Over4)
	}
	op.CondSwap(&R1.X, &R1.Z, &R2.X, &R2.Z, prevBit)
	return R1
}

// Convert the input to wire format.
//
// The output byte slice must be at least 2*bytelen(p) bytes long.
func (c *CurveOperations) Fp2ToBytes(output []byte, fp2 *Fp2Element) {
	if len(output) < 2*c.Params.Bytelen {
		panic("output byte slice too short")
	}
	var a Fp2Element
	c.Params.Op.FromMontgomery(fp2, &a)

	// convert to bytes in little endian form
	for i := 0; i < c.Params.Bytelen; i++ {
		// set i = j*8 + k
		fp2 := i / 8
		k := uint64(i % 8)
		output[i] = byte(a.A[fp2] >> (8 * k))
		output[i+c.Params.Bytelen] = byte(a.B[fp2] >> (8 * k))
	}
}

// Read 2*bytelen(p) bytes into the given ExtensionFieldElement.
//
// It is an error to call this function if the input byte slice is less than 2*bytelen(p) bytes long.
func (c *CurveOperations) Fp2FromBytes(fp2 *Fp2Element, input []byte) {
	if len(input) < 2*c.Params.Bytelen {
		panic("input byte slice too short")
	}

	for i := 0; i < c.Params.Bytelen; i++ {
		j := i / 8
		k := uint64(i % 8)
		fp2.A[j] |= uint64(input[i]) << (8 * k)
		fp2.B[j] |= uint64(input[i+c.Params.Bytelen]) << (8 * k)
	}
	c.Params.Op.ToMontgomery(fp2)
}

/* -------------------------------------------------------------------------
   Mechnisms used for isogeny calculations
   -------------------------------------------------------------------------*/

// Constructs isogeny3 objects
func Newisogeny3(op FieldOps) Isogeny {
	return &isogeny3{Field: op}
}

// Constructs isogeny4 objects
func Newisogeny4(op FieldOps) Isogeny {
	return &isogeny4{isogeny3: isogeny3{Field: op}}
}

// Given a three-torsion point p = x(PB) on the curve E_(A:C), construct the
// three-isogeny phi : E_(A:C) -> E_(A:C)/<P_3> = E_(A':C').
//
// Input: (XP_3: ZP_3), where P_3 has exact order 3 on E_A/C
// Output: * Curve coordinates (A' + 2C', A' - 2C') corresponding to E_A'/C' = A_E/C/<P3>
//         * Isogeny phi with constants in F_p^2
func (phi *isogeny3) GenerateCurve(p *ProjectivePoint) CurveCoefficientsEquiv {
	var t0, t1, t2, t3, t4 Fp2Element
	var coefEq CurveCoefficientsEquiv
	var K1, K2 = &phi.K1, &phi.K2

	op := phi.Field
	op.Sub(K1, &p.X, &p.Z)            // K1 = XP3 - ZP3
	op.Square(&t0, K1)                // t0 = K1^2
	op.Add(K2, &p.X, &p.Z)            // K2 = XP3 + ZP3
	op.Square(&t1, K2)                // t1 = K2^2
	op.Add(&t2, &t0, &t1)             // t2 = t0 + t1
	op.Add(&t3, K1, K2)               // t3 = K1 + K2
	op.Square(&t3, &t3)               // t3 = t3^2
	op.Sub(&t3, &t3, &t2)             // t3 = t3 - t2
	op.Add(&t2, &t1, &t3)             // t2 = t1 + t3
	op.Add(&t3, &t3, &t0)             // t3 = t3 + t0
	op.Add(&t4, &t3, &t0)             // t4 = t3 + t0
	op.Add(&t4, &t4, &t4)             // t4 = t4 + t4
	op.Add(&t4, &t1, &t4)             // t4 = t1 + t4
	op.Mul(&coefEq.C, &t2, &t4)       // A24m = t2 * t4
	op.Add(&t4, &t1, &t2)             // t4 = t1 + t2
	op.Add(&t4, &t4, &t4)             // t4 = t4 + t4
	op.Add(&t4, &t0, &t4)             // t4 = t0 + t4
	op.Mul(&t4, &t3, &t4)             // t4 = t3 * t4
	op.Sub(&t0, &t4, &coefEq.C)       // t0 = t4 - A24m
	op.Add(&coefEq.A, &coefEq.C, &t0) // A24p = A24m + t0
	return coefEq
}

// Given a 3-isogeny phi and a point pB = x(PB), compute x(QB), the x-coordinate
// of the image QB = phi(PB) of PB under phi : E_(A:C) -> E_(A':C').
//
// The output xQ = x(Q) is then a point on the curve E_(A':C'); the curve
// parameters are returned by the GenerateCurve function used to construct phi.
func (phi *isogeny3) EvaluatePoint(p *ProjectivePoint) ProjectivePoint {
	var t0, t1, t2 Fp2Element
	var q ProjectivePoint
	var K1, K2 = &phi.K1, &phi.K2
	var px, pz = &p.X, &p.Z

	op := phi.Field
	op.Add(&t0, px, pz)   // t0 = XQ + ZQ
	op.Sub(&t1, px, pz)   // t1 = XQ - ZQ
	op.Mul(&t0, K1, &t0)  // t2 = K1 * t0
	op.Mul(&t1, K2, &t1)  // t1 = K2 * t1
	op.Add(&t2, &t0, &t1) // t2 = t0 + t1
	op.Sub(&t0, &t1, &t0) // t0 = t1 - t0
	op.Square(&t2, &t2)   // t2 = t2 ^ 2
	op.Square(&t0, &t0)   // t0 = t0 ^ 2
	op.Mul(&q.X, px, &t2) // XQ'= XQ * t2
	op.Mul(&q.Z, pz, &t0) // ZQ'= ZQ * t0
	return q
}

// Given a four-torsion point p = x(PB) on the curve E_(A:C), construct the
// four-isogeny phi : E_(A:C) -> E_(A:C)/<P_4> = E_(A':C').
//
// Input: (XP_4: ZP_4), where P_4 has exact order 4 on E_A/C
// Output: * Curve coordinates (A' + 2C', 4C') corresponding to E_A'/C' = A_E/C/<P4>
//         * Isogeny phi with constants in F_p^2
func (phi *isogeny4) GenerateCurve(p *ProjectivePoint) CurveCoefficientsEquiv {
	var coefEq CurveCoefficientsEquiv
	var xp4, zp4 = &p.X, &p.Z
	var K1, K2, K3 = &phi.K1, &phi.K2, &phi.K3

	op := phi.Field
	op.Sub(K2, xp4, zp4)
	op.Add(K3, xp4, zp4)
	op.Square(K1, zp4)
	op.Add(K1, K1, K1)
	op.Square(&coefEq.C, K1)
	op.Add(K1, K1, K1)
	op.Square(&coefEq.A, xp4)
	op.Add(&coefEq.A, &coefEq.A, &coefEq.A)
	op.Square(&coefEq.A, &coefEq.A)
	return coefEq
}

// Given a 4-isogeny phi and a point xP = x(P), compute x(Q), the x-coordinate
// of the image Q = phi(P) of P under phi : E_(A:C) -> E_(A':C').
//
// Input: Isogeny returned by GenerateCurve and point q=(Qx,Qz) from E0_A/C
// Output: Corresponding point q from E1_A'/C', where E1 is 4-isogenous to E0
func (phi *isogeny4) EvaluatePoint(p *ProjectivePoint) ProjectivePoint {
	var t0, t1 Fp2Element
	var q = *p
	var xq, zq = &q.X, &q.Z
	var K1, K2, K3 = &phi.K1, &phi.K2, &phi.K3

	op := phi.Field
	op.Add(&t0, xq, zq)
	op.Sub(&t1, xq, zq)
	op.Mul(xq, &t0, K2)
	op.Mul(zq, &t1, K3)
	op.Mul(&t0, &t0, &t1)
	op.Mul(&t0, &t0, K1)
	op.Add(&t1, xq, zq)
	op.Sub(zq, xq, zq)
	op.Square(&t1, &t1)
	op.Square(zq, zq)
	op.Add(xq, &t0, &t1)
	op.Sub(&t0, zq, &t0)
	op.Mul(xq, xq, &t1)
	op.Mul(zq, zq, &t0)
	return q
}

/* -------------------------------------------------------------------------
   Utils
   -------------------------------------------------------------------------*/
func (point *ProjectivePoint) ToAffine(c *CurveOperations) *Fp2Element {
	var affine_x Fp2Element
	c.Params.Op.Inv(&affine_x, &point.Z)
	c.Params.Op.Mul(&affine_x, &affine_x, &point.X)
	return &affine_x
}

// Cleans data in fp
func (fp *Fp2Element) Zeroize() {
	// Zeroizing in 2 seperated loops tells compiler to
	// use fast runtime.memclr()
	for i := range fp.A {
		fp.A[i] = 0
	}
	for i := range fp.B {
		fp.B[i] = 0
	}
}
