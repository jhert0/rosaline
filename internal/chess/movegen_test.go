package chess

import (
	"fmt"
	"testing"
)

func generateLegalMovesBenchmark(b *testing.B, fen string) {
	position, _ := NewPosition(fen)
	for i := 0; i < b.N; i++ {
		position.GenerateMoves(LegalMoveGeneration)
	}
}

func BenchmarkGenerateLegalMoves(b *testing.B) {
	cases := []struct {
		Name string
		Fen  string
	}{
		{
			Name: "StartingPosition",
			Fen:  StartingFen,
		},
		{
			Name: "LargestNumberOfLegalMoves",
			Fen:  "R6R/3Q4/1Q4Q1/4Q3/2Q4Q/Q4Q2/pp1Q4/kBNNK1B1 w - - 0 1",
		},
		{
			Name: "InCheck",
			Fen:  "rn2kbnr/ppp2ppp/3pb3/4p3/2B1q3/BPN5/P1PP1PPP/R2QK1NR w KQkq - 0 6",
		},
	}

	for _, c := range cases {
		b.Run(fmt.Sprintf("%s", c.Name), func(b *testing.B) {
			generateLegalMovesBenchmark(b, c.Fen)
		})
	}
}

func BenchmarkGenerateCaptureMoves(b *testing.B) {
	position, _ := NewPosition(StartingFen)

	for i := 0; i < b.N; i++ {
		position.GenerateMoves(CaptureMoveGeneration)
	}
}
