package evaluation

import (
	"fmt"
	"rosaline/internal/chess"
)

// pieceValue returns the value of a piece in centipawns.
func pieceValue(piece chess.Piece) int {
	switch piece.Type() {
	case chess.Pawn:
		return 100
	case chess.Knight:
		return 320
	case chess.Bishop:
		return 330
	case chess.Rook:
		return 500
	case chess.Queen:
		return 900
	case chess.King:
		return 20000
	}

	panic(fmt.Sprintf("requested piece value for unknown piece type: %s", piece.Type()))
}
