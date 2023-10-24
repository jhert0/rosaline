package perft

import (
	"rosaline/internal/chess"
	"testing"
)

func perftTest(t *testing.T, position chess.Position, depth int, expectedNodes uint64) {
	nodes := Perft(position, depth, false)
	if nodes != expectedNodes {
		t.Fatalf("%s: expected '%d' nodes at depth '%d' got '%d'", t.Name(), expectedNodes, depth, nodes)
	}
}

func TestPerft(t *testing.T) {
	position, err := chess.NewPosition(chess.StartingFen)
	if err != nil {
		t.Fatalf("%s: error occured creating position for fen: %s", t.Name(), chess.StartingFen)
	}

	perftTest(t, position, 1, 20)
	perftTest(t, position, 2, 400)
}
