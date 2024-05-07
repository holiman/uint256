package uint256

import (
	"fmt"
	"math/big"
	"testing"
)

type opUnaryArgFunc func(*Int, *Int) *Int
type bigUnaryArgFunc func(*big.Int, *big.Int) *big.Int

type opCmpArgFunc func(*Int, *Int) bool
type bigCmpArgFunc func(*big.Int, *big.Int) bool

type opDualArgFunc func(*Int, *Int, *Int) *Int
type bigDualArgFunc func(*big.Int, *big.Int, *big.Int) *big.Int

type opThreeArgFunc func(*Int, *Int, *Int, *Int) *Int
type bigThreeArgFunc func(*big.Int, *big.Int, *big.Int, *big.Int) *big.Int

func FuzzSignExtend(f *testing.F) {
	f.Fuzz(func(t *testing.T, z0, z1, z2, z3 uint64, index uint8) {
		arg := &Int{z0, z1, z2, z3}
		n := new(Int).SetUint64(uint64(index))
		testSignExtend(t, arg, n)
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

var unaryOpFuncs = []struct {
	name   string
	u256Fn opUnaryArgFunc
	bigFn  bigUnaryArgFunc
}{
	{"Not", (*Int).Not, (*big.Int).Not},
	{"Neg", (*Int).Neg, (*big.Int).Neg},
	{"Sqrt", (*Int).Sqrt, (*big.Int).Sqrt},
	{"square", func(x *Int, y *Int) *Int {
		res := y.Clone()
		res.squared()
		return x.Set(res)
	}, func(b1 *big.Int, b2 *big.Int) *big.Int {
		return b1.Mul(b2, b2)
	}},
	{"Abs", (*Int).Abs, func(b *big.Int, b2 *big.Int) *big.Int {
		return b.Abs(S256(b2))
	},
	},
}

type binaryOpEntry struct {
	name   string
	u256Fn opDualArgFunc
	bigFn  bigDualArgFunc
}

var binaryOpFuncs = []binaryOpEntry{
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

	{"Exp", (*Int).Exp, func(b *big.Int, b2 *big.Int, b3 *big.Int) *big.Int {
		return b.Exp(b2, b3, bigtt256)
	}},

	{"Lsh", u256Lsh, bigLsh},
	{"Rsh", u256Rsh, bigRsh},
	{"SRsh", u256SRsh, bigSRsh},

	{"DivModDiv", divModDiv, bigDiv},
	{"DivModMod", divModMod, bigMod},
	{"udivremDiv", udivremDiv, bigDiv},
	{"udivremMod", udivremMod, bigMod},

	{"ExtendSign", (*Int).ExtendSign, bigExtendSign},
}

func lookupBinary(name string) binaryOpEntry {
	for _, tc := range binaryOpFuncs {
		if tc.name == name {
			return tc
		}
	}
	panic(fmt.Sprintf("%v not found", name))
}

var cmpOpFuncs = []struct {
	name   string
	u256Fn opCmpArgFunc
	bigFn  bigCmpArgFunc
}{
	{"Eq", (*Int).Eq, func(a, b *big.Int) bool { return a.Cmp(b) == 0 }},
	{"Lt", (*Int).Lt, func(a, b *big.Int) bool { return a.Cmp(b) < 0 }},
	{"Gt", (*Int).Gt, func(a, b *big.Int) bool { return a.Cmp(b) > 0 }},
	{"Slt", (*Int).Slt, func(a, b *big.Int) bool { return S256(a).Cmp(S256(b)) < 0 }},
	{"Sgt", (*Int).Sgt, func(a, b *big.Int) bool { return S256(a).Cmp(S256(b)) > 0 }},
	{"CmpEq", func(a, b *Int) bool { return a.Cmp(b) == 0 }, func(a, b *big.Int) bool { return a.Cmp(b) == 0 }},
	{"CmpLt", func(a, b *Int) bool { return a.Cmp(b) < 0 }, func(a, b *big.Int) bool { return a.Cmp(b) < 0 }},
	{"CmpGt", func(a, b *Int) bool { return a.Cmp(b) > 0 }, func(a, b *big.Int) bool { return a.Cmp(b) > 0 }},
	{"LtUint64", func(a, b *Int) bool { return a.LtUint64(b.Uint64()) }, func(a, b *big.Int) bool { return a.Cmp(new(big.Int).SetUint64(b.Uint64())) < 0 }},
	{"GtUint64", func(a, b *Int) bool { return a.GtUint64(b.Uint64()) }, func(a, b *big.Int) bool { return a.Cmp(new(big.Int).SetUint64(b.Uint64())) > 0 }},
}

var ternaryOpFuncs = []struct {
	name   string
	u256Fn opThreeArgFunc
	bigFn  bigThreeArgFunc
}{
	{"AddMod", (*Int).AddMod, bigAddMod},
	{"MulMod", (*Int).MulMod, bigMulMod},
	{"MulModWithReciprocal", (*Int).mulModWithReciprocalWrapper, bigMulMod},
}
