package board

import (
	"testing"
)

var initialized = false

func initialize() {
	if !initialized {
		SetupTables()
		initialized = true
	}
}

func TestBishopMoveMask(t *testing.T) {
	initialize()
	testD4 := BishopBlockerMasks[D4]
	expectedD4 := Bitboard(0x40221400142200)
	compareBitboard(t, testD4, expectedD4, "bishop D4 mask")

	testF1 := BishopBlockerMasks[F1]
	expectedF1 := Bitboard(0x204085000)
	compareBitboard(t, testF1, expectedF1, "bishop F1 mask")
}

func TestRookMoveMask(t *testing.T) {
	initialize()

	testD4 := RookBlockerMasks[D4]
	expectedD4 := Bitboard(0x8080876080800)
	compareBitboard(t, testD4, expectedD4, "rook D4 mask")

	testA1 := RookBlockerMasks[A1]
	expectedA1 := Bitboard(0x101010101017E)
	compareBitboard(t, testA1, expectedA1, "rook B2 mask")
}

func TestRookMovesFromBlockers(t *testing.T) {
	initialize()

	d4Blockers := Bitboard(0x808000082080000)
	testD4 := rookMovesFromBlockers(D4, d4Blockers)
	expectedD4 := Bitboard(0x80808F6080000)
	compareBitboard(t, testD4, expectedD4, "rook D4 moves")

	a1Blockers := Bitboard(0xFFFF00000000FFFF)
	testA1 := rookMovesFromBlockers(A1, a1Blockers)
	expectedA1 := Bitboard(0x102)
	compareBitboard(t, testA1, expectedA1, "rook A1 moves")
}

func TestBishopMovesFromBlockers(t *testing.T) {
	initialize()

	d4Blockers := Bitboard(0x41000000042000)
	testD4 := bishopMovesFromBlockers(D4, d4Blockers)
	expectedD4 := Bitboard(0x41221400142000)
	compareBitboard(t, testD4, expectedD4, "bishop D4 moves")

	c1Blockers := Bitboard(0xFFFF00000000FFFF)
	testC1 := bishopMovesFromBlockers(C1, c1Blockers)
	expectedC1 := Bitboard(0xA00)
	compareBitboard(t, testC1, expectedC1, "bishop C1 moves")
}

func TestRookMagics(t *testing.T) {
	initialize()

	d4Blockers := Bitboard(0x808000082080000) & RookBlockerMasks[D4]
	index := (uint64(d4Blockers) * RookMagics[D4]) >> RookShifts[D4]
	testD4 := RookAttacks[D4][index]
	expectedD4 := Bitboard(0x80808F6080000)
	compareBitboard(t, testD4, expectedD4, "rook D4 moves")
}

func TestBishopMagics(t *testing.T) {
	initialize()

	d4Blockers := Bitboard(0x41000000042000) & BishopBlockerMasks[D4]
	index := (uint64(d4Blockers) * BishopMagics[D4]) >> BishopShifts[D4]
	testD4 := BishopAttacks[D4][index]
	expectedD4 := Bitboard(0x41221400142000)
	compareBitboard(t, testD4, expectedD4, "bishop D4 moves")
}
