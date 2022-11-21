package chess

import "testing"

func getPieceTest(t *testing.T, algebraic string, expectedPiece Piece) {
	position, _ := NewPosition(STARTING_FEN)
	square, _ := SquareFromAlgebraic(algebraic)

	piece, err := position.GetPiece(square)
	if err != nil {
		t.Fatalf("%s: error returned trying to get piece from %s: %e", t.Name(), algebraic, err)
	}

	if piece != expectedPiece {
		t.Fatalf("%s: expected piece: %d got %d", t.Name(), expectedPiece, piece)
	}
}

func TestGetPiece(t *testing.T) {
	getPieceTest(t, "a1", WhiteRook)
	getPieceTest(t, "a2", WhitePawn)
	getPieceTest(t, "a8", BlackRook)
	getPieceTest(t, "e8", BlackKing)
}
