package chess

import (
	"fmt"
	"strconv"
	"strings"
)

// Decodes FEN notation into a GameState.  An error is returned
// if there is a parsing error.  FEN notation format:
// rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
func decodeFEN(fen string) (*Position, error) {
	fen = strings.TrimSpace(fen)
	parts := strings.Split(fen, " ")
	if len(parts) != 6 {
		return nil, fmt.Errorf("chess: fen invalid notation %s must have 6 sections", fen)
	}
	b, err := fenBoard(parts[0])
	if err != nil {
		return nil, err
	}
	turn, ok := fenTurnMap[parts[1]]
	if !ok {
		return nil, fmt.Errorf("chess: fen invalid turn %s", parts[1])
	}
	rights, err := formCastleRights(parts[2])
	if err != nil {
		return nil, err
	}
	sq, err := formEnPassant(parts[3])
	if err != nil {
		return nil, err
	}
	halfMoveClock, err := strconv.Atoi(parts[4])
	if err != nil || halfMoveClock < 0 {
		return nil, fmt.Errorf("chess: fen invalid half move clock %s", parts[4])
	}
	moveCount, err := strconv.Atoi(parts[5])
	if err != nil || moveCount < 1 {
		return nil, fmt.Errorf("chess: fen invalid move count %s", parts[5])
	}
	return &Position{
		board:           b,
		turn:            turn,
		castleRights:    rights,
		enPassantSquare: sq,
		halfMoveClock:   halfMoveClock,
		moveCount:       moveCount,
	}, nil
}

// generates board from fen format: rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR
func fenBoard(boardStr string) (*Board, error) {
	rankStrs := strings.Split(boardStr, "/")
	if len(rankStrs) != 8 {
		return nil, fmt.Errorf("chess: fen invalid board %s", boardStr)
	}
	m := map[Square]Piece{}
	for i, rankStr := range rankStrs {
		rank := Rank(7 - i)
		fileMap, err := fenFormRank(rankStr)
		if err != nil {
			return nil, err
		}
		for file, piece := range fileMap {
			m[NewSquare(file, rank)] = piece
		}
	}
	return NewBoard(m), nil
}

// fenFormRank converts a FEN rank string to a map of pieces
func fenFormRank(rankStr string) (map[File]Piece, error) {
	m := make(map[File]Piece, 8)
	var count int

	for i := 0; i < len(rankStr); i++ {
		c := rankStr[i]

		// Handle empty squares (digits 1-8)
		if c >= '1' && c <= '8' {
			count += int(c - '0')
			continue
		}

		// Get piece from lookup table
		piece := fenCharToPiece[c]
		if piece == NoPiece {
			return nil, fmt.Errorf("chess: invalid character in rank %q", rankStr)
		}

		m[File(count)] = piece
		count++
	}

	if count != 8 {
		return nil, fmt.Errorf("chess: invalid rank %q", rankStr)
	}

	return m, nil
}

func formCastleRights(castleStr string) (CastleRights, error) {
	// check for duplicates aka. KKkq right now is valid
	for _, s := range []string{"K", "Q", "k", "q", "-"} {
		if strings.Count(castleStr, s) > 1 {
			return "-", fmt.Errorf("chess: fen invalid castle rights %s", castleStr)
		}
	}
	for _, r := range castleStr {
		c := fmt.Sprintf("%c", r)
		switch c {
		case "K", "Q", "k", "q", "-":
		default:
			return "-", fmt.Errorf("chess: fen invalid castle rights %s", castleStr)
		}
	}
	return CastleRights(castleStr), nil
}

func formEnPassant(enPassant string) (Square, error) {
	if enPassant == "-" {
		return NoSquare, nil
	}
	sq := strToSquareMap[enPassant]
	if sq == NoSquare || !(sq.Rank() == Rank3 || sq.Rank() == Rank6) {
		return NoSquare, fmt.Errorf("chess: fen invalid En Passant square %s", enPassant)
	}
	return sq, nil
}

var (
	// whitePiecesToFEN provides direct mapping for white pieces to FEN characters
	whitePiecesToFEN = [7]byte{
		0,   // NoType (index 0)
		'K', // King   (index 1)
		'Q', // Queen  (index 2)
		'R', // Rook   (index 3)
		'B', // Bishop (index 4)
		'N', // Knight (index 5)
		'P', // Pawn   (index 6)
	}

	// blackPiecesToFEN provides direct mapping for black pieces to FEN characters
	blackPiecesToFEN = [7]byte{
		0,   // NoType (index 0)
		'k', // King   (index 1)
		'q', // Queen  (index 2)
		'r', // Rook   (index 3)
		'b', // Bishop (index 4)
		'n', // Knight (index 5)
		'p', // Pawn   (index 6)
	}

	fenTurnMap = map[string]Color{
		"w": White,
		"b": Black,
	}

	// Direct lookup array for FEN characters to pieces
	fenCharToPiece = [128]Piece{
		'K': WhiteKing,
		'Q': WhiteQueen,
		'R': WhiteRook,
		'B': WhiteBishop,
		'N': WhiteKnight,
		'P': WhitePawn,
		'k': BlackKing,
		'q': BlackQueen,
		'r': BlackRook,
		'b': BlackBishop,
		'n': BlackKnight,
		'p': BlackPawn,
	}
)
