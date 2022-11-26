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

	promotionPiece PieceType
	capturePiece   PieceType
}

// NewMove creates a new move with the given from and to.
func NewMove(from, to Square, moveType MoveType) Move {
	return Move{
		From:           from,
		To:             to,
		moveType:       moveType,
		promotionPiece: None,
		capturePiece:   None,
	}
}

func (m *Move) WithPromotion(pieceType PieceType) {
	m.promotionPiece = pieceType
}

func (m *Move) WithCapture(pieceType PieceType) {
	m.capturePiece = pieceType
}

func (m Move) Type() MoveType {
	return m.moveType
}

func (m Move) PromotionPiece() PieceType {
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
	return m.promotionPiece != None
}

func (m Move) Captures() bool {
	return m.capturePiece == None
}
