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
		err := z.FromBase10(s)
		if err != nil {
			return nil, false
		}
		return z, true
	case 10:
		err := z.FromBase10(s)
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
func (z *Int) FromBase10(s string) (err error) {
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

var scaleTable10 [78]Int

func init() {
	for k := range scaleTable10 {
		scaleTable10[k] = Int{}
		scaleTable10[k].Exp(NewInt(10), NewInt(uint64(k)))
	}
}

func (z *Int) fromBase10Long(bs string) error {
	z[0] = 0
	z[1] = 0
	z[2] = 0
	z[3] = 0
	if bs == "" {
		return nil
	}
	iv := 19
	c := 0
	if len(bs) >= (iv * 4) {
		nm, err := strconv.ParseUint(bs[c:(c+iv)], 10, 64)
		if err != nil {
			return ErrSyntaxBase10
		}
		z.Add(z, new(Int).Mul(&scaleTable10[len(bs)-iv], NewInt(uint64(nm))))
		c = c + iv
	}
	if len(bs) >= (iv * 3) {
		nm, err := strconv.ParseUint(bs[c:(c+iv)], 10, 64)
		if err != nil {
			return ErrSyntaxBase10
		}
		z.Add(z, new(Int).Mul(&scaleTable10[len(bs)-c-iv], NewInt(uint64(nm))))
		c = c + iv
	}
	if len(bs) >= (iv * 2) {
		nm, err := strconv.ParseUint(bs[c:(c+iv)], 10, 64)
		if err != nil {
			return ErrSyntaxBase10
		}
		z.Add(z, new(Int).Mul(&scaleTable10[len(bs)-c-iv], NewInt(uint64(nm))))
		c = c + iv
	}
	if len(bs) >= (iv * 1) {
		nm, err := strconv.ParseUint(bs[c:(c+iv)], 10, 64)
		if err != nil {
			return ErrSyntaxBase10
		}
		z.Add(z, new(Int).Mul(&scaleTable10[len(bs)-c-iv], NewInt(uint64(nm))))
		c = c + iv
	}
	if len(bs) == c {
		return nil
	}
	nm, err := strconv.ParseUint(bs[c:], 10, 64)
	if err != nil {
		return ErrSyntaxBase10
	}
	z.AddUint64(z, uint64(nm))
	return nil
}
