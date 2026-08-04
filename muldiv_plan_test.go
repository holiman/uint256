// uint256: Fixed size 256-bit math library
// Copyright 2026 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

package uint256

import (
	"math/big"
	"math/rand"
	"testing"
)

func TestMulDivOverflowRemWithPlan(t *testing.T) {
	tests := []struct {
		name string
		x    string
		y    string
		d    string
	}{
		{
			name: "zero-input",
			x:    "0x0",
			y:    "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			d:    "0xd6c5b4a3928170ffeeddccbbaa99887766554433221100ffeeddccbbaa998877",
		},
		{
			name: "zero-divisor",
			x:    "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			y:    "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			d:    "0x0",
		},
		{
			name: "barrett-underflow-correction",
			x:    "0x1000000000000000000000000000000000000000000000000",
			y:    "0x1000000000000000000000000000000000000000000000000",
			d:    "0x1000000000000000000000000000000000000000000000000",
		},
		{
			name: "five-limb-quotient",
			x:    "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			y:    "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			d:    "0x1000000000000000000000000000000000000000000000000",
		},
		{
			name: "dense-full-width-divisor",
			x:    "0xf123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9",
			y:    "0xcdef0123456789ab9876543210fedcba112233445566778899aabbccddeeff11",
			d:    "0xd6c5b4a3928170ffeeddccbbaa99887766554433221100ffeeddccbbaa998877",
		},
		{
			name: "narrow-divisor-fallback",
			x:    "0xf123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9",
			y:    "0xcdef0123456789ab9876543210fedcba112233445566778899aabbccddeeff11",
			d:    "0xffffffffffffffffffffffffffffffff",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			x := MustFromHex(tc.x)
			y := MustFromHex(tc.y)
			d := MustFromHex(tc.d)
			plan := NewMulDivPlan(d)
			checkMulDivOverflowRemWithPlan(t, x, y, d, plan)
		})
	}

	rng := rand.New(rand.NewSource(1)) // #nosec G404 -- deterministic test data
	for i := 0; i < 2000; i++ {
		x := Int{rng.Uint64(), rng.Uint64(), rng.Uint64(), rng.Uint64()}
		y := Int{rng.Uint64(), rng.Uint64(), rng.Uint64(), rng.Uint64()}
		d := Int{rng.Uint64(), rng.Uint64(), rng.Uint64(), rng.Uint64()}
		if d[3] == 0 {
			d[3] = 1
		}
		plan := NewMulDivPlan(&d)
		checkMulDivOverflowRemWithPlan(t, &x, &y, &d, plan)
	}
}

func TestMulDivPlanCopiesDivisor(t *testing.T) {
	d := MustFromHex("0xd6c5b4a3928170ffeeddccbbaa99887766554433221100ffeeddccbbaa998877")
	original := *d
	plan := NewMulDivPlan(d)
	d.Clear()

	x := MustFromHex("0xf123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9")
	y := MustFromHex("0xcdef0123456789ab9876543210fedcba112233445566778899aabbccddeeff11")
	checkMulDivOverflowRemWithPlan(t, x, y, &original, plan)
}

func checkMulDivOverflowRemWithPlan(t *testing.T, x, y, d *Int, plan *MulDivPlan) {
	t.Helper()

	var wantQuotient, wantRemainder big.Int
	if !d.IsZero() {
		product := new(big.Int).Mul(x.ToBig(), y.ToBig())
		wantQuotient.QuoRem(product, d.ToBig(), &wantRemainder)
	}

	var gotQuotient, gotRemainder Int
	quotient, remainder, overflow := gotQuotient.MulDivOverflowRemWithPlan(x, y, plan, &gotRemainder)
	if quotient != &gotQuotient || remainder != &gotRemainder {
		t.Fatalf("unexpected result pointers: quotient=%p remainder=%p", quotient, remainder)
	}

	var wantQuotientInt, wantRemainderInt Int
	wantOverflow := wantQuotientInt.SetFromBig(&wantQuotient)
	wantRemainderInt.SetFromBig(&wantRemainder)
	if overflow != wantOverflow {
		t.Fatalf("overflow = %v, want %v; x=%x y=%x d=%x", overflow, wantOverflow, x, y, d)
	}
	if !gotQuotient.Eq(&wantQuotientInt) {
		t.Fatalf("quotient = %x, want %x; x=%x y=%x d=%x", &gotQuotient, &wantQuotientInt, x, y, d)
	}
	if !gotRemainder.Eq(&wantRemainderInt) {
		t.Fatalf("remainder = %x, want %x; x=%x y=%x d=%x", &gotRemainder, &wantRemainderInt, x, y, d)
	}
}
