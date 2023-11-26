package evaluation

import (
	"rosaline/internal/chess"
	"testing"
)

func evaluateBenchmark(b *testing.B, fen string) {
	position, _ := chess.NewPosition(fen)
	evaluator := NewEvaluator()

	for i := 0; i < b.N; i++ {
		evaluator.Evaluate(position)
	}
}

func BenchmarkEvaluate(b *testing.B) {
	cases := []string{
		chess.StartingFen,
		"7R/5pkp/4pN2/4P1P1/6K1/6P1/q1r5/7r w - - 1 46",
		"Bn2kbnr/p1p1pppp/3q4/8/3P4/2N3Pb/PPP2P1P/R1BQR1K1 b k - 0 13",
	}

	for _, c := range cases {
		b.Run(c, func(b *testing.B) {
			evaluateBenchmark(b, c)
		})
	}
}
