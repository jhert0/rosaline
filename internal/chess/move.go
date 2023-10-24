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

func (m *Move) WithPromotion(piece Piece) {
	m.flags |= PromotionMoveFlag
	m.promotionPiece = piece
}

func (m *Move) WithCapture(piece Piece) {
	m.flags |= CaputureMoveFlag
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

// HasFlag checks if the move has that flag set.
func (m Move) HasFlag(flag MoveFlag) bool {
	return (m.flags & flag) > 0
}

func (m Move) String() string {
	str := m.From.ToAlgebraic() + m.To.ToAlgebraic()

	if m.promotionPiece != EmptyPiece {
		character := string(m.promotionPiece.Character())
		str += strings.ToLower(character)
	}

	return str
}
