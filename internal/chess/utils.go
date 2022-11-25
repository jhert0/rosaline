package chess

// pawnDirection returns the Direction they move in based on the given color.
func pawnDirection(color Color) direction {
	if color == White {
		return north
	}

	return south
}
