package chess

import (
	"fmt"
)

type BitBoard struct {
	value uint64
}

func NewBitBoard(value uint64) BitBoard {
	return BitBoard{
		value: value,
	}
}

func (bb BitBoard) Value() uint64 {
	return bb.value
}

func (bb *BitBoard) SetBit(number uint64) {
	bb.value = (1 << number) | bb.value
}

func (bb *BitBoard) ClearBit(number uint64) {
	bb.value = ^(1 << number) & bb.value
}

func (bb BitBoard) BitSet(number uint64) bool {
	return (bb.value & (1 << number)) > 0
}

func (bb BitBoard) Print() {
	fmt.Printf("%064b\n", bb.value)
}
