package chess

type GameStatus uint8

const (
	InProgress GameStatus = iota
	WhiteCheckmated
	BlackCheckmated
	Draw
	Stalemate
	WhiteResigned
	BlackResigned
)

// Game represents the current state of a chess game.
type Game struct {
	Position    Position   // The current position.
	status      GameStatus // The status of the game.
	drawOffered bool       // Whether a draw has been offered.
}

// NewGame returns a new chess game.
func NewGame(fen string) (Game, error) {
	position, err := NewPosition(fen)
	if err != nil {
		return Game{}, err
	}

	return Game{
		Position:    position,
		status:      InProgress,
		drawOffered: false,
	}, nil
}

// MakeUciMove makes a move from the given uci move.
func (g *Game) MakeUciMove(uci string) error {
	err := g.Position.MakeUciMove(uci)
	if err != nil {
		return err
	}

	// update status of the game if necessary
	opponent := g.Position.turn.OpposingSide()
	if g.Position.IsCheckmated(opponent) {
		if opponent == White {
			g.status = WhiteCheckmated
		} else {
			g.status = BlackCheckmated
		}
	} else if g.Position.IsStalemate(opponent) {
		g.status = Stalemate
	} else if g.Position.IsDraw() {
		g.status = Draw
	}

	return nil
}

// Status returns the status of the game.
func (g Game) Status() GameStatus {
	return g.status
}

// OfferDraw offers the opponent to end the game in a draw.
func (g *Game) OfferDraw() {
	g.drawOffered = true
}

// AcceptDraw accepts a draw offer.
func (g *Game) AcceptDraw() {
	g.status = Draw
}

// RejectDraw rejects a draw offer.
func (g *Game) RejectDraw() {
	g.drawOffered = false
}

// Resign ends the game with the resigning color losing.
func (g *Game) Resign(color Color) {
	if color == White {
		g.status = WhiteResigned
	} else {
		g.status = BlackResigned
	}
}
