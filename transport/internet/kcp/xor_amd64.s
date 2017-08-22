#include "textflag.h"

// func xorfwd(x []byte)
TEXT ·xorfwd(SB),NOSPLIT,$0
  MOVQ x+0(FP), SI  // x[i]
  MOVL x_len+8(FP), CX  // x.len
  MOVQ x+0(FP), DI
  ADDQ $4, DI       // x[i+4]
  SUBQ $4, CX
xorfwdloop:
  MOVL (SI), AX
  XORL AX, (DI)
  ADDQ $4, SI
  ADDQ $4, DI
  SUBQ $4, CX

  CMPL CX, $0
  JE xorfwddone

  JMP xorfwdloop
xorfwddone:        
  RET

// func xorbkd(x []byte)
TEXT ·xorbkd(SB),NOSPLIT,$0
  MOVQ x+0(FP), SI
  MOVL x_len+8(FP), CX  // x.len
  MOVQ x+0(FP), DI
  ADDQ CX, SI       // x[-8]
  SUBQ $8, SI
  ADDQ CX, DI       // x[-4]
  SUBQ $4, DI
  SUBQ $4, CX
xorbkdloop:
  MOVL (SI), AX
  XORL AX, (DI)
  SUBQ $4, SI
  SUBQ $4, DI
  SUBQ $4, CX

  CMPL CX, $0
  JE xorbkddone
  
  JMP xorbkdloop

xorbkddone:        
  RET
