// uint256: Fixed size 256-bit math library
// Copyright 2021 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

package uint256

import (
	"sync"
)

// Cache for reciprocal()
//
// Each cache set contains its own mutex to reduce lock contention in highly
// multithreaded settings.
//
// Adjust cacheIndexBits and cacheWays to scale the size of the cache.
// Total cache size is (24+72*cacheWays)*2^cacheIndexBits bytes, which is
// 96 KiB with cacheIndexBits = 8 and cacheWays = 5. There are also 16 bytes
// for hit and miss counters.
//
// Reasonable values are quite small, e.g. cacheIndexBits from 2 to 10, and
// cacheWays around 5. Note that cacheWays = 5+8n makes the set size an integer
// number of 64-byte cachelines.
//
// If you want to disable the cache, set cacheIndexBits and cacheWays to 0
//
// If you want to have a hardcoded modulus with precomputation, set
// fixedModulus = true and adjust the value stored in fixed_m. This can be
// used with or without the regular cache.

const (
	cacheIndexBits = 8
	cacheWays      = 5

	cacheSets = 1 << cacheIndexBits
	cacheMask = cacheSets - 1

	fixedModulus = true
)

type cacheSet struct {
	rw  sync.RWMutex
	mod [cacheWays]Int
	inv [cacheWays][5]uint64
}

type reciprocalCache struct {
	set  [cacheSets]cacheSet
	hit  uint64
	miss uint64
}

// NewCache returns a new reciprocalCache.
func NewCache() *reciprocalCache {
	return &reciprocalCache{}
}

func (c *reciprocalCache) Stats() (hit, miss uint64) {
	return c.hit, c.miss
}

var (
	// Fixed modulus and its reciprocal
	fixed_m Int
	fixed_r [5]uint64
)

func init() {
	if fixedModulus {
		// Initialise fixed modulus
		fixed_m[3] = 0xffffffff00000001
		fixed_m[2] = 0x0000000000000000
		fixed_m[1] = 0x00000000ffffffff
		fixed_m[0] = 0xffffffffffffffff

		fixed_r = reciprocal(fixed_m, nil)
	}
}

func (cache *reciprocalCache) has(m Int, index uint64, dest *[5]uint64) bool {
	if fixedModulus && m.Eq(&fixed_m) {
		dest[0] = fixed_r[0]
		dest[1] = fixed_r[1]
		dest[2] = fixed_r[2]
		dest[3] = fixed_r[3]
		dest[4] = fixed_r[4]

		return true
	}

	if cacheWays == 0 {
		return false
	}

	cache.set[index].rw.RLock()
	defer cache.set[index].rw.RUnlock()

	for w := 0; w < cacheWays; w++ {
		if cache.set[index].mod[w].Eq(&m) {
			copy(dest[:], cache.set[index].inv[w][:])
			cache.hit++
			return true
		}
	}
	cache.miss++
	return false
}

func (cache *reciprocalCache) put(m Int, index uint64, mu [5]uint64) {
	if cacheWays == 0 {
		return
	}

	cache.set[index].rw.Lock()
	defer cache.set[index].rw.Unlock()

	var w int
	for w = 0; w < cacheWays; w++ {
		if cache.set[index].mod[w].IsZero() {
			// Found an empty slot
			cache.set[index].mod[w] = m
			cache.set[index].inv[w] = mu
			return
		}
	}
	// Shift old elements, evicting the oldest
	for w = cacheWays - 1; w > 0; w-- {
		cache.set[index].mod[w] = cache.set[index].mod[w-1]
		cache.set[index].inv[w] = cache.set[index].inv[w-1]
	}
	// w == 0
	cache.set[index].mod[w] = m
	cache.set[index].inv[w] = mu
}
