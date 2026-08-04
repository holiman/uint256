// uint256: Fixed size 256-bit math library
// Copyright 2026 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

package uint256

import (
	"math/big"
	"math/rand"
	"sync"
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
			name: "product-below-divisor",
			x:    "0x1",
			y:    "0x2",
			d:    "0xd6c5b4a3928170ffeeddccbbaa99887766554433221100ffeeddccbbaa998877",
		},
		{
			name: "product-one-below-divisor",
			x:    "0x1",
			y:    "0xd6c5b4a3928170ffeeddccbbaa99887766554433221100ffeeddccbbaa998876",
			d:    "0xd6c5b4a3928170ffeeddccbbaa99887766554433221100ffeeddccbbaa998877",
		},
		{
			name: "product-equals-divisor",
			x:    "0x1",
			y:    "0xd6c5b4a3928170ffeeddccbbaa99887766554433221100ffeeddccbbaa998877",
			d:    "0xd6c5b4a3928170ffeeddccbbaa99887766554433221100ffeeddccbbaa998877",
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
	for i := 0; i < 10000; i++ {
		x := Int{rng.Uint64(), rng.Uint64(), rng.Uint64(), rng.Uint64()}
		y := Int{rng.Uint64(), rng.Uint64(), rng.Uint64(), rng.Uint64()}
		d := Int{rng.Uint64(), rng.Uint64(), rng.Uint64(), rng.Uint64()}
		if d[3] == 0 {
			d[3] = 1
		}
		plan := NewMulDivPlan(&d)
		checkMulDivOverflowRemWithPlan(t, &x, &y, &d, plan)
	}

	d := Int{rng.Uint64(), rng.Uint64(), rng.Uint64(), rng.Uint64()}
	if d[3] == 0 {
		d[3] = 1
	}
	plan := NewMulDivPlan(&d)
	for i := 0; i < 2000; i++ {
		x := Int{rng.Uint64(), rng.Uint64(), rng.Uint64(), rng.Uint64()}
		y := Int{rng.Uint64(), rng.Uint64(), rng.Uint64(), rng.Uint64()}
		checkMulDivOverflowRemWithPlan(t, &x, &y, &d, plan)
	}
}

func TestMulDivPlanBarrettCorrections(t *testing.T) {
	tests := []struct {
		name          string
		x, y, d       string
		wantDelta     uint64
		wantLowBorrow bool
	}{
		{
			name:          "estimate-exact",
			x:             "0x8f3aa6d8bef36a80069728dc67d9db5621ed4caac044316f9569f9e2cb82822f",
			y:             "0x1a634384d0ba8f10b6666b02a03da270f6bd65cefe8c20dccea06b688be116ca",
			d:             "0xee5d7243409c86374f54e36e7f3627fb903e28f8376a23b81b213e776add09fe",
			wantDelta:     0,
			wantLowBorrow: false,
		},
		{
			name:          "estimate-needs-correction",
			x:             "0x2fde08180e0af6543b5c5d890c2e47aa368baacdb181534b69c9019d510a605a",
			y:             "0x3654123a9d86ed5c966ffb3843db7890c583966ca25897642a77c7253c6fa3c2",
			d:             "0xb6eb43d7b34d5fc018e124f8c8ca446487ed664df51536f943ebf12972065d8e",
			wantDelta:     1,
			wantLowBorrow: false,
		},
		{
			name:          "low-limb-underflow",
			x:             "0x1000000000000000000000000000000000000000000000000",
			y:             "0x1000000000000000000000000000000000000000000000000",
			d:             "0x1000000000000000000000000000000000000000000000000",
			wantDelta:     1,
			wantLowBorrow: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			x := MustFromHex(tc.x)
			y := MustFromHex(tc.y)
			d := MustFromHex(tc.d)
			plan := NewMulDivPlan(d)
			product := new(big.Int).Mul(x.ToBig(), y.ToBig())
			wantQuotient := new(big.Int).Quo(product, d.ToBig())

			estimate := new(big.Int).Rsh(new(big.Int).Set(product), 192)
			estimate.Mul(estimate, bigFromUint64Words(plan.reciprocal[:]))
			estimate.Rsh(estimate, 320)
			if delta := new(big.Int).Sub(wantQuotient, estimate); !delta.IsUint64() || delta.Uint64() != tc.wantDelta {
				t.Fatalf("Barrett quotient estimate is %v below the quotient, want %d", delta, tc.wantDelta)
			}

			base := new(big.Int).Lsh(big.NewInt(1), 320)
			mask := base.Sub(base, big.NewInt(1))
			productLow := new(big.Int).And(product, mask)
			estimateProductLow := new(big.Int).Mul(estimate, d.ToBig())
			estimateProductLow.And(estimateProductLow, mask)
			if gotLowBorrow := productLow.Cmp(estimateProductLow) < 0; gotLowBorrow != tc.wantLowBorrow {
				t.Fatalf("low-limb subtraction borrow = %v, want %v", gotLowBorrow, tc.wantLowBorrow)
			}

			checkMulDivOverflowRemWithPlan(t, x, y, d, plan)
		})
	}
}

func bigFromUint64Words(words []uint64) *big.Int {
	z := new(big.Int)
	for i := len(words) - 1; i >= 0; i-- {
		z.Lsh(z, 64)
		z.Add(z, new(big.Int).SetUint64(words[i]))
	}
	return z
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

func TestMulDivPlanConcurrentUse(t *testing.T) {
	x := MustFromHex("0xf123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9")
	y := MustFromHex("0xcdef0123456789ab9876543210fedcba112233445566778899aabbccddeeff11")
	d := MustFromHex("0xd6c5b4a3928170ffeeddccbbaa99887766554433221100ffeeddccbbaa998877")
	plan := NewMulDivPlan(d)

	var wantQuotient, wantRemainder Int
	_, _, wantOverflow := wantQuotient.MulDivOverflowRem(x, y, d, &wantRemainder)

	const workers = 8
	const iterations = 1000
	var wg sync.WaitGroup
	wg.Add(workers)
	for worker := 0; worker < workers; worker++ {
		go func() {
			defer wg.Done()
			var gotQuotient, gotRemainder Int
			for i := 0; i < iterations; i++ {
				_, _, gotOverflow := gotQuotient.MulDivOverflowRemWithPlan(x, y, plan, &gotRemainder)
				if gotOverflow != wantOverflow || !gotQuotient.Eq(&wantQuotient) || !gotRemainder.Eq(&wantRemainder) {
					t.Errorf("got quotient=%x remainder=%x overflow=%v", &gotQuotient, &gotRemainder, gotOverflow)
					return
				}
			}
		}()
	}
	wg.Wait()
}

func TestMulDivOverflowRemWithPlanAliases(t *testing.T) {
	x := MustFromHex("0xf123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9")
	y := MustFromHex("0xcdef0123456789ab9876543210fedcba112233445566778899aabbccddeeff11")
	d := MustFromHex("0xd6c5b4a3928170ffeeddccbbaa99887766554433221100ffeeddccbbaa998877")
	plan := NewMulDivPlan(d)

	var wantQuotient, wantRemainder Int
	_, _, wantOverflow := wantQuotient.MulDivOverflowRem(x, y, d, &wantRemainder)

	t.Run("quotient-aliases-x", func(t *testing.T) {
		gotQuotient := *x
		var gotRemainder Int
		_, _, gotOverflow := gotQuotient.MulDivOverflowRemWithPlan(&gotQuotient, y, plan, &gotRemainder)
		if gotOverflow != wantOverflow || !gotQuotient.Eq(&wantQuotient) || !gotRemainder.Eq(&wantRemainder) {
			t.Fatalf("got quotient=%x remainder=%x overflow=%v", &gotQuotient, &gotRemainder, gotOverflow)
		}
	})

	t.Run("remainder-aliases-y", func(t *testing.T) {
		gotRemainder := *y
		var gotQuotient Int
		_, _, gotOverflow := gotQuotient.MulDivOverflowRemWithPlan(x, &gotRemainder, plan, &gotRemainder)
		if gotOverflow != wantOverflow || !gotQuotient.Eq(&wantQuotient) || !gotRemainder.Eq(&wantRemainder) {
			t.Fatalf("got quotient=%x remainder=%x overflow=%v", &gotQuotient, &gotRemainder, gotOverflow)
		}
	})

	t.Run("outputs-alias", func(t *testing.T) {
		var got Int
		quotient, remainder, gotOverflow := got.MulDivOverflowRemWithPlan(x, y, plan, &got)
		if quotient != &got || remainder != &got || gotOverflow != wantOverflow || !got.Eq(&wantQuotient) {
			t.Fatalf("got quotient=%x remainder=%x overflow=%v", quotient, remainder, gotOverflow)
		}
	})

	t.Run("direct-outputs-alias", func(t *testing.T) {
		var got Int
		quotient, remainder, gotOverflow := got.MulDivOverflowRem(x, y, d, &got)
		if quotient != &got || remainder != &got || gotOverflow != wantOverflow || !got.Eq(&wantQuotient) {
			t.Fatalf("got quotient=%x remainder=%x overflow=%v", quotient, remainder, gotOverflow)
		}
	})

	t.Run("small-product-outputs-alias", func(t *testing.T) {
		x, y := NewInt(1), NewInt(2)
		var got Int
		quotient, remainder, gotOverflow := got.MulDivOverflowRemWithPlan(x, y, plan, &got)
		if quotient != &got || remainder != &got || gotOverflow || !got.IsZero() {
			t.Fatalf("got quotient=%x remainder=%x overflow=%v", quotient, remainder, gotOverflow)
		}
	})

	t.Run("small-product-direct-outputs-alias", func(t *testing.T) {
		x, y := NewInt(1), NewInt(2)
		var got Int
		quotient, remainder, gotOverflow := got.MulDivOverflowRem(x, y, d, &got)
		if quotient != &got || remainder != &got || gotOverflow || !got.IsZero() {
			t.Fatalf("got quotient=%x remainder=%x overflow=%v", quotient, remainder, gotOverflow)
		}
	})
}

func checkMulDivOverflowRemWithPlan(t *testing.T, x, y, d *Int, plan *MulDivPlan) {
	t.Helper()

	var wantQuotient, wantRemainder big.Int
	if !d.IsZero() {
		product := new(big.Int).Mul(x.ToBig(), y.ToBig())
		wantQuotient.QuoRem(product, d.ToBig(), &wantRemainder)
	}

	var wantQuotientInt, wantRemainderInt Int
	wantOverflow := wantQuotientInt.SetFromBig(&wantQuotient)
	wantRemainderInt.SetFromBig(&wantRemainder)

	var directQuotient, directRemainder Int
	directQ, directR, directOverflow := directQuotient.MulDivOverflowRem(x, y, d, &directRemainder)
	if directQ != &directQuotient || directR != &directRemainder {
		t.Fatalf("unexpected direct result pointers: quotient=%p remainder=%p", directQ, directR)
	}
	if directOverflow != wantOverflow || !directQuotient.Eq(&wantQuotientInt) || !directRemainder.Eq(&wantRemainderInt) {
		t.Fatalf("direct result quotient=%x remainder=%x overflow=%v; want quotient=%x remainder=%x overflow=%v", &directQuotient, &directRemainder, directOverflow, &wantQuotientInt, &wantRemainderInt, wantOverflow)
	}

	var gotQuotient, gotRemainder Int
	quotient, remainder, overflow := gotQuotient.MulDivOverflowRemWithPlan(x, y, plan, &gotRemainder)
	if quotient != &gotQuotient || remainder != &gotRemainder {
		t.Fatalf("unexpected planned result pointers: quotient=%p remainder=%p", quotient, remainder)
	}
	if overflow != wantOverflow || !gotQuotient.Eq(&wantQuotientInt) || !gotRemainder.Eq(&wantRemainderInt) {
		t.Fatalf("planned result quotient=%x remainder=%x overflow=%v; want quotient=%x remainder=%x overflow=%v", &gotQuotient, &gotRemainder, overflow, &wantQuotientInt, &wantRemainderInt, wantOverflow)
	}
}
