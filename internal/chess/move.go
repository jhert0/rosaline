package chess

import "math"

type MoveType uint8

const (
	NormalMove MoveType = iota
	CastleMove
	EnPassantMove
)

type Move struct {
	From Square
	To   Square

	moveType MoveType

	promotionPiece Piece
	capturePiece   Piece
}

// NewMove creates a new move with the given from and to.
func NewMove(from, to Square, moveType MoveType) Move {
	return Move{
		From:           from,
		To:             to,
		moveType:       moveType,
		promotionPiece: EmptyPiece,
		capturePiece:   EmptyPiece,
	}
}

func (m *Move) WithPromotion(piece Piece) {
	m.promotionPiece = piece
}

func (m *Move) WithCapture(piece Piece) {
	m.capturePiece = piece
}

func (m Move) Type() MoveType {
	return m.moveType
}

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

func (m Move) IsPromotion() bool {
	return m.promotionPiece != EmptyPiece
}

func (m Move) Captures() bool {
	return m.capturePiece == EmptyPiece
}

func (m Move) String() string {
	str := m.From.ToAlgebraic() + m.To.ToAlgebraic()

	if m.promotionPiece != EmptyPiece {
		str += string(m.promotionPiece.Character())
	}

	return str
}
