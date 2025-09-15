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

