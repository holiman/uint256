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
		r3, r2, r1, r0, _ := rem5by4(x[3]>>(64-s), x[3]<<s|x[2]>>(64-s), x[2]<<s|x[1]>>(64-s), x[1]<<s|x[0]>>(64-s), x[0]<<s,
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
		r1, r0, _ := rem3by2(x[3]>>(64-s), x[3]<<s|x[2]>>(64-s), x[2]<<s|x[1]>>(64-s), v1, v0)
		r1, r0, _ = rem3by2(r1, r0, x[1]<<s|x[0]>>(64-s), v1, v0)
		r1, r0, _ = rem3by2(r1, r0, x[0]<<s, v1, v0)
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

// udiv sets z = x / y for y != 0 and x > y, in registers like umod.
func udiv(z, x, y *Int) *Int {
	switch {
	case y[3] != 0:
		if y[3]>>63 != 0 { // x < 2^256 <= 2y
			return z.SetOne()
		}
		s := uint(bits.LeadingZeros64(y[3]))
		_, _, _, _, q := rem5by4(x[3]>>(64-s), x[3]<<s|x[2]>>(64-s), x[2]<<s|x[1]>>(64-s), x[1]<<s|x[0]>>(64-s), x[0]<<s,
			y[3]<<s|y[2]>>(64-s), y[2]<<s|y[1]>>(64-s), y[1]<<s|y[0]>>(64-s), y[0]<<s)
		z[0], z[1], z[2], z[3] = q, 0, 0, 0
	case y[2] != 0:
		s := uint(bits.LeadingZeros64(y[2]))
		v2, v1, v0 := y[2]<<s|y[1]>>(64-s), y[1]<<s|y[0]>>(64-s), y[0]<<s
		r2, r1, r0, q1 := quoRem4by3(x[3]>>(64-s), x[3]<<s|x[2]>>(64-s), x[2]<<s|x[1]>>(64-s), x[1]<<s|x[0]>>(64-s), v2, v1, v0)
		_, _, _, q0 := quoRem4by3(r2, r1, r0, x[0]<<s, v2, v1, v0)
		z[0], z[1], z[2], z[3] = q0, q1, 0, 0
	case y[1] != 0:
		s := uint(bits.LeadingZeros64(y[1]))
		v1, v0 := y[1]<<s|y[0]>>(64-s), y[0]<<s
		r1, r0, q2 := rem3by2(x[3]>>(64-s), x[3]<<s|x[2]>>(64-s), x[2]<<s|x[1]>>(64-s), v1, v0)
		r1, r0, q1 := rem3by2(r1, r0, x[1]<<s|x[0]>>(64-s), v1, v0)
		_, _, q0 := rem3by2(r1, r0, x[0]<<s, v1, v0)
		z[0], z[1], z[2], z[3] = q0, q1, q2, 0
	default:
		d := y[0]
		q3 := x[3] / d
		q2, r := bits.Div64(x[3]%d, x[2], d)
		q1, r := bits.Div64(r, x[1], d)
		q0, _ := bits.Div64(r, x[0], d)
		z[0], z[1], z[2], z[3] = q0, q1, q2, q3
	}
	return z
}

// The remNbyM functions are one Knuth step: they reduce an N-limb u by a normalized M-limb v,
// where the top M limbs of u are below v. The estimate q satisfies q <= qhat <= q+2, so
// u - qhat*v lies in [-2v, v): its top limb is non-zero while it is negative, and at most
// two add-backs fix it.

func rem5by4(u4, u3, u2, u1, u0, v3, v2, v1, v0 uint64) (uint64, uint64, uint64, uint64, uint64) {
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
		q--
		u0, c = bits.Add64(u0, v0, 0)
		u1, c = bits.Add64(u1, v1, c)
		u2, c = bits.Add64(u2, v2, c)
		u3, c = bits.Add64(u3, v3, c)
		u4 += c
	}
	return u3, u2, u1, u0, q
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

// quoRem4by3 is rem4by3 that also returns the quotient digit. It is separate because carrying q through rem4by3 slows umod.
func quoRem4by3(u3, u2, u1, u0, v2, v1, v0 uint64) (uint64, uint64, uint64, uint64) {
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
		q--
		u0, c = bits.Add64(u0, v0, 0)
		u1, c = bits.Add64(u1, v1, c)
		u2, c = bits.Add64(u2, v2, c)
		u3 += c
	}
	return u2, u1, u0, q
}

func rem3by2(u2, u1, u0, v1, v0 uint64) (uint64, uint64, uint64) {
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
		q--
		u0, c = bits.Add64(u0, v0, 0)
		u1, c = bits.Add64(u1, v1, c)
		u2 += c
	}
	return u1, u0, q
}
