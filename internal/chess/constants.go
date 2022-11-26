package chess

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

type direction int8

// These directions are from white's perspective.
const (
	north direction = 8
	south direction = -8

	east direction = 1
	west direction = -1
)
