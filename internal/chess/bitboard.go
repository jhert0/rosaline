package chess

import (
	"fmt"
	"math/bits"
)

type BitBoard uint64

func NewBitBoard(value uint64) BitBoard {
	return BitBoard(value)
}

func (bb *BitBoard) SetBit(number uint64) {
	*bb |= BitBoard(uint64(1) << number)
}

func (bb *BitBoard) ClearBit(number uint64) {
	*bb &= BitBoard(^(uint64(1) << number))
}

func (bb BitBoard) BitSet(number uint64) bool {
	return (bb & BitBoard((uint64(1) << number))) > 0
}

func (bb BitBoard) TrailingZeros() int {
	return bits.TrailingZeros64(uint64(bb))
}

func (bb BitBoard) PopulationCount() int {
	return bits.OnesCount64(uint64(bb))
}

func (bb BitBoard) Print() {
	fmt.Printf("%064b\n", bb)
}
