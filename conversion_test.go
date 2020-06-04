// uint256: Fixed size 256-bit math library
// Copyright 2020 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

package uint256

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"
)

var (
	_ fmt.Formatter = &Int{} // Test if Int supports Formatter interface.
)

func TestFromBig(t *testing.T) {
	a := new(big.Int)
	b, o := FromBig(a)
	if o {
		t.Fatalf("conversion overflowed! big.Int %x", a.Bytes())
	}
	if exp, got := a.Bytes(), b.Bytes(); !bytes.Equal(got, exp) {
		t.Fatalf("got %x exp %x", got, exp)
	}

	a = big.NewInt(1)
	b, o = FromBig(a)
	if o {
		t.Fatalf("conversion overflowed! big.Int %x", a.Bytes())
	}
	if exp, got := a.Bytes(), b.Bytes(); !bytes.Equal(got, exp) {
		t.Fatalf("got %x exp %x", got, exp)
	}

	a = big.NewInt(0x1000000000000000)
	b, o = FromBig(a)
	if o {
		t.Fatalf("conversion overflowed! big.Int %x", a.Bytes())
	}
	if exp, got := a.Bytes(), b.Bytes(); !bytes.Equal(got, exp) {
		t.Fatalf("got %x exp %x", got, exp)
	}

	a = big.NewInt(0x1234)
	b, o = FromBig(a)
	if o {
		t.Fatalf("conversion overflowed! big.Int %x", a.Bytes())
	}
	if exp, got := a.Bytes(), b.Bytes(); !bytes.Equal(got, exp) {
		t.Fatalf("got %x exp %x", got, exp)
	}

	a = big.NewInt(1)
	a.Lsh(a, 256)

	b, o = FromBig(a)
	if !o {
		t.Fatalf("expected overflow")
	}
	if !b.Eq(new(Int)) {
		t.Fatalf("got %x exp 0", b.Bytes())
	}

	a.Sub(a, big.NewInt(1))
	b, o = FromBig(a)
	if o {
		t.Fatalf("conversion overflowed! big.Int %x", a.Bytes())
	}
	if exp, got := a.Bytes(), b.Bytes(); !bytes.Equal(got, exp) {
		t.Fatalf("got %x exp %x", got, exp)
	}
}

func TestFromBigOverflow(t *testing.T) {
	_, o := FromBig(new(big.Int).SetBytes(hex2Bytes("ababee444444444444ffcc333333333333ddaa222222222222bb8811111111111199")))
	if !o {
		t.Errorf("expected overflow, got %v", o)
	}
	_, o = FromBig(new(big.Int).SetBytes(hex2Bytes("ee444444444444ffcc333333333333ddaa222222222222bb8811111111111199")))
	if o {
		t.Errorf("expected no overflow, got %v", o)
	}
	b := new(big.Int).SetBytes(hex2Bytes("ee444444444444ffcc333333333333ddaa222222222222bb8811111111111199"))
	_, o = FromBig(b.Neg(b))
	if o {
		t.Errorf("expected no overflow, got %v", o)
	}
}

func TestToBig(t *testing.T) {

	if bigZero := new(Int).ToBig(); bigZero.Cmp(new(big.Int)) != 0 {
		t.Errorf("expected big.Int 0, got %x", bigZero)
	}

	for i := uint(0); i < 256; i++ {
		f := new(Int).SetUint64(1)
		f.Lsh(f, i)
		b := f.ToBig()
		expected := big.NewInt(1)
		expected.Lsh(expected, i)
		if b.Cmp(expected) != 0 {
			t.Fatalf("expected %x, got %x", expected, b)
		}
	}
}

func benchmarkSetFromBig(bench *testing.B, b *big.Int) Int {
	var f Int
	for i := 0; i < bench.N; i++ {
		f.SetFromBig(b)
	}
	return f
}

func BenchmarkSetFromBig(bench *testing.B) {
	param1 := big.NewInt(0xff)
	bench.Run("1word", func(bench *testing.B) { benchmarkSetFromBig(bench, param1) })

	param2 := new(big.Int).Lsh(param1, 64)
	bench.Run("2words", func(bench *testing.B) { benchmarkSetFromBig(bench, param2) })

	param3 := new(big.Int).Lsh(param2, 64)
	bench.Run("3words", func(bench *testing.B) { benchmarkSetFromBig(bench, param3) })

	param4 := new(big.Int).Lsh(param3, 64)
	bench.Run("4words", func(bench *testing.B) { benchmarkSetFromBig(bench, param4) })

	param5 := new(big.Int).Lsh(param4, 64)
	bench.Run("overflow", func(bench *testing.B) { benchmarkSetFromBig(bench, param5) })
}

func benchmarkToBig(bench *testing.B, f *Int) *big.Int {
	var b *big.Int
	for i := 0; i < bench.N; i++ {
		b = f.ToBig()
	}
	return b
}

func BenchmarkToBig(bench *testing.B) {
	param1 := new(Int).SetUint64(0xff)
	bench.Run("1word", func(bench *testing.B) { benchmarkToBig(bench, param1) })

	param2 := new(Int).Lsh(param1, 64)
	bench.Run("2words", func(bench *testing.B) { benchmarkToBig(bench, param2) })

	param3 := new(Int).Lsh(param2, 64)
	bench.Run("3words", func(bench *testing.B) { benchmarkToBig(bench, param3) })

	param4 := new(Int).Lsh(param3, 64)
	bench.Run("4words", func(bench *testing.B) { benchmarkToBig(bench, param4) })
}

func TestFormat(t *testing.T) {
	testCases := []string{
		"0",
		"1",
		"ffeeddccbbaa99887766554433221100ffeeddccbbaa99887766554433221100",
	}

	for i := 0; i < len(testCases); i++ {
		expected := testCases[i]
		b, _ := new(big.Int).SetString(expected, 16)
		f, o := FromBig(b)
		if o {
			t.Fatalf("too big test case %s", expected)
		}
		s := fmt.Sprintf("%x", f)
		if s != expected {
			t.Errorf("Invalid format conversion to hex: %s, expected %s", s, expected)
		}
	}
}

// TestSetBytes tests all setbyte-methods from 0 to overlong,
// - verifies that all non-set bits are properly cleared
// - verifies that overlong input is correctly cropped
func TestSetBytes(t *testing.T) {
	for i := 0; i < 35; i++ {
		buf := hex2Bytes("aaaa12131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f3031bbbb")
		exp, _ := FromBig(new(big.Int).SetBytes(buf[0:i]))
		z := NewInt().SetAllOne().SetBytes(buf[0:i])
		if !z.Eq(exp) {
			t.Errorf("testcase %d: exp %x, got %x", i, exp, z)
		}
	}
	// nil check
	exp, _ := FromBig(new(big.Int).SetBytes(nil))
	z := NewInt().SetAllOne().SetBytes(nil)
	if !z.Eq(exp) {
		t.Errorf("nil-test : exp %x, got %x", exp, z)
	}
}

func BenchmarkSetBytes(b *testing.B) {

	val := NewInt()
	bytearr := hex2Bytes("12131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f3031")
	b.Run("generic", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			val.SetBytes(bytearr[:1])
			val.SetBytes(bytearr[:2])
			val.SetBytes(bytearr[:3])
			val.SetBytes(bytearr[:4])
			val.SetBytes(bytearr[:5])
			val.SetBytes(bytearr[:6])
			val.SetBytes(bytearr[:7])
			val.SetBytes(bytearr[:8])
			val.SetBytes(bytearr[:9])
			val.SetBytes(bytearr[:10])
			val.SetBytes(bytearr[:11])
			val.SetBytes(bytearr[:12])
			val.SetBytes(bytearr[:13])
			val.SetBytes(bytearr[:14])
			val.SetBytes(bytearr[:15])
			val.SetBytes(bytearr[:16])
			val.SetBytes(bytearr[:17])
			val.SetBytes(bytearr[:18])
			val.SetBytes(bytearr[:19])
			val.SetBytes(bytearr[:20])
			val.SetBytes(bytearr[:21])
			val.SetBytes(bytearr[:22])
			val.SetBytes(bytearr[:23])
			val.SetBytes(bytearr[:24])
			val.SetBytes(bytearr[:25])
			val.SetBytes(bytearr[:26])
			val.SetBytes(bytearr[:27])
			val.SetBytes(bytearr[:28])
			val.SetBytes(bytearr[:29])
			val.SetBytes(bytearr[:20])
			val.SetBytes(bytearr[:31])
			val.SetBytes(bytearr[:32])
		}
	})
	b.Run("specific", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			val.SetBytes1(bytearr)
			val.SetBytes2(bytearr)
			val.SetBytes3(bytearr)
			val.SetBytes4(bytearr)
			val.SetBytes5(bytearr)
			val.SetBytes6(bytearr)
			val.SetBytes7(bytearr)
			val.SetBytes8(bytearr)
			val.SetBytes9(bytearr)
			val.SetBytes10(bytearr)
			val.SetBytes11(bytearr)
			val.SetBytes12(bytearr)
			val.SetBytes13(bytearr)
			val.SetBytes14(bytearr)
			val.SetBytes15(bytearr)
			val.SetBytes16(bytearr)
			val.SetBytes17(bytearr)
			val.SetBytes18(bytearr)
			val.SetBytes19(bytearr)
			val.SetBytes20(bytearr)
			val.SetBytes21(bytearr)
			val.SetBytes22(bytearr)
			val.SetBytes23(bytearr)
			val.SetBytes24(bytearr)
			val.SetBytes25(bytearr)
			val.SetBytes26(bytearr)
			val.SetBytes27(bytearr)
			val.SetBytes28(bytearr)
			val.SetBytes29(bytearr)
			val.SetBytes30(bytearr)
			val.SetBytes31(bytearr)
			val.SetBytes32(bytearr)
		}
	})
}
