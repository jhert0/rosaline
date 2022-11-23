package chess

import "testing"

func getPieceTest(t *testing.T, algebraic string, expectedType PieceType, expectedColor Color) {
	position, _ := NewPosition(StartingFen)
	square, _ := SquareFromAlgebraic(algebraic)

	piece, err := position.GetPiece(square)
	if err != nil {
		t.Fatalf("%s: error returned trying to get piece from %s: %e", t.Name(), algebraic, err)
	}

	if piece.Type() != expectedType {
		t.Fatalf("%s: expected piece type: %d at %s got %d", t.Name(), expectedType, algebraic, piece.Type())
	}

	if piece.Color() != expectedColor {
		t.Fatalf("%s: expected piece color: %d at %s got %d", t.Name(), expectedColor, algebraic, piece.Color())
	}
}

func TestGetPiece(t *testing.T) {
	getPieceTest(t, "a1", Rook, White)
	getPieceTest(t, "a2", Pawn, White)
	getPieceTest(t, "a8", Rook, Black)
	getPieceTest(t, "e8", King, Black)
}
