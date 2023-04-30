package chess

// pawnDirection returns the Direction they move in based on the given color.
func pawnDirection(color Color) direction {
	if color == White {
		return north
	}

	return south
}

// pawnPromotionRank returns the rank where a pawn will get promoted.
func pawnPromotionRank(color Color) int {
	if color == White {
		return 8
	}

	return 1
}

// pawnStartingRank returns the rank that a pawn will start on
// for the standard starting position.
func pawnStartingRank(color Color) int {
	if color == White {
		return 2;
	}

	return 7
}
