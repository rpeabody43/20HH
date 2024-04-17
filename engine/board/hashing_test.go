package board

import (
	"testing"
)

func TestTranspositionHash(t *testing.T) {
	Init()

	fenPosition := FromFEN("r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3")

	transposition1 := StartPos()
	transposition1.UCIMakeMove("e2e4")
	transposition1.UCIMakeMove("e7e5")
	transposition1.UCIMakeMove("g1f3")
	transposition1.UCIMakeMove("b8c6")

	transposition2 := StartPos()
	transposition2.UCIMakeMove("g1f3")
	transposition2.UCIMakeMove("b8c6")
	transposition2.UCIMakeMove("e2e4")
	transposition2.UCIMakeMove("e7e5")
	transposition2.UCIMakeMove("d2d3")
	transposition2.UndoMove(NewMove(D2, D3, 0))

	if fenPosition.hash != transposition1.hash {
		t.Errorf("FEN hash (0x%x) != transposition 1 hash (0x%x)",
			fenPosition.hash, transposition1.hash,
		)
	}
	if fenPosition.hash != transposition2.hash {
		t.Errorf("FEN hash (0x%x) != transposition 2 hash (0x%x)",
			fenPosition.hash, transposition2.hash,
		)
	}
}

func TestHashDifference(t *testing.T) {
	Init()

	startPosition := StartPos()

	position1 := StartPos()
	position1.UCIMakeMove("e2e4")

	if startPosition.hash == position1.hash {
		t.Errorf("Starting hash (0x%x) == position 1 hash (0x%x)",
			startPosition.hash, position1.hash,
		)
	}
}

func TestIncrementalHashUpdates(t *testing.T) {
	Init()

	position1 := StartPos()
	position1.UCIMakeMove("e2e4")
	position1.UCIMakeMove("e7e5")
	position1.UCIMakeMove("g1f3")
	position1.UCIMakeMove("b8c6")
	incrementalHash := position1.hash
	position1.genHash()
	fromScratchHash := position1.hash

	if incrementalHash != fromScratchHash {
		t.Errorf("From scratch hash (0x%x) != incremental update hash (0x%x)",
			fromScratchHash, incrementalHash,
		)
	}
}

func TestPolyGlotHashingFEN(t *testing.T) {
	Init()

	var tests = []struct {
		name     string
		fen      string
		expected uint64
	}{
		{
			"starting position",
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			0x463b96181691fc9c,
		},
		{
			"position after e2e4",
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			0x823c9b50fd114196,
		},
		{
			"position after e2e4 d7d5",
			"rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2",
			0x0756b94461c50fb0,
		},
		{
			"position after e2e4 d7d5 e4e5",
			"rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 2",
			0x662fafb965db29d4,
		},
		{
			"position after e2e4 d7d5 e4e5 f7f5",
			"rnbqkbnr/ppp1p1pp/8/3pPp2/8/8/PPPP1PPP/RNBQKBNR w KQkq f6 0 3",
			0x22a48b5a8e47ff78,
		},
		{
			"position after e2e4 d7d5 e4e5 f7f5 e1e2",
			"rnbqkbnr/ppp1p1pp/8/3pPp2/8/8/PPPPKPPP/RNBQ1BNR b kq - 0 3",
			0x652a607ca3f242c1,
		},
		{
			"position after e2e4 d7d5 e4e5 f7f5 e1e2 e8f7",
			"rnbq1bnr/ppp1pkpp/8/3pPp2/8/8/PPPPKPPP/RNBQ1BNR w - - 0 4",
			0x00fdd303c946bdd9,
		},
		{
			"position after a2a4 b7b5 h2h4 b5b4 c2c4",
			"rnbqkbnr/p1pppppp/8/8/PpP4P/8/1P1PPPP1/RNBQKBNR b KQkq c3 0 3",
			0x3c8123ea7b067637,
		},
		{
			"position after a2a4 b7b5 h2h4 b5b4 c2c4 b4c3 a1a3",
			"rnbqkbnr/p1pppppp/8/8/P6P/R1p5/1P1PPPP1/1NBQKBNR b Kkq - 0 4",
			0x5c3f9b829b279560,
		},
	}

	for _, position := range tests {
		t.Run(position.name, func(t *testing.T) {
			board := FromFEN(position.fen)
			if board.hash != position.expected {
				t.Errorf("Hash was 0x%x instead of 0x%x",
					board.hash,
					position.expected,
				)
			}
		})
	}
}

func TestPolyGlotHashing(t *testing.T) {
	Init()

	var tests = []struct {
		name     string
		moves    []string
		expected uint64
	}{
		{
			"position after e2e4",
			[]string{"e2e4"},
			0x823c9b50fd114196,
		},
		{
			"position after e2e4 d7d5",
			[]string{"e2e4", "d7d5"},
			0x0756b94461c50fb0,
		},
		{
			"position after e2e4 d7d5 e4e5",
			[]string{"e2e4", "d7d5", "e4e5"},
			0x662fafb965db29d4,
		},
		{
			"position after e2e4 d7d5 e4e5 f7f5",
			[]string{"e2e4", "d7d5", "e4e5", "f7f5"},
			0x22a48b5a8e47ff78,
		},
		{
			"position after e2e4 d7d5 e4e5 f7f5 e1e2",
			[]string{"e2e4", "d7d5", "e4e5", "f7f5", "e1e2"},
			0x652a607ca3f242c1,
		},
		{
			"position after e2e4 d7d5 e4e5 f7f5 e1e2 e8f7",
			[]string{"e2e4", "d7d5", "e4e5", "f7f5", "e1e2", "e8f7"},
			0x00fdd303c946bdd9,
		},
		{
			"position after a2a4 b7b5 h2h4 b5b4 c2c4",
			[]string{"a2a4", "b7b5", "h2h4", "b5b4", "c2c4"},
			0x3c8123ea7b067637,
		},
		{
			"position after a2a4 b7b5 h2h4 b5b4 c2c4 b4c3 a1a3",
			[]string{"a2a4", "b7b5", "h2h4", "b5b4", "c2c4", "b4c3", "a1a3"},
			0x5c3f9b829b279560,
		},
	}

	for _, position := range tests {
		t.Run(position.name, func(t *testing.T) {
			board := StartPos()
			for _, moveStr := range position.moves {
				board.UCIMakeMove(moveStr)
			}
			if board.hash != position.expected {
				t.Errorf("Hash was 0x%x instead of 0x%x",
					board.hash,
					position.expected,
				)
			}
		})
	}
}
