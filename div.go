// uint256: Fixed size 256-bit math library
// Copyright 2020 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

package uint256

import "math/bits"

// reciprocal2by1 computes <^d, ^0> / d.
func reciprocal2by1(d uint64) uint64 {
	reciprocal, _ := bits.Div64(^d, ^uint64(0), d)
	return reciprocal
}

// reciprocal3by2 computes (2^192 - 1) / <d1, d0> - 2^64 for a normalized d1.
// The 2-by-1 division of <^d1, ^d0> by d1 is at most 2 too big.
func reciprocal3by2(d1, d0 uint64) uint64 {
	v, r := bits.Div64(^d1, ^d0, d1)
	ph, pl := bits.Mul64(v, d0)
	if r < ph {
		v--
	}
	var carry uint64
	if d0 != 0 {
		carry = 1
	}
	s, c := bits.Add64(r, d1, carry)
	if c == 0 && (s < ph || (s == ph && d0-1 < pl)) {
		v--
	}
	return v
}

// udivrem3by2 divides <u2, u1, u0> / <d1, d0> for a normalized d1 and <u2, u1> < <d1, d0>.
// It uses the provided reciprocal3by2 of d.
// Implementation ported from https://github.com/chfast/intx and is based on
// "Improved division by invariant integers", Algorithm 5.
func udivrem3by2(u2, u1, u0, d1, d0, reciprocal uint64) (quot, r1, r0 uint64) {
	qh, ql := bits.Mul64(reciprocal, u2)
	ql, carry := bits.Add64(ql, u1, 0)
	qh, _ = bits.Add64(qh, u2, carry)

	r1 = u1 - qh*d1
	th, tl := bits.Mul64(d0, qh)
	r0, borrow := bits.Sub64(u0, tl, 0)
	r1, _ = bits.Sub64(r1, th, borrow)
	r0, borrow = bits.Sub64(r0, d0, 0)
	r1, _ = bits.Sub64(r1, d1, borrow)
	qh++

	if r1 >= ql {
		qh--
		r0, carry = bits.Add64(r0, d0, 0)
		r1, _ = bits.Add64(r1, d1, carry)
	}

	if r1 > d1 || (r1 == d1 && r0 >= d0) {
		qh++
		r0, borrow = bits.Sub64(r0, d0, 0)
		r1, _ = bits.Sub64(r1, d1, borrow)
	}

	return qh, r1, r0
}

// udivrem2by1 divides <uh, ul> / d and produces both quotient and remainder.
// It uses the provided d's reciprocal.
// Implementation ported from https://github.com/chfast/intx and is based on
// "Improved division by invariant integers", Algorithm 4.
func udivrem2by1(uh, ul, d, reciprocal uint64) (quot, rem uint64) {
	qh, ql := bits.Mul64(reciprocal, uh)
	ql, carry := bits.Add64(ql, ul, 0)
	qh, _ = bits.Add64(qh, uh, carry)
	qh++

	r := ul - qh*d

	if r > ql {
		qh--
		r += d
	}

	if r >= d {
		qh++
		r -= d
	}

	return qh, r
}

// umod sets z = x % y for y != 0 and x > y. Unlike udivrem it keeps the operands
// in registers: one Knuth step per divisor size, the hardware division for one-limb divisors.
func umod(z, x, y *Int) *Int {
	switch {
	case y[3] != 0:
		if y[3]>>63 != 0 { // x < 2^256 <= 2y
			return z.Sub(x, y)
		}
		s := uint(bits.LeadingZeros64(y[3]))
		r3, r2, r1, r0 := rem5by4(x[3]>>(64-s), x[3]<<s|x[2]>>(64-s), x[2]<<s|x[1]>>(64-s), x[1]<<s|x[0]>>(64-s), x[0]<<s,
			y[3]<<s|y[2]>>(64-s), y[2]<<s|y[1]>>(64-s), y[1]<<s|y[0]>>(64-s), y[0]<<s)
		z[0], z[1], z[2], z[3] = r0>>s|r1<<(64-s), r1>>s|r2<<(64-s), r2>>s|r3<<(64-s), r3>>s
	case y[2] != 0:
		s := uint(bits.LeadingZeros64(y[2]))
		v2, v1, v0 := y[2]<<s|y[1]>>(64-s), y[1]<<s|y[0]>>(64-s), y[0]<<s
		r2, r1, r0 := rem4by3(x[3]>>(64-s), x[3]<<s|x[2]>>(64-s), x[2]<<s|x[1]>>(64-s), x[1]<<s|x[0]>>(64-s), v2, v1, v0)
		r2, r1, r0 = rem4by3(r2, r1, r0, x[0]<<s, v2, v1, v0)
		z[0], z[1], z[2], z[3] = r0>>s|r1<<(64-s), r1>>s|r2<<(64-s), r2>>s, 0
	case y[1] != 0:
		s := uint(bits.LeadingZeros64(y[1]))
		v1, v0 := y[1]<<s|y[0]>>(64-s), y[0]<<s
		r1, r0 := rem3by2(x[3]>>(64-s), x[3]<<s|x[2]>>(64-s), x[2]<<s|x[1]>>(64-s), v1, v0)
		r1, r0 = rem3by2(r1, r0, x[1]<<s|x[0]>>(64-s), v1, v0)
		r1, r0 = rem3by2(r1, r0, x[0]<<s, v1, v0)
		z[0], z[1], z[2], z[3] = r0>>s|r1<<(64-s), r1>>s, 0, 0
	default:
		d := y[0]
		r := x[3] % d
		_, r = bits.Div64(r, x[2], d)
		_, r = bits.Div64(r, x[1], d)
		_, r = bits.Div64(r, x[0], d)
		z.SetUint64(r)
	}
	return z
}

// The remNbyM functions are one Knuth step: they reduce an N-limb u by a normalized M-limb v,
// where the top M limbs of u are below v. The estimate q satisfies q <= qhat <= q+2, so
// u - qhat*v lies in [-2v, v): its top limb is non-zero while it is negative, and at most
// two add-backs fix it.

func rem5by4(u4, u3, u2, u1, u0, v3, v2, v1, v0 uint64) (uint64, uint64, uint64, uint64) {
	q, _ := bits.Div64(u4, u3, v3) // the only caller has u4 < 2^63 <= v3
	h0, l0 := bits.Mul64(q, v0)
	h1, l1 := bits.Mul64(q, v1)
	h2, l2 := bits.Mul64(q, v2)
	h3, l3 := bits.Mul64(q, v3)
	p1, c := bits.Add64(l1, h0, 0)
	p2, c := bits.Add64(l2, h1, c)
	p3, c := bits.Add64(l3, h2, c)
	var b uint64
	u0, b = bits.Sub64(u0, l0, 0)
	u1, b = bits.Sub64(u1, p1, b)
	u2, b = bits.Sub64(u2, p2, b)
	u3, b = bits.Sub64(u3, p3, b)
	u4 -= h3 + c + b
	for u4 != 0 {
		u0, c = bits.Add64(u0, v0, 0)
		u1, c = bits.Add64(u1, v1, c)
		u2, c = bits.Add64(u2, v2, c)
		u3, c = bits.Add64(u3, v3, c)
		u4 += c
	}
	return u3, u2, u1, u0
}

func rem4by3(u3, u2, u1, u0, v2, v1, v0 uint64) (uint64, uint64, uint64) {
	q := ^uint64(0)
	if u3 < v2 {
		q, _ = bits.Div64(u3, u2, v2)
	}
	h0, l0 := bits.Mul64(q, v0)
	h1, l1 := bits.Mul64(q, v1)
	h2, l2 := bits.Mul64(q, v2)
	p1, c := bits.Add64(l1, h0, 0)
	p2, c := bits.Add64(l2, h1, c)
	var b uint64
	u0, b = bits.Sub64(u0, l0, 0)
	u1, b = bits.Sub64(u1, p1, b)
	u2, b = bits.Sub64(u2, p2, b)
	u3 -= h2 + c + b
	for u3 != 0 {
		u0, c = bits.Add64(u0, v0, 0)
		u1, c = bits.Add64(u1, v1, c)
		u2, c = bits.Add64(u2, v2, c)
		u3 += c
	}
	return u2, u1, u0
}

func rem3by2(u2, u1, u0, v1, v0 uint64) (uint64, uint64) {
	q := ^uint64(0)
	if u2 < v1 {
		q, _ = bits.Div64(u2, u1, v1)
	}
	h0, l0 := bits.Mul64(q, v0)
	h1, l1 := bits.Mul64(q, v1)
	p1, c := bits.Add64(l1, h0, 0)
	var b uint64
	u0, b = bits.Sub64(u0, l0, 0)
	u1, b = bits.Sub64(u1, p1, b)
	u2 -= h1 + c + b
	for u2 != 0 {
		u0, c = bits.Add64(u0, v0, 0)
		u1, c = bits.Add64(u1, v1, c)
		u2 += c
	}
	return u1, u0
}
