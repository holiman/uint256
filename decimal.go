// uint256: Fixed size 256-bit math library
// Copyright 2020 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

package uint256

import (
	"io"
	"strconv"
	"strings"
)

const twoPow256Sub1 = "115792089237316195423570985008687907853269984665640564039457584007913129639935"
const twoPow128 = "340282366920938463463374607431768211456"
const twoPow64 = "18446744073709551616"

func (z *Int) Dec() string {
	return z.ToBig().String()
}

// SetString implements a subset of (*big.Int).SetString
// ok will be true iff i == nil
// The SetString method on uint256.Int differs from big.Int when it comes to
// base 16, since uint256.Int is stricter, requiring:
// - 0x or 0X prefix
// - Non-empty string after prefix
// - No leading zeroes
// The base10 version allows leading zeroes.
func (z *Int) SetString(s string, base int) (i *Int, ok bool) {
	switch base {
	case 0:
		if strings.HasPrefix(s, "0x") {
			err := z.fromHex(s)
			if err != nil {
				return nil, false
			}
			return z, true
		}
		if err := z.SetFromDecimal(s); err != nil {
			return nil, false
		}
		return z, true
	case 10:
		if err := z.SetFromDecimal(s); err != nil {
			return nil, false
		}
		return z, true
	case 16:
		if err := z.fromHex(s); err != nil {
			return nil, false
		}
		return z, true
	}
	return nil, false
}

// FromDecimal is a convenience-constructor to create an Int from a
// decimal (base 10) string. Numbers larger than 256 bits are not accepted.
func FromDecimal(hex string) (*Int, error) {
	var z Int
	if err := z.SetFromDecimal(hex); err != nil {
		return nil, err
	}
	return &z, nil
}

// SetFromDecimal sets z from the given string, interpreted as a decimal number.
func (z *Int) SetFromDecimal(s string) (err error) {
	// Remove max one leading +
	if len(s) > 0 && s[0] == '+' {
		s = s[1:]
	}
	// Remove any number of leading zeroes
	if len(s) > 0 && s[0] == '0' {
		var i int
		var c rune
		for i, c = range s {
			if c != '0' {
				break
			}
		}
		s = s[i:]
	}
	if len(s) < len(twoPow256Sub1) {
		return z.fromDecimal(s)
	}
	if len(s) == len(twoPow256Sub1) {
		if s > twoPow256Sub1 {
			return ErrBig256Range
		}
		return z.fromDecimal(s)
	}
	return ErrBig256Range
}

// multipliers holds the values that are needed for fromDecimal
var multipliers = [5]*Int{
	nil,                                 // represents first round, no multiplication needed
	&Int{10000000000000000000, 0, 0, 0}, // 10 ^ 19
	&Int{687399551400673280, 5421010862427522170, 0, 0},                     // 10 ^ 38
	&Int{5332261958806667264, 17004971331911604867, 2938735877055718769, 0}, // 10 ^ 57
	&Int{0, 8607968719199866880, 532749306367912313, 1593091911132452277},   // 10 ^ 76
}

// fromDecimal is a helper function to only ever be called via SetFromDecimal
// this function takes a string and chunks it up, calling ParseUint on it up to 5 times
// these chunks are then multiplied by the proper power of 10, then added together.
func (z *Int) fromDecimal(bs string) error {
	// first clear the input
	z.Clear()
	// the maximum value of uint64 is 18446744073709551615, which is 20 characters
	// one less means that a string of 19 9's is always within the uint64 limit
	var (
		num       uint64
		err       error
		remaining = len(bs)
	)
	if remaining == 0 {
		return io.EOF
	}
	// We proceed in steps of 19 characters (nibbles), from least significant to most significant.
	// This means that the first (up to) 19 characters do not need to be multiplied.
	// In the second iteration, our slice of 19 characters needs to be multipleied
	// by a factor of 10^19. Et cetera.
	for i, mult := range multipliers {
		if remaining <= 0 {
			return nil // Done
		} else if remaining > 19 {
			num, err = strconv.ParseUint(bs[remaining-19:remaining], 10, 64)
		} else {
			// Final round
			num, err = strconv.ParseUint(bs, 10, 64)
		}
		if err != nil {
			return err
		}
		// add that number to our running total
		if i != 0 {
			base := NewInt(num)
			z.Add(z, base.Mul(base, mult))
		} else {
			z.SetUint64(num)
		}
		// Chop off another 19 characters
		if remaining > 19 {
			bs = bs[0 : remaining-19]
		}
		remaining -= 19
	}
	return nil
}
