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

// Test SetFromDecimal
func testSetFromDecimal(tc string) error {
	a := new(Int).SetAllOne()
	err := a.SetFromDecimal(tc)
	haveOk := (err == nil)
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
		//if err := testSetFromBase0(tc); err != nil {
		//	t.Errorf("test %d, input '%s', SetString(..,0) err: %v", i, tc, err)
		//}
		//if err := testSetFromBase10(tc); err != nil {
		//	t.Errorf("test %d, input '%s', SetString(..,0) err: %v", i, tc, err)
		//}
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
		//if err := testSetFromBase10(tc); err != nil {
		//	t.Errorf("input '%s', SetString(..,0) err: %v", tc, err)
		//}
	})
}

func BenchmarkStringBase10(b *testing.B) {
	input := twoPow256Sub1

	b.Run("bigint", func(b *testing.B) {
		val := new(big.Int)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for j := 1; j < len(input); j++ {
				val.SetString(input[:j], 10)
			}
		}
	})

	b.Run("u256", func(b *testing.B) {
		val := new(Int)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for j := 1; j < len(input); j++ {
				val.SetFromDecimal(input[:j])
			}
		}
	})
}

func BenchmarkStringBase16(b *testing.B) {
	input := "aaaa12131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f3031bbbb"
	b.Run("bigint", func(b *testing.B) {
		val := new(big.Int)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for j := 1; j < len(input); j++ {
				val.SetString(input[:j], 16)
			}
		}
	})
	b.Run("u256", func(b *testing.B) {
		val := new(Int)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for j := 1; j < len(input); j++ {
				val.SetFromHex(input[:j])
			}
		}
	})
}
