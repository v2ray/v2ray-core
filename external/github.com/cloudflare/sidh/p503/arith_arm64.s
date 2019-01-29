// +build arm64,!noasm

#include "textflag.h"

TEXT ·fp503ConditionalSwap(SB), NOSPLIT, $0-17
	MOVD	x+0(FP), R0
	MOVD	y+8(FP), R1
	MOVB	choice+16(FP), R2

	// Set flags
	// If choice is not 0 or 1, this implementation will swap completely
	CMP	$0, R2

	LDP	0(R0), (R3, R4)
	LDP	0(R1), (R5, R6)
	CSEL	EQ, R3, R5, R7
	CSEL	EQ, R4, R6, R8
	STP	(R7, R8), 0(R0)
	CSEL	NE, R3, R5, R9
	CSEL	NE, R4, R6, R10
	STP	(R9, R10), 0(R1)

	LDP	16(R0), (R3, R4)
	LDP	16(R1), (R5, R6)
	CSEL	EQ, R3, R5, R7
	CSEL	EQ, R4, R6, R8
	STP	(R7, R8), 16(R0)
	CSEL	NE, R3, R5, R9
	CSEL	NE, R4, R6, R10
	STP	(R9, R10), 16(R1)

	LDP	32(R0), (R3, R4)
	LDP	32(R1), (R5, R6)
	CSEL	EQ, R3, R5, R7
	CSEL	EQ, R4, R6, R8
	STP	(R7, R8), 32(R0)
	CSEL	NE, R3, R5, R9
	CSEL	NE, R4, R6, R10
	STP	(R9, R10), 32(R1)

	LDP	48(R0), (R3, R4)
	LDP	48(R1), (R5, R6)
	CSEL	EQ, R3, R5, R7
	CSEL	EQ, R4, R6, R8
	STP	(R7, R8), 48(R0)
	CSEL	NE, R3, R5, R9
	CSEL	NE, R4, R6, R10
	STP	(R9, R10), 48(R1)

	RET

TEXT ·fp503AddReduced(SB), NOSPLIT, $0-24
	MOVD	z+0(FP), R2
	MOVD	x+8(FP), R0
	MOVD	y+16(FP), R1

	// Load first summand into R3-R10
	// Add first summand and second summand and store result in R3-R10
	LDP	0(R0), (R3, R4)
	LDP	0(R1), (R11, R12)
	LDP	16(R0), (R5, R6)
	LDP	16(R1), (R13, R14)
	ADDS	R11, R3
	ADCS	R12, R4
	ADCS	R13, R5
	ADCS	R14, R6

	LDP	32(R0), (R7, R8)
	LDP	32(R1), (R11, R12)
	LDP	48(R0), (R9, R10)
	LDP	48(R1), (R13, R14)
	ADCS	R11, R7
	ADCS	R12, R8
	ADCS	R13, R9
	ADC	R14, R10

	// Subtract 2 * p503 in R11-R17 from the result in R3-R10
	LDP	·p503x2+0(SB), (R11, R12)
	LDP	·p503x2+24(SB), (R13, R14)
	SUBS	R11, R3
	SBCS	R12, R4
	LDP	·p503x2+40(SB), (R15, R16)
	SBCS	R12, R5
	SBCS	R13, R6
	MOVD	·p503x2+56(SB), R17
	SBCS	R14, R7
	SBCS	R15, R8
	SBCS	R16, R9
	SBCS	R17, R10
	SBC	ZR, ZR, R19

	// If x + y - 2 * p503 < 0, R19 is 1 and 2 * p503 should be added
	AND	R19, R11
	AND	R19, R12
	AND	R19, R13
	AND	R19, R14
	AND	R19, R15
	AND	R19, R16
	AND	R19, R17

	ADDS	R11, R3
	ADCS	R12, R4
	STP	(R3, R4), 0(R2)
	ADCS	R12, R5
	ADCS	R13, R6
	STP	(R5, R6), 16(R2)
	ADCS	R14, R7
	ADCS	R15, R8
	STP	(R7, R8), 32(R2)
	ADCS	R16, R9
	ADC	R17, R10
	STP	(R9, R10), 48(R2)

	RET

TEXT ·fp503SubReduced(SB), NOSPLIT, $0-24
	MOVD	z+0(FP), R2
	MOVD	x+8(FP), R0
	MOVD	y+16(FP), R1

	// Load x into R3-R10
	// Subtract y from x and store result in R3-R10
	LDP	0(R0), (R3, R4)
	LDP	0(R1), (R11, R12)
	LDP	16(R0), (R5, R6)
	LDP	16(R1), (R13, R14)
	SUBS	R11, R3
	SBCS	R12, R4
	SBCS	R13, R5
	SBCS	R14, R6

	LDP	32(R0), (R7, R8)
	LDP	32(R1), (R11, R12)
	LDP	48(R0), (R9, R10)
	LDP	48(R1), (R13, R14)
	SBCS	R11, R7
	SBCS	R12, R8
	SBCS	R13, R9
	SBCS	R14, R10
	SBC	ZR, ZR, R19

	// If x - y < 0, R19 is 1 and 2 * p503 should be added
	LDP	·p503x2+0(SB), (R11, R12)
	LDP	·p503x2+24(SB), (R13, R14)
	AND	R19, R11
	AND	R19, R12
	LDP	·p503x2+40(SB), (R15, R16)
	AND	R19, R13
	AND	R19, R14
	MOVD	·p503x2+56(SB), R17
	AND	R19, R15
	AND	R19, R16
	AND	R19, R17

	ADDS	R11, R3
	ADCS	R12, R4
	STP	(R3, R4), 0(R2)
	ADCS	R12, R5
	ADCS	R13, R6
	STP	(R5, R6), 16(R2)
	ADCS	R14, R7
	ADCS	R15, R8
	STP	(R7, R8), 32(R2)
	ADCS	R16, R9
	ADC	R17, R10
	STP	(R9, R10), 48(R2)

	RET

TEXT ·fp503AddLazy(SB), NOSPLIT, $0-24
	MOVD	z+0(FP), R2
	MOVD	x+8(FP), R0
	MOVD	y+16(FP), R1

	// Load first summand into R3-R10
	// Add first summand and second summand and store result in R3-R10
	LDP	0(R0), (R3, R4)
	LDP	0(R1), (R11, R12)
	LDP	16(R0), (R5, R6)
	LDP	16(R1), (R13, R14)
	ADDS	R11, R3
	ADCS	R12, R4
	STP	(R3, R4), 0(R2)
	ADCS	R13, R5
	ADCS	R14, R6
	STP	(R5, R6), 16(R2)

	LDP	32(R0), (R7, R8)
	LDP	32(R1), (R11, R12)
	LDP	48(R0), (R9, R10)
	LDP	48(R1), (R13, R14)
	ADCS	R11, R7
	ADCS	R12, R8
	STP	(R7, R8), 32(R2)
	ADCS	R13, R9
	ADC	R14, R10
	STP	(R9, R10), 48(R2)

	RET

TEXT ·fp503X2AddLazy(SB), NOSPLIT, $0-24
	MOVD	z+0(FP), R2
	MOVD	x+8(FP), R0
	MOVD	y+16(FP), R1

	LDP	0(R0), (R3, R4)
	LDP	0(R1), (R11, R12)
	LDP	16(R0), (R5, R6)
	LDP	16(R1), (R13, R14)
	ADDS	R11, R3
	ADCS	R12, R4
	STP	(R3, R4), 0(R2)
	ADCS	R13, R5
	ADCS	R14, R6
	STP	(R5, R6), 16(R2)

	LDP	32(R0), (R7, R8)
	LDP	32(R1), (R11, R12)
	LDP	48(R0), (R9, R10)
	LDP	48(R1), (R13, R14)
	ADCS	R11, R7
	ADCS	R12, R8
	STP	(R7, R8), 32(R2)
	ADCS	R13, R9
	ADCS	R14, R10
	STP	(R9, R10), 48(R2)

	LDP	64(R0), (R3, R4)
	LDP	64(R1), (R11, R12)
	LDP	80(R0), (R5, R6)
	LDP	80(R1), (R13, R14)
	ADCS	R11, R3
	ADCS	R12, R4
	STP	(R3, R4), 64(R2)
	ADCS	R13, R5
	ADCS	R14, R6
	STP	(R5, R6), 80(R2)

	LDP	96(R0), (R7, R8)
	LDP	96(R1), (R11, R12)
	LDP	112(R0), (R9, R10)
	LDP	112(R1), (R13, R14)
	ADCS	R11, R7
	ADCS	R12, R8
	STP	(R7, R8), 96(R2)
	ADCS	R13, R9
	ADC	R14, R10
	STP	(R9, R10), 112(R2)

	RET

TEXT ·fp503X2SubLazy(SB), NOSPLIT, $0-24
	MOVD	z+0(FP), R2
	MOVD	x+8(FP), R0
	MOVD	y+16(FP), R1

	LDP	0(R0), (R3, R4)
	LDP	0(R1), (R11, R12)
	LDP	16(R0), (R5, R6)
	LDP	16(R1), (R13, R14)
	SUBS	R11, R3
	SBCS	R12, R4
	STP	(R3, R4), 0(R2)
	SBCS	R13, R5
	SBCS	R14, R6
	STP	(R5, R6), 16(R2)

	LDP	32(R0), (R7, R8)
	LDP	32(R1), (R11, R12)
	LDP	48(R0), (R9, R10)
	LDP	48(R1), (R13, R14)
	SBCS	R11, R7
	SBCS	R12, R8
	STP	(R7, R8), 32(R2)
	SBCS	R13, R9
	SBCS	R14, R10
	STP	(R9, R10), 48(R2)

	LDP	64(R0), (R3, R4)
	LDP	64(R1), (R11, R12)
	LDP	80(R0), (R5, R6)
	LDP	80(R1), (R13, R14)
	SBCS	R11, R3
	SBCS	R12, R4
	SBCS	R13, R5
	SBCS	R14, R6

	LDP	96(R0), (R7, R8)
	LDP	96(R1), (R11, R12)
	LDP	112(R0), (R9, R10)
	LDP	112(R1), (R13, R14)
	SBCS	R11, R7
	SBCS	R12, R8
	SBCS	R13, R9
	SBCS	R14, R10
	SBC	ZR, ZR, R15

	// If x - y < 0, R15 is 1 and p503 should be added
	LDP	·p503+16(SB), (R16, R17)
	LDP	·p503+32(SB), (R19, R20)
	AND	R15, R16
	AND	R15, R17
	LDP	·p503+48(SB), (R21, R22)
	AND	R15, R19
	AND	R15, R20
	AND	R15, R21
	AND	R15, R22

	ADDS	R16, R3
	ADCS	R16, R4
	STP	(R3, R4), 64(R2)
	ADCS	R16, R5
	ADCS	R17, R6
	STP	(R5, R6), 80(R2)
	ADCS	R19, R7
	ADCS	R20, R8
	STP	(R7, R8), 96(R2)
	ADCS	R21, R9
	ADC	R22, R10
	STP	(R9, R10), 112(R2)

	RET

// Expects that X0*Y0 is already in Z0(low),Z3(high) and X0*Y1 in Z1(low),Z2(high)
// Z0 is not actually touched
// Result of (X0-X1) * (Y0-Y1) will be in Z0-Z3
// Inputs get overwritten, except for X1
#define mul128x128comba(X0, X1, Y0, Y1, Z0, Z1, Z2, Z3, T0) \
	MUL	X1, Y0, X0	\
	UMULH	X1, Y0, Y0	\
	ADDS	Z3, Z1		\
	ADC	ZR, Z2		\
				\
	MUL	Y1, X1, T0	\
	UMULH	Y1, X1, Y1	\
	ADDS	X0, Z1		\
	ADCS	Y0, Z2		\
	ADC	ZR, ZR, Z3	\
				\
	ADDS	T0, Z2		\
	ADC	Y1, Z3

// Expects that X points to (X0-X1)
// Result of (X0-X3) * (Y0-Y3) will be in Z0-Z7
// Inputs get overwritten, except X2-X3 and Y2-Y3
#define mul256x256karatsuba(X, X0, X1, X2, X3, Y0, Y1, Y2, Y3, Z0, Z1, Z2, Z3, Z4, Z5, Z6, Z7, T0, T1)\
	ADDS	X2, X0		\	// xH + xL, destroys xL
	ADCS	X3, X1		\
	ADCS	ZR, ZR, T0	\
				\
	ADDS	Y2, Y0, Z6	\	// yH + yL
	ADCS	Y3, Y1, T1	\
	ADC	ZR, ZR, Z7	\
				\
	SUB	T0, ZR, Z2	\
	SUB	Z7, ZR, Z3	\
	AND	Z7, T0		\	// combined carry
				\
	AND	Z2, Z6, Z0	\	// masked(yH + yL)
	AND	Z2, T1, Z1	\
				\
	AND	Z3, X0, Z4	\	// masked(xH + xL)
	AND	Z3, X1, Z5	\
				\
	MUL	Z6, X0, Z2	\
	MUL	T1, X0, Z3	\
				\
	ADDS	Z4, Z0		\
	UMULH	T1, X0, Z4	\
	ADCS	Z5, Z1		\
	UMULH	Z6, X0, Z5	\
	ADC	ZR, T0		\
				\	// (xH + xL) * (yH + yL)
	mul128x128comba(X0, X1, Z6, T1, Z2, Z3, Z4, Z5, Z7)\
				\
	LDP 0+X, (X0, X1)	\
				\
	ADDS	Z0, Z4		\
	UMULH	Y0, X0, Z7	\
	UMULH	Y1, X0, T1	\
	ADCS	Z1, Z5		\
	MUL	Y0, X0, Z0	\
	MUL	Y1, X0, Z1	\
	ADC	ZR, T0		\
				\	// xL * yL
	mul128x128comba(X0, X1, Y0, Y1, Z0, Z1, T1, Z7, Z6)\
				\
	MUL	Y2, X2, X0	\
	UMULH	Y2, X2, Y0	\
	SUBS	Z0, Z2		\	// (xH + xL) * (yH + yL) - xL * yL
	SBCS	Z1, Z3		\
	SBCS	T1, Z4		\
	MUL	Y3, X2, X1	\
	UMULH	Y3, X2, Z6	\
	SBCS	Z7, Z5		\
	SBCS	ZR, T0		\
				\	// xH * yH
	mul128x128comba(X2, X3, Y2, Y3, X0, X1, Z6, Y0, Y1)\
				\
	SUBS	X0, Z2		\	// (xH + xL) * (yH + yL) - xL * yL - xH * yH
	SBCS	X1, Z3		\
	SBCS	Z6, Z4		\
	SBCS	Y0, Z5		\
	SBCS	ZR, T0		\
				\
	ADDS	T1, Z2		\	// (xH * yH) * 2^256 + ((xH + xL) * (yH + yL) - xL * yL - xH * yH) * 2^128 + xL * yL
	ADCS	Z7, Z3		\
	ADCS	X0, Z4		\
	ADCS	X1, Z5		\
	ADCS	T0, Z6		\
	ADC	Y0, ZR, Z7


// This implements two-level Karatsuba with a 128x128 Comba multiplier
// at the bottom
TEXT ·fp503Mul(SB), NOSPLIT, $0-24
	MOVD	z+0(FP), R2
	MOVD	x+8(FP), R0
	MOVD	y+16(FP), R1

	// Load xL in R3-R6, xH in R7-R10
	// (xH + xL) in R25-R29
	LDP	0(R0), (R3, R4)
	LDP	32(R0), (R7, R8)
	ADDS	R3, R7, R25
	ADCS	R4, R8, R26
	LDP	16(R0), (R5, R6)
	LDP	48(R0), (R9, R10)
	ADCS	R5, R9, R27
	ADCS	R6, R10, R29
	ADC	ZR, ZR, R7

	// Load yL in R11-R14, yH in R15-19
	// (yH + yL) in R11-R14, destroys yL
	LDP	0(R1), (R11, R12)
	LDP	32(R1), (R15, R16)
	ADDS	R15, R11
	ADCS	R16, R12
	LDP	16(R1), (R13, R14)
	LDP	48(R1), (R17, R19)
	ADCS	R17, R13
	ADCS	R19, R14
	ADC	ZR, ZR, R8

	// Compute maskes and combined carry
	SUB	R7, ZR, R9
	SUB	R8, ZR, R10
	AND	R8, R7

	// masked(yH + yL)
	AND	R9, R11, R15
	AND	R9, R12, R16
	AND	R9, R13, R17
	AND	R9, R14, R19

	// masked(xH + xL)
	AND	R10, R25, R20
	AND	R10, R26, R21
	AND	R10, R27, R22
	AND	R10, R29, R23

	// masked(xH + xL) + masked(yH + yL) in R15-R19
	ADDS	R20, R15
	ADCS	R21, R16
	ADCS	R22, R17
	ADCS	R23, R19
	ADC	ZR, R7

	// Use z as temporary storage
	STP	(R25, R26), 0(R2)

	// (xH + xL) * (yH + yL)
	mul256x256karatsuba(0(R2), R25, R26, R27, R29, R11, R12, R13, R14, R8, R9, R10, R20, R21, R22, R23, R24, R0, R1)

	MOVD	x+8(FP), R0
	MOVD	y+16(FP), R1

	ADDS	R21, R15
	ADCS	R22, R16
	ADCS	R23, R17
	ADCS	R24, R19
	ADC	ZR, R7

	// Load yL in R11-R14
	LDP	0(R1), (R11, R12)
	LDP	16(R1), (R13, R14)

	// xL * yL
	mul256x256karatsuba(0(R0), R3, R4, R5, R6, R11, R12, R13, R14, R21, R22, R23, R24, R25, R26, R27, R29, R1, R2)

	MOVD	z+0(FP), R2
	MOVD	y+16(FP), R1

	// (xH + xL) * (yH + yL) - xL * yL
	SUBS	R21, R8
	SBCS	R22, R9
	STP	(R21, R22), 0(R2)
	SBCS	R23, R10
	SBCS	R24, R20
	STP	(R23, R24), 16(R2)
	SBCS	R25, R15
	SBCS	R26, R16
	SBCS	R27, R17
	SBCS	R29, R19
	SBC	ZR, R7

	// Load xH in R3-R6, yH in R11-R14
	LDP	32(R0), (R3, R4)
	LDP	48(R0), (R5, R6)
	LDP	32(R1), (R11, R12)
	LDP	48(R1), (R13, R14)

	ADDS	R25, R8
	ADCS	R26, R9
	ADCS	R27, R10
	ADCS	R29, R20
	ADC	ZR, ZR, R1

	MOVD	R20, 32(R2)

	// xH * yH
	mul256x256karatsuba(32(R0), R3, R4, R5, R6, R11, R12, R13, R14, R21, R22, R23, R24, R25, R26, R27, R29, R2, R20)
	NEG	R1, R1

	MOVD	z+0(FP), R2
	MOVD	32(R2), R20

	// (xH + xL) * (yH + yL) - xL * yL - xH * yH in R8-R10,R20,R15-R19
	// Store lower half in z, that's done
	SUBS	R21, R8
	SBCS	R22, R9
	STP	(R8, R9), 32(R2)
	SBCS	R23, R10
	SBCS	R24, R20
	STP	(R10, R20), 48(R2)
	SBCS	R25, R15
	SBCS	R26, R16
	SBCS	R27, R17
	SBCS	R29, R19
	SBC	ZR, R7

	// (xH * yH) * 2^512 + ((xH + xL) * (yH + yL) - xL * yL - xH * yH) * 2^256 + xL * yL
	// Store remaining limbs in z
	ADDS	$1, R1
	ADCS	R21, R15
	ADCS	R22, R16
	STP	(R15, R16), 64(R2)
	ADCS	R23, R17
	ADCS	R24, R19
	STP	(R17, R19), 80(R2)
	ADCS	R7, R25
	ADCS	ZR, R26
	STP	(R25, R26), 96(R2)
	ADCS	ZR, R27
	ADC	ZR, R29
	STP	(R27, R29), 112(R2)

	RET

// Expects that X0*Y0 is already in Z0(low),Z3(high) and X0*Y1 in Z1(low),Z2(high)
// Z0 is not actually touched
// Result of (X0-X1) * (Y0-Y3) will be in Z0-Z5
// Inputs remain intact
#define mul128x256comba(X0, X1, Y0, Y1, Y2, Y3, Z0, Z1, Z2, Z3, Z4, Z5, T0, T1, T2, T3)\
	MUL	X1, Y0, T0	\
	UMULH	X1, Y0, T1	\
	ADDS	Z3, Z1		\
	ADC	ZR, Z2		\
				\
	MUL	X0, Y2, T2	\
	UMULH	X0, Y2, T3	\
	ADDS	T0, Z1		\
	ADCS	T1, Z2		\
	ADC	ZR, ZR, Z3	\
				\
	MUL	X1, Y1, T0	\
	UMULH	X1, Y1, T1	\
	ADDS	T2, Z2		\
	ADCS	T3, Z3		\
	ADC	ZR, ZR, Z4	\
				\
	MUL	X0, Y3, T2	\
	UMULH	X0, Y3, T3	\
	ADDS	T0, Z2		\
	ADCS	T1, Z3		\
	ADC	ZR, Z4		\
				\
	MUL	X1, Y2, T0	\
	UMULH	X1, Y2, T1	\
	ADDS	T2, Z3		\
	ADCS	T3, Z4		\
	ADC	ZR, ZR, Z5	\
				\
	MUL	X1, Y3, T2	\
	UMULH	X1, Y3, T3	\
	ADDS	T0, Z3		\
	ADCS	T1, Z4		\
	ADC	ZR, Z5		\
	ADDS	T2, Z4		\
	ADC	T3, Z5

// This implements the shifted 2^(B*w) Montgomery reduction from
// https://eprint.iacr.org/2016/986.pdf, section Section 3.2, with
// B = 4, w = 64. Performance results were reported in
// https://eprint.iacr.org/2018/700.pdf Section 6.
TEXT ·fp503MontgomeryReduce(SB), NOSPLIT, $0-16
	MOVD	x+8(FP), R0

	// Load x0-x1
	LDP	0(R0), (R2, R3)

	// Load the prime constant in R25-R29
	LDP	·p503p1s8+32(SB), (R25, R26)
	LDP	·p503p1s8+48(SB), (R27, R29)

	// [x0,x1] * p503p1s8 to R4-R9
	MUL	R2, R25, R4		// x0 * p503p1s8[0]
	UMULH	R2, R25, R7
	MUL	R2, R26, R5		// x0 * p503p1s8[1]
	UMULH	R2, R26, R6

	mul128x256comba(R2, R3, R25, R26, R27, R29, R4, R5, R6, R7, R8, R9, R10, R11, R12, R13)

	LDP	16(R0), (R3, R11)	// x2
	LDP	32(R0), (R12, R13)
	LDP	48(R0), (R14, R15)

	// Left-shift result in R4-R9 by 56 to R4-R10
	ORR	R9>>8, ZR, R10
	LSL	$56, R9
	ORR	R8>>8, R9
	LSL	$56, R8
	ORR	R7>>8, R8
	LSL	$56, R7
	ORR	R6>>8, R7
	LSL	$56, R6
	ORR	R5>>8, R6
	LSL	$56, R5
	ORR	R4>>8, R5
	LSL	$56, R4

	ADDS	R4, R11			// x3
	ADCS	R5, R12			// x4
	ADCS	R6, R13
	ADCS	R7, R14
	ADCS	R8, R15
	LDP	64(R0), (R16, R17)
	LDP	80(R0), (R19, R20)
	MUL	R3, R25, R4		// x2 * p503p1s8[0]
	UMULH	R3, R25, R7
	ADCS	R9, R16
	ADCS	R10, R17
	ADCS	ZR, R19
	ADCS	ZR, R20
	LDP	96(R0), (R21, R22)
	LDP	112(R0), (R23, R24)
	MUL	R3, R26, R5		// x2 * p503p1s8[1]
	UMULH	R3, R26, R6
	ADCS	ZR, R21
	ADCS	ZR, R22
	ADCS	ZR, R23
	ADC	ZR, R24

	// [x2,x3] * p503p1s8 to R4-R9
	mul128x256comba(R3, R11, R25, R26, R27, R29, R4, R5, R6, R7, R8, R9, R10, R0, R1, R2)

	ORR	R9>>8, ZR, R10
	LSL	$56, R9
	ORR	R8>>8, R9
	LSL	$56, R8
	ORR	R7>>8, R8
	LSL	$56, R7
	ORR	R6>>8, R7
	LSL	$56, R6
	ORR	R5>>8, R6
	LSL	$56, R5
	ORR	R4>>8, R5
	LSL	$56, R4

	ADDS	R4, R13			// x5
	ADCS	R5, R14			// x6
	ADCS	R6, R15
	ADCS	R7, R16
	MUL	R12, R25, R4		// x4 * p503p1s8[0]
	UMULH	R12, R25, R7
	ADCS	R8, R17
	ADCS	R9, R19
	ADCS	R10, R20
	ADCS	ZR, R21
	MUL	R12, R26, R5		// x4 * p503p1s8[1]
	UMULH	R12, R26, R6
	ADCS	ZR, R22
	ADCS	ZR, R23
	ADC	ZR, R24

	// [x4,x5] * p503p1s8 to R4-R9
	mul128x256comba(R12, R13, R25, R26, R27, R29, R4, R5, R6, R7, R8, R9, R10, R0, R1, R2)

	ORR	R9>>8, ZR, R10
	LSL	$56, R9
	ORR	R8>>8, R9
	LSL	$56, R8
	ORR	R7>>8, R8
	LSL	$56, R7
	ORR	R6>>8, R7
	LSL	$56, R6
	ORR	R5>>8, R6
	LSL	$56, R5
	ORR	R4>>8, R5
	LSL	$56, R4

	ADDS	R4, R15			// x7
	ADCS	R5, R16			// x8
	ADCS	R6, R17
	ADCS	R7, R19
	MUL	R14, R25, R4		// x6 * p503p1s8[0]
	UMULH	R14, R25, R7
	ADCS	R8, R20
	ADCS	R9, R21
	ADCS	R10, R22
	MUL	R14, R26, R5		// x6 * p503p1s8[1]
	UMULH	R14, R26, R6
	ADCS	ZR, R23
	ADC	ZR, R24

	// [x6,x7] * p503p1s8 to R4-R9
	mul128x256comba(R14, R15, R25, R26, R27, R29, R4, R5, R6, R7, R8, R9, R10, R0, R1, R2)

	ORR	R9>>8, ZR, R10
	LSL	$56, R9
	ORR	R8>>8, R9
	LSL	$56, R8
	ORR	R7>>8, R8
	LSL	$56, R7
	ORR	R6>>8, R7
	LSL	$56, R6
	ORR	R5>>8, R6
	LSL	$56, R5
	ORR	R4>>8, R5
	LSL	$56, R4

	MOVD	z+0(FP), R0
	ADDS	R4, R17
	ADCS	R5, R19
	STP	(R16, R17),  0(R0)	// Store final result to z
	ADCS	R6, R20
	ADCS	R7, R21
	STP	(R19, R20), 16(R0)
	ADCS	R8, R22
	ADCS	R9, R23
	STP	(R21, R22), 32(R0)
	ADC	R10, R24
	STP	(R23, R24), 48(R0)

	RET

TEXT ·fp503StrongReduce(SB), NOSPLIT, $0-8
	MOVD	x+0(FP), R0

	// Keep x in R1-R8, p503 in R9-R14, subtract to R1-R8
	LDP	·p503+16(SB), (R9, R10)
	LDP	0(R0), (R1, R2)
	LDP	16(R0), (R3, R4)
	SUBS	R9, R1
	SBCS	R9, R2

	LDP	32(R0), (R5, R6)
	LDP	·p503+32(SB), (R11, R12)
	SBCS	R9, R3
	SBCS	R10, R4

	LDP	48(R0), (R7, R8)
	LDP	·p503+48(SB), (R13, R14)
	SBCS	R11, R5
	SBCS	R12, R6

	SBCS	R13, R7
	SBCS	R14, R8
	SBC	ZR, ZR, R15

	// Mask with the borrow and add p503
	AND	R15, R9
	AND	R15, R10
	AND	R15, R11
	AND	R15, R12
	AND	R15, R13
	AND	R15, R14

	ADDS	R9, R1
	ADCS	R9, R2
	STP	(R1, R2), 0(R0)
	ADCS	R9, R3
	ADCS	R10, R4
	STP	(R3, R4), 16(R0)
	ADCS	R11, R5
	ADCS	R12, R6
	STP	(R5, R6), 32(R0)
	ADCS	R13, R7
	ADCS	R14, R8
	STP	(R7, R8), 48(R0)

	RET
