package search

import (
	"rosaline/internal/chess"
	"rosaline/internal/evaluation"
	"testing"
)

func BenchmarkNegamax(b *testing.B) {
	position, _ := chess.NewPosition(chess.StartingFen)
	evaluator := evaluation.NewEvaluator()
	searcher := NewNegamaxSearcher(evaluator)

	for i := 0; i < b.N; i++ {
		searcher.Search(position, 4)
	}
}
