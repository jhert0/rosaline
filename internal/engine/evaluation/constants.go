package evaluation

import (
	"rosaline/internal/chess"
	"slices"
)

const (
	DrawScore int = 0
	MateScore int = 50000
)

// Bonuses

const (
	doubleBishopBonus = 20
)

// Penalties

const (
	doublePawnPenalty = -10
)

// Square scores
// All of these are from whites perspective.

var pawnSquareScores = []int{
	0, 0, 0, 0, 0, 0, 0, 0,
	5, 10, 10, -10, -10, 10, 10, 5,
	5, 0, 0, 5, 5, 0, 0, 5,
	0, 0, 0, 20, 20, 0, 0, 0,
	0, 0, 0, 25, 25, 0, 0, 0,
	30, 30, 30, 40, 40, 30, 30, 30,
	50, 50, 50, 50, 50, 50, 50, 50,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var knightSquareScores = []int{
	-50, -20, -20, -20, -20, -20, -20, -50,
	-30, 0, 0, 5, 5, 0, 0, -30,
	-20, 0, 5, 0, 0, 5, 0, -20,
	-20, 0, 20, 20, 20, 20, 0, -20,
	-20, 0, 0, 20, 20, 0, 0, -20,
	-20, 0, 0, 0, 0, 0, 0, -20,
	-30, 0, 0, 0, 0, 0, 0, -30,
	-50, -20, -20, -20, -20, -20, -20, -50,
}

var bishopSquareScores = []int{
	-5, -10, -10, -10, -10, -10, -10, -5,
	-10, 5, 0, 5, 5, 0, 5, -10,
	-10, 10, 10, 10, 10, 10, 10, -10,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-5, -10, -10, -10, -10, -10, -10, -5,
}

var rookSquareScores = []int{
	0, 0, 0, 5, 5, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	5, 10, 10, 10, 10, 10, 10, 5,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var queenSquareScores = []int{
	-5, -10, -10, -5, -5, -10, -10, -5,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 5, 5, 5, 5, 5, 5, -10,
	5, 0, 0, 10, 10, 0, 0, -10,
	-10, 0, 0, 10, 10, 0, 0, -10,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-5, -10, -10, -5, -5, -10, -10, -5,
}

var kingSquareScores = []int{
	0, 0, 30, 0, 0, 20, 30, 0,
	-10, -10, -10, 0, 0, -10, -10, -10,
	-20, -20, -20, -10, -10, -20, -20, -20,
	-30, -30, -30, -30, -30, -30, -30, -30,
	-40, -40, -40, -40, -40, -40, -40, -40,
	-50, -50, -50, -50, -50, -50, -50, -50,
	-50, -50, -50, -50, -50, -50, -50, -50,
	-50, -50, -50, -50, -50, -50, -50, -50,
}

var squareScores = map[chess.Piece][]int{}

func init() {
	colors := []chess.Color{chess.White, chess.Black}
	types := []chess.PieceType{chess.Pawn, chess.Knight, chess.Bishop, chess.Rook, chess.Queen, chess.King}

	for _, color := range colors {
		for _, pieceType := range types {
			piece := chess.NewPiece(pieceType, color)

			var scoreBoard []int
			switch pieceType {
			case chess.Pawn:
				scoreBoard = slices.Clone(pawnSquareScores)
				break
			case chess.Knight:
				scoreBoard = slices.Clone(knightSquareScores)
				break
			case chess.Bishop:
				scoreBoard = slices.Clone(bishopSquareScores)
				break
			case chess.Rook:
				scoreBoard = slices.Clone(rookSquareScores)
				break
			case chess.Queen:
				scoreBoard = slices.Clone(queenSquareScores)
				break
			case chess.King:
				scoreBoard = slices.Clone(kingSquareScores)
				break
			}

			if color == chess.Black {
				slices.Reverse(scoreBoard)
			}

			squareScores[piece] = scoreBoard
		}
	}
}
