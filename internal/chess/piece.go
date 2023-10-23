package chess

import (
	"fmt"
	"unicode"
)

type Piece uint8

type PieceType uint8

const (
	None   PieceType = 0x00
	Pawn   PieceType = 0x10
	Knight PieceType = 0x20
	Bishop PieceType = 0x30
	Rook   PieceType = 0x40
	Queen  PieceType = 0x50
	King   PieceType = 0x60

	EmptyPiece Piece = Piece(uint8(None) | uint8(NoColor))
)

var promotablePieces = []PieceType{Knight, Bishop, Rook, Queen}

// NewPiece creates a piece with the given type and color.
func NewPiece(pieceType PieceType, color Color) Piece {
	return Piece(uint8(pieceType) | uint8(color))
}

// Color returns the color of the piece.
func (p Piece) Color() Color {
	return Color(p & 0x0F)
}

// PieceType returns the type of the piece.
func (p Piece) Type() PieceType {
	return PieceType(p & 0xF0)
}

// Character returns the character associated with the piece.
func (p Piece) Character() rune {
	var character rune

	switch p.Type() {
	case Pawn:
		character = 'p'
		break
	case Knight:
		character = 'n'
		break
	case Bishop:
		character = 'b'
		break
	case Rook:
		character = 'r'
		break
	case Queen:
		character = 'q'
		break
	case King:
		character = 'k'
		break
	}

	if p.Color() == White {
		return unicode.ToUpper(character)
	}

	return character
}

// Value returns the value of the piece.
func (p Piece) Value() uint8 {
	switch p.Type() {
	case Pawn:
		return 1
	case Knight:
	case Bishop:
		return 3
	case Rook:
		return 5
	case Queen:
		return 9
	case King:
		return 127
	}

	panic(fmt.Sprintf("Unknown piece type '%v' encountered in Value()", p.Type()))
}
