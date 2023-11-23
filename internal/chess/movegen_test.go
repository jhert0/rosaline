package chess

import "testing"

func BenchmarkGenerateLegalMoves(b *testing.B) {
	position, _ := NewPosition(StartingFen)

	for i := 0; i < b.N; i++ {
		position.GenerateMoves(LegalMoveGeneration)
	}
}

func BenchmarkGenerateCaptureMoves(b *testing.B) {
	position, _ := NewPosition(StartingFen)

	for i := 0; i < b.N; i++ {
		position.GenerateMoves(CaptureMoveGeneration)
	}
}
