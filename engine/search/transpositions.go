package search

import (
	"20hh/engine/board"
)

const (
	Exact      uint8 = 0b01
	LowerBound uint8 = 0b10
	UpperBound uint8 = 0b11

	EntrySize = 14
)

type _TableEntry struct {
	zobrist    uint64
	depth      uint8
	ageAndFlag uint8
	eval       int16
	bestMove   board.Move
}

func (entry _TableEntry) getFlag() uint8 {
	return entry.ageAndFlag & 0b11
}

func (entry _TableEntry) getAge() uint8 {
	return entry.ageAndFlag >> 2
}

type TranspositionTable struct {
	entries   []_TableEntry
	size      uint64
	numFilled uint64
}

func NewTT(mbSize uint16) TranspositionTable {
	byteSize := uint64(mbSize) * 1024 * 1024
	size := byteSize / EntrySize
	return TranspositionTable{
		make([]_TableEntry, size),
		size,
		0,
	}
}

// Tries to find the position in the transposition table
// Returns evaluation, best move, and whether the position
// requires a full search.
func (tt *TranspositionTable) TryGet(
	zobrist uint64, ply, depth uint8, alpha, beta int16,
) (int16, board.Move, bool) {
	idx := zobrist % tt.size
	entry := tt.entries[idx]

	// If this position isn't in the table at all we can't do anything
	if entry.zobrist != zobrist {
		return 0, board.NullMove, true
	}

	bestMove := entry.bestMove
	eval := entry.eval
	if eval > CHECKMATE_EVAL {
		eval -= int16(ply)
	} else if eval < -CHECKMATE_EVAL {
		eval += int16(ply)
	}

	// If the saved depth is less than the search we're about to do,
	// the stored move will be useful but we still need to search
	if entry.depth < depth {
		return eval, bestMove, true
	}

	requiresSearch := true
	flag := entry.getFlag()
	if flag == Exact {
		requiresSearch = false
	} else if flag == UpperBound && eval >= beta {
		eval = beta
		requiresSearch = false
	} else if flag == LowerBound && eval <= alpha {
		eval = alpha
		requiresSearch = false
	}

	return eval, bestMove, requiresSearch
}

func (tt *TranspositionTable) TryPut(
	zobrist uint64,
	ply, depth, halfMoves, flag uint8,
	eval int16,
	bestMove board.Move,
) {
	replace := false
	idx := zobrist % tt.size
	if tt.entries[idx].zobrist == 0 {
		replace = true
		tt.numFilled++
	} else if depth > tt.entries[idx].depth {
		replace = true
	} else if halfMoves-tt.entries[idx].getAge() > 10 {
		replace = true
	}
	if eval > CHECKMATE_EVAL {
		eval += int16(ply)
	} else if eval < -CHECKMATE_EVAL {
		eval -= int16(ply)
	}
	if replace {
		tt.entries[idx] = _TableEntry{
			zobrist,
			depth,
			(halfMoves << 2) | flag,
			eval,
			bestMove,
		}
	}
}

func (tt *TranspositionTable) UpdatePVLine(b *board.Board, pvLine *[15]board.Move) {
	end := 0

	positionHash := b.Hash()
	idx := b.Hash() % tt.size
	moveFound := tt.entries[idx].zobrist == positionHash
	for moveFound && end < 15 {
		//fmt.Println(positionHash)
		move := tt.entries[idx].bestMove
		pvLine[end] = move
		end++
		b.MakeMove(move)

		positionHash = b.Hash()
		idx = positionHash % tt.size
		moveFound = tt.entries[idx].zobrist == positionHash
	}

	// Just in case the pv line is cut short and has extra moves left over
	// from a previous iteration
	for i := end + 1; i < 15 && pvLine[i] != board.NullMove; i++ {
		pvLine[i] = board.NullMove
	}

	for ; end > 0; end-- {
		b.UndoMove(pvLine[end-1])
	}
}

func (tt *TranspositionTable) PermillFull() uint16 {
	return uint16(1000 * float32(tt.numFilled) / float32(tt.size))
}
