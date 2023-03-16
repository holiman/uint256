// uint256: Fixed size 256-bit math library
// Copyright 2020 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

package uint256

import (
	"bufio"
	"bytes"
	"encoding/json"
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

func TestScanScientific(t *testing.T) {
	intsub1 := new(Int)
	_ = intsub1.fromDecimal(twoPow256Sub1)
	cases := []struct {
		in  string
		exp *Int
		err string
	}{
		{
			in:  "e30",
			err: "EOF",
		},
		{
			in:  "30e",
			err: "EOF",
		},
		{
			in:  twoPow256Sub1 + "e",
			err: "EOF",
		},
		{
			in:  "14e30",
			exp: new(Int).Mul(NewInt(14), new(Int).Exp(NewInt(10), NewInt(30))),
		},
		{ // 0xdd15fe86affad800000000000000000000000000000000000000000000000000
			in:  "1e77",
			exp: new(Int).Mul(NewInt(1), new(Int).Exp(NewInt(10), NewInt(77))),
		},
		{ // 0x8a2dbf142dfcc8000000000000000000000000000000000000000000000000000
			in:  "1e78",
			err: ErrBig256Range.Error(),
		},
		{
			in:  "1455522523e31",
			exp: new(Int).Mul(NewInt(1455522523), new(Int).Exp(NewInt(10), NewInt(31))),
		},
		{
			in:  twoPow256Sub1 + "e0",
			exp: intsub1,
		},
		{
			in:  "1e25352",
			err: ErrBig256Range.Error(),
		},
		{
			in:  "1213128763127863781263781263781263781263781263871263871268371268371263781627836128736128736127836127836127863781e0",
			err: ErrBig256Range.Error(),
		},
		{
			in:  twoPow256Sub1 + "e1",
			err: ErrBig256Range.Error(),
		},
		{
			in:  "1e253e52",
			err: `strconv.ParseUint: parsing "253e52": invalid syntax`,
		},
		{
			in:  "1e00000000000000000",
			exp: NewInt(1),
		},
	}
	for tc, v := range cases {
		have := ""
		i := new(Int)
		if err := i.Scan(v.in); err != nil {
			have = err.Error()
		}
		if want := v.err; have != want {
			t.Fatalf("test %d: wrong error, have '%s', want '%s'", tc, have, want)
		}
		if len(v.err) > 0 {
			continue
		}
		if !v.exp.Eq(i) {
			t.Fatalf("test %d: got %#x exp %#x", tc, i, v.exp)
		}
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

func BenchmarkScanScientific(b *testing.B) {
	intsub1 := new(Int)
	_ = intsub1.fromDecimal(twoPow256Sub1)
	cases := []struct {
		in  string
		exp *Int
		err string
	}{
		{
			in:  "14e30",
			exp: new(Int).Mul(NewInt(14), new(Int).Exp(NewInt(10), NewInt(30))),
		},
		{
			in:  "1455522523e31",
			exp: new(Int).Mul(NewInt(1455522523), new(Int).Exp(NewInt(10), NewInt(31))),
		},
		{
			in:  twoPow256Sub1 + "e0",
			exp: intsub1,
		},
		{
			in:  "1e00000000000000000",
			exp: NewInt(1),
		},
	}
	i := new(Int)
	b.ResetTimer()
	for idx := 0; idx < b.N; idx++ {
		for _, v := range cases {
			_ = i.Scan(v.in)
		}
	}
}

func benchSetFromBig(bench *testing.B, b *big.Int) Int {
	var f Int
	for i := 0; i < bench.N; i++ {
		f.SetFromBig(b)
	}
	return f
}

func BenchmarkSetFromBig(bench *testing.B) {
	param1 := big.NewInt(0xff)
	bench.Run("1word", func(bench *testing.B) { benchSetFromBig(bench, param1) })

	param2 := new(big.Int).Lsh(param1, 64)
	bench.Run("2words", func(bench *testing.B) { benchSetFromBig(bench, param2) })

	param3 := new(big.Int).Lsh(param2, 64)
	bench.Run("3words", func(bench *testing.B) { benchSetFromBig(bench, param3) })

	param4 := new(big.Int).Lsh(param3, 64)
	bench.Run("4words", func(bench *testing.B) { benchSetFromBig(bench, param4) })

	param5 := new(big.Int).Lsh(param4, 64)
	bench.Run("overflow", func(bench *testing.B) { benchSetFromBig(bench, param5) })
}

func benchToBig(bench *testing.B, f *Int) *big.Int {
	var b *big.Int
	for i := 0; i < bench.N; i++ {
		b = f.ToBig()
	}
	return b
}

func BenchmarkToBig(bench *testing.B) {
	param1 := new(Int).SetUint64(0xff)
	bench.Run("1word", func(bench *testing.B) { benchToBig(bench, param1) })

	param2 := new(Int).Lsh(param1, 64)
	bench.Run("2words", func(bench *testing.B) { benchToBig(bench, param2) })

	param3 := new(Int).Lsh(param2, 64)
	bench.Run("3words", func(bench *testing.B) { benchToBig(bench, param3) })

	param4 := new(Int).Lsh(param3, 64)
	bench.Run("4words", func(bench *testing.B) { benchToBig(bench, param4) })
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
		z := new(Int).SetAllOne().SetBytes(buf[0:i])
		if !z.Eq(exp) {
			t.Errorf("testcase %d: exp %x, got %x", i, exp, z)
		}
	}
	// nil check
	exp, _ := FromBig(new(big.Int).SetBytes(nil))
	z := new(Int).SetAllOne().SetBytes(nil)
	if !z.Eq(exp) {
		t.Errorf("nil-test : exp %x, got %x", exp, z)
	}
}

func BenchmarkSetBytes(b *testing.B) {

	val := new(Int)
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

func TestRlpEncode(t *testing.T) {

	type testcase struct {
		val string
		exp string
	}
	for i, tt := range []testcase{
		{"", "80"},
		{"01", "01"},
		{"02", "02"},
		{"04", "04"},
		{"08", "08"},
		{"10", "10"},
		{"20", "20"},
		{"40", "40"},
		{"80", "8180"},
		{"0100", "820100"},
		{"0200", "820200"},
		{"0400", "820400"},
		{"0800", "820800"},
		{"1000", "821000"},
		{"2000", "822000"},
		{"4000", "824000"},
		{"8000", "828000"},
		{"010000", "83010000"},
		{"020000", "83020000"},
		{"040000", "83040000"},
		{"080000", "83080000"},
		{"100000", "83100000"},
		{"200000", "83200000"},
		{"400000", "83400000"},
		{"800000", "83800000"},
		{"01000000", "8401000000"},
		{"02000000", "8402000000"},
		{"04000000", "8404000000"},
		{"08000000", "8408000000"},
		{"10000000", "8410000000"},
		{"20000000", "8420000000"},
		{"40000000", "8440000000"},
		{"80000000", "8480000000"},
		{"0100000000", "850100000000"},
		{"0200000000", "850200000000"},
		{"0400000000", "850400000000"},
		{"0800000000", "850800000000"},
		{"1000000000", "851000000000"},
		{"2000000000", "852000000000"},
		{"4000000000", "854000000000"},
		{"8000000000", "858000000000"},
		{"010000000000", "86010000000000"},
		{"020000000000", "86020000000000"},
		{"040000000000", "86040000000000"},
		{"080000000000", "86080000000000"},
		{"100000000000", "86100000000000"},
		{"200000000000", "86200000000000"},
		{"400000000000", "86400000000000"},
		{"800000000000", "86800000000000"},
		{"01000000000000", "8701000000000000"},
		{"02000000000000", "8702000000000000"},
		{"04000000000000", "8704000000000000"},
		{"08000000000000", "8708000000000000"},
		{"10000000000000", "8710000000000000"},
		{"20000000000000", "8720000000000000"},
		{"40000000000000", "8740000000000000"},
		{"80000000000000", "8780000000000000"},
		{"0100000000000000", "880100000000000000"},
		{"0200000000000000", "880200000000000000"},
		{"0400000000000000", "880400000000000000"},
		{"0800000000000000", "880800000000000000"},
		{"1000000000000000", "881000000000000000"},
		{"2000000000000000", "882000000000000000"},
		{"4000000000000000", "884000000000000000"},
		{"8000000000000000", "888000000000000000"},
		{"010000000000000000", "89010000000000000000"},
		{"020000000000000000", "89020000000000000000"},
		{"040000000000000000", "89040000000000000000"},
		{"080000000000000000", "89080000000000000000"},
		{"100000000000000000", "89100000000000000000"},
		{"200000000000000000", "89200000000000000000"},
		{"400000000000000000", "89400000000000000000"},
		{"800000000000000000", "89800000000000000000"},
		{"01000000000000000000", "8a01000000000000000000"},
		{"02000000000000000000", "8a02000000000000000000"},
		{"04000000000000000000", "8a04000000000000000000"},
		{"08000000000000000000", "8a08000000000000000000"},
		{"10000000000000000000", "8a10000000000000000000"},
		{"20000000000000000000", "8a20000000000000000000"},
		{"40000000000000000000", "8a40000000000000000000"},
		{"80000000000000000000", "8a80000000000000000000"},
		{"0100000000000000000000", "8b0100000000000000000000"},
		{"0200000000000000000000", "8b0200000000000000000000"},
		{"0400000000000000000000", "8b0400000000000000000000"},
		{"0800000000000000000000", "8b0800000000000000000000"},
		{"1000000000000000000000", "8b1000000000000000000000"},
		{"2000000000000000000000", "8b2000000000000000000000"},
		{"4000000000000000000000", "8b4000000000000000000000"},
		{"8000000000000000000000", "8b8000000000000000000000"},
		{"010000000000000000000000", "8c010000000000000000000000"},
		{"020000000000000000000000", "8c020000000000000000000000"},
		{"040000000000000000000000", "8c040000000000000000000000"},
		{"080000000000000000000000", "8c080000000000000000000000"},
		{"100000000000000000000000", "8c100000000000000000000000"},
		{"200000000000000000000000", "8c200000000000000000000000"},
		{"400000000000000000000000", "8c400000000000000000000000"},
		{"800000000000000000000000", "8c800000000000000000000000"},
		{"01000000000000000000000000", "8d01000000000000000000000000"},
		{"02000000000000000000000000", "8d02000000000000000000000000"},
		{"04000000000000000000000000", "8d04000000000000000000000000"},
		{"08000000000000000000000000", "8d08000000000000000000000000"},
		{"10000000000000000000000000", "8d10000000000000000000000000"},
		{"20000000000000000000000000", "8d20000000000000000000000000"},
		{"40000000000000000000000000", "8d40000000000000000000000000"},
		{"80000000000000000000000000", "8d80000000000000000000000000"},
		{"0100000000000000000000000000", "8e0100000000000000000000000000"},
		{"0200000000000000000000000000", "8e0200000000000000000000000000"},
		{"0400000000000000000000000000", "8e0400000000000000000000000000"},
		{"0800000000000000000000000000", "8e0800000000000000000000000000"},
		{"1000000000000000000000000000", "8e1000000000000000000000000000"},
		{"2000000000000000000000000000", "8e2000000000000000000000000000"},
		{"4000000000000000000000000000", "8e4000000000000000000000000000"},
		{"8000000000000000000000000000", "8e8000000000000000000000000000"},
		{"010000000000000000000000000000", "8f010000000000000000000000000000"},
		{"020000000000000000000000000000", "8f020000000000000000000000000000"},
		{"040000000000000000000000000000", "8f040000000000000000000000000000"},
		{"080000000000000000000000000000", "8f080000000000000000000000000000"},
		{"100000000000000000000000000000", "8f100000000000000000000000000000"},
		{"200000000000000000000000000000", "8f200000000000000000000000000000"},
		{"400000000000000000000000000000", "8f400000000000000000000000000000"},
		{"800000000000000000000000000000", "8f800000000000000000000000000000"},
		{"01000000000000000000000000000000", "9001000000000000000000000000000000"},
		{"02000000000000000000000000000000", "9002000000000000000000000000000000"},
		{"04000000000000000000000000000000", "9004000000000000000000000000000000"},
		{"08000000000000000000000000000000", "9008000000000000000000000000000000"},
		{"10000000000000000000000000000000", "9010000000000000000000000000000000"},
		{"20000000000000000000000000000000", "9020000000000000000000000000000000"},
		{"40000000000000000000000000000000", "9040000000000000000000000000000000"},
		{"80000000000000000000000000000000", "9080000000000000000000000000000000"},
		{"0100000000000000000000000000000000", "910100000000000000000000000000000000"},
		{"0200000000000000000000000000000000", "910200000000000000000000000000000000"},
		{"0400000000000000000000000000000000", "910400000000000000000000000000000000"},
		{"0800000000000000000000000000000000", "910800000000000000000000000000000000"},
		{"1000000000000000000000000000000000", "911000000000000000000000000000000000"},
		{"2000000000000000000000000000000000", "912000000000000000000000000000000000"},
		{"4000000000000000000000000000000000", "914000000000000000000000000000000000"},
		{"8000000000000000000000000000000000", "918000000000000000000000000000000000"},
		{"010000000000000000000000000000000000", "92010000000000000000000000000000000000"},
		{"020000000000000000000000000000000000", "92020000000000000000000000000000000000"},
		{"040000000000000000000000000000000000", "92040000000000000000000000000000000000"},
		{"080000000000000000000000000000000000", "92080000000000000000000000000000000000"},
		{"100000000000000000000000000000000000", "92100000000000000000000000000000000000"},
		{"200000000000000000000000000000000000", "92200000000000000000000000000000000000"},
		{"400000000000000000000000000000000000", "92400000000000000000000000000000000000"},
		{"800000000000000000000000000000000000", "92800000000000000000000000000000000000"},
		{"01000000000000000000000000000000000000", "9301000000000000000000000000000000000000"},
		{"02000000000000000000000000000000000000", "9302000000000000000000000000000000000000"},
		{"04000000000000000000000000000000000000", "9304000000000000000000000000000000000000"},
		{"08000000000000000000000000000000000000", "9308000000000000000000000000000000000000"},
		{"10000000000000000000000000000000000000", "9310000000000000000000000000000000000000"},
		{"20000000000000000000000000000000000000", "9320000000000000000000000000000000000000"},
		{"40000000000000000000000000000000000000", "9340000000000000000000000000000000000000"},
		{"80000000000000000000000000000000000000", "9380000000000000000000000000000000000000"},
		{"0100000000000000000000000000000000000000", "940100000000000000000000000000000000000000"},
		{"0200000000000000000000000000000000000000", "940200000000000000000000000000000000000000"},
		{"0400000000000000000000000000000000000000", "940400000000000000000000000000000000000000"},
		{"0800000000000000000000000000000000000000", "940800000000000000000000000000000000000000"},
		{"1000000000000000000000000000000000000000", "941000000000000000000000000000000000000000"},
		{"2000000000000000000000000000000000000000", "942000000000000000000000000000000000000000"},
		{"4000000000000000000000000000000000000000", "944000000000000000000000000000000000000000"},
		{"8000000000000000000000000000000000000000", "948000000000000000000000000000000000000000"},
		{"010000000000000000000000000000000000000000", "95010000000000000000000000000000000000000000"},
		{"020000000000000000000000000000000000000000", "95020000000000000000000000000000000000000000"},
		{"040000000000000000000000000000000000000000", "95040000000000000000000000000000000000000000"},
		{"080000000000000000000000000000000000000000", "95080000000000000000000000000000000000000000"},
		{"100000000000000000000000000000000000000000", "95100000000000000000000000000000000000000000"},
		{"200000000000000000000000000000000000000000", "95200000000000000000000000000000000000000000"},
		{"400000000000000000000000000000000000000000", "95400000000000000000000000000000000000000000"},
		{"800000000000000000000000000000000000000000", "95800000000000000000000000000000000000000000"},
		{"01000000000000000000000000000000000000000000", "9601000000000000000000000000000000000000000000"},
		{"02000000000000000000000000000000000000000000", "9602000000000000000000000000000000000000000000"},
		{"04000000000000000000000000000000000000000000", "9604000000000000000000000000000000000000000000"},
		{"08000000000000000000000000000000000000000000", "9608000000000000000000000000000000000000000000"},
		{"10000000000000000000000000000000000000000000", "9610000000000000000000000000000000000000000000"},
		{"20000000000000000000000000000000000000000000", "9620000000000000000000000000000000000000000000"},
		{"40000000000000000000000000000000000000000000", "9640000000000000000000000000000000000000000000"},
		{"80000000000000000000000000000000000000000000", "9680000000000000000000000000000000000000000000"},
		{"0100000000000000000000000000000000000000000000", "970100000000000000000000000000000000000000000000"},
		{"0200000000000000000000000000000000000000000000", "970200000000000000000000000000000000000000000000"},
		{"0400000000000000000000000000000000000000000000", "970400000000000000000000000000000000000000000000"},
		{"0800000000000000000000000000000000000000000000", "970800000000000000000000000000000000000000000000"},
		{"1000000000000000000000000000000000000000000000", "971000000000000000000000000000000000000000000000"},
		{"2000000000000000000000000000000000000000000000", "972000000000000000000000000000000000000000000000"},
		{"4000000000000000000000000000000000000000000000", "974000000000000000000000000000000000000000000000"},
		{"8000000000000000000000000000000000000000000000", "978000000000000000000000000000000000000000000000"},
		{"010000000000000000000000000000000000000000000000", "98010000000000000000000000000000000000000000000000"},
		{"020000000000000000000000000000000000000000000000", "98020000000000000000000000000000000000000000000000"},
		{"040000000000000000000000000000000000000000000000", "98040000000000000000000000000000000000000000000000"},
		{"080000000000000000000000000000000000000000000000", "98080000000000000000000000000000000000000000000000"},
		{"100000000000000000000000000000000000000000000000", "98100000000000000000000000000000000000000000000000"},
		{"200000000000000000000000000000000000000000000000", "98200000000000000000000000000000000000000000000000"},
		{"400000000000000000000000000000000000000000000000", "98400000000000000000000000000000000000000000000000"},
		{"800000000000000000000000000000000000000000000000", "98800000000000000000000000000000000000000000000000"},
		{"01000000000000000000000000000000000000000000000000", "9901000000000000000000000000000000000000000000000000"},
		{"02000000000000000000000000000000000000000000000000", "9902000000000000000000000000000000000000000000000000"},
		{"04000000000000000000000000000000000000000000000000", "9904000000000000000000000000000000000000000000000000"},
		{"08000000000000000000000000000000000000000000000000", "9908000000000000000000000000000000000000000000000000"},
		{"10000000000000000000000000000000000000000000000000", "9910000000000000000000000000000000000000000000000000"},
		{"20000000000000000000000000000000000000000000000000", "9920000000000000000000000000000000000000000000000000"},
		{"40000000000000000000000000000000000000000000000000", "9940000000000000000000000000000000000000000000000000"},
		{"80000000000000000000000000000000000000000000000000", "9980000000000000000000000000000000000000000000000000"},
		{"0100000000000000000000000000000000000000000000000000", "9a0100000000000000000000000000000000000000000000000000"},
		{"0200000000000000000000000000000000000000000000000000", "9a0200000000000000000000000000000000000000000000000000"},
		{"0400000000000000000000000000000000000000000000000000", "9a0400000000000000000000000000000000000000000000000000"},
		{"0800000000000000000000000000000000000000000000000000", "9a0800000000000000000000000000000000000000000000000000"},
		{"1000000000000000000000000000000000000000000000000000", "9a1000000000000000000000000000000000000000000000000000"},
		{"2000000000000000000000000000000000000000000000000000", "9a2000000000000000000000000000000000000000000000000000"},
		{"4000000000000000000000000000000000000000000000000000", "9a4000000000000000000000000000000000000000000000000000"},
		{"8000000000000000000000000000000000000000000000000000", "9a8000000000000000000000000000000000000000000000000000"},
		{"010000000000000000000000000000000000000000000000000000", "9b010000000000000000000000000000000000000000000000000000"},
		{"020000000000000000000000000000000000000000000000000000", "9b020000000000000000000000000000000000000000000000000000"},
		{"040000000000000000000000000000000000000000000000000000", "9b040000000000000000000000000000000000000000000000000000"},
		{"080000000000000000000000000000000000000000000000000000", "9b080000000000000000000000000000000000000000000000000000"},
		{"100000000000000000000000000000000000000000000000000000", "9b100000000000000000000000000000000000000000000000000000"},
		{"200000000000000000000000000000000000000000000000000000", "9b200000000000000000000000000000000000000000000000000000"},
		{"400000000000000000000000000000000000000000000000000000", "9b400000000000000000000000000000000000000000000000000000"},
		{"800000000000000000000000000000000000000000000000000000", "9b800000000000000000000000000000000000000000000000000000"},
		{"01000000000000000000000000000000000000000000000000000000", "9c01000000000000000000000000000000000000000000000000000000"},
		{"02000000000000000000000000000000000000000000000000000000", "9c02000000000000000000000000000000000000000000000000000000"},
		{"04000000000000000000000000000000000000000000000000000000", "9c04000000000000000000000000000000000000000000000000000000"},
		{"08000000000000000000000000000000000000000000000000000000", "9c08000000000000000000000000000000000000000000000000000000"},
		{"10000000000000000000000000000000000000000000000000000000", "9c10000000000000000000000000000000000000000000000000000000"},
		{"20000000000000000000000000000000000000000000000000000000", "9c20000000000000000000000000000000000000000000000000000000"},
		{"40000000000000000000000000000000000000000000000000000000", "9c40000000000000000000000000000000000000000000000000000000"},
		{"80000000000000000000000000000000000000000000000000000000", "9c80000000000000000000000000000000000000000000000000000000"},
		{"0100000000000000000000000000000000000000000000000000000000", "9d0100000000000000000000000000000000000000000000000000000000"},
		{"0200000000000000000000000000000000000000000000000000000000", "9d0200000000000000000000000000000000000000000000000000000000"},
		{"0400000000000000000000000000000000000000000000000000000000", "9d0400000000000000000000000000000000000000000000000000000000"},
		{"0800000000000000000000000000000000000000000000000000000000", "9d0800000000000000000000000000000000000000000000000000000000"},
		{"1000000000000000000000000000000000000000000000000000000000", "9d1000000000000000000000000000000000000000000000000000000000"},
		{"2000000000000000000000000000000000000000000000000000000000", "9d2000000000000000000000000000000000000000000000000000000000"},
		{"4000000000000000000000000000000000000000000000000000000000", "9d4000000000000000000000000000000000000000000000000000000000"},
		{"8000000000000000000000000000000000000000000000000000000000", "9d8000000000000000000000000000000000000000000000000000000000"},
		{"010000000000000000000000000000000000000000000000000000000000", "9e010000000000000000000000000000000000000000000000000000000000"},
		{"020000000000000000000000000000000000000000000000000000000000", "9e020000000000000000000000000000000000000000000000000000000000"},
		{"040000000000000000000000000000000000000000000000000000000000", "9e040000000000000000000000000000000000000000000000000000000000"},
		{"080000000000000000000000000000000000000000000000000000000000", "9e080000000000000000000000000000000000000000000000000000000000"},
		{"100000000000000000000000000000000000000000000000000000000000", "9e100000000000000000000000000000000000000000000000000000000000"},
		{"200000000000000000000000000000000000000000000000000000000000", "9e200000000000000000000000000000000000000000000000000000000000"},
		{"400000000000000000000000000000000000000000000000000000000000", "9e400000000000000000000000000000000000000000000000000000000000"},
		{"800000000000000000000000000000000000000000000000000000000000", "9e800000000000000000000000000000000000000000000000000000000000"},
		{"01000000000000000000000000000000000000000000000000000000000000", "9f01000000000000000000000000000000000000000000000000000000000000"},
		{"02000000000000000000000000000000000000000000000000000000000000", "9f02000000000000000000000000000000000000000000000000000000000000"},
		{"04000000000000000000000000000000000000000000000000000000000000", "9f04000000000000000000000000000000000000000000000000000000000000"},
		{"08000000000000000000000000000000000000000000000000000000000000", "9f08000000000000000000000000000000000000000000000000000000000000"},
		{"10000000000000000000000000000000000000000000000000000000000000", "9f10000000000000000000000000000000000000000000000000000000000000"},
		{"20000000000000000000000000000000000000000000000000000000000000", "9f20000000000000000000000000000000000000000000000000000000000000"},
		{"40000000000000000000000000000000000000000000000000000000000000", "9f40000000000000000000000000000000000000000000000000000000000000"},
		{"80000000000000000000000000000000000000000000000000000000000000", "9f80000000000000000000000000000000000000000000000000000000000000"},
		{"0100000000000000000000000000000000000000000000000000000000000000", "a00100000000000000000000000000000000000000000000000000000000000000"},
		{"0200000000000000000000000000000000000000000000000000000000000000", "a00200000000000000000000000000000000000000000000000000000000000000"},
		{"0400000000000000000000000000000000000000000000000000000000000000", "a00400000000000000000000000000000000000000000000000000000000000000"},
		{"0800000000000000000000000000000000000000000000000000000000000000", "a00800000000000000000000000000000000000000000000000000000000000000"},
		{"1000000000000000000000000000000000000000000000000000000000000000", "a01000000000000000000000000000000000000000000000000000000000000000"},
		{"2000000000000000000000000000000000000000000000000000000000000000", "a02000000000000000000000000000000000000000000000000000000000000000"},
		{"4000000000000000000000000000000000000000000000000000000000000000", "a04000000000000000000000000000000000000000000000000000000000000000"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "a08000000000000000000000000000000000000000000000000000000000000000"},
	} {
		z := new(Int).SetBytes(hex2Bytes(tt.val))
		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		if err := z.EncodeRLP(w); err != nil {
			t.Fatal(err)
		}
		w.Flush()
		if got, exp := b.Bytes(), hex2Bytes(tt.exp); !bytes.Equal(got, exp) {
			t.Fatalf("testcase %d got:\n%x\nexp:%x\n", i, got, exp)
		}
	}
	// And test nil
	{
		var z *Int
		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		if err := z.EncodeRLP(w); err != nil {
			t.Fatal(err)
		}
		w.Flush()
		if got, exp := b.Bytes(), hex2Bytes("80"); !bytes.Equal(got, exp) {
			t.Fatalf("nil-test got:\n%x\nexp:%x\n", got, exp)
		}
	}
}

type nilWriter struct{}

func (*nilWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// BenchmarkRLPEncoding writes 255 Ints ranging in bitsize from 0-255 in each op
func BenchmarkRLPEncoding(b *testing.B) {
	z := new(Int)
	devnull := &nilWriter{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		z.SetUint64(1)
		for bit := 0; bit < 255; bit++ {
			_ = z.EncodeRLP(devnull)
			z.Lsh(z, 1)
		}
	}
}

func referenceBig(s string) *big.Int {
	b, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic("invalid")
	}
	return b
}

type marshalTest struct {
	input interface{}
	want  string
}

type unmarshalTest struct {
	input   string
	want    interface{}
	wantErr error // if set, decoding must fail on any platform
}

var (
	encodeBigTests = []marshalTest{
		{referenceBig("0"), "0x0"},
		{referenceBig("1"), "0x1"},
		{referenceBig("ff"), "0xff"},
		{referenceBig("112233445566778899aabbccddeeff"), "0x112233445566778899aabbccddeeff"},
		{referenceBig("80a7f2c1bcc396c00"), "0x80a7f2c1bcc396c00"},
	}

	decodeBigTests = []unmarshalTest{
		// invalid
		{input: ``, wantErr: ErrEmptyString},
		{input: `0`, wantErr: ErrMissingPrefix},
		{input: `0x`, wantErr: ErrEmptyNumber},
		{input: `0x01`, wantErr: ErrLeadingZero},
		{input: `0xx`, wantErr: ErrSyntax},
		{input: `0x1zz01`, wantErr: ErrSyntax},
		{
			input:   `0x10000000000000000000000000000000000000000000000000000000000000000`,
			wantErr: ErrBig256Range,
		},
		// valid
		{input: `0x0`, want: big.NewInt(0)},
		{input: `0x2`, want: big.NewInt(0x2)},
		{input: `0x2F2`, want: big.NewInt(0x2f2)},
		{input: `0X2F2`, want: big.NewInt(0x2f2)},
		{input: `0x1122aaff`, want: big.NewInt(0x1122aaff)},
		{input: `0xbBb`, want: big.NewInt(0xbbb)},
		{input: `0xfffffffff`, want: big.NewInt(0xfffffffff)},
		{
			input: `0x112233445566778899aabbccddeeff`,
			want:  referenceBig("112233445566778899aabbccddeeff"),
		},
		{
			input: `0xffffffffffffffffffffffffffffffffffff`,
			want:  referenceBig("ffffffffffffffffffffffffffffffffffff"),
		},
		{
			input: `0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff`,
			want:  referenceBig("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		},
	}
)

func checkError(t *testing.T, input string, got, want error) bool {
	if got == nil {
		if want != nil {
			t.Errorf("input %s: got no error, want %q", input, want)
			return false
		}
		return true
	}
	if want == nil {
		t.Errorf("input %s: unexpected error %q", input, got)
	} else if got.Error() != want.Error() {
		t.Errorf("input %s: got error %q, want %q", input, got, want)
	}
	return false
}

func TestEncode(t *testing.T) {
	for _, test := range encodeBigTests {
		z, _ := FromBig(test.input.(*big.Int))
		enc := z.Hex()
		if enc != test.want {
			t.Errorf("input %x: wrong encoding %s (exp %s)", test.input, enc, test.want)
		}
	}

}

func TestDecode(t *testing.T) {
	for _, test := range decodeBigTests {
		dec, err := FromHex(test.input)
		if !checkError(t, test.input, err, test.wantErr) {
			continue
		}
		b := dec.ToBig()
		if b.Cmp(test.want.(*big.Int)) != 0 {
			t.Errorf("input %s: value mismatch: got %x, want %x", test.input, dec, test.want)
			continue
		}
	}
	// Some remaining json-tests
	type jsonStruct struct {
		Foo *Int
	}
	var jsonDecoded jsonStruct
	if err := json.Unmarshal([]byte(`{"Foo":0x1}`), &jsonDecoded); err == nil {
		t.Fatal("Expected error")
	}
	if err := json.Unmarshal([]byte(`{"Foo":1}`), &jsonDecoded); err == nil {
		t.Fatal("Expected error")
	}
	if err := json.Unmarshal([]byte(`{"Foo":""}`), &jsonDecoded); err == nil {
		t.Fatal("Expected error")
	}
	if err := json.Unmarshal([]byte(`{"Foo":"0x1"}`), &jsonDecoded); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	} else if jsonDecoded.Foo.Uint64() != 1 {
		t.Fatal("Expected 1")
	}
}

func TestEnDecode(t *testing.T) {
	type jsonStruct struct {
		Foo *Int
	}
	var testSample = func(i int, bigSample big.Int, intSample Int) {
		// Encoding
		wantHex := fmt.Sprintf("0x%s", bigSample.Text(16))
		wantDec := bigSample.Text(10)

		if got := intSample.Hex(); wantHex != got {
			t.Fatalf("test %d #1, got %v, exp %v", i, got, wantHex)
		}
		if got := intSample.String(); wantHex != got {
			t.Fatalf("test %d #2, got %v, exp %v", i, got, wantHex)
		}
		if got, _ := intSample.MarshalText(); wantHex != string(got) {
			t.Fatalf("test %d #3, got %v, exp %v", i, got, wantHex)
		}
		if got, _ := intSample.Value(); wantDec != got.(string) {
			t.Fatalf("test %d #4, got %v, exp %v", i, got, wantHex)
		}
		{ // Json
			jsonEncoded, err := json.Marshal(&jsonStruct{&intSample})
			if err != nil {
				t.Fatalf("test %d #4, err: %v", i, err)
			}
			var jsonDecoded jsonStruct
			err = json.Unmarshal(jsonEncoded, &jsonDecoded)
			if err != nil {
				t.Fatalf("test %d #5, err: %v", i, err)
			}
			if jsonDecoded.Foo.Cmp(&intSample) != 0 {
				t.Fatalf("test %d #6, got %v, exp %v", i, jsonDecoded.Foo, intSample)
			}
		}
		// Decoding
		//
		// FromHex
		decoded, err := FromHex(wantHex)
		{
			if err != nil {
				t.Fatalf("test %d #5, err: %v", i, err)
			}
			if decoded.Cmp(&intSample) != 0 {
				t.Fatalf("test %d #6, got %v, exp %v", i, decoded, intSample)
			}
		}
		// z.SetFromHex
		err = decoded.SetFromHex(wantHex)
		{
			if err != nil {
				t.Fatalf("test %d #5, err: %v", i, err)
			}
			if decoded.Cmp(&intSample) != 0 {
				t.Fatalf("test %d #6, got %v, exp %v", i, decoded, intSample)
			}
		}
		// UnmarshalText
		decoded = new(Int)
		{
			if err := decoded.UnmarshalText([]byte(wantHex)); err != nil {
				t.Fatalf("test %d #7, err: %v", i, err)
			}
			if decoded.Cmp(&intSample) != 0 {
				t.Fatalf("test %d #8, got %v, exp %v", i, decoded, intSample)
			}
		}
		// FromDecimal
		decoded, err = FromDecimal(wantDec)
		{
			if err != nil {
				t.Fatalf("test %d #9, err: %v", i, err)
			}
			if decoded.Cmp(&intSample) != 0 {
				t.Fatalf("test %d #10, got %v, exp %v", i, decoded, intSample)
			}
		}
		// Scan w string
		err = decoded.Scan(wantDec)
		{
			if err != nil {
				t.Fatalf("test %d #9, err: %v", i, err)
			}
			if decoded.Cmp(&intSample) != 0 {
				t.Fatalf("test %d #10, got %v, exp %v", i, decoded, intSample)
			}
		}
		// Scan w byte slice
		err = decoded.Scan([]byte(wantDec))
		{
			if err != nil {
				t.Fatalf("test %d #9, err: %v", i, err)
			}
			if decoded.Cmp(&intSample) != 0 {
				t.Fatalf("test %d #10, got %v, exp %v", i, decoded, intSample)
			}
		}
		// Scan with neither string nor byte
		err = decoded.Scan(5)
		{
			if err == nil {
				t.Fatalf("test %d #11, want error", i)
			}
		}
	}
	for i, bigSample := range big256Samples {
		intSample := int256Samples[i]
		testSample(i, bigSample, intSample)
	}

	for i, bigSample := range big256SamplesLt {
		intSample := int256SamplesLt[i]
		testSample(i, bigSample, intSample)
	}
}

func TestNil(t *testing.T) {
	a := NewInt(1337)
	if err := a.Scan(nil); err != nil {
		t.Fatal(err)
	}
	if !a.IsZero() {
		t.Fatal("want zero")
	}
}
