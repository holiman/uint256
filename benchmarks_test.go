// Copyright 2019-2020 uint256 Authors.
// Use of this source code is governed by a BSD-style license that can be found
// in the COPYING file.
//

package uint256

import (
	"math/big"
	"math/rand"
	"testing"
)

const numSamples = 1024

var (
	int32Samples    [numSamples]Int
	int32SamplesLt  [numSamples]Int
	int64Samples    [numSamples]Int
	int128Samples   [numSamples]Int
	int192Samples   [numSamples]Int
	int256Samples   [numSamples]Int
	int256SamplesLt [numSamples]Int // int256SamplesLt[i] <= int256Samples[i]

	big32Samples    [numSamples]big.Int
	big32SamplesLt  [numSamples]big.Int
	big64Samples    [numSamples]big.Int
	big128Samples   [numSamples]big.Int
	big192Samples   [numSamples]big.Int
	big256Samples   [numSamples]big.Int
	big256SamplesLt [numSamples]big.Int // big256SamplesLt[i] <= big256Samples[i]

	_ = initSamples()
)

func initSamples() bool {
	rnd := rand.New(rand.NewSource(0))

	// newRandInt creates new Int with so many highly likely non-zero random words.
	newRandInt := func(numWords int) Int {
		var z Int
		for i := 0; i < numWords; i++ {
			z[i] = rnd.Uint64()
		}
		return z
	}

	for i := 0; i < numSamples; i++ {
		x32g := rnd.Uint32()
		x32l := rnd.Uint32()
		if x32g < x32l {
			x32g, x32l = x32l, x32g
		}
		int32Samples[i].SetUint64(uint64(x32g))
		big32Samples[i] = *int32Samples[i].ToBig()
		int32SamplesLt[i].SetUint64(uint64(x32l))
		big32SamplesLt[i] = *int32SamplesLt[i].ToBig()

		int64Samples[i] = newRandInt(1)
		big64Samples[i] = *int64Samples[i].ToBig()

		int128Samples[i] = newRandInt(2)
		big128Samples[i] = *int128Samples[i].ToBig()

		int192Samples[i] = newRandInt(3)
		big192Samples[i] = *int192Samples[i].ToBig()

		int256Samples[i] = newRandInt(4)
		int256SamplesLt[i] = newRandInt(4)
		if int256Samples[i].Lt(&int256SamplesLt[i]) {
			int256Samples[i], int256SamplesLt[i] = int256SamplesLt[i], int256Samples[i]
		}
		big256Samples[i] = *int256Samples[i].ToBig()
		big256SamplesLt[i] = *int256SamplesLt[i].ToBig()
	}

	return true
}

func benchmark_Add_Bit(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	f, _ := FromBig(b1)
	f2, _ := FromBig(b2)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f.Add(f, f2)
	}
}
func benchmark_Add_Big(bench *testing.B) {
	b := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b.Add(b, b2)
	}
}
func Benchmark_Add(bench *testing.B) {
	bench.Run("big", benchmark_Add_Big)
	bench.Run("uint256", benchmark_Add_Bit)
}

func benchmark_SubOverflow_Bit(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	f, _ := FromBig(b1)
	f2, _ := FromBig(b2)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f.SubOverflow(f, f2)
	}
}
func benchmark_Sub_Bit(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	f, _ := FromBig(b1)
	f2, _ := FromBig(b2)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f.Sub(f, f2)
	}
}

func benchmark_Sub_Big(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1.Sub(b1, b2)
	}
}
func Benchmark_Sub(bench *testing.B) {
	bench.Run("big", benchmark_Sub_Big)
	bench.Run("uint256", benchmark_Sub_Bit)
	bench.Run("uint256_of", benchmark_SubOverflow_Bit)
}

func benchmark_Mul_Big(bench *testing.B) {
	a := big.NewInt(0).SetBytes(hex2Bytes("f123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b := big.NewInt(0).SetBytes(hex2Bytes("f123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1 := big.NewInt(0)
		b1.Mul(a, b)
		U256(b1)
	}
}

func benchmark_Mul_Bit(bench *testing.B) {
	a := big.NewInt(0).SetBytes(hex2Bytes("f123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b := big.NewInt(0).SetBytes(hex2Bytes("f123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	fa, _ := FromBig(a)
	fb, _ := FromBig(b)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f := NewInt()
		f.Mul(fa, fb)
	}
}

func benchmark_Squared_Bit(bench *testing.B) {
	a := big.NewInt(0).SetBytes(hex2Bytes("f123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	fa, _ := FromBig(a)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f := NewInt().Copy(fa)
		f.Squared()
	}
}
func benchmark_Squared_Big(bench *testing.B) {
	a := big.NewInt(0).SetBytes(hex2Bytes("f123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1 := big.NewInt(0)
		b1.Mul(a, a)
		U256(b1)
	}
}

func Benchmark_Mul(bench *testing.B) {
	bench.Run("big", benchmark_Mul_Big)
	bench.Run("uint256", benchmark_Mul_Bit)
}
func Benchmark_Square(bench *testing.B) {
	bench.Run("big", benchmark_Squared_Big)
	bench.Run("uint256", benchmark_Squared_Bit)
}

func benchmark_And_Big(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1.And(b1, b2)
	}
}
func benchmark_And_Bit(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	f, _ := FromBig(b1)
	f2, _ := FromBig(b2)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f.And(f, f2)
	}
}
func Benchmark_And(bench *testing.B) {
	bench.Run("big", benchmark_And_Big)
	bench.Run("uint256", benchmark_And_Bit)
}

func benchmark_Or_Big(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1.Or(b1, b2)
	}
}
func benchmark_Or_Bit(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	f, _ := FromBig(b1)
	f2, _ := FromBig(b2)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f.Or(f, f2)
	}
}
func Benchmark_Or(bench *testing.B) {
	bench.Run("big", benchmark_Or_Big)
	bench.Run("uint256", benchmark_Or_Bit)
}

func benchmark_Xor_Big(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1.Xor(b1, b2)
	}
}
func benchmark_Xor_Bit(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	f, _ := FromBig(b1)
	f2, _ := FromBig(b2)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f.Xor(f, f2)
	}
}

func Benchmark_Xor(bench *testing.B) {
	bench.Run("big", benchmark_Xor_Big)
	bench.Run("uint256", benchmark_Xor_Bit)
}

func benchmark_Cmp_Big(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1.Cmp(b2)
	}
}
func benchmark_Cmp_Bit(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	f, _ := FromBig(b1)
	f2, _ := FromBig(b2)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f.Cmp(f2)
	}
}
func Benchmark_Cmp(bench *testing.B) {
	bench.Run("big", benchmark_Cmp_Big)
	bench.Run("uint256", benchmark_Cmp_Bit)
}

func benchmark_Lsh_Big(n uint, bench *testing.B) {
	original := big.NewInt(0).SetBytes(hex2Bytes("FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"))
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1 := big.NewInt(0)
		b1.Lsh(original, n)
	}
}
func benchmark_Lsh_Big_N_EQ_0(bench *testing.B) {
	benchmark_Lsh_Big(0, bench)
}
func benchmark_Lsh_Big_N_GT_192(bench *testing.B) {
	benchmark_Lsh_Big(193, bench)
}
func benchmark_Lsh_Big_N_GT_128(bench *testing.B) {
	benchmark_Lsh_Big(129, bench)
}
func benchmark_Lsh_Big_N_GT_64(bench *testing.B) {
	benchmark_Lsh_Big(65, bench)
}
func benchmark_Lsh_Big_N_GT_0(bench *testing.B) {
	benchmark_Lsh_Big(1, bench)
}
func benchmark_Lsh_Bit(n uint, bench *testing.B) {
	original := big.NewInt(0).SetBytes(hex2Bytes("FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"))
	f2, _ := FromBig(original)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f1 := NewInt()
		f1.Lsh(f2, n)
	}
}
func benchmark_Lsh_Bit_N_EQ_0(bench *testing.B) {
	benchmark_Lsh_Bit(0, bench)
}
func benchmark_Lsh_Bit_N_GT_192(bench *testing.B) {
	benchmark_Lsh_Bit(193, bench)
}
func benchmark_Lsh_Bit_N_GT_128(bench *testing.B) {
	benchmark_Lsh_Bit(129, bench)
}
func benchmark_Lsh_Bit_N_GT_64(bench *testing.B) {
	benchmark_Lsh_Bit(65, bench)
}
func benchmark_Lsh_Bit_N_GT_0(bench *testing.B) {
	benchmark_Lsh_Bit(1, bench)
}
func Benchmark_Lsh(bench *testing.B) {
	bench.Run("big/n_eq_0", benchmark_Lsh_Big_N_EQ_0)
	bench.Run("big/n_gt_192", benchmark_Lsh_Big_N_GT_192)
	bench.Run("big/n_gt_128", benchmark_Lsh_Big_N_GT_128)
	bench.Run("big/n_gt_64", benchmark_Lsh_Big_N_GT_64)
	bench.Run("big/n_gt_0", benchmark_Lsh_Big_N_GT_0)

	bench.Run("uint256/n_eq_0", benchmark_Lsh_Bit_N_EQ_0)
	bench.Run("uint256/n_gt_192", benchmark_Lsh_Bit_N_GT_192)
	bench.Run("uint256/n_gt_128", benchmark_Lsh_Bit_N_GT_128)
	bench.Run("uint256/n_gt_64", benchmark_Lsh_Bit_N_GT_64)
	bench.Run("uint256/n_gt_0", benchmark_Lsh_Bit_N_GT_0)
}

func benchmark_Rsh_Big(n uint, bench *testing.B) {
	original := big.NewInt(0).SetBytes(hex2Bytes("FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"))
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1 := big.NewInt(0)
		b1.Rsh(original, n)
	}
}
func benchmark_Rsh_Big_N_EQ_0(bench *testing.B) {
	benchmark_Rsh_Big(0, bench)
}
func benchmark_Rsh_Big_N_GT_192(bench *testing.B) {
	benchmark_Rsh_Big(193, bench)
}
func benchmark_Rsh_Big_N_GT_128(bench *testing.B) {
	benchmark_Rsh_Big(129, bench)
}
func benchmark_Rsh_Big_N_GT_64(bench *testing.B) {
	benchmark_Rsh_Big(65, bench)
}
func benchmark_Rsh_Big_N_GT_0(bench *testing.B) {
	benchmark_Rsh_Big(1, bench)
}
func benchmark_Rsh_Bit(n uint, bench *testing.B) {
	original := big.NewInt(0).SetBytes(hex2Bytes("FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"))
	f2, _ := FromBig(original)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f1 := NewInt()
		f1.Rsh(f2, n)
	}
}
func benchmark_Rsh_Bit_N_EQ_0(bench *testing.B) {
	benchmark_Rsh_Bit(0, bench)
}
func benchmark_Rsh_Bit_N_GT_192(bench *testing.B) {
	benchmark_Rsh_Bit(193, bench)
}
func benchmark_Rsh_Bit_N_GT_128(bench *testing.B) {
	benchmark_Rsh_Bit(129, bench)
}
func benchmark_Rsh_Bit_N_GT_64(bench *testing.B) {
	benchmark_Rsh_Bit(65, bench)
}
func benchmark_Rsh_Bit_N_GT_0(bench *testing.B) {
	benchmark_Rsh_Bit(1, bench)
}
func Benchmark_Rsh(bench *testing.B) {
	bench.Run("big/n_eq_0", benchmark_Rsh_Big_N_EQ_0)
	bench.Run("big/n_gt_192", benchmark_Rsh_Big_N_GT_192)
	bench.Run("big/n_gt_128", benchmark_Rsh_Big_N_GT_128)
	bench.Run("big/n_gt_64", benchmark_Rsh_Big_N_GT_64)
	bench.Run("big/n_gt_0", benchmark_Rsh_Big_N_GT_0)

	bench.Run("uint256/n_eq_0", benchmark_Rsh_Bit_N_EQ_0)
	bench.Run("uint256/n_gt_192", benchmark_Rsh_Bit_N_GT_192)
	bench.Run("uint256/n_gt_128", benchmark_Rsh_Bit_N_GT_128)
	bench.Run("uint256/n_gt_64", benchmark_Rsh_Bit_N_GT_64)
	bench.Run("uint256/n_gt_0", benchmark_Rsh_Bit_N_GT_0)
}

func benchmark_Exp_Big(bench *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	orig := big.NewInt(0).SetBytes(hex2Bytes(x))
	base := big.NewInt(0).SetBytes(hex2Bytes(x))
	exp := big.NewInt(0).SetBytes(hex2Bytes(y))

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		Exp(base, exp)
		base.Set(orig)
	}
}
func benchmark_Exp_Bit(bench *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	base := big.NewInt(0).SetBytes(hex2Bytes(x))
	exp := big.NewInt(0).SetBytes(hex2Bytes(y))

	f_base, _ := FromBig(base)
	f_orig, _ := FromBig(base)
	f_exp, _ := FromBig(exp)
	f_res := Int{}

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f_res.Exp(f_base, f_exp)
		f_base.Copy(f_orig)
	}
}
func benchmark_ExpSmall_Big(bench *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "8abcdef"

	orig := big.NewInt(0).SetBytes(hex2Bytes(x))
	base := big.NewInt(0).SetBytes(hex2Bytes(x))
	exp := big.NewInt(0).SetBytes(hex2Bytes(y))

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		Exp(base, exp)
		base.Set(orig)
	}
}
func benchmark_ExpSmall_Bit(bench *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "8abcdef"

	base := big.NewInt(0).SetBytes(hex2Bytes(x))
	exp := big.NewInt(0).SetBytes(hex2Bytes(y))

	f_base, _ := FromBig(base)
	f_orig, _ := FromBig(base)
	f_exp, _ := FromBig(exp)
	f_res := Int{}

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f_res.Exp(f_base, f_exp)
		f_base.Copy(f_orig)
	}
}
func Benchmark_Exp(bench *testing.B) {
	bench.Run("large/big", benchmark_Exp_Big)
	bench.Run("large/uint256", benchmark_Exp_Bit)
	bench.Run("small/big", benchmark_ExpSmall_Big)
	bench.Run("small/uint256", benchmark_ExpSmall_Bit)
}

func Benchmark_Div(bench *testing.B) {
	bench.Run("large/big", benchmark_DivLarge_Big)
	bench.Run("large/uint256", benchmark_DivLarge_Bit)

	bench.Run("small/big", benchmark_DivSmall_Big)
	bench.Run("small/uint256", benchmark_DivSmall_Bit)
}

func benchMulModBigint(a, b, m string) func(*testing.B) {
	return func(bench *testing.B) {
		x := big.NewInt(0).SetBytes(hex2Bytes(a))
		y := big.NewInt(0).SetBytes(hex2Bytes(b))
		z := big.NewInt(0).SetBytes(hex2Bytes(m))
		bench.ResetTimer()
		for i := 0; i < bench.N; i++ {
			b1 := big.NewInt(0)
			b1.Mul(x, y)
			b1.Mod(b1, z)
			U256(b1)
		}
	}
}
func benchMulModUint256(a, b, m string) func(*testing.B) {
	return func(bench *testing.B) {
		x, _ := FromBig(big.NewInt(0).SetBytes(hex2Bytes(a)))
		y, _ := FromBig(big.NewInt(0).SetBytes(hex2Bytes(b)))
		z, _ := FromBig(big.NewInt(0).SetBytes(hex2Bytes(m)))
		bench.ResetTimer()
		for i := 0; i < bench.N; i++ {
			f := NewInt()
			f.MulMod(x, y, z)
		}
	}
}

func Benchmark_MulMod(bench *testing.B) {
	a := "fefefefefefefefefefefefefefefefefefefefefefefefefefefefefefeffef"
	b := "efefefefefefefefefefefefefefefefefefefefefefefefefefefefefefeff9"
	m := "defefefefefefefefefefefefefefefefefefefefefefefefefefefefefefefe"
	bench.Run("large/big", benchMulModBigint(a, b, m))
	bench.Run("large/uint256", benchMulModUint256(a, b, m))
	a = "00000000000000000000000000000000000000000000000000000000fefeffef"
	b = "00000000000000000000000000000000000000000000000000000000000feff9"
	m = "00000000000000000000000000000000000000000000000000000000000000fe"
	bench.Run("small/big", benchMulModBigint(a, b, m))
	bench.Run("small/uint256", benchMulModUint256(a, b, m))
}

func BenchmarkDiv(b *testing.B) {
	benchmarkDivUint256 := func(b *testing.B, xSamples, modSamples *[numSamples]Int) {
		var sink Int
		for j := 0; j < b.N; j += numSamples {
			for i := 0; i < numSamples; i++ {
				sink.Div(&xSamples[i], &modSamples[i])
			}
		}
	}
	benchmarkDivBig := func(b *testing.B, xSamples, modSamples *[numSamples]big.Int) {
		var sink big.Int
		for j := 0; j < b.N; j += numSamples {
			for i := 0; i < numSamples; i++ {
				sink.Div(&xSamples[i], &modSamples[i])
			}
		}
	}

	b.Run("small/uint256", func(b *testing.B) { benchmarkDivUint256(b, &int32Samples, &int32SamplesLt) })
	b.Run("small/big", func(b *testing.B) { benchmarkDivBig(b, &big32Samples, &big32SamplesLt) })
	b.Run("mod64/uint256", func(b *testing.B) { benchmarkDivUint256(b, &int256Samples, &int64Samples) })
	b.Run("mod64/big", func(b *testing.B) { benchmarkDivBig(b, &big256Samples, &big64Samples) })
	b.Run("mod128/uint256", func(b *testing.B) { benchmarkDivUint256(b, &int256Samples, &int128Samples) })
	b.Run("mod128/big", func(b *testing.B) { benchmarkDivBig(b, &big256Samples, &big128Samples) })
	b.Run("mod192/uint256", func(b *testing.B) { benchmarkDivUint256(b, &int256Samples, &int192Samples) })
	b.Run("mod192/big", func(b *testing.B) { benchmarkDivBig(b, &big256Samples, &big192Samples) })
	b.Run("mod256/uint256", func(b *testing.B) { benchmarkDivUint256(b, &int256Samples, &int256SamplesLt) })
	b.Run("mod256/big", func(b *testing.B) { benchmarkDivBig(b, &big256Samples, &big256SamplesLt) })
}

func BenchmarkMod(b *testing.B) {
	benchmarkModUint256 := func(b *testing.B, xSamples, modSamples *[numSamples]Int) {
		var sink Int
		for j := 0; j < b.N; j += numSamples {
			for i := 0; i < numSamples; i++ {
				sink.Mod(&xSamples[i], &modSamples[i])
			}
		}
	}
	benchmarkModBig := func(b *testing.B, xSamples, modSamples *[numSamples]big.Int) {
		var sink big.Int
		for j := 0; j < b.N; j += numSamples {
			for i := 0; i < numSamples; i++ {
				sink.Mod(&xSamples[i], &modSamples[i])
			}
		}
	}

	b.Run("small/uint256", func(b *testing.B) { benchmarkModUint256(b, &int32Samples, &int32SamplesLt) })
	b.Run("small/big", func(b *testing.B) { benchmarkModBig(b, &big32Samples, &big32SamplesLt) })
	b.Run("mod64/uint256", func(b *testing.B) { benchmarkModUint256(b, &int256Samples, &int64Samples) })
	b.Run("mod64/big", func(b *testing.B) { benchmarkModBig(b, &big256Samples, &big64Samples) })
	b.Run("mod128/uint256", func(b *testing.B) { benchmarkModUint256(b, &int256Samples, &int128Samples) })
	b.Run("mod128/big", func(b *testing.B) { benchmarkModBig(b, &big256Samples, &big128Samples) })
	b.Run("mod192/uint256", func(b *testing.B) { benchmarkModUint256(b, &int256Samples, &int192Samples) })
	b.Run("mod192/big", func(b *testing.B) { benchmarkModBig(b, &big256Samples, &big192Samples) })
	b.Run("mod256/uint256", func(b *testing.B) { benchmarkModUint256(b, &int256Samples, &int256SamplesLt) })
	b.Run("mod256/big", func(b *testing.B) { benchmarkModBig(b, &big256Samples, &big256SamplesLt) })
}

func BenchmarkMulMod(b *testing.B) {
	benchmarkMulModUint256 := func(b *testing.B, factorsSamples, modSamples *[numSamples]Int) {
		var sink, x Int
		for j := 0; j < b.N; j += numSamples {
			for i := 0; i < numSamples; i++ {
				y := factorsSamples[i]
				sink.MulMod(&x, &y, &modSamples[i])
				x = y
			}
		}
	}
	benchmarkMulModBig := func(b *testing.B, factorsSamples, modSamples *[numSamples]big.Int) {
		var sink, x big.Int
		for j := 0; j < b.N; j += numSamples {
			for i := 0; i < numSamples; i++ {
				y := factorsSamples[i]
				sink.Mul(&x, &y)
				sink.Mod(&sink, &modSamples[i])
				x = y
			}
		}
	}

	b.Run("small/uint256", func(b *testing.B) { benchmarkMulModUint256(b, &int32Samples, &int32SamplesLt) })
	b.Run("small/big", func(b *testing.B) { benchmarkMulModBig(b, &big32Samples, &big32SamplesLt) })
	b.Run("mod64/uint256", func(b *testing.B) { benchmarkMulModUint256(b, &int256Samples, &int64Samples) })
	b.Run("mod64/big", func(b *testing.B) { benchmarkMulModBig(b, &big256Samples, &big64Samples) })
	b.Run("mod128/uint256", func(b *testing.B) { benchmarkMulModUint256(b, &int256Samples, &int128Samples) })
	b.Run("mod128/big", func(b *testing.B) { benchmarkMulModBig(b, &big256Samples, &big128Samples) })
	b.Run("mod192/uint256", func(b *testing.B) { benchmarkMulModUint256(b, &int256Samples, &int192Samples) })
	b.Run("mod192/big", func(b *testing.B) { benchmarkMulModBig(b, &big256Samples, &big192Samples) })
	b.Run("mod256/uint256", func(b *testing.B) { benchmarkMulModUint256(b, &int256Samples, &int256SamplesLt) })
	b.Run("mod256/big", func(b *testing.B) { benchmarkMulModBig(b, &big256Samples, &big256SamplesLt) })
}

func benchModBigint(a, b string) func(*testing.B) {
	return func(bench *testing.B) {
		x := big.NewInt(0).SetBytes(hex2Bytes(a))
		y := big.NewInt(0).SetBytes(hex2Bytes(b))
		z := big.NewInt(0)
		bench.ResetTimer()
		for i := 0; i < bench.N; i++ {
			z.Mod(x, y)
		}
	}
}
func benchModUint256(a, b string) func(*testing.B) {
	return func(bench *testing.B) {
		x, _ := FromBig(big.NewInt(0).SetBytes(hex2Bytes(a)))
		y, _ := FromBig(big.NewInt(0).SetBytes(hex2Bytes(b)))
		z := NewInt()

		bench.ResetTimer()
		for i := 0; i < bench.N; i++ {
			z.Mod(x, y)
		}
	}
}
func Benchmark_Mod(bench *testing.B) {
	a := "fefefefefefefefefefefefefefefefefefefefefefefefefefefefefefeffef"
	b := "efefefefefefefefefefefefefefefefefefefefefefefefefefefefefefeff9"
	bench.Run("large/big", benchModBigint(a, b))
	bench.Run("large/uint256", benchModUint256(a, b))

	a = "0000000000000000000000000000000000000000000000000000000000feffef"
	b = "00000000000000000000000000000000000000000000000000000000000000f9"
	bench.Run("small/big", benchModBigint(a, b))
	bench.Run("small/uint256", benchModUint256(a, b))

}

func benchmark_DivSmall_Big(bench *testing.B) {
	a := big.NewInt(0).SetBytes(hex2Bytes("1fc2bad1e611"))
	b := big.NewInt(0).SetBytes(hex2Bytes("12bad1e611"))

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1 := big.NewInt(0)
		b1.Div(a, b)
		U256(b1)
	}
}

func benchmark_DivSmall_Bit(bench *testing.B) {
	a := big.NewInt(0).SetBytes(hex2Bytes("1fc2bad1e611"))
	b := big.NewInt(0).SetBytes(hex2Bytes("12bad1e611"))
	fa, _ := FromBig(a)
	fb, _ := FromBig(b)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f := NewInt()
		f.Div(fa, fb)
	}
}
func benchmark_DivLarge_Big(bench *testing.B) {
	a := big.NewInt(0).SetBytes(hex2Bytes("fe7fb0d1f59dfe9492ffbf73683fd1e870eec79504c60144cc7f5fc2bad1e611"))
	b := big.NewInt(0).SetBytes(hex2Bytes("ff3f9014f20db29ae04af2c2d265de17"))

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1 := big.NewInt(0)
		b1.Div(a, b)
		U256(b1)
	}
}

func benchmark_DivLarge_Bit(bench *testing.B) {
	a := big.NewInt(0).SetBytes(hex2Bytes("fe7fb0d1f59dfe9492ffbf73683fd1e870eec79504c60144cc7f5fc2bad1e611"))
	b := big.NewInt(0).SetBytes(hex2Bytes("ff3f9014f20db29ae04af2c2d265de17"))
	fa, _ := FromBig(a)
	fb, _ := FromBig(b)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f := NewInt()
		f.Div(fa, fb)
	}
}
func benchmark_SdivLarge_Big(bench *testing.B) {
	a := big.NewInt(0).SetBytes(hex2Bytes("800fffffffffffffffffffffffffd1e870eec79504c60144cc7f5fc2bad1e611"))
	b := big.NewInt(0).SetBytes(hex2Bytes("ff3f9014f20db29ae04af2c2d265de17"))

	bench.ResetTimer()

	var (
		x = big.NewInt(0)
		y = big.NewInt(0)
	)

	for i := 0; i < bench.N; i++ {
		x.Set(a)
		y.Set(b)
		U256(Sdiv(S256(x), S256(y)))
	}
}

func benchmark_SdivLarge_Bit(bench *testing.B) {
	a := big.NewInt(0).SetBytes(hex2Bytes("800fffffffffffffffffffffffffd1e870eec79504c60144cc7f5fc2bad1e611"))
	b := big.NewInt(0).SetBytes(hex2Bytes("ff3f9014f20db29ae04af2c2d265de17"))
	fa, _ := FromBig(a)
	fb, _ := FromBig(b)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f := NewInt()
		f.Sdiv(fa, fb)
	}
}

func Benchmark_SDiv(bench *testing.B) {
	bench.Run("large/big", benchmark_SdivLarge_Big)
	bench.Run("large/uint256", benchmark_SdivLarge_Bit)
}

func benchAddModBigint(a, b, m string) func(*testing.B) {
	return func(bench *testing.B) {
		x := big.NewInt(0).SetBytes(hex2Bytes(a))
		y := big.NewInt(0).SetBytes(hex2Bytes(b))
		z := big.NewInt(0).SetBytes(hex2Bytes(m))
		b1 := big.NewInt(0)
		bench.ResetTimer()
		for i := 0; i < bench.N; i++ {

			b1.Add(x, y)
			b1.Mod(b1, z)
			U256(b1)
		}
	}
}
func benchAddModUint256(a, b, m string) func(*testing.B) {
	return func(bench *testing.B) {
		x, _ := FromBig(big.NewInt(0).SetBytes(hex2Bytes(a)))
		y, _ := FromBig(big.NewInt(0).SetBytes(hex2Bytes(b)))
		z, _ := FromBig(big.NewInt(0).SetBytes(hex2Bytes(m)))
		bench.ResetTimer()
		for i := 0; i < bench.N; i++ {
			f := NewInt()
			f.AddMod(x, y, z)
		}
	}
}

func Benchmark_AddMod(bench *testing.B) {
	a := "fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0"
	b := "fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa"
	m := "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	bench.Run("large/big", benchAddModBigint(a, b, m))
	bench.Run("large/uint256", benchAddModUint256(a, b, m))
	a = "00000000000000000000000000000000000000000000000000000000fefeffef"
	b = "00000000000000000000000000000000000000000000000000000000000feff9"
	m = "00000000000000000000000000000000000000000000000000000000000000fe"
	bench.Run("small/big", benchAddModBigint(a, b, m))
	bench.Run("small/uint256", benchAddModUint256(a, b, m))
}
