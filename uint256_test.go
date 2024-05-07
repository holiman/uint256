// uint256: Fixed size 256-bit math library
// Copyright 2018-2020 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

package uint256

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"testing"
)

var (
	bigtt256 = new(big.Int).Lsh(big.NewInt(1), 256)
	bigtt255 = new(big.Int).Lsh(big.NewInt(1), 255)

	unTestCases = []string{
		"0x0",
		"0x1",
		"0x80000000000000000000000000000000",
		"0x80000000000000010000000000000000",
		"0x80000000000000000000000000000001",
		"0x12cbafcee8f60f9f3fa308c90fde8d298772ffea667aa6bc109d5c661e7929a5",
		"0x8000000000000000000000000000000000000000000000000000000000000000",
		"0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe",
		"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
	}

	// A collection of interesting input values for binary operators (especially for division).
	// No expected results as big.Int can be used as the source of truth.
	binTestCases = [][2]string{
		{"0x0", "0x0"},
		{"0x1", "0x0"},
		{"0x1", "0x767676767676767676000000767676767676"},
		{"0x2", "0x0"},
		{"0x2", "0x1"},
		{"0x12cbafcee8f60f9f3fa308c90fde8d298772ffea667aa6bc109d5c661e7929a5", "0xc76f4afb041407a8ea478d65024f5c3dfe1db1a1bb10c5ea8bec314ccf9"},
		{"0x10000000000000000", "0x2"},
		{"0x7000000000000000", "0x8000000000000000"},
		{"0x8000000000000000", "0x8000000000000000"},
		{"0x8000000000000001", "0x8000000000000000"},
		{"0x80000000000000010000000000000000", "0x80000000000000000000000000000000"},
		{"0x80000000000000000000000000000000", "0x80000000000000000000000000000001"},
		{"0x478392145435897052", "0x111"},
		{"0x767676767676767676000000767676767676", "0x2900760076761e00020076760000000076767676000000"},
		{"0x12121212121212121212121212121212", "0x232323232323232323"},
		{"0xfffff716b61616160b0b0b2b0b0b0becf4bef50a0df4f48b090b2b0bc60a0a00", "0xfffff716b61616160b0b0b2b0b230b000008010d0a2b00"},
		{"0x50beb1c60141a0000dc2b0b0b0b0b0b410a0a0df4f40b090b2b0bc60a0a00", "0x2000110000000d0a300e750a000000090a0a"},
		{"0x4b00000b41000b0b0b2b0b0b0b0b0b410a0aeff4f40b090b2b0bc60a0a1000", "0x4b00000b41000b0b0b2b0b0b0b0b0b410a0aeff4f40b0a0a"},
		{"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", "0x7"},
		{"0xf6376770abd3a36b20394c5664afef1194c801c3f05e42566f085ed24d002bb0", "0xb368d219438b7f3f"},
		{"0x0", "0x10900000000000000000000000000000000000000000000000000"},
		{"0x77676767676760000000000000001002e000000000000040000000e000000000", "0xfffc000000000000767676240000000000002b0576047"},
		{"0x767676767676000000000076000000000000005600000000000000000000", "0x767676767676000000000076000000760000"},
		{"0x8200000000000000000000000000000000000000000000000000000000000000", "0x8200000000000000fe000004000000ffff000000fffff700"},
		{"0xdac7fff9ffd9e1322626262626262600", "0xd021262626262626"},
		{"0x8000000000000001800000000000000080000000000000008000000000000000", "0x800000000000000080000000000000008000000000000000"},
		{"0xe8e8e8e2000100000009ea02000000000000ff3ffffff80000001000220000", "0xe8e8e8e2000100000009ea02000000000000ff3ffffff800000010002280ff"},
		{"0xc9700000000000000000023f00c00014ff000000000000000022300805", "0xc9700000000000000000023f00c00014ff002c000000000000223108"},
		{"0x40000000fd000000db0000000000000000000000000000000000000000000001", "0x40000000fd000000db0000000000000000000040000000fd000000db000001"},
		{"0x40000000fd000000db0000000000000000000000000000000000000000000001", "0x40000000fd000000db0000000000000000000040000000fd000000db0000d3"},
		{"0x1f000000000000000000000000000000200000000100000000000000000000", "0x100000000ffffffffffffffff0000000000002e000000"},
		{"0x7effffff80000000000000000000000000020000440000000000000000000001", "0x7effffff800000007effffff800000008000ff0000010000"},
		{"0x5fd8fffffffffffffffffffffffffffffc090000ce700004d0c9ffffff000001", "0x2ffffffffffffffffffffffffffffffffff000000030000"},
		{"0x62d8fffffffffffffffffffffffffffffc18000000000000000000ca00000001", "0x2ffffffffffffffffffffffffffffffffff200000000000"},
		{"0x7effffff8000000000000000000000000000000000000000d900000000000001", "0x7effffff8000000000000000000000000000000000008001"},
		{"0x6400aff20ff00200004e7fd1eff08ffca0afd1eff08ffca0a", "0x210000000000000022"},
		{"0x6d5adef08547abf7eb", "0x13590cab83b779e708b533b0eef3561483ddeefc841f5"},
		{"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe", "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe", "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"0xe8e8e8e2000100000009ea02000000000000ff3ffffff80000001000220000", "0xffffffffffffffff7effffff800000007effffff800000008000ff0000010000"},
		{"0x1ce97e1ab91a", "0x66aa0a5319bcf5cb4"}, // regression test for udivrem() where len(x) < len(y)
	}

	// A collection of interesting input values for ternary operators (addmod, mulmod).
	ternTestCases = [][3]string{
		{"0x0", "0x0", "0x0"},
		{"0x1", "0x0", "0x0"},
		{"0x1", "0x1", "0x0"},
		{"0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd", "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe", "0x0"},
		{"0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd", "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe", "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd", "0x3", "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"},
		{"0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0x2"},
		{"0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0x1"},
		{"0xffffffffffffffffffffffffffffffff", "0xffffffffffffffffffffffffffffffff", "0xfffffffffffffffffffffffffffffffe00000000000000000000000000000002"},
		{"0xffffffffffffffffffffffffffffffff", "0xffffffffffffffffffffffffffffffff", "0xfffffffffffffffffffffffffffffffe00000000000000000000000000000001"},
		{"0xffffffffffffffffffffffffffff000004020041fffffffffc00000060000020", "0xffffffffffffffffffffffffffffffe6000000ffffffe60000febebeffffffff", "0xffffffffffffffffffe6000000ffffffe60000febebeffffffffffffffffffff"},
		{"0xffffffffffffffffffffffffffffffff00ffffe6ff0000000000000060000020", "0xffffffffffffffffffffffffffffffffffe6000000ffff00e60000febebeffff", "0xffffffffffffffffffe6000000ffff00e60000fe0000ffff00e60000febebeff"},
		{"0xfffffffffffffffffffffffff600000000005af50100bebe000000004a00be0a", "0xffffffffffffffffffffffffffffeaffdfd9fffffffffffff5f60000000000ff", "0xffffffffffffffffffffffeaffdfd9fffffffffffffff60000000000ffffffff"},
		{"0x8000000000000001000000000000000000000000000000000000000000000000", "0x800000000000000100000000000000000000000000000000000000000000000b", "0x8000000000000000000000000000000000000000000000000000000000000000"},
		{"0x8000000000000000000000000000000000000000000000000000000000000000", "0x8000000000000001000000000000000000000000000000000000000000000000", "0x8000000000000000000000000000000000000000000000000000000000000000"},
		{"0x8000000000000000000000000000000000000000000000000000000000000000", "0x8000000000000001000000000000000000000000000000000000000000000000", "0x8000000000000001000000000000000000000000000000000000000000000000"},
		{"0x8000000000000000000000000000000000000000000000000000000000000000", "0x8000000000000000000000000000000100000000000000000000000000000000", "0x8000000000000000000000000000000000000000000000000000000000000001"},
		{"0x1", "0x1", "0xffffffff00000001000000000000000000000000ffffffffffffffffffffffff"},
		{"0x1", "0x1", "0x1000000003030303030303030303030303030303030303030303030303030"},
		{"0x1", "0x1", "0x4000000000000000130303030303030303030303030303030303030303030"},
		{"0x1", "0x1", "0x8000000000000000000000000000000043030303000000000"},
		{"0x1", "0x1", "0x8000000000000000000000000000000003030303030303030"},
	}
)

func hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

func checkOverflow(b *big.Int, f *Int, overflow bool) error {
	max := big.NewInt(0).SetBytes(hex2Bytes("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"))
	shouldOverflow := b.Cmp(max) > 0
	if overflow != shouldOverflow {
		return fmt.Errorf("Overflow should be %v, was %v\nf= %x\nb= %x\b", shouldOverflow, overflow, f, b)
	}
	return nil
}

func randNum() *Int {
	//How many bits? 0-256
	nbits, _ := rand.Int(rand.Reader, big.NewInt(257))
	//Max random value, a 130-bits integer, i.e 2^130
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(nbits.Int64()), nil)
	b, _ := rand.Int(rand.Reader, max)
	f, _ := FromBig(b)
	return f
}

func randNums() (*big.Int, *Int) {
	//How many bits? 0-256
	nbits, _ := rand.Int(rand.Reader, big.NewInt(257))
	//Max random value, a 130-bits integer, i.e 2^130
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(nbits.Int64()), nil)
	b, _ := rand.Int(rand.Reader, max)
	f, overflow := FromBig(b)
	if err := checkOverflow(b, f, overflow); err != nil {
		panic(err)
	}
	return b, f
}

func randHighNums() (*big.Int, *Int) {
	//How many bits? 0-256
	nbits := int64(256)
	//Max random value, a 130-bits integer, i.e 2^130
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(nbits), nil)
	//Generate cryptographically strong pseudo-random between 0 - max
	b, _ := rand.Int(rand.Reader, max)
	f, overflow := FromBig(b)
	if err := checkOverflow(b, f, overflow); err != nil {
		panic(err)
	}
	return b, f
}

func checkEq(b *big.Int, f *Int) bool {
	f2, _ := FromBig(b)
	return f.Eq(f2)
}

func requireEq(t *testing.T, exp *big.Int, got *Int, txt string) bool {
	expF, _ := FromBig(exp)
	if !expF.Eq(got) {
		t.Errorf("got %x expected %x: %v\n", got, expF, txt)
		return false
	}
	return true
}

func TestRandomSubOverflow(t *testing.T) {
	for i := 0; i < 10000; i++ {
		b, f1 := randNums()
		b2, f2 := randNums()
		f1a, f2a := f1.Clone(), f2.Clone()
		_, overflow := f1.SubOverflow(f1, f2)
		b.Sub(b, b2)

		// check overflow
		if have, want := overflow, b.Cmp(big.NewInt(0)) < 0; have != want {
			t.Fatalf("underflow should be %v, was %v\nf= %x\nb= %x\b", have, want, f1, b)
		}
		if eq := checkEq(b, f1); !eq {
			t.Fatalf("Expected equality:\nf1= %x\nf2= %x\n[ - ]==\nf= %x\nb= %x\n", f1a, f2a, f1, b)
		}
	}
}

func TestRandomMulOverflow(t *testing.T) {
	for i := 0; i < 10000; i++ {
		b, f1 := randNums()
		b2, f2 := randNums()

		f1a, f2a := f1.Clone(), f2.Clone()
		_, overflow := f1.MulOverflow(f1, f2)
		b.Mul(b, b2)
		if err := checkOverflow(b, f1, overflow); err != nil {
			t.Fatal(err)
		}
		if eq := checkEq(b, f1); !eq {
			t.Fatalf("Expected equality:\nf1= %x\nf2= %x\n[ - ]==\nf= %x\nb= %x\n", f1a, f2a, f1, b)
		}
	}
}

// divModDiv wraps DivMod and returns quotient only
func divModDiv(z, x, y *Int) *Int {
	var m Int
	z.DivMod(x, y, &m)
	return z
}

// divModMod wraps DivMod and returns modulus only
func divModMod(z, x, y *Int) *Int {
	new(Int).DivMod(x, y, z)
	return z
}

// udivremDiv wraps udivrem and returns quotient
func udivremDiv(z, x, y *Int) *Int {
	var quot Int
	if !y.IsZero() {
		udivrem(quot[:], x[:], y)
	}
	return z.Set(&quot)
}

// udivremMod wraps udivrem and returns remainder
func udivremMod(z, x, y *Int) *Int {
	if y.IsZero() {
		return z.Clear()
	}
	var quot Int
	rem := udivrem(quot[:], x[:], y)
	return z.Set(&rem)
}

func set3Int(s1, s2, s3, d1, d2, d3 *Int) {
	d1.Set(s1)
	d2.Set(s2)
	d3.Set(s3)
}

func set3Big(s1, s2, s3, d1, d2, d3 *big.Int) {
	d1.Set(s1)
	d2.Set(s2)
	d3.Set(s3)
}

func TestRandomMulMod(t *testing.T) {
	for i := 0; i < 10000; i++ {
		b1, f1 := randNums()
		b2, f2 := randNums()
		b3, f3 := randNums()
		b4, f4 := randNums()
		for b4.Cmp(big.NewInt(0)) == 0 {
			b4, f4 = randNums()
		}

		f1.MulMod(f2, f3, f4)
		b1.Mod(big.NewInt(0).Mul(b2, b3), b4)

		if !checkEq(b1, f1) {
			t.Fatalf("Expected equality:\nf2= %x\nf3= %x\nf4= %x\n[ op ]==\nf = %x\nb = %x\n", f2, f3, f4, f1, b1)
		}

		f1.mulModWithReciprocalWrapper(f2, f3, f4)

		if !checkEq(b1, f1) {
			t.Fatalf("Expected equality:\nf2= %x\nf3= %x\nf4= %x\n[ op ]==\nf = %x\nb = %x\n", f2, f3, f4, f1, b1)
		}
	}

	// Tests related to powers of 2

	f_minusone := &Int{^uint64(0), ^uint64(0), ^uint64(0), ^uint64(0)}

	b_one := big.NewInt(1)
	b_minusone := big.NewInt(0)
	b_minusone = b_minusone.Sub(b_minusone, b_one)

	for i := uint(0); i < 256; i++ {
		b := big.NewInt(0)
		f := NewInt(0)

		t1, t2, t3 := b, b, b
		u1, u2, u3 := f, f, f

		b1 := b.Lsh(b, i)
		f1 := f.Lsh(f, i)

		b2, f2 := randNums()
		for b2.Cmp(big.NewInt(0)) == 0 {
			b2, f2 = randNums()
		}

		b3, f3 := randNums()
		for b3.Cmp(big.NewInt(0)) == 0 {
			b3, f3 = randNums()
		}

		// Tests with one operand a power of 2

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t1, t2), t3)
		f.MulMod(u1, u2, u3)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf1= 0x%x\nf2= 0x%x\nf3= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f1, f2, f3, f, b)
		}

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t1, t3), t2)
		f.MulMod(u1, u3, u2)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf1= 0x%x\nf3= 0x%x\nf2= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f1, f3, f2, f, b)
		}

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t2, t1), t3)
		f.MulMod(u2, u1, u3)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf2= 0x%x\nf1= 0x%x\nf3= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f2, f1, f3, f, b)
		}

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t2, t3), t1)
		f.MulMod(u2, u3, u1)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf2= 0x%x\nf3= 0x%x\nf1= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f2, f3, f1, f, b)
		}

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t3, t1), t2)
		f.MulMod(u3, u1, u2)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf3= 0x%x\nf1= 0x%x\nf2= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f3, f1, f2, f, b)
		}

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t3, t2), t1)
		f.MulMod(u3, u2, u1)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf3= 0x%x\nf2= 0x%x\nf1= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f3, f2, f1, f, b)
		}

		// Tests with one operand 2^256 minus a power of 2

		f1.Xor(f1, f_minusone)
		b1.Xor(b1, b_minusone)

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t1, t2), t3)
		f.MulMod(u1, u2, u3)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf1= 0x%x\nf2= 0x%x\nf3= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f1, f2, f3, f, b)
		}

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t1, t3), t2)
		f.MulMod(u1, u3, u2)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf1= 0x%x\nf3= 0x%x\nf2= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f1, f3, f2, f, b)
		}

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t2, t1), t3)
		f.MulMod(u2, u1, u3)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf2= 0x%x\nf1= 0x%x\nf3= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f2, f1, f3, f, b)
		}

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t2, t3), t1)
		f.MulMod(u2, u3, u1)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf2= 0x%x\nf3= 0x%x\nf1= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f2, f3, f1, f, b)
		}

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t3, t1), t2)
		f.MulMod(u3, u1, u2)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf3= 0x%x\nf1= 0x%x\nf2= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f3, f1, f2, f, b)
		}

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t3, t2), t1)
		f.MulMod(u3, u2, u1)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf3= 0x%x\nf2= 0x%x\nf1= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f3, f2, f1, f, b)
		}

		f1.Xor(f1, f_minusone)
		b1.Xor(b1, b_minusone)

		// Tests with one operand a power of 2 plus 1

		b1.Add(b1, b_one)
		f1.AddUint64(f1, 1)

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t1, t2), t3)
		f.MulMod(u1, u2, u3)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf1= 0x%x\nf2= 0x%x\nf3= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f1, f2, f3, f, b)
		}

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t1, t3), t2)
		f.MulMod(u1, u3, u2)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf1= 0x%x\nf3= 0x%x\nf2= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f1, f3, f2, f, b)
		}

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t2, t1), t3)
		f.MulMod(u2, u1, u3)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf2= 0x%x\nf1= 0x%x\nf3= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f2, f1, f3, f, b)
		}

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t2, t3), t1)
		f.MulMod(u2, u3, u1)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf2= 0x%x\nf3= 0x%x\nf1= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f2, f3, f1, f, b)
		}

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t3, t1), t2)
		f.MulMod(u3, u1, u2)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf3= 0x%x\nf1= 0x%x\nf2= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f3, f1, f2, f, b)
		}

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t3, t2), t1)
		f.MulMod(u3, u2, u1)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf3= 0x%x\nf2= 0x%x\nf1= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f3, f2, f1, f, b)
		}

		// Tests with one operand a power of 2 minus 1

		if i == 0 {
			continue // skip zero operand
		}

		b1.Sub(b1, b_one)
		b1.Sub(b1, b_one)
		f1.SubUint64(f1, 2)

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t1, t2), t3)
		f.MulMod(u1, u2, u3)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf1= 0x%x\nf2= 0x%x\nf3= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f1, f2, f3, f, b)
		}

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t1, t3), t2)
		f.MulMod(u1, u3, u2)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf1= 0x%x\nf3= 0x%x\nf2= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f1, f3, f2, f, b)
		}

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t2, t1), t3)
		f.MulMod(u2, u1, u3)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf2= 0x%x\nf1= 0x%x\nf3= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f2, f1, f3, f, b)
		}

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t2, t3), t1)
		f.MulMod(u2, u3, u1)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf2= 0x%x\nf3= 0x%x\nf1= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f2, f3, f1, f, b)
		}

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t3, t1), t2)
		f.MulMod(u3, u1, u2)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf3= 0x%x\nf1= 0x%x\nf2= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f3, f1, f2, f, b)
		}

		set3Big(b1, b2, b3, t1, t2, t3)
		set3Int(f1, f2, f3, u1, u2, u3)

		b.Mod(b.Mul(t3, t2), t1)
		f.MulMod(u3, u2, u1)

		if !checkEq(b, f) {
			t.Fatalf("Expected equality:\nf3= 0x%x\nf2= 0x%x\nf1= 0x%x\n[ op ]==\nf = %x\nb = %x\n", f3, f2, f1, f, b)
		}
	}
}

func TestRandomMulDivOverflow(t *testing.T) {
	for i := 0; i < 10000; i++ {
		b1, f1 := randNums()
		b2, f2 := randNums()
		b3, f3 := randNums()

		f1a, f2a, f3a := f1.Clone(), f2.Clone(), f3.Clone()

		_, overflow := f1.MulDivOverflow(f1, f2, f3)
		if b3.BitLen() == 0 {
			b1.SetInt64(0)
		} else {
			b1.Div(b1.Mul(b1, b2), b3)
		}

		if err := checkOverflow(b1, f1, overflow); err != nil {
			t.Fatal(err)
		}
		if eq := checkEq(b1, f1); !eq {
			t.Fatalf("Expected equality:\nf1= %x\nf2= %x\nf3= %x\n[ - ]==\nf= %x\nb= %x\n", f1a, f2a, f3a, f1, b1)
		}
	}
}

func S256(x *big.Int) *big.Int {
	if x.Cmp(bigtt255) < 0 {
		return x
	} else {
		return new(big.Int).Sub(x, bigtt256)
	}
}

func TestRandomAbs(t *testing.T) {
	for i := 0; i < 10000; i++ {
		b, f1 := randHighNums()
		b2 := S256(big.NewInt(0).Set(b))
		b2.Abs(b2)
		f1a := new(Int).Abs(f1)

		if eq := checkEq(b2, f1a); !eq {
			bf, _ := FromBig(b2)
			t.Fatalf("Expected equality:\nf1= %x\n[ abs ]==\nf = %x\nbf= %x\nb = %x\n", f1, f1a, bf, b2)
		}
	}
}

func TestUdivremQuick(t *testing.T) {
	var (
		u        = []uint64{1, 0, 0, 0, 0}
		expected = new(Int)
	)
	rem := udivrem([]uint64{}, u, &Int{0, 1, 0, 0})
	copy(expected[:], u)
	if !rem.Eq(expected) {
		t.Errorf("Wrong remainder: %x, expected %x", rem, expected)
	}
}

func Test10KRandomSDiv(t *testing.T) { test10KRandom(t, "SDiv") }
func Test10KRandomLsh(t *testing.T)  { test10KRandom(t, "Lsh") }
func Test10KRandomRsh(t *testing.T)  { test10KRandom(t, "Rsh") }
func Test10KRandomSRsh(t *testing.T) { test10KRandom(t, "SRsh") }
func Test10KRandomExp(t *testing.T)  { test10KRandom(t, "Exp") }

func test10KRandom(t *testing.T, name string) {
	tc := lookupBinary(name)
	for i := 0; i < 10000; i++ {
		f1 := randNum()
		f2 := randNum()
		checkBinaryOperation(t, tc.name, tc.u256Fn, tc.bigFn, *f1, *f2)
	}
}

func TestSRsh(t *testing.T) {
	type testCase struct {
		arg string
		n   uint64
	}
	testCases := []testCase{
		{"0xFFFFEEEEDDDDCCCCBBBBAAAA9999888877776666555544443333222211110000", 0},
		{"0xFFFFEEEEDDDDCCCCBBBBAAAA9999888877776666555544443333222211110000", 16},
		{"0xFFFFEEEEDDDDCCCCBBBBAAAA9999888877776666555544443333222211110000", 64},
		{"0xFFFFEEEEDDDDCCCCBBBBAAAA9999888877776666555544443333222211110000", 96},
		{"0xFFFFEEEEDDDDCCCCBBBBAAAA9999888877776666555544443333222211110000", 127},
		{"0xFFFFEEEEDDDDCCCCBBBBAAAA9999888877776666555544443333222211110000", 128},
		{"0xFFFFEEEEDDDDCCCCBBBBAAAA9999888877776666555544443333222211110000", 129},
		{"0xFFFFEEEEDDDDCCCCBBBBAAAA9999888877776666555544443333222211110000", 192},
		{"0x8000000000000000000000000000000000000000000000000000000000000000", 254},
		{"0x8000000000000000000000000000000000000000000000000000000000000000", 255},
		{"0xFFFFEEEEDDDDCCCCBBBBAAAA9999888877776666555544443333222211110000", 256},
		{"0xFFFFEEEEDDDDCCCCBBBBAAAA9999888877776666555544443333222211110000", 300},
		{"0x7FFFEEEEDDDDCCCCBBBBAAAA9999888877776666555544443333222211110000", 16},
		{"0x7FFFEEEEDDDDCCCCBBBBAAAA9999888877776666555544443333222211110000", 256},
	}
	op := lookupBinary("SRsh")
	for _, tc := range testCases {
		arg := MustFromHex(tc.arg)
		n := NewInt(tc.n)
		checkBinaryOperation(t, op.name, op.u256Fn, op.bigFn, *arg, *n)
	}
}

func TestByte(t *testing.T) {
	input, err := FromHex("0x102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f")
	if err != nil {
		t.Fatal(err)
	}
	for i := uint64(0); i < 35; i++ {
		var (
			z     = input.Clone()
			index = NewInt(i)
			have  = z.Byte(index)
			want  = NewInt(i)
		)
		if i >= 32 {
			want.Clear()
		}
		if !have.Eq(want) {
			t.Fatalf("index %d: have %#x want %#x", i, have, want)
		}
		// Also check that we indeed modified the z
		if z != have {
			t.Fatalf("index %d: should return self %v %v", i, z, have)
		}
	}
}

func TestSignExtend(t *testing.T) {
	type testCase struct {
		arg string
		n   uint64
	}
	testCases := []testCase{
		{"10000000000000000", 2},
		{"8080808080808080808080808080808080808080808080808080808080808080", 0},
		{"8080808080808080808080808080808080808080808080808080808080808080", 1},
		{"8080808080808080808080808080808080808080808080808080808080808080", 2},
		{"8080808080808080808080808080808080808080808080808080808080808080", 3},
		{"8080808080808080808080808080808080808080808080808080808080808080", 8},
		{"8080808080808080808080808080808080808080808080808080808080808080", 18},
		{"8080808080808080808080808080808080808080808080808080808080808080", 30},
		{"8080808080808080808080808080808080808080808080808080808080808080", 31},
		{"8080808080808080808080808080808080808080808080808080808080808080", 32},
		{"4040404040404040404040404040404040404040404040404040404040404040", 0},
		{"4040404040404040404040404040404040404040404040404040404040404040", 1},
		{"4040404040404040404040404040404040404040404040404040404040404040", 15},
		{"4040404040404040404040404040404040404040404040404040404040404040", 19},
		{"4040404040404040404040404040404040404040404040404040404040404040", 30},
		{"4040404040404040404040404040404040404040404040404040404040404040", 31},
		{"4040404040404040404040404040404040404040404040404040404040404040", 32},
	}
	op := lookupBinary("ExtendSign")
	for _, tc := range testCases {
		arg := new(Int).SetBytes(hex2Bytes(tc.arg))
		n := new(Int).SetUint64(tc.n)
		checkBinaryOperation(t, op.name, op.u256Fn, op.bigFn, *arg, *n)
	}
}

func TestAddSubUint64(t *testing.T) {
	type testCase struct {
		arg string
		n   uint64
	}
	testCases := []testCase{
		{"0", 1},
		{"1", 0},
		{"1", 1},
		{"1", 3},
		{"0x10000000000000000", 1},
		{"0x100000000000000000000000000000000", 1},
		{"0", 0xffffffffffffffff},
		{"1", 0xffffffff},
		{"0xffffffffffffffff", 1},
		{"0xffffffffffffffff", 0xffffffffffffffff},
		{"0x10000000000000000", 1},
		{"0xfffffffffffffffffffffffffffffffff", 1},
		{"0xfffffffffffffffffffffffffffffffff", 2},
	}

	for i := 0; i < len(testCases); i++ {
		tc := &testCases[i]
		bigArg, _ := new(big.Int).SetString(tc.arg, 0)
		arg, _ := FromBig(bigArg)
		{ // SubUint64
			want, _ := FromBig(u256(new(big.Int).Sub(bigArg, new(big.Int).SetUint64(tc.n))))
			have := new(Int).SetAllOne().SubUint64(arg, tc.n)
			if !have.Eq(want) {
				t.Logf("args: %s, %d\n", tc.arg, tc.n)
				t.Logf("want : %x\n", want)
				t.Logf("have : %x\n\n", have)
				t.Fail()
			}
		}
		{ // AddUint64
			want, _ := FromBig(u256(new(big.Int).Add(bigArg, new(big.Int).SetUint64(tc.n))))
			have := new(Int).AddUint64(arg, tc.n)
			if !have.Eq(want) {
				t.Logf("args: %s, %d\n", tc.arg, tc.n)
				t.Logf("want : %x\n", want)
				t.Logf("have : %x\n\n", have)
				t.Fail()
			}
		}
	}
}

func TestSGT(t *testing.T) {

	x := new(Int).SetBytes(hex2Bytes("fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"))
	y := new(Int).SetBytes(hex2Bytes("00"))
	actual := x.Sgt(y)
	if actual {
		t.Fatalf("Expected %v false", actual)
	}

	x = new(Int).SetBytes(hex2Bytes("00"))
	y = new(Int).SetBytes(hex2Bytes("fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"))
	actual = x.Sgt(y)
	if !actual {
		t.Fatalf("Expected %v true", actual)
	}
}

const (
	// number of bits in a big.Word
	wordBits = 32 << (uint64(^big.Word(0)) >> 63)
)

var (
	tt256m1 = new(big.Int).Sub(bigtt256, big.NewInt(1))
)

// u256 encodes as a 256 bit two's complement number. This operation is destructive.
func u256(x *big.Int) *big.Int {
	return x.And(x, tt256m1)
}

// TestFixedExpReusedArgs tests the cases in Exp() where the arguments (including result) alias the same objects.
func TestFixedExpReusedArgs(t *testing.T) {
	f2 := Int{2, 0, 0, 0}
	f2.Exp(&f2, &f2)
	requireEq(t, big.NewInt(2*2), &f2, "")

	// TODO: This is tested in TestBinOp().
	f3 := Int{3, 0, 0, 0}
	f4 := Int{4, 0, 0, 0}
	f3.Exp(&f4, &f3)
	requireEq(t, big.NewInt(4*4*4), &f3, "")

	// TODO: This is tested in TestBinOp().
	f5 := Int{5, 0, 0, 0}
	f6 := Int{6, 0, 0, 0}
	f6.Exp(&f6, &f5)
	requireEq(t, big.NewInt(6*6*6*6*6), &f6, "")

	f3 = Int{3, 0, 0, 0}
	fr := new(Int).Exp(&f3, &f3)
	requireEq(t, big.NewInt(3*3*3), fr, "")
}

func TestPaddingRepresentation(t *testing.T) {
	a := big.NewInt(0xFF0AFcafe)
	aa := new(Int).SetUint64(0xFF0afcafe)
	bb := new(Int).SetBytes(a.Bytes())
	if !aa.Eq(bb) {
		t.Fatal("aa != bb")
	}

	check := func(padded []byte, expectedHex string) {
		if expected := hex2Bytes(expectedHex); !bytes.Equal(padded, expected) {
			t.Errorf("incorrect padded bytes: %x, expected: %x", padded, expected)
		}
	}

	check(aa.PaddedBytes(32), "0000000000000000000000000000000000000000000000000000000ff0afcafe")
	check(aa.PaddedBytes(20), "0000000000000000000000000000000ff0afcafe")
	check(aa.PaddedBytes(40), "00000000000000000000000000000000000000000000000000000000000000000000000ff0afcafe")

	bytearr := hex2Bytes("0e320219838e859b2f9f18b72e3d4073ca50b37d")
	a = new(big.Int).SetBytes(bytearr)
	aa = new(Int).SetBytes(bytearr)
	bb = new(Int).SetBytes(a.Bytes())
	if !aa.Eq(bb) {
		t.Fatal("aa != bb")
	}

	check(aa.PaddedBytes(32), "0000000000000000000000000e320219838e859b2f9f18b72e3d4073ca50b37d")
	check(aa.PaddedBytes(20), "0e320219838e859b2f9f18b72e3d4073ca50b37d")
	check(aa.PaddedBytes(40), "00000000000000000000000000000000000000000e320219838e859b2f9f18b72e3d4073ca50b37d")
}

func TestWriteToSlice(t *testing.T) {
	x1 := hex2Bytes("fe7fb0d1f59dfe9492ffbf73683fd1e870eec79504c60144cc7f5fc2bad1e611")

	a := big.NewInt(0).SetBytes(x1)
	fa, _ := FromBig(a)

	dest := make([]byte, 32)
	fa.WriteToSlice(dest)
	if !bytes.Equal(dest, x1) {
		t.Errorf("got %x, expected %x", dest, x1)
	}

	fb := new(Int)
	exp := make([]byte, 32)
	fb.WriteToSlice(dest)
	if !bytes.Equal(dest, exp) {
		t.Errorf("got %x, expected %x", dest, exp)
	}
	// a too small buffer
	// Should fill the lower parts, masking upper bytes
	exp = hex2Bytes("683fd1e870eec79504c60144cc7f5fc2bad1e611")
	dest = make([]byte, 20)
	fa.WriteToSlice(dest)
	if !bytes.Equal(dest, exp) {
		t.Errorf("got %x, expected %x", dest, exp)
	}

	// a too large buffer, already filled with stuff
	// Should fill the leftmost 32 bytes, not touch the other things
	dest = hex2Bytes("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	exp = hex2Bytes("fe7fb0d1f59dfe9492ffbf73683fd1e870eec79504c60144cc7f5fc2bad1e611ffffffffffffffff")

	fa.WriteToSlice(dest)
	if !bytes.Equal(dest, exp) {
		t.Errorf("got %x, expected %x", dest, x1)
	}

	// an empty slice, no panics please
	dest = []byte{}
	exp = []byte{}

	fa.WriteToSlice(dest)
	if !bytes.Equal(dest, exp) {
		t.Errorf("got %x, expected %x", dest, x1)
	}

}
func TestInt_WriteToArray(t *testing.T) {
	x1 := hex2Bytes("0000000000000000000000000000d1e870eec79504c60144cc7f5fc2bad1e611")
	a := big.NewInt(0).SetBytes(x1)
	fa, _ := FromBig(a)

	{
		dest := [20]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
		fa.WriteToArray20(&dest)
		exp := hex2Bytes("0000d1e870eec79504c60144cc7f5fc2bad1e611")
		if !bytes.Equal(dest[:], exp) {
			t.Errorf("got %x, expected %x", dest, exp)
		}

	}

	{
		dest := [32]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
		fa.WriteToArray32(&dest)
		exp := hex2Bytes("0000000000000000000000000000d1e870eec79504c60144cc7f5fc2bad1e611")
		if !bytes.Equal(dest[:], exp) {
			t.Errorf("got %x, expected %x", dest, exp)
		}

	}
}

type gethAddress [20]byte

// SetBytes sets the address to the value of b.
// If b is larger than len(a) it will panic.
func (a *gethAddress) setBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-20:]
	}
	copy(a[20-len(b):], b)
}

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func bytesToAddress(b []byte) gethAddress {
	var a gethAddress
	a.setBytes(b)
	return a
}

type gethHash [32]byte

// SetBytes sets the address to the value of b.
// If b is larger than len(a) it will panic.
func (a *gethHash) setBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-32:]
	}
	copy(a[32-len(b):], b)
}

// BytesToHash returns gethHash with value b.
// If b is larger than len(h), b will be cropped from the left.
func bytesToHash(b []byte) gethHash {
	var a gethHash
	a.setBytes(b)
	return a
}

func TestByteRepresentation(t *testing.T) {
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
		bytearr := hex2Bytes(tt)
		// big.Int -> address, hash
		a := big.NewInt(0).SetBytes(bytearr)
		want20 := bytesToAddress(a.Bytes())
		want32 := bytesToHash(a.Bytes())

		// uint256.Int -> address
		b := new(Int).SetBytes(bytearr)
		have20 := gethAddress(b.Bytes20())
		have32 := gethHash(b.Bytes32())

		if have, want := want20, have20; have != want {
			t.Errorf("test %d: have %x want %x", i, have, want)
		}
		if have, want := want32, have32; have != want {
			t.Errorf("test %d: have %x want %x", i, have, want)
		}
	}
}

func testLog10(t *testing.T, z *Int) {
	want := uint(len(z.Dec()))
	if want > 0 {
		want--
	}
	if have := z.Log10(); have != want {
		t.Errorf("%s (%s): have %v want %v", z.Hex(), z.Dec(), have, want)
	}
}

func TestLog10(t *testing.T) {
	testLog10(t, new(Int))
	for i := uint(0); i < 255; i++ {
		z := NewInt(1)
		testLog10(t, z.Lsh(z, i))
		if i != 0 {
			testLog10(t, new(Int).SubUint64(z, 1))
		}
	}
	z := NewInt(1)
	ten := NewInt(10)
	for i := uint(0); i < 80; i++ {
		testLog10(t, z.Mul(z, ten))
		testLog10(t, new(Int).SubUint64(z, 1))
	}
}

func FuzzLog10(f *testing.F) {
	f.Fuzz(func(t *testing.T, aa, bb, cc, dd uint64) {
		testLog10(t, &Int{aa, bb, cc, dd})
	})
}

func BenchmarkLog10(b *testing.B) {
	var u256Ints []*Int
	var bigints []*big.Int

	for i := uint(0); i < 255; i++ {
		a := NewInt(1)
		a.Lsh(a, i)
		u256Ints = append(u256Ints, a)
		bigints = append(bigints, a.ToBig())
	}
	b.Run("Log10/uint256", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, z := range u256Ints {
				_ = z.Log10()
			}
		}
	})
	b.Run("Log10/big", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, z := range bigints {
				f, _ := new(big.Float).SetInt(z).Float64()
				_ = int(math.Log10(f))
			}
		}
	})
}

func TestCmpUint64(t *testing.T) {
	check := func(z *Int, x uint64) {
		want := z.ToBig().Cmp(new(big.Int).SetUint64(x))
		have := z.CmpUint64(x)
		if have != want {
			t.Errorf("%s.CmpUint64( %x ) : have %v want %v", z.Hex(), x, have, want)
		}
	}
	for i := uint(0); i < 255; i++ {
		z := NewInt(1)
		z.Lsh(z, i)
		check(z, new(Int).Set(z).Uint64())                               // z, z
		check(z, new(Int).AddUint64(z, 1).Uint64())                      // z, z + 1
		check(z, new(Int).SubUint64(z, 1).Uint64())                      // z, z - 1
		check(z, new(big.Int).Rsh(new(Int).Set(z).ToBig(), 64).Uint64()) // z, z >> 64
	}
}

func TestCmpBig(t *testing.T) {
	check := func(z *Int, x *big.Int) {
		want := z.ToBig().Cmp(x)
		have := z.CmpBig(x)
		if have != want {
			t.Errorf("%s.CmpBig( %x ) : have %v want %v", z.Hex(), x, have, want)
		}
	}
	for i := uint(0); i < 255; i++ {
		z := NewInt(1)
		z.Lsh(z, i)
		check(z, new(Int).Set(z).ToBig())                        // z, z
		check(z, new(Int).AddUint64(z, 1).ToBig())               // z, z + 1
		check(z, new(Int).SubUint64(z, 1).ToBig())               // z, z - 1
		check(z, new(big.Int).Neg(new(Int).Set(z).ToBig()))      // z, -z
		check(z, new(big.Int).Lsh(new(Int).Set(z).ToBig(), 256)) // z, z << 256
	}
}

func testSetFromDecForFuzzing(tc string) error {
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
	if _, err := a.Value(); err != nil {
		return fmt.Errorf("fail to Value() %s, got err %s", tc, err)
	}
	return nil
}

func FuzzSetString(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 512 {
			return
		}
		if err := testSetFromDecForFuzzing(string(data)); err != nil {
			t.Fatal(err)
		}
	})
}
