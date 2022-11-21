package chess

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Position is a representation of the current state of the game.
type Position struct {
	turn Color

	whiteBB BitBoard
	blackBB BitBoard

	pawnBB   BitBoard
	bishopBB BitBoard
	knightBB BitBoard
	rookBB   BitBoard
	queenBB  BitBoard
	kingBB   BitBoard

	enPassant      Square
	castlingRights CastlingRights
	fiftyMoveClock int
	fullMoves      int

	previous *Position
}

// NewPositions creates a Position from the given FEN.
func NewPosition(fen string) (Position, error) {
	fenParts := strings.Split(fen, " ")
	if len(fenParts) < 6 {
		return Position{}, errors.New("too few sections in fen")
	}

	ranks := strings.Split(fenParts[0], "/")
	if len(ranks) < 8 {
		return Position{}, errors.New("too few ranks in board")
	}

	position := Position{}
	position.whiteBB = NewBitBoard(0)
	position.blackBB = NewBitBoard(0)

	position.pawnBB = NewBitBoard(0)
	position.knightBB = NewBitBoard(0)
	position.bishopBB = NewBitBoard(0)
	position.rookBB = NewBitBoard(0)
	position.queenBB = NewBitBoard(0)
	position.kingBB = NewBitBoard(0)

	rankNumber := 8
	for _, rank := range ranks {
		fileNumber := 1

		for _, character := range rank {
			index := uint64(SquareFromRankFile(rankNumber, fileNumber))
			fileIncrement := 1

			switch character {
			case 'p':
				position.blackBB.SetBit(index)
				position.pawnBB.SetBit(index)
				break
			case 'n':
				position.blackBB.SetBit(index)
				position.knightBB.SetBit(index)
				break
			case 'b':
				position.blackBB.SetBit(index)
				position.bishopBB.SetBit(index)
				break
			case 'r':
				position.blackBB.SetBit(index)
				position.rookBB.SetBit(index)
				break
			case 'q':
				position.blackBB.SetBit(index)
				position.queenBB.SetBit(index)
				break
			case 'k':
				position.blackBB.SetBit(index)
				position.kingBB.SetBit(index)
				break

			case 'P':
				position.whiteBB.SetBit(index)
				position.pawnBB.SetBit(index)
				break
			case 'N':
				position.whiteBB.SetBit(index)
				position.knightBB.SetBit(index)
				break
			case 'B':
				position.whiteBB.SetBit(index)
				position.bishopBB.SetBit(index)
				break
			case 'R':
				position.whiteBB.SetBit(index)
				position.rookBB.SetBit(index)
				break
			case 'Q':
				position.whiteBB.SetBit(index)
				position.queenBB.SetBit(index)
				break
			case 'K':
				position.whiteBB.SetBit(index)
				position.kingBB.SetBit(index)
				break

			case '1':
			case '2':
			case '3':
			case '4':
			case '5':
			case '6':
			case '7':
			case '8':
				fileIncrement, _ = strconv.Atoi(string(character))
				break
			}

			fileNumber += fileIncrement
		}

		rankNumber -= 1
	}

	// parse who's turn it is in the current position
	if fenParts[1] == "w" {
		position.turn = WHITE
	} else if fenParts[2] == "b" {
		position.turn = BLACK
	} else {
		return Position{}, errors.New(fmt.Sprintf("invalid character: %s for turn", fenParts[1]))
	}

	// parse castling rights
	castlingRights := fenParts[2]
	if castlingRights == "-" {
		position.castlingRights = 0
	} else {
		if strings.Contains(castlingRights, "K") {
			position.castlingRights |= WHITE_CASTLE_KINGSIDE
		}

		if strings.Contains(castlingRights, "Q") {
			position.castlingRights |= WHITE_CASTLE_QUEENSIDE
		}

		if strings.Contains(castlingRights, "k") {
			position.castlingRights |= BLACK_CASTLE_KINGSIDE
		}

		if strings.Contains(castlingRights, "q") {
			position.castlingRights |= BLACK_CASTLE_QUEENSIDE
		}
	}

	// parse en passant square
	if fenParts[3] == "-" {
		position.enPassant = -1
	} else {
		square, err := SquareFromAlgebraic(fenParts[3])
		if err != nil {
			return Position{}, err
		}

		position.enPassant = square
	}

	// parse half moves
	moveClock, err := strconv.Atoi(fenParts[4])
	if err != nil {
		return Position{}, errors.New(fmt.Sprintf("invalid value: %s for fifty move clock", fenParts[4]))
	}

	position.fiftyMoveClock = moveClock

	// parse full moves
	fullMoves, err := strconv.Atoi(fenParts[5])
	if err != nil {
		return Position{}, errors.New(fmt.Sprintf("invalid value: %s for full moves", fenParts[5]))
	}

	position.fullMoves = fullMoves

	position.previous = nil

	return position, nil
}

func (p Position) Turn() Color {
	return p.turn
}

func (p Position) FullMoves() int {
	return p.fullMoves
}

func (p Position) EnPassantPossible() bool {
	return p.enPassant > -1
}
