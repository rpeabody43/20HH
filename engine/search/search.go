package search

import (
	"20hh/engine/board"

	"math/rand"
)

// Just for testing
func RandomMove(board *board.Board) board.Move {
	moves, movesAmt := board.GenMoves()
	randIdx := rand.Intn(movesAmt)
	return moves[randIdx]
}
