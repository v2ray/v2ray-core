// +build amd64,!noasm

#include "textflag.h"

TEXT ·cpuid(SB), NOSPLIT, $0-4
    MOVL eaxArg+0(FP), AX
    MOVL ecxArg+4(FP), CX
    CPUID
    MOVL AX, eax+8(FP)
    MOVL BX, ebx+12(FP)
    MOVL CX, ecx+16(FP)
    MOVL DX, edx+20(FP)
    RET
