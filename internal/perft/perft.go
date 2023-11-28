package perft

import (
	"fmt"
	"rosaline/internal/chess"
)

func Perft(position chess.Position, depth int, print bool) uint64 {
	if depth == 0 {
		return 1
	}

	leaf := depth == 2

	var nodes uint64 = 0
	moves := position.GenerateMoves(chess.LegalMoveGeneration)
	for _, move := range moves {
		err := position.MakeMove(move)
		if err != nil {
			panic(err)
		}

		var count uint64
		if leaf {
			count = uint64(len(position.GenerateMoves(chess.LegalMoveGeneration)))
		} else {
			count = Perft(position, depth-1, false)
		}

		nodes += count

		if print {
			fmt.Printf("%s: %d\n", move, count)
		}

		position.Undo()
	}

	return nodes
}
