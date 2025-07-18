/*
Package chess provides position representation and manipulation for chess games.
The package implements complete position tracking including piece placement,
castling rights, en passant squares, and move counts. It supports standard chess
formats (FEN) and provides methods for position analysis and move validation.
Example usage:

	// Create starting position
	pos := StartingPosition()

	// Check valid moves
	moves := pos.ValidMoves()

	// Update position with move
	newPos := pos.Update(move)

	// Get FEN string
	fen := pos.String()
*/
package chess

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

// Side represents a side of the board.
type Side int

const (
	// KingSide is the right side of the board from white's perspective.
	KingSide Side = iota + 1
	// QueenSide is the left side of the board from white's perspective.
	QueenSide
)

// CastleRights holds the state of both sides castling abilities.
type CastleRights string

// CanCastle returns true if the given color and side combination can castle.
//
// Example:
//
//	if rights.CanCastle(White, KingSide) {
//	    // White can castle kingside
//	}
func (cr CastleRights) CanCastle(c Color, side Side) bool {
	char := "k"
	if side == QueenSide {
		char = "q"
	}
	if c == White {
		char = strings.ToUpper(char)
	}
	return strings.Contains(string(cr), char)
}

// String implements the fmt.Stringer interface and returns
// a FEN compatible string.  Ex. KQq.
func (cr CastleRights) String() string {
	return string(cr)
}

// Position represents a complete chess position state.
// It includes piece placement, castling rights, en passant squares,
// move counts, and side to move.
type Position struct {
	board           *Board       // Current board state
	castleRights    CastleRights // Available castling options
	validMoves      []Move       // Cache of legal moves
	halfMoveClock   int          // Half-move counter
	moveCount       int          // Full move counter
	turn            Color        // Side to move
	enPassantSquare Square       // En passant target square
	inCheck         bool         // Whether current side is in check
}

const (
	startFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" // Starting position FEN
)

// StartingPosition returns the starting position
// rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1.
func StartingPosition() *Position {
	pos, _ := decodeFEN(startFEN)
	return pos
}

// Update returns a new position resulting from the given move.
// The move isn't validated - use Game.Move() for validation.
// This method is optimized for move generation where validation
// is handled separately.
//
// Example:
//
//	newPos := pos.Update(move)
func (pos *Position) Update(m *Move) *Position {
	moveCount := pos.moveCount
	if pos.turn == Black {
		moveCount++
	}

	if m == nil {
		return &Position{
			board:           pos.board.copy(),
			turn:            pos.turn.Other(),
			castleRights:    pos.castleRights,
			enPassantSquare: NoSquare,
			halfMoveClock:   pos.halfMoveClock + 1,
			moveCount:       moveCount,
			inCheck:         false,
		}
	}

	ncr := pos.updateCastleRights(m)
	p := pos.board.Piece(m.s1)
	halfMove := pos.halfMoveClock
	if p.Type() == Pawn || m.HasTag(Capture) {
		halfMove = 0
	} else {
		halfMove++
	}
	b := pos.board.copy()
	b.update(m)
	return &Position{
		board:           b,
		turn:            pos.turn.Other(),
		castleRights:    ncr,
		enPassantSquare: pos.updateEnPassantSquare(m),
		halfMoveClock:   halfMove,
		moveCount:       moveCount,
		inCheck:         m.HasTag(Check),
	}
}

// ValidMoves returns all legal moves in the current position.
// The moves are cached for performance.
// TODO: Can we make this more efficient? Maybe using an iterator?
func (pos *Position) ValidMoves() []Move {
	if pos.validMoves != nil {
		return append([]Move(nil), pos.validMoves...)
	}
	pos.validMoves = engine{}.CalcMoves(pos, false)
	return append([]Move(nil), pos.validMoves...)
}

// Status returns the position's status as one of the outcome methods.
// Possible returns values include Checkmate, Stalemate, and NoMethod.
func (pos *Position) Status() Method {
	return engine{}.Status(pos)
}

// Board returns the position's board.
func (pos *Position) Board() *Board {
	return pos.board
}

// Turn returns the color to move next.
func (pos *Position) Turn() Color {
	return pos.turn
}

// ChangeTurn returns a new position with the turn changed.
func (pos *Position) ChangeTurn() *Position {
	pos.turn = pos.turn.Other()
	return pos
}

// HalfMoveClock returns the half-move clock (50-rule).
func (pos *Position) HalfMoveClock() int {
	return pos.halfMoveClock
}

// EnPassantSquare returns the en-passant square.
func (pos *Position) EnPassantSquare() Square {
	return pos.enPassantSquare
}

// CastleRights returns the castling rights of the position.
func (pos *Position) CastleRights() CastleRights {
	return pos.castleRights
}

// Ply returns the half-move number (increments every move).
func (pos *Position) Ply() int {
	if pos == nil {
		return 0
	}
	if pos.moveCount == 0 {
		return 0
	}

	if pos.turn == White {
		return (pos.moveCount-1)*2 + 1
	} else {
		return (pos.moveCount) * 2
	}
}

// String implements the fmt.Stringer interface and returns a
// string with the FEN format: rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1.
func (pos *Position) String() string {
	b := pos.board.String()
	t := pos.turn.String()
	c := pos.castleRights.String()
	sq := "-"
	if pos.enPassantSquare != NoSquare {
		sq = pos.enPassantSquare.String()
	}
	return fmt.Sprintf("%s %s %s %s %d %d", b, t, c, sq, pos.halfMoveClock, pos.moveCount)
}

// XFENString() is similar to String() except that it returns a string with
// the X-FEN format
func (pos *Position) XFENString() string {
	b := pos.board.String()
	t := pos.turn.String()
	c := pos.castleRights.String()
	sq := "-"
	if pos.enPassantSquare != NoSquare {
		// Check if there is a pawn in a position to capture en passant
		var rank Rank
		if pos.turn == White {
			rank = Rank5
		} else {
			rank = Rank4
		}
		// The en passant target square will always be on the rank opposite the current turn's pawns
		file := pos.enPassantSquare.File()
		potentialPawnFiles := []File{file - 1, file + 1} // Pawns that could capture en passant will be on an adjacent file

		for _, f := range potentialPawnFiles {
			if f < FileA || f > FileH { // Ensure file is within bounds
				continue
			}

			potentialPawnSquare := NewSquare(f, rank)
			potentialPawn := pos.board.Piece(potentialPawnSquare)
			if potentialPawn == NoPiece {
				continue
			}
			if potentialPawn.Type() != Pawn {
				continue
			}
			if potentialPawn.Color() == pos.turn {
				sq = pos.enPassantSquare.String()
				break
			}
		}
	}
	return fmt.Sprintf("%s %s %s %s %d %d", b, t, c, sq, pos.halfMoveClock, pos.moveCount)
}

// Hash returns a unique hash of the position.
func (pos *Position) Hash() [16]byte {
	b, _ := pos.MarshalBinary()
	return md5.Sum(b)
}

// MarshalText implements the encoding.TextMarshaler interface and
// encodes the position's FEN.
func (pos *Position) MarshalText() ([]byte, error) {
	return []byte(pos.String()), nil
}

// UnmarshalText implements the encoding.TextUnarshaler interface and
// assumes the data is in the FEN format.
func (pos *Position) UnmarshalText(text []byte) error {
	cp, err := decodeFEN(string(text))
	if err != nil {
		return err
	}
	pos.board = cp.board
	pos.castleRights = cp.castleRights
	pos.turn = cp.turn
	pos.enPassantSquare = cp.enPassantSquare
	pos.halfMoveClock = cp.halfMoveClock
	pos.moveCount = cp.moveCount
	pos.inCheck = isInCheck(cp)
	return nil
}

const (
	bitsCastleWhiteKing uint8 = 1 << iota
	bitsCastleWhiteQueen
	bitsCastleBlackKing
	bitsCastleBlackQueen
	bitsTurn
	bitsHasEnPassant
)

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (pos *Position) MarshalBinary() ([]byte, error) {
	boardBytes, err := pos.board.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(boardBytes)
	if err = binary.Write(buf, binary.BigEndian, uint8(pos.halfMoveClock)); err != nil {
		return nil, err
	}
	if err = binary.Write(buf, binary.BigEndian, uint16(pos.moveCount)); err != nil {
		return nil, err
	}
	if err = binary.Write(buf, binary.BigEndian, pos.enPassantSquare); err != nil {
		return nil, err
	}
	var b uint8
	if pos.castleRights.CanCastle(White, KingSide) {
		b |= bitsCastleWhiteKing
	}
	if pos.castleRights.CanCastle(White, QueenSide) {
		b |= bitsCastleWhiteQueen
	}
	if pos.castleRights.CanCastle(Black, KingSide) {
		b |= bitsCastleBlackKing
	}
	if pos.castleRights.CanCastle(Black, QueenSide) {
		b |= bitsCastleBlackQueen
	}
	if pos.turn == Black {
		b |= bitsTurn
	}
	if pos.enPassantSquare != NoSquare {
		b |= bitsHasEnPassant
	}
	if err = binary.Write(buf, binary.BigEndian, b); err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

// UnmarshalBinary implements the encoding.BinaryMarshaler interface.
func (pos *Position) UnmarshalBinary(data []byte) error {
	const size = 101
	if len(data) != size {
		return errors.New("chess: position binary data should consist of 101 bytes")
	}
	board := &Board{}
	if err := board.UnmarshalBinary(data[:96]); err != nil {
		return err
	}
	pos.board = board
	buf := bytes.NewBuffer(data[96:])
	halfMove := uint8(pos.halfMoveClock)
	if err := binary.Read(buf, binary.BigEndian, &halfMove); err != nil {
		return err
	}
	pos.halfMoveClock = int(halfMove)
	moveCount := uint16(pos.moveCount)
	if err := binary.Read(buf, binary.BigEndian, &moveCount); err != nil {
		return err
	}
	pos.moveCount = int(moveCount)
	if err := binary.Read(buf, binary.BigEndian, &pos.enPassantSquare); err != nil {
		return err
	}
	var b uint8
	if err := binary.Read(buf, binary.BigEndian, &b); err != nil {
		return err
	}
	pos.castleRights = ""
	pos.turn = White
	if b&bitsCastleWhiteKing != 0 {
		pos.castleRights += "K"
	}
	if b&bitsCastleWhiteQueen != 0 {
		pos.castleRights += "Q"
	}
	if b&bitsCastleBlackKing != 0 {
		pos.castleRights += "k"
	}
	if b&bitsCastleBlackQueen != 0 {
		pos.castleRights += "q"
	}
	if pos.castleRights == "" {
		pos.castleRights = "-"
	}
	if b&bitsTurn != 0 {
		pos.turn = Black
	}
	if b&bitsHasEnPassant == 0 {
		pos.enPassantSquare = NoSquare
	}
	pos.inCheck = isInCheck(pos)
	return nil
}

func (pos *Position) copy() *Position {
	return &Position{
		board:           pos.board.copy(),
		turn:            pos.turn,
		castleRights:    pos.castleRights,
		enPassantSquare: pos.enPassantSquare,
		halfMoveClock:   pos.halfMoveClock,
		moveCount:       pos.moveCount,
		inCheck:         pos.inCheck,
	}
}

func (pos *Position) updateCastleRights(m *Move) CastleRights {
	cr := string(pos.castleRights)
	p := pos.board.Piece(m.s1)
	if p == WhiteKing || m.s1 == H1 || m.s2 == H1 {
		cr = strings.ReplaceAll(cr, "K", "")
	}
	if p == WhiteKing || m.s1 == A1 || m.s2 == A1 {
		cr = strings.ReplaceAll(cr, "Q", "")
	}
	if p == BlackKing || m.s1 == H8 || m.s2 == H8 {
		cr = strings.ReplaceAll(cr, "k", "")
	}
	if p == BlackKing || m.s1 == A8 || m.s2 == A8 {
		cr = strings.ReplaceAll(cr, "q", "")
	}
	if cr == "" {
		cr = "-"
	}
	return CastleRights(cr)
}

func (pos *Position) updateEnPassantSquare(m *Move) Square {
	const squaresPerRank = 8
	p := pos.board.Piece(m.s1)
	if p.Type() != Pawn {
		return NoSquare
	}
	if pos.turn == White &&
		(bbForSquare(m.s1)&bbRank2) != 0 &&
		(bbForSquare(m.s2)&bbRank4) != 0 {
		return m.s2 - squaresPerRank
	} else if pos.turn == Black &&
		(bbForSquare(m.s1)&bbRank7) != 0 &&
		(bbForSquare(m.s2)&bbRank5) != 0 {
		return m.s2 + squaresPerRank
	}
	return NoSquare
}

// samePosition returns true if the two positions are the same.
func (pos *Position) samePosition(pos2 *Position) bool {
	return pos.board.String() == pos2.board.String() &&
		pos.turn == pos2.turn &&
		pos.castleRights.String() == pos2.castleRights.String() &&
		pos.enPassantSquare == pos2.enPassantSquare
}
