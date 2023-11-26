package search

import (
	"cmp"
	"math"
	"rosaline/internal/chess"
	"rosaline/internal/evaluation"
	"slices"
	"time"
)

const (
	initialAlpha = math.MinInt + 1
	initialBeta  = math.MaxInt - 1

	maxNumberKillerMoves = 128

	nullMovePruningReduction = 2
)

type SearchResults struct {
	BestMove chess.Move
	Score    int
	Depth    int
	Nodes    int
	Time     time.Duration
	NPS      float64
	Hits     int
	Misses   int
}

type NegamaxSearcher struct {
	evaluator evaluation.Evaluator
	drawTable drawTable

	killerMoves     map[chess.Color][]chess.Move
	killerMoveIndex int

	ttable TranspositionTable

	stop bool

	nodes int
}

func NewNegamaxSearcher(evaluator evaluation.Evaluator) NegamaxSearcher {
	return NegamaxSearcher{
		evaluator:       evaluator,
		drawTable:       newDrawTable(),
		killerMoves:     make(map[chess.Color][]chess.Move),
		killerMoveIndex: 0,
		ttable:          NewTranspositionTable(),
		nodes:           0,
	}
}

func (s NegamaxSearcher) Search(position *chess.Position, depth int) SearchResults {
	s.nodes = 0
	s.stop = false

	s.ttable.ResetCounters()

	start := time.Now()

	bestMove := chess.Move{}
	bestScore := math.MinInt + 1

	moves := position.GenerateMoves(chess.LegalMoveGeneration)

	slices.SortFunc(moves, func(m1, m2 chess.Move) int {
		return cmp.Compare(s.scoreMove(position, m1), s.scoreMove(position, m2))
	})

	for _, move := range moves {
		position.MakeMove(move)
		score := -s.doSearch(position, initialAlpha, initialBeta, depth-1, 0, 0)
		position.Undo()

		if score > bestScore {
			bestScore = score
			bestMove = move
		}

		if s.stop {
			break
		}
	}

	elapsed := time.Since(start)
	nps := float64(s.nodes) / float64(elapsed.Milliseconds())

	return SearchResults{
		BestMove: bestMove,
		Score:    bestScore,
		Depth:    depth,
		Nodes:    s.nodes,
		Time:     elapsed,
		NPS:      nps,
		Hits:     s.ttable.hits,
		Misses:   s.ttable.misses,
	}
}

func (s NegamaxSearcher) scoreMove(position *chess.Position, move chess.Move) int {
	turn := position.Turn()
	if slices.Contains(s.killerMoves[turn], move) {
		return 1000
	}

	return 0
}

func (s *NegamaxSearcher) doSearch(position *chess.Position, alpha int, beta int, depth int, ply int, extensions int) int {
	if s.stop {
		return 0
	}

	if s.drawTable.IsRepeat(position.Hash()) {
		return evaluation.DrawScore
	}

	if position.IsDraw() {
		return evaluation.DrawScore
	}

	pvNode := beta-alpha != 1

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
		if inCheck { // don't go in quiescence search when in check
			return s.evaluator.AbsoluteEvaluation(position)
		} else {
			return s.quiescence(position, alpha, beta)
		}
	}

	s.nodes++

	entry, ok := s.ttable.Get(position.Hash())
	if ok {
		if entry.Depth >= depth {
			switch entry.Type {
			case ExactNode:
				return entry.Score
			case UpperNode:
				if entry.Score < alpha {
					return alpha
				}

				break
			case LowerNode:
				if entry.Score > beta {
					return beta
				}

				break
			}
		}
	}

	// null move pruning
	doNullPruning := !inCheck && !pvNode
	if doNullPruning && depth >= 3 && ply != 0 {
		s.drawTable.Push(position.Hash())

		position.MakeNullMove()
		score := -s.doSearch(position, -beta, -beta+1, depth-1-nullMovePruningReduction, ply+1, extensions)
		position.Undo()

		s.drawTable.Pop()

		if score >= beta {
			return beta
		}
	}

	slices.SortFunc(moves, func(m1, m2 chess.Move) int {
		return cmp.Compare(s.scoreMove(position, m1), s.scoreMove(position, m2))
	})

	bestMove := chess.NullMove
	bestScore := math.MinInt
	nodeType := UpperNode

	for _, move := range moves {
		s.drawTable.Push(position.Hash())

		position.MakeMove(move)
		score := -s.doSearch(position, -beta, -alpha, depth-1, ply+1, extensions)
		position.Undo()

		s.drawTable.Pop()

		if score > bestScore {
			bestScore = score
			bestMove = move
		}

		if score >= beta {
			nodeType = LowerNode

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

			break
		}

		if score > alpha {
			alpha = score
			nodeType = ExactNode
		}
	}

	if !s.stop {
		entry := NewTableEntry(nodeType, bestMove, bestScore, depth, position.Plies())
		s.ttable.Insert(position.Hash(), entry)
	}

	return bestScore
}

func (s NegamaxSearcher) quiescence(position *chess.Position, alpha int, beta int) int {
	evaluation := s.evaluator.AbsoluteEvaluation(position)
	if evaluation >= beta {
		return beta
	}

	if alpha < evaluation {
		alpha = evaluation
	}

	captures := position.GenerateMoves(chess.CaptureMoveGeneration)
	for _, capture := range captures {
		position.MakeMove(capture)
		score := -s.quiescence(position, -beta, -alpha)
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

func (s *NegamaxSearcher) Stop() {
	s.stop = true
}

func (s NegamaxSearcher) Stopped() bool {
	return s.stop
}

func (s *NegamaxSearcher) ClearPreviousSearch() {
	clear(s.killerMoves)
	s.killerMoveIndex = 0
}

// Reset clears any information about searched positions.
func (s *NegamaxSearcher) Reset() {
	s.drawTable.Clear()
	s.ClearPreviousSearch()
	s.ttable.Clear()
}
