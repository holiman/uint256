package uint256

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"testing"
)

func TestStringScanBase10(t *testing.T) {
	z := new(Int)
	type testCase struct {
		input string
		err   error
	}
	for i, tc := range []testCase{
		{input: twoPow256Sub1 + "1", err: ErrBig256Range},
		{input: "2" + twoPow256Sub1[1:], err: ErrBig256Range},
		{input: twoPow256Sub1[1:]},
		{input: twoPow128},
		{input: twoPow128 + "1"},
		{input: twoPow128[1:]},
		{input: twoPow64 + "1"},
		{input: twoPow64[1:]},
		{input: "banana", err: strconv.ErrSyntax},
		{input: "0xab", err: strconv.ErrSyntax},
		{input: "ab", err: strconv.ErrSyntax},
		{input: "0"},
		{input: "000"},
		{input: "010"},
		{input: "01"},
		{input: "-0", err: strconv.ErrSyntax},
		{input: "-10", err: strconv.ErrSyntax},
		{input: "115792089237316195423570985008687907853269984665640564039457584007913129639936", err: ErrBig256Range},
		{input: "115792089237316195423570985008687907853269984665640564039457584007913129639935"},
	} {
		z.SetAllOne() // Set to ensure all bits are cleared after
		err := z.SetFromBase10(tc.input)
		if !errors.Is(err, tc.err) {
			t.Errorf("test %d, input %v: want err %s, have %s", i, tc.input, tc.err, err)
		}
		if err != nil {
			continue
		}
		var want string
		if w, ok := big.NewInt(0).SetString(tc.input, 10); !ok {
			panic(fmt.Sprintf("test %d error", i))
		} else {
			want = w.String()
		}
		if have := z.ToBig().String(); have != want {
			t.Errorf("test %d: input %v,  want %v: have %s", i, tc.input, want, have)
		}
	}
}

func FuzzBase10StringCompare(f *testing.F) {
	var (
		bi        = new(big.Int)
		z         = new(Int)
		max256, _ = FromBase10(twoPow256Sub1)
	)
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
		"-10",
		"115792089237316195423570985008687907853269984665640564039457584007913129639936",
		"115792089237316195423570985008687907853269984665640564039457584007913129639935",
		"apple",
		"04112401274120741204712xxxxxz00",
		"0x10101011010",
		"熊熊熊熊熊熊熊熊",
	} {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, orig string) {
		err := z.SetFromBase10(orig)
		val, ok := bi.SetString(orig, 10)
		// if fail, make sure that we failed too
		if !ok {
			if err == nil {
				t.Errorf("expected base 10 parse to fail: %s", orig)
			}
			return
		}
		// if its negative number, we should err
		if len(orig) > 0 && (orig[0] == '-') {
			if !errors.Is(err, strconv.ErrSyntax) {
				t.Errorf("should have errored at negative number: %s", orig)
			}
			return
		}
		// if its too large, ignore it also
		if val.Cmp(max256.ToBig()) > 0 {
			if !errors.Is(err, ErrBig256Range) {
				t.Errorf("should have errored at number overflow: %s", orig)
			}
			return
		}
		// so here, if it errors, it means that we failed
		if err != nil {
			t.Errorf("should have parsed %s to %s, but err'd instead", orig, val.String())
			return
		}
		// otherwise, make sure that the values are equal
		if z.ToBig().Cmp(bi) != 0 {
			t.Errorf("should have parsed %s to %s, but got %s", orig, bi.String(), z.Base10())
			return
		}
		// make sure that bigint base 10 string is equal to base10 string
		if z.Base10() != bi.String() {
			t.Errorf("should have parsed %s to %s, but got %s", orig, bi.String(), z.Base10())
			return
		}
		value, err := z.Value()
		if err != nil {
			t.Errorf("fail to Value() %s, got err %s", val, err)
			return
		}
		if z.Base10()+"e0" != fmt.Sprint(value) {
			t.Errorf("value of %s did not match base 10 encoding %s", value, z.Base10())
			return
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
