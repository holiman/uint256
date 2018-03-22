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

package uint256

import (
	"fmt"
	"math/big"
	"math/bits"
)

var (
	bigtt256  = bigPow(2, 256)
	SignedMax = &Int{
		0xffffffffffffffff,
		0xffffffffffffffff,
		0xffffffffffffffff,
		0x7fffffffffffffff,
	}
	SignedMin = &Int{
		0x0000000000000000,
		0x0000000000000000,
		0x0000000000000000,
		0x8000000000000000,
	}
	zero = &Int{}
)

// bigPow returns a ** b as a big integer.
func bigPow(a, b int64) *big.Int {
	r := big.NewInt(a)
	return r.Exp(r, big.NewInt(b), nil)
}

// Int is represented as an array of 4 uint64, in big-endian order,
// so that int[4] is the most significant, and int[0] is the least significant
type Int [4]uint64

func NewInt() *Int {
	return &Int{}
}

// FromBig is a convenience-constructor from big.Int.
// returns a new Int and whether overflow occurred
func FromBig(int *big.Int) (*Int, bool) {
	// Let's not ruin the argument
	z := &Int{}
	overflow := z.SetFromBig(int)
	return z, overflow
}

// SetFromBig is a convenience-setter from big.Int. Not optimized for speed, mainly for easy testing
func (z *Int) SetFromBig(int *big.Int) bool {
	z.SetBytes(int.Bytes())
	if int.Sign() == -1 {
		z.Neg()
	}
	return len(int.Bits()) > 64
}

// SetBytes interprets buf as the bytes of a big-endian unsigned
// integer, sets z to that value, and returns z.
func (z *Int) SetBytes(buf []byte) *Int {
	var d uint64
	k := 0
	s := uint64(0)
	i := len(buf)
	for ; i > 0; i-- {
		//fmt.Printf("buf[i-1] %x \n",buf[i-1])
		d |= uint64(buf[i-1]) << s
		if s += 8; s == 64 {
			z[k] = d
			k++
			s, d = 0, 0
			if k >= len(z) {
				break
			}
		}
	}
	if k < len(z) {
		z[k] = d
	}
	//fmt.Printf("z %v \n", z.Hex())
	return z
}

// GetBytes returns a the 32 bytes of z (little-endian)
func (z *Int) Bytes() [32]byte {
	var b [32]byte
	for i := 0; i < 32; i++ {
		b[32-i] = byte(z[i/8] >> uint64(8*(i%8)))
	}
	return b
}

// Uint64 returns the lower 64-bits of z
func (z *Int) Uint64() uint64 {
	return z[0]
}

// Uint64 returns the lower 64-bits of z and bool whether overflow occurred
func (z *Int) Uint64WithOverflow() (uint64, bool) {
	return z[0], (z[1] != 0 || z[2] != 0 || z[3] != 0)
}

// Uint64 returns the lower 63-bits of z as int64
func (z *Int) Int64() int64 {
	return int64(z[0] & 0x7fffffffffffffff)
}

// Clone create a new Int identical to z
func (z *Int) Clone() *Int {
	return &Int{z[0], z[1], z[2], z[3]}
}

const bitmask32 = 0x00000000ffffffff

// u64Add adds return a+b+carry and whether overflow occurred
func u64Add(a, b uint64, c bool) (uint64, bool) {
	if c {
		e := a + b + 1
		return e, (e <= a)
	}
	e := a + b
	return e, (e < a)
}

// u64Add adds return a-b-carry and whether underflow occurred
func u64Sub(a, b uint64, c bool) (uint64, bool) {
	if c {
		return a - b - 1, b >= a
	}
	return a - b, b > a
}

// Add sets z to the sum x+y
func (z *Int) Add(x, y *Int) {
	var (
		carry bool
	)
	z[0], carry = u64Add(x[0], y[0], carry)
	z[1], carry = u64Add(x[1], y[1], carry)
	z[2], carry = u64Add(x[2], y[2], carry)
	// Last group
	z[3] = x[3] + y[3]
	if carry {
		z[3]++
	}
}

// AddOverflow sets z to the sum x+y, and returns whether overflow occurred
func (z *Int) AddOverflow(x, y *Int) bool {
	var carry bool
	for i, _ := range z {
		z[i], carry = u64Add(x[i], y[i], carry)
	}
	return carry
}

// addLow128 adds two uint64 integers to the lower half of z ( y is the least significant)
func (z *Int) addLow128(x, y uint64) {
	var carry bool
	z[0], carry = u64Add(z[0], y, carry)
	z[1], carry = u64Add(z[1], x, carry)
	if carry {
		if z[2]++; z[2] == 0 {
			z[3]++
		}
	}
}

// addMiddle128 adds two uint64 integers to the middle part of z
func (z *Int) addMiddle128(x, y uint64) {
	var carry bool
	z[1], carry = u64Add(z[1], y, carry)
	z[2], carry = u64Add(z[2], x, carry)
	if carry {
		z[3]++
	}
}

// addMiddle128 adds two uint64 integers to the upper part of z
func (z *Int) addHigh128(x, y uint64) {
	var carry bool
	z[2], carry = u64Add(z[2], y, carry)
	if carry {
		z[3]++
	}
	z[3] += x
}

// PaddedBytes encodes a Int as a 0-padded byte slice. The length
// of the slice is at least n bytes.
// Example, z =1, n = 20 => [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1]
func (z *Int) PaddedBytes(n int) []byte {
	b := make([]byte, n)

	for i := 0; i < 32 && i < n; i++ {
		b[i] = byte(z[4-i/8] >> uint64(8*(i%8)))
	}
	return b
}

// Sub64 set z to the difference x - y, where y is a 64 bit uint
func (z *Int) Sub64(x *Int, y uint64) {
	var underflow bool

	if z[0], underflow = u64Sub(z[0], y, underflow); !underflow {
		return
	}
	if z[1], underflow = u64Sub(z[1], 0, underflow); !underflow {
		return
	}
	if z[2], underflow = u64Sub(z[2], 0, underflow); !underflow {
		return
	}
	z[3]--
}

// Sub sets z to the difference x-y and returns true if the operation underflowed
func (z *Int) SubOverflow(x, y *Int) bool {
	var (
		underflow bool
	)
	if z[0], underflow = u64Sub(z[0], y, underflow); !underflow {
		return false
	}
	if z[1], underflow = u64Sub(z[1], 0, underflow); !underflow {
		return false
	}
	if z[2], underflow = u64Sub(z[2], 0, underflow); !underflow {
		return false
	}
	z[3], underflow = u64Sub(z[3], 0, underflow)
	return underflow
}

// Sub sets z to the difference x-y
func (z *Int) Sub(x, y *Int) {
	var underflow bool

	z[0], underflow = u64Sub(x[0], y[0], underflow)
	z[1], underflow = u64Sub(x[1], y[1], underflow)
	z[2], underflow = u64Sub(x[2], y[2], underflow)
	if underflow {
		z[3] = x[3] - y[3] - 1
	} else {
		z[3] = x[3] - y[3]
	}
}

// mulIntoLower128 multiplies two 64-bit uints and sets the result as the lower two uint64s (c,d) in x.
// This method does not touch the upper two (a,b)
func (z *Int) mulIntoLower128(x, y uint64) *Int {

	if x == 0 || y == 0 {
		z[0], z[1] = 0, 0
		return z
	}
	low32x, low32y := x&bitmask32, y&bitmask32
	high32x, high32y := x>>32, y>>32

	z[1], z[0] = high32x*high32y, low32x*low32y

	d := low32x * high32y // Needs up 32
	z.addLow128(d>>32, (d&bitmask32)<<32)

	d = high32x * low32y // Needs up 32
	z.addLow128(d>>32, (d&bitmask32)<<32)

	return z
}

// mulIntoMiddle128 multiplies two 64-bit uints and sets the result as the middle two uint64s (b,c) in x.
// This method does not touch the other two (a,d)
func (z *Int) mulIntoMiddle128(x, y uint64) *Int {

	if x == 0 || y == 0 {
		z[1], z[2] = 0, 0
		return z
	}
	low32x, low32y := x&bitmask32, y&bitmask32
	high32x, high32y := x>>32, y>>32

	z[2], z[1] = high32x*high32y, low32x*low32y

	d := low32x * high32y // Needs up 32
	z.addMiddle128(d>>32, (d&bitmask32)<<32)

	d = high32x * low32y // Needs up 32
	z.addMiddle128(d>>32, (d&bitmask32)<<32)

	return z
}

// mulIntoUpper128 multiplies two 64-bit uints and sets the result as the upper two uint64s (a,b) in x.
// This method does not touch the other two (c,d)
func (z *Int) mulIntoUpper128(x, y uint64) *Int {

	if x == 0 || y == 0 {
		z[2], z[3] = 0, 0
		return z
	}
	low32x, low32y := x&bitmask32, y&bitmask32
	high32x, high32y := x>>32, y>>32

	z[3], z[2] = high32x*high32y, low32x*low32y

	d := low32x * high32y // Needs up 32
	z.addHigh128(d>>32, (d&bitmask32)<<32)

	d = high32x * low32y // Needs up 32
	z.addHigh128(d>>32, (d&bitmask32)<<32)

	return z
}

// Mul sets z to the sum x*y
func (z *Int) Mul(x, y *Int) {

	var (
		alfa = &Int{} // Aggregate results
		beta = &Int{} // Calculate intermediate
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

	alfa.mulIntoLower128(x[0], y[0])
	alfa.mulIntoUpper128(x[0], y[2])
	alfa[3] += x[0]*y[3] + x[1]*y[2] + x[2]*y[1] + x[3]*y[0] // Top ones, ignore overflow

	beta.mulIntoMiddle128(x[0], y[1])
	alfa.Add(alfa, beta)

	beta.Clear().mulIntoMiddle128(x[1], y[0])
	alfa.Add(alfa, beta)

	beta.Clear().mulIntoUpper128(x[1], y[1])
	alfa.addHigh128(beta[3], beta[2])

	beta.Clear().mulIntoUpper128(x[2], y[0])
	alfa.addHigh128(beta[3], beta[2])
	z.Copy(alfa)

}
func (x *Int) Squared() {

	var (
		alfa = &Int{} // Aggregate results
		beta = &Int{} // Calculate intermediate
	)
	// This algo is based on Mul, but since it's squaring, we know that
	// e.g. x.b*y.c + x.c*y.c == 2 * x.b * x.c, and can save some calculations
	// 2 * d * b
	alfa.mulIntoUpper128(x[0], x[2]).lshOne()
	alfa.mulIntoLower128(x[0], x[0])

	// 2 * a * d + 2 * b * c
	alfa[3] += (x[0]*x[3] + x[1]*x[2]) << 1

	// 2 * d * c
	beta.mulIntoMiddle128(x[0], x[1]).lshOne()
	alfa.Add(alfa, beta)

	// c * c
	beta.Clear().mulIntoUpper128(x[1], x[1])
	alfa.addHigh128(beta[3], beta[2])
	x.Copy(alfa)
}

func (z *Int) setBit(n uint) *Int {
	// n == 0 -> LSB
	// n == 255 -> MSB
	if n < 256 {
		z[n>>6] |= 1 << (n & 0x3f)
	}
	return z
}

// isBitSet returns true if bit n is set, where n = 0 eq LSB
func (z *Int) isBitSet(n uint) bool {
	if n > 255 {
		return false
	}
	// z [ n / 64] & 1 << (n % 64)
	return (z[n>>6] & (1 << (n & 0x3f))) != 0
}

// Div sets z to the quotient n/d for returns z.
// If d == 0, z is set to 0
func (z *Int) Div(n, d *Int) *Int {
	if d.IsZero() || d.Gt(n) {
		return z.Clear()
	}
	if n.Eq(d) {
		return z.SetOne()
	}
	// Shortcut some cases
	if n.IsUint64() {
		return z.SetUint64(n.Uint64() / d.Uint64())
	}
	// At this point, we know
	// n/d ; n > d > 0

	// The rest is a pretty un-optimized implementation of "Long division"
	// from https://en.wikipedia.org/wiki/Division_algorithm.
	// Could probably be improved upon (it's very slow now)

	r := &Int{}
	q := &Int{}

	for i := n.Bitlen() - 1; i >= 0; i-- {
		// Left-shift r by 1 bit
		r.lshOne()
		// SetFromBig the least-significant bit of r equal to bit i of the numerator
		if ni := n.isBitSet(uint(i)); ni {
			r[0] |= 1
		}
		if !r.Lt(d) {
			r.Sub(r, d)
			q.setBit(uint(i))
		}
	}
	z.Copy(q)
	return z
}

// Abs interprets x as a a signed number, and sets z to the Abs value
//   S256(0)        = 0
//   S256(1)        = 1
//   S256(2**255)   = -2**255
//   S256(2**256-1) = -1

func (z *Int) Abs() *Int {
	if z.Lt(SignedMin) {
		return z
	}
	z.Sub(zero, z)
	return z
}
func (z *Int) Neg() *Int {
	z.Sub(zero, z)
	return z
}

// Sdiv interprets n and d as signed integers, does a
// signed division on the two operands and sets z to the result
// If d == 0, z is set to 0
func (z *Int) Sdiv(n, d *Int) *Int {
	if d.IsZero() || n.IsZero() {
		return z.Clear()
	}
	if n.Eq(d) {
		return z.SetOne()
	}
	// Shortcut some cases
	if n.IsUint64() && d.IsUint64() {
		return z.SetUint64(n.Uint64() / d.Uint64())
	}
	if n.Sign() > 0 {
		if d.Sign() > 0 {
			// pos / pos
			z.Div(n, d)
			return z
		} else {
			// pos / neg
			z.Div(n, d.Neg())
			return z.Neg()
		}
	}

	if d.Sign() < 0 {
		// neg / neg
		z.Div(n.Neg(), d.Neg())
		return z
	}
	// neg / pos
	z.Div(n.Neg(), d)
	return z.Neg()
}

// Mod sets z to the modulus x%y for y != 0 and returns z.
// If y == 0, z is set to 0 (OBS: differs from the big.Int)
func (z *Int) Mod(x, y *Int) *Int {
	panic("TODO implement me")
}

// Smod interprets x and y as signed integers sets z to
// (sign x) * { abs(x) modulus abs(y) }
// If y == 0, z is set to 0 (OBS: differs from the big.Int)
func (z *Int) Smod(x, y *Int) *Int {
	panic("TODO implement me")
}

// Sign returns:
//
//	-1 if z <  0
//	 0 if z == 0
//	+1 if z >  0
// Where z is interpreted as a signed number
func (z *Int) Sign() int {
	if z.IsZero() {
		return 0
	}
	if z.Lt(SignedMin) {
		return 1
	}
	return -1
}

// Bitlen returns the number of bits required to represent x
func (x *Int) Bitlen() int {
	switch {
	case x[3] != 0:
		return 192 + bits.Len64(x[3])
	case x[2] != 0:
		return 128 + bits.Len64(x[2])
	case x[1] != 0:
		return 64 + bits.Len64(x[1])
	default:
		return bits.Len64(x[0])
	}
}

func (z *Int) lsh64(x *Int) *Int {
	z[3], z[2], z[1], z[0] = x[2], x[1], x[0], 0
	return z
}
func (z *Int) lsh128(x *Int) *Int {
	z[3], z[2], z[1], z[0] = x[1], x[0], 0, 0
	return z
}
func (z *Int) lsh192(x *Int) *Int {
	z[3], z[2], z[1], z[0] = x[0], 0, 0, 0
	return z
}
func (z *Int) rsh64(x *Int) *Int {
	z[3], z[2], z[1], z[0] = 0, x[3], x[2], x[1]
	return z
}
func (z *Int) rsh128(x *Int) *Int {
	z[3], z[2], z[1], z[0] = 0, 0, x[3], x[2]
	return z
}
func (z *Int) rsh192(x *Int) *Int {
	z[3], z[2], z[1], z[0] = 0, 0, 0, x[3]
	return z
}

// Not sets z = ^x and returns z.
func (z *Int) Not() *Int {
	z[3], z[2], z[1], z[0] = ^z[3], ^z[2], ^z[1], ^z[0]
	return z
}

// Gt returns true if z > x
func (z *Int) Gt(x *Int) bool {
	if z[3] > x[3] {
		return true
	}
	if z[3] < x[3] {
		return false
	}
	if z[2] > x[2] {
		return true
	}
	if z[2] < x[2] {
		return false
	}
	if z[1] > x[1] {
		return true
	}
	if z[1] < x[1] {
		return false
	}
	return z[0] > z[0]
}

// Slt interprets x and y as signed integers, and returns
// true if x < y
func (x *Int) Slt(y *Int) bool {
	x.Sign()
	panic("TODO")
}

// Sgt interprets x and y as signed integers, and returns
// true if x > y
func (x *Int) Sgt(y *Int) bool {
	panic("TODO")
}

// SetIfGt sets f to 1 if f > g
func (f *Int) SetIfGt(g *Int) {
	if f.Gt(g) {
		f.SetOne()
	} else {
		f.Clear()
	}
}

// Lt returns true if l < g
func (f *Int) Lt(g *Int) bool {
	if f[3] < g[3] {
		return true
	}
	if f[3] > g[3] {
		return false
	}
	if f[2] < g[2] {
		return true
	}
	if f[2] > g[2] {
		return false
	}
	if f[1] < g[1] {
		return true
	}
	if f[1] > g[1] {
		return false
	}
	return f[0] < g[0]
}

// SetIfLt sets f to 1 if f < g
func (f *Int) SetIfLt(g *Int) {
	if f.Lt(g) {
		f.SetOne()
	} else {
		f.Clear()
	}
}

// SetUint64 sets f to the value a
func (z *Int) SetUint64(n uint64) *Int {
	z[3], z[2], z[1], z[0] = 0, 0, 0, n
	return z
}

// Eq returns true if z == x
func (z *Int) Eq(x *Int) bool {
	return (z[0] == x[0]) && (z[1] == x[1]) && (z[2] == x[2]) && (z[3] == x[3])
}

// SetIfEq sets x to
// 1 if z == x
// 0 if Z != x
func (z *Int) SetIfEq(x *Int) {
	if z.Eq(x) {
		z.SetOne()
	} else {
		z.Clear()
	}
}

// Cmp compares x and y and returns:
//
//   -1 if x <  y
//    0 if x == y
//   +1 if x >  y
//
func (x *Int) Cmp(y *Int) (r int) {
	if x.Gt(y) {
		return 1
	}
	if x.Lt(y) {
		return -1
	}
	return 0
}

// LtUint64 returns true if x is smaller than n
func (z *Int) LtUint64(n uint64) bool {
	return (z[3] == 0) && (z[2] == 0) && (z[1] == 0) && z[0] < n
}

// LtUint64 returns true if x is larger than n
func (z *Int) GtUint64(n uint64) bool {
	return (z[3] != 0) || (z[2] != 0) || (z[1] != 0) || z[0] > n
}

// IsUint64 reports whether z can be represented as a uint64.
func (z *Int) IsUint64() bool {
	return (z[3] == 0) && (z[2] == 0) && (z[1] == 0)
}

// IsZero returns true if z == 0
func (z *Int) IsZero() bool {
	return (z[3] == 0) && (z[2] == 0) && (z[1] == 0) && (z[0] == 0)
}

// IsOne returns true if z == 1
func (z *Int) IsOne() bool {
	return (z[3] == 0) && (z[2] == 0) && (z[1] == 0) && (z[0] == 1)
}

// Clear sets z to 0
func (z *Int) Clear() *Int {
	z[3], z[2], z[1], z[0] = 0, 0, 0, 0
	return z
}

// SetOne sets z to 1
func (z *Int) SetOne() *Int {
	z[3], z[2], z[1], z[0] = 0, 0, 0, 1
	return z
}

// Lsh shifts z by 1 bit.
func (z *Int) lshOne() {
	var (
		a, b uint64
	)
	a = z[0] >> 63
	b = z[1] >> 63

	z[0] = z[0] << 1
	z[1] = z[1]<<1 | a

	a = z[2] >> 63
	z[2] = z[2]<<1 | b
	z[3] = z[3]<<1 | a
}

// Lsh sets z = x << n and returns z.
func (z *Int) Lsh(x *Int, n uint) *Int {
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
	case n > 192:
		if n > 256 {
			return z.Clear()
		}
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
	a = z[0] >> (64 - n)
	z[0] = z[0] << n

sh64:
	b = z[1] >> (64 - n)
	z[1] = (z[1] << n) | a

sh128:
	a = z[2] >> (64 - n)
	z[2] = (z[2] << n) | b

sh192:
	z[3] = (z[3] << n) | a

	return z
}

// Rsh sets z = x >> n and returns z.
func (z *Int) Rsh(x *Int, n uint) *Int {
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
	case n > 192:
		if n > 256 {
			return z.Clear()
		}
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
	a = z[3] << (64 - n)
	z[3] = z[3] >> n

sh64:
	b = z[2] << (64 - n)
	z[2] = (z[2] >> n) | a

sh128:
	a = z[1] << (64 - n)
	z[1] = (z[1] >> n) | b

sh192:
	z[0] = (z[0] >> n) | a

	return z
}

// Srsh (Signed/Arithmetic right shift)
// considers z to be a signed integer, during right-shift
// and sets z = x >> n and returns z.
func (z *Int) Srsh(x *Int, n uint) *Int {
	panic("implement me")
}

// Copy copies the value x into z, and returns z
func (z *Int) Copy(x *Int) *Int {
	z[0], z[1], z[2], z[3] = x[0], x[1], x[2], x[3]
	return z
}

// Or sets z = x | y and returns z.
func (z *Int) Or(x, y *Int) *Int {
	z[0] = x[0] | y[0]
	z[1] = x[1] | y[1]
	z[2] = x[2] | y[2]
	z[3] = x[3] | y[3]
	return z
}

// And sets z = x & y and returns z.
func (z *Int) And(x, y *Int) *Int {
	z[0] = x[0] & y[0]
	z[1] = x[1] & y[1]
	z[2] = x[2] & y[2]
	z[3] = x[3] & y[3]
	return z
}

// Xor sets z = x ^ y and returns z.
func (z *Int) Xor(x, y *Int) *Int {
	z[0] = x[0] ^ y[0]
	z[1] = x[1] ^ y[1]
	z[2] = x[2] ^ y[2]
	z[3] = x[3] ^ y[3]
	return z
}

// Byte sets z to the value of the byte at position n,
// with 'z' considered as a big-endian 32-byte integer
// if 'n' > 32, f is set to 0
// Example: f = '5', n=31 => 5
func (z *Int) Byte(n *Int) *Int {
	if number, overflow := n.Uint64WithOverflow(); !overflow {
		if number < 32 {
			number := z[number/8]
			offset := (n[0] & 0x7) << 3 // 8*(n.d % 8)
			z[0] = (number & (0xff00000000000000 >> offset)) >> (56 - offset)
			z[3], z[2], z[1] = 0, 0, 0
			return z
		}
	}
	return z.Clear()
}

// Hex returns a hex representation of f
func (f *Int) Hex() string {
	return fmt.Sprintf("%016x.%016x.%016x.%016x", f[3], f[2], f[1], f[0])
}

// Exp implements exponentiation by squaring, and sets
// z to base^exp
func (z *Int) Exp(base, exponent *Int) *Int {
	return z.Copy(ExpF(base, exponent))
}

// ExpF returns a newly-allocated big integer and does not change
// base or exponent.
func ExpF(base, exponent *Int) *Int {
	z := &Int{1, 0, 0, 0}
	// b^0 == 1
	if exponent.IsZero() || base.IsOne() {
		return z
	}
	// b^1 == b
	if exponent.IsOne() {
		z.Copy(base)
		return z
	}
	var (
		word uint64
		bits int
	)
	expBitlen := exponent.Bitlen()

	word = exponent[0]
	bits = 0
	for ; bits < expBitlen && bits < 64; bits++ {
		if word&1 == 1 {
			z.Mul(z, base)
		}
		base.Squared()
		word >>= 1
	}

	word = exponent[1]
	for ; bits < expBitlen && bits < 128; bits++ {
		if word&1 == 1 {
			z.Mul(z, base)
		}
		base.Squared()
		word >>= 1
	}

	word = exponent[2]
	for ; bits < expBitlen && bits < 192; bits++ {
		if word&1 == 1 {
			z.Mul(z, base)
		}
		base.Squared()
		word >>= 1
	}

	word = exponent[3]
	for ; bits < expBitlen && bits < 256; bits++ {
		if word&1 == 1 {
			z.Mul(z, base)
		}
		base.Squared()
		word >>= 1
	}
	return z
}

//Extend length of two’s complement signed integer
// sets z to
//  - num if back  > 31
//  - num interpreted as a signed number with sign-bit at (back*8+7), extended to the full 256 bits
func (z *Int) SignExtend(back, num *Int) {
	panic("implement me")
}
