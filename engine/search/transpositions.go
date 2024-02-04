package search

import (
	"20hh/engine/board"
)

const (
	Exact      = 0
	LowerBound = 1
	UpperBound = 2

	EntrySize = 14
)

type TableEntry struct {
	zobrist    uint64
	depth      uint8
	ageAndFlag uint8
	eval       int16
	bestMove   board.Move
}

func (entry TableEntry) GetFlag() uint8 {
	return entry.ageAndFlag & 0b11
}

func (entry TableEntry) GetAge() uint8 {
	return entry.ageAndFlag >> 2
}

type TranspositionTable struct {
	entries []TableEntry
	size    uint64
}

func NewTT(mbSize uint64) TranspositionTable {
	byteSize := mbSize * 1024 * 1024
	size := byteSize / EntrySize
	return TranspositionTable{
		make([]TableEntry, size),
		size,
	}
}
