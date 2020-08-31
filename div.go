// uint256: Fixed size 256-bit math library
// Copyright 2020 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

package uint256

import "math/bits"

// reciprocal2by1 computes <^d, ^0> / d.
func reciprocal2by1(d uint64) uint64 {
	d9 := d >> 55
	v0 := uint32(reciprocalTable[d9-256])

	d40 := (d >> 24) + 1
	v1 := uint64((v0 << 11) - uint32(uint64(v0*v0)*d40>>40) - 1)

	v2 := (v1 << 13) + (v1 * (0x1000000000000000 - v1*d40) >> 47)

	d0 := d & 1
	d63 := (d >> 1) + d0 // ceil(d/2)
	e := ((v2 >> 1) & (0 - d0)) - v2*d63
	p1h, _ := bits.Mul64(v2, e)
	v3 := (p1h >> 1) + (v2 << 31)

	p2h, p2l := bits.Mul64(v3, d)
	_, carry := bits.Add64(p2l, d, 0)
	p2h += carry

	v4 := v3 - p2h - d
	return v4
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

func reciprocal3by2(dh, dl uint64) uint64 {
	// d[1] is d.hi, d[0] is d.lo
	var d [2]uint64
	d[0] = dl
	d[1] = dh
	v := reciprocal2by1(d[1])
	p := d[1] * v
	p += d[0]
	if p < d[0] {
		v--
		if p >= d[1] {
			v--
			p -= d[1]
		}
		p -= d[1]
	}

	th, tl := bits.Mul64(v, d[0])

	p += th
	if p < th {
		v--
		if p >= d[1] {
			if p > d[1] || tl >= d[0] {
				v--
			}
		}
	}

	return v
}

func udivrem3by2(u2, u1, u0, dh, dl, reciprocal uint64) (quot uint64, rem [2]uint64) {
	var carry uint64
	qh, ql := bits.Mul64(reciprocal, u2)

	// q = fast_add(q, u)
	ql, carry = bits.Add64(ql, u1, 0)
	qh, _ = bits.Add64(qh, u2, carry)

	r1 := u1 - qh*dh

	th, tl := bits.Mul64(dl, qh)

	// for udivrem2by1
	// auto r = u.lo - q.hi * d;
	//r := ul - qh * d

	// for udivrem3by2
	// auto r = uint128{r1, u0} - t - d;
	//r =: 
	// first do (uint128{r1, u0} - t)

	var r1u0minust [2]uint64
	r1u0minust[0], carry = bits.Sub64(u0, tl, 0)
	r1u0minust[1], _ = bits.Sub64(r1, th, carry)
	// then do (r1u0minust - d)
	var r [2]uint64
	r[0], carry = bits.Sub64(r1u0minust[0], dl, 0)
	r[1], _ = bits.Sub64(r1u0minust[1], dh, carry)

	qh++ // ++q.hi

	if r[1] >= ql {
		qh--
		// r is a uint128
		//r += d
		r[0], carry = bits.Add64(r[0], dl, 0)
		r[1], _ = bits.Add64(r[1], dh, carry)
	}

	//if r >= d {
	if r[1] > dh || (r[1] == dh && r[0] >= dl) {
		qh++
		// r is a uint128
		//r -= d
		r[0], carry = bits.Sub64(r[0], dl, 0)
		r[1], _ = bits.Sub64(r[1], dh, carry)
	}

	return qh, r
}
