//go:build amd64 && !purego

#include "textflag.h"

// func andAVX2(z, x, y *[32]byte)
// andAVX2 performs z = x&y. It modifies z in-place and has no return value
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
// orAVX2 performs z = x|y. It modifies z in-place and has no return value
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
// xorAVX2 performs z = x^y. It modifies z in-place and has no return value
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
// notAVX2 performs z = ^x. It modifies z in-place and has no return value
TEXT ·notAVX2(SB), NOSPLIT, $0-16
    MOVQ z+0(FP), AX
    MOVQ x+8(FP), BX
    // Compare Y0 with itself, each 32-bit element is equal,
    // so all bits in each element are set to 1 (0xFFFFFFFF).
    VPCMPEQD Y0, Y0, Y0

    VMOVDQU (BX), Y1
    VPXOR Y0, Y1, Y2     // Y2 = Y0 XOR Y1

    VMOVDQU Y2, (AX)
    VZEROUPPER
    RET
