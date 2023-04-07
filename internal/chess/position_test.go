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

func TestFen(t *testing.T) {
	fens := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"rnbq1rk1/pppp1ppp/4pn2/8/1bB1P3/2N2N2/PPPP1PPP/R1BQK2R w KQ - 6 5",
		"rn1qkbnr/pp2pppp/8/7b/2BQ4/2N1PP2/PP4PP/R1B1K1NR b KQkq - 0 7",
	}

	for _, fen := range fens {
		position, err := NewPosition(fen)
		if err != nil {
			t.Fatalf("%s: fen %s returned error: %e", t.Name(), fen, err)
		}

		if position.Fen() != fen {
			t.Fatalf("%s: expected %s got %s", t.Name(), fen, position.Fen())
		}
	}
}

func makeMoveTest(t *testing.T, position Position, move string, fen string) {
	err := position.MakeUciMove(move)
	if err != nil {
		t.Fatalf("%s: move: %s returned an error: %e", t.Name(), move, err)
	}

	if position.Fen() != fen {
		t.Fatalf("%s: after move: %s expecteed fen %s got %s", t.Name(), move, fen, position.Fen())
	}
}

func TestMakeUciMove(t *testing.T) {
	position, err := NewPosition(StartingFen)
	if err != nil {
		t.Fatalf("%s: fen %s returned error: %s", t.Name(), StartingFen, err)
	}

	makeMoveTest(
		t,
		position,
		"e2e4",
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
	)
}
