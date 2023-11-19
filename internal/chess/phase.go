package chess

import "fmt"

type Phase uint8

const (
	OpeningPhase Phase = iota
	EndgamePhase
)

func (p Phase) String() string {
	switch p {
	case OpeningPhase:
		return "Opening"
	case EndgamePhase:
		return "Endgame"
	}

	panic(fmt.Sprintf("trying to convert unknown phase '%d' to string", p))
}
