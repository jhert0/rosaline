package interfaces

import (
	"bufio"
	"fmt"
	"os"
	"rosaline/internal/chess"
	"rosaline/internal/evaluation"
	"rosaline/internal/perft"
	"rosaline/internal/search"
	"rosaline/internal/utils"
	"strconv"
	"strings"
)

type cliInterface struct {
	searcher  search.NegamaxSearcher
	evaluator evaluation.Evaluator
}

func NewCliProtocolHandler() cliInterface {
	evaluator := evaluation.NewEvaluator()
	return cliInterface{
		searcher:  search.NewNegamaxSearcher(evaluator),
		evaluator: evaluator,
	}
}

func (i cliInterface) Loop() {
	scanner := bufio.NewScanner(os.Stdin)

	position, _ := chess.NewPosition(chess.StartingFen)

	for {
		fmt.Print(position.Turn())
		fmt.Print("> ")

		scanner.Scan()
		cmd, args := utils.ParseCommand(scanner.Text())

		if cmd == "quit" {
			break
		} else if cmd == "display" {
			position.Print()
		} else if cmd == "perft" {
			depth := 1
			if len(args) > 0 {
				var err error
				depth, err = strconv.Atoi(args[0])
				if err != nil {
					fmt.Println("invalid argument provided for depth")
					break
				}
			}

			perft.Perft(position, depth, true)
		} else if cmd == "moves" {
			moves := position.GenerateMoves(chess.LegalMoveGeneration)
			for _, move := range moves {
				fmt.Println(move)
			}
		} else if cmd == "move" {
			if len(args) < 1 {
				fmt.Println("move requires a uci formated move as an argument")
				continue
			}

			err := position.MakeUciMove(args[0])
			if err != nil {
				fmt.Println(err)
			}
		} else if cmd == "undo" {
			position.Undo()
		} else if cmd == "go" {
			depth := DefaultDepth
			if len(args) >= 1 {
				var err error
				depth, err = strconv.Atoi(args[0])
				if err != nil {
					fmt.Println(err)
					continue
				}
			}

			results := i.searcher.Search(position, depth)
			fmt.Println("best move:", results.BestMove)
			fmt.Println("score:", results.Score)
			fmt.Println("elapsed:", results.Time.Seconds())
		} else if cmd == "evaluate" {
			score := i.evaluator.Evaluate(position)
			fmt.Println("score:", score)
		} else if cmd == "play" {
			results := i.searcher.Search(position, DefaultDepth)
			position.MakeMove(results.BestMove)
			fmt.Println("played:", results.BestMove)
		} else if cmd == "fen" {
			fmt.Println(position.Fen())
		} else if cmd == "setfen" {
			if len(args) < 1 {
				fmt.Println("setfen requires a fen as an argument")
				continue
			}

			fen := strings.Join(args, " ")
			if fen == "startpos" {
				fen = chess.StartingFen
			}

			p, err := chess.NewPosition(fen)
			if err != nil {
				fmt.Println(err)
				continue
			}

			i.searcher.Reset()
			position = p
		} else if cmd == "help" {
			fmt.Println("display                      displays the current position")
			fmt.Println("fen                          displays the current positions fen")
			fmt.Println("setfen [fen | startpos]      changes the position to the given fen")
			fmt.Println("perft [depth]                runs move generation test code to the specified depth")
			fmt.Println("moves                        displays the legal moves for the current position")
			fmt.Println("move [uci]                   make the given uci formatted move")
			fmt.Println("undo                         undos the last move")
			fmt.Println("go                           searches for the best move in the current position")
			fmt.Println("evaluate                     evaluates the current position")
			fmt.Println("play                         finds and plays the best move")
			fmt.Println("help                         displays this message")
			fmt.Println("quit                         exits the program")
		} else {
			fmt.Println("unknown command:", cmd)
		}
	}
}
