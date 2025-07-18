package chess

import (
	"testing"
)

var (
	//noline:gochecknoglobals // This is a test file
	// TODO: This is a legacy counter for generating unique labels. (will be removed in the future)
	validFENs = []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
		"7k/8/8/8/8/8/8/R6K w - - 0 1",
		"7k/8/8/8/8/8/8/2B1KB2 w - - 0 1",
		"8/8/8/4k3/8/8/8/R3K2R w KQ - 0 1",
		"8/8/8/8/4k3/8/3KP3/8 w - - 0 1",
		"8/8/5k2/8/5K2/8/4P3/8 w - - 0 1",
		"r4rk1/1b2bppp/ppq1p3/2pp3n/5P2/1P1BP3/PBPPQ1PP/R4RK1 w - - 0 1",
		"3r1rk1/p3qppp/2bb4/2p5/3p4/1P2P3/PBQN1PPP/2R2RK1 w - - 0 1",
		"4r1k1/1b3p1p/ppq3p1/2p5/8/1P3R1Q/PBP3PP/7K w - - 0 1",
		"5k2/ppp5/4P3/3R3p/6P1/1K2Nr2/PP3P2/8 b - - 1 32",
		"rnbqkbnr/pp1ppppp/8/8/1Pp1PP2/8/P1PP2PP/RNBQKBNR b KQkq b3 0 3",
		"rnbqkbnr/p1ppppp1/7p/Pp6/8/8/1PPPPPPP/RNBQKBNR w KQkq b6 0 3",
		"rnbqkbnr/1pppppp1/7p/pP6/8/8/P1PPPPPP/RNBQKBNR w KQkq a6 0 3",
		"rnbqkbnr/1pppppp1/7p/pP6/4P3/8/P1PP1PPP/RNBQKBNR b KQkq e3 0 3",
	}

	validXFENs = []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
		"7k/8/8/8/8/8/8/R6K w - - 0 1",
		"7k/8/8/8/8/8/8/2B1KB2 w - - 0 1",
		"8/8/8/4k3/8/8/8/R3K2R w KQ - 0 1",
		"8/8/8/8/4k3/8/3KP3/8 w - - 0 1",
		"8/8/5k2/8/5K2/8/4P3/8 w - - 0 1",
		"r4rk1/1b2bppp/ppq1p3/2pp3n/5P2/1P1BP3/PBPPQ1PP/R4RK1 w - - 0 1",
		"3r1rk1/p3qppp/2bb4/2p5/3p4/1P2P3/PBQN1PPP/2R2RK1 w - - 0 1",
		"4r1k1/1b3p1p/ppq3p1/2p5/8/1P3R1Q/PBP3PP/7K w - - 0 1",
		"5k2/ppp5/4P3/3R3p/6P1/1K2Nr2/PP3P2/8 b - - 1 32",
		"rnbqkbnr/pp1ppppp/8/8/1Pp1PP2/8/P1PP2PP/RNBQKBNR b KQkq b3 0 3",
		"rnbqkbnr/p1ppppp1/7p/Pp6/8/8/1PPPPPPP/RNBQKBNR w KQkq b6 0 3",
		"rnbqkbnr/1pppppp1/7p/pP6/8/8/P1PPPPPP/RNBQKBNR w KQkq a6 0 3",
		"rnbqkbnr/1pppppp1/7p/pP6/4P3/8/P1PP1PPP/RNBQKBNR b KQkq - 0 3",
	}

	//nolint:gochecknoglobals // test data
	// TODO: This is a legacy counter for generating unique labels. (will be removed in the future)
	invalidFENs = []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPP/RNBQKBNR w KQkq - 0 1",
		"rnbqkbnr/pppppppp/8/8/4P2/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KKkq c6 0 2",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq c12 1 2",
		"7k/8/8/8/8/8/8/R6K w - - 0 -1",
		"7k/8/8/8/8/8/8/2B1KB2 w - - -1 1",
		"8/8/8/4k3/8/8/8/R3K2R w KQ - 0 0",
		"8/8/8/8/4k3/8/3KP3/8 c - - 0 1",
		"8/8/5k2/8/5K2/8/4P3P/8 w - - 0 1",
		"r4rk1/1b2bppp/ppq1p3/2pp3n/5P2/1P1BP3/PBPPQ1PP/R4RK1 w e4 - 0 1",
	}
)

func TestValidFENs(t *testing.T) {
	for idx, f := range validFENs {
		state, err := decodeFEN(f)
		if err != nil {
			t.Fatal("recieved unexpected error", err)
		}
		if f != state.String() {
			t.Fatalf("fen expected board string %s but got %s", f, state.String())
		}
		xfen := state.XFENString()
		if xfen != validXFENs[idx] {
			t.Fatalf("xfen for fen %v (%v) was %v but expected %v", idx, f,
				xfen, validXFENs[idx])
		}
	}
}

func TestInvalidFENs(t *testing.T) {
	for _, f := range invalidFENs {
		if _, err := decodeFEN(f); err == nil {
			t.Fatal("fen expected error from ", f)
		}
	}
}

func BenchmarkFenBoard(b *testing.B) {
	// Test cases representing different scenarios
	benchmarks := []struct {
		name string
		fen  string
	}{
		{
			name: "StartingPosition",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR",
		},
		{
			name: "EmptyBoard",
			fen:  "8/8/8/8/8/8/8/8",
		},
		{
			name: "MidGame",
			fen:  "r1bqk2r/pppp1ppp/2n2n2/2b1p3/2B1P3/2N2N2/PPPP1PPP/R1BQK2R",
		},
		{
			name: "ComplexPosition",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R",
		},
		{
			name: "EndGame",
			fen:  "4k3/8/8/8/8/8/4P3/4K3",
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			// Reset the timer to exclude setup time
			b.ResetTimer()
			b.ReportAllocs()
			// Run the benchmark
			for i := 0; i < b.N; i++ {
				board, err := fenBoard(bm.fen)
				if err != nil {
					b.Fatal(err)
				}
				// Prevent compiler optimizations
				_ = board
			}
		})
	}
}
