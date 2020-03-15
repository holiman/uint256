package uint256

import (
	"math/bits"
	"testing"
)

func referenceReciprocal2by1(d uint64) uint64 {
	reciprocal, _ := bits.Div64(^d, ^uint64(0), d)
	return reciprocal
}

func TestReciprocal2by1(t *testing.T) {
	for d := uint64(0x8000000000000000); d < 0x800000000000ffff; d++ {
		reciprocal := reciprocal2by1(d)
		expected := referenceReciprocal2by1(d)
		if reciprocal != expected {
			t.Fatalf("wrong reciprocal 2by1: %x, expected: %x", reciprocal, expected)
		}
	}
}
