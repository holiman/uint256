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
	_S          = bits.UintSize / 8   // word size in bytes
	u256_nWords = 256 / bits.UintSize // number of Words in 256-bit
)

// marshallBigintGeneric is a platform-independent implementation of MarshallBigInt
func NewFromBig(b *big.Int) (*Int, bool) {
	var overflow bool
	z := b.Bits()
	numWords := len(z)
	if numWords == 0 {
		return &Int{}, overflow
	}
	// If there's more than 64 bits, we can skip all higher words
	// z consists of 64 or 32-bit words. So we only care about the last
	// (or last two)
	if numWords > u256_nWords {
		z = z[:u256_nWords]
		numWords = len(z)
		overflow = true
	}
	// Code below is for 64-bit platforms only (numWords: [1-4] )
	fixed := &Int{}
	fixed[0] = uint64(z[0])
	if numWords > 0 {
		fixed[1] = uint64(z[1])
		if numWords > 1 {
			fixed[2] = uint64(z[2])
			if numWords > 2 {
				fixed[3] = uint64(z[3])
			}
		}
	}
	if b.Sign() == -1 {
		fixed.Neg()
	}
	return fixed, overflow

}
