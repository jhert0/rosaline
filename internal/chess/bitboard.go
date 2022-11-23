package chess

import (
	"fmt"
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

func (bb BitBoard) Print() {
	fmt.Printf("%064b\n", bb)
}
