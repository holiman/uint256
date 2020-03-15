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
