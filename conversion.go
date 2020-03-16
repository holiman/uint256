// Copyright 2020 Martin Holst Swende. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the COPYING file.
//

// Package math provides integer math utilities.

package uint256

import (
	"math/big"
	"math/bits"
)

const (
	maxWords = 256 / bits.UintSize // number of big.Words in 256-bit
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
		return new(big.Int).SetBits(words[:])
	default:
		panic("unsupported architecture")
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
	default:
		panic("unsupported architecture")
	}

	if b.Sign() == -1 {
		z.Neg()
	}
	return overflow
}
