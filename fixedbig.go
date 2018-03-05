// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package math provides integer math utilities.

package math

import (
	"fmt"
	"math/big"
	"math/bits"
)

var (
	tt256 = BigPow(2, 256)
)

// BigPow returns a ** b as a big integer.
func BigPow(a, b int64) *big.Int {
	r := big.NewInt(a)
	return r.Exp(r, big.NewInt(b), nil)
}

type Fixed256bit struct {
	a uint64 // Most significant
	b uint64
	c uint64
	d uint64 // Least significant
}

// newFixedFromBig is a convenience-constructor from big.Int. Not optimized for speed, mainly for easy testing
func NewFixedFromBig(int *big.Int) (*Fixed256bit, bool) {
	// Let's not ruin the argument
	z := &Fixed256bit{}
	overflow := z.Set(int)
	return z, overflow
}

func NewFixed() *Fixed256bit {
	return &Fixed256bit{}
}

// Set is a convenience-setter from big.Int. Not optimized for speed, mainly for easy testing
func (z *Fixed256bit) Set(int *big.Int) bool {
	// Let's not ruin the argument
	x := new(big.Int).Set(int)

	for x.Cmp(new(big.Int)) < 0 {
		// Below 0
		x.Add(tt256, x)
	}
	z.d = x.Uint64()
	z.c = x.Rsh(x, 64).Uint64()
	z.b = x.Rsh(x, 64).Uint64()
	z.a = x.Rsh(x, 64).Uint64()
	x.Rsh(x, 64).Uint64()
	return len(x.Bits()) != 0
}
func (z *Fixed256bit) Clone() *Fixed256bit {
	return &Fixed256bit{z.a, z.b, z.c, z.d}
}

const bitmask32 = 0x00000000ffffffff

func add64(a uint64, b uint64, carry uint64) (uint64, uint64) {

	var (
		q, sum uint64
	)
	sum = carry + (a & bitmask32) + (b & bitmask32)
	q = sum & bitmask32
	carry = sum >> 32
	sum = carry + (a >> 32) + (b >> 32)
	q |= (sum & bitmask32) << 32
	carry = sum >> 32
	return q, carry
}

// Add2 sets z to the sum x+y
func (z *Fixed256bit) Add2(x, y *Fixed256bit) {

	var (
		carry, sum uint64
		q          uint64
	)
	// Least significant
	sum = (y.d & bitmask32) + (x.d & bitmask32)
	carry = sum >> 32
	sum = carry + (y.d >> 32) + (x.d >> 32)
	carry = sum >> 32
	z.d = x.d + y.d

	//	z.c, carry = add64(y.c, x.c, carry)
	// Written out as:

	sum = carry + (y.c & bitmask32) + (x.c & bitmask32)
	q = sum & bitmask32
	carry = sum >> 32
	sum = carry + (y.c >> 32) + (x.c >> 32)
	q |= (sum & bitmask32) << 32
	carry = sum >> 32
	z.c = q

	//	z.b, carry = add64(y.b, x.b, carry)
	sum = carry + (y.b & bitmask32) + (x.b & bitmask32)
	q = sum & bitmask32
	carry = sum >> 32
	sum = carry + (y.b >> 32) + (x.b >> 32)
	q |= (sum & bitmask32) << 32
	carry = sum >> 32
	z.b = q

	z.a = x.a + y.a + carry
}

// Add sets z to the sum x+y
func (z *Fixed256bit) Add(x, y *Fixed256bit) {

	var (
		sum   uint64
		carry uint64
	)
	// LSB
	sum = x.d + y.d
	if sum < x.d {
		carry = 1
	}
	z.d = sum

	// Next 64 bits
	sum = x.c + y.c
	if sum < x.c {
		sum += carry
		carry = 1
	} else {
		sum += carry
		if sum < x.c {
			carry = 1
		} else {
			carry = 0
		}
	}
	z.c = sum
	// Second to last group
	sum = x.b + y.b

	if sum < x.b {
		sum += carry
		carry = 1
	} else {
		sum += carry
		if sum < x.b {
			carry = 1
		} else {
			carry = 0
		}
	}
	z.b = sum

	// Last group
	z.a = x.a + y.a + carry

}

// addLow128 adds two uint64 integers to x, as c and d ( d is the least significant)
func (x *Fixed256bit) addLow128(c, d uint64) {

	var (
		sum   uint64
		carry uint64
	)
	// LSB
	sum = x.d + d
	if sum < x.d {
		carry = 1
	}
	x.d = sum

	// Next 64 bits
	sum = x.c + c

	if sum < x.c {
		sum += carry
		carry = 1
	} else {
		sum += carry
		if sum < x.c {
			carry = 1
		} else {
			// done
			x.c = sum
			return
		}
	}
	x.c = sum
	sum = x.b + carry
	if sum < x.b {
		x.a = x.a + 1
	}
	x.b = sum
}

// addMiddle128 adds two uint64 integers to x, as b and c ( d is the least significant)
func (x *Fixed256bit) addMiddle128(b, c uint64) {

	var (
		sum   uint64
		carry uint64
	)
	sum = x.c + c
	if sum < x.c {
		carry = 1
	}
	x.c = sum

	// Next 64 bits
	sum = x.b + b

	if sum < x.b {
		sum += carry
		carry = 1
	} else {
		sum += carry
		if sum < x.b {
			carry = 1
		} else {
			// done
			x.b = sum
			return
		}
	}
	x.a = x.a + 1
}

// addMiddle128 adds two uint64 integers to x, as a and b ( a is the most significant)
func (x *Fixed256bit) addHigh128(a, b uint64) {

	var (
		sum uint64
	)
	sum = x.b + b
	if sum < x.b {
		x.b = sum
		x.a += a + 1
		return
	}
	x.b = sum
	x.a += a
}

// Sub sets z to the difference x-y
func (z *Fixed256bit) Sub(x, y *Fixed256bit) {

	var (
		underflow uint64
		q         uint64
	)

	q = x.d - y.d
	if q > x.d {
		underflow = 1
	}
	z.d = q

	q = x.c - y.c
	if q > x.c { // underflow again
		q -= underflow
		underflow = 1
	} else {
		// No underflow, we can decrement it
		q -= underflow
		// May cause another underflow
		if q > x.c {
			underflow = 1
		} else {
			underflow = 0
		}
	}
	z.c = q

	q = x.b - y.b
	if q > x.b { // underflow again
		q -= underflow
		underflow = 1
	} else {
		// No underflow, we can decrement it
		q -= underflow
		// May cause another underflow
		if q > x.b {
			underflow = 1
		} else {
			underflow = 0
		}
	}
	z.b = q
	z.a = x.a - y.a - underflow
}

// Sub sets z to the difference x-y and returns true if the operation underflowed
func (z *Fixed256bit) SubOverflow(x, y *Fixed256bit) bool {

	var (
		underflow bool
		q         uint64
	)

	q = x.d - y.d
	underflow = (x.d < y.d)
	z.d = q

	q = x.c - y.c
	if q > x.c { // underflow again
		if underflow {
			q--
		}
		underflow = true
	} else if underflow {
		// No underflow, we can decrement it
		q--
		// May cause another underflow
		underflow = q > x.c
	}
	z.c = q

	q = x.b - y.b
	if q > x.b { // underflow again
		if underflow {
			q--
		}
		underflow = true
	} else if underflow {
		// No underflow, we can decrement it
		q--
		// May cause another underflow
		underflow = q > x.b
	}
	z.b = q

	q = x.a - y.a
	if q > x.a { // underflow again
		if underflow {
			q--
		}
		underflow = true
	} else if underflow {
		// No underflow, we can decrement it
		q--
		// May cause another underflow
		underflow = q > x.a
	}

	z.a = q
	return underflow
}

// mulIntoLower64 multiplies two 64-bit uints and sets the result in x. The parameter y
// is used as a buffer, and will be overwritten (does not have to be cleared prior
// to usage.
func (x *Fixed256bit) mulIntoLower64(a, b uint64) *Fixed256bit {

	if a == 0 || b == 0 {
		return x.Clear()
	}
	low_a := a & bitmask32
	low_b := b & bitmask32
	high_a := a >> 32
	high_b := b >> 32

	x.a, x.b, x.c, x.d = 0, 0, high_a*high_b, low_a*low_b

	d := low_a * high_b // Needs up 32
	x.addLow128(d>>32, (d&bitmask32)<<32)

	d = high_a * low_b // Needs up 32
	x.addLow128(d>>32, (d&bitmask32)<<32)

	return x
}

// mulIntoMiddle64 equals mulIntoLower(..).lsh64()
func (x *Fixed256bit) mulIntoMiddle64(a, b uint64) *Fixed256bit {

	if a == 0 || b == 0 {
		return x.Clear()
	}
	low_a := a & bitmask32
	low_b := b & bitmask32
	high_a := a >> 32
	high_b := b >> 32

	x.a, x.b, x.c, x.d = 0, high_a*high_b, low_a*low_b, 0

	d := low_a * high_b // Needs up 32
	x.addMiddle128(d>>32, (d&bitmask32)<<32)

	d = high_a * low_b // Needs up 32
	x.addMiddle128(d>>32, (d&bitmask32)<<32)

	return x
}

// mulIntoMiddle64 equals mulIntoLower(..).lsh128()
func (x *Fixed256bit) mulIntoUpper64(a, b uint64) *Fixed256bit {

	if a == 0 || b == 0 {
		return x.Clear()
	}
	low_a := a & bitmask32
	low_b := b & bitmask32
	high_a := a >> 32
	high_b := b >> 32

	x.a, x.b, x.c, x.d = high_a*high_b, low_a*low_b, 0, 0

	d := low_a * high_b // Needs up 32
	x.addHigh128(d>>32, (d&bitmask32)<<32)

	d = high_a * low_b // Needs up 32
	x.addHigh128(d>>32, (d&bitmask32)<<32)

	return x
}

// Mul sets z to the sum x*y
func (z *Fixed256bit) Mul(x, y *Fixed256bit) {

	var (
		alfa = &Fixed256bit{} // Aggregate results
		beta = &Fixed256bit{} // Calculate intermediate
	)
	// The numbers are internally represented as [ a, b, c, d ]
	// We do the following operations
	//
	// d1 * d2
	// d1 * c2 (upshift 64)
	// d1 * b2 (upshift 128)
	// d1 * a2 (upshift 192)
	//
	// c1 * d2 (upshift 64)
	// c1 * c2 (upshift 128)
	// c1 * b2 (upshift 192)
	//
	// b1 * d2 (upshift 128)
	// b1 * c2 (upshift 192)
	//
	// a1 * d2 (upshift 192)
	//
	// And we aggregate results into 'alfa'

	// One optimization, however, is reordering.
	// For these ones, we don't care about if they overflow, thus we can use native multiplication
	// and set the result immediately into `a` of the result.
	// b1 * c2 (upshift 192)
	// a1 * d2 (upshift 192)
	// d1 * a2 (upshift 192)
	// c1 * b2 11(upshift 192)

	// Remaining ops:
	//
	// d1 * d2
	// d1 * c2 (upshift 64)
	// d1 * b2 (upshift 128)
	//
	// c1 * d2 (upshift 64)
	// c1 * c2 (upshift 128)
	//
	// b1 * d2 (upshift 128)

	alfa.mulIntoLower64(x.d, y.d)
	alfa.a = x.d*y.a + x.c*y.b + x.b*y.c + x.a*y.d // Top ones, ignore overflow

	beta.mulIntoMiddle64(x.d, y.c) //.lsh64(beta)

	alfa.Add(alfa, beta)

	beta.mulIntoUpper64(x.d, y.b) //.lsh128(beta)

	alfa.Add(alfa, beta)

	beta.mulIntoMiddle64(x.c, y.d) //.lsh64(beta)

	alfa.Add(alfa, beta)
	beta.mulIntoUpper64(x.c, y.c) //.lsh128(beta)
	alfa.Add(alfa, beta)

	beta.mulIntoUpper64(x.b, y.d) //.lsh128(beta)
	z.Add(alfa, beta)

}
func (z *Fixed256bit) setBit(n uint) {
	// n == 0 -> LSB
	// n == 256 -> MSB
	var w *uint64
	if n < 64 {
		w = &z.d
	} else if n < 128 {
		w = &z.c
	} else if n < 192 {
		w = &z.b
	} else if n < 256 {
		w = &z.a
	} else {
		return
	}

	//n %= 64
	n &= 0x3f
	//	mask := 0x1 << n
	*w |= (1 << n)

}
func (z *Fixed256bit) isBitSet(n uint) bool {
	// n == 0 -> LSB
	// n == 256 -> MSB
	var w uint64
	if n < 64 {
		w = z.d
	} else if n < 128 {
		w = z.c
	} else if n < 192 {
		w = z.b
	} else if n < 256 {
		w = z.a
	} else {
		w = 0
	}

	//n %= 64
	n &= 0x3f
	//	mask := 0x1 << n
	return w&(1<<n) != 0
}

// Div sets z to the quotient n/d for returns z.
// If d == 0, z is set to 0
// Div implements Euclidean division (unlike Go); see DivMod for more details.
func (z *Fixed256bit) Div(n, d *Fixed256bit) *Fixed256bit {
	if d.IsZero() || d.Gt(n) {
		return z.Clear()
	}
	if n.Eq(d) {
		return z.SetOne()
	}
	// Shortcut some cases
	if n.IsUint64() {
		return z.SetUint64(n.d / d.d)
	}
	// At this point, we know
	// n/d ; n > d > 0

	// The rest is a pretty un-optimized implementation of "Long division"
	// from https://en.wikipedia.org/wiki/Division_algorithm.
	// Could probably be improved upon (it's very slow now)

	r := &Fixed256bit{}
	q := &Fixed256bit{}

	for i := n.Bitlen() - 1; i >= 0; i-- {
		// Left-shift r by 1 bit
		r.lshOne()
		// Set the least-significant bit of r equal to bit i of the numerator
		if ni := n.isBitSet(uint(i)); ni {
			r.d |= 1
		}
		if !r.Lt(d) {
			r.Sub(r, d)
			q.setBit(uint(i))
		}
	}
	z.Copy(q)
	return z
}

func (x *Fixed256bit) Bitlen() int {
	switch {
	case x.a != 0:
		return 192 + bits.Len64(x.a)
	case x.b != 0:
		return 128 + bits.Len64(x.b)
	case x.c != 0:
		return 64 + bits.Len64(x.c)
	default:
		return bits.Len64(x.d)
	}
}

// Mod sets z to the modulus x%y for y != 0 and returns z.
// If y == 0, z is set to 0 (OBS: differs from the big.Int)
// Mod implements Euclidean modulus (unlike Go); see DivMod for more details.
func (z *Fixed256bit) Mod(x, y *Fixed256bit) *Fixed256bit {
	if y.IsZero() {
		return z.Clear()
	}
	panic("TODO! Implement me")
	return z
}

func (z *Fixed256bit) lsh64(x *Fixed256bit) *Fixed256bit {
	z.a = x.b
	z.b = x.c
	z.c = x.d
	z.d = 0
	return z
}
func (z *Fixed256bit) lsh128(x *Fixed256bit) *Fixed256bit {
	z.a = x.c
	z.b = x.d
	z.c, z.d = 0, 0
	return z
}
func (z *Fixed256bit) lsh192(x *Fixed256bit) *Fixed256bit {
	z.a, z.b, z.c, z.d = x.d, 0, 0, 0
	return z
}
func (z *Fixed256bit) rsh128(x *Fixed256bit) *Fixed256bit {
	z.d = x.b
	z.c = x.a
	z.b = 0
	z.a = 0
	return z
}
func (z *Fixed256bit) rsh64(x *Fixed256bit) *Fixed256bit {
	z.d = x.c
	z.c = x.b
	z.b = x.a
	z.a = 0
	return z
}
func (z *Fixed256bit) rsh192(x *Fixed256bit) *Fixed256bit {
	z.d, z.c, z.b, z.a = x.a, 0, 0, 0
	return z
}

// Not sets z = ^x and returns z.
func (z *Fixed256bit) Not() *Fixed256bit {
	z.a, z.b, z.c, z.d = ^z.a, ^z.b, ^z.c, ^z.d
	return z
}

// Gt returns true if f > g
func (f *Fixed256bit) Gt(g *Fixed256bit) bool {
	if f.a > g.a {
		return true
	}
	if f.a < g.a {
		return false
	}
	if f.b > g.b {
		return true
	}
	if f.b < g.b {
		return false
	}
	if f.c > g.c {
		return true
	}
	if f.c < g.c {
		return false
	}
	if f.d > g.d {
		return true
	}
	return false
}

// SetIfGt sets f to 1 if f > g
func (f *Fixed256bit) SetIfGt(g *Fixed256bit) {
	if f.Gt(g) {
		f.SetOne()
	} else {
		f.Clear()
	}
}

// Lt returns true if l < g
func (f *Fixed256bit) Lt(g *Fixed256bit) bool {
	if f.a < g.a {
		return true
	}
	if f.a > g.a {
		return false
	}
	if f.b < g.b {
		return true
	}
	if f.b > g.b {
		return false
	}
	if f.c < g.c {
		return true
	}
	if f.c > g.c {
		return false
	}
	if f.d < g.d {
		return true
	}
	return false
}

// SetIfLt sets f to 1 if f < g
func (f *Fixed256bit) SetIfLt(g *Fixed256bit) {
	if f.Lt(g) {
		f.SetOne()
	} else {
		f.Clear()
	}
}
func (f *Fixed256bit) SetUint64(a uint64) *Fixed256bit {
	f.a, f.b, f.c, f.d = 0, 0, 0, a
	return f
}

// Eq returns true if f == g
func (f *Fixed256bit) Eq(g *Fixed256bit) bool {
	return (f.a == g.a) && (f.b == g.b) && (f.c == g.c) && (f.d == g.d)
}

// Eq returns true if f == g
func (f *Fixed256bit) SetIfEq(g *Fixed256bit) {
	if (f.a == g.a) && (f.b == g.b) && (f.c == g.c) && (f.d == g.d) {
		f.SetOne()
	} else {
		f.Clear()
	}
}

// Cmp compares x and y and returns:
//
//   -1 if x <  y
//    0 if x == y
//   +1 if x >  y
//
func (x *Fixed256bit) Cmp(y *Fixed256bit) (r int) {
	if x.Gt(y) {
		return 1
	}
	if x.Lt(y) {
		return -1
	}
	return 0
}

// ltsmall can be used to check if x is smaller than n
func (x *Fixed256bit) ltSmall(n uint64) bool {
	return x.a == 0 && x.b == 0 && x.c == 0 && x.d < n
}

// IsUint64 reports whether x can be represented as a uint64.
func (x *Fixed256bit) IsUint64() bool {
	return (x.a == 0) && (x.b == 0) && (x.c == 0)
}

// IsZero returns true if f == 0
func (f *Fixed256bit) IsZero() bool {
	return (f.a == 0) && (f.b == 0) && (f.c == 0) && (f.d == 0)
}

// IsOne returns true if f == 1
func (f *Fixed256bit) IsOne() bool {
	return f.a == 0 && f.b == 0 && f.c == 0 && f.d == 1
}

// Clear sets z to 0
func (z *Fixed256bit) Clear() *Fixed256bit {
	z.a, z.b, z.c, z.d = 0, 0, 0, 0
	return z
}

// SetOne sets z to 1
func (z *Fixed256bit) SetOne() *Fixed256bit {
	z.a, z.b, z.c, z.d = 0, 0, 0, 1
	return z
}

// Lsh shifts z by 1 bit.
func (z *Fixed256bit) lshOne() {
	var (
		a, b uint64
	)
	a = z.d >> 63
	z.d = z.d << 1

	b = z.c >> 63
	z.c = (z.c << 1) | a

	a = z.b >> 63
	z.b = (z.b << 1) | b

	b = z.a >> 63
	z.a = (z.a << 1) | a
}

// Lsh sets z = x << n and returns z.
func (z *Fixed256bit) Lsh(x *Fixed256bit, n uint) *Fixed256bit {
	// n % 64 == 0
	if n&0x3f == 0 {
		switch n {
		case 0:
			return z.Copy(x)
		case 64:
			return z.lsh64(x)
		case 128:
			return z.lsh128(x)
		case 192:
			return z.lsh192(x)
		default:
			return z.Clear()
		}
	}
	var (
		a, b uint64
	)
	// Big swaps first
	switch {
	case n > 256:
		return z.Clear()
	case n > 192:
		z.lsh192(x)
		n -= 192
		goto sh192
	case n > 128:
		z.lsh128(x)
		n -= 128
		goto sh128
	case n > 64:
		z.lsh64(x)
		n -= 64
		goto sh64
	default:
		z.Copy(x)
	}

	// remaining shifts
	a = z.d >> (64 - n)
	z.d = z.d << n

sh64:
	b = z.c >> (64 - n)
	z.c = (z.c << n) | a

sh128:
	a = z.b >> (64 - n)
	z.b = (z.b << n) | b

sh192:
	z.a = (z.a << n) | a

	return z
}

// Rsh sets z = x >> n and returns z.
func (z *Fixed256bit) Rsh(x *Fixed256bit, n uint) *Fixed256bit {
	// n % 64 == 0
	if n&0x3f == 0 {
		switch n {
		case 0:
			return z.Copy(x)
		case 64:
			return z.rsh64(x)
		case 128:
			return z.rsh128(x)
		case 192:
			return z.rsh192(x)
		default:
			return z.Clear()
		}
	}
	var (
		a, b uint64
	)
	// Big swaps first
	switch {
	case n > 256:
		return z.Clear()
	case n > 192:
		z.rsh192(x)
		n -= 192
		goto sh192
	case n > 128:
		z.rsh128(x)
		n -= 128
		goto sh128
	case n > 64:
		z.rsh64(x)
		n -= 64
		goto sh64
	default:
		z.Copy(x)
	}

	// remaining shifts
	a = z.a << (64 - n)
	z.a = z.a >> n

sh64:
	b = z.b << (64 - n)
	z.b = (z.b >> n) | a

sh128:
	a = z.c << (64 - n)
	z.c = (z.c >> n) | b

sh192:
	z.d = (z.d >> n) | a

	return z
}
func (z *Fixed256bit) Copy(x *Fixed256bit) *Fixed256bit {
	z.a, z.b, z.c, z.d = x.a, x.b, x.c, x.d
	return z
}

// Or sets z = x | y and returns z.
func (z *Fixed256bit) Or(x, y *Fixed256bit) *Fixed256bit {
	z.a = x.a | y.a
	z.b = x.b | y.b
	z.c = x.c | y.c
	z.d = x.d | y.d
	return z
}

// And sets z = x & y and returns z.
func (z *Fixed256bit) And(x, y *Fixed256bit) *Fixed256bit {
	z.a = x.a & y.a
	z.b = x.b & y.b
	z.c = x.c & y.c
	z.d = x.d & y.d
	return z
}

// Xor sets z = x ^ y and returns z.
func (z *Fixed256bit) Xor(x, y *Fixed256bit) *Fixed256bit {
	z.a = x.a ^ y.a
	z.b = x.b ^ y.b
	z.c = x.c ^ y.c
	z.d = x.d ^ y.d
	return z
}

// Byte sets f to the value of the byte at position n,
// Example: f = '5', n=31 => 5
func (f *Fixed256bit) Byte(n *Fixed256bit) *Fixed256bit {
	var number uint64
	if n.ltSmall(32) {
		if n.d > 24 {
			// f.d holds bytes [24 .. 31]
			number = f.d
		} else if n.d > 15 {
			// f.c holds bytes [16 .. 23]
			number = f.c
		} else if n.d > 7 {
			// f.b holds bytes [8 .. 15]
			number = f.b
		} else {
			// f.a holds MSB, bytes [0 .. 7]
			number = f.a
		}
		offset := (n.d & 0x7) << 3 // 8*(n.d % 8)
		number = (number & (0xff00000000000000 >> offset)) >> (56 - offset)
	}

	f.a, f.b, f.c, f.d = 0, 0, 0, number
	return f
}

func (f *Fixed256bit) Hex() string {
	return fmt.Sprintf("%016x.%016x.%016x.%016x", f.a, f.b, f.c, f.d)
}

// Exp implements exponentiation by squaring.
// Exp returns a newly-allocated big integer and does not change
// base or exponent.
//
// Courtesy @karalabe and @chfast, with improvements by @holiman
func ExpF(base, exponent *Fixed256bit) *Fixed256bit {
	z := &Fixed256bit{a: 0, b: 0, c: 0, d: 1}
	// b^0 == 1
	if exponent.IsZero() || base.IsOne() {
		return z
	}
	// b^1 == 1
	if exponent.IsOne() {
		z.Copy(base)
		return z
	}
	var (
		word uint64
		bits int
	)
	exp_bitlen := exponent.Bitlen()

	word = exponent.d
	bits = 0
	for ; bits < exp_bitlen && bits < 64; bits++ {
		if word&1 == 1 {
			z.Mul(z, base)
		}
		base.Mul(base, base)
		word >>= 1
	}

	word = exponent.c
	for ; bits < exp_bitlen && bits < 128; bits++ {
		if word&1 == 1 {
			z.Mul(z, base)
		}
		base.Mul(base, base)
		word >>= 1
	}

	word = exponent.b
	for ; bits < exp_bitlen && bits < 192; bits++ {
		if word&1 == 1 {
			z.Mul(z, base)
		}
		base.Mul(base, base)
		word >>= 1
	}

	word = exponent.a
	for ; bits < exp_bitlen && bits < 256; bits++ {
		if word&1 == 1 {
			z.Mul(z, base)
		}
		base.Mul(base, base)
		word >>= 1
	}
	return z
}
