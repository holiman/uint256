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

const (
	cacheIndexBits = 8
	cacheWays      = 5

	cacheSets = 1 << cacheIndexBits
	cacheMask = cacheSets - 1
)

type cacheSet struct {
	rw  sync.RWMutex
	mod [cacheWays]Int
	inv [cacheWays][5]uint64
}

type reciprocalCache struct {
	set             [cacheSets]cacheSet
	hit             uint64
	miss            uint64
	fixedModulus    *Int
	fixedReciprocal [5]uint64
}

// NewCache returns a new reciprocalCache.
func NewCache(fixedModulus *Int) *reciprocalCache {
	if fixedModulus != nil {
		return &reciprocalCache{
			fixedModulus:    fixedModulus,
			fixedReciprocal: reciprocal(*fixedModulus, nil),
		}
	}
	return &reciprocalCache{}
}

func (c *reciprocalCache) Stats() (hit, miss uint64) {
	return c.hit, c.miss
}

var (
	// FixedModulusCurveFoo is the fixed modulus for a curve.
	FixedModulusCurveFoo, _ = FromHex("0xffffffff00000001000000000000000000000000ffffffffffffffffffffffff")
)

func (cache *reciprocalCache) has(m Int, index uint64, dest *[5]uint64) bool {
	if cache.fixedModulus != nil && m.Eq(cache.fixedModulus) {
		dest[0] = cache.fixedReciprocal[0]
		dest[1] = cache.fixedReciprocal[1]
		dest[2] = cache.fixedReciprocal[2]
		dest[3] = cache.fixedReciprocal[3]
		dest[4] = cache.fixedReciprocal[4]

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
