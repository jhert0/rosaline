package evaluation

import "rosaline/internal/chess"

func evaluationMultiplier(color chess.Color) int {
	if color == chess.White {
		return 1
	}

	return -1
}
