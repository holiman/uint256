package math

import (
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"
)

func checkOverflow(b *big.Int, f *Fixed256bit, overflow bool) error {
	max := big.NewInt(0).SetBytes(common.Hex2Bytes("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"))
	shouldOverflow := (b.Cmp(max) > 0)
	if overflow != shouldOverflow {
		return fmt.Errorf("Overflow should be %v, was %v\nf= %v\nb= %x\b", shouldOverflow, overflow, f.Hex(), b)
	}
	return nil
}

func randNums() (*big.Int, *Fixed256bit, error) {
	//How many bits? 0-256
	nbits, _ := rand.Int(rand.Reader, big.NewInt(256))
	//Max random value, a 130-bits integer, i.e 2^130
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(nbits.Int64()), nil)
	//Generate cryptographically strong pseudo-random between 0 - max
	b, _ := rand.Int(rand.Reader, max)
	f, overflow := NewFixedFromBig(b)
	err := checkOverflow(b, f, overflow)
	return b, f, err
}
func checkEq(b *big.Int, f *Fixed256bit) bool {
	f2, _ := NewFixedFromBig(b)
	return f.Eq(f2)
}

func TestRandomSubOverflow(t *testing.T) {
	for i := 0; i < 10000; i++ {
		b, f1, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		b2, f2, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		f1a, f2a := f1.Clone(), f2.Clone()
		overflow := f1.SubOverflow(f1, f2)
		b.Sub(b, b2)
		checkOverflow(b, f1, overflow)
		if eq := checkEq(b, f1); !eq {
			t.Fatalf("Expected equality:\nf1= %v\nf2= %v\n[ - ]==\nf= %v\nb= %x\n", f1a.Hex(), f2a.Hex(), f1.Hex(), b)
		}
	}
}
func TestRandomSub(t *testing.T) {
	for i := 0; i < 10000; i++ {
		b, f1, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		b2, f2, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		f1a, f2a := f1.Clone(), f2.Clone()
		f1.Sub(f1, f2)
		b.Sub(b, b2)
		if eq := checkEq(b, f1); !eq {
			t.Fatalf("Expected equality:\nf1= %v\nf2= %v\n[ - ]==\nf= %v\nb= %x\n", f1a.Hex(), f2a.Hex(), f1.Hex(), b)
		}
	}
}
func TestRandomAdd(t *testing.T) {
	for i := 0; i < 10000; i++ {
		b, f1, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		b2, f2, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		f1a, f2a := f1.Clone(), f2.Clone()
		f1.Add(f1, f2)
		b.Add(b, b2)
		//checkOverflow(b, f1, overflow)
		if eq := checkEq(b, f1); !eq {
			t.Fatalf("Expected equality:\nf1= %v\nf2= %v\n[ + ]==\nf= %v\nb= %x\n", f1a.Hex(), f2a.Hex(), f1.Hex(), b)
		}
	}
}
func TestRandomAdd2(t *testing.T) {
	for i := 0; i < 10000; i++ {
		b, f1, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		b2, f2, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		f1a, f2a := f1.Clone(), f2.Clone()
		f1.Add2(f1, f2)
		b.Add(b, b2)
		//checkOverflow(b, f1, overflow)
		if eq := checkEq(b, f1); !eq {
			t.Fatalf("Expected equality:\nf1= %v\nf2= %v\n[ + ]==\nf= %v\nb= %x\n", f1a.Hex(), f2a.Hex(), f1.Hex(), b)
		}
	}
}
func TestRandomMul(t *testing.T) {
	for i := 0; i < 100000; i++ {
		b, f1, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		b2, f2, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		f1a, f2a := f1.Clone(), f2.Clone()
		f1.Mul(f1, f2)
		b.Mul(b, b2)
		if eq := checkEq(b, f1); !eq {
			bf, _ := NewFixedFromBig(b)
			t.Fatalf("Expected equality:\nf1= %v\nf2= %v\n[ * ]==\nf = %v\nbf= %v\nb = %x\n", f1a.Hex(), f2a.Hex(), f1.Hex(), bf.Hex(), b)
		}
	}
}
func TestRandomDiv(t *testing.T) {
	for i := 0; i < 10000; i++ {
		b, f1, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		b2, f2, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		f1a, f2a := f1.Clone(), f2.Clone()
		f1.Div(f1, f2)
		if b2.BitLen() == 0 {
			// zero
			b = big.NewInt(0)
		} else {
			b.Div(b, b2)
		}
		if eq := checkEq(b, f1); !eq {
			bf, _ := NewFixedFromBig(b)
			t.Fatalf("Expected equality:\nf1= %v\nf2= %v\n[ / ]==\nf = %v\nbf= %v\nb = %x\n", f1a.Hex(), f2a.Hex(), f1.Hex(), bf.Hex(), b)
		}
	}
}

func TestRandomLsh(t *testing.T) {
	for i := 0; i < 10000; i++ {
		b, f1, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		f1a := f1.Clone()
		nbits, _ := rand.Int(rand.Reader, big.NewInt(256))
		n := uint(nbits.Uint64())
		f1.Lsh(f1, n)
		b.Lsh(b, n)
		if eq := checkEq(b, f1); !eq {
			bf, _ := NewFixedFromBig(b)
			t.Fatalf("Expected equality:\nf1= %v\n n= %v\n[ << ]==\nf = %v\nbf= %v\nb = %x\n", f1a.Hex(), n, f1.Hex(), bf.Hex(), b)
		}
	}
}

func TestRandomRsh(t *testing.T) {
	for i := 0; i < 10000; i++ {
		b, f1, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		f1a := f1.Clone()
		nbits, _ := rand.Int(rand.Reader, big.NewInt(256))
		n := uint(nbits.Uint64())
		f1.Rsh(f1, n)
		b.Rsh(b, n)
		if eq := checkEq(b, f1); !eq {
			t.Fatalf("Expected equality:\nf1= %v\n n= %v\n[ << ]==\nf= %v\nb= %x\n", f1a.Hex(), n, f1.Hex(), b)
		}
	}
}

const (
	// number of bits in a big.Word
	wordBits = 32 << (uint64(^big.Word(0)) >> 63)
	// number of bytes in a big.Word
	wordBytes = wordBits / 8
)

var (
	tt256m1 = new(big.Int).Sub(tt256, big.NewInt(1))
)

// U256 encodes as a 256 bit two's complement number. This operation is destructive.
func U256(x *big.Int) *big.Int {
	return x.And(x, tt256m1)
}

// Exp implements exponentiation by squaring.
// Exp returns a newly-allocated big integer and does not change
// base or exponent. The result is truncated to 256 bits.
//
// Courtesy @karalabe and @chfast
func Exp(base, exponent *big.Int) *big.Int {
	result := big.NewInt(1)

	for _, word := range exponent.Bits() {
		for i := 0; i < wordBits; i++ {
			if word&1 == 1 {
				U256(result.Mul(result, base))
			}
			U256(base.Mul(base, base))
			word >>= 1
		}
	}
	return result
}

func TestRandomExp(t *testing.T) {
	for i := 0; i < 10000; i++ {
		b_base, base, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		b_exp, exp, err := randNums()
		if err != nil {
			t.Fatal(err)
		}
		basecopy, expcopy := base.Clone(), exp.Clone()

		f_res := ExpF(base, exp)
		b_res := Exp(b_base, b_exp)
		if eq := checkEq(b_res, f_res); !eq {
			bf, _ := NewFixedFromBig(b_res)
			t.Fatalf("Expected equality:\nbase= %v\nexp = %v\n[ ^ ]==\nf = %v\nbf= %v\nb = %x\n", basecopy.Hex(), expcopy.Hex(), f_res.Hex(), bf.Hex(), b_res)
		}
	}
}

func TestFixed256bit_Add(t *testing.T) {
	//	b := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	//	b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))

	/*
			f= 0000000000000000.00000001d1f4043c.ffc96b146ec8f42e.780bfa023e3ffb94
			b= 1d1f4043cffc96b146ec8f42f780bfa023e3ffb94
			f1= 0000000000000000.000000000f91a6e3.514d614d40e14ca3.d285eea405bd4b42
			f2= 0000000000000000.00000001c2625d59.ae7c09c72de7a78b.a5860b5e3882b052

		b := big.NewInt(0).SetBytes(common.Hex2Bytes( "0f91a6e3514d614d40e14ca3d285eea405bd4b42"))
		b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("01c2625d59ae7c09c72de7a78ba5860b5e3882b052"))

	*/

	/*
		f1= 000282209f633a3c.a040e862bb69d925.73449d21bce09ea3.a74348fbf1ced62e
		f2= 0000000000000000.0000000000000000.000000000000003a.fd56300e26f61922
		++
		f= 000282209f633a3c.a040e862bb69d925.73449d21bce09edd.a499790a18c4ef50
		b= 282209f633a3ca040e862bb69d92573449d21bce09edea499790a18c4ef50
	*/
	b1 := big.NewInt(0).SetBytes(common.Hex2Bytes("000282209f633a3ca040e862bb69d92573449d21bce09ea3a74348fbf1ced62e"))
	b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("00000000000000000000000000000000000000000000003afd56300e26f61922"))

	f, _ := NewFixedFromBig(b1)
	f2, _ := NewFixedFromBig(b2)
	fmt.Printf("B1: %x\n", b1)
	fmt.Printf("B2: %x\n", b2)
	fmt.Printf("F1: %s\n", f.Hex())
	fmt.Printf("F2: %s\n", f2.Hex())
	fmt.Println("--")
	b1.Add(b1, b2)
	f.Add(f, f2)
	fmt.Printf("B: %x\n", b1)
	fmt.Printf("F: %s\n", f.Hex())
	/*
		b	000282209f633a3c.a040e862bb69d925.73449d21bce09ede.a499790a18c4ef50
		f	000282209f633a3c.a040e862bb69d925.73449d21bce09edd.a499790a18c4ef50
	*/

}

func TestFixed256bit_Sub(t *testing.T) {

	b1 := big.NewInt(0).SetBytes(common.Hex2Bytes("00000000000000000000000000000000000000000002a3f8ba829e365f479526"))
	b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("000000000000000000000000000000004ffab28fa389b141ce4876fa1965c937"))

	f, _ := NewFixedFromBig(b1)
	f2, _ := NewFixedFromBig(b2)

	fmt.Printf("B1   : %x\n", b1)
	fmt.Printf("B2   : %x\n", b2)
	fmt.Printf("F1   : %s\n", f.Hex())
	fmt.Printf("F2   : %s\n", f2.Hex())
	fmt.Println("--")
	b1.Sub(b1, b2)
	f.Sub(f, f2)
	fmt.Printf("B   : %x\n", b1)
	fmt.Printf("F   : %s\n", f.Hex())

	res, _ := NewFixedFromBig(b1)
	fmt.Printf("b->f: %s\n", res.Hex())
	fmt.Printf("EQ  : %v\n", f.Eq(res))

}

func TestFixed256bit_Mul(t *testing.T) {

	b1 := big.NewInt(0).SetBytes(common.Hex2Bytes("00000000000000000000000000000000000000000002a3f8ba829e365f479526"))
	b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001"))

	/*
		f1= 0000000000000000.0000000000000000.000000000003e6da.5c61238a298b16ec
		f2= 0000000000000000.0000000000000000.0000000000000000.000000000000206c
		[ * ]==
		f = 1d6c3e387e80a3f8.0000000000000bb3.1d6c3e387e80afab.1d6c437ae98b2b90
		bf= 0000000000000000.0000000000000000.000000007e80afab.1d6c437ae98b2b90
		b = 7e80afab1d6c437ae98b2b90

	*/

	f, _ := NewFixedFromBig(b1)
	f2, _ := NewFixedFromBig(b2)

	fmt.Printf("B1   : %x\n", b1)
	fmt.Printf("B2   : %x\n", b2)
	fmt.Printf("F1   : %s\n", f.Hex())
	fmt.Printf("F2   : %s\n", f2.Hex())
	fmt.Println("--")
	b1.Mul(b1, b2)
	f.Mul(f, f2)
	fmt.Printf("B   : %x\n", b1)
	fmt.Printf("F   : %s\n", f.Hex())

	res, _ := NewFixedFromBig(b1)
	fmt.Printf("b->f: %s\n", res.Hex())
	fmt.Printf("EQ  : %v\n", f.Eq(res))
}

func TestFixed256bit_Div(t *testing.T) {

	/*
		b1 := big.NewInt(0).SetBytes(common.Hex2Bytes("000000000000000000000000000000000000000000000000000000000000000C"))
		b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000004"))

		f1= 0000000000000000.0000000000000000.0000000000000000.0000015fa035e510
		f2= 0000000000000000.0000000000000000.0000000000000000.00000000001d0209
		[ / ]==
		f = 0000000000000000.0000000000000000.0000000000000000.0000015fa035e510
		bf= 0000000000000000.0000000000000000.0000000000000000.00000000000c1f28

		b1 := big.NewInt(0).SetBytes(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000015fa035e510"))
		b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("00000000000000000000000000000000000000000000000000000000001d0209"))


		f1= 12cbafcee8f60f9f.3fa308c90fde8d29.8772ffea667aa6bc.109d5c661e7929a5
		f2= 00000c76f4afb041.407a8ea478d65024.f5c3dfe1db1a1bb1.0c5ea8bec314ccf9
		[ / ]==
		f = 0000000000000000.0000000000000000.0000000000000000.0000000000000000
		bf= 0000000000000000.0000000000000000.0000000000000000.0000000000018206

	*/
	b1 := big.NewInt(0).SetBytes(common.Hex2Bytes("12cbafcee8f60f9f3fa308c90fde8d298772ffea667aa6bc109d5c661e7929a5"))
	b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("00000c76f4afb041407a8ea478d65024f5c3dfe1db1a1bb10c5ea8bec314ccf9"))

	f, _ := NewFixedFromBig(b1)
	f2, _ := NewFixedFromBig(b2)

	fmt.Printf("B1   : %x\n", b1)
	fmt.Printf("B2   : %x\n", b2)
	fmt.Printf("F1   : %s\n", f.Hex())
	fmt.Printf("F2   : %s\n", f2.Hex())
	fmt.Println("--")
	b1.Div(b1, b2)
	f.Div(f, f2)
	fmt.Printf("B   : %x\n", b1)
	fmt.Printf("F   : %s\n", f.Hex())

	res, _ := NewFixedFromBig(b1)
	fmt.Printf("b->f: %s\n", res.Hex())
	fmt.Printf("EQ  : %v\n", f.Eq(res))
}

func TestFixedExp(t *testing.T) {

	b_base := big.NewInt(0).SetBytes(common.Hex2Bytes("00000000000000000000000000000000000000000000006d5adef08547abf7eb"))
	b_exp := big.NewInt(0).SetBytes(common.Hex2Bytes("000000000000000000013590cab83b779e708b533b0eef3561483ddeefc841f5"))

	base, _ := NewFixedFromBig(b_base)
	exp, _ := NewFixedFromBig(b_exp)

	fmt.Printf("B1   : %x\n", b_base)
	fmt.Printf("B2   : %x\n", b_exp)
	fmt.Printf("F1   : %s\n", base.Hex())
	fmt.Printf("F2   : %s\n", exp.Hex())
	fmt.Println("--")

	//	base.d = 2
	//	exp.d = 255
	res := ExpF(base, exp)

	b_res := Exp(b_base, b_exp)

	want, _ := NewFixedFromBig(b_res)
	fmt.Printf("B: %x\n", b_res)
	fmt.Printf("want : %s\n", want.Hex())
	fmt.Printf("got  : %s\n", res.Hex())

	fb_res, _ := NewFixedFromBig(b_res)
	fmt.Printf("EQ  : %v\n", res.Eq(fb_res))

}

func benchmark_Add_Bit(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	f, _ := NewFixedFromBig(b1)
	f2, _ := NewFixedFromBig(b2)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f.Add(f, f2)
	}
}
func benchmark_Add_Big(bench *testing.B) {
	b := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b.Add(b, b2)
	}
}
func Benchmark_Add(bench *testing.B){
	bench.Run("big", benchmark_Add_Big)
	bench.Run("fixedbit", benchmark_Add_Bit)
}


func benchmark_SubOverflow_Bit(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	f, _ := NewFixedFromBig(b1)
	f2, _ := NewFixedFromBig(b2)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f.SubOverflow(f, f2)
	}
}
func benchmark_Sub_Bit(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	f, _ := NewFixedFromBig(b1)
	f2, _ := NewFixedFromBig(b2)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f.Sub(f, f2)
	}
}

func benchmark_Sub_Big(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1.Sub(b1, b2)
	}
}
func Benchmark_Sub(bench *testing.B){
	bench.Run("big", benchmark_Sub_Big)
	bench.Run("fixedbit", benchmark_Sub_Bit)
	bench.Run("fixedbit_of", benchmark_SubOverflow_Bit)
}

func benchmark_Mul_Big(bench *testing.B) {
	a := big.NewInt(0).SetBytes(common.Hex2Bytes("f123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b := big.NewInt(0).SetBytes(common.Hex2Bytes("f123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1 := big.NewInt(0)
		b1.Mul(a, b)
		U256(b1)
	}
}

func benchmark_Mul_Bit(bench *testing.B) {
	a := big.NewInt(0).SetBytes(common.Hex2Bytes("f123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b := big.NewInt(0).SetBytes(common.Hex2Bytes("f123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	fa, _ := NewFixedFromBig(a)
	fb, _ := NewFixedFromBig(b)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f := NewFixed()
		f.Mul(fa, fb)
	}
}
func Benchmark_Mul(bench *testing.B) {
	bench.Run("big", benchmark_Mul_Big)
	bench.Run("fixedbit", benchmark_Mul_Bit)
}

func benchmark_And_Big(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1.And(b1, b2)
	}
}
func benchmark_And_Bit(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	f, _ := NewFixedFromBig(b1)
	f2, _ := NewFixedFromBig(b2)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f.And(f, f2)
	}
}
func Benchmark_And(bench *testing.B) {
	bench.Run("big", benchmark_And_Big)
	bench.Run("fixedbit", benchmark_And_Bit)
}

func benchmark_Or_Big(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1.Or(b1, b2)
	}
}
func benchmark_Or_Bit(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	f, _ := NewFixedFromBig(b1)
	f2, _ := NewFixedFromBig(b2)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f.Or(f, f2)
	}
}
func Benchmark_Or(bench *testing.B) {
	bench.Run("big", benchmark_Or_Big)
	bench.Run("fixedbit", benchmark_Or_Bit)
}

func benchmark_Xor_Big(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1.Xor(b1, b2)
	}
}
func benchmark_Xor_Bit(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	f, _ := NewFixedFromBig(b1)
	f2, _ := NewFixedFromBig(b2)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f.Xor(f, f2)
	}
}

func Benchmark_Xor(bench *testing.B) {
	bench.Run("big", benchmark_Xor_Big)
	bench.Run("fixedbit", benchmark_Xor_Bit)
}

func benchmark_Cmp_Big(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1.Cmp(b2)
	}
}
func benchmark_Cmp_Bit(bench *testing.B) {
	b1 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdeffedcba9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	b2 := big.NewInt(0).SetBytes(common.Hex2Bytes("0123456789abcdefaaaaaa9876543210f2f3f4f5f6f7f8f9fff3f4f5f6f7f8f9"))
	f, _ := NewFixedFromBig(b1)
	f2, _ := NewFixedFromBig(b2)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f.Cmp(f2)
	}
}
func Benchmark_Cmp(bench *testing.B) {
	bench.Run("big", benchmark_Cmp_Big)
	bench.Run("fixedbit", benchmark_Cmp_Bit)
}


func benchmark_Lsh_Big(bench *testing.B) {
	original := big.NewInt(0).SetBytes(common.Hex2Bytes("FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"))
	n := uint(255)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1 := big.NewInt(0)
		b1.Lsh(original, n)
	}
}
func benchmark_Lsh_Bit(bench *testing.B) {
	original := big.NewInt(0).SetBytes(common.Hex2Bytes("FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"))
	n := uint(255)
	f2, _ := NewFixedFromBig(original)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f1 := NewFixed()
		f1.Lsh(f2, n)
	}
}
func Benchmark_Lsh(bench *testing.B) {
	bench.Run("big", benchmark_Lsh_Big)
	bench.Run("fixedbit", benchmark_Lsh_Bit)
}

func benchmark_Rsh_Big(bench *testing.B) {
	original := big.NewInt(0).SetBytes(common.Hex2Bytes("FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"))
	n := uint(255)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1 := big.NewInt(0)
		b1.Rsh(original, n)
	}
}
func benchmark_Rsh_Bit(bench *testing.B) {
	original := big.NewInt(0).SetBytes(common.Hex2Bytes("FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"))
	f2, _ := NewFixedFromBig(original)
	n := uint(255)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f1 := NewFixed()
		f1.Rsh(f2, n)
	}
}
func Benchmark_Rsh(bench *testing.B) {
	bench.Run("big", benchmark_Rsh_Big)
	bench.Run("fixedbit", benchmark_Rsh_Bit)
}

func benchmark_Exp_Big(bench *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	orig := big.NewInt(0).SetBytes(common.Hex2Bytes(x))
	base := big.NewInt(0).SetBytes(common.Hex2Bytes(x))
	exp := big.NewInt(0).SetBytes(common.Hex2Bytes(y))

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		Exp(base, exp)
		base.Set(orig)
	}
}
func benchmark_Exp_Bit(bench *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	base := big.NewInt(0).SetBytes(common.Hex2Bytes(x))
	exp := big.NewInt(0).SetBytes(common.Hex2Bytes(y))

	f_base, _ := NewFixedFromBig(base)
	f_orig, _ := NewFixedFromBig(base)
	f_exp, _ := NewFixedFromBig(exp)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		ExpF(f_base, f_exp)
		f_base.Copy(f_orig)
	}
}
func benchmark_ExpSmall_Big(bench *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "8abcdef"

	orig := big.NewInt(0).SetBytes(common.Hex2Bytes(x))
	base := big.NewInt(0).SetBytes(common.Hex2Bytes(x))
	exp := big.NewInt(0).SetBytes(common.Hex2Bytes(y))

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		Exp(base, exp)
		base.Set(orig)
	}
}
func benchmark_ExpSmall_Bit(bench *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "8abcdef"

	base := big.NewInt(0).SetBytes(common.Hex2Bytes(x))
	exp := big.NewInt(0).SetBytes(common.Hex2Bytes(y))

	f_base, _ := NewFixedFromBig(base)
	f_orig, _ := NewFixedFromBig(base)
	f_exp, _ := NewFixedFromBig(exp)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		ExpF(f_base, f_exp)
		f_base.Copy(f_orig)
	}
}
func Benchmark_Exp(bench *testing.B) {
	bench.Run("large/big", benchmark_Exp_Big)
	bench.Run("large/fixedbit", benchmark_Exp_Bit)
	bench.Run("small/big", benchmark_ExpSmall_Big)
	bench.Run("small/fixedbit", benchmark_ExpSmall_Bit)
}

func Benchmark_Div(bench *testing.B) {
	bench.Run("large/big", benchmark_DivLarge_Big)
	bench.Run("large/fixedbit", benchmark_DivLarge_Bit)
	bench.Run("small/big", benchmark_DivSmall_Big)
	bench.Run("small/fixedbit", benchmark_DivSmall_Bit)
}

func benchmark_DivSmall_Big(bench *testing.B) {
	a := big.NewInt(0).SetBytes(common.Hex2Bytes("1fc2bad1e611"))
	b := big.NewInt(0).SetBytes(common.Hex2Bytes("12bad1e611"))

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1 := big.NewInt(0)
		b1.Div(a, b)
		U256(b1)
	}
}

func benchmark_DivSmall_Bit(bench *testing.B) {
	a := big.NewInt(0).SetBytes(common.Hex2Bytes("1fc2bad1e611"))
	b := big.NewInt(0).SetBytes(common.Hex2Bytes("12bad1e611"))
	fa, _ := NewFixedFromBig(a)
	fb, _ := NewFixedFromBig(b)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f := NewFixed()
		f.Div(fa, fb)
	}
}
func benchmark_DivLarge_Big(bench *testing.B) {
	a := big.NewInt(0).SetBytes(common.Hex2Bytes("fe7fb0d1f59dfe9492ffbf73683fd1e870eec79504c60144cc7f5fc2bad1e611"))
	b := big.NewInt(0).SetBytes(common.Hex2Bytes("ff3f9014f20db29ae04af2c2d265de17"))

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		b1 := big.NewInt(0)
		b1.Div(a, b)
		U256(b1)
	}
}

func benchmark_DivLarge_Bit(bench *testing.B) {
	a := big.NewInt(0).SetBytes(common.Hex2Bytes("fe7fb0d1f59dfe9492ffbf73683fd1e870eec79504c60144cc7f5fc2bad1e611"))
	b := big.NewInt(0).SetBytes(common.Hex2Bytes("ff3f9014f20db29ae04af2c2d265de17"))
	fa, _ := NewFixedFromBig(a)
	fb, _ := NewFixedFromBig(b)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		f := NewFixed()
		f.Div(fa, fb)
	}
}
