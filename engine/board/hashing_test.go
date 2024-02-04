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

func TestIncrementalUpdates(t *testing.T) {
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
