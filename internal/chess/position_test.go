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
		"rnbqkb1r/ppp1pppp/5n2/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3",
		"r2qkb1r/ppp2ppp/2np1n2/8/2B3b1/5N2/PPPP1PPP/RNBQ1RK1 b kq - 5 6",
		"3r4/6k1/3n2p1/pNnB4/P4P2/1P6/7P/5K2 w - - 2 49",
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

func makeMoveTest(t *testing.T, position *Position, move string, fen string) {
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
		&position,
		"e2e4",
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
	)
}

func TestGetKingSquare(t *testing.T) {
	position, err := NewPosition(StartingFen)
	if err != nil {
		t.Fatalf("%s: fen %s returned error: %s", t.Name(), StartingFen, err)
	}

	whiteKing := position.GetKingSquare(White)
	if whiteKing != E1 {
		t.Fatalf("%s: expected white king to be on %d but got %d", t.Name(), E1, whiteKing)
	}

	blackKing := position.GetKingSquare(Black)
	if blackKing != E8 {
		t.Fatalf("%s: expected black king to be on %d but got %d", t.Name(), E8, blackKing)
	}
}

func TestIsValid(t *testing.T) {
	position, err := NewPosition(StartingFen)
	if err != nil {
		t.Fatalf("%s: fen %s returned error: %s", t.Name(), StartingFen, err)
	}

	if !position.IsValid() {
		t.Fatalf("%s: fen %s should be valid but IsValid returned false", t.Name(), StartingFen)
	}
}
