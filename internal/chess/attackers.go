package chess

func generatePawnAttacks(position Position, color Color) []Move {
	attacks := []Move{}
	direction := pawnDirection(color)

	pawnBB := position.pawnBB & position.GetColorBB(color)
	for pawnBB > 0 {
		square := Square(pawnBB.PopLsb())

		if square.File() != 1 {
			move := NewMove(square, square+Square(direction)+Square(west), CaptureMove)
			attacks = append(attacks, move)
		}

		if square.File() != 8 {
			move := NewMove(square, square+Square(direction)+Square(east), CaptureMove)
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
			move := NewMove(square, toSquare, CaptureMove)
			attacks = append(attacks, move)
		}
	}

	return attacks
}

func generateBishopAttacks(pieceBB BitBoard, occupied BitBoard) []Move {
	attacks := []Move{}

	for pieceBB > 0 {
		fromSquare := Square(pieceBB.PopLsb())

		attackBB := getBishopAttacks(occupied, fromSquare)
		for attackBB > 0 {
			toSquare := Square(attackBB.PopLsb())
			move := NewMove(fromSquare, toSquare, CaptureMove)
			attacks = append(attacks, move)
		}
	}

	return attacks
}

func generateRookAttacks(pieceBB BitBoard, occupied BitBoard) []Move {
	attacks := []Move{}

	for pieceBB > 0 {
		fromSquare := Square(pieceBB.PopLsb())

		attackBB := getRookAttacks(occupied, fromSquare)
		for attackBB > 0 {
			toSquare := Square(attackBB.PopLsb())
			move := NewMove(fromSquare, toSquare, CaptureMove)
			attacks = append(attacks, move)
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
			move := NewMove(fromSquare, toSquare, CaptureMove)
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
