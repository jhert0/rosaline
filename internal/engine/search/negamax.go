package search

import (
	"math"
	"rosaline/internal/chess"
	"rosaline/internal/engine/evaluation"
)

type negamaxSearcher struct {
	evaluator evaluation.Evaluator
}

func NewNegamaxSearcher(evaluator evaluation.Evaluator) negamaxSearcher {
	return negamaxSearcher{
		evaluator: evaluator,
	}
}

func (s negamaxSearcher) Search(position chess.Position, depth int) ScoredMove {
	moves := position.GenerateMoves(chess.LegalMoveGeneration)
	bestMove := ScoredMove{}
	max := math.MinInt

	for _, move := range moves {
		position.MakeMove(move)

		score := s.doSearch(position, depth)
		if score > max {
			max = score
			bestMove = NewScoredMove(move, score)
		}

		position.Undo()
	}

	return bestMove
}

func (s *negamaxSearcher) doSearch(position chess.Position, depth int) int {
	if depth == 0 {
		return s.evaluator.Evaluate(position)
	}

	moves := position.GenerateMoves(chess.LegalMoveGeneration)
	if len(moves) == 0 {
		return s.evaluator.Evaluate(position)
	}

	max := math.MinInt

	for _, move := range moves {
		position.MakeMove(move)

		score := -s.doSearch(position, depth-1)
		if score > max {
			max = score
		}

		position.Undo()
	}

	return max
}
