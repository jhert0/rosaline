package chess

import "math"

func generatePawnAttacks(square Square, color Color) BitBoard {
	attacks := BitBoard(0)
	direction := pawnDirection(color)

	attacks.SetBit(uint64(square + Square(direction) + Square(east)))
	attacks.SetBit(uint64(square + Square(direction) + Square(west)))

	return attacks
}

func generateKnightAttacks(square Square) BitBoard {
	return knightMoves[square]
}

func generateBishopAttacks(square Square, occupied BitBoard) BitBoard {
	attacks := BitBoard(0)

	directions := []direction{
		north + east,
		north + west,
		south + east,
		south + west,
	}

directionLoop:
	for _, direction := range directions {
		toSquare := square

		for {
			if !toSquare.IsValid() {
				continue directionLoop
			}

			var rankDifference = math.Abs(float64(toSquare.Rank()) - float64(square.Rank()))
			var fileDifference = math.Abs(float64(toSquare.File()) - float64(square.File()))
			if rankDifference != fileDifference {
				continue directionLoop
			}

			attacks.SetBit(uint64(toSquare))

			if occupied.BitSet(uint64(toSquare)) {
				continue directionLoop
			}

			toSquare += Square(direction)
		}
	}

	return attacks
}

func generateRookAttacks(square Square, occupied BitBoard) BitBoard {
	// TODO: use bit board magic to generate this

	attacks := BitBoard(0)

	directions := []direction{north, south, east, west}

diretionLoop:
	for _, direction := range directions {
		toSquare := square

		for {
			if !toSquare.IsValid() {
				continue diretionLoop
			}

			if (direction == north || direction == south) && toSquare.File() != square.File() {
				continue diretionLoop
			}

			if (direction == east || direction == west) && toSquare.Rank() != square.Rank() {
				continue diretionLoop
			}

			attacks.SetBit(uint64(toSquare))

			if occupied.BitSet(uint64(toSquare)) {
				continue diretionLoop
			}

			toSquare += Square(direction)
		}
	}

	return attacks
}

func generateQueenAttacks(square Square, occupied BitBoard) BitBoard {
	return generateRookAttacks(square, occupied) | generateBishopAttacks(square, occupied)
}

func generateKingAttacks(square Square) BitBoard {
	return kingMoves[square]
}

func getAttackers(position Position, color Color) BitBoard {
	attackers := BitBoard(0)
	occupied := position.whiteBB | position.blackBB

	colorBB := position.GetColorBB(color)
	for colorBB > 0 {
		square := Square(colorBB.TrailingZeros())

		piece, _ := position.GetPiece(square)
		switch piece.Type() {
		case Pawn:
			attackers = generatePawnAttacks(square, piece.Color())
			break
		case Knight:
			attackers = generateKnightAttacks(square)
			break
		case Bishop:
			attackers = generateBishopAttacks(square, occupied)
			break
		case Rook:
			attackers = generateRookAttacks(square, occupied)
			break
		case Queen:
			attackers = generateQueenAttacks(square, occupied)
			break
		case King:
			attackers = generateKingAttacks(square)
			break
		}

		colorBB.ClearBit(uint64(square))
	}

	return attackers
}
