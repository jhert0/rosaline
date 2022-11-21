package chess

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Square int

// SquareFromAlgebraic creates a Square from the given algebraic string.
func SquareFromAlgebraic(algebraic string) (Square, error) {
	if len(algebraic) > 2 || len(algebraic) < 2 {
		return 0, errors.New("algebraic length is invalid")
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
		return 0, errors.New(fmt.Sprintf("invalid value: %c for rank", algebraic[1]))
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

// ToAlgebraic returns the algebraic string for the Square.
func (s Square) ToAlgebraic() string {
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
