package chess

import (
	"slices"
	"testing"
)

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

func sliceEqual(s1, s2 []Square) bool {
	if len(s1) != len(s2) {
		return false
	}

	slices.Sort(s1)
	slices.Sort(s2)

	for i, val := range s1 {
		if val != s2[i] {
			return false
		}
	}

	return true
}

func surroundingSquareTest(t *testing.T, square Square, expectedSquares []Square) {
	squares := SurroundingSquares(square)
	if !sliceEqual(squares, expectedSquares) {
		t.Fatalf("%s: expected squares %v got %v", t.Name(), expectedSquares, squares)
	}
}

func TestSurroundingSquares(t *testing.T) {
	surroundingSquareTest(t, A1, []Square{A2, B2, B1})
	surroundingSquareTest(t, C5, []Square{C4, D4, D5, D6, C6, B6, B5, B4})
}

func fileDistanceTest(t *testing.T, sq1, sq2 Square, expectedDistance int) {
	distance := FileDistance(sq1, sq2)
	if distance != expectedDistance {
		t.Fatalf("%s: expected file distance betwen %s and %s to be '%d' got '%d'", t.Name(), sq1, sq1, expectedDistance, distance)
	}
}

func TestFileDistance(t *testing.T) {
	fileDistanceTest(t, A1, A2, 0)
	fileDistanceTest(t, A1, B2, 1)
	fileDistanceTest(t, C3, G8, 4)
}

func rankDistanceTest(t *testing.T, sq1, sq2 Square, expectedDistance int) {
	distance := RankDistance(sq1, sq2)
	if distance != expectedDistance {
		t.Fatalf("%s: expected rank distance between %s and %s to be '%d' got '%d'", t.Name(), sq1, sq2, expectedDistance, distance)
	}
}

func TestRankDistance(t *testing.T) {
	rankDistanceTest(t, B3, C3, 0)
	rankDistanceTest(t, A1, A2, 1)
	rankDistanceTest(t, A1, A8, 7)
}

func colorTest(t *testing.T, square Square, expectedColor Color) {
	color := square.Color()
	if color != expectedColor {
		t.Fatalf("%s: expected square color for %s to be '%s' got '%s'", t.Name(), square, expectedColor, color)
	}
}

func TestColor(t *testing.T) {
	colorTest(t, A1, Black)
	colorTest(t, H8, Black)
	colorTest(t, D6, Black)
	colorTest(t, C7, Black)
	colorTest(t, A7, Black)

	colorTest(t, A2, White)
	colorTest(t, E4, White)
	colorTest(t, A8, White)
	colorTest(t, D7, White)
	colorTest(t, H1, White)
}
