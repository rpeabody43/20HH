package search

import (
	"20hh/engine/board"
)

func orderMoves(moves []board.Move, ttMove board.Move) {
	ttMoveIdx := -1

	// Extremely basic move ordering at the moment
	// Just find the move stored in the transposition table
	// and put it at the front
	for i, move := range moves {
		if move == ttMove {
			ttMoveIdx = i
			break
		}
	}

	if ttMoveIdx > -1 {
		moves[ttMoveIdx], moves[0] = moves[0], moves[ttMoveIdx]
	}
}
