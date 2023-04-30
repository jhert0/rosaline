package chess

// IsSquareAttacked returns whether the given square is being attacked
func (p Position) IsSquareAttacked(square Square) bool {
	return p.GetAttackers(square) != BitBoard(0)
}

// GetAttackers returns a BitBoard containing all pieces attacking the given Square.
func (p Position) GetAttackers(square Square) BitBoard {
	return BitBoard(0)
}
