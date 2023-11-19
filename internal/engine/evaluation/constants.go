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
	semiOpenFileBonus = 10
	openFileBonus     = semiOpenFileBonus * 2
)

// Penalties

const (
	doublePawnPenalty = -10
)

// Square scores
// All of these are from whites perspective.

var pawnSquareScores = [2][]int{
	// opening
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		5, 10, 10, -10, -10, 10, 10, 5,
		5, 0, 0, 5, 5, 0, 0, 5,
		0, 0, 0, 20, 20, 0, 0, 0,
		0, 0, 0, 25, 25, 0, 0, 0,
		30, 30, 30, 40, 40, 30, 30, 30,
		50, 50, 50, 50, 50, 50, 50, 50,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	// endgame
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		10, 10, 10, 10, 10, 10, 10, 10,
		20, 20, 20, 25, 25, 20, 20, 20,
		30, 30, 30, 35, 35, 30, 30, 30,
		30, 30, 30, 40, 40, 30, 30, 30,
		50, 50, 50, 50, 50, 50, 50, 50,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
}

var knightSquareScores = [2][]int{
	// opening
	{
		-50, -20, -20, -20, -20, -20, -20, -50,
		-30, 0, 0, 5, 5, 0, 0, -30,
		-20, 0, 5, 0, 0, 5, 0, -20,
		-20, 0, 20, 20, 20, 20, 0, -20,
		-20, 0, 0, 20, 20, 0, 0, -20,
		-20, 0, 0, 0, 0, 0, 0, -20,
		-30, 0, 0, 0, 0, 0, 0, -30,
		-50, -20, -20, -20, -20, -20, -20, -50,
	},
	// endgame
	{
		-50, -20, -20, -20, -20, -20, -20, -50,
		-30, 0, 0, 5, 5, 0, 0, -30,
		-20, 0, 5, 0, 0, 5, 0, -20,
		-20, 0, 20, 20, 20, 20, 0, -20,
		-20, 0, 0, 20, 20, 0, 0, -20,
		-20, 0, 0, 0, 0, 0, 0, -20,
		-30, 0, 0, 0, 0, 0, 0, -30,
		-50, -20, -20, -20, -20, -20, -20, -50,
	},
}

var bishopSquareScores = [2][]int{
	// opening
	{
		-5, -10, -10, -10, -10, -10, -10, -5,
		-10, 5, 0, 5, 5, 0, 5, -10,
		-10, 10, 10, 10, 10, 10, 10, -10,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-5, -10, -10, -10, -10, -10, -10, -5,
	},
	// endgame
	{
		-5, -10, -10, -10, -10, -10, -10, -5,
		-10, 5, 0, 5, 5, 0, 5, -10,
		-10, 10, 10, 10, 10, 10, 10, -10,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-5, -10, -10, -10, -10, -10, -10, -5,
	},
}

var rookSquareScores = [2][]int{
	// opening
	{
		0, 0, 0, 5, 5, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		5, 10, 10, 10, 10, 10, 10, 5,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	// endgame
	{
		0, 0, 0, 5, 5, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		5, 10, 10, 10, 10, 10, 10, 5,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
}

var queenSquareScores = [2][]int{
	// opening
	{
		-5, -10, -10, -5, -5, -10, -10, -5,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 5, 5, 5, 5, 5, 5, -10,
		5, 0, 0, 10, 10, 0, 0, -10,
		-10, 0, 0, 10, 10, 0, 0, -10,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-5, -10, -10, -5, -5, -10, -10, -5,
	},
	// endgame
	{
		-5, -10, -10, -5, -5, -10, -10, -5,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 5, 5, 5, 5, 5, 5, -10,
		5, 0, 0, 10, 10, 0, 0, -10,
		-10, 0, 0, 10, 10, 0, 0, -10,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-5, -10, -10, -5, -5, -10, -10, -5,
	},
}

var kingSquareScores = [2][]int{
	// opening
	{
		0, 0, 30, 0, 0, 20, 30, 0,
		-10, -10, -10, 0, 0, -10, -10, -10,
		-20, -20, -20, -10, -10, -20, -20, -20,
		-30, -30, -30, -30, -30, -30, -30, -30,
		-40, -40, -40, -40, -40, -40, -40, -40,
		-50, -50, -50, -50, -50, -50, -50, -50,
		-50, -50, -50, -50, -50, -50, -50, -50,
		-50, -50, -50, -50, -50, -50, -50, -50,
	},
	// endgame
	{
		-50, -40, -30, -30, -30, -30, -40, -50,
		-30, -20, -10, -10, -10, -10, -20, -30,
		-30, -10, -10, -10, -10, -10, -10, -30,
		-30, -5, -5, 40, 40, -5, -5, -30,
		-30, -5, -5, 40, 40, -5, -5, -30,
		-30, 0, 0, 0, 0, 0, 0, -30,
		-30, 0, 0, 0, 0, 0, 0, -30,
		-50, -40, 0, 0, 0, 0, -40, -50,
	},
}

var openingScores = map[chess.Piece][]int{}
var endgameScores = map[chess.Piece][]int{}

func init() {
	initializeScoreBoard(openingScores, 0) // initialize opening square score board
	initializeScoreBoard(endgameScores, 1) // initialize endgame square score board
}

func initializeScoreBoard(board map[chess.Piece][]int, index int) {
	colors := []chess.Color{chess.White, chess.Black}
	types := []chess.PieceType{chess.Pawn, chess.Knight, chess.Bishop, chess.Rook, chess.Queen, chess.King}

	for _, color := range colors {
		for _, pieceType := range types {
			piece := chess.NewPiece(pieceType, color)

			var scoreBoard []int
			switch pieceType {
			case chess.Pawn:
				scoreBoard = slices.Clone(pawnSquareScores[index])
				break
			case chess.Knight:
				scoreBoard = slices.Clone(knightSquareScores[index])
				break
			case chess.Bishop:
				scoreBoard = slices.Clone(bishopSquareScores[index])
				break
			case chess.Rook:
				scoreBoard = slices.Clone(rookSquareScores[index])
				break
			case chess.Queen:
				scoreBoard = slices.Clone(queenSquareScores[index])
				break
			case chess.King:
				scoreBoard = slices.Clone(kingSquareScores[index])
				break
			}

			if color == chess.Black {
				slices.Reverse(scoreBoard)
			}

			board[piece] = scoreBoard
		}
	}
}
