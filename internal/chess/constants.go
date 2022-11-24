package chess

const (
	StartingFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

type CastlingRights uint8

const (
	WhiteCastleKingside  CastlingRights = 0b0001
	WhiteCastleQueenside CastlingRights = 0b0010
	WhiteCastleBoth      CastlingRights = WhiteCastleKingside | WhiteCastleQueenside

	BlackCastleKingside  CastlingRights = 0b0100
	BlackCastleQueenside CastlingRights = 0b1000
	BlackCastleBoth      CastlingRights = BlackCastleKingside | BlackCastleQueenside
)
