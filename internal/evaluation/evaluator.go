package evaluation

import (
	"rosaline/internal/chess"
)

type Evaluator struct {
}

func NewEvaluator() Evaluator {
	return Evaluator{}
}

// evaluateSide determines the score for a side based:
//   - The value of it's pieces.
//   - The location of it's pieces.
//   - Bonuses/penalties for advantageous/disadvantageous positions.
func (e Evaluator) evaluateSide(position chess.Position, color chess.Color) int {
	score := 0

	colorBB := position.GetColorBB(color)
	for colorBB > 0 {
		square := chess.Square(colorBB.PopLsb())
		piece, _ := position.GetPieceAt(square)
		score += PieceValue(piece)
		score += e.getSquareScore(position, piece, square)
	}

	pawnBB := position.GetPieceBB(chess.Pawn)
	ourPawns := colorBB & pawnBB
	for _, bb := range chess.FileBitBoards {
		pawns := ourPawns & bb

		// give a penalty for two or more pawns of the same color being in a file
		if pawns.PopulationCount() >= 2 {
			score += doublePawnPenalty
		} else if pawns.PopulationCount() == 0 {
			score += semiOpenFileBonus
		}

		pawns = pawnBB & bb
		if pawns.PopulationCount() == 0 { // add a bonus for no pawns being in a file
			score += openFileBonus
		}
	}

	// add bonus for having 2 or more bishops
	bishopBB := colorBB & position.GetPieceBB(chess.Bishop)
	if bishopBB.PopulationCount() >= 2 {
		score += doubleBishopBonus
	}

	return score
}

// getSquareScore returns the score for the given piece being on the give square.
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

// Evaluate returns the score of the current position.
//
// The score will be:
//   - positive: if white is winning
//   - zero: if the position is a draw
//   - negative: if black is winning
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

// AbsoluteEvaluation evaluates the position relative to the player to move
// and returns the score.
//
// The score will be:
//   - positive: if the player to move is winning
//   - zero: if the position is a draw
//   - negative: if the player to move is losing
func (e Evaluator) AbsoluteEvaluation(position chess.Position) int {
	evaluation := e.Evaluate(position)
	return evaluation * evaluationMultiplier(position.Turn())
}
