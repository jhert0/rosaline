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
		return NOCOLOR
	}

	if p <= 6 {
		return WHITE
	}

	return BLACK
}

// Character returns the character associated with the piece.
func (p Piece) Character() rune {
	var character rune

	switch p {
	case WhitePawn:
	case BlackPawn:
		character = 'p'
		break
	case WhiteKnight:
	case BlackKnight:
		character = 'n'
		break
	case WhiteBishop:
	case BlackBishop:
		character = 'b'
		break
	case WhiteRook:
	case BlackRook:
		character = 'r'
		break
	case WhiteQueen:
	case BlackQueen:
		character = 'q'
		break
	case WhiteKing:
	case BlackKing:
		character = 'k'
		break
	}

	if p.Color() == WHITE {
		return unicode.ToUpper(character)
	}

	return character
}
