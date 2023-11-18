package search

import (
	"math"
	"rosaline/internal/chess"
	"rosaline/internal/engine/evaluation"
)

const (
	initialAlpha = math.MinInt + 1
	initialBeta  = math.MaxInt - 1
)

type NegamaxSearcher struct {
	evaluator evaluation.Evaluator
}

func NewNegamaxSearcher(evaluator evaluation.Evaluator) NegamaxSearcher {
	return NegamaxSearcher{
		evaluator: evaluator,
	}
}

func (s NegamaxSearcher) Search(position chess.Position, depth int) ScoredMove {
	moves := position.GenerateMoves(chess.LegalMoveGeneration)
	bestMove := ScoredMove{}
	bestScore := math.MinInt + 1

	for _, move := range moves {
		position.MakeMove(move)

		score := s.doSearch(position, initialAlpha, initialBeta, depth-1)
		scoredMove := NewScoredMove(move, score)

		if score > bestScore {
			bestScore = score
			bestMove = scoredMove
		}

		position.Undo()
	}

	return bestMove
}

func (s *NegamaxSearcher) doSearch(position chess.Position, alpha int, beta int, depth int) int {
	if depth == 0 {
		return s.evaluator.Evaluate(position)
	}

	moves := position.GenerateMoves(chess.LegalMoveGeneration)
	for _, move := range moves {
		position.MakeMove(move)

		score := -s.doSearch(position, -beta, -alpha, depth-1)

		position.Undo()

		if score >= beta {
			return beta
		}

		if score > alpha {
			alpha = score
		}
	}

	return alpha
}
