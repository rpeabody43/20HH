package search

import (
	"fmt"
	"time"

	"20hh/engine/board"
)

const (
	INFINITY       = 100_000
	NEG_INFINITY   = -100_000
	CHECKMATE_EVAL = INFINITY - 1000
)

type Searcher struct {
	searchCancelled    bool
	totalNodesSearched int
	maxNodes           int
	timeSearchingMs    uint64
	timeStarted        time.Time
	searchedOneMove    bool
}

type PVNode struct {
	move    board.Move
	moveIdx int
}

func (s *Searcher) printInfo(depth, score int, elapsed uint64, pv *[50]PVNode) {
	nps := float64(s.totalNodesSearched*1000) / float64(elapsed)
	if elapsed == 0 {
		nps = float64(s.totalNodesSearched)
	}

	// Format score
	scoreString := fmt.Sprintf("cp %d", score)
	if score > CHECKMATE_EVAL {
		plies := INFINITY - score
		mateIn := (plies / 2) + (plies % 2)
		scoreString = fmt.Sprintf("mate %d", mateIn)
	} else if score < -CHECKMATE_EVAL {
		plies := NEG_INFINITY - score
		mateIn := (plies / 2) + (plies % 2)
		scoreString = fmt.Sprintf("mate %d", mateIn)
	}

	// Format principle variation
	pvString := ""
	for i := 0; pv[i].move != 0; i++ {
		pvString += fmt.Sprintf(" %s", pv[i].move)
	}

	fmt.Printf(
		"info depth %d score %s nodes %d nps %d time %d pv%s\n",
		depth, scoreString, s.totalNodesSearched, uint64(nps), elapsed, pvString,
	)
}

func (s *Searcher) CancelSearch() {
	s.searchCancelled = true
}

func (s *Searcher) StartSearch(b *board.Board, out chan board.Move, maxNodes int) {
	s.searchCancelled = false
	s.totalNodesSearched = 0
	s.maxNodes = maxNodes
	timeSearchingMs := uint64(0)

	pv := [50]PVNode{}

	eval := NEG_INFINITY
	for searchDepth := 1; searchDepth < 256; searchDepth++ {
		s.searchedOneMove = false

		timeStarted := time.Now()
		evalAtDepth := s.search(b, &pv, searchDepth, -INFINITY, INFINITY, 0)
		timeSearchingMs += uint64(time.Since(timeStarted).Milliseconds())

		// This new eval is only good if the search wasn't cancelled
		// before getting through one move
		if s.searchedOneMove {
			eval = evalAtDepth
		}
		s.printInfo(searchDepth, eval, timeSearchingMs, &pv)
		if s.searchCancelled {
			break
		}
		// If a checkmate or there's one legal move stop here
		if eval > CHECKMATE_EVAL || eval < -CHECKMATE_EVAL {
			break
		}
	}

	out <- pv[0].move
}

func (s *Searcher) search(b *board.Board, pv *[50]PVNode, depth, alpha, beta, ply int) int {
	if s.searchCancelled {
		return 0
	}
	s.totalNodesSearched++
	if s.totalNodesSearched >= s.maxNodes {
		s.searchCancelled = true
	}
	if depth == 0 {
		return s.qSearch(b, alpha, beta, ply)
	}

	allMoves, _ := b.GenMoves(false)
	prevBestIdx := -1
	// If we haven't fully searched a move, then we're still on the pv line
	if !s.searchedOneMove {
		prevBestIdx = pv[ply].moveIdx
	}
	orderMoves(allMoves, prevBestIdx)

	legalMoves := 0
	for moveIdx, move := range allMoves {
		if !b.MakeMove(move) {
			continue
		}
		score := -s.search(b, pv, depth-1, -beta, -alpha, ply+1)
		b.UndoMove(move)
		legalMoves++

		if score >= beta {
			return beta
		} else if score > alpha {
			alpha = score
			s.searchedOneMove = true
			newPVNode := PVNode{
				move,
				moveIdx,
			}
			pv[ply] = newPVNode
		}
	}
	if legalMoves == 0 {
		if b.InCheck() {
			return NEG_INFINITY + ply // CHECKMATE
		} else {
			return 0 // STALEMATE
		}
	}
	return alpha
}

// Searches until it finds a "quiet" position for a better eval
func (s *Searcher) qSearch(b *board.Board, alpha, beta, ply int) int {
	if s.searchCancelled {
		return 0
	}
	s.totalNodesSearched++
	if s.totalNodesSearched >= s.maxNodes {
		s.searchCancelled = true
	}

	// eval anyway in case it's a bad capture
	eval := evalPosition(b)
	if eval >= beta {
		return beta
	}
	if eval > alpha {
		alpha = eval
	}

	quietMoves, _ := b.GenMoves(true)
	for _, move := range quietMoves {
		if !b.MakeMove(move) {
			continue
		}
		score := -s.qSearch(b, -beta, -alpha, ply+1)
		b.UndoMove(move)
		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}
	return alpha
}
