// Copyright 2020 Martin Holst Swende. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the COPYING file.
//

package uint256

import (
	"bytes"
	"math/big"
	"testing"
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
