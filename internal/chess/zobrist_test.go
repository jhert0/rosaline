package chess

import "testing"

func generateHashBenchmark(b *testing.B, fen string) {
	position, _ := NewPosition(fen)
	for i := 0; i < b.N; i++ {
		generateHash(position)
	}
}

func BenchmarkGenerateHash(b *testing.B) {
	cases := []string{
		StartingFen,
		"1b2k2r/Pppp1ppp/5nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w k - 0 1",
	}

	for _, c := range cases {
		b.Run(c, func(b *testing.B) {
			generateHashBenchmark(b, c)
		})
	}
}
