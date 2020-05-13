// uint256: Fixed size 256-bit math library
// Copyright 2020 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

package uint256

import (
	"fmt"
	"math/big"
	"math/bits"
)

const (
	maxWords = 256 / bits.UintSize // number of big.Words in 256-bit

	// The constants below work as compile-time checks: in case evaluated to
	// negative value it cannot be assigned to uint type and compilation fails.
	// These particular expressions check if maxWords either 4 or 8 matching
	// 32-bit and 64-bit architectures.
	_ uint = -(maxWords & (maxWords - 1)) // maxWords is power of two.
	_ uint = -(maxWords & ^(4 | 8))       // maxWords is 4 or 8.
)

// ToBig returns a big.Int version of z.
func (z *Int) ToBig() *big.Int {
	b := new(big.Int)
	switch maxWords { // Compile-time check.
	case 4: // 64-bit architectures.
		words := [4]big.Word{big.Word(z[0]), big.Word(z[1]), big.Word(z[2]), big.Word(z[3])}
		b.SetBits(words[:])
	case 8: // 32-bit architectures.
		words := [8]big.Word{
			big.Word(z[0]), big.Word(z[0] >> 32),
			big.Word(z[1]), big.Word(z[1] >> 32),
			big.Word(z[2]), big.Word(z[2] >> 32),
			big.Word(z[3]), big.Word(z[3] >> 32),
		}
		b.SetBits(words[:])
	}
	return b
}

// FromBig is a convenience-constructor from big.Int.
// Returns a new Int and whether overflow occurred.
func FromBig(b *big.Int) (*Int, bool) {
	z := &Int{}
	overflow := z.SetFromBig(b)
	return z, overflow
}

// SetFromBig converts a big.Int to Int and sets the value to z.
// TODO: Ensure we have sufficient testing, esp for negative bigints.
func (z *Int) SetFromBig(b *big.Int) bool {
	z.Clear()
	words := b.Bits()
	overflow := len(words) > maxWords

	switch maxWords { // Compile-time check.
	case 4: // 64-bit architectures.
		if len(words) > 0 {
			z[0] = uint64(words[0])
			if len(words) > 1 {
				z[1] = uint64(words[1])
				if len(words) > 2 {
					z[2] = uint64(words[2])
					if len(words) > 3 {
						z[3] = uint64(words[3])
					}
				}
			}
		}
	case 8: // 32-bit architectures.
		numWords := len(words)
		if overflow {
			numWords = maxWords
		}
		for i := 0; i < numWords; i++ {
			if i%2 == 0 {
				z[i/2] = uint64(words[i])
			} else {
				z[i/2] |= uint64(words[i]) << 32
			}
		}
	}

	if b.Sign() == -1 {
		z.Neg(z)
	}
	return overflow
}

// Format implements fmt.Formatter. It accepts the formats
// 'b' (binary), 'o' (octal with 0 prefix), 'O' (octal with 0o prefix),
// 'd' (decimal), 'x' (lowercase hexadecimal), and
// 'X' (uppercase hexadecimal).
// Also supported are the full suite of package fmt's format
// flags for integral types, including '+' and ' ' for sign
// control, '#' for leading zero in octal and for hexadecimal,
// a leading "0x" or "0X" for "%#x" and "%#X" respectively,
// specification of minimum digits precision, output field
// width, space or zero padding, and '-' for left or right
// justification.
//
func (z *Int) Format(s fmt.State, ch rune) {
	z.ToBig().Format(s, ch)
}
