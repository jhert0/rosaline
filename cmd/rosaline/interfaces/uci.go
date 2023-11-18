package interfaces

import (
	"bufio"
	"fmt"
	"os"
	"rosaline/internal/chess"
	"rosaline/internal/engine"
	"rosaline/internal/utils"
)

type uciInterface struct {
}

func NewUciProtocolHandler() uciInterface {
	return uciInterface{}
}

func (i uciInterface) Loop() {
	engine := engine.NewEngine()

	scanner := bufio.NewScanner(os.Stdin)

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
			engine.NewGame(chess.StartingFen)
			break
		case "position":
			lastMove := args[len(args)-1]
			engine.Game.MakeUciMove(lastMove)
			break
		case "go":
			move := engine.PlayBestMove()
			fmt.Println("bestmove", move)
			break
		case "quit":
			break loop
		}
	}
}
