package board

import (
	"fmt"
	"testing"
)

// https://www.chessprogramming.org/Perft_Results
func TestWithPerft(t *testing.T) {
	SetupTables()
	var tests = []struct {
		name     string
		fen      string
		depth    int
		expected int
	}{
		{"Starting pos", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 5, 4865609},
		{"Position 2", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 5, 193690690},
		{"Position 3", "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1", 5, 674624},
		{"Position 4", "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", 5, 15833292},
		{"Position 5", "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8", 5, 89941194},
		{"Position 6", "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10", 5, 164075551},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := FromFEN(tt.fen)
			actual := board.perft(tt.depth, false)
			if actual != tt.expected {
				t.Errorf(
					"In %s perft returned %d instead of %d at depth %d",
					tt.name,
					actual,
					tt.expected,
					tt.depth,
				)
			}
		})
	}
}

func (board *Board) perft(depth int, print bool) int {
	if depth == 0 {
		return 1
	}

	nodes := 0
	moves, movesCount := board.GenMoves()
	for i := 0; i <= movesCount; i++ {
		if print {
			fmt.Println(moves[i])
		}
		if !board.MakeMove(moves[i]) {
			continue
		}
		nodes += board.perft(depth-1, print)
		board.UndoMove(moves[i])
		if print {
			fmt.Println("^")
		}
	}

	return nodes
}

func (board *Board) divide(depth int) int {
	totalNodes := 0
	moves, movesCount := board.GenMoves()
	for i := 0; i <= movesCount; i++ {
		if !board.MakeMove(moves[i]) {
			continue
		}
		fmt.Printf("%s: ", moves[i])
		var moveNodes int
		if moves[i].GetFrom() == E8 && moves[i].GetTo() == C8 {
			moveNodes = board.divide(depth - 1)
		} else {
			moveNodes = board.perft(depth-1, false)
		}
		fmt.Println(moveNodes)
		totalNodes += moveNodes
		board.UndoMove(moves[i])
	}

	return totalNodes
}
