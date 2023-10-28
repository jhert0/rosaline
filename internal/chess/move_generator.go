package chess

type MoveGenerationType uint8

const (
	LegalMoveGeneration MoveGenerationType = iota
	AttackMoveGeneration
)

// generatePawnMoves generates the moves for the pawns on the board
func generatePawnMoves(position Position, genType MoveGenerationType) []Move {
	moves := []Move{}
	dir := Square(pawnDirection(position.turn))

	pawnBB := position.pawnBB & position.GetColorBB(position.turn)
	for pawnBB > 0 {
		square := Square(pawnBB.PopLsb())

		if !position.PieceAt(square+dir) && genType != AttackMoveGeneration {
			toSquare := square + dir

			if toSquare.Rank() == pawnPromotionRank(position.Turn()) {
				for _, pieceType := range promotablePieces {
					move := NewMove(square, toSquare, NormalMove, PawnPushMoveFlag)
					move.WithPromotion(NewPiece(pieceType, position.turn))

					moves = append(moves, move)
				}
			} else {
				moves = append(moves, NewMove(square, toSquare, NormalMove, QuietMoveFlag))
			}

			if square.Rank() == pawnStartingRank(position.turn) {
				toSquare := square + (dir * 2)
				if !position.PieceAt(toSquare) {
					moves = append(moves, NewMove(square, toSquare, NormalMove, QuietMoveFlag))
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

			capturePiece, _ := position.GetPiece(captureSquare)
			if capturePiece != EmptyPiece && capturePiece.Color() != position.turn {
				if captureSquare.Rank() == pawnPromotionRank(position.Turn()) {
					for _, pieceType := range promotablePieces {
						move := NewMove(square, captureSquare, NormalMove, QuietMoveFlag)
						move.WithCapture(capturePiece)
						move.WithPromotion(NewPiece(pieceType, position.Turn()))
						moves = append(moves, move)
					}
				} else {
					move := NewMove(square, captureSquare, NormalMove, QuietMoveFlag)
					move.WithCapture(capturePiece)
					moves = append(moves, move)
				}
			}
		}

		if position.EnPassantPossible() {
			captureSquare := position.EnPassant() + Square(pawnDirection(position.turn.OpposingSide()))

			capturePiece, _ := position.GetPiece(captureSquare)
			if capturePiece.Type() == Pawn && capturePiece.Color() == position.Turn().OpposingSide() {
				move := NewMove(square, position.EnPassant(), EnPassantMove, QuietMoveFlag)
				move.WithCapture(capturePiece)
				moves = append(moves, move)
			}
		}
	}

	return moves
}

// generateKnightMoves generates the moves for the knights on the board
func generateKnightMoves(position Position) []Move {
	moves := []Move{}

	knightBB := position.knightBB & position.GetColorBB(position.turn)
	for knightBB > 0 {
		fromSquare := Square(knightBB.PopLsb())

		moveBB := knightMoves[fromSquare]
		for moveBB > 0 {
			toSquare := Square(moveBB.PopLsb())

			piece, _ := position.GetPiece(toSquare)
			if piece == EmptyPiece {
				moves = append(moves, NewMove(fromSquare, toSquare, NormalMove, QuietMoveFlag))
			} else if piece.Color() != position.turn {
				move := NewMove(fromSquare, toSquare, NormalMove, QuietMoveFlag)
				move.WithCapture(piece)
				moves = append(moves, move)
			}
		}
	}

	return moves
}

// generateBishopMoves generates the moves for the bishops on the board
func generateBishopMoves(position Position, pieceBB BitBoard) []Move {
	// TODO: generate moves using bitboard magic

	moves := []Move{}

	directions := []direction{
		north + east,
		north + west,
		south + east,
		south + west,
	}

	for pieceBB > 0 {
		fromSquare := Square(pieceBB.PopLsb())

	directionLoop:
		for _, direction := range directions {
			toSquare := fromSquare + Square(direction)

			for {
				if !toSquare.IsValid() {
					continue directionLoop
				}

				rankDifference, fileDifference := RankFileDifference(toSquare, fromSquare)
				if rankDifference != fileDifference {
					continue directionLoop
				}

				piece, _ := position.GetPiece(toSquare)

				// we hit one of our pieces, stop looking for moves in this direction
				if piece.Color() == position.turn {
					break
				}

				move := NewMove(fromSquare, toSquare, NormalMove, QuietMoveFlag)

				if piece.Color() == position.turn.OpposingSide() {
					move.WithCapture(piece)
				}

				moves = append(moves, move)

				if piece.Color() == position.turn.OpposingSide() {
					continue directionLoop
				}

				toSquare += Square(direction)
			}
		}
	}

	return moves
}

// generateRookMoves generates the moves for the rooks on the board
func generateRookMoves(position Position, pieceBB BitBoard) []Move {
	// TODO: use bitboard magic to generate rook moves

	moves := []Move{}

	directions := []direction{north, south, east, west}

	for pieceBB > 0 {
		fromSquare := Square(pieceBB.PopLsb())

	diretionLoop:
		for _, direction := range directions {
			toSquare := fromSquare + Square(direction)

			for {
				if !toSquare.IsValid() {
					continue diretionLoop
				}

				if (direction == north || direction == south) && toSquare.File() != fromSquare.File() {
					continue diretionLoop
				}

				if (direction == east || direction == west) && toSquare.Rank() != fromSquare.Rank() {
					continue diretionLoop
				}

				piece, _ := position.GetPiece(toSquare)

				// we hit one of our pieces, stop looking for moves in this direction
				if piece.Color() == position.turn {
					break
				}

				move := NewMove(fromSquare, toSquare, NormalMove, QuietMoveFlag)

				if piece.Color() == position.turn.OpposingSide() {
					move.WithCapture(piece)
				}

				moves = append(moves, move)

				if piece.Color() == position.turn.OpposingSide() {
					continue diretionLoop
				}

				toSquare += Square(direction)
			}
		}
	}

	return moves
}

// generateQueenMoves generates the moves for the queens on the board
func generateQueenMoves(position Position) []Move {
	queenBB := position.queenBB & position.GetColorBB(position.turn)

	moves := generateRookMoves(position, queenBB)

	bishopMoves := generateBishopMoves(position, queenBB)
	moves = append(moves, bishopMoves...)

	return moves
}

// generateKingMoves generates the moves for the kings on the board
func generateKingMoves(position Position, includeCastling bool) []Move {
	moves := []Move{}

	kingBB := position.kingBB & position.GetColorBB(position.turn)
	for kingBB > 0 {
		fromSquare := Square(kingBB.PopLsb())

		moveBB := kingMoves[fromSquare]
		for moveBB > 0 {
			toSquare := Square(moveBB.PopLsb())

			piece, _ := position.GetPiece(toSquare)
			if piece == EmptyPiece {
				moves = append(moves, NewMove(fromSquare, toSquare, NormalMove, QuietMoveFlag))
			} else if piece.Color() != position.Turn() {
				move := NewMove(fromSquare, toSquare, NormalMove, QuietMoveFlag)
				move.WithCapture(piece)
				moves = append(moves, move)
			}
		}
	}

	if includeCastling {
		kingSquare := position.GetKingSquare(position.turn)
		if position.turn == White {
			if position.HasCastlingRights(WhiteCastleKingside) && position.squaresEmpty([]Square{F1, G1}) {
				move := NewMove(kingSquare, kingSquare+Square(east*2), CastleMove, QuietMoveFlag)
				moves = append(moves, move)
			}

			if position.HasCastlingRights(WhiteCastleQueenside) && position.squaresEmpty([]Square{D1, C1, B1}) {
				move := NewMove(kingSquare, kingSquare+Square(west*2), CastleMove, QuietMoveFlag)
				moves = append(moves, move)
			}
		} else {
			if position.HasCastlingRights(BlackCastleKingside) && position.squaresEmpty([]Square{F8, G8}) {
				move := NewMove(kingSquare, kingSquare+Square(east*2), CastleMove, QuietMoveFlag)
				moves = append(moves, move)
			}

			if position.HasCastlingRights(BlackCastleQueenside) && position.squaresEmpty([]Square{D8, C8, B8}) {
				move := NewMove(kingSquare, kingSquare+Square(west*2), CastleMove, QuietMoveFlag)
				moves = append(moves, move)
			}
		}
	}

	return moves
}

// isLegalMove checks that the move would not result in an illegal position.
func (p Position) isLegalMove(move Move) bool {
	piece, _ := p.GetPiece(move.From)

	// the only way to get out of a double check is to move the king, therefore any other move is illegal
	if p.NumberOfCheckers(p.Turn()) == 2 && piece.Type() != King {
		return false
	}

	// check that the squares in between the king and rook are not attacked
	if move.Type() == CastleMove {
		direction := east
		if move.To == C1 || move.To == C8 {
			direction = west
		}

		difference := move.FileDifference()
		for i := 0; i <= difference; i++ {
			square := move.To + Square(difference*int(direction))
			if p.IsSquareAttackedBy(square, p.turn.OpposingSide()) {
				return false
			}
		}

		return true
	}

	// check that the king is not moving into an attacked square
	if piece.Type() == King {
		return !p.IsSquareAttackedBy(move.To, p.turn.OpposingSide())
	}

	return true
}

func (position Position) generateLegalMoves() []Move {
	moves := []Move{}

	colorBB := position.GetColorBB(position.turn)

	pawnMoves := generatePawnMoves(position, LegalMoveGeneration)
	moves = append(moves, pawnMoves...)

	knightMoves := generateKnightMoves(position)
	moves = append(moves, knightMoves...)

	bishopBB := position.GetPieceBB(Bishop)
	bishopMoves := generateBishopMoves(position, bishopBB&colorBB)
	moves = append(moves, bishopMoves...)

	rookBB := position.GetPieceBB(Rook)
	rookMoves := generateRookMoves(position, rookBB&colorBB)
	moves = append(moves, rookMoves...)

	queenMoves := generateQueenMoves(position)
	moves = append(moves, queenMoves...)

	kingMoves := generateKingMoves(position, !position.IsKingInCheck(position.turn))
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

	pawnMoves := generatePawnMoves(position, AttackMoveGeneration)
	moves = append(moves, pawnMoves...)

	knightMoves := generateKnightMoves(position)
	moves = append(moves, knightMoves...)

	bishopBB := position.GetPieceBB(Bishop)
	bishopMoves := generateBishopMoves(position, bishopBB&colorBB)
	moves = append(moves, bishopMoves...)

	rookBB := position.GetPieceBB(Rook)
	rookMoves := generateRookMoves(position, rookBB&colorBB)
	moves = append(moves, rookMoves...)

	queenMoves := generateQueenMoves(position)
	moves = append(moves, queenMoves...)

	kingMoves := generateKingMoves(position, false)
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
	case AttackMoveGeneration:
		return position.generateAttackMoves()
	default:
		panic("Unknown move generation type '%d' passed to GenerateMoves")
	}
}
