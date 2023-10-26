package chess

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Square int8

const (
	A1 Square = iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1

	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2

	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3

	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4

	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5

	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6

	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7

	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8
)

// SquareFromAlgebraic creates a Square from the given algebraic string.
func SquareFromAlgebraic(algebraic string) (Square, error) {
	if len(algebraic) > 2 || len(algebraic) < 2 {
		return -1, errors.New("algebraic length is invalid")
	}

	// parse the file from the algebraic string
	var file int
	switch algebraic[0] {
	case 'a':
		file = 1
		break
	case 'b':
		file = 2
		break
	case 'c':
		file = 3
		break
	case 'd':
		file = 4
		break
	case 'e':
		file = 5
		break
	case 'f':
		file = 6
		break
	case 'g':
		file = 7
		break
	case 'h':
		file = 8
		break
	}

	// parse the rank from the algebraic string
	rank, err := strconv.Atoi(string(algebraic[1]))
	if err != nil {
		return -1, errors.New(fmt.Sprintf("invalid value: %c for rank", algebraic[1]))
	}

	return SquareFromRankFile(rank, file), nil
}

// SquareFromRankFile creates a Square from the rank file combination.
// It expects the actual rank/file numbers not zero index ones.
func SquareFromRankFile(rank, file int) Square {
	return Square(((rank - 1) * 8) + (file - 1))
}

// Rank calculates the rank for the Square.
func (s Square) Rank() int {
	return int(s/8) + 1
}

// File calculates the file for the Square.
func (s Square) File() int {
	return int(s%8) + 1
}

// IsValid returns if the Square is a valid square index on the board.
func (s Square) IsValid() bool {
	return s >= 0 && s < 64
}

// ToAlgebraic returns the algebraic string for the Square.
func (s Square) ToAlgebraic() string {
	if !s.IsValid() {
		return "<invalid square>"
	}

	algebraic := [2]string{}

	switch s.File() {
	case 1:
		algebraic[0] = "a"
		break
	case 2:
		algebraic[0] = "b"
		break
	case 3:
		algebraic[0] = "c"
		break
	case 4:
		algebraic[0] = "d"
		break
	case 5:
		algebraic[0] = "e"
		break
	case 6:
		algebraic[0] = "f"
		break
	case 7:
		algebraic[0] = "g"
		break
	case 8:
		algebraic[0] = "h"
		break
	}

	algebraic[1] = fmt.Sprintf("%d", s.Rank())

	return strings.Join(algebraic[:], "")
}

func (s Square) String() string {
	return s.ToAlgebraic()
}

// FileDistance returns the distance in file between two sqaures.
func FileDistance(sq1, sq2 Square) int {
	return int(math.Abs(float64(sq1.File() - sq2.File())))
}

// RankDistance returns the distance in ranks between two squares.
func RankDistance(sq1, sq2 Square) int {
	return int(math.Abs(float64(sq1.Rank() - sq2.Rank())))
}

// RankFileDifference returns the distnace in rank and file between two sqaures.
func RankFileDifference(sq1, sq2 Square) (int, int) {
	return RankDistance(sq1, sq2), FileDistance(sq1, sq2)
}

// SurroundingSquares returns the squares that surround the given square.
func SurroundingSquares(square Square) []Square {
	squares := make([]Square, 0, 8)

	for _, direction := range directions {
		if square.File() == 1 && (direction == west || direction == southwest || direction == northwest) {
			continue
		}

		if square.File() == 8 && (direction == east || direction == southeast || direction == northeast) {
			continue
		}

		if square.Rank() == 1 && (direction == south || direction == southwest || direction == southeast) {
			continue
		}

		if square.Rank() == 8 && (direction == north || direction == northwest || direction == northeast) {
			continue
		}

		newSquare := square + Square(direction)
		squares = append(squares, newSquare)
	}

	return squares
}
