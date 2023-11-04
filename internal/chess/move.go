package chess

import (
	"math"
	"strings"
)

type MoveType uint8

const (
	NormalMove MoveType = iota
	CastleMove
	EnPassantMove
)

func (t MoveType) String() string {
	switch t {
	case NormalMove:
		return "Normal"
	case CastleMove:
		return "Castle"
	case EnPassantMove:
		return "En Passant"
	}

	return "<unknown>"
}

type MoveFlag uint16

const (
	NoMoveFlag        MoveFlag = 0x0000
	QuietMoveFlag     MoveFlag = 0x0001
	CaputureMoveFlag  MoveFlag = 0x0010
	PawnPushMoveFlag  MoveFlag = 0x0100
	PromotionMoveFlag MoveFlag = 0x1000
)

type Move struct {
	From Square
	To   Square

	moveType MoveType
	flags    MoveFlag

	promotionPiece Piece
	capturePiece   Piece
}

// NewMove creates a new move with the given from and to.
func NewMove(from, to Square, moveType MoveType, flags MoveFlag) Move {
	return Move{
		From:           from,
		To:             to,
		moveType:       moveType,
		flags:          flags,
		promotionPiece: EmptyPiece,
		capturePiece:   EmptyPiece,
	}
}

// WithPromotion sets that the move will result with the moving piece being promoted to the given piece.
func (m *Move) WithPromotion(piece Piece) {
	m.flags |= PromotionMoveFlag
	m.promotionPiece = piece
}

// WithCapture sets that the move will result in the given piece being captured.
func (m *Move) WithCapture(piece Piece) {
	m.flags |= CaputureMoveFlag
	m.capturePiece = piece
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
	from := float64(m.From.Rank())
	to := float64(m.To.Rank())
	return int(math.Abs(from - to))
}

// FileDifference calculates the difference in files between the from
// square and the to square.
func (m Move) FileDifference() int {
	from := float64(m.From.File())
	to := float64(m.To.File())
	return int(math.Abs(from - to))
}

// IsPromotion returns whether the move results in a promotion.
func (m Move) IsPromotion() bool {
	return m.promotionPiece != EmptyPiece
}

// Captures returns whether the move results in a capture.
func (m Move) Captures() bool {
	return m.capturePiece != EmptyPiece
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
	return m.HasFlag(PawnPushMoveFlag) || m.HasFlag(CaputureMoveFlag) || m.HasFlag(PromotionMoveFlag)
}

func (m Move) String() string {
	str := m.From.ToAlgebraic() + m.To.ToAlgebraic()

	if m.promotionPiece != EmptyPiece {
		character := string(m.promotionPiece.Character())
		str += strings.ToLower(character)
	}

	return str
}
