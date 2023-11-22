package search

import "rosaline/internal/chess"

type ScoredMove struct {
	Move  chess.Move // The move that needs to be made to get the attached score.
	Score int        // The score of the position after the move has been made.
}

// NewScoredMove creates a new ScoredMove.
func NewScoredMove(move chess.Move, score int) ScoredMove {
	return ScoredMove{
		Move:  move,
		Score: score,
	}
}

func (m ScoredMove) String() string {
	return m.Move.String()
}
