package main

import (
	"bufio"
	"fmt"
	"os"
	"rosaline/internal/chess"
	"strings"
)

func parseCommand(line string) (string, []string) {
	parts := strings.Split(line, " ")
	if len(parts) < 1 {
		return "", []string{}
	}

	cmd := parts[0]
	args := []string{}
	if len(parts) > 1 {
		args = parts[1:len(parts)]
	}

	return cmd, args
}

func main() {
	position, err := chess.NewPosition(chess.StartingFen)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
loop:
	for {
		position.Print()

		turn := position.Turn()
		fmt.Printf("%s> ", turn)

		scanner.Scan()

		command, args := parseCommand(scanner.Text())
		switch command {
		case "fen":
			fmt.Println(position.Fen())
			break
		case "startpos":
			if len(args) < 1 {
				fmt.Println("startpos requires an argument")
				break
			}

			fen := strings.Join(args, " ")
			newPosition, err := chess.NewPosition(fen)
			if err != nil {
				fmt.Println("invalid fen string:", err)
				break
			}

			position = newPosition

			break
		case "move":
			if len(args) < 1 {
				fmt.Println("move requires a uci formatted move as an argument")
				break
			}

			err := position.MakeUciMove(args[0])
			if err != nil {
				fmt.Println(err)
			}

			break
		case "debug":
			fmt.Println("Turn:", position.Turn())
			fmt.Printf("Castling Rights: %04b\n", position.CastlingRights())

			enPassant := position.EnPassant()
			fmt.Printf("En Passant: %d ", enPassant)
			if enPassant.IsValid() {
				fmt.Printf("(%s)", enPassant.ToAlgebraic())
			} else {
				fmt.Printf("(-)")
			}
			fmt.Println()

			fmt.Println("FEN:", position.Fen())
			break
		case "quit":
			break loop
		}
	}
}
