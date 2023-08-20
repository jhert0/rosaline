package chess

import (
	"math"
)

// generatePawnMoves generates the moves for the pawns on the board
func generatePawnMoves(position Position, pieceBB BitBoard) []Move {
	moves := []Move{}

	for pieceBB > 0 {
		square := Square(pieceBB.TrailingZeros())

		direction := Square(pawnDirection(position.turn))
		if !position.PieceAt(square + direction) {
			toSquare := square + direction

			if toSquare.Rank() == 8 && position.Turn() == White {
				for _, pieceType := range promotablePieces {
					move := NewMove(square, toSquare, NormalMove, PawnPushMoveFlag)
					move.WithPromotion(NewPiece(pieceType, position.turn))

					moves = append(moves, move)
				}
			} else if toSquare.Rank() == 1 && position.Turn() == Black {
				for _, pieceType := range promotablePieces {
					move := NewMove(square, toSquare, NormalMove, QuietMoveFlag)
					move.WithPromotion(NewPiece(pieceType, position.turn))

					moves = append(moves, move)
				}
			} else {
				moves = append(moves, NewMove(square, toSquare, NormalMove, QuietMoveFlag))
			}

			if !position.PieceAt(square+(direction*2)) && square.Rank() == pawnStartingRank(position.turn) {
				moves = append(moves, NewMove(square, square+(direction*2), NormalMove, QuietMoveFlag))
			}
		}

		capturePiece, _ := position.GetPiece(square + direction + Square(east))
		if capturePiece != EmptyPiece && capturePiece.Color() != position.turn {
			move := NewMove(square, square+direction+Square(east), NormalMove, QuietMoveFlag)
			move.WithCapture(capturePiece)
			moves = append(moves, move)
		}

		capturePiece, _ = position.GetPiece(square + direction + Square(west))
		if capturePiece != EmptyPiece && capturePiece.Color() != position.turn {
			move := NewMove(square, square+direction+Square(west), NormalMove, QuietMoveFlag)
			move.WithCapture(capturePiece)
			moves = append(moves, move)
		}

		if position.EnPassantPossible() {
			if position.IsPieceAt(square+Square(west), Pawn, position.turn.OpposingSide()) {
				move := NewMove(square, square+Square(west)+direction, EnPassantMove, QuietMoveFlag)
				moves = append(moves, move)
			}

			if position.IsPieceAt(square+Square(east), Pawn, position.turn.OpposingSide()) {
				move := NewMove(square, square+Square(east)+direction, EnPassantMove, QuietMoveFlag)
				moves = append(moves, move)
			}
		}

		pieceBB.ClearBit(uint64(square))
	}

	return moves
}

// generateKnightMoves generates the moves for the knights on the board
func generateKnightMoves(position Position, pieceBB BitBoard) []Move {
	moves := []Move{}

	for pieceBB > 0 {
		fromSquare := Square(pieceBB.TrailingZeros())

		moveBB := knightMoves[fromSquare]
		for moveBB > 0 {
			toSquare := Square(moveBB.TrailingZeros())

			piece, _ := position.GetPiece(toSquare)
			if piece == EmptyPiece {
				moves = append(moves, NewMove(fromSquare, toSquare, NormalMove, QuietMoveFlag))
			} else if piece.Color() != position.turn {
				move := NewMove(fromSquare, toSquare, NormalMove, QuietMoveFlag)
				move.WithCapture(piece)
				moves = append(moves, move)
			}

			moveBB.ClearBit(uint64(toSquare))
		}

		pieceBB.ClearBit(uint64(fromSquare))
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
		fromSquare := Square(pieceBB.TrailingZeros())

	directionLoop:
		for _, direction := range directions {
			toSquare := fromSquare + Square(direction)

			for {
				if !toSquare.IsValid() {
					continue directionLoop
				}

				var rankDifference = math.Abs(float64(toSquare.Rank()) - float64(fromSquare.Rank()))
				var fileDifference = math.Abs(float64(toSquare.File()) - float64(fromSquare.File()))
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

		pieceBB.ClearBit(uint64(fromSquare))
	}

	return moves
}

// generateRookMoves generates the moves for the rooks on the board
func generateRookMoves(position Position, pieceBB BitBoard) []Move {
	// TODO: use bitboard magic to generate rook moves

	moves := []Move{}

	directions := []direction{north, south, east, west}

	for pieceBB > 0 {
		fromSquare := Square(pieceBB.TrailingZeros())

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

		pieceBB.ClearBit(uint64(fromSquare))
	}

	return moves
}

// generateQueenMoves generates the moves for the queens on the board
func generateQueenMoves(position Position, pieceBB BitBoard) []Move {
	moves := generateRookMoves(position, pieceBB)

	bishopMoves := generateBishopMoves(position, pieceBB)
	moves = append(moves, bishopMoves...)

	return moves
}

// generateKingMoves generates the moves for the kings on the board
func generateKingMoves(position Position, pieceBB BitBoard, includeCastling bool) []Move {
	moves := []Move{}

	for pieceBB > 0 {
		fromSquare := Square(pieceBB.TrailingZeros())

		moveBB := kingMoves[fromSquare]
		for moveBB > 0 {
			toSquare := Square(moveBB.TrailingZeros())

			piece, _ := position.GetPiece(toSquare)
			if piece == EmptyPiece {
				moves = append(moves, NewMove(fromSquare, toSquare, NormalMove, QuietMoveFlag))
			} else if piece.Color() != position.Turn() {
				move := NewMove(fromSquare, toSquare, NormalMove, QuietMoveFlag)
				move.WithCapture(piece)
				moves = append(moves, move)
			}

			moveBB.ClearBit(uint64(toSquare))
		}

		pieceBB.ClearBit(uint64(fromSquare))
	}

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

	return moves
}

// isLegalMove checks that the move would not result in an illegal position.
func (p Position) isLegalMove(move Move) bool {
	piece, _ := p.GetPiece(move.From)

	// check that the squares in between the king and rook are not attacked
	if move.Type() == CastleMove {
		attacked := false

		direction := east
		if move.To == C1 || move.To == C8 {
			direction = west
		}

		difference := move.FileDifference()
		for i := 0; i <= difference; i++ {
			square := move.To + Square(difference*int(direction))
			if p.IsSquareAttacked(square) {
				attacked = true
				break
			}
		}

		return !attacked
	}

	// check that the king is not moving into an attacked square
	if piece.Type() == King {
		return !p.IsSquareAttacked(move.To)
	}

	return true
}

// GenerateMoves generates all legal moves in the position.
func (position Position) GenerateMoves() []Move {
	moves := []Move{}

	colorBB := position.GetColorBB(position.turn)

	pawnBB := position.GetPieceBB(Pawn)
	pawnMoves := generatePawnMoves(position, pawnBB&colorBB)
	moves = append(moves, pawnMoves...)

	knightBB := position.GetPieceBB(Knight)
	knightMoves := generateKnightMoves(position, knightBB&colorBB)
	moves = append(moves, knightMoves...)

	bishopBB := position.GetPieceBB(Bishop)
	bishopMoves := generateBishopMoves(position, bishopBB&colorBB)
	moves = append(moves, bishopMoves...)

	rookBB := position.GetPieceBB(Rook)
	rookMoves := generateRookMoves(position, rookBB&colorBB)
	moves = append(moves, rookMoves...)

	queenBB := position.GetPieceBB(Queen)
	queenMoves := generateQueenMoves(position, queenBB&colorBB)
	moves = append(moves, queenMoves...)

	kingBB := position.GetPieceBB(King)
	kingMoves := generateKingMoves(position, kingBB&colorBB, true)
	moves = append(moves, kingMoves...)

	legalMoves := []Move{}
	for _, move := range moves {
		if position.isLegalMove(move) {
			legalMoves = append(legalMoves, move)
		}
	}

	return legalMoves
}
