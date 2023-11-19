package evaluation

import (
	"rosaline/internal/chess"
)

type Evaluator struct {
}

func NewEvaluator() Evaluator {
	return Evaluator{}
}

// evaluateSide determines a score for a side based on the value of the pieces it has.
func (e Evaluator) evaluateSide(position chess.Position, color chess.Color) int {
	score := 0

	colorBB := position.GetColorBB(color)
	for colorBB > 0 {
		square := chess.Square(colorBB.PopLsb())
		piece, _ := position.GetPieceAt(square)
		score += pieceValue(piece)
		score += e.getSquareScore(position, piece, square)
	}

	pawnBB := colorBB & position.GetPieceBB(chess.Pawn)
	for _, bb := range chess.FileBitBoards {
		pawns := pawnBB & bb
		if pawns.PopulationCount() >= 2 {
			score += doublePawnPenalty
		}
	}

	bishopBB := colorBB & position.GetPieceBB(chess.Bishop)
	if bishopBB.PopulationCount() >= 2 {
		score += doubleBishopBonus
	}

	return score
}

func (e Evaluator) getSquareScore(position chess.Position, piece chess.Piece, square chess.Square) int {
	if position.Phase() == chess.EndgamePhase {
		scoreBoard, ok := openingScores[piece]
		if ok {
			return scoreBoard[square]
		}
	} else {
		scoreBoard, ok := endgameScores[piece]
		if ok {
			return scoreBoard[square]
		}
	}

	return 0
}

func (e Evaluator) Evaluate(position chess.Position) int {
	if position.IsDraw() {
		return DrawScore
	}

	turn := position.Turn()
	if position.IsCheckmated(turn) {
		return MateScore * evaluationMultiplier(turn.OpposingSide())
	}

	if position.IsCheckmated(turn.OpposingSide()) {
		return MateScore * evaluationMultiplier(turn)
	}

	return e.evaluateSide(position, chess.White) - e.evaluateSide(position, chess.Black)
}
