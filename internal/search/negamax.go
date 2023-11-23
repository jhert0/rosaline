package search

import (
	"math"
	"rosaline/internal/chess"
	"rosaline/internal/evaluation"
)

const (
	initialAlpha = math.MinInt + 1
	initialBeta  = math.MaxInt - 1
)

type NegamaxSearcher struct {
	evaluator evaluation.Evaluator
	drawTable drawTable
}

func NewNegamaxSearcher(evaluator evaluation.Evaluator) NegamaxSearcher {
	return NegamaxSearcher{
		evaluator: evaluator,
		drawTable: newDrawTable(),
	}
}

func (s NegamaxSearcher) Search(position chess.Position, depth int) ScoredMove {
	moves := position.GenerateMoves(chess.LegalMoveGeneration)
	bestMove := ScoredMove{}
	bestScore := math.MinInt + 1

	for _, move := range moves {
		position.MakeMove(move)

		score := -s.doSearch(position, initialAlpha, initialBeta, depth-1, 0)
		scoredMove := NewScoredMove(move, score)

		if score > bestScore {
			bestScore = score
			bestMove = scoredMove
		}

		position.Undo()
	}

	return bestMove
}

func (s *NegamaxSearcher) doSearch(position chess.Position, alpha int, beta int, depth int, ply int) int {
	if s.drawTable.IsRepeat(position.Hash()) {
		return evaluation.DrawScore
	}

	if position.IsDraw() {
		return evaluation.DrawScore
	}

	moves := position.GenerateMoves(chess.LegalMoveGeneration)
	if len(moves) == 0 {
		if position.IsKingInCheck(position.Turn()) {
			return -evaluation.MateScore + ply
		} else {
			return evaluation.DrawScore
		}
	}

	if depth == 0 {
		return s.evaluator.AbsoluteEvaluation(position)
	}

	for _, move := range moves {
		s.drawTable.Push(position.Hash())

		position.MakeMove(move)
		score := -s.doSearch(position, -beta, -alpha, depth-1, ply+1)
		position.Undo()

		s.drawTable.Pop()

		if score >= beta {
			return beta
		}

		if score > alpha {
			alpha = score
		}
	}

	return alpha
}

// Reset clears any information about searched positions.
func (s *NegamaxSearcher) Reset() {
	s.drawTable.Clear()
}
