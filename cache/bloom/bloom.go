package bloom

import (
	"sync"
)

type BloomFilter struct {
	bits []uint64
	k    uint // number of hash functions
	m    uint // size of bitset
	mu   sync.RWMutex
}

func New(n uint, falsePositiveRate float64) *BloomFilter {
	m := uint(float64(n) * 1.44 * 10)
	if m < 64 {
		m = 64
	}
	k := uint(7)
	if falsePositiveRate > 0 {
		k = uint(-1.44*falsePositiveRate + 7)
	}
	if k < 1 {
		k = 1
	}
	if k > 30 {
		k = 30
	}

	// round m up to multiple of 64
	m = ((m + 63) / 64) * 64

	return &BloomFilter{
		bits: make([]uint64, m/64),
		k:    k,
		m:    m,
	}
}

func (bf *BloomFilter) Add(key string) {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	for i := uint(0); i < bf.k; i++ {
		idx := bf.hash(key, i) % bf.m
		bf.bits[idx/64] |= 1 << (idx % 64)
	}
}

func (bf *BloomFilter) MightContain(key string) bool {
	bf.mu.RLock()
	defer bf.mu.RUnlock()

	for i := uint(0); i < bf.k; i++ {
		idx := bf.hash(key, i) % bf.m
		if bf.bits[idx/64]&(1<<(idx%64)) == 0 {
			return false
		}
	}
	return true
}

func (bf *BloomFilter) Size() uint {
	return bf.m
}

func (bf *BloomFilter) hash(key string, seed uint) uint {
	h := uint(0)
	for i := 0; i < len(key); i++ {
		h = h*31 + uint(key[i]) + seed*0x9e3779b9
	}
	return h
}
