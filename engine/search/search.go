package search

import (
	"fmt"
	"time"

	"20hh/engine/board"
)

const (
	INFINITY       int16 = 10_000
	NEG_INFINITY   int16 = -10_000
	CHECKMATE_EVAL int16 = INFINITY - 1000
)

type Searcher struct {
	searchCancelled    bool
	totalNodesSearched int
	maxNodes           int
	timeSearchingMs    uint64
	timeStarted        time.Time
	searchedOneMove    bool

	tt TranspositionTable
}

func (s *Searcher) Init(ttSizeMb uint16) {
	s.tt = NewTT(ttSizeMb)
}

func (s *Searcher) printInfo(depth uint8, score int16, elapsed uint64, pv *[15]board.Move) {
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
	for i := 0; i < 15 && pv[i] != board.NullMove; i++ {
		pvString += fmt.Sprintf(" %s", pv[i])
	}

	fmt.Printf(
		"info depth %d score %s nodes %d nps %d time %d hashfull %d pv%s\n",
		depth, scoreString, s.totalNodesSearched, uint64(nps), elapsed,
		s.tt.PermillFull(), pvString,
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

	eval := NEG_INFINITY
	pvLine := [15]board.Move{}
	for searchDepth := uint8(1); searchDepth <= 255; searchDepth++ {
		s.searchedOneMove = false

		timeStarted := time.Now()
		evalAtDepth := s.search(b, NEG_INFINITY, INFINITY, searchDepth, 0)
		timeSearchingMs += uint64(time.Since(timeStarted).Milliseconds())

		// This new eval is only good if the search wasn't cancelled
		// before getting through one move
		if s.searchedOneMove || !s.searchCancelled {
			eval = evalAtDepth
		}
		s.tt.UpdatePVLine(b, &pvLine)
		s.printInfo(searchDepth, eval, timeSearchingMs, &pvLine)
		if s.searchCancelled {
			break
		}
		// If a checkmate or there's one legal move stop here
		if eval > CHECKMATE_EVAL || eval < -CHECKMATE_EVAL {
			break
		}
	}

	out <- pvLine[0]
}

func (s *Searcher) search(b *board.Board, alpha, beta int16, depth, ply uint8) int16 {
	if s.searchCancelled {
		return 0
	}
	if IsRepetition(b) {
		return 0
	}
	s.totalNodesSearched++
	if s.totalNodesSearched >= s.maxNodes {
		s.searchCancelled = true
	}
	if depth == 0 {
		return s.qSearch(b, alpha, beta, ply)
	}

	ttEval, ttMove, needsSearch := s.tt.TryGet(
		b.Hash(), ply, depth, alpha, beta,
	)
	if !needsSearch {
		return ttEval
	}

	allMoves, _ := b.GenMoves(false)

	orderMoves(allMoves, ttMove)

	ttFlag := LowerBound

	bestMove := board.NullMove
	legalMoves := 0
	for _, move := range allMoves {
		if !b.MakeMove(move) {
			continue
		}
		score := -s.search(b, -beta, -alpha, depth-1, ply+1)
		b.UndoMove(move)
		legalMoves++

		if s.searchCancelled {
			break
		}

		if score >= beta {
			ttFlag = UpperBound
			s.tt.TryPut(
				b.Hash(), ply, depth,
				uint8(b.TotalHalfMoves()),
				ttFlag, beta, move,
			)
			return beta
		} else if score > alpha {
			ttFlag = Exact
			alpha = score
			bestMove = move
			if ply == 0 {
				s.searchedOneMove = true
			}
		}
	}
	if legalMoves == 0 {
		if b.InCheck() {
			return NEG_INFINITY + int16(ply) // CHECKMATE
		} else {
			return 0 // STALEMATE
		}
	}
	if bestMove != board.NullMove {
		s.tt.TryPut(
			b.Hash(), ply, depth,
			uint8(b.TotalHalfMoves()),
			ttFlag, alpha, bestMove,
		)
	}
	return alpha
}

// Searches until it finds a "quiet" position for a better eval
func (s *Searcher) qSearch(b *board.Board, alpha, beta int16, ply uint8) int16 {
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
		if s.searchCancelled {
			return 0
		}
		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}
	return alpha
}
