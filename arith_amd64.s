//go:build amd64 && !purego

#include "textflag.h"

// func andAVX2(z, x, y *[32]byte)
TEXT ·andAVX2(SB), NOSPLIT, $0-24
    MOVQ z+0(FP), AX
    MOVQ x+8(FP), BX
    MOVQ y+16(FP), CX

    VMOVDQU (BX), Y0
    VMOVDQU (CX), Y1

    VPAND Y0, Y1, Y2    // Y2 = Y0 & Y1

    VMOVDQU Y2, (AX)

    VZEROUPPER
    RET

// func orAVX2(z, x, y *[32]byte)
TEXT ·orAVX2(SB), NOSPLIT, $0-24
    MOVQ z+0(FP), AX
    MOVQ x+8(FP), BX
    MOVQ y+16(FP), CX

    VMOVDQU (BX), Y0
    VMOVDQU (CX), Y1
    VPOR Y0, Y1, Y2
    VMOVDQU Y2, (AX)
    VZEROUPPER
    RET

// func xorAVX2(z, x, y *[32]byte)
TEXT ·xorAVX2(SB), NOSPLIT, $0-24
    MOVQ z+0(FP), AX
    MOVQ x+8(FP), BX
    MOVQ y+16(FP), CX

    VMOVDQU (BX), Y0
    VMOVDQU (CX), Y1
    VPXOR Y0, Y1, Y2
    VMOVDQU Y2, (AX)
    VZEROUPPER
    RET

// func notAVX2(z, x *[32]byte)
TEXT ·notAVX2(SB), NOSPLIT, $0-16
    MOVQ z+0(FP), AX
    MOVQ x+8(FP), BX

    VPCMPEQD Y0, Y0, Y0  // Y0 = 0xFFFFFFFF...

    VMOVDQU (BX), Y1
    VPXOR Y0, Y1, Y2     // Y2 = Y0 XOR Y1

    VMOVDQU Y2, (AX)
    VZEROUPPER
    RET

// func addAVX2(z, x, y *[32]byte) uint64
TEXT ·addAVX2(SB), NOSPLIT, $0-32
    MOVQ z+0(FP), AX
    MOVQ x+8(FP), BX
    MOVQ y+16(FP), CX

    MOVQ 0(BX), R8      // x[0]
    MOVQ 8(BX), R9      // x[1]
    MOVQ 16(BX), R10    // x[2]
    MOVQ 24(BX), R11    // x[3]

    ADDQ 0(CX), R8      // x[0] + y[0]
    ADCQ 8(CX), R9      // x[1] + y[1] + carry
    ADCQ 16(CX), R10    // x[2] + y[2] + carry
    ADCQ 24(CX), R11    // x[3] + y[3] + carry

    MOVQ R8, 0(AX)
    MOVQ R9, 8(AX)
    MOVQ R10, 16(AX)
    MOVQ R11, 24(AX)

    MOVQ $0, R12
    ADCQ $0, R12
    MOVQ R12, ret+24(FP)
    RET
