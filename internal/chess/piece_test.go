package chess

import "testing"

func pieceColorTest(t *testing.T, piece Piece, expectedColor Color) {
	color := piece.Color()
	if color != expectedColor {
		t.Fatalf("%s: expected piece color to be '%s' got '%s'", t.Name(), expectedColor, color)
	}
}

func TestPieceColor(t *testing.T) {
	pieceColorTest(t, NewPiece(Pawn, Black), Black)
	pieceColorTest(t, NewPiece(Knight, Black), Black)
	pieceColorTest(t, NewPiece(King, Black), Black)

	pieceColorTest(t, NewPiece(Pawn, White), White)
	pieceColorTest(t, NewPiece(Knight, White), White)
	pieceColorTest(t, NewPiece(King, White), White)
}

func pieceTypeTest(t *testing.T, piece Piece, expectedType PieceType) {
	pieceType := piece.Type()
	if pieceType != expectedType {
		t.Fatalf("%s: expected piece type to be '%s' got '%s'", t.Name(), expectedType, pieceType)
	}
}

func TestPieceType(t *testing.T) {
	pieceTypeTest(t, NewPiece(Pawn, Black), Pawn)
	pieceTypeTest(t, NewPiece(Knight, Black), Knight)
	pieceTypeTest(t, NewPiece(King, Black), King)

	pieceTypeTest(t, NewPiece(Pawn, White), Pawn)
	pieceTypeTest(t, NewPiece(Knight, White), Knight)
	pieceTypeTest(t, NewPiece(King, White), King)
}
