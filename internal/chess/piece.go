package chess

import "unicode"

type Piece uint8

const (
	NoPiece Piece = iota

	WhitePawn
	WhiteKnight
	WhiteBishop
	WhiteRook
	WhiteQueen
	WhiteKing

	BlackPawn
	BlackKnight
	BlackBishop
	BlackRook
	BlackQueen
	BlackKing
)

// Color returns the color of the piece.
func (p Piece) Color() Color {
	if p == NoPiece {
		return NoColor
	}

	if p <= 6 {
		return White
	}

	return Black
}

// Character returns the character associated with the piece.
func (p Piece) Character() rune {
	var character rune

	switch p {
	case WhitePawn, BlackPawn:
		character = 'p'
		break
	case WhiteKnight, BlackKnight:
		character = 'n'
		break
	case WhiteBishop, BlackBishop:
		character = 'b'
		break
	case WhiteRook, BlackRook:
		character = 'r'
		break
	case WhiteQueen, BlackQueen:
		character = 'q'
		break
	case WhiteKing, BlackKing:
		character = 'k'
		break
	}

	if p.Color() == White {
		return unicode.ToUpper(character)
	}

	return character
}
