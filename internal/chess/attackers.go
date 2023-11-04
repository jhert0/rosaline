package chess

func generatePawnAttacks(position Position, color Color) []Move {
	attacks := []Move{}
	direction := pawnDirection(color)

	pawnBB := position.pawnBB & position.GetColorBB(color)
	for pawnBB > 0 {
		square := Square(pawnBB.PopLsb())

		if square.File() != 1 {
			move := NewMove(square, square+Square(direction)+Square(west), NormalMove, NoMoveFlag)
			attacks = append(attacks, move)
		}

		if square.File() != 8 {
			move := NewMove(square, square+Square(direction)+Square(east), NormalMove, NoMoveFlag)
			attacks = append(attacks, move)
		}
	}

	return attacks
}

func generateKnightAttacks(position Position, color Color) []Move {
	attacks := []Move{}

	knightBB := position.knightBB & position.GetColorBB(color)

	for knightBB > 0 {
		square := Square(knightBB.PopLsb())
		moveBB := knightMoves[square]

		for moveBB > 0 {
			toSquare := Square(moveBB.PopLsb())
			move := NewMove(square, toSquare, NormalMove, NoMoveFlag)
			attacks = append(attacks, move)
		}
	}

	return attacks
}

func generateBishopAttacks(pieceBB BitBoard, occupied BitBoard) []Move {
	attacks := []Move{}

	directions := []direction{
		north + east,
		north + west,
		south + east,
		south + west,
	}

	for pieceBB > 0 {
		square := Square(pieceBB.PopLsb())

	directionLoop:
		for _, direction := range directions {
			toSquare := square + Square(direction)

			for {
				if !toSquare.IsValid() {
					continue directionLoop
				}

				var rankDifference = RankDistance(square, toSquare)
				var fileDifference = FileDistance(square, toSquare)
				if rankDifference != fileDifference {
					continue directionLoop
				}

				move := NewMove(square, toSquare, NormalMove, NoMoveFlag)
				attacks = append(attacks, move)

				if occupied.IsBitSet(uint64(toSquare)) {
					continue directionLoop
				}

				toSquare += Square(direction)
			}
		}
	}

	return attacks
}

func generateRookAttacks(pieceBB BitBoard, occupied BitBoard) []Move {
	// TODO: use bit board magic to generate this

	attacks := []Move{}

	directions := []direction{north, south, east, west}

	for pieceBB > 0 {
		square := Square(pieceBB.PopLsb())

	diretionLoop:
		for _, direction := range directions {
			toSquare := square + Square(direction)

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

				move := NewMove(square, toSquare, NormalMove, NoMoveFlag)
				attacks = append(attacks, move)

				if occupied.IsBitSet(uint64(toSquare)) {
					continue diretionLoop
				}

				toSquare += Square(direction)
			}
		}
	}

	return attacks
}

func generateQueenAttacks(position Position, occupied BitBoard, color Color) []Move {
	attacks := []Move{}

	queenBB := position.queenBB & position.GetColorBB(color)

	rookAttacks := generateRookAttacks(queenBB, occupied)
	attacks = append(attacks, rookAttacks...)
	bishopAttacks := generateBishopAttacks(queenBB, occupied)
	attacks = append(attacks, bishopAttacks...)

	return attacks
}

func generateKingAttacks(position Position, color Color) []Move {
	attacks := []Move{}

	kingBB := position.kingBB & position.GetColorBB(color)
	for kingBB > 0 {
		fromSquare := Square(kingBB.PopLsb())

		moveBB := kingMoves[fromSquare]
		for moveBB > 0 {
			toSquare := Square(moveBB.PopLsb())
			move := NewMove(fromSquare, toSquare, NormalMove, NoMoveFlag)
			attacks = append(attacks, move)
		}
	}

	return attacks
}

func getAttackers(position Position, color Color) []Move {
	attacks := []Move{}

	occupied := position.whiteBB | position.blackBB
	colorBB := position.GetColorBB(color)

	pawnAttacks := generatePawnAttacks(position, color)
	attacks = append(attacks, pawnAttacks...)

	knightAttacks := generateKnightAttacks(position, color)
	attacks = append(attacks, knightAttacks...)

	bishopBB := position.bishopBB & colorBB
	bishopAttacks := generateBishopAttacks(bishopBB, occupied)
	attacks = append(attacks, bishopAttacks...)

	rookBB := position.rookBB & colorBB
	rookAttacks := generateRookAttacks(rookBB, occupied)
	attacks = append(attacks, rookAttacks...)

	queenAttacks := generateQueenAttacks(position, occupied, color)
	attacks = append(attacks, queenAttacks...)

	kingAttacks := generateKingAttacks(position, color)
	attacks = append(attacks, kingAttacks...)

	return attacks
}
