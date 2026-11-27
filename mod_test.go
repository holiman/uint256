// uint256: Fixed size 256-bit math library
// Copyright 2026 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

package uint256

import (
	"math/big"
	"math/rand"
	"testing"
)

func reciprocalToBig(mu [5]uint64) *big.Int {
	z := new(big.Int)
	for i := len(mu) - 1; i >= 0; i-- {
		z.Lsh(z, 64)
		z.Add(z, new(big.Int).SetUint64(mu[i]))
	}
	return z
}

func reciprocalTestModuli() []Int {
	const max = ^uint64(0)

	moduli := []Int{
		{0, 0, 0, 1},
		{1, 0, 0, 1},
		{max, 0, 0, 1},
		{1 << 63, 0, 0, 1},
		{max, max, max, 1 << 63},
		{0, 0, 0, 1 << 63},
		{1, 0, 0, 1 << 63},
		{max, max, max, max},
	}

	// Exercise every normalization boundary, including 64n+1-bit moduli.
	for bit := uint(192); bit < 256; bit++ {
		power := new(Int).Lsh(new(Int).SetUint64(1), bit)
		moduli = append(moduli, *power)
		if bit > 192 {
			moduli = append(moduli, *new(Int).SubUint64(power, 1))
		}
		moduli = append(moduli, *new(Int).AddUint64(power, 1))

		// This family normalizes to 2^255 + 2^126, which saturates the
		// initial reciprocal seed.
		var correction Int
		correction[bit/64] |= uint64(1) << (bit % 64)
		low := bit - 129
		correction[low/64] |= uint64(1) << (low % 64)
		moduli = append(moduli, correction)
	}

	rnd := rand.New(rand.NewSource(1)) // #nosec G404 -- deterministic test corpus.
	for i := 0; i < 1024; i++ {
		moduli = append(moduli, Int{rnd.Uint64(), rnd.Uint64(), rnd.Uint64(), rnd.Uint64() | 1})
	}
	// Vary every lower limb while retaining the saturated normalized seed.
	saturated := rand.New(rand.NewSource(2)) // #nosec G404 -- deterministic test corpus.
	for i := 0; i < 64; i++ {
		moduli = append(moduli, Int{
			saturated.Uint64(), saturated.Uint64(), saturated.Uint64(),
			1<<63 | (saturated.Uint64() & 0xffffffff),
		})
	}
	return moduli
}

func TestReciprocalApproximationBounds(t *testing.T) {
	max := new(big.Int).Lsh(big.NewInt(1), 512)
	max.Sub(max, big.NewInt(1))
	two := big.NewInt(2)

	for _, m := range reciprocalTestModuli() {
		got := reciprocalToBig(Reciprocal(&m))
		want := new(big.Int).Quo(new(big.Int).Set(max), m.ToBig())
		lower := new(big.Int).Sub(new(big.Int).Set(want), two)

		// The downward-rounded Newton candidate never exceeds the Barrett
		// reciprocal floor((2^512-1)/m). Alignment can lose one more unit.
		if got.Cmp(lower) < 0 || got.Cmp(want) > 0 {
			t.Fatalf("Reciprocal(%#x) = %#x, want %#x or %#x", &m, got, want, lower)
		}
	}
}

func TestReciprocalCorrectionBoundary(t *testing.T) {
	max := new(big.Int).Lsh(big.NewInt(1), 512)
	max.Sub(max, big.NewInt(1))

	// These moduli all normalize to 2^255 + 2^126, where the initial
	// 32-bit reciprocal estimate is saturated and needs correction to preserve
	// the Barrett upper bound.
	for top := uint(192); top < 256; top++ {
		var m Int
		m[top/64] |= uint64(1) << (top % 64)
		low := top - 129
		m[low/64] |= uint64(1) << (low % 64)

		got := reciprocalToBig(Reciprocal(&m))
		want := new(big.Int).Quo(new(big.Int).Set(max), m.ToBig())
		if got.Cmp(want) > 0 {
			t.Fatalf("Reciprocal(%#x) = %#x, exceeds %#x", &m, got, want)
		}
	}
}

func TestReciprocalSaturatedSeedBounds(t *testing.T) {
	max := new(big.Int).Lsh(big.NewInt(1), 512)
	max.Sub(max, big.NewInt(1))
	two := big.NewInt(2)
	rnd := rand.New(rand.NewSource(2)) // #nosec G404 -- deterministic test corpus.

	// The top 32 normalized bits are exactly 0x80000000 for every value
	// here, so they exercise the saturated-seed correction path.
	for i := 0; i < 1024; i++ {
		a := rnd.Uint32()
		if a == 0 {
			a = 1
		}
		m := Int{0, uint64(a) << 32, 0, 1 << 63}
		got := reciprocalToBig(Reciprocal(&m))
		want := new(big.Int).Quo(new(big.Int).Set(max), m.ToBig())
		lower := new(big.Int).Sub(new(big.Int).Set(want), two)
		if got.Cmp(lower) < 0 || got.Cmp(want) > 0 {
			t.Fatalf("Reciprocal(%#x) = %#x, want %#x or %#x", &m, got, want, lower)
		}
	}
}

func TestMulModWithReciprocalBoundaryDifferential(t *testing.T) {
	const max = ^uint64(0)
	operands := []Int{
		{},
		{1},
		{max, max, max, max},
		{max - 1, max, max, max},
		{max, 0, 0, 0},
		{0, max, 0, 0},
		{0, 0, max, 0},
		{0, 0, 0, max},
		{1, 0, 0, 1 << 63},
		{max, max, max, 1 << 63},
	}

	for modulusIndex, modulus := range reciprocalTestModuli() {
		mu := Reciprocal(&modulus)
		values := append([]Int(nil), operands...)
		values = append(values, *new(Int).SubUint64(&modulus, 1), modulus)
		modulusBig := modulus.ToBig()

		for xIndex, x := range values {
			for yIndex, y := range values {
				want := new(big.Int).Mul(x.ToBig(), y.ToBig())
				want.Mod(want, modulusBig)

				var got Int
				got.MulModWithReciprocal(&x, &y, &modulus, &mu)
				if got.ToBig().Cmp(want) != 0 {
					t.Fatalf("modulus %d, operands %d and %d: MulModWithReciprocal(%#x, %#x, %#x) = %#x, want %#x", modulusIndex, xIndex, yIndex, &x, &y, &modulus, &got, want)
				}
			}
		}
	}
}
