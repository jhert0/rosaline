package engine

import (
	"rosaline/internal/chess"
	"rosaline/internal/engine/evaluation"
	"rosaline/internal/engine/search"
)

type Engine struct {
	Name string // The name of the engine.

	Game      chess.Game // The current game the engine is playing/analyzing.
	searcher  search.NegamaxSearcher
	evaluator evaluation.Evaluator

	DefaultDepth int
}

// NewEngine creates a new Engine.
func NewEngine() Engine {
	evaluator := evaluation.NewEvaluator()

	return Engine{
		Name:         "Rosaline",
		Game:         chess.Game{},
		searcher:     search.NewNegamaxSearcher(evaluator),
		evaluator:    evaluator,
		DefaultDepth: 4,
	}
}

// NewGame starts a new game with the given fen.
func (e *Engine) NewGame(fen string) error {
	game, err := chess.NewGame(fen)
	if err != nil {
		return err
	}

	e.Game = game
	e.searcher.Reset()

	return nil
}

// Evaluate evaluates the current position and returns the score for it.
//
// A negative score means black is winning, a positive score means white is winning, and a score of 0 means the position is a draw.
func (e Engine) Evaluate() int {
	return e.evaluator.Evaluate(e.Game.Position)
}

// Search searches for the best move found with the given depth and returns it.
func (e Engine) Search(depth int) search.ScoredMove {
	return e.searcher.Search(e.Game.Position, depth)
}

// PlayBestMove plays the best the move it found after searching and then returns it.
func (e *Engine) PlayBestMove() search.ScoredMove {
	bestMove := e.Search(e.DefaultDepth)

	err := e.Game.MakeMove(bestMove.Move)
	if err != nil {
		panic(err)
	}

	return bestMove
}
