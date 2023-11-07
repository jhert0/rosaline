package search

import "rosaline/internal/chess"

type Searcher interface {
	Search(position chess.Position, depth int) ScoredMove
}
