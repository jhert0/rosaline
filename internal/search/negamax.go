package search

import (
	"cmp"
	"math"
	"rosaline/internal/chess"
	"rosaline/internal/evaluation"
	"slices"
)

const (
	initialAlpha = math.MinInt + 1
	initialBeta  = math.MaxInt - 1

	maxNumberKillerMoves = 128
)

type NegamaxSearcher struct {
	evaluator evaluation.Evaluator
	drawTable drawTable

	killerMoves     map[chess.Color][]chess.Move
	killerMoveIndex int
}

func NewNegamaxSearcher(evaluator evaluation.Evaluator) NegamaxSearcher {
	return NegamaxSearcher{
		evaluator:       evaluator,
		drawTable:       newDrawTable(),
		killerMoves:     make(map[chess.Color][]chess.Move),
		killerMoveIndex: 0,
	}
}

func (s NegamaxSearcher) Search(position chess.Position, depth int) ScoredMove {
	moves := position.GenerateMoves(chess.LegalMoveGeneration)
	bestMove := ScoredMove{}
	bestScore := math.MinInt + 1

	for _, move := range moves {
		position.MakeMove(move)

		score := -s.doSearch(position, initialAlpha, initialBeta, depth-1, 0, 0)
		scoredMove := NewScoredMove(move, score)

		if score > bestScore {
			bestScore = score
			bestMove = scoredMove
		}

		position.Undo()
	}

	return bestMove
}

func (s NegamaxSearcher) scoreMove(position chess.Position, move chess.Move) int {
	turn := position.Turn()
	if slices.Contains(s.killerMoves[turn], move) {
		return 1000
	}

	return 0
}

func (s *NegamaxSearcher) doSearch(position chess.Position, alpha int, beta int, depth int, ply int, extensions int) int {
	if s.drawTable.IsRepeat(position.Hash()) {
		return evaluation.DrawScore
	}

	if position.IsDraw() {
		return evaluation.DrawScore
	}

	inCheck := position.IsKingInCheck(position.Turn())
	if inCheck && extensions < 2 { // limit the number of depth increases we will do to 2
		depth++
		extensions++
	}

	moves := position.GenerateMoves(chess.LegalMoveGeneration)
	if len(moves) == 0 {
		if inCheck {
			return -evaluation.MateScore + ply
		} else {
			return evaluation.DrawScore
		}
	}

	if depth == 0 {
		return s.evaluator.AbsoluteEvaluation(position)
	}

	slices.SortFunc(moves, func(m1, m2 chess.Move) int {
		return cmp.Compare(s.scoreMove(position, m1), s.scoreMove(position, m2))
	})

	for _, move := range moves {
		s.drawTable.Push(position.Hash())

		position.MakeMove(move)
		score := -s.doSearch(position, -beta, -alpha, depth-1, ply+1, extensions)
		position.Undo()

		s.drawTable.Pop()

		if score >= beta {
			if !move.HasFlag(chess.CaputureMoveFlag) {
				turn := position.Turn()
				length := len(s.killerMoves[turn])
				if length >= maxNumberKillerMoves {
					if s.killerMoveIndex >= (maxNumberKillerMoves - 1) {
						s.killerMoveIndex = 0
					}

					if !slices.Contains(s.killerMoves[turn], move) {
						s.killerMoves[turn][s.killerMoveIndex] = move
						s.killerMoveIndex++
					}
				} else {
					if !slices.Contains(s.killerMoves[turn], move) {
						s.killerMoves[turn] = append(s.killerMoves[turn], move)
					}
				}
			}

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
	clear(s.killerMoves)
}
