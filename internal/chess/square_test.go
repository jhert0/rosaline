package chess

import "testing"

func algebraicTest(t *testing.T, algebraic string, expectedRank int, expectedFile int) {
	square, err := SquareFromAlgebraic(algebraic)
	if err != nil {
		t.Fatalf("%s: error: %e returned", t.Name(), err)
	}

	if square.Rank() != expectedRank || square.File() != expectedFile {
		t.Fatalf("%s: expected rank: %d and file: %d but got rank: %d and file: %d", t.Name(), expectedRank, expectedFile, square.Rank(), square.File())
	}
}

func TestSquareFromAlgebraic(t *testing.T) {
	algebraicTest(t, "a1", 1, 1)
	algebraicTest(t, "h1", 1, 8)
	algebraicTest(t, "a8", 8, 1)
}

func rankFileTest(t *testing.T, rank int, file int, expectedSquare int) {
	square := SquareFromRankFile(rank, file)
	if square != Square(expectedSquare) {
		t.Fatalf("%s: expected %d  but got %d ", t.Name(), expectedSquare, square)
	}
}

func TestSquareFromRankFile(t *testing.T) {
	rankFileTest(t, 1, 1, 0)
	rankFileTest(t, 1, 8, 7)
	rankFileTest(t, 8, 1, 56)
}
