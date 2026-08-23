// uint256: Fixed size 256-bit math library
// Copyright 2026 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

package uint256

import "math/bits"

// correctReciprocal verifies a reciprocal made from the saturated initial
// estimate and decrements it when it overestimates 1/(4y).
func correctReciprocal(r4l, r4k, r4j, r4i, r4h uint64, y *Int) (uint64, uint64, uint64, uint64, uint64) {
	x0, x1, x2, x3, x4 := r4l, r4k, r4j, r4i, r4h
	var q0, q1, q2, q3, q4, q5, q6, q7, q8, t0, t1, c, b, t uint64

	q1, q0 = bits.Mul64(x0, y[0])
	q3, q2 = bits.Mul64(x2, y[0])
	q5, q4 = bits.Mul64(x4, y[0])

	t1, t0 = bits.Mul64(x1, y[0])
	q1, c = bits.Add64(q1, t0, 0)
	q2, c = bits.Add64(q2, t1, c)
	q3, c = bits.Add64(q3, 0, c)
	t1, t0 = bits.Mul64(x3, y[0])
	q3, c = bits.Add64(q3, t0, c)
	q4, c = bits.Add64(q4, t1, c)
	q5, _ = bits.Add64(q5, 0, c)

	t1, t0 = bits.Mul64(x0, y[1])
	q1, c = bits.Add64(q1, t0, 0)
	q2, c = bits.Add64(q2, t1, c)
	t1, t0 = bits.Mul64(x2, y[1])
	q3, c = bits.Add64(q3, t0, c)
	q4, c = bits.Add64(q4, t1, c)
	q6, t0 = bits.Mul64(x4, y[1])
	q5, c = bits.Add64(q5, t0, c)
	q6, _ = bits.Add64(q6, 0, c)

	t1, t0 = bits.Mul64(x1, y[1])
	q2, c = bits.Add64(q2, t0, 0)
	q3, c = bits.Add64(q3, t1, c)
	t1, t0 = bits.Mul64(x3, y[1])
	q4, c = bits.Add64(q4, t0, c)
	q5, c = bits.Add64(q5, t1, c)
	q6, _ = bits.Add64(q6, 0, c)

	t1, t0 = bits.Mul64(x0, y[2])
	q2, c = bits.Add64(q2, t0, 0)
	q3, c = bits.Add64(q3, t1, c)
	t1, t0 = bits.Mul64(x2, y[2])
	q4, c = bits.Add64(q4, t0, c)
	q5, c = bits.Add64(q5, t1, c)
	q7, t0 = bits.Mul64(x4, y[2])
	q6, c = bits.Add64(q6, t0, c)
	q7, _ = bits.Add64(q7, 0, c)

	t1, t0 = bits.Mul64(x1, y[2])
	q3, c = bits.Add64(q3, t0, 0)
	q4, c = bits.Add64(q4, t1, c)
	t1, t0 = bits.Mul64(x3, y[2])
	q5, c = bits.Add64(q5, t0, c)
	q6, c = bits.Add64(q6, t1, c)
	q7, _ = bits.Add64(q7, 0, c)

	t1, t0 = bits.Mul64(x0, y[3])
	q3, c = bits.Add64(q3, t0, 0)
	q4, c = bits.Add64(q4, t1, c)
	t1, t0 = bits.Mul64(x2, y[3])
	q5, c = bits.Add64(q5, t0, c)
	q6, c = bits.Add64(q6, t1, c)
	q8, t0 = bits.Mul64(x4, y[3])
	q7, c = bits.Add64(q7, t0, c)
	q8, _ = bits.Add64(q8, 0, c)

	t1, t0 = bits.Mul64(x1, y[3])
	q4, c = bits.Add64(q4, t0, 0)
	q5, c = bits.Add64(q5, t1, c)
	t1, t0 = bits.Mul64(x3, y[3])
	q6, c = bits.Add64(q6, t0, c)
	q7, c = bits.Add64(q7, t1, c)
	q8, _ = bits.Add64(q8, 0, c)

	_, b = bits.Sub64(0, q0, 0)
	_, b = bits.Sub64(0, q1, b)
	_, b = bits.Sub64(0, q2, b)
	_, b = bits.Sub64(0, q3, b)
	_, b = bits.Sub64(0, q4, b)
	_, b = bits.Sub64(0, q5, b)
	_, b = bits.Sub64(0, q6, b)
	_, b = bits.Sub64(0, q7, b)
	_, b = bits.Sub64(uint64(1)<<62, q8, b)

	x0, t = bits.Sub64(r4l, 1, 0)
	x1, t = bits.Sub64(r4k, 0, t)
	x2, t = bits.Sub64(r4j, 0, t)
	x3, t = bits.Sub64(r4i, 0, t)
	x4, _ = bits.Sub64(r4h, 0, t)

	if b != 0 {
		return x0, x1, x2, x3, x4
	}
	return r4l, r4k, r4j, r4i, r4h
}
