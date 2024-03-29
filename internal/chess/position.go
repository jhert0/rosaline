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

	attackersBB [64]BitBoard // BitBoard for attackers of each square.

	squares [64]Piece // Keeps track of what piece is on each square. Used for faster lookups.

	enPassant               Square         // The square where en passant is posssible.
	castlingRights          CastlingRights // The current castling rights for both players.
	lastIrreversibleMovePly int            // The ply at which the last irreversible move happened. An irreversible move is a pawn move or capture.
	fiftyMoveClock          int            // Number of moves since a capture or a pawn has moved. This is stored in half moves.
	plies                   int            // Number of half moves in the game.

	hash uint64 // The zobrist hash of the current position.

	repetitions int // The number of times the current position has ocurred.

	previous *Position // The previous Position.
}

// NewPositions creates a Position from the given FEN.
func NewPosition(fen string) (Position, error) {
	fenParts := strings.Split(fen, " ")
	if len(fenParts) < 6 {
		return Position{}, fmt.Errorf("%w: too few sections in fen", ErrInvalidFen)
	}

	ranks := strings.Split(fenParts[0], "/")
	if len(ranks) < 8 {
		return Position{}, fmt.Errorf("%w: too few ranks in board", ErrInvalidFen)
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

	position.squares = [64]Piece{}

	rankNumber := 8
	for _, rank := range ranks {
		fileNumber := 1

		for _, character := range rank {
			index := SquareFromRankFile(rankNumber, fileNumber)
			fileIncrement := 1

			switch character {
			case 'p':
				position.setPiece(index, NewPiece(Pawn, Black))
				break
			case 'n':
				position.setPiece(index, NewPiece(Knight, Black))
				break
			case 'b':
				position.setPiece(index, NewPiece(Bishop, Black))
				break
			case 'r':
				position.setPiece(index, NewPiece(Rook, Black))
				break
			case 'q':
				position.setPiece(index, NewPiece(Queen, Black))
				break
			case 'k':
				position.setPiece(index, NewPiece(King, Black))
				break

			case 'P':
				position.setPiece(index, NewPiece(Pawn, White))
				break
			case 'N':
				position.setPiece(index, NewPiece(Knight, White))
				break
			case 'B':
				position.setPiece(index, NewPiece(Bishop, White))
				break
			case 'R':
				position.setPiece(index, NewPiece(Rook, White))
				break
			case 'Q':
				position.setPiece(index, NewPiece(Queen, White))
				break
			case 'K':
				position.setPiece(index, NewPiece(King, White))
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
		return Position{}, fmt.Errorf("%w: invalid character '%s' for turn", ErrInvalidFen, fenParts[1])
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
			return Position{}, fmt.Errorf("%w: invalid en passant square: %s", ErrInvalidFen, fenParts[3])
		}

		position.enPassant = square
	}

	// parse half moves
	moveClock, err := strconv.Atoi(fenParts[4])
	if err != nil {
		return Position{}, fmt.Errorf("%w: invalid value '%s' for fifty move clock", ErrInvalidFen, fenParts[4])
	}

	position.fiftyMoveClock = moveClock

	// parse full moves
	fullMoves, err := strconv.Atoi(fenParts[5])
	if err != nil {
		return Position{}, fmt.Errorf("%w: invalid value '%s' for full moves", ErrInvalidFen, fenParts[5])
	}

	position.plies = (fullMoves - 1) * 2
	if position.turn == Black {
		position.plies++
	}

	position.lastIrreversibleMovePly = position.plies

	position.hash = generateHash(position)
	position.previous = nil

	if ok, err := position.IsValid(); !ok {
		return Position{}, err
	}

	position.updateAttackers()

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
			if p.IsSquareOccupied(square) {
				piece, _ := p.GetPieceAt(square)
				builder.WriteString(string(piece.Character()))
			} else {
				for f := file + 1; f <= 8; f++ {
					nextSquare := SquareFromRankFile(rank, f)
					if p.IsSquareOccupied(nextSquare) {
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
	if p.castlingRights == 0 {
		builder.WriteString("-")
	} else {
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
	builder.WriteString(fmt.Sprintf("%d", p.FullMoves()))

	return builder.String()
}

// HasCastlingRights checks if the given castling rights are available.
func (p Position) HasCastlingRights(rights CastlingRights) bool {
	return (p.castlingRights & rights) > 0
}

// PieceAt checks if there is a piece at the give square.
func (p Position) IsSquareOccupied(square Square) bool {
	return p.whiteBB.IsBitSet(uint64(square)) || p.blackBB.IsBitSet(uint64(square))
}

// GetPieceColor gets the color of the piece at the give square.
func (p Position) GetPieceColor(square Square) (Color, error) {
	if p.whiteBB.IsBitSet(uint64(square)) {
		return White, nil
	}

	if p.blackBB.IsBitSet(uint64(square)) {
		return Black, nil
	}

	return NoColor, errors.New(fmt.Sprintf("no piece exists at: %s", square.ToAlgebraic()))
}

// GetPieceAt gets the piece at the given square.
func (p Position) GetPieceAt(square Square) (Piece, error) {
	if !p.IsSquareOccupied(square) {
		return EmptyPiece, errors.New(fmt.Sprintf("no piece exists at: %s", square.ToAlgebraic()))
	}

	return p.squares[square], nil
}

// GetKingSquare returns the square of the specified color's King is on.
func (p Position) GetKingSquare(color Color) Square {
	kingBB := p.GetColorBB(color) & p.kingBB
	return Square(kingBB.Lsb())
}

// IsPieceAt returns whether a piece matching the piece type and color are
func (p Position) IsPieceAt(square Square, pieceType PieceType, color Color) bool {
	piece, _ := p.GetPieceAt(square)
	return piece.Type() == pieceType && piece.Color() == color
}

func (p Position) Turn() Color {
	return p.turn
}

// Plies return the number of plies that have been made.
func (p Position) Plies() int {
	return p.plies
}

// FullMoves returns the number of full moves that have been made.
func (p Position) FullMoves() int {
	return (p.plies / 2) + 1
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

// squaresEmpty returns whether all squares are empty.
func (p Position) squaresEmpty(squares []Square) bool {
	for _, square := range squares {
		piece, _ := p.GetPieceAt(square)
		if piece != EmptyPiece {
			return false
		}
	}

	return true
}

// IsValid returns whether the position is playable, i.e no more than 8 pawns, one king, etc.
func (p Position) IsValid() (bool, error) {
	// check that neither side has more than 16 pieces or zero pieces
	whitePieces := p.whiteBB.PopulationCount()
	if whitePieces > 16 || whitePieces == 0 {
		return false, fmt.Errorf("%w: white has an invald number of pieces: %d", ErrInvalidPosition, whitePieces)
	}

	blackPieces := p.blackBB.PopulationCount()
	if blackPieces > 16 || blackPieces == 0 {
		return false, fmt.Errorf("%w: black has an invalid number of pieces: %d", ErrInvalidPosition, blackPieces)
	}

	// check that neither side has more than 8 pawns
	whitePawns := p.pawnBB & p.whiteBB
	if whitePawns.PopulationCount() > 8 {
		return false, fmt.Errorf("%w: white has too many pawns: %d", ErrInvalidPosition, whitePawns)
	}

	blackPawns := p.pawnBB & p.blackBB
	if blackPawns.PopulationCount() > 8 {
		return false, fmt.Errorf("%w: black has too many pawns: %d", ErrInvalidPosition, blackPawns)
	}

	// check that both sides only have one king
	whiteKing := p.kingBB & p.whiteBB
	if whiteKing.PopulationCount() != 1 {
		return false, fmt.Errorf("%w: white has an invalid number of kings: %d", ErrInvalidPosition, whiteKing.PopulationCount())
	}

	blackKings := p.kingBB & p.blackBB
	if blackKings.PopulationCount() != 1 {
		return false, fmt.Errorf("%w: black has an invalid number of kings: %d", ErrInvalidPosition, blackKings.PopulationCount())
	}

	// the opposing side can't be in check
	if p.IsKingInCheck(p.turn.OpposingSide()) {
		return false, fmt.Errorf("%w: the opponent can't be in check", ErrInvalidPosition)
	}

	// it is not possible to be in check by 3 or more pieces
	checkers := p.NumberOfCheckers(p.turn)
	if p.NumberOfCheckers(p.turn) >= 3 {
		return false, fmt.Errorf("%w: the %s king is in check by %d pieces which is not possible", ErrInvalidPosition, p.turn, checkers)
	}

	return true, nil
}

func (p Position) Print() {
	for rank := 8; rank >= 1; rank-- {
		fmt.Printf("%d | ", rank)

		for file := 1; file <= 8; file++ {
			square := SquareFromRankFile(rank, file)
			piece, _ := p.GetPieceAt(square)
			if p.IsSquareOccupied(square) {
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
	if !square.IsValid() {
		panic(fmt.Sprintf("invalid square '%d' passed to setPiece", square))
	}

	if piece == EmptyPiece {
		panic(fmt.Sprintf("trying to set an empty piece on square %s", square))
	}

	p.squares[square] = piece

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
	default:
		panic(fmt.Sprintf("setPiece: unkown piece type '%d' encountered", piece.Type()))
	}
}

func (p *Position) clearPiece(square Square) {
	if !square.IsValid() {
		panic(fmt.Sprintf("invalid square '%d' passed to clearPiece", square))
	}

	piece, _ := p.GetPieceAt(square)
	if piece == EmptyPiece {
		panic(fmt.Sprintf("trying to clear a square %s with no piece", square))
	}

	p.squares[square] = EmptyPiece

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
	default:
		panic(fmt.Sprintf("clearPiece: unkown piece type '%d' encountered", piece.Type()))
	}
}

// MakeMove applies the move to the current position.
func (p *Position) MakeMove(move Move) error {
	movingPiece, err := p.GetPieceAt(move.From())
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidMove, err)
	}

	if movingPiece.Color() != p.turn {
		return fmt.Errorf("%w: tyring to move opponent's piece with %s", ErrInvalidMove, move)
	}

	capturePiece, _ := p.GetPieceAt(move.To())
	if movingPiece.Color() == capturePiece.Color() {
		return fmt.Errorf("%w: trying to capture piece of same color with %s", ErrInvalidMove, move)
	}

	to := move.To()
	from := move.From()

	copy := p.Copy()

	p.enPassant = -1   // clear en passant square, this will be set later if needed
	p.fiftyMoveClock++ // increment the fifty move clock, this will be cleared later if needed

	switch move.Type() {
	case QuietMove:
		if movingPiece.Type() == Pawn && move.RankDifference() == 2 {
			opposingSide := p.turn.OpposingSide()

			// check to see if a pawn is on a valid square for en passant
			westPawn, _ := p.GetPieceAt(to + Square(west))
			eastPawn, _ := p.GetPieceAt(to + Square(east))
			if (westPawn.Type() == Pawn && westPawn.Color() == opposingSide && to.File() != 1) || (eastPawn.Type() == Pawn && eastPawn.Color() == opposingSide && to.File() != 8) {
				p.enPassant = to + Square(pawnDirection(opposingSide))
			}
		}

		if movingPiece.Type() == Rook {
			switch from {
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

		p.clearPiece(from)
		p.setPiece(to, movingPiece)
		break
	case CaptureMove:
		p.clearPiece(to)
		p.clearPiece(from)

		p.setPiece(to, movingPiece)

		if capturePiece.Type() == Rook {
			switch to {
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
		break
	case CastleMove:
		// move the king to it's new square
		p.clearPiece(from)
		p.setPiece(to, movingPiece)

		// move the rook to it's new square
		switch to {
		case C1:
			rook, _ := p.GetPieceAt(A1)
			p.clearPiece(A1)
			p.setPiece(D1, rook)
			break
		case G1:
			rook, _ := p.GetPieceAt(H1)
			p.clearPiece(H1)
			p.setPiece(F1, rook)
			break
		case C8:
			rook, _ := p.GetPieceAt(A8)
			p.clearPiece(A8)
			p.setPiece(D8, rook)
			break
		case G8:
			rook, _ := p.GetPieceAt(H8)
			p.clearPiece(H8)
			p.setPiece(F8, rook)
			break
		}

		if p.turn == White {
			p.castlingRights &= ^WhiteCastleBoth
		} else {
			p.castlingRights &= ^BlackCastleBoth
		}
		break
	case EnPassantMove:
		// move the pawn to it's new square
		p.clearPiece(from)
		p.setPiece(to, movingPiece)

		// remove the captured pawn
		captureSquare := to + Square(pawnDirection(p.turn.OpposingSide()))
		p.clearPiece(captureSquare)
		break
	}

	if move.IsPromotion() {
		p.clearPiece(to) // remove the original piece

		// place the newly promoted piece
		p.setPiece(to, move.PromotionPiece())
	}

	// clear the fifty move clock, a pawn has moved or a capture has happened
	if movingPiece.Type() == Pawn || move.IsCapture() {
		p.fiftyMoveClock = 0
	}

	p.updateAttackers()

	p.plies++

	if move.IsIrreversible() {
		p.lastIrreversibleMovePly = p.plies
	}

	p.repetitions = 0
	p.hash = generateHash(*p)
	p.turn = p.turn.OpposingSide()
	p.previous = &copy

	// determine the number of times this position has been reached
	numPlies := p.plies - p.lastIrreversibleMovePly
	previous := p.previous
	for i := 0; i < numPlies; i++ {
		if previous == nil {
			break
		}

		if previous.hash == p.hash {
			p.repetitions++
		}

		previous = previous.previous
	}

	return nil
}

// MakeUciMove makes a move from the given uci string.
func (p *Position) MakeUciMove(uci string) error {
	if len(uci) < 4 {
		return errors.New(fmt.Sprintf("invalid move: provided uci: '%s' is too short", uci))
	}

	from, err := SquareFromAlgebraic(uci[:2])
	if err != nil {
		return err
	}

	to, err := SquareFromAlgebraic(uci[2:4])
	if err != nil {
		return err
	}

	movingPiece, err := p.GetPieceAt(from)
	if err != nil {
		return err
	}

	moveType := QuietMove
	flags := NoMoveFlag

	if movingPiece.Type() == King {
		switch uci {
		case "e1g1", "e1c1", "e8g8", "e8c8":
			moveType = CastleMove
			break
		}
	} else if movingPiece.Type() == Pawn && to == p.enPassant {
		moveType = EnPassantMove
	} else if movingPiece.Type() == Pawn && FileDistance(to, from) == 0 { // if a pawn is moving and not moving to another file then it's a pawn push
		flags |= PawnPushMoveFlag
	}

	capturePiece, _ := p.GetPieceAt(to)
	if capturePiece != EmptyPiece {
		moveType = CaptureMove
	}

	move := NewMove(from, to, moveType)

	if flags != NoMoveFlag {
		move.WithFlags(flags)
	}

	if len(uci) > 4 {
		switch uci[4] {
		case 'n':
			move.WithPromotion(NewPiece(Knight, p.turn))
			break
		case 'b':
			move.WithPromotion(NewPiece(Bishop, p.turn))
			break
		case 'r':
			move.WithPromotion(NewPiece(Rook, p.turn))
			break
		case 'q':
			move.WithPromotion(NewPiece(Queen, p.turn))
			break
		}
	}

	return p.MakeMove(move)
}

// MakeNullMove switches sides without making an actual move.
func (p *Position) MakeNullMove() {
	copy := p.Copy()

	p.enPassant = -1
	p.plies++

	p.turn = p.turn.OpposingSide()

	p.previous = &copy
}

// updateAttackers updates the BitBoard that keeps track of attackers.
func (p *Position) updateAttackers() {
	p.attackersBB = [64]BitBoard{}

	ourAttacks := getAttackers(*p, p.turn)
	theirAttacks := getAttackers(*p, p.turn.OpposingSide())

	moves := append(ourAttacks, theirAttacks...)
	for _, move := range moves {
		p.attackersBB[move.To()].SetBit(uint64(move.From()))
	}
}

// IsSquareAttackedBy returns whether the given square is being attacked by the given color.
func (p Position) IsSquareAttackedBy(square Square, color Color) bool {
	return (p.GetAttackers(square) & p.GetColorBB(color)) != BitBoard(0)
}

// IsSquareAttacked returns whether the given square is attacked.
func (p Position) IsSquareAttacked(square Square) bool {
	return p.GetAttackers(square) != BitBoard(0)
}

// IsKingInCheck returns whether the given color's king is being attacked.
func (p Position) IsKingInCheck(color Color) bool {
	square := p.GetKingSquare(color)
	return p.IsSquareAttackedBy(square, color.OpposingSide())
}

// NumberOfCheckers returns the number of opposing pieces attacking the given color's king.
func (p Position) NumberOfCheckers(color Color) int {
	if !p.IsKingInCheck(color) {
		return 0
	}

	kingSquare := p.GetKingSquare(color)
	attackers := p.attackersBB[kingSquare] & p.GetColorBB(color.OpposingSide())

	return attackers.PopulationCount()
}

// IsCheckmated returns whether the specified color has been checkmated.
func (p Position) IsCheckmated(color Color) bool {
	if !p.IsKingInCheck(color) { // king must be in check to be checkmated
		return false
	}

	kingSquare := p.GetKingSquare(color)
	squares := SurroundingSquares(kingSquare)
	for _, square := range squares {
		piece, _ := p.GetPieceAt(square)

		// check that the surrounding square is not attacked and not occupied by one of our pieces
		if !p.IsSquareAttackedBy(square, color.OpposingSide()) && piece.Color() != color {
			return false
		}
	}

	return true
}

// IsDraw returns whether the position is a draw.
func (p Position) IsDraw() bool {
	if p.fiftyMoveClock >= 100 || p.repetitions >= 3 {
		return true
	}

	return p.whiteBB.PopulationCount() <= 1 && p.blackBB.PopulationCount() <= 1
}

// IsStalemate returns whether position is a stalemate due to the given color having no legal moves.
func (p Position) IsStalemate(color Color) bool {
	if p.IsKingInCheck(color) { // king can't be in check and also in stalemate
		return false
	}

	// temporarily switch to the color's turn if it is not currently their turn
	if p.turn != color {
		p.MakeNullMove()
		defer p.Undo()
	}

	moves := p.GenerateMoves(LegalMoveGeneration)
	return len(moves) == 0
}

// GetAttackers returns a BitBoard containing all pieces attacking the given Square.
func (p Position) GetAttackers(square Square) BitBoard {
	return p.attackersBB[square]
}

// Hash returns the hash for the current position.
func (p Position) Hash() uint64 {
	return p.hash
}

// Phase returns the phase the position is in.
func (p Position) Phase() Phase {
	if p.queenBB == 0 {
		return EndgamePhase
	}

	whiteQueens := p.whiteBB & p.queenBB
	whiteMinor := p.whiteBB & (p.knightBB | p.bishopBB)
	if whiteQueens.PopulationCount() <= 1 && whiteMinor <= 1 {
		return EndgamePhase
	}

	blackQueens := p.blackBB & p.queenBB
	blackMinor := p.blackBB & (p.knightBB | p.bishopBB)
	if blackQueens.PopulationCount() <= 1 && blackMinor <= 1 {
		return EndgamePhase
	}

	return OpeningPhase
}

// Copy creates a copy of the current position.
func (p Position) Copy() Position {
	copy := Position{
		turn:                    p.turn,
		whiteBB:                 p.whiteBB,
		blackBB:                 p.blackBB,
		pawnBB:                  p.pawnBB,
		bishopBB:                p.bishopBB,
		knightBB:                p.knightBB,
		rookBB:                  p.rookBB,
		queenBB:                 p.queenBB,
		kingBB:                  p.kingBB,
		attackersBB:             p.attackersBB,
		squares:                 p.squares,
		enPassant:               p.enPassant,
		castlingRights:          p.castlingRights,
		fiftyMoveClock:          p.fiftyMoveClock,
		lastIrreversibleMovePly: p.lastIrreversibleMovePly,
		plies:                   p.plies,
		hash:                    p.hash,
		repetitions:             p.repetitions,
		previous:                p.previous,
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
	p.attackersBB = p.previous.attackersBB
	p.squares = p.previous.squares
	p.enPassant = p.previous.enPassant
	p.castlingRights = p.previous.castlingRights
	p.fiftyMoveClock = p.previous.fiftyMoveClock
	p.lastIrreversibleMovePly = p.previous.lastIrreversibleMovePly
	p.plies = p.previous.plies
	p.hash = p.previous.hash
	p.repetitions = p.previous.repetitions
	p.previous = p.previous.previous
}
