package chess

import "strings"

// Color represents the color of a chess piece.
type Color int8

const (
	// NoColor represents no color.
	NoColor Color = iota
	// White represents the color white.
	White
	// Black represents the color black.
	Black
)

func ColorFromString(s string) Color {
	switch strings.ToLower(s) {
	case "w":
		return White
	case "b":
		return Black
	}
	return NoColor
}

// Other returns the opposite color of the receiver.
func (c Color) Other() Color {
	switch c {
	case White:
		return Black
	case Black:
		return White
	}
	return NoColor
}

// String implements the fmt.Stringer interface and returns.
// the color's FEN compatible notation.
func (c Color) String() string {
	switch c {
	case White:
		return "w"
	case Black:
		return "b"
	}
	return "-"
}

// Name returns a display friendly name.
func (c Color) Name() string {
	switch c {
	case White:
		return "White"
	case Black:
		return "Black"
	}
	return "No Color"
}

// PieceType is the type of a piece.
type PieceType int8

const (
	// NoPieceType represents a lack of piece type.
	NoPieceType PieceType = iota
	// King represents a king.
	King
	// Queen represents a queen.
	Queen
	// Rook represents a rook.
	Rook
	// Bishop represents a bishop.
	Bishop
	// Knight represents a knight.
	Knight
	// Pawn represents a pawn.
	Pawn
)

// PieceTypes returns a slice of all piece types.
func PieceTypes() [6]PieceType {
	return [6]PieceType{King, Queen, Rook, Bishop, Knight, Pawn}
}

func PieceTypeFromByte(b byte) PieceType {
	switch b {
	case 'k':
		return King
	case 'q':
		return Queen
	case 'r':
		return Rook
	case 'b':
		return Bishop
	case 'n':
		return Knight
	case 'p':
		return Pawn
	}
	return NoPieceType
}

func PieceTypeFromString(s string) PieceType {
	if len(s) != 1 {
		return NoPieceType
	}
	return PieceTypeFromByte(strings.ToLower(s)[0])
}

func (p PieceType) String() string {
	switch p {
	case King:
		return "k"
	case Queen:
		return "q"
	case Rook:
		return "r"
	case Bishop:
		return "b"
	case Knight:
		return "n"
	case Pawn:
		return "p"
	}
	return ""
}

func (p PieceType) Bytes() []byte {
	switch p {
	case King:
		return []byte{'k'}
	case Queen:
		return []byte{'q'}
	case Rook:
		return []byte{'r'}
	case Bishop:
		return []byte{'b'}
	case Knight:
		return []byte{'n'}
	case Pawn:
		return []byte{'p'}
	case NoPieceType:
		return []byte{}
	}
	return []byte{}
}

func (p PieceType) ToPolyglotPromotionValue() int {
	switch p {
	case Knight:
		return 1
	case Bishop:
		return 2
	case Rook:
		return 3
	case Queen:
		return 4
	default:
		return 0
	}
}

// Piece is a piece type with a color.
type Piece int8

const (
	// NoPiece represents no piece.
	NoPiece Piece = iota
	// WhiteKing is a white king.
	WhiteKing
	// WhiteQueen is a white queen.
	WhiteQueen
	// WhiteRook is a white rook.
	WhiteRook
	// WhiteBishop is a white bishop.
	WhiteBishop
	// WhiteKnight is a white knight.
	WhiteKnight
	// WhitePawn is a white pawn.
	WhitePawn
	// BlackKing is a black king.
	BlackKing
	// BlackQueen is a black queen.
	BlackQueen
	// BlackRook is a black rook.
	BlackRook
	// BlackBishop is a black bishop.
	BlackBishop
	// BlackKnight is a black knight.
	BlackKnight
	// BlackPawn is a black pawn.
	BlackPawn
)

// TODO: This is a constant slice
//
//nolint:gochecknoglobals // This is a constant slice.
var allPieces = []Piece{
	WhiteKing, WhiteQueen, WhiteRook, WhiteBishop, WhiteKnight, WhitePawn,
	BlackKing, BlackQueen, BlackRook, BlackBishop, BlackKnight, BlackPawn,
}

// NewPiece returns the piece matching the PieceType and Color.
// NoPiece is returned if the PieceType or Color isn't valid.
func NewPiece(t PieceType, c Color) Piece {
	for _, p := range allPieces {
		if p.Color() == c && p.Type() == t {
			return p
		}
	}
	return NoPiece
}

// Type returns the type of the piece.
func (p Piece) Type() PieceType {
	switch p {
	case WhiteKing, BlackKing:
		return King
	case WhiteQueen, BlackQueen:
		return Queen
	case WhiteRook, BlackRook:
		return Rook
	case WhiteBishop, BlackBishop:
		return Bishop
	case WhiteKnight, BlackKnight:
		return Knight
	case WhitePawn, BlackPawn:
		return Pawn
	}
	return NoPieceType
}

// Color returns the color of the piece.
func (p Piece) Color() Color {
	switch p {
	case WhiteKing, WhiteQueen, WhiteRook, WhiteBishop, WhiteKnight, WhitePawn:
		return White
	case BlackKing, BlackQueen, BlackRook, BlackBishop, BlackKnight, BlackPawn:
		return Black
	}
	return NoColor
}

// String implements the fmt.Stringer interface.
func (p Piece) String() string {
	return pieceUnicodes[int(p)]
}

// DarkString is equivalent to String() except colors reversed for terminal
// windows in dark mode
func (p Piece) DarkString() string {
	return pieceDarkUnicodes[int(p)]
}

// TODO: These are constant slices
var (
	//nolint:gochecknoglobals // This is a constant slice.
	pieceUnicodes = []string{" ", "♔", "♕", "♖", "♗", "♘", "♙", "♚", "♛", "♜", "♝", "♞", "♟"}
	//nolint:gochecknoglobals // This is a constant slice.
	pieceDarkUnicodes = []string{" ", "♚", "♛", "♜", "♝", "♞", "♟", "♔", "♕", "♖", "♗", "♘", "♙"}
)

// getFENChar returns the FEN character representation of a piece
// Returns a single byte representing the piece.
func (p Piece) getFENChar() byte {
	pieceType := p.Type()
	if pieceType < 0 || pieceType > 6 {
		return 0 // Invalid piece type
	}

	if p.Color() == White {
		return whitePiecesToFEN[pieceType]
	}
	return blackPiecesToFEN[pieceType]
}
