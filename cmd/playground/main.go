package main

import (
	"bufio"
	"fmt"
	"os"
	"rosaline/internal/chess"
	"rosaline/internal/perft"
	"rosaline/internal/utils"
	"strconv"
	"strings"
	"time"
)

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

		command, args := utils.ParseCommand(scanner.Text())
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

			if args[0] == "null" {
				position.MakeNullMove()
			} else {
				err := position.MakeUciMove(args[0])
				if err != nil {
					fmt.Println(err)
				}
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
			moves := position.GenerateMoves(chess.LegalMoveGeneration)
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
			number := perft.Perft(position, depth, true)
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
			fmt.Println("Hash:", position.Hash())
			break
		case "attackers":
			for i := 0; i < 64; i++ {
				square := chess.Square(i)
				fmt.Printf("%s: ", square.ToAlgebraic())
				attackers := position.GetAttackers(square)

				fmt.Print("[")
				for attackers > 0 {
					square := chess.Square(attackers.PopLsb())
					fmt.Printf("%s", square)

					if attackers != 0 {
						fmt.Print(", ")
					}
				}
				fmt.Println("]")
			}

			break
		case "quit":
			break loop
		}

		fmt.Println()
	}
}
