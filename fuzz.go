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
	switch len(data) {
	case 64:
		return fuzzBinaryOp(data)
	case 96:
		return fuzzTernaryOp(data)
	}
	return -1
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

func bigMulMod(b1, b2, b3, b4 *big.Int) *big.Int {
	return b1.Mod(big.NewInt(0).Mul(b2, b3), b4)
}

func intMulMod(f1, f2, f3, f4 *Int) *Int {
	return f1.MulMod(f2, f3, f4)
}

func bigAddMod(b1, b2, b3, b4 *big.Int) *big.Int {
	return b1.Mod(big.NewInt(0).Add(b2, b3), b4)
}

func intAddMod(f1, f2, f3, f4 *Int) *Int {
	return f1.AddMod(f2, f3, f4)
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
		checkThreeArgOp(intMulMod, bigMulMod, x, y, z)
	}
	{ // addMod
		checkThreeArgOp(intAddMod, bigAddMod, x, y, z)
	}
	return 1
}
