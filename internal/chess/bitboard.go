package chess

import (
	"fmt"
	"math/bits"
)

type BitBoard uint64

// NewBitBoard creates a new BitBoard with the give value.
func NewBitBoard(value uint64) BitBoard {
	return BitBoard(value)
}

// SetBit sets the bit at the given bit number.
func (bb *BitBoard) SetBit(number uint64) {
	if number > 63 {
		panic(fmt.Sprintf("tried to set invalid bit number: %d", number))
	}

	*bb |= BitBoard(uint64(1) << number)
}

// ClearBit clears the bit at the given bit number.
func (bb *BitBoard) ClearBit(number uint64) {
	if number > 63 {
		panic(fmt.Sprintf("tried to clear invalid bit number: %d", number))
	}

	*bb &= BitBoard(^(uint64(1) << number))
}

// IsBitSet checks if the bit is set at the given bit number.
func (bb BitBoard) IsBitSet(number uint64) bool {
	if number > 63 {
		panic(fmt.Sprintf("tried to check value of invalid bit number: %d", number))
	}

	return (bb & BitBoard((uint64(1) << number))) > 0
}

// Lsb returns the bit number of the last bit that is set.
func (bb BitBoard) Lsb() int {
	return bits.TrailingZeros64(uint64(bb))
}

// Msb returns the bit number of the first bit that is set.
func (bb BitBoard) Msb() int {
	return bits.LeadingZeros64(uint64(bb))
}

// PopulationCount returns the number of bits that are set.
func (bb BitBoard) PopulationCount() int {
	return bits.OnesCount64(uint64(bb))
}

// PopLsb clears the least significant bit and returns the bit number it was at.
func (bb *BitBoard) PopLsb() int {
	lsb := bb.Lsb()
	bb.ClearBit(uint64(lsb))
	return lsb
}

// PopMsb clears the most significant bit and the returns the bit number it was at.
func (bb *BitBoard) PopMsb() int {
	msb := bb.Lsb()
	bb.ClearBit(uint64(msb))
	return msb
}

// NorthOne fills the square north of the current square on an otherwise empty BitBoard.
func (bb *BitBoard) NorthOne() {
	*bb <<= 8
}

// SouthOne fills the square south of the current square on an otherwise empty BitBoard.
func (bb *BitBoard) SouthOne() {
	*bb >>= 8
}

// FillNorth fills the squares above the set square on an otherwise empty BitBoard.
func (bb *BitBoard) FillNorth() {
	*bb |= *bb << 8
	*bb |= *bb << 16
	*bb |= *bb << 32
}

// FillSouth fills the squares below the set square on an otherwise empty BitBoard.
func (bb *BitBoard) FillSouth() {
	*bb |= *bb >> 8
	*bb |= *bb >> 16
	*bb |= *bb >> 32
}

// FillFile fills north and south squares on an otherwise empty BitBoard.
func (bb *BitBoard) FillFile() {
	bb.FillNorth()
	bb.FillSouth()
}

// Print prints the BitBoard value in binary.
func (bb BitBoard) Print() {
	fmt.Printf("%064b\n", bb)
}
