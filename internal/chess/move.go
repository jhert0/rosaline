package chess

import (
	"math"
	"strings"
)

type MoveType uint8

const (
	QuietMove MoveType = iota
	CaptureMove
	EnPassantMove
	CastleMove
	Null
)

func (t MoveType) String() string {
	switch t {
	case QuietMove:
		return "Quiet"
	case CaptureMove:
		return "Capture"
	case CastleMove:
		return "Castle"
	case EnPassantMove:
		return "En Passant"
	case Null:
		return "Null"
	}

	return "<unknown>"
}

type MoveFlag uint8

const (
	NoMoveFlag       MoveFlag = 0b0000
	PawnPushMoveFlag MoveFlag = 0b0001
)

type Move struct {
	from Square
	to   Square

	moveType MoveType
	flags    MoveFlag

	promotionPiece Piece
}

var NullMove = NewMove(A1, A1, Null)

// NewMove creates a new move with the given from and to.
func NewMove(from, to Square, moveType MoveType) Move {
	return Move{
		from:           from,
		to:             to,
		moveType:       moveType,
		flags:          NoMoveFlag,
		promotionPiece: EmptyPiece,
	}
}

// From returns the square the piece is moving from.
func (m Move) From() Square {
	return m.from
}

// To returns the square the piece is moving to.
func (m Move) To() Square {
	return m.to
}

// WithPromotion sets that the move will result with the moving piece being promoted to the given piece.
func (m *Move) WithPromotion(piece Piece) {
	m.promotionPiece = piece
}

func (m *Move) WithFlags(flags MoveFlag) {
	m.flags |= flags
}

// Type returns the type of the move.
func (m Move) Type() MoveType {
	return m.moveType
}

// PromotionPiece returns the piece that the moving piece will be promoted to.
func (m Move) PromotionPiece() Piece {
	return m.promotionPiece
}

// RankDifference calculates the difference in ranks between the from
// square and the to square.
func (m Move) RankDifference() int {
	from := float64(m.from.Rank())
	to := float64(m.to.Rank())
	return int(math.Abs(from - to))
}

// FileDifference calculates the difference in files between the from
// square and the to square.
func (m Move) FileDifference() int {
	from := float64(m.from.File())
	to := float64(m.to.File())
	return int(math.Abs(from - to))
}

// IsPromotion returns whether the move results in a promotion.
func (m Move) IsPromotion() bool {
	return m.promotionPiece != EmptyPiece
}

// IsCapture returns whether the move results in a capture.
func (m Move) IsCapture() bool {
	return m.Type() == CaptureMove || m.Type() == EnPassantMove
}

// HasFlag checks if the move has that flag set.
func (m Move) HasFlag(flag MoveFlag) bool {
	return (m.flags & flag) > 0
}

// IsIrreversible returns whether the move can be reversed.
//
// Moves such as pawn moves and captures can't be reversed once done meaning
// the position before the was made is no longer possible.
func (m Move) IsIrreversible() bool {
	return m.HasFlag(PawnPushMoveFlag) || m.IsCapture() || m.IsPromotion()
}

func (m Move) String() string {
	if m.Type() == Null {
		return "0000"
	}

	str := m.from.ToAlgebraic() + m.to.ToAlgebraic()

	if m.promotionPiece != EmptyPiece {
		character := string(m.promotionPiece.Character())
		str += strings.ToLower(character)
	}

	return str
}
