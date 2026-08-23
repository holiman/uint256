// uint256: Fixed size 256-bit math library
// Copyright 2026 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

package uint256

import "math/bits"

// reduce4WithQuotient computes the quotient and least non-negative residue of x modulo m.
//
// It is the quotient-preserving form of reduce4. Keeping it separate leaves the
// existing MulMod reduction path unchanged.
//
// requires a four-word modulus (m[3] != 0) and its inverse (mu)
func (z *Int) reduce4WithQuotient(x *[8]uint64, m *Int, mu *[5]uint64, quotient *[5]uint64) *Int {

	// NB: Most variable names in the comments match the pseudocode for
	// 	Barrett reduction in the Handbook of Applied Cryptography.

	// q1 = x/2^192

	x0 := x[3]
	x1 := x[4]
	x2 := x[5]
	x3 := x[6]
	x4 := x[7]

	// q2 = q1 * mu; q3 = q2 / 2^320

	var q0, q1, q2, q3, q4, q5, t0, t1, c uint64

	q0, _ = bits.Mul64(x3, mu[0])
	q1, t0 = bits.Mul64(x4, mu[0])
	q0, c = bits.Add64(q0, t0, 0)
	q1, _ = bits.Add64(q1, 0, c)

	t1, _ = bits.Mul64(x2, mu[1])
	q0, c = bits.Add64(q0, t1, 0)
	q2, t0 = bits.Mul64(x4, mu[1])
	q1, c = bits.Add64(q1, t0, c)
	q2, _ = bits.Add64(q2, 0, c)

	t1, t0 = bits.Mul64(x3, mu[1])
	q0, c = bits.Add64(q0, t0, 0)
	q1, c = bits.Add64(q1, t1, c)
	q2, _ = bits.Add64(q2, 0, c)

	t1, t0 = bits.Mul64(x2, mu[2])
	q0, c = bits.Add64(q0, t0, 0)
	q1, c = bits.Add64(q1, t1, c)
	q3, t0 = bits.Mul64(x4, mu[2])
	q2, c = bits.Add64(q2, t0, c)
	q3, _ = bits.Add64(q3, 0, c)

	t1, _ = bits.Mul64(x1, mu[2])
	q0, c = bits.Add64(q0, t1, 0)
	t1, t0 = bits.Mul64(x3, mu[2])
	q1, c = bits.Add64(q1, t0, c)
	q2, c = bits.Add64(q2, t1, c)
	q3, _ = bits.Add64(q3, 0, c)

	t1, _ = bits.Mul64(x0, mu[3])
	q0, c = bits.Add64(q0, t1, 0)
	t1, t0 = bits.Mul64(x2, mu[3])
	q1, c = bits.Add64(q1, t0, c)
	q2, c = bits.Add64(q2, t1, c)
	q4, t0 = bits.Mul64(x4, mu[3])
	q3, c = bits.Add64(q3, t0, c)
	q4, _ = bits.Add64(q4, 0, c)

	t1, t0 = bits.Mul64(x1, mu[3])
	q0, c = bits.Add64(q0, t0, 0)
	q1, c = bits.Add64(q1, t1, c)
	t1, t0 = bits.Mul64(x3, mu[3])
	q2, c = bits.Add64(q2, t0, c)
	q3, c = bits.Add64(q3, t1, c)
	q4, _ = bits.Add64(q4, 0, c)

	t1, t0 = bits.Mul64(x0, mu[4])
	_, c = bits.Add64(q0, t0, 0)
	q1, c = bits.Add64(q1, t1, c)
	t1, t0 = bits.Mul64(x2, mu[4])
	q2, c = bits.Add64(q2, t0, c)
	q3, c = bits.Add64(q3, t1, c)
	q5, t0 = bits.Mul64(x4, mu[4])
	q4, c = bits.Add64(q4, t0, c)
	q5, _ = bits.Add64(q5, 0, c)

	t1, t0 = bits.Mul64(x1, mu[4])
	q1, c = bits.Add64(q1, t0, 0)
	q2, c = bits.Add64(q2, t1, c)
	t1, t0 = bits.Mul64(x3, mu[4])
	q3, c = bits.Add64(q3, t0, c)
	q4, c = bits.Add64(q4, t1, c)
	q5, _ = bits.Add64(q5, 0, c)

	// Drop the fractional part of q3

	q0 = q1
	q1 = q2
	q2 = q3
	q3 = q4
	q4 = q5

	quotient[0], quotient[1], quotient[2], quotient[3], quotient[4] = q0, q1, q2, q3, q4

	// r1 = x mod 2^320

	x0 = x[0]
	x1 = x[1]
	x2 = x[2]
	x3 = x[3]
	x4 = x[4]

	// r2 = q3 * m mod 2^320

	var r0, r1, r2, r3, r4 uint64

	r4, r3 = bits.Mul64(q0, m[3])
	_, t0 = bits.Mul64(q1, m[3])
	r4, _ = bits.Add64(r4, t0, 0)

	t1, r2 = bits.Mul64(q0, m[2])
	r3, c = bits.Add64(r3, t1, 0)
	_, t0 = bits.Mul64(q2, m[2])
	r4, _ = bits.Add64(r4, t0, c)

	t1, t0 = bits.Mul64(q1, m[2])
	r3, c = bits.Add64(r3, t0, 0)
	r4, _ = bits.Add64(r4, t1, c)

	t1, r1 = bits.Mul64(q0, m[1])
	r2, c = bits.Add64(r2, t1, 0)
	t1, t0 = bits.Mul64(q2, m[1])
	r3, c = bits.Add64(r3, t0, c)
	r4, _ = bits.Add64(r4, t1, c)

	t1, t0 = bits.Mul64(q1, m[1])
	r2, c = bits.Add64(r2, t0, 0)
	r3, c = bits.Add64(r3, t1, c)
	_, t0 = bits.Mul64(q3, m[1])
	r4, _ = bits.Add64(r4, t0, c)

	t1, r0 = bits.Mul64(q0, m[0])
	r1, c = bits.Add64(r1, t1, 0)
	t1, t0 = bits.Mul64(q2, m[0])
	r2, c = bits.Add64(r2, t0, c)
	r3, c = bits.Add64(r3, t1, c)
	_, t0 = bits.Mul64(q4, m[0])
	r4, _ = bits.Add64(r4, t0, c)

	t1, t0 = bits.Mul64(q1, m[0])
	r1, c = bits.Add64(r1, t0, 0)
	r2, c = bits.Add64(r2, t1, c)
	t1, t0 = bits.Mul64(q3, m[0])
	r3, c = bits.Add64(r3, t0, c)
	r4, _ = bits.Add64(r4, t1, c)

	// r = r1 - r2

	var b uint64

	r0, b = bits.Sub64(x0, r0, 0)
	r1, b = bits.Sub64(x1, r1, b)
	r2, b = bits.Sub64(x2, r2, b)
	r3, b = bits.Sub64(x3, r3, b)
	r4, b = bits.Sub64(x4, r4, b)

	// if r<0 then r+=m

	if b != 0 {
		r0, c = bits.Add64(r0, m[0], 0)
		r1, c = bits.Add64(r1, m[1], c)
		r2, c = bits.Add64(r2, m[2], c)
		r3, c = bits.Add64(r3, m[3], c)
		r4, _ = bits.Add64(r4, 0, c)

		// Compensate for the extra m above; each correction below increments
		// the quotient once.
		quotient[0], c = bits.Sub64(quotient[0], 1, 0)
		quotient[1], c = bits.Sub64(quotient[1], 0, c)
		quotient[2], c = bits.Sub64(quotient[2], 0, c)
		quotient[3], c = bits.Sub64(quotient[3], 0, c)
		quotient[4], _ = bits.Sub64(quotient[4], 0, c)
	}

	// while (r>=m) r-=m

	for {
		// q = r - m
		q0, b = bits.Sub64(r0, m[0], 0)
		q1, b = bits.Sub64(r1, m[1], b)
		q2, b = bits.Sub64(r2, m[2], b)
		q3, b = bits.Sub64(r3, m[3], b)
		q4, b = bits.Sub64(r4, 0, b)

		// if borrow break
		if b != 0 {
			break
		}

		// r = q
		r4, r3, r2, r1, r0 = q4, q3, q2, q1, q0
		quotient[0], c = bits.Add64(quotient[0], 1, 0)
		quotient[1], c = bits.Add64(quotient[1], 0, c)
		quotient[2], c = bits.Add64(quotient[2], 0, c)
		quotient[3], c = bits.Add64(quotient[3], 0, c)
		quotient[4], _ = bits.Add64(quotient[4], 0, c)
	}

	z[3], z[2], z[1], z[0] = r3, r2, r1, r0

	return z
}
