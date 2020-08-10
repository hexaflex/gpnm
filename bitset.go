package pnm

import (
	"math"
)

const (
	wordBitSize  = 8
	wordByteSize = 1
)

// bitset defines a set of bit values.
type bitset []uint8

// newBitset creates a new bitset of the given size.
func newBitset(bits uint) bitset {
	size := int(math.Ceil((float64(bits) / wordBitSize)))
	return make(bitset, size)
}

// Set sets the bit at the given index.
func (b bitset) Set(i int) {
	w := i / wordBitSize
	if i < 0 || w >= len(b) {
		return
	}

	bit := uint8(1 << (8 - uint(i%wordBitSize)))
	b[w] &^= bit
	b[w] ^= bit
}

// Test returns true if the bit at the given index is set.
func (b bitset) Test(i int) bool {
	w := i / wordBitSize
	return i >= 0 && w < len(b) && ((b[w]>>(8-uint(i%wordBitSize)))&1) == 1
}
