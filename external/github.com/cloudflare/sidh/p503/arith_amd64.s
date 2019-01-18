// +build amd64,!noasm

#include "textflag.h"

// p503
#define P503_0     $0xFFFFFFFFFFFFFFFF
#define P503_1     $0xFFFFFFFFFFFFFFFF
#define P503_2     $0xFFFFFFFFFFFFFFFF
#define P503_3     $0xABFFFFFFFFFFFFFF
#define P503_4     $0x13085BDA2211E7A0
#define P503_5     $0x1B9BF6C87B7E7DAF
#define P503_6     $0x6045C6BDDA77A4D0
#define P503_7     $0x004066F541811E1E

// p503+1
#define P503P1_3   $0xAC00000000000000
#define P503P1_4   $0x13085BDA2211E7A0
#define P503P1_5   $0x1B9BF6C87B7E7DAF
#define P503P1_6   $0x6045C6BDDA77A4D0
#define P503P1_7   $0x004066F541811E1E

// p503x2
#define P503X2_0   $0xFFFFFFFFFFFFFFFE
#define P503X2_1   $0xFFFFFFFFFFFFFFFF
#define P503X2_2   $0xFFFFFFFFFFFFFFFF
#define P503X2_3   $0x57FFFFFFFFFFFFFF
#define P503X2_4   $0x2610B7B44423CF41
#define P503X2_5   $0x3737ED90F6FCFB5E
#define P503X2_6   $0xC08B8D7BB4EF49A0
#define P503X2_7   $0x0080CDEA83023C3C

#define REG_P1 DI
#define REG_P2 SI
#define REG_P3 DX

// Performs schoolbook multiplication of 2 256-bit numbers. This optimized version
// uses MULX instruction. Macro smashes value in DX.
// Input: I0 and I1.
// Output: O
// All the other arguments are resgisters, used for storing temporary values
#define MULS256_MULX(O, I0, I1, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) \
	MOVQ    I0, DX          \
	MULXQ   I1, T1, T0      \   // T0:T1 = A0*B0
	MOVQ    T1, O           \   // O[0]
	MULXQ   8+I1, T2, T1    \   // T1:T2 = U0*V1
	ADDQ    T2, T0          \
	MULXQ   16+I1, T3, T2   \   // T2:T3 = U0*V2
	ADCQ    T3, T1          \
	MULXQ   24+I1, T4, T3   \   // T3:T4 = U0*V3
	ADCQ    T4, T2          \
	\ // Column U1
	MOVQ    8+I0, DX        \
	ADCQ    $0, T3          \
	MULXQ   0+I1, T4, T5    \   // T5:T4 = U1*V0
	MULXQ   8+I1, T7, T6    \   // T6:T7 = U1*V1
	ADDQ    T7, T5          \
	MULXQ   16+I1, T8, T7   \   // T7:T8 = U1*V2
	ADCQ    T8, T6          \
	MULXQ   24+I1, T9, T8   \   // T8:T9 = U1*V3
	ADCQ    T9, T7          \
	ADCQ    $0, T8          \
	ADDQ    T0, T4          \
	MOVQ    T4, 8+O         \   // O[1]
	ADCQ    T1, T5          \
	ADCQ    T2, T6          \
	ADCQ    T3, T7          \
	\ // Column U2
	MOVQ    16+I0, DX       \
	ADCQ    $0, T8          \
	MULXQ   0+I1, T0, T1    \   // T1:T0 = U2*V0
	MULXQ   8+I1, T3, T2    \   // T2:T3 = U2*V1
	ADDQ    T3, T1          \
	MULXQ   16+I1, T4, T3   \   // T3:T4 = U2*V2
	ADCQ    T4, T2          \
	MULXQ   24+I1, T9, T4   \   // T4:T9 = U2*V3
	ADCQ    T9, T3          \
	\ // Column U3
	MOVQ    24+I0, DX       \
	ADCQ    $0, T4          \
	ADDQ    T5, T0          \
	MOVQ    T0, 16+O        \   // O[2]
	ADCQ    T6, T1          \
	ADCQ    T7, T2          \
	ADCQ    T8, T3          \
	ADCQ    $0, T4          \
	MULXQ   0+I1, T0, T5    \   // T5:T0 = U3*V0
	MULXQ   8+I1, T7, T6    \   // T6:T7 = U3*V1
	ADDQ    T7, T5          \
	MULXQ   16+I1, T8, T7   \   // T7:T8 = U3*V2
	ADCQ    T8, T6          \
	MULXQ   24+I1, T9, T8   \   // T8:T9 = U3*V3
	ADCQ    T9, T7          \
	ADCQ    $0, T8          \
	\ // Add values in remaining columns
	ADDQ    T0, T1          \
	MOVQ    T1, 24+O        \   // O[3]
	ADCQ    T5, T2          \
	MOVQ    T2, 32+O        \   // O[4]
	ADCQ    T6, T3          \
	MOVQ    T3, 40+O        \   // O[5]
	ADCQ    T7, T4          \
	MOVQ    T4, 48+O        \   // O[6]
	ADCQ    $0, T8          \   // O[7]
	MOVQ    T8, 56+O

// Performs schoolbook multiplication of 2 256-bit numbers. This optimized version
// uses ADOX, ADCX and MULX instructions. Macro smashes values in AX and DX.
// Input: I0 and I1.
// Output: O
// All the other arguments resgisters are used for storing temporary values
#define MULS256_MULX_ADCX_ADOX(O, I0, I1, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) \
							\   // U0[0]
	MOVQ     0+I0, DX       \   // MULX requires multiplayer in DX
							\   // T0:T1 = I1*DX
	MULXQ    I1, T1, T0     \   // T0:T1 = U0*V0 (low:high)
	MOVQ     T1, O          \   // O0[0]
	MULXQ     8+I1, T2, T1  \   // T2:T1 = U0*V1
	XORQ     AX, AX         \
	ADOXQ    T2, T0         \
	MULXQ    16+I1, T3, T2  \   // T2:T3 = U0*V2
	ADOXQ    T3, T1         \
	MULXQ    24+I1, T4, T3  \   // T3:T4 = U0*V3
	ADOXQ    T4, T2         \
	\  // Column U1
	MOVQ      8+I0, DX      \
	MULXQ    I1, T4, T5     \   // T5:T4 = U1*V0
	ADOXQ    AX, T3         \
	XORQ     AX, AX         \
	MULXQ     8+I1, T7, T6  \   // T6:T7 = U1*V1
	ADOXQ    T0, T4         \
	MOVQ     T4, 8+O        \   // O[1]
	ADCXQ    T7, T5         \
	MULXQ    16+I1, T8, T7  \   // T7:T8 = U1*V2
	ADCXQ    T8, T6         \
	ADOXQ    T1, T5 \
	MULXQ    24+I1, T9, T8  \   // T8:T9 = U1*V3
	ADCXQ    T9, T7         \
	ADCXQ    AX, T8         \
	ADOXQ    T2, T6         \
	\ // Column U2
	MOVQ     16+I0, DX      \
	MULXQ    I1, T0, T1     \   // T1:T0 = U2*V0
	ADOXQ    T3, T7         \
	ADOXQ    AX, T8         \
	XORQ     AX, AX         \
	MULXQ    8+I1, T3, T2   \   // T2:T3 = U2*V1
	ADOXQ    T5, T0         \
	MOVQ     T0, 16+O       \   // O[2]
	ADCXQ    T3, T1         \
	MULXQ    16+I1, T4, T3  \   // T3:T4 = U2*V2
	ADCXQ    T4, T2         \
	ADOXQ    T6, T1         \
	MULXQ    24+I1, T9, T4  \   // T9:T4 = U2*V3
	ADCXQ    T9, T3         \
	MOVQ     24+I0, DX      \
	ADCXQ    AX, T4         \
	\
	ADOXQ    T7, T2         \
	ADOXQ    T8, T3         \
	ADOXQ    AX, T4         \
	\ // Column U3
	MULXQ    I1, T0, T5     \   // T5:T0 = U3*B0
	XORQ     AX, AX         \
	MULXQ    8+I1, T7, T6   \   // T6:T7 = U3*B1
	ADCXQ    T7, T5         \
	ADOXQ    T0, T1         \
	MULXQ    16+I1, T8, T7  \   // T7:T8 = U3*V2
	ADCXQ    T8, T6         \
	ADOXQ    T5, T2         \
	MULXQ    24+I1, T9, T8  \   // T8:T9 = U3*V3
	ADCXQ    T9, T7         \
	ADCXQ    AX, T8         \
	\
	ADOXQ   T6, T3          \
	ADOXQ   T7, T4          \
	ADOXQ   AX, T8          \
	MOVQ    T1, 24+O        \   // O[3]
	MOVQ    T2, 32+O        \   // O[4]
	MOVQ    T3, 40+O        \   // O[5]
	MOVQ    T4, 48+O        \   // O[6] and O[7] below
	MOVQ    T8, 56+O

// Template of a macro that performs schoolbook multiplication of 128-bit with 320-bit
// number. It uses MULX instruction This template must be customized with functions
// performing ADD (add1, add2) and ADD-with-carry (adc1, adc2). addX/adcX may or may
// not be instructions that use two independent carry chains.
// Input:
//   * I0 128-bit number
//   * I1 320-bit number
//   * add1, add2: instruction performing integer addition and starting carry chain
//   * adc1, adc2: instruction performing integer addition with carry
// Output: T[0-6] registers
#define MULS_128x320(I0, I1, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, add1, add2, adc1, adc2) \
	\ // Column 0
	MOVQ    I0, DX              \
	MULXQ   I1+24(SB), T0, T1   \
	MULXQ   I1+32(SB), T4, T2   \
	XORQ    AX, AX              \
	MULXQ   I1+40(SB), T5, T3   \
	add1    T4, T1              \
	adc1    T5, T2              \
	MULXQ   I1+48(SB), T7, T4   \
	adc1    T7, T3              \
	MULXQ   I1+56(SB), T6, T5   \
	adc1    T6, T4              \
	adc1    AX, T5              \
	\ // Column 1
	MOVQ    8+I0, DX            \
	MULXQ   I1+24(SB), T6, T7   \
	add2    T6, T1              \
	adc2    T7, T2              \
	MULXQ   I1+32(SB), T8, T6   \
	adc2    T6, T3              \
	MULXQ   I1+40(SB), T7, T9   \
	adc2    T9, T4              \
	MULXQ   I1+48(SB), T9, T6   \
	adc2    T6, T5              \
	MULXQ   I1+56(SB), DX, T6   \
	adc2    AX, T6              \
	\ // Output
	XORQ    AX, AX              \
	add1    T8, T2              \
	adc1    T7, T3              \
	adc1    T9, T4              \
	adc1    DX, T5              \
	adc1    AX, T6

// Multiplies 128-bit with 320-bit integer. Optimized with MULX instruction.
#define MULS_128x320_MULX(I0, I1, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) \
	MULS_128x320(I0, I1, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, ADDQ, ADDQ, ADCQ, ADCQ)

// Multiplies 128-bit with 320-bit integer. Optimized with  MULX, ADOX and ADCX instructions
#define MULS_128x320_MULX_ADCX_ADOX(I0, I1, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) \
	MULS_128x320(I0, I1, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, ADOXQ, ADCXQ, ADOXQ, ADCXQ)

// Template of a macro performing multiplication of two 512-bit numbers. It uses one
// level of Karatsuba and one level of schoolbook multiplication. Template must be
// customized with macro performing schoolbook multiplication.
// Input:
//  * I0, I1 - two 512-bit numbers
//  * MULS - either MULS256_MULX or MULS256_MULX_ADCX_ADOX
// Output: OUT - 1024-bit long
#define MUL(OUT, I0, I1, MULS) \
	\ // R[8-11]: U1+U0
	XORQ    AX, AX  \
	MOVQ    ( 0)(I0), R8    \
	MOVQ    ( 8)(I0), R9    \
	MOVQ    (16)(I0), R10   \
	MOVQ    (24)(I0), R11   \
	ADDQ    (32)(I0), R8    \
	ADCQ    (40)(I0), R9    \
	ADCQ    (48)(I0), R10   \
	ADCQ    (56)(I0), R11   \
	SBBQ    $0, AX          \ // store mask
	MOVQ    R8,  ( 0)(SP)   \
	MOVQ    R9,  ( 8)(SP)   \
	MOVQ    R10, (16)(SP)   \
	MOVQ    R11, (24)(SP)   \
	\
	\ // R[12-15]: V1+V0
	XORQ    BX, BX          \
	MOVQ    ( 0)(I1), R12   \
	MOVQ    ( 8)(I1), R13   \
	MOVQ    (16)(I1), R14   \
	MOVQ    (24)(I1), R15   \
	ADDQ    (32)(I1), R12   \
	ADCQ    (40)(I1), R13   \
	ADCQ    (48)(I1), R14   \
	ADCQ    (56)(I1), R15   \
	SBBQ    $0, BX          \ // store mask
	MOVQ    R12, (32)(SP)   \
	MOVQ    R13, (40)(SP)   \
	MOVQ    R14, (48)(SP)   \
	MOVQ    R15, (56)(SP)   \
	\ // Prepare mask for U0+U1 (U1+U0 mod 256^4 if U1+U0 sets carry flag, otherwise 0)
	ANDQ    AX, R12         \
	ANDQ    AX, R13         \
	ANDQ    AX, R14         \
	ANDQ    AX, R15         \
	\ // Prepare mask for V0+V1 (V1+V0 mod 256^4 if U1+U0 sets carry flag, otherwise 0)
	ANDQ    BX, R8          \
	ANDQ    BX, R9          \
	ANDQ    BX, R10         \
	ANDQ    BX, R11         \
	\ // res = masked(U0+U1) + masked(V0 + V1)
	ADDQ    R12, R8         \
	ADCQ    R13, R9         \
	ADCQ    R14, R10        \
	ADCQ    R15, R11        \
	\ // SP[64-96] <- res
	MOVQ     R8, (64)(SP)   \
	MOVQ     R9, (72)(SP)   \
	MOVQ    R10, (80)(SP)   \
	MOVQ    R11, (88)(SP)   \
	\ // BP will be used for schoolbook multiplication below
	MOVQ    BP, 96(SP)  \
	\ // (U1+U0)*(V1+V0)
	MULS((64)(OUT), 0(SP), 32(SP), R8, R9, R10, R11, R12, R13, R14, R15, BX, BP)    \
	\ // U0 x V0
	MULS(0(OUT), 0(I0), 0(I1), R8, R9, R10, R11, R12, R13, R14, R15, BX, BP)    \
	\ // U1 x V1
	MULS(0(SP), 32(I0), 32(I1), R8, R9, R10, R11, R12, R13, R14, R15, BX, BP)  \
	\ // Recover BP
	MOVQ    96(SP), BP  \
	\ // Final part of schoolbook multiplication; R[8-11] = (U0+U1) x (V0+V1)
	MOVQ    (64)(SP), R8    \
	MOVQ    (72)(SP), R9    \
	MOVQ    (80)(SP), R10   \
	MOVQ    (88)(SP), R11   \
	MOVQ    (96)(OUT), AX   \
	ADDQ    AX, R8          \
	MOVQ    (104)(OUT), AX  \
	ADCQ    AX, R9          \
	MOVQ    (112)(OUT), AX  \
	ADCQ    AX, R10         \
	MOVQ    (120)(OUT), AX  \
	ADCQ    AX, R11 \
	\ // R[12-15, 8-11] = (U0+U1) x (V0+V1) - U0xV0
	MOVQ    (64)(OUT), R12  \
	MOVQ    (72)(OUT), R13  \
	MOVQ    (80)(OUT), R14  \
	MOVQ    (88)(OUT), R15  \
	SUBQ    ( 0)(OUT), R12  \
	SBBQ    ( 8)(OUT), R13  \
	SBBQ    (16)(OUT), R14  \
	SBBQ    (24)(OUT), R15  \
	SBBQ    (32)(OUT), R8   \
	SBBQ    (40)(OUT), R9   \
	SBBQ    (48)(OUT), R10  \
	SBBQ    (56)(OUT), R11  \
	\ // r8-r15 <- (U0+U1) x (V0+V1) - U0xV0 - U1xV1
	SUBQ    ( 0)(SP), R12   \
	SBBQ    ( 8)(SP), R13   \
	SBBQ    (16)(SP), R14   \
	SBBQ    (24)(SP), R15   \
	SBBQ    (32)(SP), R8    \
	SBBQ    (40)(SP), R9    \
	SBBQ    (48)(SP), R10   \
	SBBQ    (56)(SP), R11   \
	\
	;                       ADDQ   (32)(OUT), R12; MOVQ    R12, ( 32)(OUT) \
	;                       ADCQ   (40)(OUT), R13; MOVQ    R13, ( 40)(OUT) \
	;                       ADCQ   (48)(OUT), R14; MOVQ    R14, ( 48)(OUT) \
	;                       ADCQ   (56)(OUT), R15; MOVQ    R15, ( 56)(OUT) \
	MOVQ    ( 0)(SP), AX;   ADCQ    AX,  R8;       MOVQ     R8, ( 64)(OUT) \
	MOVQ    ( 8)(SP), AX;   ADCQ    AX,  R9;       MOVQ     R9, ( 72)(OUT) \
	MOVQ    (16)(SP), AX;   ADCQ    AX, R10;       MOVQ    R10, ( 80)(OUT) \
	MOVQ    (24)(SP), AX;   ADCQ    AX, R11;       MOVQ    R11, ( 88)(OUT) \
	MOVQ    (32)(SP), R12;  ADCQ    $0, R12;       MOVQ    R12, ( 96)(OUT) \
	MOVQ    (40)(SP), R13;  ADCQ    $0, R13;       MOVQ    R13, (104)(OUT) \
	MOVQ    (48)(SP), R14;  ADCQ    $0, R14;       MOVQ    R14, (112)(OUT) \
	MOVQ    (56)(SP), R15;  ADCQ    $0, R15;       MOVQ    R15, (120)(OUT)

// Template for calculating the Montgomery reduction algorithm described in
// section 5.2.3 of https://eprint.iacr.org/2017/1015.pdf. Template must be
// customized with schoolbook multiplicaton for 128 x 320-bit number.
// This macro reuses memory of IN value and *changes* it. Smashes registers
// R[8-15], BX, CX
// Input:
//    * IN: 1024-bit number to be reduced
//    * MULS: either MULS_128x320_MULX or MULS_128x320_MULX_ADCX_ADOX
// Output: OUT 512-bit
#define REDC(OUT, IN, MULS) \
	MULS(0(IN), ·p503p1, R8, R9, R10, R11, R12, R13, R14, BX, CX, R15) \
	XORQ    R15, R15        \
	ADDQ    (24)(IN), R8    \
	ADCQ    (32)(IN), R9    \
	ADCQ    (40)(IN), R10   \
	ADCQ    (48)(IN), R11   \
	ADCQ    (56)(IN), R12   \
	ADCQ    (64)(IN), R13   \
	ADCQ    (72)(IN), R14   \
	ADCQ    (80)(IN), R15   \
	MOVQ    R8, (24)(IN)    \
	MOVQ    R9, (32)(IN)    \
	MOVQ    R10, (40)(IN)   \
	MOVQ    R11, (48)(IN)   \
	MOVQ    R12, (56)(IN)   \
	MOVQ    R13, (64)(IN)   \
	MOVQ    R14, (72)(IN)   \
	MOVQ    R15, (80)(IN)   \
	MOVQ    (88)(IN), R8    \
	MOVQ    (96)(IN), R9    \
	MOVQ    (104)(IN), R10  \
	MOVQ    (112)(IN), R11  \
	MOVQ    (120)(IN), R12  \
	ADCQ    $0, R8          \
	ADCQ    $0, R9          \
	ADCQ    $0, R10         \
	ADCQ    $0, R11         \
	ADCQ    $0, R12         \
	MOVQ    R8, (88)(IN)    \
	MOVQ    R9, (96)(IN)    \
	MOVQ    R10, (104)(IN)  \
	MOVQ    R11, (112)(IN)  \
	MOVQ    R12, (120)(IN)  \
	\
	MULS(16(IN), ·p503p1, R8, R9, R10, R11, R12, R13, R14, BX, CX, R15)    \
	XORQ    R15, R15        \
	ADDQ    (40)(IN), R8    \
	ADCQ    (48)(IN), R9    \
	ADCQ    (56)(IN), R10   \
	ADCQ    (64)(IN), R11   \
	ADCQ    (72)(IN), R12   \
	ADCQ    (80)(IN), R13   \
	ADCQ    (88)(IN), R14   \
	ADCQ    (96)(IN), R15   \
	MOVQ    R8, (40)(IN)    \
	MOVQ    R9, (48)(IN)    \
	MOVQ    R10, (56)(IN)   \
	MOVQ    R11, (64)(IN)   \
	MOVQ    R12, (72)(IN)   \
	MOVQ    R13, (80)(IN)   \
	MOVQ    R14, (88)(IN)   \
	MOVQ    R15, (96)(IN)   \
	MOVQ    (104)(IN), R8   \
	MOVQ    (112)(IN), R9   \
	MOVQ    (120)(IN), R10  \
	ADCQ    $0, R8          \
	ADCQ    $0, R9          \
	ADCQ    $0, R10         \
	MOVQ    R8, (104)(IN)   \
	MOVQ    R9, (112)(IN)   \
	MOVQ    R10, (120)(IN)  \
	\
	MULS(32(IN), ·p503p1, R8, R9, R10, R11, R12, R13, R14, BX, CX, R15)    \
	XORQ    R15, R15        \
	XORQ    BX, BX          \
	ADDQ    ( 56)(IN), R8   \
	ADCQ    ( 64)(IN), R9   \
	ADCQ    ( 72)(IN), R10  \
	ADCQ    ( 80)(IN), R11  \
	ADCQ    ( 88)(IN), R12  \
	ADCQ    ( 96)(IN), R13  \
	ADCQ    (104)(IN), R14  \
	ADCQ    (112)(IN), R15  \
	ADCQ    (120)(IN), BX   \
	MOVQ    R8,  ( 56)(IN)  \
	MOVQ    R10, ( 72)(IN)  \
	MOVQ    R11, ( 80)(IN)  \
	MOVQ    R12, ( 88)(IN)  \
	MOVQ    R13, ( 96)(IN)  \
	MOVQ    R14, (104)(IN)  \
	MOVQ    R15, (112)(IN)  \
	MOVQ    BX,  (120)(IN)  \
	MOVQ    R9,  (  0)(OUT) \ // Result: OUT[0]
	\
	MULS(48(IN), ·p503p1, R8, R9, R10, R11, R12, R13, R14, BX, CX, R15)    \
	ADDQ    ( 72)(IN), R8   \
	ADCQ    ( 80)(IN), R9   \
	ADCQ    ( 88)(IN), R10  \
	ADCQ    ( 96)(IN), R11  \
	ADCQ    (104)(IN), R12  \
	ADCQ    (112)(IN), R13  \
	ADCQ    (120)(IN), R14  \
	MOVQ    R8,  ( 8)(OUT)  \ // Result: OUT[1]
	MOVQ    R9,  (16)(OUT)  \ // Result: OUT[2]
	MOVQ    R10, (24)(OUT)  \ // Result: OUT[3]
	MOVQ    R11, (32)(OUT)  \ // Result: OUT[4]
	MOVQ    R12, (40)(OUT)  \ // Result: OUT[5]
	MOVQ    R13, (48)(OUT)  \ // Result: OUT[6] and OUT[7]
	MOVQ    R14, (56)(OUT)

TEXT ·fp503StrongReduce(SB), NOSPLIT, $0-8
	MOVQ	x+0(FP), REG_P1

	// Zero AX for later use:
	XORQ	AX, AX

	// Load p into registers:
	MOVQ	P503_0, R8
	// P503_{1,2} = P503_0, so reuse R8
	MOVQ	P503_3, R9
	MOVQ	P503_4, R10
	MOVQ	P503_5, R11
	MOVQ	P503_6, R12
	MOVQ	P503_7, R13

	// Set x <- x - p
	SUBQ	R8,  ( 0)(REG_P1)
	SBBQ	R8,  ( 8)(REG_P1)
	SBBQ	R8,  (16)(REG_P1)
	SBBQ	R9,  (24)(REG_P1)
	SBBQ	R10, (32)(REG_P1)
	SBBQ	R11, (40)(REG_P1)
	SBBQ	R12, (48)(REG_P1)
	SBBQ	R13, (56)(REG_P1)

	// Save carry flag indicating x-p < 0 as a mask
	SBBQ	$0, AX

	// Conditionally add p to x if x-p < 0
	ANDQ	AX, R8
	ANDQ	AX, R9
	ANDQ	AX, R10
	ANDQ	AX, R11
	ANDQ	AX, R12
	ANDQ	AX, R13

	ADDQ	R8, ( 0)(REG_P1)
	ADCQ	R8, ( 8)(REG_P1)
	ADCQ	R8, (16)(REG_P1)
	ADCQ	R9, (24)(REG_P1)
	ADCQ	R10,(32)(REG_P1)
	ADCQ	R11,(40)(REG_P1)
	ADCQ	R12,(48)(REG_P1)
	ADCQ	R13,(56)(REG_P1)

	RET

TEXT ·fp503ConditionalSwap(SB),NOSPLIT,$0-17

	MOVQ	x+0(FP), REG_P1
	MOVQ	y+8(FP), REG_P2
	MOVB	choice+16(FP), AL	// AL = 0 or 1
	MOVBLZX	AL, AX				// AX = 0 or 1
	NEGQ	AX					// AX = 0x00..00 or 0xff..ff

#ifndef CSWAP_BLOCK
#define CSWAP_BLOCK(idx) 	\
	MOVQ	(idx*8)(REG_P1), BX	\ // BX = x[idx]
	MOVQ 	(idx*8)(REG_P2), CX	\ // CX = y[idx]
	MOVQ	CX, DX				\ // DX = y[idx]
	XORQ	BX, DX				\ // DX = y[idx] ^ x[idx]
	ANDQ	AX, DX				\ // DX = (y[idx] ^ x[idx]) & mask
	XORQ	DX, BX				\ // BX = (y[idx] ^ x[idx]) & mask) ^ x[idx] = x[idx] or y[idx]
	XORQ	DX, CX				\ // CX = (y[idx] ^ x[idx]) & mask) ^ y[idx] = y[idx] or x[idx]
	MOVQ	BX, (idx*8)(REG_P1)	\
	MOVQ	CX, (idx*8)(REG_P2)
#endif

	CSWAP_BLOCK(0)
	CSWAP_BLOCK(1)
	CSWAP_BLOCK(2)
	CSWAP_BLOCK(3)
	CSWAP_BLOCK(4)
	CSWAP_BLOCK(5)
	CSWAP_BLOCK(6)
	CSWAP_BLOCK(7)

#ifdef CSWAP_BLOCK
#undef CSWAP_BLOCK
#endif

	RET

TEXT ·fp503AddReduced(SB),NOSPLIT,$0-24

	MOVQ	z+0(FP), REG_P3
	MOVQ	x+8(FP), REG_P1
	MOVQ	y+16(FP), REG_P2

	// Used later to calculate a mask
	XORQ    CX, CX

	// [R8-R15]: z = x + y
	MOVQ	( 0)(REG_P1), R8
	MOVQ	( 8)(REG_P1), R9
	MOVQ	(16)(REG_P1), R10
	MOVQ	(24)(REG_P1), R11
	MOVQ	(32)(REG_P1), R12
	MOVQ	(40)(REG_P1), R13
	MOVQ	(48)(REG_P1), R14
	MOVQ	(56)(REG_P1), R15
	ADDQ	( 0)(REG_P2), R8
	ADCQ	( 8)(REG_P2), R9
	ADCQ	(16)(REG_P2), R10
	ADCQ	(24)(REG_P2), R11
	ADCQ	(32)(REG_P2), R12
	ADCQ	(40)(REG_P2), R13
	ADCQ	(48)(REG_P2), R14
	ADCQ	(56)(REG_P2), R15

	MOVQ    P503X2_0, AX
	SUBQ    AX, R8
	MOVQ    P503X2_1, AX
	SBBQ    AX, R9
	SBBQ    AX, R10
	MOVQ    P503X2_3, AX
	SBBQ    AX, R11
	MOVQ    P503X2_4, AX
	SBBQ    AX, R12
	MOVQ    P503X2_5, AX
	SBBQ    AX, R13
	MOVQ    P503X2_6, AX
	SBBQ    AX, R14
	MOVQ    P503X2_7, AX
	SBBQ    AX, R15

	// mask
	SBBQ    $0, CX

	// move z to REG_P3
	MOVQ    R8,  ( 0)(REG_P3)
	MOVQ    R9,  ( 8)(REG_P3)
	MOVQ    R10, (16)(REG_P3)
	MOVQ    R11, (24)(REG_P3)
	MOVQ    R12, (32)(REG_P3)
	MOVQ    R13, (40)(REG_P3)
	MOVQ    R14, (48)(REG_P3)
	MOVQ    R15, (56)(REG_P3)

	// if z<0 add p503x2 back
	MOVQ    P503X2_0,   R8
	MOVQ    P503X2_1,   R9
	MOVQ    P503X2_3,   R10
	MOVQ    P503X2_4,   R11
	MOVQ    P503X2_5,   R12
	MOVQ    P503X2_6,   R13
	MOVQ    P503X2_7,   R14
	ANDQ    CX, R8
	ANDQ    CX, R9
	ANDQ    CX, R10
	ANDQ    CX, R11
	ANDQ    CX, R12
	ANDQ    CX, R13
	ANDQ    CX, R14
	MOVQ    ( 0)(REG_P3), AX; ADDQ    R8,  AX; MOVQ    AX, ( 0)(REG_P3)
	MOVQ    ( 8)(REG_P3), AX; ADCQ    R9,  AX; MOVQ    AX, ( 8)(REG_P3)
	MOVQ    (16)(REG_P3), AX; ADCQ    R9,  AX; MOVQ    AX, (16)(REG_P3)
	MOVQ    (24)(REG_P3), AX; ADCQ    R10, AX; MOVQ    AX, (24)(REG_P3)
	MOVQ    (32)(REG_P3), AX; ADCQ    R11, AX; MOVQ    AX, (32)(REG_P3)
	MOVQ    (40)(REG_P3), AX; ADCQ    R12, AX; MOVQ    AX, (40)(REG_P3)
	MOVQ    (48)(REG_P3), AX; ADCQ    R13, AX; MOVQ    AX, (48)(REG_P3)
	MOVQ    (56)(REG_P3), AX; ADCQ    R14, AX; MOVQ    AX, (56)(REG_P3)
	RET

TEXT ·fp503SubReduced(SB), NOSPLIT, $0-24

	MOVQ    z+0(FP), REG_P3
	MOVQ    x+8(FP), REG_P1
	MOVQ    y+16(FP), REG_P2

	// Used later to calculate a mask
	XORQ    CX, CX

	MOVQ    ( 0)(REG_P1), R8
	MOVQ    ( 8)(REG_P1), R9
	MOVQ    (16)(REG_P1), R10
	MOVQ    (24)(REG_P1), R11
	MOVQ    (32)(REG_P1), R12
	MOVQ    (40)(REG_P1), R13
	MOVQ    (48)(REG_P1), R14
	MOVQ    (56)(REG_P1), R15

	SUBQ    ( 0)(REG_P2), R8
	SBBQ    ( 8)(REG_P2), R9
	SBBQ    (16)(REG_P2), R10
	SBBQ    (24)(REG_P2), R11
	SBBQ    (32)(REG_P2), R12
	SBBQ    (40)(REG_P2), R13
	SBBQ    (48)(REG_P2), R14
	SBBQ    (56)(REG_P2), R15

	// mask
	SBBQ    $0, CX

	// store x-y in REG_P3
	MOVQ    R8,  ( 0)(REG_P3)
	MOVQ    R9,  ( 8)(REG_P3)
	MOVQ    R10, (16)(REG_P3)
	MOVQ    R11, (24)(REG_P3)
	MOVQ    R12, (32)(REG_P3)
	MOVQ    R13, (40)(REG_P3)
	MOVQ    R14, (48)(REG_P3)
	MOVQ    R15, (56)(REG_P3)

	// if z<0 add p503x2 back
	MOVQ    P503X2_0, R8
	MOVQ    P503X2_1, R9
	MOVQ    P503X2_3, R10
	MOVQ    P503X2_4, R11
	MOVQ    P503X2_5, R12
	MOVQ    P503X2_6, R13
	MOVQ    P503X2_7, R14
	ANDQ    CX, R8
	ANDQ    CX, R9
	ANDQ    CX, R10
	ANDQ    CX, R11
	ANDQ    CX, R12
	ANDQ    CX, R13
	ANDQ    CX, R14
	MOVQ    ( 0)(REG_P3), AX; ADDQ    R8,  AX; MOVQ    AX, ( 0)(REG_P3)
	MOVQ    ( 8)(REG_P3), AX; ADCQ    R9,  AX; MOVQ    AX, ( 8)(REG_P3)
	MOVQ    (16)(REG_P3), AX; ADCQ    R9,  AX; MOVQ    AX, (16)(REG_P3)
	MOVQ    (24)(REG_P3), AX; ADCQ    R10, AX; MOVQ    AX, (24)(REG_P3)
	MOVQ    (32)(REG_P3), AX; ADCQ    R11, AX; MOVQ    AX, (32)(REG_P3)
	MOVQ    (40)(REG_P3), AX; ADCQ    R12, AX; MOVQ    AX, (40)(REG_P3)
	MOVQ    (48)(REG_P3), AX; ADCQ    R13, AX; MOVQ    AX, (48)(REG_P3)
	MOVQ    (56)(REG_P3), AX; ADCQ    R14, AX; MOVQ    AX, (56)(REG_P3)

	RET

TEXT ·fp503Mul(SB), NOSPLIT, $104-24
	MOVQ    z+ 0(FP), CX
	MOVQ    x+ 8(FP), REG_P1
	MOVQ    y+16(FP), REG_P2

	// Check wether to use optimized implementation
	CMPB    ·HasADXandBMI2(SB), $1
	JE      mul_with_mulx_adcx_adox
	CMPB    ·HasBMI2(SB), $1
	JE      mul_with_mulx

	// Generic x86 implementation (below) uses variant of Karatsuba method.
	//
	// Here we store the destination in CX instead of in REG_P3 because the
	// multiplication instructions use DX as an implicit destination
	// operand: MULQ $REG sets DX:AX <-- AX * $REG.

	// RAX and RDX will be used for a mask (0-borrow)
	XORQ	AX, AX

	// RCX[0-3]: U1+U0
	MOVQ	(32)(REG_P1), R8
	MOVQ	(40)(REG_P1), R9
	MOVQ	(48)(REG_P1), R10
	MOVQ	(56)(REG_P1), R11
	ADDQ	( 0)(REG_P1), R8
	ADCQ	( 8)(REG_P1), R9
	ADCQ	(16)(REG_P1), R10
	ADCQ	(24)(REG_P1), R11
	MOVQ	R8,  ( 0)(CX)
	MOVQ	R9,  ( 8)(CX)
	MOVQ	R10, (16)(CX)
	MOVQ	R11, (24)(CX)

	SBBQ	$0, AX

	// R12-R15: V1+V0
	XORQ	DX, DX
	MOVQ	(32)(REG_P2), R12
	MOVQ	(40)(REG_P2), R13
	MOVQ	(48)(REG_P2), R14
	MOVQ	(56)(REG_P2), R15
	ADDQ	( 0)(REG_P2), R12
	ADCQ	( 8)(REG_P2), R13
	ADCQ	(16)(REG_P2), R14
	ADCQ	(24)(REG_P2), R15

	SBBQ	$0, DX

	// Store carries on stack
	MOVQ	AX, (64)(SP)
	MOVQ	DX, (72)(SP)

	// (SP[0-3],R8,R9,R10,R11) <- (U0+U1)*(V0+V1).
	// MUL using comba; In comments below U=U0+U1 V=V0+V1

	// U0*V0
	MOVQ    (CX), AX
	MULQ    R12
	MOVQ    AX, (SP)        // C0
	MOVQ    DX, R8

	// U0*V1
	XORQ    R9, R9
	MOVQ    (CX), AX
	MULQ    R13
	ADDQ    AX, R8
	ADCQ    DX, R9

	// U1*V0
	XORQ    R10, R10
	MOVQ    (8)(CX), AX
	MULQ    R12
	ADDQ    AX, R8
	MOVQ    R8, (8)(SP)     // C1
	ADCQ    DX, R9
	ADCQ    $0, R10

	// U0*V2
	XORQ    R8, R8
	MOVQ    (CX), AX
	MULQ    R14
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	// U2*V0
	MOVQ    (16)(CX), AX
	MULQ    R12
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	// U1*V1
	MOVQ    (8)(CX), AX
	MULQ    R13
	ADDQ    AX, R9
	MOVQ    R9, (16)(SP)        // C2
	ADCQ    DX, R10
	ADCQ    $0, R8

	// U0*V3
	XORQ    R9, R9
	MOVQ    (CX), AX
	MULQ    R15
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	// U3*V0
	MOVQ    (24)(CX), AX
	MULQ    R12
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	// U1*V2
	MOVQ    (8)(CX), AX
	MULQ    R14
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	// U2*V1
	MOVQ    (16)(CX), AX
	MULQ    R13
	ADDQ    AX, R10
	MOVQ    R10, (24)(SP)       // C3
	ADCQ    DX, R8
	ADCQ    $0, R9

	// U1*V3
	XORQ    R10, R10
	MOVQ    (8)(CX), AX
	MULQ    R15
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	// U3*V1
	MOVQ    (24)(CX), AX
	MULQ    R13
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	// U2*V2
	MOVQ    (16)(CX), AX
	MULQ    R14
	ADDQ    AX, R8
	MOVQ    R8, (32)(SP)        // C4
	ADCQ    DX, R9
	ADCQ    $0, R10

	// U2*V3
	XORQ    R11, R11
	MOVQ    (16)(CX), AX
	MULQ    R15
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R11

	// U3*V2
	MOVQ    (24)(CX), AX
	MULQ    R14
	ADDQ    AX, R9              // C5
	ADCQ    DX, R10
	ADCQ    $0, R11

	// U3*V3
	MOVQ    (24)(CX), AX
	MULQ    R15
	ADDQ    AX, R10             // C6
	ADCQ    DX, R11             // C7

	MOVQ    (64)(SP), AX
	ANDQ    AX, R12
	ANDQ    AX, R13
	ANDQ    AX, R14
	ANDQ    AX, R15
	ADDQ    R8, R12
	ADCQ    R9, R13
	ADCQ    R10, R14
	ADCQ    R11, R15

	MOVQ    (72)(SP), AX
	MOVQ    (CX), R8
	MOVQ    (8)(CX), R9
	MOVQ    (16)(CX), R10
	MOVQ    (24)(CX), R11
	ANDQ    AX, R8
	ANDQ    AX, R9
	ANDQ    AX, R10
	ANDQ    AX, R11
	ADDQ    R12, R8
	ADCQ    R13, R9
	ADCQ    R14, R10
	ADCQ    R15, R11
	MOVQ    R8, (32)(SP)
	MOVQ    R9, (40)(SP)
	MOVQ    R10, (48)(SP)
	MOVQ    R11, (56)(SP)

	// CX[0-7] <- AL*BL

	// U0*V0
	MOVQ    (REG_P1), R11
	MOVQ    (REG_P2), AX
	MULQ    R11
	XORQ    R9, R9
	MOVQ    AX, (CX)            // C0
	MOVQ    DX, R8

	// U0*V1
	MOVQ    (16)(REG_P1), R14
	MOVQ    (8)(REG_P2), AX
	MULQ    R11
	XORQ    R10, R10
	ADDQ    AX, R8
	ADCQ    DX, R9

	// U1*V0
	MOVQ    (8)(REG_P1), R12
	MOVQ    (REG_P2), AX
	MULQ    R12
	ADDQ    AX, R8
	MOVQ    R8, (8)(CX)         // C1
	ADCQ    DX, R9
	ADCQ    $0, R10

	// U0*V2
	XORQ    R8, R8
	MOVQ    (16)(REG_P2), AX
	MULQ    R11
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	// U2*V0
	MOVQ    (REG_P2), R13
	MOVQ    R14, AX
	MULQ    R13
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	// U1*V1
	MOVQ    (8)(REG_P2), AX
	MULQ    R12
	ADDQ    AX, R9
	MOVQ    R9, (16)(CX)        // C2
	ADCQ    DX, R10
	ADCQ    $0, R8

	// U0*V3
	XORQ    R9, R9
	MOVQ    (24)(REG_P2), AX
	MULQ    R11
	MOVQ    (24)(REG_P1), R15
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	// U3*V1
	MOVQ    R15, AX
	MULQ    R13
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	// U2*V2
	MOVQ    (16)(REG_P2), AX
	MULQ    R12
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	// U2*V3
	MOVQ    (8)(REG_P2), AX
	MULQ    R14
	ADDQ    AX, R10
	MOVQ    R10, (24)(CX)       // C3
	ADCQ    DX, R8
	ADCQ    $0, R9

	// U3*V2
	XORQ    R10, R10
	MOVQ    (24)(REG_P2), AX
	MULQ    R12
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	// U3*V1
	MOVQ    (8)(REG_P2), AX
	MULQ    R15
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	// U2*V2
	MOVQ    (16)(REG_P2), AX
	MULQ    R14
	ADDQ    AX, R8
	MOVQ    R8, (32)(CX)		// C4
	ADCQ    DX, R9
	ADCQ    $0, R10

	// U2*V3
	XORQ    R8, R8
	MOVQ    (24)(REG_P2), AX
	MULQ    R14
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	// U3*V2
	MOVQ    (16)(REG_P2), AX
	MULQ    R15
	ADDQ    AX, R9
	MOVQ    R9, (40)(CX)		// C5
	ADCQ    DX, R10
	ADCQ    $0, R8

	// U3*V3
	MOVQ    (24)(REG_P2), AX
	MULQ    R15
	ADDQ    AX, R10
	MOVQ    R10, (48)(CX)		// C6
	ADCQ    DX, R8
	MOVQ    R8, (56)(CX)		// C7

	// CX[8-15] <- U1*V1
	MOVQ    (32)(REG_P1), R11
	MOVQ    (32)(REG_P2), AX
	MULQ    R11
	XORQ    R9, R9
	MOVQ    AX, (64)(CX)        // C0
	MOVQ    DX, R8

	MOVQ    (48)(REG_P1), R14
	MOVQ    (40)(REG_P2), AX
	MULQ    R11
	XORQ    R10, R10
	ADDQ    AX, R8
	ADCQ    DX, R9

	MOVQ    (40)(REG_P1), R12
	MOVQ    (32)(REG_P2), AX
	MULQ    R12
	ADDQ    AX, R8
	MOVQ    R8, (72)(CX)        // C1
	ADCQ    DX, R9
	ADCQ    $0, R10

	XORQ    R8, R8
	MOVQ    (48)(REG_P2), AX
	MULQ    R11
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	MOVQ    (32)(REG_P2), R13
	MOVQ    R14, AX
	MULQ    R13
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	MOVQ    (40)(REG_P2), AX
	MULQ    R12
	ADDQ    AX, R9
	MOVQ    R9, (80)(CX)        // C2
	ADCQ    DX, R10
	ADCQ    $0, R8

	XORQ    R9, R9
	MOVQ    (56)(REG_P2), AX
	MULQ    R11
	MOVQ    (56)(REG_P1), R15
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	MOVQ    R15, AX
	MULQ    R13
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	MOVQ    (48)(REG_P2), AX
	MULQ    R12
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	MOVQ    (40)(REG_P2), AX
	MULQ    R14
	ADDQ    AX, R10
	MOVQ    R10, (88)(CX)       // C3
	ADCQ    DX, R8
	ADCQ    $0, R9

	XORQ    R10, R10
	MOVQ    (56)(REG_P2), AX
	MULQ    R12
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	MOVQ    (40)(REG_P2), AX
	MULQ    R15
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	MOVQ    (48)(REG_P2), AX
	MULQ    R14
	ADDQ    AX, R8
	MOVQ    R8, (96)(CX)        // C4
	ADCQ    DX, R9
	ADCQ    $0, R10

	XORQ    R8, R8
	MOVQ    (56)(REG_P2), AX
	MULQ    R14
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	MOVQ    (48)(REG_P2), AX
	MULQ    R15
	ADDQ    AX, R9
	MOVQ    R9, (104)(CX)       // C5
	ADCQ    DX, R10
	ADCQ    $0, R8

	MOVQ    (56)(REG_P2), AX
	MULQ    R15
	ADDQ    AX, R10
	MOVQ    R10, (112)(CX)      // C6
	ADCQ    DX, R8
	MOVQ    R8, (120)(CX)       // C7

	// [R8-R15] <- (U0+U1)*(V0+V1) - U1*V1
	MOVQ    (SP), R8
	SUBQ    (CX), R8
	MOVQ    (8)(SP), R9
	SBBQ    (8)(CX), R9
	MOVQ    (16)(SP), R10
	SBBQ    (16)(CX), R10
	MOVQ    (24)(SP), R11
	SBBQ    (24)(CX), R11
	MOVQ    (32)(SP), R12
	SBBQ    (32)(CX), R12
	MOVQ    (40)(SP), R13
	SBBQ    (40)(CX), R13
	MOVQ    (48)(SP), R14
	SBBQ    (48)(CX), R14
	MOVQ    (56)(SP), R15
	SBBQ    (56)(CX), R15

	// [R8-R15] <- (U0+U1)*(V0+V1) - U1*V0 - U0*U1
	MOVQ    ( 64)(CX), AX;	SUBQ    AX, R8
	MOVQ    ( 72)(CX), AX;	SBBQ    AX, R9
	MOVQ    ( 80)(CX), AX;	SBBQ    AX, R10
	MOVQ    ( 88)(CX), AX;	SBBQ    AX, R11
	MOVQ    ( 96)(CX), AX;	SBBQ    AX, R12
	MOVQ    (104)(CX), DX;	SBBQ    DX, R13
	MOVQ    (112)(CX), DI;	SBBQ    DI, R14
	MOVQ    (120)(CX), SI;	SBBQ    SI, R15

	// Final result
	ADDQ    (32)(CX), R8;	MOVQ    R8,  (32)(CX)
	ADCQ    (40)(CX), R9;	MOVQ    R9,  (40)(CX)
	ADCQ    (48)(CX), R10;	MOVQ    R10, (48)(CX)
	ADCQ    (56)(CX), R11;	MOVQ    R11, (56)(CX)
	ADCQ    (64)(CX), R12;	MOVQ    R12, (64)(CX)
	ADCQ    (72)(CX), R13;	MOVQ    R13, (72)(CX)
	ADCQ    (80)(CX), R14;	MOVQ    R14, (80)(CX)
	ADCQ    (88)(CX), R15;	MOVQ    R15, (88)(CX)
	ADCQ    $0, AX;        	MOVQ    AX,  (96)(CX)
	ADCQ    $0, DX;        	MOVQ    DX, (104)(CX)
	ADCQ    $0, DI;         MOVQ    DI, (112)(CX)
	ADCQ    $0, SI;     	MOVQ    SI, (120)(CX)
	RET

mul_with_mulx_adcx_adox:
	// Mul implementation for CPUs supporting two independent carry chain
	// (ADOX/ADCX) instructions and carry-less MULX multiplier
	MUL(CX, REG_P1, REG_P2, MULS256_MULX_ADCX_ADOX)
	RET

mul_with_mulx:
	// Mul implementation for CPUs supporting carry-less MULX multiplier.
	MUL(CX, REG_P1, REG_P2, MULS256_MULX)
	RET

TEXT ·fp503MontgomeryReduce(SB), $0-16
	MOVQ    z+0(FP), REG_P2
	MOVQ    x+8(FP), REG_P1

	// Check wether to use optimized implementation
	CMPB    ·HasADXandBMI2(SB), $1
	JE      redc_with_mulx_adcx_adox
	CMPB    ·HasBMI2(SB), $1
	JE      redc_with_mulx

	MOVQ    (REG_P1), R11
	MOVQ    P503P1_3, AX
	MULQ    R11
	XORQ    R8, R8
	ADDQ    (24)(REG_P1), AX
	MOVQ    AX, (24)(REG_P2)
	ADCQ    DX, R8

	XORQ    R9, R9
	MOVQ    P503P1_4, AX
	MULQ    R11
	XORQ    R10, R10
	ADDQ    AX, R8
	ADCQ    DX, R9

	MOVQ    (8)(REG_P1), R12
	MOVQ    P503P1_3, AX
	MULQ    R12
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10
	ADDQ    (32)(REG_P1), R8
	MOVQ    R8, (32)(REG_P2)       // Z4
	ADCQ    $0, R9
	ADCQ    $0, R10

	XORQ    R8, R8
	MOVQ    P503P1_5, AX
	MULQ    R11
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	MOVQ    P503P1_4, AX
	MULQ    R12
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	MOVQ    (16)(REG_P1), R13
	MOVQ    P503P1_3, AX
	MULQ    R13
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8
	ADDQ    (40)(REG_P1), R9
	MOVQ    R9, (40)(REG_P2)       // Z5
	ADCQ    $0, R10
	ADCQ    $0, R8

	XORQ    R9, R9
	MOVQ    P503P1_6, AX
	MULQ    R11
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	MOVQ    P503P1_5, AX
	MULQ    R12
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	MOVQ    P503P1_4, AX
	MULQ    R13
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	MOVQ    (24)(REG_P2), R14
	MOVQ    P503P1_3, AX
	MULQ    R14
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9
	ADDQ    (48)(REG_P1), R10
	MOVQ    R10, (48)(REG_P2)      // Z6
	ADCQ    $0, R8
	ADCQ    $0, R9

	XORQ    R10, R10
	MOVQ    P503P1_7, AX
	MULQ    R11
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	MOVQ    P503P1_6, AX
	MULQ    R12
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	MOVQ    P503P1_5, AX
	MULQ    R13
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	MOVQ    P503P1_4, AX
	MULQ    R14
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	MOVQ    (32)(REG_P2), R15
	MOVQ    P503P1_3, AX
	MULQ    R15
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10
	ADDQ    (56)(REG_P1), R8
	MOVQ    R8, (56)(REG_P2)       // Z7
	ADCQ    $0, R9
	ADCQ    $0, R10

	XORQ    R8, R8
	MOVQ    P503P1_7, AX
	MULQ    R12
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	MOVQ    P503P1_6, AX
	MULQ    R13
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	MOVQ    P503P1_5, AX
	MULQ    R14
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	MOVQ    P503P1_4, AX
	MULQ    R15
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	MOVQ    (40)(REG_P2), CX
	MOVQ    P503P1_3, AX
	MULQ    CX
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8
	ADDQ    (64)(REG_P1), R9
	MOVQ    R9, (REG_P2)           // Z0
	ADCQ    $0, R10
	ADCQ    $0, R8

	XORQ    R9, R9
	MOVQ    P503P1_7, AX
	MULQ    R13
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	MOVQ    P503P1_6, AX
	MULQ    R14
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	MOVQ    P503P1_5, AX
	MULQ    R15
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	MOVQ    P503P1_4, AX
	MULQ    CX
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	MOVQ    (48)(REG_P2), R13
	MOVQ    P503P1_3, AX
	MULQ    R13
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9
	ADDQ    (72)(REG_P1), R10
	MOVQ    R10, (8)(REG_P2)       // Z1
	ADCQ    $0, R8
	ADCQ    $0, R9

	XORQ    R10, R10
	MOVQ    P503P1_7, AX
	MULQ    R14
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	MOVQ    P503P1_6, AX
	MULQ    R15
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	MOVQ    P503P1_5, AX
	MULQ    CX
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	MOVQ    P503P1_4, AX
	MULQ    R13
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	MOVQ    (56)(REG_P2), R14
	MOVQ    P503P1_3, AX
	MULQ    R14
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10
	ADDQ    (80)(REG_P1), R8
	MOVQ    R8, (16)(REG_P2)       // Z2
	ADCQ    $0, R9
	ADCQ    $0, R10

	XORQ    R8, R8
	MOVQ    P503P1_7, AX
	MULQ    R15
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	MOVQ    P503P1_6, AX
	MULQ    CX
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	MOVQ    P503P1_5, AX
	MULQ    R13
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	MOVQ    P503P1_4, AX
	MULQ    R14
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8
	ADDQ    (88)(REG_P1), R9
	MOVQ    R9, (24)(REG_P2)       // Z3
	ADCQ    $0, R10
	ADCQ    $0, R8

	XORQ    R9, R9
	MOVQ    P503P1_7, AX
	MULQ    CX
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	MOVQ    P503P1_6, AX
	MULQ    R13
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	MOVQ    P503P1_5, AX
	MULQ    R14
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9
	ADDQ    (96)(REG_P1), R10
	MOVQ    R10, (32)(REG_P2)      // Z4
	ADCQ    $0, R8
	ADCQ    $0, R9

	XORQ    R10, R10
	MOVQ    P503P1_7, AX
	MULQ    R13
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	MOVQ    P503P1_6, AX
	MULQ    R14
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10
	ADDQ    (104)(REG_P1), R8      // Z5
	MOVQ    R8, (40)(REG_P2)       // Z5
	ADCQ    $0, R9
	ADCQ    $0, R10

	MOVQ    P503P1_7, AX
	MULQ    R14
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADDQ    (112)(REG_P1), R9      // Z6
	MOVQ    R9, (48)(REG_P2)       // Z6
	ADCQ    $0, R10
	ADDQ    (120)(REG_P1), R10     // Z7
	MOVQ    R10, (56)(REG_P2)      // Z7
	RET

redc_with_mulx_adcx_adox:
	// Implementation of the Montgomery reduction for CPUs
	// supporting two independent carry chain (ADOX/ADCX)
	// instructions and carry-less MULX multiplier
	REDC(REG_P2, REG_P1, MULS_128x320_MULX_ADCX_ADOX)
	RET

redc_with_mulx:
	// Implementation of the Montgomery reduction for CPUs
	// supporting carry-less MULX multiplier.
	REDC(REG_P2, REG_P1, MULS_128x320_MULX)
	RET

TEXT ·fp503AddLazy(SB), NOSPLIT, $0-24

	MOVQ z+0(FP), REG_P3
	MOVQ x+8(FP), REG_P1
	MOVQ y+16(FP), REG_P2

	MOVQ	(REG_P1), R8
	MOVQ	(8)(REG_P1), R9
	MOVQ	(16)(REG_P1), R10
	MOVQ	(24)(REG_P1), R11
	MOVQ	(32)(REG_P1), R12
	MOVQ	(40)(REG_P1), R13
	MOVQ	(48)(REG_P1), R14
	MOVQ	(56)(REG_P1), R15

	ADDQ	(REG_P2), R8
	ADCQ	(8)(REG_P2), R9
	ADCQ	(16)(REG_P2), R10
	ADCQ	(24)(REG_P2), R11
	ADCQ	(32)(REG_P2), R12
	ADCQ	(40)(REG_P2), R13
	ADCQ	(48)(REG_P2), R14
	ADCQ	(56)(REG_P2), R15

	MOVQ	R8, (REG_P3)
	MOVQ	R9, (8)(REG_P3)
	MOVQ	R10, (16)(REG_P3)
	MOVQ	R11, (24)(REG_P3)
	MOVQ	R12, (32)(REG_P3)
	MOVQ	R13, (40)(REG_P3)
	MOVQ	R14, (48)(REG_P3)
	MOVQ	R15, (56)(REG_P3)

	RET

TEXT ·fp503X2AddLazy(SB), NOSPLIT, $0-24

	MOVQ	z+0(FP), REG_P3
	MOVQ	x+8(FP), REG_P1
	MOVQ	y+16(FP), REG_P2

	MOVQ	(REG_P1), R8
	MOVQ	(8)(REG_P1), R9
	MOVQ	(16)(REG_P1), R10
	MOVQ	(24)(REG_P1), R11
	MOVQ	(32)(REG_P1), R12
	MOVQ	(40)(REG_P1), R13
	MOVQ	(48)(REG_P1), R14
	MOVQ	(56)(REG_P1), R15
	MOVQ	(64)(REG_P1), AX
	MOVQ	(72)(REG_P1), BX
	MOVQ	(80)(REG_P1), CX

	ADDQ	(REG_P2), R8
	ADCQ	(8)(REG_P2), R9
	ADCQ	(16)(REG_P2), R10
	ADCQ	(24)(REG_P2), R11
	ADCQ	(32)(REG_P2), R12
	ADCQ	(40)(REG_P2), R13
	ADCQ	(48)(REG_P2), R14
	ADCQ	(56)(REG_P2), R15
	ADCQ	(64)(REG_P2), AX
	ADCQ	(72)(REG_P2), BX
	ADCQ	(80)(REG_P2), CX

	MOVQ	R8, (REG_P3)
	MOVQ	R9, (8)(REG_P3)
	MOVQ	R10, (16)(REG_P3)
	MOVQ	R11, (24)(REG_P3)
	MOVQ	R12, (32)(REG_P3)
	MOVQ	R13, (40)(REG_P3)
	MOVQ	R14, (48)(REG_P3)
	MOVQ	R15, (56)(REG_P3)
	MOVQ	AX, (64)(REG_P3)
	MOVQ	BX, (72)(REG_P3)
	MOVQ	CX, (80)(REG_P3)

	MOVQ	(88)(REG_P1), R8
	MOVQ	(96)(REG_P1), R9
	MOVQ	(104)(REG_P1), R10
	MOVQ	(112)(REG_P1), R11
	MOVQ	(120)(REG_P1), R12

	ADCQ	(88)(REG_P2), R8
	ADCQ	(96)(REG_P2), R9
	ADCQ	(104)(REG_P2), R10
	ADCQ	(112)(REG_P2), R11
	ADCQ	(120)(REG_P2), R12

	MOVQ	R8, (88)(REG_P3)
	MOVQ	R9, (96)(REG_P3)
	MOVQ	R10, (104)(REG_P3)
	MOVQ	R11, (112)(REG_P3)
	MOVQ	R12, (120)(REG_P3)

	RET

TEXT ·fp503X2SubLazy(SB), NOSPLIT, $0-24

	MOVQ z+0(FP), REG_P3
	MOVQ x+8(FP), REG_P1
	MOVQ y+16(FP), REG_P2
	// Used later to store result of 0-borrow
	XORQ CX, CX

	// SUBC for first 11 limbs
	MOVQ	(REG_P1), R8
	MOVQ	(8)(REG_P1), R9
	MOVQ	(16)(REG_P1), R10
	MOVQ	(24)(REG_P1), R11
	MOVQ	(32)(REG_P1), R12
	MOVQ	(40)(REG_P1), R13
	MOVQ	(48)(REG_P1), R14
	MOVQ	(56)(REG_P1), R15
	MOVQ	(64)(REG_P1), AX
	MOVQ	(72)(REG_P1), BX

	SUBQ	(REG_P2), R8
	SBBQ	(8)(REG_P2), R9
	SBBQ	(16)(REG_P2), R10
	SBBQ	(24)(REG_P2), R11
	SBBQ	(32)(REG_P2), R12
	SBBQ	(40)(REG_P2), R13
	SBBQ	(48)(REG_P2), R14
	SBBQ	(56)(REG_P2), R15
	SBBQ	(64)(REG_P2), AX
	SBBQ	(72)(REG_P2), BX

	MOVQ	R8, (REG_P3)
	MOVQ	R9, (8)(REG_P3)
	MOVQ	R10, (16)(REG_P3)
	MOVQ	R11, (24)(REG_P3)
	MOVQ	R12, (32)(REG_P3)
	MOVQ	R13, (40)(REG_P3)
	MOVQ	R14, (48)(REG_P3)
	MOVQ	R15, (56)(REG_P3)
	MOVQ	AX, (64)(REG_P3)
	MOVQ	BX, (72)(REG_P3)

	// SUBC for last 5 limbs
	MOVQ	(80)(REG_P1), 	R8
	MOVQ	(88)(REG_P1), 	R9
	MOVQ	(96)(REG_P1), 	R10
	MOVQ	(104)(REG_P1), 	R11
	MOVQ	(112)(REG_P1), 	R12
	MOVQ	(120)(REG_P1), 	R13

	SBBQ	(80)(REG_P2), R8
	SBBQ	(88)(REG_P2), R9
	SBBQ	(96)(REG_P2), R10
	SBBQ	(104)(REG_P2), R11
	SBBQ	(112)(REG_P2), R12
	SBBQ	(120)(REG_P2), R13

	MOVQ	R8, (80)(REG_P3)
	MOVQ	R9, (88)(REG_P3)
	MOVQ	R10, (96)(REG_P3)
	MOVQ	R11, (104)(REG_P3)
	MOVQ	R12, (112)(REG_P3)
	MOVQ	R13, (120)(REG_P3)

	// Now the carry flag is 1 if x-y < 0.  If so, add p*2^512.
	SBBQ	$0, CX

	// Load p into registers:
	MOVQ	P503_0, R8
	// P503_{1,2} = P503_0, so reuse R8
	MOVQ	P503_3, R9
	MOVQ	P503_4, R10
	MOVQ	P503_5, R11
	MOVQ	P503_6, R12
	MOVQ	P503_7, R13

	ANDQ	CX, R8
	ANDQ	CX, R9
	ANDQ	CX, R10
	ANDQ	CX, R11
	ANDQ	CX, R12
	ANDQ	CX, R13

	MOVQ   (64   )(REG_P3), AX; ADDQ R8,  AX; MOVQ AX, (64   )(REG_P3)
	MOVQ   (64+ 8)(REG_P3), AX; ADCQ R8,  AX; MOVQ AX, (64+ 8)(REG_P3)
	MOVQ   (64+16)(REG_P3), AX; ADCQ R8,  AX; MOVQ AX, (64+16)(REG_P3)
	MOVQ   (64+24)(REG_P3), AX; ADCQ R9,  AX; MOVQ AX, (64+24)(REG_P3)
	MOVQ   (64+32)(REG_P3), AX; ADCQ R10, AX; MOVQ AX, (64+32)(REG_P3)
	MOVQ   (64+40)(REG_P3), AX; ADCQ R11, AX; MOVQ AX, (64+40)(REG_P3)
	MOVQ   (64+48)(REG_P3), AX; ADCQ R12, AX; MOVQ AX, (64+48)(REG_P3)
	MOVQ   (64+56)(REG_P3), AX; ADCQ R13, AX; MOVQ AX, (64+56)(REG_P3)

	RET
