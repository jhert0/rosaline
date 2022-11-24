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

			case '1', '2', '3', '4', '5', '6', '7', '8':
				fileIncrement, _ = strconv.Atoi(string(character))
				break
			}

			fileNumber += fileIncrement
		}

		rankNumber -= 1
	}

	// parse who's turn it is in the current position
	if fenParts[1] == "w" {
		position.turn = White
	} else if fenParts[1] == "b" {
		position.turn = Black
	} else {
		return Position{}, errors.New(fmt.Sprintf("invalid character: %s for turn", fenParts[1]))
	}

	// parse castling rights
	castlingRights := fenParts[2]
	if castlingRights == "-" {
		position.castlingRights = 0
	} else {
		if strings.Contains(castlingRights, "K") {
			position.castlingRights |= WhiteCastleKingside
		}

		if strings.Contains(castlingRights, "Q") {
			position.castlingRights |= WhiteCastleQueenside
		}

		if strings.Contains(castlingRights, "k") {
			position.castlingRights |= BlackCastleKingside
		}

		if strings.Contains(castlingRights, "q") {
			position.castlingRights |= BlackCastleQueenside
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

// HasCastlingRights checks if the given castling rights are available.
func (p Position) HasCastlingRights(rights CastlingRights) bool {
	return (p.castlingRights & rights) > 0
}

// PieceAt checks if there is a piece at the give square.
func (p Position) PieceAt(square Square) bool {
	return p.whiteBB.BitSet(uint64(square)) || p.blackBB.BitSet(uint64(square))
}

// GetPieceColor gets the color of the piece at the give square.
func (p Position) GetPieceColor(square Square) (Color, error) {
	if p.whiteBB.BitSet(uint64(square)) {
		return White, nil
	}

	if p.blackBB.BitSet(uint64(square)) {
		return Black, nil
	}

	return NoColor, errors.New(fmt.Sprintf("no piece exists at: %s", square.ToAlgebraic()))
}

// GetPiece gets the piece at the given square.
func (p Position) GetPiece(square Square) (Piece, error) {
	if !p.PieceAt(square) {
		return NewEmptyPiece(), errors.New(fmt.Sprintf("no piece exists at: %s", square.ToAlgebraic()))
	}

	index := uint64(square)
	color, _ := p.GetPieceColor(square)

	if p.pawnBB.BitSet(index) {
		return NewPiece(Pawn, color), nil
	}

	if p.knightBB.BitSet(index) {
		return NewPiece(Knight, color), nil
	}

	if p.bishopBB.BitSet(index) {
		return NewPiece(Bishop, color), nil
	}

	if p.rookBB.BitSet(index) {
		return NewPiece(Rook, color), nil
	}

	if p.queenBB.BitSet(index) {
		return NewPiece(Queen, color), nil
	}

	if p.kingBB.BitSet(index) {
		return NewPiece(King, color), nil
	}

	return NewEmptyPiece(), nil
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

func (p Position) Print() {
	for rank := 8; rank >= 1; rank-- {
		fmt.Printf("%d | ", rank)

		for file := 1; file <= 8; file++ {
			square := SquareFromRankFile(rank, file)
			piece, _ := p.GetPiece(square)
			if p.PieceAt(square) {
				fmt.Printf(" %c ", piece.Character())
			} else {
				fmt.Print(" - ")
			}
		}

		fmt.Println()
	}

	fmt.Println("   +------------------------")
	fmt.Println("     a  b  c  d  e  f  g  h")
}
