// uint256: Fixed size 256-bit math library
// Copyright 2026 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

package uint256

import (
	"math/big"
	"testing"
)

type addModAliasPattern struct {
	name  string
	roles [4]int // z, x, y, m
}

var addModAliasPatterns = []addModAliasPattern{
	{"distinct", [4]int{0, 1, 2, 3}},
	{"z=x", [4]int{0, 0, 1, 2}},
	{"z=y", [4]int{0, 1, 0, 2}},
	{"z=m", [4]int{0, 1, 2, 0}},
	{"x=y", [4]int{0, 1, 1, 2}},
	{"x=m", [4]int{0, 1, 2, 1}},
	{"y=m", [4]int{0, 1, 2, 2}},
	{"z=x=y", [4]int{0, 0, 0, 1}},
	{"z=x=m", [4]int{0, 0, 1, 0}},
	{"z=y=m", [4]int{0, 1, 0, 0}},
	{"x=y=m", [4]int{0, 1, 1, 1}},
	{"z=x,y=m", [4]int{0, 0, 1, 1}},
	{"z=y,x=m", [4]int{0, 1, 0, 1}},
	{"z=m,x=y", [4]int{0, 1, 1, 0}},
	{"all", [4]int{0, 0, 0, 0}},
}

func checkAddMod(t *testing.T, z, x, y, m *Int) {
	t.Helper()

	wantBig := new(big.Int).Add(x.ToBig(), y.ToBig())
	if m.IsZero() {
		wantBig.SetUint64(0)
	} else {
		wantBig.Mod(wantBig, m.ToBig())
	}
	want, overflow := FromBig(wantBig)
	if overflow {
		t.Fatal("unexpected big.Int overflow")
	}

	got := z.AddMod(x, y, m)
	if got != z {
		t.Fatalf("unexpected result pointer: have %p want %p", got, z)
	}
	if !got.Eq(want) {
		t.Fatalf("AddMod(%#x, %#x, %#x): have %#x want %#x", x, y, m, got, want)
	}
}

func checkAddModAliasPatterns(t *testing.T, values [4]Int) {
	t.Helper()
	for _, pattern := range addModAliasPatterns {
		values := values
		roles := pattern.roles
		t.Run(pattern.name, func(t *testing.T) {
			checkAddMod(t, &values[roles[0]], &values[roles[1]], &values[roles[2]], &values[roles[3]])
		})
	}
}

func checkAddModResultAliases(t *testing.T, x, y, m Int) {
	t.Helper()
	t.Run("distinct", func(t *testing.T) {
		x, y, m := x, y, m
		var z Int
		checkAddMod(t, &z, &x, &y, &m)
	})
	t.Run("z=x", func(t *testing.T) {
		x, y, m := x, y, m
		checkAddMod(t, &x, &x, &y, &m)
	})
	t.Run("z=y", func(t *testing.T) {
		x, y, m := x, y, m
		checkAddMod(t, &y, &x, &y, &m)
	})
	t.Run("z=m", func(t *testing.T) {
		x, y, m := x, y, m
		checkAddMod(t, &m, &x, &y, &m)
	})
}

func TestAddModSmallModulus(t *testing.T) {
	tests := []struct {
		name   string
		values [4]Int
	}{
		{
			name: "zero",
			values: [4]Int{
				{0xfedcba9876543210, 0xfedcba9876543210, 0xfedcba9876543210, 0xfedcba9876543210},
				{0x0123456789abcdef, 0x0123456789abcdef, 0x0123456789abcdef, 0x0123456789abcdef},
				{0x1111111111111111, 0x2222222222222222, 0x3333333333333333, 0x4444444444444444},
				{},
			},
		},
		{
			name: "one",
			values: [4]Int{
				{0xfedcba9876543210, 0xfedcba9876543210, 0xfedcba9876543210, 0xfedcba9876543210},
				{0x0123456789abcdef, 0x0123456789abcdef, 0x0123456789abcdef, 0x0123456789abcdef},
				{0x1111111111111111, 0x2222222222222222, 0x3333333333333333, 0x4444444444444444},
				{1},
			},
		},
		{
			name: "mod64",
			values: [4]Int{
				{0x0123456789abcdef, 0, 0, 0},
				{0xfffffffffffffffe, 0, 0, 0},
				{0x123456789abcdef0, 0, 0, 0},
				{0xfffffffffffffff1, 0, 0, 0},
			},
		},
		{
			name: "mod128-boundary",
			values: [4]Int{
				{0x0123456789abcdef, 0, 0, 0},
				{0, 0, 1, 0}, // 2^128
				{2, 0, 0, 0},
				{1, 0, 1, 0}, // 2^128 + 1
			},
		},
		{
			name: "mod64-wide-unreduced",
			values: [4]Int{
				{0x0123456789abcdef, 0, 0, 0},
				{0, 1 << 63, 0, 0}, // 2^127
				{0, 1 << 63, 0, 0}, // 2^127
				{1, 1, 0, 0},
			},
		},
		{
			name: "mod64-double",
			values: [4]Int{
				{0x0123456789abcdef, 0, 0, 0},
				{0xfffffffffffffffe, 1, 0, 0}, // 2 * (2^64 - 1)
				{},
				{0xffffffffffffffff, 0, 0, 0},
			},
		},
		{
			name: "mod64-at-modulus",
			values: [4]Int{
				{0x0123456789abcdef, 0, 0, 0},
				{0xfffffffffffffff1, 0, 0, 0},
				{0xfffffffffffffff2, 0, 0, 0},
				{0xfffffffffffffff1, 0, 0, 0},
			},
		},
		{
			name: "mod128-double",
			values: [4]Int{
				{0x0123456789abcdef, 0, 0, 0},
				{2, 0, 2, 0}, // 2 * (2^128 + 1)
				{},
				{1, 0, 1, 0},
			},
		},
		{
			name: "mod128-near-reduced",
			values: [4]Int{
				{0x0123456789abcdef, 0, 0, 0},
				{0xfffffffffffffff2, 0xffffffffffffffff, 0, 0},
				{0xfffffffffffffff2, 0xffffffffffffffff, 0, 0},
				{0xfffffffffffffff1, 0xffffffffffffffff, 0, 0},
			},
		},
		{
			name: "mod128-carry",
			values: [4]Int{
				{0x0123456789abcdef, 0, 0, 0},
				{0xfffffffffffffffd, 0xffffffffffffffff, 0, 0},
				{0xfffffffffffffffd, 0xffffffffffffffff, 0, 0},
				{0xffffffffffffffff, 0xffffffffffffffff, 0, 0},
			},
		},
		{
			name: "mod192-double",
			values: [4]Int{
				{0x0123456789abcdef, 0, 0, 0},
				{0xfffffffffffffffe, 0xffffffffffffffff, 0xffffffffffffffff, 1}, // 2 * (2^192 - 1)
				{},
				{0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 0},
			},
		},
		{
			name: "mod192-near-reduced",
			values: [4]Int{
				{0x0123456789abcdef, 0, 0, 0},
				{0xfffffffffffffff2, 0xffffffffffffffff, 0xffffffffffffffff, 0},
				{0xfffffffffffffff2, 0xffffffffffffffff, 0xffffffffffffffff, 0},
				{0xfffffffffffffff1, 0xffffffffffffffff, 0xffffffffffffffff, 0},
			},
		},
		{
			name: "mod192-carry",
			values: [4]Int{
				{0x0123456789abcdef, 0, 0, 0},
				{0xfffffffffffffffd, 0xffffffffffffffff, 0xffffffffffffffff, 0},
				{0xfffffffffffffffd, 0xffffffffffffffff, 0xffffffffffffffff, 0},
				{0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 0},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			checkAddModResultAliases(t, test.values[1], test.values[2], test.values[3])
			checkAddModAliasPatterns(t, test.values)
		})
	}
}

func FuzzAddModSmallModulus(f *testing.F) {
	for _, seed := range [][11]uint64{
		// m == 0.
		{0, 0, 0, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 1, 2, 3, 4},
		// m == 1.
		{1, 0, 0, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 1, 2, 3, 4},
		// m == 2^64 + 1.
		{1, 1, 0, 0, 1, 0, 0, 2, 0, 0, 0},
		// m == 2^128 + 1, x == 2^128, y == 2.
		{1, 0, 1, 0, 0, 1, 0, 2, 0, 0, 0},
		// m == 2^128 + 1, x == 2m.
		{1, 0, 1, 2, 0, 2, 0, 0, 0, 0, 0},
		// m == 2^64 + 1, x == y == 2^127.
		{1, 1, 0, 0, 1 << 63, 0, 0, 0, 1 << 63, 0, 0},
		// x == m and y == m + 1 at each small-modulus width.
		{0xfffffffffffffff1, 0, 0, 0xfffffffffffffff1, 0, 0, 0, 0xfffffffffffffff2, 0, 0, 0},
		{0xfffffffffffffff1, 0xffffffffffffffff, 0, 0xfffffffffffffff1, 0xffffffffffffffff, 0, 0, 0xfffffffffffffff2, 0xffffffffffffffff, 0, 0},
		{0xfffffffffffffff1, 0xffffffffffffffff, 0xffffffffffffffff, 0xfffffffffffffff1, 0xffffffffffffffff, 0xffffffffffffffff, 0, 0xfffffffffffffff2, 0xffffffffffffffff, 0xffffffffffffffff, 0},
		// m == 2^192 - 1 with a carry-producing reduced sum.
		{0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 0xfffffffffffffffd, 0xffffffffffffffff, 0xffffffffffffffff, 0, 0xfffffffffffffffd, 0xffffffffffffffff, 0xffffffffffffffff, 0},
	} {
		f.Add(seed[0], seed[1], seed[2], seed[3], seed[4], seed[5], seed[6], seed[7], seed[8], seed[9], seed[10])
	}

	f.Fuzz(func(t *testing.T,
		m0, m1, m2,
		x0, x1, x2, x3,
		y0, y1, y2, y3 uint64) {
		x := Int{x0, x1, x2, x3}
		y := Int{y0, y1, y2, y3}
		m := Int{m0, m1, m2, 0}

		var z Int
		checkAddMod(t, &z, &x, &y, &m)

		x = Int{x0, x1, x2, x3}
		y = Int{y0, y1, y2, y3}
		m = Int{m0, m1, m2, 0}
		checkAddMod(t, &x, &x, &y, &m)

		x = Int{x0, x1, x2, x3}
		y = Int{y0, y1, y2, y3}
		m = Int{m0, m1, m2, 0}
		checkAddMod(t, &y, &x, &y, &m)

		x = Int{x0, x1, x2, x3}
		y = Int{y0, y1, y2, y3}
		m = Int{m0, m1, m2, 0}
		checkAddMod(t, &m, &x, &y, &m)
	})
}
