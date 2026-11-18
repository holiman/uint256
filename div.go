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
		if s, c := bits.Add64(r, d1, 1); c == 0 && (s < ph || (s == ph && d0-1 < pl)) {
			v--
		}
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

// umod512 sets z = p % m for m < 2^192, m != 0, in registers like umod.
// The divisor reciprocal replaces a hardware division in each of the many Knuth steps.
func umod512(z *Int, p *[8]uint64, m *Int) *Int {
	switch {
	case m[2] != 0:
		s := uint(bits.LeadingZeros64(m[2]))
		v2, v1, v0 := m[2]<<s|m[1]>>(64-s), m[1]<<s|m[0]>>(64-s), m[0]<<s
		rec := reciprocal3by2(v2, v1)
		r2, r1, r0 := rem4by3r(p[7]>>(64-s), p[7]<<s|p[6]>>(64-s), p[6]<<s|p[5]>>(64-s), p[5]<<s|p[4]>>(64-s), v2, v1, v0, rec)
		for i := 4; i > 0; i-- {
			r2, r1, r0 = rem4by3r(r2, r1, r0, p[i]<<s|p[i-1]>>(64-s), v2, v1, v0, rec)
		}
		r2, r1, r0 = rem4by3r(r2, r1, r0, p[0]<<s, v2, v1, v0, rec)
		z[0], z[1], z[2], z[3] = r0>>s|r1<<(64-s), r1>>s|r2<<(64-s), r2>>s, 0
	case m[1] != 0:
		s := uint(bits.LeadingZeros64(m[1]))
		v1, v0 := m[1]<<s|m[0]>>(64-s), m[0]<<s
		rec := reciprocal3by2(v1, v0)
		_, r1, r0 := udivrem3by2(p[7]>>(64-s), p[7]<<s|p[6]>>(64-s), p[6]<<s|p[5]>>(64-s), v1, v0, rec)
		for i := 5; i > 0; i-- {
			_, r1, r0 = udivrem3by2(r1, r0, p[i]<<s|p[i-1]>>(64-s), v1, v0, rec)
		}
		_, r1, r0 = udivrem3by2(r1, r0, p[0]<<s, v1, v0, rec)
		z[0], z[1], z[2], z[3] = r0>>s|r1<<(64-s), r1>>s, 0, 0
	default:
		s := uint(bits.LeadingZeros64(m[0]))
		d := m[0] << s
		rec := reciprocal2by1(d)
		r := p[7] >> (64 - s)
		for i := 7; i > 0; i-- {
			_, r = udivrem2by1(r, p[i]<<s|p[i-1]>>(64-s), d, rec)
		}
		_, r = udivrem2by1(r, p[0]<<s, d, rec)
		z.SetUint64(r >> s)
	}
	return z
}

// The remNbyM functions are one Knuth step: they reduce an N-limb u by a normalized M-limb v,
// where the top M limbs of u are below v. The estimate q satisfies q <= qhat <= q+2, so
// u - qhat*v lies in [-2v, v): its top limb is non-zero while it is negative, and at most
// two add-backs fix it.

func rem5by4(u4, u3, u2, u1, u0, v3, v2, v1, v0 uint64) (uint64, uint64, uint64, uint64, uint64) {
	q, _ := bits.Div64(u4, u3, v3) // callers have u4 < 2^63 <= v3
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

// rem4by3r is rem4by3 with the quotient digit from udivrem3by2, given rec = reciprocal3by2(v2, v1).
// The digit is at most one too big, so a single add-back fixes it.
func rem4by3r(u3, u2, u1, u0, v2, v1, v0, rec uint64) (uint64, uint64, uint64) {
	if u3 == v2 && u2 == v1 { // the digit is 2^64 - 1, out of udivrem3by2 range
		return rem4by3(u3, u2, u1, u0, v2, v1, v0)
	}
	q, r2, r1 := udivrem3by2(u3, u2, u1, v2, v1, rec)
	h, l := bits.Mul64(q, v0)
	r0, b := bits.Sub64(u0, l, 0)
	r1, b = bits.Sub64(r1, h, b)
	r2, b = bits.Sub64(r2, 0, b)
	if b != 0 {
		var c uint64
		r0, c = bits.Add64(r0, v0, 0)
		r1, c = bits.Add64(r1, v1, c)
		r2, _ = bits.Add64(r2, v2, c)
	}
	return r2, r1, r0
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
