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

func rankFileTest(t *testing.T, rank int, file int, expectedSquare Square) {
	square := SquareFromRankFile(rank, file)
	if square != Square(expectedSquare) {
		t.Fatalf("%s: expected %d  but got %d ", t.Name(), expectedSquare, square)
	}
}

func TestSquareFromRankFile(t *testing.T) {
	rankFileTest(t, 1, 1, A1)
	rankFileTest(t, 1, 8, H1)
	rankFileTest(t, 8, 1, A8)
}

func squareFileTest(t *testing.T, square Square, file int) {
	if square.File() != file {
		t.Fatalf("%s: expected %d but got %d for %s", t.Name(), file, square.File(), square.ToAlgebraic())
	}
}

func TestSquareFile(t *testing.T) {
	squareFileTest(t, A1, 1)
	squareFileTest(t, A2, 1)
	squareFileTest(t, G5, 7)
}

func squareRankTest(t *testing.T, square Square, rank int) {
	if square.Rank() != rank {
		t.Fatalf("%s: expected %d but got %d for %s", t.Name(), rank, square.Rank(), square.ToAlgebraic())
	}
}

func TestSquareRank(t *testing.T) {
	squareRankTest(t, A1, 1)
	squareRankTest(t, A2, 2)
	squareRankTest(t, G5, 5)
}
