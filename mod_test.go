// uint256: Fixed size 256-bit math library
// Copyright 2021 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

package uint256

import "testing"

func TestLeadingZeros(t *testing.T) {
	one := Int{value: [4]uint64{1, 0, 0, 0}, tainted: false}

	testCases := []Int{
		Int{value: [4]uint64{0, 0, 0, 0}, tainted: false},
		Int{value: [4]uint64{1, 0, 0, 0}, tainted: false},
		Int{value: [4]uint64{0x7fffffffffffffff, 0, 0, 0}, tainted: false},
		Int{value: [4]uint64{0x8000000000000000, 0, 0, 0}, tainted: false},
		Int{value: [4]uint64{0xffffffffffffffff, 0, 0, 0}, tainted: false},
		Int{value: [4]uint64{0, 1, 0, 0}, tainted: false},
		Int{value: [4]uint64{0, 0x7fffffffffffffff, 0, 0}, tainted: false},
		Int{value: [4]uint64{0, 0x8000000000000000, 0, 0}, tainted: false},
		Int{value: [4]uint64{0, 0xffffffffffffffff, 0, 0}, tainted: false},
		Int{value: [4]uint64{0, 0, 1, 0}, tainted: false},
		Int{value: [4]uint64{0, 0, 0x7fffffffffffffff, 0}, tainted: false},
		Int{value: [4]uint64{0, 0, 0x8000000000000000, 0}, tainted: false},
		Int{value: [4]uint64{0, 0, 0xffffffffffffffff, 0}, tainted: false},
		Int{value: [4]uint64{0, 0, 0, 1}, tainted: false},
		Int{value: [4]uint64{0, 0, 0, 0x7fffffffffffffff}, tainted: false},
		Int{value: [4]uint64{0, 0, 0, 0x8000000000000000}, tainted: false},
		Int{value: [4]uint64{0, 0, 0, 0xffffffffffffffff}, tainted: false},
	}

	for _, x := range testCases {
		z := leadingZeros(&x)
		if z >= 0 && z < 256 {
			allZeros := new(Int).Rsh(&x, uint(256-z))
			oneBit := new(Int).Rsh(&x, uint(255-z))
			if allZeros.IsZero() && oneBit.Eq(&one) {
				continue
			}
		} else if z == 256 {
			if x.IsZero() {
				continue
			}
		}
		t.Errorf("wrong leading zeros %d of %x", z, x.value)
	}
}
