package chess

type Color uint8

const (
	NoColor Color = 0
	White   Color = 1
	Black   Color = 2
)

const (
	StartingFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

type CastlingRights uint8

const (
	WHITE_CASTLE_KINGSIDE  CastlingRights = 0b0001
	WHITE_CASTLE_QUEENSIDE CastlingRights = 0b0010
	WHITE_CASTLE_BOTH      CastlingRights = WHITE_CASTLE_KINGSIDE | WHITE_CASTLE_QUEENSIDE

	BLACK_CASTLE_KINGSIDE  CastlingRights = 0b0100
	BLACK_CASTLE_QUEENSIDE CastlingRights = 0b1000
	BLACK_CASTLE_BOTH      CastlingRights = BLACK_CASTLE_KINGSIDE | BLACK_CASTLE_QUEENSIDE
)
