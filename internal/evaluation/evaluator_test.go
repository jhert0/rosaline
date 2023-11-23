package evaluation

import (
	"rosaline/internal/chess"
	"testing"
)

func BenchmarkEvaluate(b *testing.B) {
	position, _ := chess.NewPosition(chess.StartingFen)
	evaluator := NewEvaluator()

	for i := 0; i < b.N; i++ {
		evaluator.Evaluate(position)
	}
}
