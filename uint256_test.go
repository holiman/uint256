// Copyright 2019 Martin Holst Swende. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the COPYING file.
//

package uint256

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"math/big"
	"testing"
)

func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

func checkOverflow(b *big.Int, f *Int, overflow bool) error {
	max := big.NewInt(0).SetBytes(Hex2Bytes("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"))
	shouldOverflow := (b.Cmp(max) > 0)
	if overflow != shouldOverflow {
		return fmt.Errorf("Overflow should be %v, was %v\nf= %v\nb= %x\b", shouldOverflow, overflow, f.Hex(), b)
	}
	return nil
}

func randNums() (*big.Int, *Int, error) {
	//How many bits? 0-256
	nbits, _ := rand.Int(rand.Reader, big.NewInt(256))
	//Max random value, a 130-bits integer, i.e 2^130
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(nbits.Int64()), nil)
	b, _ := rand.Int(rand.Reader, max)
	f, overflow := FromBig(b)
	err := checkOverflow(b, f, overflow)
	return b, f, err
}
func randHighNums() (*big.Int, *Int, error) {
	//How many bits? 0-256
	nbits := int64(256)
	//Max random value, a 130-bits integer, i.e 2^130
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(nbits), nil)
	//Generate cryptographically strong pseudo-random between 0 - max
	b, _ := rand.Int(rand.Reader, max)
	f, overflow := FromBig(b)
	//fmt.Printf("f %v\n", f.Hex())
	err := checkOverflow(b, f, overflow)
	return b, f, err
}
func checkEq(b *big.Int, f *Int) bool {
	f2, _ := FromBig(b)
	return f.Eq(f2)
}

func requireEq(t *testing.T, exp *big.Int, got *Int, txt string) bool {
	expF, _ := FromBig(exp)
	if !expF.Eq(got) {
		t.Errorf("got %v expected %v: %v\n", got.Hex(), expF.Hex(), txt)
		return false
	}
	return true
}

func TestBasicStuff(t *testing.T) {
	i, _ := FromBig(big.NewInt(1))
	fmt.Printf("1 %v\n", i.Hex())
	i, _ = FromBig(big.NewInt(-1))
	fmt.Printf("-1 %v\n", i.Hex())
	b := big.NewInt(0)
	b.SetBytes(Hex2Bytes("39d81aff56a841bea668f4c67599a0e1467b49e2e66674cbe36f2d"))
	i, _ = FromBig(b)
	fmt.Printf("%x \n%s\n", b, i.Hex())

	b.SetBytes(Hex2Bytes("dead432298f4ab7ff3fbdbe642972dbbb78835f8ecbea7d3a39dc183d1edbee39787336d1136"))
	i, _ = FromBig(b)
	fmt.Printf("%x \n%s\n", b, i.Hex())

	fmt.Printf("%s \n", NewInt().setBit(255).Hex())
	fmt.Printf("%s \n", NewInt().setBit(254).Hex())
	fmt.Printf("%s \n", NewInt().setBit(32).Hex())
	fmt.Printf("%s \n", NewInt().setBit(1).Hex())
	fmt.Printf("%s \n", NewInt().setBit(0).Hex())

	fmt.Printf("%v \n", NewInt().setBit(0).isBitSet(0))
	fmt.Printf("%v \n", NewInt().setBit(64).isBitSet(64))
	fmt.Printf("%v \n", NewInt().setBit(254).isBitSet(254))

}
func testRandomOp(t *testing.T, nativeFunc func(a, b, c *Int), bigintFunc func(a, b, c *big.Int)) {
	for i := 0; i < 10000; i++ {
		b1, f1, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		b2, f2, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		f1a, f2a := f1.Clone(), f2.Clone()
		nativeFunc(f1, f1, f2)
		bigintFunc(b1, b1, b2)
		//checkOverflow(b, f1, overflow)
		if eq := checkEq(b1, f1); !eq {
			bf, _ := FromBig(b1)
			t.Fatalf("Expected equality:\nf1= %v\nf2= %v\n[ op ]==\nf = %v\nbf= %v\nb = %x\n", f1a.Hex(), f2a.Hex(), f1.Hex(), bf.Hex(), b1)
		}
	}
}
func TestRandomSubOverflow(t *testing.T) {
	for i := 0; i < 10000; i++ {
		b, f1, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		b2, f2, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		f1a, f2a := f1.Clone(), f2.Clone()
		overflow := f1.SubOverflow(f1, f2)
		b.Sub(b, b2)
		checkOverflow(b, f1, overflow)
		if eq := checkEq(b, f1); !eq {
			t.Fatalf("Expected equality:\nf1= %v\nf2= %v\n[ - ]==\nf= %v\nb= %x\n", f1a.Hex(), f2a.Hex(), f1.Hex(), b)
		}
	}
}
func TestRandomSub(t *testing.T) {
	testRandomOp(t,
		func(f1, f2, f3 *Int) {
			f1.Sub(f2, f3)
		},
		func(b1, b2, b3 *big.Int) {
			b1.Sub(b2, b3)
		},
	)
}

func TestRandomAdd(t *testing.T) {
	testRandomOp(t,
		func(f1, f2, f3 *Int) {
			f1.Add(f2, f3)
		},
		func(b1, b2, b3 *big.Int) {
			b1.Add(b2, b3)
		},
	)
}
func TestRandomMul(t *testing.T) {

	testRandomOp(t,
		func(f1, f2, f3 *Int) {
			f1.Mul(f2, f3)
		},
		func(b1, b2, b3 *big.Int) {
			b1.Mul(b2, b3)
		},
	)
}
func TestRandomSquare(t *testing.T) {
	testRandomOp(t,
		func(f1, f2, f3 *Int) {
			f1.Squared()
		},
		func(b1, b2, b3 *big.Int) {
			b1.Mul(b1, b1)
		},
	)
}
func TestRandomDiv(t *testing.T) {
	testRandomOp(t,
		func(f1, f2, f3 *Int) {
			f1.Div(f2, f3)
		},
		func(b1, b2, b3 *big.Int) {
			if b3.Sign() == 0 {
				b1.SetUint64(0)
			} else {
				b1.Div(b2, b3)
			}
		},
	)
}

func TestRandomMod(t *testing.T) {
	testRandomOp(t,
		func(f1, f2, f3 *Int) {
			f1.Mod(f2, f3)
		},
		func(b1, b2, b3 *big.Int) {
			if b3.Sign() == 0 {
				b1.SetUint64(0)
			} else {
				b1.Mod(b2, b3)
			}
		},
	)
}
func TestRandomSMod(t *testing.T) {
	testRandomOp(t,
		func(f1, f2, f3 *Int) {
			f1.Smod(f2, f3)
		},
		func(b1, b2, b3 *big.Int) {
			b1.Set(Smod(b2, b3))
		},
	)
}

func TestRandomMulMod(t *testing.T) {
	for i := 0; i < 10000; i++ {
		b1, f1, err := randNums()
		if err != nil {
			t.Fatalf("Error getting a random number: %v", err)
		}

		b2, f2, err := randNums()
		if err != nil {
			t.Fatalf("Error getting a random number: %v", err)
		}

		b3, f3, err := randNums()
		if err != nil {
			t.Fatalf("Error getting a random number: %v", err)
		}

		b4, f4, err := randNums()
		for b4.Cmp(big.NewInt(0)) == 0 {
			b4, f4, err = randNums()
			if err != nil {
				t.Fatalf("Error getting a random number: %v", err)
			}
		}

		f1.MulMod(f2, f3, f4)
		b1.Mod(big.NewInt(0).Mul(b2, b3), b4)

		if !checkEq(b1, f1) {
			t.Fatalf("Expected equality:\nf2= %v\nf3= %v\nf4= %v\n[ op ]==\nf = %v\nb = %x\n", f2.Hex(), f3.Hex(), f4.Hex(), f1.Hex(), b1)
		}
	}
}

var bigtt255 = bigPow(2, 255)

func S256(x *big.Int) *big.Int {
	if x.Cmp(bigtt255) < 0 {
		return x
	} else {
		return new(big.Int).Sub(x, bigtt256)
	}
}

func TestRandomAbs(t *testing.T) {
	fmt.Printf("SignedMin %x\n", bigtt255)
	fmt.Printf("tt256 %x\n", bigtt256)
	for i := 0; i < 10000; i++ {
		b, f1, err := randHighNums()
		if err != nil {
			t.Fatal(err)
		}
		U256(b)
		b2 := S256(big.NewInt(0).Set(b))
		b2.Abs(b2)
		f1a := f1.Clone().Abs()

		if eq := checkEq(b2, f1a); !eq {
			bf, _ := FromBig(b2)
			t.Fatalf("Expected equality:\nf1= %v\n[ abs ]==\nf = %v\nbf= %v\nb = %x\n", f1.Hex(), f1a.Hex(), bf.Hex(), b2)
		}
	}
}

func TestRandomSDiv(t *testing.T) {
	for i := 0; i < 10000; i++ {
		b, f1, err := randHighNums()
		if err != nil {
			t.Fatal(err)
		}
		b2, f2, err := randHighNums()
		if err != nil {
			t.Fatal(err)
		}
		U256(b)
		U256(b2)

		f1a, f2a := f1.Clone(), f2.Clone()

		f1aAbs, f2aAbs := f1.Clone().Abs(), f2.Clone().Abs()

		f1.Sdiv(f1, f2)
		if b2.BitLen() == 0 {
			// zero
			b = big.NewInt(0)
		} else {
			bb1 := S256(big.NewInt(0).Set(b))
			bb2 := S256(big.NewInt(0).Set(b2))

			b = Sdiv(bb1, bb2)
		}
		if eq := checkEq(b, f1); !eq {
			bf, _ := FromBig(b)
			t.Fatalf("Expected equality:\nf1  = %v\nf2  = %v\n\n\nabs1= %v\nabs2= %v\n[sdiv]==\nf   = %v\nbf  = %v\nb   = %x\n",
				f1a.Hex(), f2a.Hex(), f1aAbs.Hex(), f2aAbs.Hex(), f1.Hex(), bf.Hex(), b)
		}
	}
}

func TestRandomLsh(t *testing.T) {
	for i := 0; i < 10000; i++ {
		b, f1, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		f1a := f1.Clone()
		nbits, _ := rand.Int(rand.Reader, big.NewInt(256))
		n := uint(nbits.Uint64())
		f1.Lsh(f1, n)
		b.Lsh(b, n)
		if eq := checkEq(b, f1); !eq {
			bf, _ := FromBig(b)
			t.Fatalf("Expected equality:\nf1= %v\n n= %v\n[ << ]==\nf = %v\nbf= %v\nb = %x\n", f1a.Hex(), n, f1.Hex(), bf.Hex(), b)
		}
	}
}

func TestRandomRsh(t *testing.T) {
	for i := 0; i < 10000; i++ {
		b, f1, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		f1a := f1.Clone()
		nbits, _ := rand.Int(rand.Reader, big.NewInt(256))
		n := uint(nbits.Uint64())
		f1.Rsh(f1, n)
		b.Rsh(b, n)
		if eq := checkEq(b, f1); !eq {
			t.Fatalf("Expected equality:\nf1= %v\n n= %v\n[ << ]==\nf= %v\nb= %x\n", f1a.Hex(), n, f1.Hex(), b)
		}
	}
}

func TestSrsh(t *testing.T) {
	var n uint = 16
	actual := new(Int).SetBytes(Hex2Bytes("FFFFEEEEDDDDCCCCBBBBAAAA9999888877776666555544443333222211110000"))
	actual.Srsh(actual, n)
	expected := new(Int).SetBytes(Hex2Bytes("FFFFFFFFEEEEDDDDCCCCBBBBAAAA999988887777666655554444333322221111"))
	if !actual.Eq(expected) {
		t.Fatalf("Expected %v, got %v", expected.Hex(), actual.Hex())
	}

	n = 64
	actual = new(Int).SetBytes(Hex2Bytes("FFFFEEEEDDDDCCCCBBBBAAAA9999888877776666555544443333222211110000"))
	actual.Srsh(actual, n)
	expected = new(Int).SetBytes(Hex2Bytes("FFFFFFFFFFFFFFFFFFFFEEEEDDDDCCCCBBBBAAAA999988887777666655554444"))
	if !actual.Eq(expected) {
		t.Fatalf("Expected %v, got %v", expected.Hex(), actual.Hex())
	}

	n = 96
	actual = new(Int).SetBytes(Hex2Bytes("FFFFEEEEDDDDCCCCBBBBAAAA9999888877776666555544443333222211110000"))
	actual.Srsh(actual, n)
	expected = new(Int).SetBytes(Hex2Bytes("FFFFFFFFFFFFFFFFFFFFFFFFFFFFEEEEDDDDCCCCBBBBAAAA9999888877776666"))
	if !actual.Eq(expected) {
		t.Fatalf("Expected %v, got %v", expected.Hex(), actual.Hex())
	}

	n = 256
	actual = new(Int).SetBytes(Hex2Bytes("FFFFEEEEDDDDCCCCBBBBAAAA9999888877776666555544443333222211110000"))
	actual.Srsh(actual, n)
	expected = new(Int).SetBytes(Hex2Bytes("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"))
	if !actual.Eq(expected) {
		t.Fatalf("Expected %v, got %v", expected.Hex(), actual.Hex())
	}

	n = 300
	actual = new(Int).SetBytes(Hex2Bytes("FFFFEEEEDDDDCCCCBBBBAAAA9999888877776666555544443333222211110000"))
	actual.Srsh(actual, n)
	expected = new(Int).SetBytes(Hex2Bytes("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"))
	if !actual.Eq(expected) {
		t.Fatalf("Expected %v, got %v", expected.Hex(), actual.Hex())
	}

	n = 16
	actual = new(Int).SetBytes(Hex2Bytes("7FFFEEEEDDDDCCCCBBBBAAAA9999888877776666555544443333222211110000"))
	actual.Srsh(actual, n)
	expected = new(Int).SetBytes(Hex2Bytes("7FFFEEEEDDDDCCCCBBBBAAAA999988887777666655554444333322221111"))
	if !actual.Eq(expected) {
		t.Fatalf("Expected %v, got %v", expected.Hex(), actual.Hex())
	}

	n = 64
	actual = new(Int).SetBytes(Hex2Bytes("7FFFEEEEDDDDCCCCBBBBAAAA9999888877776666555544443333222211110000"))
	actual.Srsh(actual, n)
	expected = new(Int).SetBytes(Hex2Bytes("7FFFEEEEDDDDCCCCBBBBAAAA999988887777666655554444"))
	if !actual.Eq(expected) {
		t.Fatalf("Expected %v, got %v", expected.Hex(), actual.Hex())
	}

	n = 256
	actual = new(Int).SetBytes(Hex2Bytes("7FFFEEEEDDDDCCCCBBBBAAAA9999888877776666555544443333222211110000"))
	actual.Srsh(actual, n)
	expected = new(Int).SetBytes(nil)
	if !actual.Eq(expected) {
		t.Fatalf("Expected %v, got %v", expected.Hex(), actual.Hex())
	}
}

func TestByte(t *testing.T) {
	z := new(Int).SetBytes(Hex2Bytes("ABCDEF09080706050403020100000000000000000000000000000000000000ef"))
	actual := z.Byte(NewInt().SetUint64(0))
	expected := new(Int).SetBytes(Hex2Bytes("00000000000000000000000000000000000000000000000000000000000000ab"))
	if !actual.Eq(expected) {
		t.Fatalf("Expected %v, got %v", expected.Hex(), actual.Hex())
	}

	z = new(Int).SetBytes(Hex2Bytes("ABCDEF09080706050403020100000000000000000000000000000000000000ef"))
	actual = z.Byte(NewInt().SetUint64(31))
	expected = new(Int).SetBytes(Hex2Bytes("00000000000000000000000000000000000000000000000000000000000000ef"))
	if !actual.Eq(expected) {
		t.Fatalf("Expected %v, got %v", expected.Hex(), actual.Hex())
	}

	z = new(Int).SetBytes(Hex2Bytes("ABCDEF09080706050403020100000000000000000000000000000000000000ef"))
	actual = z.Byte(NewInt().SetUint64(32))
	expected = new(Int).SetBytes(Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"))
	if !actual.Eq(expected) {
		t.Fatalf("Expected %v, got %v", expected.Hex(), actual.Hex())
	}

	z = new(Int).SetBytes(Hex2Bytes("ABCDEF0908070605040302011111111111111111111111111111111111111111"))
	actual = z.Byte(new(Int).SetBytes(Hex2Bytes("f000000000000000000000000000000000000000000000000000000000000001")))
	expected = new(Int).SetBytes(Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"))
	if !actual.Eq(expected) {
		t.Fatalf("Expected %v, got %v", expected.Hex(), actual.Hex())
	}

}
func TestSGT(t *testing.T) {

	x := new(Int).SetBytes(Hex2Bytes("fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"))
	y := new(Int).SetBytes(Hex2Bytes("00"))
	actual := x.Sgt(y)
	expected := false
	if actual != expected {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}

	x = new(Int).SetBytes(Hex2Bytes("00"))
	y = new(Int).SetBytes(Hex2Bytes("fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"))
	actual = x.Sgt(y)
	expected = true
	if actual != expected {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}

const (
	// number of bits in a big.Word
	wordBits = 32 << (uint64(^big.Word(0)) >> 63)
	// number of bytes in a big.Word
	wordBytes = wordBits / 8
)

var (
	tt256m1 = new(big.Int).Sub(bigtt256, big.NewInt(1))
)

// U256 encodes as a 256 bit two's complement number. This operation is destructive.
func U256(x *big.Int) *big.Int {
	return x.And(x, tt256m1)
}

// Exp implements exponentiation by squaring.
// Exp returns a newly-allocated big integer and does not change
// base or exponent. The result is truncated to 256 bits.
//
// Courtesy @karalabe and @chfast
func Exp(base, exponent *big.Int) *big.Int {
	result := big.NewInt(1)

	for _, word := range exponent.Bits() {
		for i := 0; i < wordBits; i++ {
			if word&1 == 1 {
				U256(result.Mul(result, base))
			}
			U256(base.Mul(base, base))
			word >>= 1
		}
	}
	return result
}

func Sdiv(x, y *big.Int) *big.Int {
	if y.Sign() == 0 {
		return new(big.Int)

	}
	n := new(big.Int)
	if x.Sign() == y.Sign() {
		//	if n.Mul(x, y).Sign() < 0 {
		n.SetInt64(1)
	} else {
		n.SetInt64(-1)
	}
	res := x.Div(x.Abs(x), y.Abs(y))
	res.Mul(res, n)
	return res
}
func Smod(x, y *big.Int) *big.Int {
	res := new(big.Int)
	if y.Sign() == 0 {
		return res
	}

	if x.Sign() < 0 {
		res.Mod(x.Abs(x), y.Abs(y))
		res.Neg(res)
	} else {
		res.Mod(x.Abs(x), y.Abs(y))
	}
	return U256(res)

}
func TestRandomExp(t *testing.T) {
	for i := 0; i < 10000; i++ {
		b_base, base, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		b_exp, exp, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		basecopy, expcopy := base.Clone(), exp.Clone()

		f_res := ExpF(base, exp)
		b_res := Exp(b_base, b_exp)
		if eq := checkEq(b_res, f_res); !eq {
			bf, _ := FromBig(b_res)
			t.Fatalf("Expected equality:\nbase= %v\nexp = %v\n[ ^ ]==\nf = %v\nbf= %v\nb = %x\n", basecopy.Hex(), expcopy.Hex(), f_res.Hex(), bf.Hex(), b_res)
		}
	}
}

func TestFixed256bit_Add(t *testing.T) {
	b1 := big.NewInt(0).SetBytes(Hex2Bytes("000282209f633a3ca040e862bb69d92573449d21bce09ea3a74348fbf1ced62e"))
	b2 := big.NewInt(0).SetBytes(Hex2Bytes("00000000000000000000000000000000000000000000003afd56300e26f61922"))

	f, _ := FromBig(b1)
	f2, _ := FromBig(b2)
	fmt.Printf("B1: %x\n", b1)
	fmt.Printf("B2: %x\n", b2)
	fmt.Printf("F1: %s\n", f.Hex())
	fmt.Printf("F2: %s\n", f2.Hex())
	fmt.Println("--")
	b1.Add(b1, b2)
	f.Add(f, f2)
	fmt.Printf("B: %x\n", b1)
	fmt.Printf("F: %s\n", f.Hex())

}

func TestFixed256bit_Sub(t *testing.T) {

	b1 := big.NewInt(0).SetBytes(Hex2Bytes("00000000000000000000000000000000000000000002a3f8ba829e365f479526"))
	b2 := big.NewInt(0).SetBytes(Hex2Bytes("000000000000000000000000000000004ffab28fa389b141ce4876fa1965c937"))

	f, _ := FromBig(b1)
	f2, _ := FromBig(b2)

	fmt.Printf("B1   : %x\n", b1)
	fmt.Printf("B2   : %x\n", b2)
	fmt.Printf("F1   : %s\n", f.Hex())
	fmt.Printf("F2   : %s\n", f2.Hex())
	fmt.Println("--")
	b1.Sub(b1, b2)
	f.Sub(f, f2)
	fmt.Printf("B   : %x\n", b1)
	fmt.Printf("F   : %s\n", f.Hex())

	res, _ := FromBig(b1)
	fmt.Printf("b->f: %s\n", res.Hex())
	fmt.Printf("EQ  : %v\n", f.Eq(res))

}

func TestFixed256bit_Mul(t *testing.T) {

	b1 := big.NewInt(0).SetBytes(Hex2Bytes("00000000000000000000000000000000000000000002a3f8ba829e365f479526"))
	b2 := big.NewInt(0).SetBytes(Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001"))

	/*
		f1= 0000000000000000.0000000000000000.000000000003e6da.5c61238a298b16ec
		f2= 0000000000000000.0000000000000000.0000000000000000.000000000000206c
		[ * ]==
		f = 1d6c3e387e80a3f8.0000000000000bb3.1d6c3e387e80afab.1d6c437ae98b2b90
		bf= 0000000000000000.0000000000000000.000000007e80afab.1d6c437ae98b2b90
		b = 7e80afab1d6c437ae98b2b90

	*/

	f, _ := FromBig(b1)
	f2, _ := FromBig(b2)

	fmt.Printf("B1   : %x\n", b1)
	fmt.Printf("B2   : %x\n", b2)
	fmt.Printf("F1   : %s\n", f.Hex())
	fmt.Printf("F2   : %s\n", f2.Hex())
	fmt.Println("--")
	b1.Mul(b1, b2)
	f.Mul(f, f2)
	fmt.Printf("B   : %x\n", b1)
	fmt.Printf("F   : %s\n", f.Hex())

	res, _ := FromBig(b1)
	fmt.Printf("b->f: %s\n", res.Hex())
	fmt.Printf("EQ  : %v\n", f.Eq(res))
}

func TestFixed256bit_Div(t *testing.T) {

	b1 := big.NewInt(0).SetBytes(Hex2Bytes("12cbafcee8f60f9f3fa308c90fde8d298772ffea667aa6bc109d5c661e7929a5"))
	b2 := big.NewInt(0).SetBytes(Hex2Bytes("00000c76f4afb041407a8ea478d65024f5c3dfe1db1a1bb10c5ea8bec314ccf9"))

	f, _ := FromBig(b1)
	f2, _ := FromBig(b2)

	fmt.Printf("B1   : %x\n", b1)
	fmt.Printf("B2   : %x\n", b2)
	fmt.Printf("F1   : %s\n", f.Hex())
	fmt.Printf("F2   : %s\n", f2.Hex())
	fmt.Println("--")
	b1.Div(b1, b2)
	f.Div(f, f2)
	fmt.Printf("B   : %x\n", b1)
	fmt.Printf("F   : %s\n", f.Hex())

	res, _ := FromBig(b1)
	fmt.Printf("b->f: %s\n", res.Hex())
	fmt.Printf("EQ  : %v\n", f.Eq(res))
}

func TestFixedExp(t *testing.T) {

	b_base := big.NewInt(0).SetBytes(Hex2Bytes("00000000000000000000000000000000000000000000006d5adef08547abf7eb"))
	b_exp := big.NewInt(0).SetBytes(Hex2Bytes("000000000000000000013590cab83b779e708b533b0eef3561483ddeefc841f5"))

	base, _ := FromBig(b_base)
	exp, _ := FromBig(b_exp)

	fmt.Printf("B1   : %x\n", b_base)
	fmt.Printf("B2   : %x\n", b_exp)
	fmt.Printf("F1   : %s\n", base.Hex())
	fmt.Printf("F2   : %s\n", exp.Hex())
	fmt.Println("--")

	//	base.d = 2
	//	exp.d = 255
	res := ExpF(base, exp)

	b_res := Exp(b_base, b_exp)

	want, _ := FromBig(b_res)
	fmt.Printf("B: %x\n", b_res)
	fmt.Printf("want : %s\n", want.Hex())
	fmt.Printf("got  : %s\n", res.Hex())

	fb_res, _ := FromBig(b_res)
	fmt.Printf("EQ  : %v\n", res.Eq(fb_res))

}
func TestAddmod(t *testing.T) {
	b1 := big.NewInt(0).SetBytes(Hex2Bytes("fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd"))
	b2 := big.NewInt(0).SetBytes(Hex2Bytes("fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"))
	b3 := big.NewInt(0).SetBytes(Hex2Bytes("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"))
	ex := big.NewInt(0).SetBytes(Hex2Bytes("fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc"))

	f1, _ := FromBig(b1)
	f2, _ := FromBig(b2)
	f3, _ := FromBig(b3)
	res := NewInt()
	res.AddMod(f1, f2, f3)
	requireEq(t, ex, res, "1 res wrong")
	requireEq(t, b1, f1, "1 f1 changed")
	requireEq(t, b2, f2, "1 f2 changed")
	requireEq(t, b3, f3, "1 f3 changed")

	f1, _ = FromBig(b1)
	f2, _ = FromBig(b2)
	f3, _ = FromBig(b3)
	f1.AddMod(f1, f2, f3)
	requireEq(t, ex, f1, "2 f1 wrong")
	requireEq(t, b2, f2, "2 f2 changed")
	requireEq(t, b3, f3, "2 f3 changed")

	f1, _ = FromBig(b1)
	f2, _ = FromBig(b2)
	f3, _ = FromBig(b3)
	f2.AddMod(f1, f2, f3)
	requireEq(t, ex, f2, "3 f2 wrong")
	requireEq(t, b1, f1, "3 f1 changed")
	requireEq(t, b3, f3, "3 f3 changed")

	f1, _ = FromBig(b1)
	f2, _ = FromBig(b2)
	f3, _ = FromBig(b3)
	f3.AddMod(f1, f2, f3)
	requireEq(t, ex, f3, "4 f3 wrong")
	requireEq(t, b1, f1, "4 f1 changed")
	requireEq(t, b2, f2, "4 f2 changed")

	b1 = big.NewInt(0).SetBytes(Hex2Bytes("fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd"))
	b2 = big.NewInt(0).SetBytes(Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000003"))
	b3 = big.NewInt(0).SetBytes(Hex2Bytes("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"))
	ex = big.NewInt(0).SetBytes(Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001"))

	f1, _ = FromBig(b1)
	f2, _ = FromBig(b2)
	f3, _ = FromBig(b3)
	f3.AddMod(f1, f2, f3)
	requireEq(t, ex, f3, "5 f3 wrong")
}

func benchmark_Add_Bit(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	f, _ := FromBig(b1)
	f2, _ := FromBig(b2)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f.Add(f, f2)
	}
}
func benchmark_Add_Big(bench *testing.B) {
	b := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
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
	b1 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	f, _ := FromBig(b1)
	f2, _ := FromBig(b2)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f.SubOverflow(f, f2)
	}
}
func benchmark_Sub_Bit(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	f, _ := FromBig(b1)
	f2, _ := FromBig(b2)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f.Sub(f, f2)
	}
}

func benchmark_Sub_Big(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))

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
	a := big.NewInt(0).SetBytes(Hex2Bytes("f123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b := big.NewInt(0).SetBytes(Hex2Bytes("f123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1 := big.NewInt(0)
		b1.Mul(a, b)
		U256(b1)
	}
}

func benchmark_Mul_Bit(bench *testing.B) {
	a := big.NewInt(0).SetBytes(Hex2Bytes("f123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b := big.NewInt(0).SetBytes(Hex2Bytes("f123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	fa, _ := FromBig(a)
	fb, _ := FromBig(b)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f := NewInt()
		f.Mul(fa, fb)
	}
}

func benchmark_Squared_Bit(bench *testing.B) {
	a := big.NewInt(0).SetBytes(Hex2Bytes("f123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	fa, _ := FromBig(a)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f := NewInt().Copy(fa)
		f.Squared()
	}
}
func benchmark_Squared_Big(bench *testing.B) {
	a := big.NewInt(0).SetBytes(Hex2Bytes("f123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))

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
	b1 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1.And(b1, b2)
	}
}
func benchmark_And_Bit(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
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
	b1 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1.Or(b1, b2)
	}
}
func benchmark_Or_Bit(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
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
	b1 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1.Xor(b1, b2)
	}
}
func benchmark_Xor_Bit(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
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
	b1 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1.Cmp(b2)
	}
}
func benchmark_Cmp_Bit(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
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
	original := big.NewInt(0).SetBytes(Hex2Bytes("FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"))
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
	original := big.NewInt(0).SetBytes(Hex2Bytes("FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"))
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
	original := big.NewInt(0).SetBytes(Hex2Bytes("FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"))
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
	original := big.NewInt(0).SetBytes(Hex2Bytes("FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"))
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

	orig := big.NewInt(0).SetBytes(Hex2Bytes(x))
	base := big.NewInt(0).SetBytes(Hex2Bytes(x))
	exp := big.NewInt(0).SetBytes(Hex2Bytes(y))

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		Exp(base, exp)
		base.Set(orig)
	}
}
func benchmark_Exp_Bit(bench *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	base := big.NewInt(0).SetBytes(Hex2Bytes(x))
	exp := big.NewInt(0).SetBytes(Hex2Bytes(y))

	f_base, _ := FromBig(base)
	f_orig, _ := FromBig(base)
	f_exp, _ := FromBig(exp)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		ExpF(f_base, f_exp)
		f_base.Copy(f_orig)
	}
}
func benchmark_ExpSmall_Big(bench *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "8abcdef"

	orig := big.NewInt(0).SetBytes(Hex2Bytes(x))
	base := big.NewInt(0).SetBytes(Hex2Bytes(x))
	exp := big.NewInt(0).SetBytes(Hex2Bytes(y))

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		Exp(base, exp)
		base.Set(orig)
	}
}
func benchmark_ExpSmall_Bit(bench *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "8abcdef"

	base := big.NewInt(0).SetBytes(Hex2Bytes(x))
	exp := big.NewInt(0).SetBytes(Hex2Bytes(y))

	f_base, _ := FromBig(base)
	f_orig, _ := FromBig(base)
	f_exp, _ := FromBig(exp)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		ExpF(f_base, f_exp)
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
		x := big.NewInt(0).SetBytes(Hex2Bytes(a))
		y := big.NewInt(0).SetBytes(Hex2Bytes(b))
		z := big.NewInt(0).SetBytes(Hex2Bytes(m))
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
		x, _ := FromBig(big.NewInt(0).SetBytes(Hex2Bytes(a)))
		y, _ := FromBig(big.NewInt(0).SetBytes(Hex2Bytes(b)))
		z, _ := FromBig(big.NewInt(0).SetBytes(Hex2Bytes(m)))
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

func benchModBigint(a, b string) func(*testing.B) {
	return func(bench *testing.B) {
		x := big.NewInt(0).SetBytes(Hex2Bytes(a))
		y := big.NewInt(0).SetBytes(Hex2Bytes(b))
		z := big.NewInt(0)
		bench.ResetTimer()
		for i := 0; i < bench.N; i++ {
			z.Mod(x, y)
		}
	}
}
func benchModUint256(a, b string) func(*testing.B) {
	return func(bench *testing.B) {
		x, _ := FromBig(big.NewInt(0).SetBytes(Hex2Bytes(a)))
		y, _ := FromBig(big.NewInt(0).SetBytes(Hex2Bytes(b)))
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
	a := big.NewInt(0).SetBytes(Hex2Bytes("1fc2bad1e611"))
	b := big.NewInt(0).SetBytes(Hex2Bytes("12bad1e611"))

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1 := big.NewInt(0)
		b1.Div(a, b)
		U256(b1)
	}
}

func benchmark_DivSmall_Bit(bench *testing.B) {
	a := big.NewInt(0).SetBytes(Hex2Bytes("1fc2bad1e611"))
	b := big.NewInt(0).SetBytes(Hex2Bytes("12bad1e611"))
	fa, _ := FromBig(a)
	fb, _ := FromBig(b)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f := NewInt()
		f.Div(fa, fb)
	}
}
func benchmark_DivLarge_Big(bench *testing.B) {
	a := big.NewInt(0).SetBytes(Hex2Bytes("fe7fb0d1f59dfe9492ffbf73683fd1e870eec79504c60144cc7f5fc2bad1e611"))
	b := big.NewInt(0).SetBytes(Hex2Bytes("ff3f9014f20db29ae04af2c2d265de17"))

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1 := big.NewInt(0)
		b1.Div(a, b)
		U256(b1)
	}
}

func benchmark_DivLarge_Bit(bench *testing.B) {
	a := big.NewInt(0).SetBytes(Hex2Bytes("fe7fb0d1f59dfe9492ffbf73683fd1e870eec79504c60144cc7f5fc2bad1e611"))
	b := big.NewInt(0).SetBytes(Hex2Bytes("ff3f9014f20db29ae04af2c2d265de17"))
	fa, _ := FromBig(a)
	fb, _ := FromBig(b)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f := NewInt()
		f.Div(fa, fb)
	}
}
func benchmark_SdivLarge_Big(bench *testing.B) {
	a := big.NewInt(0).SetBytes(Hex2Bytes("800fffffffffffffffffffffffffd1e870eec79504c60144cc7f5fc2bad1e611"))
	b := big.NewInt(0).SetBytes(Hex2Bytes("ff3f9014f20db29ae04af2c2d265de17"))

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
	a := big.NewInt(0).SetBytes(Hex2Bytes("800fffffffffffffffffffffffffd1e870eec79504c60144cc7f5fc2bad1e611"))
	b := big.NewInt(0).SetBytes(Hex2Bytes("ff3f9014f20db29ae04af2c2d265de17"))
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

func TestByteRepresentation(t *testing.T) {

	//0e320219838e859b2f9f18b72e3d4073ca50b37d
	a := big.NewInt(0xFF0AFcafe)
	aa := NewInt().SetUint64(0xFF0afcafe)
	bb := NewInt().SetBytes(a.Bytes())
	x := a.Bytes()
	y := aa.Bytes()
	z := bb.Bytes()
	fmt.Printf("x %x  y %x z %x \n", x, y, z)
	fmt.Printf("padded %x\n", bb.PaddedBytes(32))
	fmt.Printf("padded %x\n", bb.PaddedBytes(20))
	fmt.Printf("padded %x\n", bb.PaddedBytes(40))
}

func TestByteRepresentation2(t *testing.T) {

	bytearr := Hex2Bytes("0e320219838e859b2f9f18b72e3d4073ca50b37d")
	a := big.NewInt(0).SetBytes(bytearr)
	aa := NewInt().SetBytes(bytearr)
	bb := NewInt().SetBytes(a.Bytes())
	x := a.Bytes()
	y := aa.Bytes()
	z := bb.Bytes()
	fmt.Printf("x %x\ny %x\nz %x\n", x, y, z)
	fmt.Printf("padded %x\n", bb.PaddedBytes(32))
	fmt.Printf("padded %x\n", bb.PaddedBytes(20))
	fmt.Printf("padded %x\n", bb.PaddedBytes(40))
}

func benchAddModBigint(a, b, m string) func(*testing.B) {
	return func(bench *testing.B) {
		x := big.NewInt(0).SetBytes(Hex2Bytes(a))
		y := big.NewInt(0).SetBytes(Hex2Bytes(b))
		z := big.NewInt(0).SetBytes(Hex2Bytes(m))
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
		x, _ := FromBig(big.NewInt(0).SetBytes(Hex2Bytes(a)))
		y, _ := FromBig(big.NewInt(0).SetBytes(Hex2Bytes(b)))
		z, _ := FromBig(big.NewInt(0).SetBytes(Hex2Bytes(m)))
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

func TestWriteToSlice(t *testing.T) {
	x1 := fromHex("fe7fb0d1f59dfe9492ffbf73683fd1e870eec79504c60144cc7f5fc2bad1e611")

	a := big.NewInt(0).SetBytes(x1)
	fa, _ := FromBig(a)

	dest := make([]byte, 32)
	fa.WriteToSlice(dest)
	if bytes.Compare(dest, x1) != 0 {
		t.Errorf("got %x, expected %x", dest, x1)
	}

	fb := NewInt()
	exp := make([]byte, 32)
	fb.WriteToSlice(dest)
	if bytes.Compare(dest, exp) != 0 {
		t.Errorf("got %x, expected %x", dest, exp)
	}
	// a too small buffer
	// Should fill the lower parts, masking upper bytes
	exp = fromHex("683fd1e870eec79504c60144cc7f5fc2bad1e611")
	dest = make([]byte, 20)
	fa.WriteToSlice(dest)
	if bytes.Compare(dest, exp) != 0 {
		t.Errorf("got %x, expected %x", dest, exp)
	}

	// a too large buffer, already filled with stuff
	// Should fill the leftmost 32 bytes, not touch the other things
	dest = fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	exp = fromHex("fe7fb0d1f59dfe9492ffbf73683fd1e870eec79504c60144cc7f5fc2bad1e611ffffffffffffffff")

	fa.WriteToSlice(dest)
	if bytes.Compare(dest, exp) != 0 {
		t.Errorf("got %x, expected %x", dest, x1)
	}

	// an empty slice, no panics please
	dest = []byte{}
	exp = []byte{}

	fa.WriteToSlice(dest)
	if bytes.Compare(dest, exp) != 0 {
		t.Errorf("got %x, expected %x", dest, x1)
	}

}
func TestInt_WriteToArray(t *testing.T) {
	x1 := fromHex("0000000000000000000000000000d1e870eec79504c60144cc7f5fc2bad1e611")
	a := big.NewInt(0).SetBytes(x1)
	fa, _ := FromBig(a)

	{
		dest := [20]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
		fa.WriteToArray20(&dest)
		exp := fromHex("0000d1e870eec79504c60144cc7f5fc2bad1e611")
		if bytes.Compare(dest[:], exp) != 0 {
			t.Errorf("got %x, expected %x", dest, exp)
		}

	}

	{
		dest := [32]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
		fa.WriteToArray32(&dest)
		exp := fromHex("0000000000000000000000000000d1e870eec79504c60144cc7f5fc2bad1e611")
		if bytes.Compare(dest[:], exp) != 0 {
			t.Errorf("got %x, expected %x", dest, exp)
		}

	}
}

func TestFromBFromBigigOld(t *testing.T) {
	z, o := FromBig(new(big.Int).SetBytes(Hex2Bytes("ababee444444444444ffcc333333333333ddaa222222222222bb8811111111111199")))
	fmt.Printf("z %x\n", z)
	if !o {
		t.Errorf("expected overflow, got %v", o)
	}
	z, o = FromBig(new(big.Int).SetBytes(Hex2Bytes("ee444444444444ffcc333333333333ddaa222222222222bb8811111111111199")))
	if o {
		t.Errorf("expected no overflow, got %v", o)
	}
	b := new(big.Int).SetBytes(Hex2Bytes("ee444444444444ffcc333333333333ddaa222222222222bb8811111111111199"))
	z, o = FromBig(b.Neg(b))
	fmt.Printf("z %x\n", z)
	if o {
		t.Errorf("expected no overflow, got %v", o)
	}

}

func TestFromBig(t *testing.T) {
	z, o := NewFromBig(new(big.Int).SetBytes(Hex2Bytes("ababee444444444444ffcc333333333333ddaa222222222222bb8811111111111199")))
	fmt.Printf("z %x\n", z)
	if !o {
		t.Errorf("expected overflow, got %v", o)
	}
	z, o = NewFromBig(new(big.Int).SetBytes(Hex2Bytes("ee444444444444ffcc333333333333ddaa222222222222bb8811111111111199")))
	if o {
		t.Errorf("expected no overflow, got %v", o)
	}
	b := new(big.Int).SetBytes(Hex2Bytes("ee444444444444ffcc333333333333ddaa222222222222bb8811111111111199"))
	z, o = NewFromBig(b.Neg(b))
	fmt.Printf("z %x\n", z)
	if o {
		t.Errorf("expected no overflow, got %v", o)
	}
}

const addrSize = 20

type gethAddress [addrSize]byte

// SetBytes sets the address to the value of b.
// If b is larger than len(a) it will panic.
func (a *gethAddress) setBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-addrSize:]
	}
	copy(a[addrSize-len(b):], b)
}

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func bytesToAddress(b []byte) gethAddress {
	var a gethAddress
	a.setBytes(b)
	return a
}

// BigToAddress returns Address with byte values of b.
// If b is larger than len(h), b will be cropped from the left.
func bigToAddress(b *big.Int) gethAddress { return bytesToAddress(b.Bytes()) }

func TestByte20Representation(t *testing.T) {

	for i, tt := range []string{
		"1337fafafa0e320219838e859b2f9f18b72e3d4073ca50b37d",
		"fafafa0e320219838e859b2f9f18b72e3d4073ca50b37d",
		"0e320219838e859b2f9f18b72e3d4073ca50b37d",
		"320219838e859b2f9f18b72e3d4073ca50b37d",
		"838e859b2f9f18b72e3d4073ca50b37d",
		"38e859b2f9f18b72e3d4073ca50b37d",
		"f18b72e3d4073ca50b37d",
		"b37d",
		"01",
		"",
		"00",
	} {
		bytearr := Hex2Bytes(tt)
		// big.Int -> address
		a := big.NewInt(0).SetBytes(bytearr)
		exp := bytesToAddress(a.Bytes())

		// uint256.Int -> address
		b := NewInt().SetBytes(bytearr)
		got := gethAddress(b.Bytes20())

		if got != exp {
			t.Errorf("testcase %d: got %x exp %x", i, got, exp)
		}
		fmt.Printf("got %x \n", got)

	}

}
