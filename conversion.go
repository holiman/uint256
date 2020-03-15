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

// NewFromBig creates new Int from big.Int.
func NewFromBig(b *big.Int) (*Int, bool) {
	z := &Int{}
	overflow := z.SetFromBig(b)
	return z, overflow
}

// SetFromBig
// TODO: finish implementation by adding 32-bit platform support,
// ensure we have sufficient testing, esp for negative bigints
func (z *Int) SetFromBig(b *big.Int) bool {
	var overflow bool
	z.Clear()
	words := b.Bits()
	numWords := len(words)
	if numWords == 0 {
		return overflow
	}
	// If there's more than 64 bits, we can skip all higher words
	// words consists of 64 or 32-bit words. So we only care about the last
	// (or last two)
	if numWords > maxWords {
		words = words[:maxWords]
		numWords = len(words)
		overflow = true
	}
	// Code below is for 64-bit platforms only (numWords: [1-4] )
	z[0] = uint64(words[0])
	if numWords > 1 {
		z[1] = uint64(words[1])
		if numWords > 2 {
			z[2] = uint64(words[2])
			if numWords > 3 {
				z[3] = uint64(words[3])
			}
		}
	}
	if b.Sign() == -1 {
		z.Neg()
	}
	return overflow
}
