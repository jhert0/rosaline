package evaluation

import "rosaline/internal/chess"

type basicEvaluator struct {
}

func NewBasicEvaluator() basicEvaluator {
	return basicEvaluator{}
}

// evaluateSide determines a score for a side based on the value of the pieces it has.
func (e basicEvaluator) evaluateSide(position chess.Position, color chess.Color) int {
	score := 0

	colorBB := position.GetColorBB(color)
	for colorBB > 0 {
		square := chess.Square(colorBB.PopLsb())
		piece, _ := position.GetPieceAt(square)
		score += pieceValue(piece)
	}

	return score
}

func (e basicEvaluator) Evaluate(position chess.Position) int {
	return e.evaluateSide(position, chess.White) - e.evaluateSide(position, chess.Black)
}
