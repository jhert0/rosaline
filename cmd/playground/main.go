package main

import (
	"bufio"
	"fmt"
	"os"
	"rosaline/internal/chess"
	"strconv"
	"strings"
	"time"
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

		fmt.Println()

		command, args := parseCommand(scanner.Text())
		switch command {
		case "reset":
			position, err = chess.NewPosition(chess.StartingFen)
			if err != nil {
				panic(err)
			}
			break
		case "fen":
			if len(args) < 1 {
				fmt.Println(position.Fen())
			} else {
				fen := strings.Join(args, " ")

				if fen == "startpos" {
					fen = chess.StartingFen
				}

				newPosition, err := chess.NewPosition(fen)
				if err != nil {
					fmt.Println("invalid fen string:", err)
					break
				}

				position = newPosition
			}
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
		case "undo":
			if !position.CanUndo() {
				fmt.Println("unable to undo, no previous position availible")
				break
			}

			position.Undo()
			break
		case "?":
			moves := position.GenerateMoves()
			for _, move := range moves {
				fmt.Println(move)
			}
			break
		case "perft":
			depth := 1
			if len(args) >= 1 {
				depth, err = strconv.Atoi(args[0])
				if err != nil {
					fmt.Println("invalid argument provided for depth")
					break
				}
			}

			start := time.Now()
			number := perft(position, depth, true)
			elapsed := time.Since(start)

			fmt.Println()
			fmt.Println("Depth:", depth)
			fmt.Println("Nodes:", number)
			fmt.Printf("Total time: %f seconds\n", elapsed.Seconds())
			break
		case "debug":
			fmt.Println("Turn:", position.Turn())
			fmt.Printf("Castling Rights: %04b (%s)\n", position.CastlingRights(), position.CastlingRights())

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

		fmt.Println()
	}
}

func perft(position chess.Position, depth int, print bool) uint64 {
	if depth == 0 {
		return 1
	}

	leaf := depth == 2

	var nodes uint64 = 0
	moves := position.GenerateMoves()
	for _, move := range moves {
		err := position.MakeUciMove(move.String())
		if err != nil {
			panic(err)
		}

		var count uint64
		if leaf {
			count = uint64(len(position.GenerateMoves()))
		} else {
			count = perft(position, depth-1, false)
		}

		nodes += count

		if print {
			fmt.Printf("%s: %d\n", move, count)
		}

		position.Undo()
	}

	return nodes
}
