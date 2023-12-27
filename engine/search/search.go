package search

import (
	"20hh/engine/board"

	"math/rand"
)

const INFINITY = 100_000
const NEG_INFINITY = -100_000

// Just for testing
func RandomMove(b *board.Board) board.Move {
	moves, movesAmt := b.GenMoves()
	randIdx := rand.Intn(movesAmt)
	return moves[randIdx]
}

func BestMoveInPosition(b *board.Board) board.Move {
	rootMoves, _ := b.GenMoves()

	searchDepth := 5

	bestScore := NEG_INFINITY
	bestMoveIdx := -1
	for moveIdx, move := range rootMoves {
		if !b.MakeMove(move) {
			continue
		}
		score := -search(b, searchDepth-1, -INFINITY, INFINITY)
		b.UndoMove(move)
		if score > bestScore {
			bestScore = score
			bestMoveIdx = moveIdx
		}
	}

	return rootMoves[bestMoveIdx]
}

func search(b *board.Board, depth, alpha, beta int) int {
	if depth == 0 {
		// TODO implement Q Search
		return evalPosition(b)
	}
	allMoves, _ := b.GenMoves()
	legalMoves := 0
	for _, move := range allMoves {
		if !b.MakeMove(move) {
			continue
		}
		legalMoves++
		score := -search(b, depth-1, -beta, -alpha)
		b.UndoMove(move)
		if score >= beta {
			return beta
		} else if score > alpha {
			alpha = score
		}
	}
	if legalMoves == 0 {
		if b.InCheck() {
			return NEG_INFINITY // CHECKMATE
		} else {
			return 0 // STALEMATE
		}
	}
	return alpha
}
