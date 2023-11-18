package interfaces

import (
	"bufio"
	"fmt"
	"os"
	"rosaline/internal/chess"
	"rosaline/internal/engine"
	"rosaline/internal/perft"
	"rosaline/internal/utils"
	"strconv"
	"strings"
)

type cliInterface struct {
	scanner *bufio.Scanner
}

func NewCliProtocolHandler() cliInterface {
	return cliInterface{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (i cliInterface) Loop() {
	scanner := bufio.NewScanner(os.Stdin)

	engine := engine.NewEngine()
	engine.NewGame(chess.StartingFen)

	for {
		fmt.Print(engine.Game.Position.Turn())
		fmt.Print("> ")

		scanner.Scan()
		cmd, args := utils.ParseCommand(scanner.Text())

		if cmd == "quit" {
			break
		} else if cmd == "display" {
			engine.Game.Position.Print()
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

			perft.Perft(engine.Game.Position, depth, true)
		} else if cmd == "moves" {
			moves := engine.Game.Position.GenerateMoves(chess.LegalMoveGeneration)
			for _, move := range moves {
				fmt.Println(move)
			}
		} else if cmd == "move" {
			if len(args) < 1 {
				fmt.Println("move requires a uci formated move as an argument")
				continue
			}

			err := engine.Game.MakeUciMove(args[0])
			if err != nil {
				fmt.Println(err)
			}
		} else if cmd == "undo" {
			engine.Game.Position.Undo()
		} else if cmd == "go" {
			move := engine.Search(4)
			fmt.Println("best move:", move)
			fmt.Println("score:", move.Score)
		} else if cmd == "evaluate" {
			score := engine.Evaluate()
			fmt.Println("score:", score)
		} else if cmd == "play" {
			move := engine.PlayBestMove()
			fmt.Println("played:", move.Move)
		} else if cmd == "fen" {
			fmt.Println(engine.Game.Position.Fen())
		} else if cmd == "setfen" {
			if len(args) < 1 {
				fmt.Println("setfen requires a fen as an argument")
				continue
			}

			fen := strings.Join(args, " ")
			if fen == "startpos" {
				fen = chess.StartingFen
			}

			err := engine.NewGame(fen)
			if err != nil {
				fmt.Println(err)
			}
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
