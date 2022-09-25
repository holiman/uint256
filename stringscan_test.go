package uint256

import (
	"errors"
	"testing"
)

func TestStringScanBase10(t *testing.T) {
	z := new(Int)

	type testCase struct {
		i   string
		err error
	}

	cases := []testCase{
		{i: twoPow256 + "1", err: ErrBig256Range},
		{i: "2" + twoPow256[1:], err: ErrBig256Range},
		{i: twoPow256[1:]},
		{i: twoPow128},
		{i: twoPow128 + "1"},
		{i: twoPow128[1:]},
		{i: twoPow64 + "1"},
		{i: twoPow64[1:]},
	}

	for _, v := range cases {
		err := z.FromBase10(v.i)
		if !errors.Is(err, v.err) {
			t.Errorf("expect err %s, got %s", v.err, err)
		}
		if err == nil {
			got := z.ToBig().String()
			if got != v.i {
				t.Errorf("expect val %s, got %s", v.i, got)
			}
		}
	}
}

func BenchmarkStringBase10(b *testing.B) {
	val := new(Int)
	bytearr := twoPow256
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
