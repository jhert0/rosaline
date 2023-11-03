package chess

import (
	"errors"
)

var ErrInvalidPosition = errors.New("invalid postion")
var ErrInvalidFen = errors.New("invalid fen")
var ErrInvalidMove = errors.New("invalid move")
