package chess

import (
	"testing"
)

// Common test positions
var (
	// Starting position
	startingPos = mustPosition("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	// Middle game position with lots of possible moves
	middlePos = mustPosition("r1bqk2r/pppp1ppp/2n2n2/2b1p3/2B1P3/2N2N2/PPPP1PPP/R1BQK2R w KQkq - 0 1")

	// Endgame position with few pieces
	endPos = mustPosition("4k3/8/8/8/8/8/4P3/4K3 w - - 0 1")

	// Position with multiple possible pawn promotions
	promoPos = mustPosition("4k3/PPPP4/8/8/8/8/4pppp/4K3 w - - 0 1")
)

func BenchmarkStandardMoves(b *testing.B) {
	benchmarks := []struct {
		name      string
		pos       *Position
		wantFirst bool
	}{
		{"StartingPos_AllMoves", startingPos, false},
		{"StartingPos_FirstMove", startingPos, true},
		{"MiddleGame_AllMoves", middlePos, false},
		{"MiddleGame_FirstMove", middlePos, true},
		{"Endgame_AllMoves", endPos, false},
		{"Endgame_FirstMove", endPos, true},
		{"Promotions_AllMoves", promoPos, false},
		{"Promotions_FirstMove", promoPos, true},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			// Reset timer to exclude setup
			b.ResetTimer()

			// Enable allocation tracking
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				moves := standardMoves(bm.pos, bm.wantFirst)
				// Prevent compiler optimization
				if len(moves) == 0 {
					b.Fatal("unexpected zero moves")
				}
			}
		})
	}
}

// Benchmark specific scenarios
func BenchmarkStandardMoves_PawnPromotions(b *testing.B) {
	pos := promoPos
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		moves := standardMoves(pos, false)
		if len(moves) == 0 {
			b.Fatal("unexpected zero moves")
		}
	}
}

// Benchmark with different board sizes to understand allocation scaling
func BenchmarkStandardMoves_BoardDensity(b *testing.B) {
	positions := []struct {
		name string
		fen  string
	}{
		{"Empty", "4k3/8/8/8/8/8/8/4K3 w - - 0 1"},
		{"QuarterFull", "rnbqk3/pppp4/8/8/8/8/8/4K3 w - - 0 1"},
		{"HalfFull", "rnbqkbnr/pppp4/8/8/8/8/8/4K3 w - - 0 1"},
		{"Full", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},
	}

	for _, p := range positions {
		pos := mustPosition(p.fen)
		b.Run(p.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				moves := standardMoves(pos, false)
				if len(moves) == 0 && p.name != "Empty" {
					b.Fatal("unexpected zero moves")
				}
			}
		})
	}
}

// Helper function to convert FEN to Position
func mustPosition(fen string) *Position {
	fenObject, err := FEN(fen)
	pos := NewGame(fenObject).Position()
	if err != nil {
		panic(err)
	}
	return pos
}
