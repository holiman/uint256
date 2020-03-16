// Copyright 2019 Martin Holst Swende. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the COPYING file.
//

// Package math provides integer math utilities.

package uint256

import (
	"fmt"
	"math"
	"math/big"
	"math/bits"
)

var (
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

// Int is represented as an array of 4 uint64, in little-endian order,
// so that Int[3] is the most significant, and Int[0] is the least significant
type Int [4]uint64

func NewInt() *Int {
	return &Int{}
}

// SetBytes interprets buf as the bytes of a big-endian unsigned
// integer, sets z to that value, and returns z.
func (z *Int) SetBytes(buf []byte) *Int {
	var d uint64
	k := 0
	s := uint64(0)
	i := len(buf)
	z[0], z[1], z[2], z[3] = 0, 0, 0, 0
	for ; i > 0; i-- {
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
	return z
}

// Bytes32 returns a the a 32 byte big-endian array.
func (z *Int) Bytes32() [32]byte {
	var b [32]byte
	for i := 0; i < 32; i++ {
		b[31-i] = byte(z[i/8] >> uint64(8*(i%8)))
	}
	return b
}

// Bytes20 returns a the a 32 byte big-endian array.
func (z *Int) Bytes20() [20]byte {
	var b [20]byte
	for i := 0; i < 20; i++ {
		b[19-i] = byte(z[i/8] >> uint64(8*(i%8)))
	}
	return b
}

// Bytes returns the value of z as a big-endian byte slice.
func (z *Int) Bytes() []byte {
	length := z.ByteLen()
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		buf[length-1-i] = byte(z[i/8] >> uint64(8*(i%8)))
	}
	return buf
}

// WriteToSlice writes the content of z into the given byteslice.
// If dest is larger than 32 bytes, z will fill the first parts, and leave
// the end untouched.
// OBS! If dest is smaller than 32 bytes, only the end parts of z will be used
// for filling the array, making it useful for filling an Address object
func (z *Int) WriteToSlice(dest []byte) {
	// ensure 32 bytes
	// A too large buffer. Fill last 32 bytes
	end := len(dest) - 1
	if end > 31 {
		end = 31
	}
	for i := 0; i <= end; i++ {
		dest[end-i] = byte(z[i/8] >> uint64(8*(i%8)))
	}
}

// WriteToArray32 writes all 32 bytes of z to the destination array, including zero-bytes
func (z *Int) WriteToArray32(dest *[32]byte) {
	for i := 0; i < 32; i++ {
		dest[31-i] = byte(z[i/8] >> uint64(8*(i%8)))
	}
}

// WriteToArray20 writes the last 20 bytes of z to the destination array, including zero-bytes
func (z *Int) WriteToArray20(dest *[20]byte) {
	for i := 0; i < 20; i++ {
		dest[19-i] = byte(z[i/8] >> uint64(8*(i%8)))
	}
}

//func (z *Int) WriteToArr32(dest [32]bytes){
//	for i := 0; i < 32; i++ {
//		dest[31-i] = byte(z[i/8] >> uint64(8*(i%8)))
//	}
//}
// Uint64 returns the lower 64-bits of z
func (z *Int) Uint64() uint64 {
	return z[0]
}

// Uint64 returns the lower 64-bits of z and bool whether overflow occurred
func (z *Int) Uint64WithOverflow() (uint64, bool) {
	return z[0], z[1] != 0 || z[2] != 0 || z[3] != 0
}

// Uint64 returns the lower 63-bits of z as int64
func (z *Int) Int64() int64 {
	return int64(z[0] & 0x7fffffffffffffff)
}

// Clone create a new Int identical to z
func (z *Int) Clone() *Int {
	return &Int{z[0], z[1], z[2], z[3]}
}

// Add sets z to the sum x+y
func (z *Int) Add(x, y *Int) {
	z.AddOverflow(x, y) // Inlined.
}

// AddOverflow sets z to the sum x+y, and returns whether overflow occurred
func (z *Int) AddOverflow(x, y *Int) bool {
	var carry uint64
	z[0], carry = bits.Add64(x[0], y[0], 0)
	z[1], carry = bits.Add64(x[1], y[1], carry)
	z[2], carry = bits.Add64(x[2], y[2], carry)
	z[3], carry = bits.Add64(x[3], y[3], carry)
	return carry != 0
}

// Add sets z to the sum ( x+y ) mod m
func (z *Int) AddMod(x, y, m *Int) {

	if z == m { //z is an alias for m
		m = m.Clone()
	}
	if overflow := z.AddOverflow(x, y); overflow {
		// It overflowed. the actual value is
		// 0x10 00..0 + 0x???..??
		//
		// We can split it into
		// 0xffff...f + 0x1 + 0x???..??
		// And mod each item individually
		a := NewInt().SetAllOne()
		a.Mod(a, m)
		z.Mod(z, m)
		z.Add(z, a)
		// reuse a
		a.SetOne()
		z.Add(z, a)

	}
	z.Mod(z, m)
}

// addMiddle128 adds two uint64 integers to the upper part of z
func addTo128(z []uint64, x0, x1 uint64) {
	var carry uint64
	z[0], carry = bits.Add64(z[0], x0, carry) // TODO: The order of adding x, y is confusing.
	z[1], _ = bits.Add64(z[1], x1, carry)
}

// PaddedBytes encodes a Int as a 0-padded byte slice. The length
// of the slice is at least n bytes.
// Example, z =1, n = 20 => [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1]
func (z *Int) PaddedBytes(n int) []byte {
	b := make([]byte, n)

	for i := 0; i < 32 && i < n; i++ {
		b[n-1-i] = byte(z[i/8] >> uint64(8*(i%8)))
	}
	return b
}

// Sub64 set z to the difference x - y, where y is a 64 bit uint
func (z *Int) Sub64(x *Int, y uint64) {
	var carry uint64

	if z[0], carry = bits.Sub64(x[0], y, carry); carry == 0 {
		return
	}
	if z[1], carry = bits.Sub64(x[1], 0, carry); carry == 0 {
		return
	}
	if z[2], carry = bits.Sub64(x[2], 0, carry); carry == 0 {
		return
	}
	z[3]--
}

// Sub sets z to the difference x-y and returns true if the operation underflowed
func (z *Int) SubOverflow(x, y *Int) bool {
	var carry uint64
	z[0], carry = bits.Sub64(x[0], y[0], carry)
	z[1], carry = bits.Sub64(x[1], y[1], carry)
	z[2], carry = bits.Sub64(x[2], y[2], carry)
	z[3], carry = bits.Sub64(x[3], y[3], carry)
	return carry != 0
}

// Sub sets z to the difference x-y
func (z *Int) Sub(x, y *Int) {
	z.SubOverflow(x, y) // Inlined.
}

// umulStep computes (carry, z) = z + (x * y) + carry.
func umulStep(z, x, y, carry uint64) (uint64, uint64) {
	ph, p := bits.Mul64(x, y)
	p, carry = bits.Add64(p, carry, 0)
	carry, _ = bits.Add64(ph, 0, carry)
	p, carry1 := bits.Add64(p, z, 0)
	carry, _ = bits.Add64(carry, 0, carry1)
	return p, carry
}

// umul computes full 256 x 256 -> 512 multiplication.
func umul(x, y *Int) [8]uint64 {
	var res [8]uint64
	for j := 0; j < len(y); j++ {
		var carry uint64
		res[j+0], carry = umulStep(res[j+0], x[0], y[j], carry)
		res[j+1], carry = umulStep(res[j+1], x[1], y[j], carry)
		res[j+2], carry = umulStep(res[j+2], x[2], y[j], carry)
		res[j+3], carry = umulStep(res[j+3], x[3], y[j], carry)
		res[j+4] = carry
	}
	return res
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

	alfa[1], alfa[0] = bits.Mul64(x[0], y[0])
	alfa[3], alfa[2] = bits.Mul64(x[0], y[2])
	alfa[3] += x[0]*y[3] + x[1]*y[2] + x[2]*y[1] + x[3]*y[0] // Top ones, ignore overflow

	beta[2], beta[1] = bits.Mul64(x[0], y[1])
	alfa.Add(alfa, beta)

	beta[2], beta[1] = bits.Mul64(x[1], y[0])
	alfa.Add(alfa, beta)

	beta[3], beta[2] = bits.Mul64(x[1], y[1])
	addTo128(alfa[2:], beta[2], beta[3])

	beta[3], beta[2] = bits.Mul64(x[2], y[0])
	addTo128(alfa[2:], beta[2], beta[3])
	z.Copy(alfa)
}

func (z *Int) Squared() {

	var (
		alfa = &Int{} // Aggregate results
		beta = &Int{} // Calculate intermediate
	)
	// This algo is based on Mul, but since it's squaring, we know that
	// e.g. z.b*y.c + z.c*y.c == 2 * z.b * z.c, and can save some calculations
	// 2 * d * b
	alfa[3], alfa[2] = bits.Mul64(z[0], z[2])
	alfa.lshOne()
	alfa[1], alfa[0] = bits.Mul64(z[0], z[0])

	// 2 * a * d + 2 * b * c
	alfa[3] += (z[0]*z[3] + z[1]*z[2]) << 1

	// 2 * d * c
	beta[2], beta[1] = bits.Mul64(z[0], z[1])
	beta.lshOne()
	alfa.Add(alfa, beta)

	// c * c
	beta[3], beta[2] = bits.Mul64(z[1], z[1])
	addTo128(alfa[2:], beta[2], beta[3])
	z.Copy(alfa)
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

func nlz(d *Int) uint {
	for i := 3; i >= 0; i-- {
		if d[i] != 0 {
			return uint(bits.LeadingZeros64(d[i]) % 32)
		}
	}
	return 0
}

// Normalized form of d.
func shl(d *Int, s uint, isdividend bool) []uint32 {
	dn := make([]uint32, 9)
	for i := 0; i < 4; i++ {
		dn[2*i] = uint32(d[i])
		dn[2*i+1] = uint32(d[i] >> 32)
	}
	var n int
	for i := 7; i >= 0; i-- {
		if dn[i] != 0 {
			n = i
			break
		}
	}
	var prev, t uint32
	for i := 0; i <= n; i++ {
		t = dn[i]
		dn[i] = prev | (dn[i] << s)
		prev = t >> (32 - s)
	}
	if isdividend {
		n = n + 1
		dn[n] = prev
	}
	return dn[:n+1]
}

func divKnuth(x, y []uint32) []uint32 {
	m, n := len(x)-1, len(y)
	q := make([]uint32, m-n+1)
	// Number base (2**32)
	var b uint64 = 0x100000000

	if n <= 2 {
		panic("Should have been handled by udivremBy1()")
	}

	// Main Loop
	var qhat, rhat uint64
	for j := m - n; j >= 0; j-- {
		qhat = (uint64(x[j+n])*b + uint64(x[j+n-1])) / uint64(y[n-1])
		rhat = uint64(x[j+n])*b + uint64(x[j+n-1]) - qhat*uint64(y[n-1])

	AGAIN:
		if qhat >= b || (qhat*uint64(y[n-2]) > b*rhat+uint64(x[j+n-2])) {
			qhat = qhat - 1
			rhat = rhat + uint64(y[n-1])
			if rhat < b {
				goto AGAIN
			}
		}

		// Multiply and subtract.
		var p uint64
		var t, k int64
		for i := 0; i < n; i++ {
			p = qhat * uint64(y[i])
			t = int64(x[i+j]) - k - int64(p&0xffffffff)
			x[i+j] = uint32(t)
			k = int64(p>>32) - (t >> 32)
		}
		t = int64(x[j+n]) - k
		x[j+n] = uint32(t)

		q[j] = uint32(qhat)
		if t < 0 {
			// If we subtracted too much, add back.
			q[j] = q[j] - 1
			var k, t uint64
			for i := 0; i < n; i++ {
				t = uint64(x[i+j]) + uint64(y[i]) + k
				x[i+j] = uint32(t)
				k = t >> 32
			}
			x[j+n] = x[j+n] + uint32(k)
		}
	}
	return q
}

// udivremBy1 divides u by single normalized word d and produces both quotient and remainder.
func udivremBy1(u []uint64, d uint64) (quot []uint64, rem uint64) {
	quot = make([]uint64, len(u)-1)
	rem = u[len(u)-1] // Set the top word as remainder.

	for j := len(u) - 2; j >= 0; j-- {
		quot[j], rem = bits.Div64(rem, u[j], d)
	}

	return quot, rem
}

func udivrem(u []uint64, d *Int) (quot []uint64, rem *Int, err error) {
	var dLen int
	for i := len(d) - 1; i >= 0; i-- {
		if d[i] != 0 {
			dLen = i + 1
			break
		}
	}

	shift := bits.LeadingZeros64(d[dLen-1])

	var dnStorage Int
	dn := dnStorage[:dLen]
	for i := dLen - 1; i > 0; i-- {
		dn[i] = (d[i] << shift) | (d[i-1] >> (64 - shift))
	}
	dn[0] = d[0] << shift

	var uLen int
	for i := len(u) - 1; i >= 0; i-- {
		if u[i] != 0 {
			uLen = i + 1
			break
		}
	}

	var unStorage [9]uint64
	un := unStorage[:uLen+1]
	un[uLen] = u[uLen-1] >> (64 - shift)
	for i := uLen - 1; i > 0; i-- {
		un[i] = (u[i] << shift) | (u[i-1] >> (64 - shift))
	}
	un[0] = u[0] << shift

	// TODO: Skip the highest word of numerator if not significant.

	if dLen == 1 {
		quot, r := udivremBy1(un, dn[0])
		return quot, new(Int).SetUint64(r >> shift), nil
	}

	return quot, rem, fmt.Errorf("not implemented")
}

// Div sets z to the quotient x/y for returns z.
// If d == 0, z is set to 0
func (z *Int) Div(x, y *Int) *Int {
	if y.IsZero() || y.Gt(x) {
		return z.Clear()
	}
	if x.Eq(y) {
		return z.SetOne()
	}
	// Shortcut some cases
	if x.IsUint64() {
		return z.SetUint64(x.Uint64() / y.Uint64())
	}

	// At this point, we know
	// x/y ; x > y > 0

	if quot, _, err := udivrem(x[:], y); err == nil {
		z.Clear()
		copy(z[:len(quot)], quot)
		return z
	}

	// See Knuth, Volume 2, section 4.3.1, Algorithm D.

	// Normalize by shifting divisor left just enough so that its high-order
	// bit is on and u left the same amount.
	// function nlz do the caculating of the amount and shl do the left operation.
	s := nlz(y)
	xn := shl(x, s, true)
	yn := shl(y, s, false)

	// divKnuth do the division of normalized dividend and divisor with Knuth Algorithm D.
	q := divKnuth(xn, yn)

	z.Clear()
	for i := 0; i < len(q); i++ {
		z[i/2] = z[i/2] | uint64(q[i])<<(32*(uint64(i)%2))
	}

	return z
}

// Mod sets z to the modulus x%y for y != 0 and returns z.
// If y == 0, z is set to 0 (OBS: differs from the big.Int)
func (z *Int) Mod(x, y *Int) *Int {
	if x.IsZero() || y.IsZero() {
		return z.Clear()
	}
	switch x.Cmp(y) {
	case -1:
		// x < y
		copy(z[:], x[:])
		return z
	case 0:
		// x == y
		return z.Clear() // They are equal
	}

	// At this point:
	// x != 0
	// y != 0
	// x > y

	// Shortcut trivial case
	if x.IsUint64() {
		return z.SetUint64(x.Uint64() % y.Uint64())
	}

	if _, rem, err := udivrem(x[:], y); err == nil {
		return z.Copy(rem)
	}

	q := NewInt()
	q.Div(x, y)
	q.Mul(q, y)
	z.Sub(x, q)
	return z
}

// Smod interprets x and y as signed integers sets z to
// (sign x) * { abs(x) modulus abs(y) }
// If y == 0, z is set to 0 (OBS: differs from the big.Int)
// OBS! Modifies x and y
func (z *Int) Smod(x, y *Int) *Int {
	ys := y.Sign()
	xs := x.Sign()

	// abs x
	if xs == -1 {
		x.Neg()
	}
	// abs y
	if ys == -1 {
		y.Neg()
	}
	z.Mod(x, y)
	if xs == -1 {
		z.Neg()
	}
	return z
}

// MulMod calculates the modulo-n multiplication of x and y and
// returns z
func (z *Int) MulMod(x, y, m *Int) *Int {
	p := umul(x, y)
	var (
		pl Int
		ph Int
	)
	copy(pl[:], p[:4])
	copy(ph[:], p[4:])

	// If the multiplication is within 256 bits use Mod().
	if ph.IsZero() {
		if z == m { //z is an alias for m; TODO: This should not be needed.
			m = m.Clone()
		}
		z.Mod(&pl, m)
		return z
	}

	if _, rem, err := udivrem(p[:], m); err == nil {
		return z.Copy(rem)
	}

	var pbytes [len(p) * 8]byte
	for i := 0; i < len(pbytes); i++ {
		pbytes[len(pbytes)-1-i] = byte(p[i/8] >> uint64(8*(i%8)))
	}

	// At this point, we _could_ do x=x mod m, y = y mod m, and test again
	// if they fit within 256 bytes. But for now just wrap big.Int instead
	bp := new(big.Int)
	bp.SetBytes(pbytes[:])
	z.SetFromBig(bp.Mod(bp, m.ToBig()))
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
// OBS! This method (potentially) modifies both n and d
func (z *Int) Sdiv(n, d *Int) *Int {
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

// BitLen returns the number of bits required to represent x
func (z *Int) BitLen() int {
	switch {
	case z[3] != 0:
		return 192 + bits.Len64(z[3])
	case z[2] != 0:
		return 128 + bits.Len64(z[2])
	case z[1] != 0:
		return 64 + bits.Len64(z[1])
	default:
		return bits.Len64(z[0])
	}
}
func (z *Int) ByteLen() int {
	return (z.BitLen() + 7) / 8
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
func (z *Int) srsh64(x *Int) *Int {
	z[3], z[2], z[1], z[0] = math.MaxUint64, x[3], x[2], x[1]
	return z
}
func (z *Int) srsh128(x *Int) *Int {
	z[3], z[2], z[1], z[0] = math.MaxUint64, math.MaxUint64, x[3], x[2]
	return z
}
func (z *Int) srsh192(x *Int) *Int {
	z[3], z[2], z[1], z[0] = math.MaxUint64, math.MaxUint64, math.MaxUint64, x[3]
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
	return z[0] > x[0]
}

// Slt interprets z and x as signed integers, and returns
// true if z < x
func (z *Int) Slt(x *Int) bool {

	zSign := z.Sign()
	xSign := x.Sign()

	switch {
	case zSign >= 0 && xSign < 0:
		return false
	case zSign < 0 && xSign >= 0:
		return true
	default:
		return z.Lt(x)
	}
}

// Sgt interprets z and x as signed integers, and returns
// true if z > x
func (z *Int) Sgt(x *Int) bool {
	zSign := z.Sign()
	xSign := x.Sign()

	switch {
	case zSign >= 0 && xSign < 0:
		return true
	case zSign < 0 && xSign >= 0:
		return false
	default:
		return z.Gt(x)
	}
}

// SetIfGt sets z to 1 if z > x
func (z *Int) SetIfGt(x *Int) {
	if z.Gt(x) {
		z.SetOne()
	} else {
		z.Clear()
	}
}

// Lt returns true if z < x
func (z *Int) Lt(x *Int) bool {
	if z[3] < x[3] {
		return true
	}
	if z[3] > x[3] {
		return false
	}
	if z[2] < x[2] {
		return true
	}
	if z[2] > x[2] {
		return false
	}
	if z[1] < x[1] {
		return true
	}
	if z[1] > x[1] {
		return false
	}
	return z[0] < x[0]
}

// SetIfLt sets z to 1 if z < x
func (z *Int) SetIfLt(x *Int) {
	if z.Lt(x) {
		z.SetOne()
	} else {
		z.Clear()
	}
}

// SetUint64 sets z to the value x
func (z *Int) SetUint64(x uint64) *Int {
	z[3], z[2], z[1], z[0] = 0, 0, 0, x
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

// Cmp compares z and x and returns:
//
//   -1 if z <  x
//    0 if z == x
//   +1 if z >  x
//
func (z *Int) Cmp(x *Int) (r int) {
	if z.Gt(x) {
		return 1
	}
	if z.Lt(x) {
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

// IsUint128 reports whether z can be represented in 128 bits.
func (z *Int) IsUint128() bool {
	return (z[3] == 0) && (z[2] == 0)
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

// SetAllOne sets all the bits of z to 1
func (z *Int) SetAllOne() *Int {
	z[3], z[2], z[1], z[0] = math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64
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
	// If the MSB is 0, Srsh is same as Rsh.
	if !z.isBitSet(255) {
		return z.Rsh(x, n)
	}
	// n % 64 == 0
	if n&0x3f == 0 {
		switch n {
		case 0:
			return z.Copy(x)
		case 64:
			return z.srsh64(x)
		case 128:
			return z.srsh128(x)
		case 192:
			return z.srsh192(x)
		default:
			return z.SetAllOne()
		}
	}
	var (
		a uint64 = math.MaxUint64 << (64 - n%64)
	)
	// Big swaps first
	switch {
	case n > 192:
		if n > 256 {
			return z.SetAllOne()
		}
		z.srsh192(x)
		n -= 192
		goto sh192
	case n > 128:
		z.srsh128(x)
		n -= 128
		goto sh128
	case n > 64:
		z.srsh64(x)
		n -= 64
		goto sh64
	default:
		z.Copy(x)
	}

	// remaining shifts
	z[3], a = (z[3]>>n)|a, z[3]<<(64-n)

sh64:
	z[2], a = (z[2]>>n)|a, z[2]<<(64-n)

sh128:
	z[1], a = (z[1]>>n)|a, z[1]<<(64-n)

sh192:
	z[0] = (z[0] >> n) | a

	return z
}

// Copy copies the value x into z, and returns z
func (z *Int) Copy(x *Int) *Int {
	*z = *x
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
	// in z, z[0] is the least significant
	//
	if number, overflow := n.Uint64WithOverflow(); !overflow {
		if number < 32 {
			number := z[4-1-number/8]
			offset := (n[0] & 0x7) << 3 // 8*(n.d % 8)
			z[0] = (number & (0xff00000000000000 >> offset)) >> (56 - offset)
			z[3], z[2], z[1] = 0, 0, 0
			return z
		}
	}
	return z.Clear()
}

// Hex returns a hex representation of z
func (z *Int) Hex() string {
	return fmt.Sprintf("%016x.%016x.%016x.%016x", z[3], z[2], z[1], z[0])
}

// Exp sets z = base**exponent mod 2**256, and returns z.
func (z *Int) Exp(base, exponent *Int) *Int {
	res := Int{1, 0, 0, 0}
	// b^0 == 1
	if exponent.IsZero() || base.IsOne() {
		return z.Copy(&res)
	}
	// b^1 == b
	if exponent.IsOne() {
		return z.Copy(base)
	}
	var (
		word       uint64
		bits       int
		multiplier = *base
	)
	expBitlen := exponent.BitLen()

	word = exponent[0]
	bits = 0
	for ; bits < expBitlen && bits < 64; bits++ {
		if word&1 == 1 {
			res.Mul(&res, &multiplier)
		}
		multiplier.Squared()
		word >>= 1
	}

	word = exponent[1]
	for ; bits < expBitlen && bits < 128; bits++ {
		if word&1 == 1 {
			res.Mul(&res, &multiplier)
		}
		multiplier.Squared()
		word >>= 1
	}

	word = exponent[2]
	for ; bits < expBitlen && bits < 192; bits++ {
		if word&1 == 1 {
			res.Mul(&res, &multiplier)
		}
		multiplier.Squared()
		word >>= 1
	}

	word = exponent[3]
	for ; bits < expBitlen && bits < 256; bits++ {
		if word&1 == 1 {
			res.Mul(&res, &multiplier)
		}
		multiplier.Squared()
		word >>= 1
	}
	return z.Copy(&res)
}

//Extend length of two’s complement signed integer
// sets z to
//  - num if back  > 31
//  - num interpreted as a signed number with sign-bit at (back*8+7), extended to the full 256 bits
func (z *Int) SignExtend(back, num *Int) {
	if back.GtUint64(31) {
		z.Copy(num)
		return
	}
	bit := uint(back.Uint64()*8 + 7)

	mask := back.Lsh(back.SetOne(), bit)
	mask.Sub64(mask, 1)
	if num.isBitSet(bit) {
		num.Or(num, mask.Not())
	} else {
		num.And(num, mask)
	}

}

var _ fmt.Formatter = zero

func (z *Int) Format(s fmt.State, ch rune) {
	z.ToBig().Format(s, ch)
}
