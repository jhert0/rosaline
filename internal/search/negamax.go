package search

import (
	"cmp"
	"fmt"
	"math"
	"rosaline/internal/chess"
	"rosaline/internal/evaluation"
	"slices"
	"strings"
	"time"
)

const (
	initialAlpha = math.MinInt + 1
	initialBeta  = math.MaxInt - 1

	maxNumberKillerMoves = 128

	nullMovePruningReduction = 2

	window = 500

	MaxDepth = 16
)

type NegamaxSearcher struct {
	evaluator evaluation.Evaluator
	drawTable drawTable

	killerMoves     map[chess.Color][]chess.Move
	killerMoveIndex int

	ttable TranspositionTable

	pvtable  [MaxDepth][MaxDepth]chess.Move
	pvlength [MaxDepth]int

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

func (s *NegamaxSearcher) Search(position *chess.Position, depth int, print bool) chess.Move {
	s.ClearPreviousSearch()

	bestMove := chess.NullMove

	for d := 1; d <= depth; d++ {
		start := time.Now()

		score := s.doSearch(position, initialAlpha, initialBeta, d, 0, 0)

		elapsed := time.Since(start)

		bestMove = s.pvtable[0][0]

		if print {
			nps := float64(s.nodes) / float64(elapsed.Seconds())
			fmt.Printf("info depth %d score cp %d nodes %d nps %f pv %s time %d tbhits %d\n", d, score, s.nodes, nps, s.getPV(), elapsed.Milliseconds(), s.ttable.Hits())
		}

		if s.stop {
			break
		}
	}

	return bestMove
}

func (s NegamaxSearcher) getPV() string {
	var builder strings.Builder

	length := s.pvlength[0]
	for i := 0; i < length; i++ {
		builder.WriteString(s.pvtable[0][i].String())

		if i != length {
			builder.WriteString(" ")
		}
	}

	return builder.String()
}

func (s NegamaxSearcher) scoreMove(position *chess.Position, move chess.Move, ply int) int {
	turn := position.Turn()

	if s.pvtable[0][ply] == move {
		return 2000
	}

	if slices.Contains(s.killerMoves[turn], move) {
		return 1000
	}

	return 0
}

func (s *NegamaxSearcher) doSearch(position *chess.Position, alpha int, beta int, depth int, ply int, extensions int) int {
	s.pvlength[ply] = ply

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
		if entry.Depth >= depth && entry.Hash == position.Hash() && ply != 0 {
			switch entry.Type {
			case ExactNode:
				s.pvlength[ply] = ply + 1
				s.pvtable[ply][ply] = entry.Move
				return entry.Score
			case UpperNode:
				if entry.Score <= alpha {
					return alpha
				}

				break
			case LowerNode:
				if entry.Score >= beta {
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
		return cmp.Compare(s.scoreMove(position, m1, ply), s.scoreMove(position, m2, ply))
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

			if !move.IsCapture() {
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

			s.pvtable[ply][ply] = move

			for i := ply + 1; i < s.pvlength[ply+1]; i++ {
				pvMove := s.pvtable[ply+1][i]
				s.pvtable[ply][i] = pvMove
			}

			s.pvlength[ply] = s.pvlength[ply+1]
		}
	}

	if !s.stop {
		entry := NewTableEntry(position.Hash(), nodeType, bestMove, bestScore, depth, position.Plies())
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
	s.nodes = 0
	s.stop = false

	s.ttable.ResetCounters()

	clear(s.killerMoves)
	s.killerMoveIndex = 0

	s.pvtable = [MaxDepth][MaxDepth]chess.Move{}
	s.pvlength = [MaxDepth]int{}
}

// Reset clears any information about searched positions.
func (s *NegamaxSearcher) Reset() {
	s.drawTable.Clear()
	s.ClearPreviousSearch()
	s.ttable.Clear()
}
