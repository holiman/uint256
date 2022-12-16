// uint256: Fixed size 256-bit math library
// Copyright 2020 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

package uint256

import (
	"fmt"
	"math/big"
	"testing"
)

const (
	twoPow64  = "18446744073709551616"
	twoPow128 = "340282366920938463463374607431768211456"
)

// Test SetFromDecimal
func testSetFromDec(tc string) error {
	a := new(Int).SetAllOne()
	err := a.SetFromDecimal(tc)
	// If input is negative, we should eror
	if len(tc) > 0 && tc[0] == '-' {
		if err == nil {
			return fmt.Errorf("want error on negative input")
		}
		return nil
	}
	// Need to compare with big.Int
	bigA, ok := big.NewInt(0).SetString(tc, 10)
	if !ok {
		if err == nil {
			return fmt.Errorf("want error")
		}
		return nil // both agree that input is bad
	}
	if bigA.BitLen() > 256 {
		if err == nil {
			return fmt.Errorf("want error (bitlen > 256)")
		}
		return nil
	}
	want := bigA.String()
	have := a.Dec()
	if want != have {
		return fmt.Errorf("want %v, have %v", want, have)
	}
	return nil
}

// Test SetString base 0
func testSetFromBase0(tc string) error {
	a := new(Int).SetAllOne()
	a, haveOk := a.SetString(tc, 0)
	// If input is negative, we should eror
	if len(tc) > 0 && tc[0] == '-' {
		if haveOk {
			return fmt.Errorf("want error on negative input")
		}
		return nil
	}
	// Need to compare with big.Int
	bigA, ok := big.NewInt(0).SetString(tc, 0)
	if !ok {
		if haveOk {
			return fmt.Errorf("want error")
		}
		return nil // both agree that input is bad
	}
	if bigA.BitLen() > 256 {
		if haveOk {
			return fmt.Errorf("want error (bitlen > 256)")
		}
		return nil
	}
	if !haveOk {
		return fmt.Errorf("want no err, have err")
	}
	want := bigA.String()
	have := a.Dec()
	if want != have {
		return fmt.Errorf("want %v, have %v", want, have)
	}
	return nil
}

// Test SetString base 10
func testSetFromBase10(tc string) error {
	a := new(Int).SetAllOne()
	a, haveOk := a.SetString(tc, 10)
	// If input is negative, we should eror
	if len(tc) > 0 && tc[0] == '-' {
		if haveOk {
			return fmt.Errorf("want error on negative input")
		}
		return nil
	}
	// Need to compare with big.Int
	bigA, ok := big.NewInt(0).SetString(tc, 10)
	if !ok {
		if haveOk {
			return fmt.Errorf("want error")
		}
		return nil // both agree that input is bad
	}
	if bigA.BitLen() > 256 {
		if haveOk {
			return fmt.Errorf("want error (bitlen > 256)")
		}
		return nil
	}
	if !haveOk {
		return fmt.Errorf("want no err, have err")
	}
	want := bigA.String()
	have := a.Dec()
	if want != have {
		return fmt.Errorf("want %v, have %v", want, have)
	}
	return nil
}

func TestStringScan(t *testing.T) {
	for i, tc := range []string{
		"0000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"0000000000000000000000000000000000000000000000000000000000000000000000000000097",
		"-000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"1157920892373161954235709850086879078532699846656405640394575840079131296399351",
		"215792089237316195423570985008687907853269984665640564039457584007913129639935",
		"115792089237316195423570985008687907853269984665640564039457584007913129639935",
		"15792089237316195423570985008687907853269984665640564039457584007913129639935",
		"+115792089237316195423570985008687907853269984665640564039457584007913129639935",
		"115792089237316195423570985008687907853269984665640564039457584007913129639936",
		"115792089237316195423570985008687907853269984665640564039457584007913129639935",
		"+0b00000000000000000000000000000000000000000000000000000000000000010",
		"340282366920938463463374607431768211456",
		"3402823669209384634633746074317682114561",
		"+3402823669209384634633746074317682114561",
		"+-3402823669209384634633746074317682114561",
		"40282366920938463463374607431768211456",
		"00000000000000000000000097",
		"184467440737095516161",
		"8446744073709551616",
		"banana",
		"+0x10",
		"000",
		"+000",
		"010",
		"0xab",
		"-10",
		"01",
		"ab",
		"0",
		"-0",
		"+0",
		"",
		"熊熊熊熊熊熊熊熊",
		"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"-0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
	} {
		if err := testSetFromDec(tc); err != nil {
			t.Errorf("test %d, input '%s', SetFromDecimal err: %v", i, tc, err)
		}
		if err := testSetFromBase0(tc); err != nil {
			t.Errorf("test %d, input '%s', SetString(..,0) err: %v", i, tc, err)
		}
		if err := testSetFromBase10(tc); err != nil {
			t.Errorf("test %d, input '%s', SetString(..,0) err: %v", i, tc, err)
		}
		// TODO test SetString(.., 16)
	}
}

func FuzzBase10StringCompare(f *testing.F) {

	for _, tc := range []string{
		twoPow256Sub1 + "1",
		"2" + twoPow256Sub1[1:],
		twoPow256Sub1,
		twoPow128,
		twoPow128,
		twoPow128,
		twoPow64,
		twoPow64,
		"banana",
		"0xab",
		"ab",
		"0",
		"000",
		"010",
		"01",
		"-0",
		"+0",
		"-10",
		"115792089237316195423570985008687907853269984665640564039457584007913129639936",
		"115792089237316195423570985008687907853269984665640564039457584007913129639935",
		"apple",
		"04112401274120741204712xxxxxz00",
		"0x10101011010",
		"熊熊熊熊熊熊熊熊",
		"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"-0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"-0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"0x10000000000000000000000000000000000000000000000000000000000000000",
		"+0x10000000000000000000000000000000000000000000000000000000000000000",
		"+0x00000000000000000000000000000000000000000000000000000000000000000",
		"-0x00000000000000000000000000000000000000000000000000000000000000000",
	} {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, tc string) {
		if err := testSetFromDec(tc); err != nil {
			t.Errorf("input '%s', SetFromDecimal err: %v", tc, err)
		}
		// Our base16 parsing differs, so inputs like this would be rejected
		// where big.Int accepts them:
		// +0x00000000000000000000000000000000000000000000000000000000000000000
		//if err := testSetFromBase0(tc); err != nil {
		//	t.Errorf("input '%s', SetString(..,0) err: %v", tc, err)
		//}
		if err := testSetFromBase10(tc); err != nil {
			t.Errorf("input '%s', SetString(..,0) err: %v", tc, err)
		}
	})
}

func BenchmarkStringBase10BigInt(b *testing.B) {
	val := new(big.Int)
	bytearr := twoPow256Sub1
	b.Run("generic", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			val.SetString(bytearr[:2], 10)
			val.SetString(bytearr[:4], 10)
			val.SetString(bytearr[:6], 10)
			val.SetString(bytearr[:8], 10)
			val.SetString(bytearr[:16], 10)
			val.SetString(bytearr[:12], 10)
			val.SetString(bytearr[:14], 10)
			val.SetString(bytearr[:16], 10)
			val.SetString(bytearr[:18], 10)
			val.SetString(bytearr[:20], 10)
			val.SetString(bytearr[:22], 10)
			val.SetString(bytearr[:24], 10)
			val.SetString(bytearr[:26], 10)
			val.SetString(bytearr[:28], 10)
			val.SetString(bytearr[:30], 10)
			val.SetString(bytearr[:32], 10)
			val.SetString(bytearr[:34], 10)
			val.SetString(bytearr[:36], 10)
			val.SetString(bytearr[:38], 10)
			val.SetString(bytearr[:40], 10)
			val.SetString(bytearr[:42], 10)
			val.SetString(bytearr[:44], 10)
			val.SetString(bytearr[:46], 10)
			val.SetString(bytearr[:48], 10)
			val.SetString(bytearr[:50], 10)
			val.SetString(bytearr[:52], 10)
			val.SetString(bytearr[:54], 10)
			val.SetString(bytearr[:56], 10)
			val.SetString(bytearr[:58], 10)
			val.SetString(bytearr[:60], 10)
			val.SetString(bytearr[:62], 10)
			val.SetString(bytearr[:64], 10)
		}
	})
}

func BenchmarkStringBase10(b *testing.B) {
	val := new(Int)
	bytearr := twoPow256Sub1
	b.Run("generic", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			val.SetString(bytearr[:2], 10)
			val.SetString(bytearr[:4], 10)
			val.SetString(bytearr[:6], 10)
			val.SetString(bytearr[:8], 10)
			val.SetString(bytearr[:16], 10)
			val.SetString(bytearr[:12], 10)
			val.SetString(bytearr[:14], 10)
			val.SetString(bytearr[:16], 10)
			val.SetString(bytearr[:18], 10)
			val.SetString(bytearr[:20], 10)
			val.SetString(bytearr[:22], 10)
			val.SetString(bytearr[:24], 10)
			val.SetString(bytearr[:26], 10)
			val.SetString(bytearr[:28], 10)
			val.SetString(bytearr[:30], 10)
			val.SetString(bytearr[:32], 10)
			val.SetString(bytearr[:34], 10)
			val.SetString(bytearr[:36], 10)
			val.SetString(bytearr[:38], 10)
			val.SetString(bytearr[:40], 10)
			val.SetString(bytearr[:42], 10)
			val.SetString(bytearr[:44], 10)
			val.SetString(bytearr[:46], 10)
			val.SetString(bytearr[:48], 10)
			val.SetString(bytearr[:50], 10)
			val.SetString(bytearr[:52], 10)
			val.SetString(bytearr[:54], 10)
			val.SetString(bytearr[:56], 10)
			val.SetString(bytearr[:58], 10)
			val.SetString(bytearr[:60], 10)
			val.SetString(bytearr[:62], 10)
			val.SetString(bytearr[:64], 10)
		}
	})
}
func BenchmarkStringBase16(b *testing.B) {
	val := new(Int)
	bytearr := "aaaa12131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f3031bbbb"
	b.Run("generic", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			val.SetString(bytearr[:1], 16)
			val.SetString(bytearr[:2], 16)
			val.SetString(bytearr[:3], 16)
			val.SetString(bytearr[:4], 16)
			val.SetString(bytearr[:5], 16)
			val.SetString(bytearr[:6], 16)
			val.SetString(bytearr[:7], 16)
			val.SetString(bytearr[:8], 16)
			val.SetString(bytearr[:9], 16)
			val.SetString(bytearr[:10], 16)
			val.SetString(bytearr[:11], 16)
			val.SetString(bytearr[:12], 16)
			val.SetString(bytearr[:13], 16)
			val.SetString(bytearr[:14], 16)
			val.SetString(bytearr[:15], 16)
			val.SetString(bytearr[:16], 16)
			val.SetString(bytearr[:17], 16)
			val.SetString(bytearr[:18], 16)
			val.SetString(bytearr[:19], 16)
			val.SetString(bytearr[:20], 16)
			val.SetString(bytearr[:21], 16)
			val.SetString(bytearr[:22], 16)
			val.SetString(bytearr[:23], 16)
			val.SetString(bytearr[:24], 16)
			val.SetString(bytearr[:25], 16)
			val.SetString(bytearr[:26], 16)
			val.SetString(bytearr[:27], 16)
			val.SetString(bytearr[:28], 16)
			val.SetString(bytearr[:29], 16)
			val.SetString(bytearr[:20], 16)
			val.SetString(bytearr[:31], 16)
			val.SetString(bytearr[:32], 16)
		}
	})
}

func TestFoo(t *testing.T) {
	s := "+0b00000000000000000000000000000000000000000000000000000000000000010"
	b, ok := new(big.Int).SetString(s, 0)
	fmt.Printf("b: %v, ok : %v\n", b, ok)
	z, ok := new(Int).SetString(s, 0)
	fmt.Printf("z: %v, ok : %v\n", z, ok)
}
