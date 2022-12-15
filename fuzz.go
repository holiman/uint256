// uint256: Fixed size 256-bit math library
// Copyright 2020 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

//go:build gofuzz
// +build gofuzz

package uint256

import (
	"fmt"
	"math/big"
	"reflect"
	"runtime"
	"strings"
)

const (
	opUdivrem = iota
	opMul
	opLsh
	opAdd
	opSub
	opMulmod
)

type opUnaryArgFunc func(*Int, *Int) *Int
type bigUnaryArgFunc func(*big.Int, *big.Int) *big.Int

type opDualArgFunc func(*Int, *Int, *Int) *Int
type bigDualArgFunc func(*big.Int, *big.Int, *big.Int) *big.Int

type opThreeArgFunc func(*Int, *Int, *Int, *Int) *Int
type bigThreeArgFunc func(*big.Int, *big.Int, *big.Int, *big.Int) *big.Int

func crash(op interface{}, msg string, args ...Int) {
	fn := runtime.FuncForPC(reflect.ValueOf(op).Pointer())
	fnName := fn.Name()
	fnFile, fnLine := fn.FileLine(fn.Entry())
	var strArgs []string
	for i, arg := range args {
		strArgs = append(strArgs, fmt.Sprintf("%d: %x", i, &arg))
	}
	panic(fmt.Sprintf("%s\nfor %s (%s:%d)\n%v",
		msg, fnName, fnFile, fnLine, strings.Join(strArgs, "\n")))
}

func checkUnaryOp(op opUnaryArgFunc, bigOp bigUnaryArgFunc, x Int) {
	origX := x
	var result Int
	ret := op(&result, &x)
	if ret != &result {
		crash(op, "returned not the pointer receiver", x)
	}
	if x != origX {
		crash(op, "argument modified", x)
	}
	expected, _ := FromBig(bigOp(new(big.Int), x.ToBig()))
	if result != *expected {
		crash(op, "unexpected result", x)
	}
	// Test again when the receiver is not zero.
	var garbage Int
	garbage.Sub(&garbage, NewInt(1))
	ret = op(&garbage, &x)
	if ret != &garbage {
		crash(op, "returned not the pointer receiver", x)
	}
	if garbage != *expected {
		crash(op, "unexpected result", x)
	}
	// Test again with the receiver aliasing arguments.
	ret = op(&x, &x)
	if ret != &x {
		crash(op, "returned not the pointer receiver", x)
	}
	if x != *expected {
		crash(op, "unexpected result", x)
	}
}

func checkDualArgOp(op opDualArgFunc, bigOp bigDualArgFunc, x, y Int) {
	origX := x
	origY := y

	var result Int
	ret := op(&result, &x, &y)
	if ret != &result {
		crash(op, "returned not the pointer receiver", x, y)
	}
	if x != origX {
		crash(op, "first argument modified", x, y)
	}
	if y != origY {
		crash(op, "second argument modified", x, y)
	}

	expected, _ := FromBig(bigOp(new(big.Int), x.ToBig(), y.ToBig()))
	if result != *expected {
		crash(op, "unexpected result", x, y)
	}

	// Test again when the receiver is not zero.
	var garbage Int
	garbage.Xor(&x, &y)
	ret = op(&garbage, &x, &y)
	if ret != &garbage {
		crash(op, "returned not the pointer receiver", x, y)
	}
	if garbage != *expected {
		crash(op, "unexpected result", x, y)
	}
	if x != origX {
		crash(op, "first argument modified", x, y)
	}
	if y != origY {
		crash(op, "second argument modified", x, y)
	}

	// Test again with the receiver aliasing arguments.
	ret = op(&x, &x, &y)
	if ret != &x {
		crash(op, "returned not the pointer receiver", x, y)
	}
	if x != *expected {
		crash(op, "unexpected result", x, y)
	}

	ret = op(&y, &origX, &y)
	if ret != &y {
		crash(op, "returned not the pointer receiver", x, y)
	}
	if y != *expected {
		crash(op, "unexpected result", x, y)
	}
}

func checkThreeArgOp(op opThreeArgFunc, bigOp bigThreeArgFunc, x, y, z Int) {
	origX := x
	origY := y
	origZ := z

	var result Int
	ret := op(&result, &x, &y, &z)
	if ret != &result {
		crash(op, "returned not the pointer receiver", x, y, z)
	}
	switch {
	case x != origX:
		crash(op, "first argument modified", x, y, z)
	case y != origY:
		crash(op, "second argument modified", x, y, z)
	case z != origZ:
		crash(op, "third argument modified", x, y, z)
	}
	expected, _ := FromBig(bigOp(new(big.Int), x.ToBig(), y.ToBig(), z.ToBig()))
	if have, want := result, *expected; have != want {
		crash(op, fmt.Sprintf("unexpected result: have %v want %v", have, want), x, y, z)
	}

	// Test again when the receiver is not zero.
	var garbage Int
	garbage.Xor(&x, &y)
	ret = op(&garbage, &x, &y, &z)
	if ret != &garbage {
		crash(op, "returned not the pointer receiver", x, y, z)
	}
	if have, want := garbage, *expected; have != want {
		crash(op, fmt.Sprintf("unexpected result: have %v want %v", have, want), x, y, z)
	}
	switch {
	case x != origX:
		crash(op, "first argument modified", x, y, z)
	case y != origY:
		crash(op, "second argument modified", x, y, z)
	case z != origZ:
		crash(op, "third argument modified", x, y, z)
	}

	// Test again with the receiver aliasing arguments.
	ret = op(&x, &x, &y, &z)
	if ret != &x {
		crash(op, "returned not the pointer receiver", x, y, z)
	}
	if have, want := x, *expected; have != want {
		crash(op, fmt.Sprintf("unexpected result: have %v want %v", have, want), x, y, z)
	}

	ret = op(&y, &origX, &y, &z)
	if ret != &y {
		crash(op, "returned not the pointer receiver", x, y, z)
	}
	if y != *expected {
		crash(op, "unexpected result", x, y, z)
	}
	ret = op(&z, &origX, &origY, &z)
	if ret != &z {
		crash(op, "returned not the pointer receiver", x, y, z)
	}
	if z != *expected {
		crash(op, fmt.Sprintf("unexpected result: have %v want %v", z.ToBig(), expected), x, y, z)
	}
}

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

func bigintMulMod(b1, b2, b3, b4 *big.Int) *big.Int {
	return b1.Mod(big.NewInt(0).Mul(b2, b3), b4)
}

func intMulMod(f1, f2, f3, f4 *Int) *Int {
	return f1.MulMod(f2, f3, f4)
}

func bigintAddMod(b1, b2, b3, b4 *big.Int) *big.Int {
	return b1.Mod(big.NewInt(0).Add(b2, b3), b4)
}

func intAddMod(f1, f2, f3, f4 *Int) *Int {
	return f1.AddMod(f2, f3, f4)
}

func bigintMulDiv(b1, b2, b3, b4 *big.Int) *big.Int {
	b1.Mul(b2, b3)
	return b1.Div(b1, b4)
}

func intMulDiv(f1, f2, f3, f4 *Int) *Int {
	f1.MulDivOverflow(f2, f3, f4)
	return f1
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

func FuzzSetString(data []byte) int {
	if len(data) > 512 {
		// Too large, makes no sense
		return -1
	}
	var (
		orig = string(data)
		bi   = new(big.Int)
		z    = new(Int)
	)
	z, haveOk := z.SetString(orig, 10)
	bi, wantOk := bi.SetString(orig, 10)
	// if bigint parsing fail, make sure that we failed too
	if !wantOk {
		if haveOk {
			panic(fmt.Sprintf("parsing status, want ok=%v, have ok=%v. Input: %s", haveOk, wantOk, orig))
		}
		return 1
	}
	// if its a negative number, we should err
	if len(orig) > 0 && (orig[0] == '-') {
		if haveOk {
			panic(fmt.Sprintf("should have errored at negative number: %s", orig))
		}
		return 1
	}
	// if its too large, ignore it also
	if bi.BitLen() > 256 {
		if haveOk {
			panic(fmt.Sprintf("should have errored at number overflow: %s", orig))
		}
		return 1
	}
	// No more reasons not to succeed
	if !haveOk {
		panic(fmt.Sprintf("should have parsed '%s' to '%s', but errored instead", orig, bi.String()))
	}
	// otherwise, make sure that the values are equal
	if z.ToBig().Cmp(bi) != 0 {
		panic(fmt.Sprintf("should have parsed %s to %s, but got %s", orig, bi.String(), z.Dec()))
	}
	// make sure that bigint base 10 string is equal to base10 string
	if z.Dec() != bi.String() {
		panic(fmt.Sprintf("should have parsed %s to %s, but got %s", orig, bi.String(), z.Dec()))
	}
	value, err := z.Value()
	if err != nil {
		panic(fmt.Sprintf("fail to Value() %s, got err %s", bi, err))
	}
	if z.Dec()+"e0" != fmt.Sprint(value) {
		panic(fmt.Sprintf("value of %s did not match base 10 encoding %s", value, z.Dec()))
	}
	return 1
}
