package uint256

import (
	"fmt"
	"math/big"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type opUnaryArgFunc func(*Int, *Int) *Int
type bigUnaryArgFunc func(*big.Int, *big.Int) *big.Int

type opDualArgFunc func(*Int, *Int, *Int) *Int
type bigDualArgFunc func(*big.Int, *big.Int, *big.Int) *big.Int

type opThreeArgFunc func(*Int, *Int, *Int, *Int) *Int
type bigThreeArgFunc func(*big.Int, *big.Int, *big.Int, *big.Int) *big.Int

func crash(op interface{}, msg string, args ...Int) {
	fn := runtime.FuncForPC(reflect.ValueOf(op).Pointer())
	fnName := fn.Name()
	fnFile, fnLine := fn.FileLine(fn.Entry())
	var strArgs []string
	for i, arg := range args {
		strArgs = append(strArgs, fmt.Sprintf("%d: %x", i, &arg))
	}
	panic(fmt.Sprintf("%s\nfor %s (%s:%d)\n%v",
		msg, fnName, fnFile, fnLine, strings.Join(strArgs, "\n")))
}

func checkUnaryOp(op opUnaryArgFunc, bigOp bigUnaryArgFunc, x Int) {
	origX := x
	var result Int
	ret := op(&result, &x)
	if ret != &result {
		crash(op, "returned not the pointer receiver", x)
	}
	if x != origX {
		crash(op, "argument modified", x)
	}
	expected, _ := FromBig(bigOp(new(big.Int), x.ToBig()))
	if result != *expected {
		crash(op, "unexpected result", x)
	}
	// Test again when the receiver is not zero.
	var garbage Int
	garbage.Sub(&garbage, NewInt(1))
	ret = op(&garbage, &x)
	if ret != &garbage {
		crash(op, "returned not the pointer receiver", x)
	}
	if garbage != *expected {
		crash(op, "unexpected result", x)
	}
	// Test again with the receiver aliasing arguments.
	ret = op(&x, &x)
	if ret != &x {
		crash(op, "returned not the pointer receiver", x)
	}
	if x != *expected {
		crash(op, "unexpected result", x)
	}
}

func checkThreeArgOp(op opThreeArgFunc, bigOp bigThreeArgFunc, x, y, z Int) {
	origX := x
	origY := y
	origZ := z

	var result Int
	ret := op(&result, &x, &y, &z)
	if ret != &result {
		crash(op, "returned not the pointer receiver", x, y, z)
	}
	switch {
	case x != origX:
		crash(op, "first argument modified", x, y, z)
	case y != origY:
		crash(op, "second argument modified", x, y, z)
	case z != origZ:
		crash(op, "third argument modified", x, y, z)
	}
	expected, _ := FromBig(bigOp(new(big.Int), x.ToBig(), y.ToBig(), z.ToBig()))
	if have, want := result, *expected; have != want {
		crash(op, fmt.Sprintf("unexpected result: have %v want %v", have, want), x, y, z)
	}

	// Test again when the receiver is not zero.
	var garbage Int
	garbage.Xor(&x, &y)
	ret = op(&garbage, &x, &y, &z)
	if ret != &garbage {
		crash(op, "returned not the pointer receiver", x, y, z)
	}
	if have, want := garbage, *expected; have != want {
		crash(op, fmt.Sprintf("unexpected result: have %v want %v", have, want), x, y, z)
	}
	switch {
	case x != origX:
		crash(op, "first argument modified", x, y, z)
	case y != origY:
		crash(op, "second argument modified", x, y, z)
	case z != origZ:
		crash(op, "third argument modified", x, y, z)
	}

	// Test again with the receiver aliasing arguments.
	ret = op(&x, &x, &y, &z)
	if ret != &x {
		crash(op, "returned not the pointer receiver", x, y, z)
	}
	if have, want := x, *expected; have != want {
		crash(op, fmt.Sprintf("unexpected result: have %v want %v", have, want), x, y, z)
	}

	ret = op(&y, &origX, &y, &z)
	if ret != &y {
		crash(op, "returned not the pointer receiver", x, y, z)
	}
	if y != *expected {
		crash(op, "unexpected result", x, y, z)
	}
	ret = op(&z, &origX, &origY, &z)
	if ret != &z {
		crash(op, "returned not the pointer receiver", x, y, z)
	}
	if z != *expected {
		crash(op, fmt.Sprintf("unexpected result: have %v want %v", z.ToBig(), expected), x, y, z)
	}
}

func FuzzSignExtend(f *testing.F) {
	f.Fuzz(func(t *testing.T, z0, z1, z2, z3 uint64, index uint8) {
		arg := &Int{z0, z1, z2, z3}
		n := new(Int).SetUint64(uint64(index))
		testSignExtend(t, arg, n)
	})
}

func FuzzUnaryOp(f *testing.F) {
	f.Fuzz(func(t *testing.T, x0, x1, x2, x3 uint64) {
		x := Int{x0, x1, x2, x3}
		checkUnaryOp((*Int).Sqrt, (*big.Int).Sqrt, x)
	})
}

func u256Lsh(z, x, y *Int) *Int {
	return z.Lsh(x, uint(y.Uint64()&0x1FF))
}
func u256Rsh(z, x, y *Int) *Int {
	return z.Rsh(x, uint(y.Uint64()&0x1FF))
}
func u256SRsh(z, x, y *Int) *Int {
	return z.SRsh(x, uint(y.Uint64()&0x1FF))
}

func bigLsh(z, x, y *big.Int) *big.Int {
	return z.Lsh(x, uint(y.Uint64()&0x1FF))
}

func bigRsh(z, x, y *big.Int) *big.Int {
	return z.Rsh(x, uint(y.Uint64()&0x1FF))
}
func bigSRsh(z, x, y *big.Int) *big.Int {
	n := uint(y.Uint64() & 0x1FF)
	x = S256(x)
	return z.Rsh(x, n)
}

var binaryOpFuncs = []struct {
	name   string
	u256Fn opDualArgFunc
	bigFn  bigDualArgFunc
}{
	{"Add", (*Int).Add, (*big.Int).Add},
	{"Sub", (*Int).Sub, (*big.Int).Sub},
	{"Mul", (*Int).Mul, (*big.Int).Mul},
	{"Div", (*Int).Div, bigDiv},
	{"Mod", (*Int).Mod, bigMod},
	{"SDiv", (*Int).SDiv, bigSDiv},
	{"SMod", (*Int).SMod, bigSMod},
	{"And", (*Int).And, (*big.Int).And},
	{"Or", (*Int).Or, (*big.Int).Or},
	{"Xor", (*Int).Xor, (*big.Int).Xor},

	{"Lsh", u256Lsh, bigLsh},
	{"Rsh", u256Rsh, bigRsh},
	{"SRsh", u256SRsh, bigSRsh},

	{"DivModDiv", divModDiv, bigDiv},
	{"DivModMod", divModMod, bigMod},
	{"udivremDiv", udivremDiv, bigDiv},
	{"udivremMod", udivremMod, bigMod},

	{"ExtendSign", (*Int).ExtendSign, bigExtendSign},
}

//func FuzzBinaryOperations(f *testing.F) {
//	f.Fuzz(func(t *testing.T, x0, x1, x2, x3, y0, y1, y2, y3 uint64) {
//
//		x := Int{x0, x1, x2, x3}
//		y := Int{y0, y1, y2, y3}
//
//		for _, tc := range binaryOpFuncs {
//			checkDualArgOp(tc.u256Fn, tc.bigFn, x, y)
//		}
//	})
//}

func bigintMulMod(b1, b2, b3, b4 *big.Int) *big.Int {
	return b1.Mod(big.NewInt(0).Mul(b2, b3), b4)
}

func bigintAddMod(b1, b2, b3, b4 *big.Int) *big.Int {
	return b1.Mod(big.NewInt(0).Add(b2, b3), b4)
}

func bigintMulDiv(b1, b2, b3, b4 *big.Int) *big.Int {
	b1.Mul(b2, b3)
	return b1.Div(b1, b4)
}

func intMulDiv(f1, f2, f3, f4 *Int) *Int {
	f1.MulDivOverflow(f2, f3, f4)
	return f1
}

func FuzzTernaryOp(f *testing.F) {
	f.Fuzz(func(t *testing.T,
		x0, x1, x2, x3,
		y0, y1, y2, y3,
		z0, z1, z2, z3 uint64) {
		x := Int{x0, x1, x2, x3}
		y := Int{y0, y1, y2, y3}
		z := Int{z0, z1, z2, z3}
		if z.IsZero() {
			return
		}
		checkThreeArgOp((*Int).MulMod, bigintMulMod, x, y, z)
		checkThreeArgOp((*Int).AddMod, bigintAddMod, x, y, z)
		checkThreeArgOp(intMulDiv, bigintMulDiv, x, y, z)
	})
}
