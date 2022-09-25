package uint256

import (
	"strconv"
	"strings"
)

const twoPow256 = "115792089237316195423570985008687907853269984665640564039457584007913129639935"
const twoPow128 = "340282366920938463463374607431768211456"
const twoPow64 = "18446744073709551616"

func (z *Int) SetString(b string, base int) (err error) {
	switch base {
	case 0:
		if strings.HasPrefix(b, "0x") {
			return z.fromHex(b)
		}
		return z.FromBase10(b)
	case 10:
		return z.FromBase10(b)
	case 16:
		return z.fromHex(b)
	}
	return ErrSyntax
}
func (z *Int) FromBase10(s string) (err error) {
	if len(s) < len(twoPow256) {
		return z.fromBase10Long(s)
	}
	if len(s) == len(twoPow256) {
		if s[0] > '1' {
			return ErrBig256Range
		}
		return z.fromBase10Long(s)
	}
	return ErrBig256Range
}

var scaleTable10 [78]*Int

func init() {
	for k := range scaleTable10 {
		scaleTable10[k] = new(Int)
		scaleTable10[k] = scaleTable10[k].Exp(NewInt(10), NewInt(uint64(k)))
	}
}

func (z *Int) fromBase10Long(bs string) error {
	z[0] = 0
	z[1] = 0
	z[2] = 0
	z[3] = 0
	if len(bs) == 0 {
		return nil
	}
	iv := 19
	c := 0
	if len(bs) >= (iv * 4) {
		nm, err := strconv.Atoi(bs[c:(c + iv)])
		if err != nil {
			return err
		}
		z.Add(z, new(Int).Mul(scaleTable10[len(bs)-iv], NewInt(uint64(nm))))
		c = c + iv
	}
	if len(bs) >= (iv * 3) {
		nm, err := strconv.Atoi(bs[c:(c + iv)])
		if err != nil {
			return err
		}
		z.Add(z, new(Int).Mul(scaleTable10[len(bs)-c-iv], NewInt(uint64(nm))))
		c = c + iv
	}
	if len(bs) >= (iv * 2) {
		nm, err := strconv.Atoi(bs[c:(c + iv)])
		if err != nil {
			return err
		}
		z.Add(z, new(Int).Mul(scaleTable10[len(bs)-c-iv], NewInt(uint64(nm))))
		c = c + iv
	}
	if len(bs) >= (iv * 1) {
		nm, err := strconv.Atoi(bs[c:(c + iv)])
		if err != nil {
			return err
		}
		z.Add(z, new(Int).Mul(scaleTable10[len(bs)-c-iv], NewInt(uint64(nm))))
		c = c + iv
	}
	if len(bs) == c {
		return nil
	}
	nm, err := strconv.Atoi(bs[c:])
	if err != nil {
		return err
	}
	z.AddUint64(z, uint64(nm))
	return nil
}
