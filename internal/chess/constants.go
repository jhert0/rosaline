package chess

import "strings"

const (
	StartingFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

type CastlingRights uint8

const (
	WhiteCastleKingside  CastlingRights = 0b0100
	WhiteCastleQueenside CastlingRights = 0b1000
	WhiteCastleBoth      CastlingRights = WhiteCastleKingside | WhiteCastleQueenside

	BlackCastleKingside  CastlingRights = 0b0001
	BlackCastleQueenside CastlingRights = 0b0010
	BlackCastleBoth      CastlingRights = BlackCastleKingside | BlackCastleQueenside
)

func (rights CastlingRights) String() string {
	var builder strings.Builder

	if (rights & WhiteCastleKingside) > 0 {
		builder.WriteString("K")
	}

	if (rights & WhiteCastleQueenside) > 0 {
		builder.WriteString("Q")
	}

	if (rights & BlackCastleKingside) > 0 {
		builder.WriteString("k")
	}

	if (rights & BlackCastleQueenside) > 0 {
		builder.WriteString("q")
	}

	return builder.String()
}

type direction int8

// These directions are from white's perspective.
const (
	north direction = 8
	south direction = -8

	east direction = 1
	west direction = -1
)
