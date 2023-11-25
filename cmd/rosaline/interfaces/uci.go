package interfaces

import (
	"bufio"
	"fmt"
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
				results := i.searcher.Search(position, DefaultDepth)
				position.MakeMove(results.BestMove)
				fmt.Println("bestmove", results.BestMove)
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
