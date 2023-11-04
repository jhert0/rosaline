package chess

import "math/rand"

const (
	numSquares    = 64
	numPieceTypes = 6
	numSides      = 2
)

var zobristTable [numSquares][numPieceTypes][numSides]uint64

func init() {
	for i := 0; i < numSquares; i++ {
		for j := 0; j < numPieceTypes; j++ {
			for k := 0; k < numSides; k++ {
				zobristTable[i][j][k] = rand.Uint64()
			}
		}
	}
}

func generateHash(p Position) uint64 {
	var hash uint64 = 0

	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			square := SquareFromRankFile(rank+1, file+1)
			piece, err := p.GetPieceAt(square)
			if err == nil {
				var pieceIndex int
				switch piece.Type() {
				case Pawn:
					pieceIndex = 0
					break
				case Knight:
					pieceIndex = 1
					break
				case Bishop:
					pieceIndex = 2
					break
				case Rook:
					pieceIndex = 3
					break
				case Queen:
					pieceIndex = 4
					break
				case King:
					pieceIndex = 5
					break
				}

				var colorIndex = 0
				if piece.Color() == Black {
					colorIndex = 1
				}

				hash ^= zobristTable[square][pieceIndex][colorIndex]
			}
		}
	}

	return hash
}
