package chess

type MoveGenerationType uint8

const (
	LegalMoveGeneration MoveGenerationType = iota
	CaptureMoveGeneration
)

// generatePawnMoves generates the moves for the pawns on the board
func generatePawnMoves(position Position, genType MoveGenerationType) []Move {
	moves := []Move{}
	dir := Square(pawnDirection(position.turn))

	pawnBB := position.pawnBB & position.GetColorBB(position.turn)
	for pawnBB > 0 {
		square := Square(pawnBB.PopLsb())

		if !position.IsSquareOccupied(square+dir) && genType != CaptureMoveGeneration {
			toSquare := square + dir

			if toSquare.Rank() == pawnPromotionRank(position.Turn()) {
				for _, pieceType := range promotablePieces {
					move := NewMove(square, toSquare, QuietMove)
					move.WithFlags(PawnPushMoveFlag)
					move.WithPromotion(NewPiece(pieceType, position.turn))

					moves = append(moves, move)
				}
			} else {
				move := NewMove(square, toSquare, QuietMove)
				move.WithFlags(PawnPushMoveFlag)
				moves = append(moves, move)
			}

			if square.Rank() == pawnStartingRank(position.turn) {
				toSquare := square + (dir * 2)
				if !position.IsSquareOccupied(toSquare) {
					move := NewMove(square, toSquare, QuietMove)
					move.WithFlags(PawnPushMoveFlag)
					moves = append(moves, move)
				}
			}
		}

		captureOffsets := [2]direction{east, west}
		for _, offset := range captureOffsets {
			captureSquare := square + dir + Square(offset)

			if offset == west && square.File() == 1 {
				continue
			}

			if offset == east && square.File() == 8 {
				continue
			}

			capturePiece, _ := position.GetPieceAt(captureSquare)
			if capturePiece != EmptyPiece && capturePiece.Color() != position.turn {
				if captureSquare.Rank() == pawnPromotionRank(position.Turn()) {
					for _, pieceType := range promotablePieces {
						move := NewMove(square, captureSquare, CaptureMove)
						move.WithPromotion(NewPiece(pieceType, position.Turn()))
						moves = append(moves, move)
					}
				} else {
					move := NewMove(square, captureSquare, CaptureMove)
					moves = append(moves, move)
				}
			}
		}

		enPassantSquare := position.EnPassant()
		if square+dir+Square(east) == enPassantSquare || square+dir+Square(west) == enPassantSquare {
			captureSquare := enPassantSquare + Square(pawnDirection(position.turn.OpposingSide()))

			capturePiece, _ := position.GetPieceAt(captureSquare)
			if capturePiece.Type() == Pawn && capturePiece.Color() == position.Turn().OpposingSide() {
				move := NewMove(square, position.EnPassant(), EnPassantMove)
				moves = append(moves, move)
			}
		}
	}

	return moves
}

// generateKnightMoves generates the moves for the knights on the board
func generateKnightMoves(position Position, genType MoveGenerationType) []Move {
	moves := []Move{}

	knightBB := position.knightBB & position.GetColorBB(position.turn)

	occupied := position.whiteBB | position.blackBB
	opponent := position.GetColorBB(position.turn.OpposingSide())

	for knightBB > 0 {
		fromSquare := Square(knightBB.PopLsb())
		attackBB := knightMoves[fromSquare]

		if genType != CaptureMoveGeneration {
			moveBB := attackBB & ^occupied
			for moveBB > 0 {
				toSquare := Square(moveBB.PopLsb())
				move := NewMove(fromSquare, toSquare, QuietMove)
				moves = append(moves, move)
			}
		}

		capturesBB := attackBB & opponent
		for capturesBB > 0 {
			toSquare := Square(capturesBB.PopLsb())
			move := NewMove(fromSquare, toSquare, CaptureMove)
			moves = append(moves, move)
		}
	}

	return moves
}

// generateBishopMoves generates the moves for the bishops on the board
func generateBishopMoves(position Position, pieceBB BitBoard, genType MoveGenerationType) []Move {
	moves := []Move{}

	occupied := position.whiteBB | position.blackBB
	opponent := position.GetColorBB(position.turn.OpposingSide())

	for pieceBB > 0 {
		fromSquare := Square(pieceBB.PopLsb())
		attackBB := getBishopAttacks(occupied, fromSquare)

		if genType != CaptureMoveGeneration {
			moveBB := attackBB & ^occupied
			for moveBB > 0 {
				toSquare := Square(moveBB.PopLsb())
				move := NewMove(fromSquare, toSquare, QuietMove)
				moves = append(moves, move)
			}
		}

		capturesBB := attackBB & opponent
		for capturesBB > 0 {
			toSquare := Square(capturesBB.PopLsb())
			move := NewMove(fromSquare, toSquare, CaptureMove)
			moves = append(moves, move)
		}
	}

	return moves
}

// generateRookMoves generates the moves for the rooks on the board
func generateRookMoves(position Position, pieceBB BitBoard, genType MoveGenerationType) []Move {
	moves := []Move{}

	occupied := position.whiteBB | position.blackBB
	opponent := position.GetColorBB(position.turn.OpposingSide())

	for pieceBB > 0 {
		fromSquare := Square(pieceBB.PopLsb())
		attacks := getRookAttacks(occupied, fromSquare)

		if genType != CaptureMoveGeneration {
			moveBB := attacks & ^occupied
			for moveBB > 0 {
				toSquare := Square(moveBB.PopLsb())
				move := NewMove(fromSquare, toSquare, QuietMove)
				moves = append(moves, move)
			}
		}

		capturesBB := attacks & opponent
		for capturesBB > 0 {
			toSquare := Square(capturesBB.PopLsb())
			move := NewMove(fromSquare, toSquare, CaptureMove)
			moves = append(moves, move)
		}
	}

	return moves
}

// generateQueenMoves generates the moves for the queens on the board
func generateQueenMoves(position Position, genType MoveGenerationType) []Move {
	queenBB := position.queenBB & position.GetColorBB(position.turn)

	moves := generateRookMoves(position, queenBB, genType)

	bishopMoves := generateBishopMoves(position, queenBB, genType)
	moves = append(moves, bishopMoves...)

	return moves
}

// generateKingMoves generates the moves for the kings on the board
func generateKingMoves(position Position, genType MoveGenerationType, includeCastling bool) []Move {
	moves := []Move{}

	kingSquare := position.GetKingSquare(position.turn)

	attacks := kingMoves[kingSquare]
	occupied := position.whiteBB | position.blackBB
	opponent := position.GetColorBB(position.turn.OpposingSide())

	if genType != CaptureMoveGeneration {
		moveBB := attacks & ^occupied
		for moveBB > 0 {
			toSquare := Square(moveBB.PopLsb())
			move := NewMove(kingSquare, toSquare, QuietMove)
			moves = append(moves, move)
		}
	}

	capturesBB := attacks & opponent
	for capturesBB > 0 {
		toSquare := Square(capturesBB.PopLsb())
		move := NewMove(kingSquare, toSquare, CaptureMove)
		moves = append(moves, move)
	}

	if includeCastling {
		if position.turn == White {
			if position.HasCastlingRights(WhiteCastleKingside) && position.squaresEmpty([]Square{F1, G1}) {
				move := NewMove(kingSquare, kingSquare+Square(east*2), CastleMove)
				moves = append(moves, move)
			}

			if position.HasCastlingRights(WhiteCastleQueenside) && position.squaresEmpty([]Square{D1, C1, B1}) {
				move := NewMove(kingSquare, kingSquare+Square(west*2), CastleMove)
				moves = append(moves, move)
			}
		} else {
			if position.HasCastlingRights(BlackCastleKingside) && position.squaresEmpty([]Square{F8, G8}) {
				move := NewMove(kingSquare, kingSquare+Square(east*2), CastleMove)
				moves = append(moves, move)
			}

			if position.HasCastlingRights(BlackCastleQueenside) && position.squaresEmpty([]Square{D8, C8, B8}) {
				move := NewMove(kingSquare, kingSquare+Square(west*2), CastleMove)
				moves = append(moves, move)
			}
		}
	}

	return moves
}

// isLegalMove checks that the move would not result in an illegal position.
func (p Position) isLegalMove(move Move) bool {
	// check that the squares in between the king and rook are not attacked
	if move.Type() == CastleMove {
		direction := east
		if move.To() == C1 || move.To() == C8 {
			direction = west
		}

		difference := move.FileDifference()
		for i := 0; i <= difference; i++ {
			square := move.From() + Square(difference*int(direction))
			if p.IsSquareAttackedBy(square, p.turn.OpposingSide()) {
				return false
			}
		}

		return true
	}

	// check that after the move is made that the king is not in check
	p.MakeMove(move)
	inCheck := p.IsKingInCheck(p.turn.OpposingSide())
	p.Undo()

	return !inCheck
}

func (position Position) generateLegalMoves() []Move {
	moves := []Move{}

	checkers := position.NumberOfCheckers(position.turn)
	if checkers < 2 {
		colorBB := position.GetColorBB(position.turn)

		pawnMoves := generatePawnMoves(position, LegalMoveGeneration)
		moves = append(moves, pawnMoves...)

		knightMoves := generateKnightMoves(position, LegalMoveGeneration)
		moves = append(moves, knightMoves...)

		bishopBB := position.GetPieceBB(Bishop)
		bishopMoves := generateBishopMoves(position, bishopBB&colorBB, LegalMoveGeneration)
		moves = append(moves, bishopMoves...)

		rookBB := position.GetPieceBB(Rook)
		rookMoves := generateRookMoves(position, rookBB&colorBB, LegalMoveGeneration)
		moves = append(moves, rookMoves...)

		queenMoves := generateQueenMoves(position, LegalMoveGeneration)
		moves = append(moves, queenMoves...)
	}

	inCheck := checkers != 0
	kingMoves := generateKingMoves(position, LegalMoveGeneration, !inCheck)
	moves = append(moves, kingMoves...)

	legalMoves := []Move{}
	for _, move := range moves {
		if position.isLegalMove(move) {
			legalMoves = append(legalMoves, move)
		}
	}

	return legalMoves
}

func (position Position) generateAttackMoves() []Move {
	moves := []Move{}

	colorBB := position.GetColorBB(position.turn)

	pawnMoves := generatePawnMoves(position, CaptureMoveGeneration)
	moves = append(moves, pawnMoves...)

	knightMoves := generateKnightMoves(position, CaptureMoveGeneration)
	moves = append(moves, knightMoves...)

	bishopBB := position.GetPieceBB(Bishop)
	bishopMoves := generateBishopMoves(position, bishopBB&colorBB, CaptureMoveGeneration)
	moves = append(moves, bishopMoves...)

	rookBB := position.GetPieceBB(Rook)
	rookMoves := generateRookMoves(position, rookBB&colorBB, CaptureMoveGeneration)
	moves = append(moves, rookMoves...)

	queenMoves := generateQueenMoves(position, CaptureMoveGeneration)
	moves = append(moves, queenMoves...)

	kingMoves := generateKingMoves(position, CaptureMoveGeneration, false)
	moves = append(moves, kingMoves...)

	legalMoves := []Move{}
	for _, move := range moves {
		if position.isLegalMove(move) {
			legalMoves = append(legalMoves, move)
		}
	}

	return legalMoves
}

// GenerateMoves generates all legal moves in the position.
func (position Position) GenerateMoves(genType MoveGenerationType) []Move {
	switch genType {
	case LegalMoveGeneration:
		return position.generateLegalMoves()
	case CaptureMoveGeneration:
		return position.generateAttackMoves()
	default:
		panic("Unknown move generation type '%d' passed to GenerateMoves")
	}
}
