package uint256

import (
	"strconv"
	"strings"
)

const twoPow256Sub1 = "115792089237316195423570985008687907853269984665640564039457584007913129639935"
const twoPow128 = "340282366920938463463374607431768211456"
const twoPow64 = "18446744073709551616"

func (z *Int) Base10() string {
	return z.ToBig().String()
}

// SetString implements a subset of (*big.Int).SetString
// ok will be true iff i == nil
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
		err := z.SetFromBase10(s)
		if err != nil {
			return nil, false
		}
		return z, true
	case 10:
		err := z.SetFromBase10(s)
		if err != nil {
			return nil, false
		}
		return z, true
	case 16:
		err := z.fromHex(s)
		if err != nil {
			return nil, false
		}
		return z, true
	}
	return nil, false
}

// FromBase10 is a convenience-constructor to create an Int from a
// decimal (base 10) string. Numbers larger than 256 bits are not accepted.
func FromBase10(hex string) (*Int, error) {
	var z Int
	if err := z.SetFromBase10(hex); err != nil {
		return nil, err
	}
	return &z, nil
}

// SetFromBase10 sets z from the given string, interpreted as a decimal number.
func (z *Int) SetFromBase10(s string) (err error) {
	if len(s) < len(twoPow256Sub1) {
		return z.fromBase10Long(s)
	}
	if len(s) == len(twoPow256Sub1) {
		if s > twoPow256Sub1 {
			return ErrBig256Range
		}
		return z.fromBase10Long(s)
	}
	return ErrBig256Range
}

// scaleTable10 contains the 10-exponents, 
// 0: 10 ^ 0
// 1: 10 ^ 1
// .. 
// 77: 10 ^ 77
var scaleTable10 [78]Int

func init() {
	for k := range scaleTable10 {
		scaleTable10[k].Exp(NewInt(10), NewInt(uint64(k)))
	}
}

// helper function to only ever be called via SetFromBase10
// this function takes a string and chunks it up, calling ParseUint on it up to 5 times
// these chunks are then multiplied by the proper power of 10, then added together.
func (z *Int) fromBase10Long(bs string) error {
	// first clear the input
	z.Clear()
	// if the input value is empty string, just do nothing. effectively, empty string sets to 0
	if bs == "" {
		return nil
	}
	// the maximum value of uint64 is 18446744073709551615, which is 20 characters
	// one less means that a string of 19 9's is always within the uint64 limit
	cutLength := 19
	// cutStart tracks the current position of the string that we are in
	cutStart := 0
	// start iterating from 4 to 1. This is because the maximum value of uint256 is 78 characters,
	// which can be divided into 5 integers of up to 19 characters.
	// however, the last number will always be below 19 characters, so i=0 is dealt with as special case
	for i := 4; i >= 1; i-- {
		// check if the length of the string is larger than cutLength * i
		if len(bs) >= (cutLength * i) {
			// cut the string from the cutStart to the cutLength.
			nm, err := strconv.ParseUint(bs[cutStart:(cutStart+cutLength)], 10, 64)
			if err != nil {
				return ErrSyntaxBase10
			}
			// create a new int with that number as the value
			base := NewInt(nm)
			// pointer to the exponent. We need to multiply our number by 10^(len-cutStart-cutLength)
			// len-cutStart-cutLength is index of the last character in our cutset, counting from the right.
			exp := &scaleTable10[len(bs)-cutStart-cutLength]
			// add that number to our running total
			z.Add(z, base.Mul(exp, base))
			// increase the cut start point, since we have now read from cutStart to cutStart + length
			cutStart = cutStart + cutLength
		}
	}
	// if we have read every character of the string, we are done, and can return
	// this is a short circuit that we can do if the length of the string is a multiple of 19
	if len(bs) == cutStart {
		return nil
	}
	// finally, there are a remaining set of characters.
	// SetFromBase10 already did the check that this remaining cutset, after 4 cuts, will be lower than 19 charactes
	nm, err := strconv.ParseUint(bs[cutStart:], 10, 64)
	if err != nil {
		return ErrSyntaxBase10
	}
	// and add it! no need to multiply by 10^0
	z.AddUint64(z, uint64(nm))
	return nil
}
