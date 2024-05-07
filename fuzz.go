// uint256: Fixed size 256-bit math library
// Copyright 2020 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

//go:build gofuzz
// +build gofuzz

package uint256

import (
	"fmt"
	"math/big"
)

func Fuzz(data []byte) int {
	if len(data) < 32 {
		return 0
	}
	switch {
	case len(data) < 64:
		return fuzzUnaryOp(data) // needs 32 byte
	case len(data) < 96:
		return fuzzBinaryOp(data) // needs 64 byte
	case len(data) < 128:
		return fuzzTernaryOp(data) // needs 96 byte
	}
	// Too large input
	return -1
}

func fuzzUnaryOp(data []byte) int {
	var x Int
	x.SetBytes(data[0:32])
	checkUnaryOp((*Int).Sqrt, (*big.Int).Sqrt, x)
	return 1
}

func fuzzBinaryOp(data []byte) int {
	var x, y Int
	x.SetBytes(data[0:32])
	y.SetBytes(data[32:])
	if !y.IsZero() { // uDivrem
		checkDualArgOp((*Int).Div, (*big.Int).Div, x, y)
		checkDualArgOp((*Int).Mod, (*big.Int).Mod, x, y)
	}
	{ // opMul
		checkDualArgOp((*Int).Mul, (*big.Int).Mul, x, y)
	}
	{ // opLsh
		lsh := func(z, x, y *Int) *Int {
			return z.Lsh(x, uint(y[0]))
		}
		bigLsh := func(z, x, y *big.Int) *big.Int {
			n := uint(y.Uint64())
			if n > 256 {
				n = 256
			}
			return z.Lsh(x, n)
		}
		checkDualArgOp(lsh, bigLsh, x, y)
	}
	{ // opAdd
		checkDualArgOp((*Int).Add, (*big.Int).Add, x, y)
	}
	{ // opSub
		checkDualArgOp((*Int).Sub, (*big.Int).Sub, x, y)
	}
	return 1
}

func fuzzTernaryOp(data []byte) int {
	var x, y, z Int
	x.SetBytes(data[:32])
	y.SetBytes(data[32:64])
	z.SetBytes(data[64:])
	if z.IsZero() {
		return 0
	}

	{ // mulMod
		checkThreeArgOp(intMulMod, bigintMulMod, x, y, z)
	}
	{ // addMod
		checkThreeArgOp(intAddMod, bigintAddMod, x, y, z)
	}
	{ // mulDiv
		checkThreeArgOp(intMulDiv, bigintMulDiv, x, y, z)
	}
	return 1
}

// Test SetFromDecimal
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

func FuzzSetString(data []byte) int {
	if len(data) > 512 {
		// Too large, makes no sense
		return -1
	}
	if err := testSetFromDecForFuzzing(string(data)); err != nil {
		panic(err)
	}
	return 1
}
