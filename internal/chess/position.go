package chess

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Position is a representation of the current state of the game.
type Position struct {
	turn Color // Whos turn it is currently.

	whiteBB BitBoard // BitBoard for all wihte pieces.
	blackBB BitBoard // BitBoard for all black pieces.

	pawnBB   BitBoard // BitBoard for all pawns.
	bishopBB BitBoard // BitBoard for all bishops.
	knightBB BitBoard // BitBoard for all knights.
	rookBB   BitBoard // BitBoard for all rooks.
	queenBB  BitBoard // BitBoard for all queens.
	kingBB   BitBoard // BitBoard for all kings.

	enPassant      Square         // The square where en passant is posssible.
	castlingRights CastlingRights // The current castling rights for both players.
	fiftyMoveClock int            // Number of moves since a capture or a pawn has moved. This is stored in half moves.
	fullMoves      int            // Number of moves played in the game.

	previous *Position // The previous Position.
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

// Fen gets the FEN for the current position.
func (p Position) Fen() string {
	var builder strings.Builder

	// write the board
	for rank := 8; rank >= 1; rank-- {
		for file := 1; file <= 8; {
			fileIncrement := 1

			square := SquareFromRankFile(rank, file)
			if p.PieceAt(square) {
				piece, _ := p.GetPiece(square)
				builder.WriteString(string(piece.Character()))
			} else {
				for f := file + 1; f <= 8; f++ {
					nextSquare := SquareFromRankFile(rank, f)
					if p.PieceAt(nextSquare) {
						break
					}

					fileIncrement++
				}

				builder.WriteString(fmt.Sprintf("%d", fileIncrement))
			}

			file += fileIncrement
		}

		if rank != 1 {
			builder.WriteString("/")
		}
	}

	builder.WriteString(" ")

	// write whos turn it is
	if p.turn == White {
		builder.WriteString("w")
	} else {
		builder.WriteString("b")
	}

	builder.WriteString(" ")

	// write castling rights
	if p.HasCastlingRights(WhiteCastleKingside) {
		builder.WriteString("K")
	}

	if p.HasCastlingRights(WhiteCastleQueenside) {
		builder.WriteString("Q")
	}

	if p.HasCastlingRights(BlackCastleKingside) {
		builder.WriteString("k")
	}

	if p.HasCastlingRights(BlackCastleQueenside) {
		builder.WriteString("q")
	}

	builder.WriteString(" ")

	// write en passant square
	if p.enPassant == -1 {
		builder.WriteString("-")
	} else {
		builder.WriteString(p.enPassant.ToAlgebraic())
	}

	builder.WriteString(" ")

	// write fifty move clock
	builder.WriteString(fmt.Sprintf("%d", p.fiftyMoveClock))

	builder.WriteString(" ")

	// write full moves
	builder.WriteString(fmt.Sprintf("%d", p.fullMoves))

	return builder.String()
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

func (p Position) EnPassant() Square {
	return p.enPassant
}

func (p Position) CastlingRights() CastlingRights {
	return p.castlingRights
}

func (p Position) EnPassantPossible() bool {
	return p.enPassant > -1
}

// GetColorBB returns the bitboard for the given color.
func (p Position) GetColorBB(color Color) BitBoard {
	if color == White {
		return p.whiteBB
	}

	if color == Black {
		return p.blackBB
	}

	panic("requested bitboard for NoColor")
}

// GetPieceBB returns the bitboard for the given PieceType.
func (p Position) GetPieceBB(pieceType PieceType) BitBoard {
	switch pieceType {
	case Pawn:
		return p.pawnBB
	case Knight:
		return p.knightBB
	case Bishop:
		return p.bishopBB
	case Rook:
		return p.rookBB
	case Queen:
		return p.queenBB
	case King:
		return p.kingBB
	}

	panic(fmt.Sprintf("requested bitboard for unknown piece type: %d", pieceType))
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

func (p *Position) setPiece(square Square, piece Piece) {
	index := uint64(square)

	if piece.Color() == White {
		p.whiteBB.SetBit(index)
	} else {
		p.blackBB.SetBit(index)
	}

	switch piece.Type() {
	case Pawn:
		p.pawnBB.SetBit(index)
		break
	case Knight:
		p.knightBB.SetBit(index)
		break
	case Bishop:
		p.bishopBB.SetBit(index)
		break
	case Rook:
		p.rookBB.SetBit(index)
		break
	case Queen:
		p.queenBB.SetBit(index)
		break
	case King:
		p.kingBB.SetBit(index)
		break
	}
}

func (p *Position) clearPiece(square Square, piece Piece) {
	index := uint64(square)

	if piece.Color() == White {
		p.whiteBB.ClearBit(index)
	} else {
		p.blackBB.ClearBit(index)
	}

	switch piece.Type() {
	case Pawn:
		p.pawnBB.ClearBit(index)
		break
	case Knight:
		p.knightBB.ClearBit(index)
		break
	case Bishop:
		p.bishopBB.ClearBit(index)
		break
	case Rook:
		p.rookBB.ClearBit(index)
		break
	case Queen:
		p.queenBB.ClearBit(index)
		break
	case King:
		p.kingBB.ClearBit(index)
		break
	}
}

// makeMove applies the move to the current position.
func (p *Position) makeMove(move Move) error {
	movingPiece, err := p.GetPiece(move.From)
	if err != nil {
		return err
	}

	capturePiece, _ := p.GetPiece(move.To)
	if movingPiece.Color() == capturePiece.Color() {
		return errors.New("trying to capture piece of same color")
	}

	copy := p.Copy()

	p.enPassant = -1   // clear en passant square, this will be set later if needed
	p.fiftyMoveClock++ // increment the fifty move clock, this will be cleared later if needed

	switch move.Type() {
	case NormalMove:
		// clear the fifty move clock, a pawn has moved or a capture has happened
		if movingPiece.Type() == Pawn || capturePiece.Type() != None {
			p.fiftyMoveClock = 0
		}

		if movingPiece.Type() == Pawn && move.RankDifference() == 2 {
			opposingSide := p.turn.OpposingSide()

			// check to see if a pawn is on a valid square for en passant
			pawn1, _ := p.GetPiece(move.To + Square(west))
			pawn2, _ := p.GetPiece(move.To + Square(east))
			if (pawn1.Type() == Pawn && pawn1.Color() == opposingSide) || (pawn2.Type() == Pawn && pawn2.Color() == opposingSide) {
				p.enPassant = move.To + Square(pawnDirection(opposingSide))
			}
		}

		if movingPiece.Type() == Rook {
			switch move.From {
			case A1:
				p.castlingRights &= ^WhiteCastleQueenside
				break
			case A8:
				p.castlingRights &= ^BlackCastleQueenside
				break
			case H1:
				p.castlingRights &= ^WhiteCastleKingside
				break
			case H8:
				p.castlingRights &= ^BlackCastleKingside
				break
			}
		}

		if movingPiece.Type() == King {
			if p.turn == White {
				p.castlingRights &= ^WhiteCastleBoth
			} else {
				p.castlingRights &= ^BlackCastleBoth
			}
		}

		p.clearPiece(move.From, movingPiece)
		p.setPiece(move.To, movingPiece)

		if capturePiece.Type() != None {
			p.clearPiece(move.To, capturePiece)

			if capturePiece.Type() == Rook {
				switch move.To {
				case A1:
					p.castlingRights &= ^WhiteCastleQueenside
					break
				case A8:
					p.castlingRights &= ^BlackCastleQueenside
					break
				case H1:
					p.castlingRights &= ^WhiteCastleKingside
					break
				case H8:
					p.castlingRights &= ^BlackCastleKingside
					break
				}
			}
		}

		if move.IsPromotion() {
			p.clearPiece(move.To, movingPiece) // remove the original piece

			// place the newly promoted piece
			promotedPiece := NewPiece(move.PromotionPiece(), p.turn)
			p.setPiece(move.To, promotedPiece)
		}
		break
	case CastleMove:
		// move the king to its new square
		p.clearPiece(move.From, movingPiece)
		p.setPiece(move.To, movingPiece)

		// move the rook to its new square
		switch move.To {
		case C1:
			rook, _ := p.GetPiece(A1)
			p.clearPiece(A1, rook)
			p.setPiece(D1, rook)
			break
		case G1:
			rook, _ := p.GetPiece(H1)
			p.clearPiece(H1, rook)
			p.setPiece(F1, rook)
			break
		case C8:
			rook, _ := p.GetPiece(C8)
			p.clearPiece(C8, rook)
			p.setPiece(F8, rook)
			break
		case G8:
			rook, _ := p.GetPiece(A8)
			p.clearPiece(A8, rook)
			p.setPiece(D8, rook)
			break
		}

		if p.turn == White {
			p.castlingRights &= ^WhiteCastleBoth
		} else {
			p.castlingRights &= ^BlackCastleBoth
		}
		break
	case EnPassantMove:
		// move the pawn to its new square
		p.clearPiece(move.From, movingPiece)
		p.setPiece(move.To, movingPiece)

		// remove the captured pawn
		captureSquare := move.To + Square(pawnDirection(p.turn.OpposingSide()))
		capturedPawn, _ := p.GetPiece(captureSquare)
		p.clearPiece(captureSquare, capturedPawn)
		break
	}

	if p.turn == Black {
		p.fullMoves++
	}

	p.turn = p.turn.OpposingSide()
	p.previous = &copy

	return nil
}

// MakeUciMove makes a move from the give uci string.
func (p *Position) MakeUciMove(uci string) error {
	if len(uci) < 4 {
		return errors.New("uci move is too short")
	}

	from, err := SquareFromAlgebraic(uci[:2])
	if err != nil {
		return err
	}

	to, err := SquareFromAlgebraic(uci[2:4])
	if err != nil {
		return err
	}

	moveType := NormalMove

	switch uci {
	case "e1g1", "e1c1", "e8g8", "e8c8":
		moveType = CastleMove
		break
	}

	movingPiece, err := p.GetPiece(from)
	if err != nil {
		return err
	}

	if movingPiece.Type() == Pawn && to == p.enPassant {
		moveType = EnPassantMove
	}

	move := NewMove(from, to, moveType)

	capturePiece, err := p.GetPiece(to)
	if err == nil {
		move.WithCapture(capturePiece.Type())
	}

	if len(uci) > 4 {
		switch uci[4] {
		case 'n':
			move.WithPromotion(Knight)
			break
		case 'b':
			move.WithPromotion(Bishop)
			break
		case 'r':
			move.WithPromotion(Rook)
			break
		case 'q':
			move.WithPromotion(Queen)
			break
		}
	}

	return p.makeMove(move)
}

// Copy creates a copy of the current position.
func (p Position) Copy() Position {
	copy := Position{
		turn:           p.turn,
		whiteBB:        p.whiteBB,
		blackBB:        p.blackBB,
		pawnBB:         p.pawnBB,
		bishopBB:       p.bishopBB,
		knightBB:       p.knightBB,
		rookBB:         p.rookBB,
		queenBB:        p.queenBB,
		kingBB:         p.kingBB,
		enPassant:      p.enPassant,
		castlingRights: p.castlingRights,
		fiftyMoveClock: p.fiftyMoveClock,
		fullMoves:      p.fullMoves,
		previous:       p.previous,
	}

	return copy
}

// CanUndo returns if there is a previous Position that can be restored to.
func (p Position) CanUndo() bool {
	return p.previous != nil
}

// Undo restores the Position to the previous Position.
func (p *Position) Undo() {
	p.turn = p.previous.turn
	p.whiteBB = p.previous.whiteBB
	p.blackBB = p.previous.blackBB
	p.pawnBB = p.previous.pawnBB
	p.bishopBB = p.previous.bishopBB
	p.knightBB = p.previous.knightBB
	p.rookBB = p.previous.rookBB
	p.queenBB = p.previous.queenBB
	p.kingBB = p.previous.kingBB
	p.enPassant = p.previous.enPassant
	p.castlingRights = p.previous.castlingRights
	p.fiftyMoveClock = p.previous.fiftyMoveClock
	p.fullMoves = p.previous.fullMoves
	p.previous = p.previous.previous
}
