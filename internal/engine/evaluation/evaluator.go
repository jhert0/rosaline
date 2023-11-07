package evaluation

import "rosaline/internal/chess"

type Evaluator interface {
	Evaluate(position chess.Position) int
}
