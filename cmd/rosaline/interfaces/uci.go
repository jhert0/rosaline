package interfaces

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"rosaline/internal/chess"
	"rosaline/internal/evaluation"
	"rosaline/internal/search"
	"rosaline/internal/utils"
)

type uciInterface struct {
	searcher  search.NegamaxSearcher
	evaluator evaluation.Evaluator
}

func NewUciProtocolHandler() uciInterface {
	evaluator := evaluation.NewEvaluator()
	return uciInterface{
		searcher:  search.NewNegamaxSearcher(evaluator),
		evaluator: evaluator,
	}
}

func (i uciInterface) Loop() {
	scanner := bufio.NewScanner(os.Stdin)

	position, _ := chess.NewPosition(chess.StartingFen)

loop:
	for {
		scanner.Scan()
		cmd, args := utils.ParseCommand(scanner.Text())

		switch cmd {
		case "uci":
			fmt.Println("id name rosaline")
			fmt.Println("id author rosaline contributors")
			fmt.Println("uciok")
			break
		case "isready":
			fmt.Println("readyok")
			break
		case "ucinewgame":
			i.searcher.Reset()
			position, _ = chess.NewPosition(chess.StartingFen)
			break
		case "position":
			lastMove := args[len(args)-1]
			position.MakeUciMove(lastMove)
			break
		case "go":
			go func() {
				var bestMove chess.Move
				var bestScore int = math.MinInt

				for depth := 1; depth <= 4; depth++ {
					results := i.searcher.Search(position, depth)

					if i.searcher.Stopped() {
						break
					}

					if results.Score > bestScore {
						bestMove = results.BestMove
						bestScore = results.Score
					}

					fmt.Printf("info depth %d score cp %d nodes %d nps %d time %d\n", depth, results.Score, results.Nodes, results.NPS, results.Time.Milliseconds())
				}

				position.MakeMove(bestMove)
				fmt.Println("bestmove", bestMove)
			}()
			break
		case "stop":
			i.searcher.Stop()
			break
		case "quit":
			i.searcher.Stop()
			break loop
		}
	}
}
