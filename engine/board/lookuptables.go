package board

import (
	"math"
)

var RankMasks [64]Bitboard
var FileMasks [64]Bitboard
var FileMasksNoBorder [64]Bitboard
var RankMasksNoBorder [64]Bitboard

var PawnAttacks [2][64]Bitboard
var PawnQuietMoves [2][64]Bitboard
var KnightAttacks [64]Bitboard
var KingAttacks [64]Bitboard

var RookBlockerMasks [64]Bitboard
var BishopBlockerMasks [64]Bitboard

var RookAttacks [64][]Bitboard
var BishopAttacks [64][]Bitboard

// Cardinal directions
const (
	North = 8
	South = 8
	East  = 1
	West  = 1
)

func SetupTables() {
	// We need to define these two first so we can generate FileMasksNoBorder
	RankMasks[0] = 0xFF << 56
	RankMasks[7] = 0xFF
	for i := uint8(0); i < 8; i++ {
		bb := Bitboard(0)
		for idx := i; idx < 64; idx += 8 {
			bb.SetSquare(idx)
		}
		FileMasks[i] = bb
		FileMasksNoBorder[i] = (bb &^ RankMasks[0]) &^ RankMasks[7]
	}

	for i := uint8(0); i < 8; i++ {
		bb := Bitboard(0)
		for idx := i * 8; idx < (i+1)*8; idx++ {
			bb.SetSquare(idx)
		}
		RankMasks[i] = bb
		RankMasksNoBorder[i] = (bb &^ FileMasks[A]) &^ FileMasks[H]
	}

	for sq := Square(0); sq < 64; sq++ {
		sqBB := Bitboard(0)
		sqBB.SetSquare(sq)

		// Used to prevent wraparounds
		notA := ^FileMasks[A]
		notAB := notA &^ FileMasks[B]
		notH := ^FileMasks[H]
		notGH := notH &^ FileMasks[G]

		// Pawn moves
		whitePawnMove := sqBB << North
		blackPawnMove := sqBB >> South
		PawnQuietMoves[White][sq] = whitePawnMove
		PawnQuietMoves[Black][sq] = blackPawnMove
		// Double moves
		rank := sq / 8
		if rank == 1 {
			PawnQuietMoves[White][sq] |= sqBB << North << North
		}
		if rank == 6 {
			PawnQuietMoves[Black][sq] |= sqBB >> South >> South
		}
		// Pawn attacks
		PawnAttacks[White][sq] = ((whitePawnMove << East) & notA) |
			((whitePawnMove >> West) & notH)
		PawnAttacks[Black][sq] = ((blackPawnMove << East) & notA) |
			((blackPawnMove >> West) & notH)

		// King moves/attacks
		kingMoves := sqBB << North
		kingMoves |= sqBB >> South
		kingMoves |= (sqBB << East) & notA
		kingMoves |= (sqBB >> West) & notH
		kingMoves |= (sqBB << North << East) & notA
		kingMoves |= (sqBB << North >> West) & notH
		kingMoves |= (sqBB >> South << East) & notA
		kingMoves |= (sqBB >> South >> West) & notH

		KingAttacks[sq] = kingMoves

		// Knight moves/attacks
		knightMoves := (sqBB << North << North << East) & notA
		knightMoves |= (sqBB << North << North >> West) & notH
		knightMoves |= (sqBB << North << East << East) & notAB
		knightMoves |= (sqBB << North >> West >> West) & notGH

		knightMoves |= (sqBB >> South >> South << East) & notA
		knightMoves |= (sqBB >> South >> South >> West) & notH
		knightMoves |= (sqBB >> South << East << East) & notAB
		knightMoves |= (sqBB >> South >> West >> West) & notGH

		KnightAttacks[sq] = knightMoves

		// XOR so sq itself isn't included
		RookBlockerMasks[sq] = RankMasksNoBorder[sq/8] ^ FileMasksNoBorder[sq%8]
		BishopBlockerMasks[sq] = bishopBlockerMask(sq)

		/*// Magics and Shifts are stored as literals
		RookMagics[sq], RookAttacks[sq] = findMagicAtSq(sq, true)
		BishopMagics[sq], BishopAttacks[sq] = findMagicAtSq(sq, false)
		*/

		setupRookTable(sq)
		setupBishopTable(sq)
	}
	// printAllMagics()
}

func bishopBlockerMask(sq Square) Bitboard {
	mask := Bitboard(0)
	sqBB := Bitboard(0)
	sqBB.SetSquare(sq)

	shifts := diagShifts()
	distances := diagDistancesFromBorders(sq)

	for i, shift := range shifts {
		distance := distances[i]
		for j := uint8(1); j < distance; j++ {
			mask |= sqBB.ShiftWithNegative(shift * int16(j))
		}
	}

	return mask
}

func getSliderMovesFromBlockers(sq Square, blockers Bitboard,
	shifts [4]int16, distances [4]uint8) Bitboard {
	ret := Bitboard(0)
	sqBB := Bitboard(0)
	sqBB.SetSquare(sq)
	for i, shift := range shifts {
		distance := distances[i]

		for j := uint8(1); j <= distance; j++ {
			newSqBB := sqBB.ShiftWithNegative(shift * int16(j))
			ret |= newSqBB
			// If this is a blocking square don't add any more moves
			if newSqBB&blockers > 0 {
				break
			}
		}
	}
	return ret
}

func rookMovesFromBlockers(sq Square, blockers Bitboard) Bitboard {
	blockers = blockers & RookBlockerMasks[sq]
	distances := distancesFromBorders(sq)
	shifts := [4]int16{North, -South, West, -East}
	return getSliderMovesFromBlockers(sq, blockers, shifts, distances)
}

func bishopMovesFromBlockers(sq Square, blockers Bitboard) Bitboard {
	blockers = blockers & BishopBlockerMasks[sq]
	distances := diagDistancesFromBorders(sq)
	shifts := diagShifts()
	return getSliderMovesFromBlockers(sq, blockers, shifts, distances)
}

func diagShifts() [4]int16 {
	return [4]int16{North + West, North - East, -South + West, -South - East}
}

// Returns an array containing distances from each side of the board
// {top, bottom, right, left}
func distancesFromBorders(sq Square) [4]uint8 {
	rank := uint8(sq / 8)
	file := uint8(sq % 8)
	distanceFromTop := 7 - rank
	distanceFromRight := 7 - file
	return [4]uint8{distanceFromTop, rank, distanceFromRight, file}
}

// Returns an array containing distances from each side of the board,
// as if moving diagonally like a bishop
// {up & left, up & right, down & left, down & right}
func diagDistancesFromBorders(sq Square) [4]uint8 {
	rank := uint8(sq / 8)
	file := uint8(sq % 8)
	distanceFromTop := 7 - rank
	distanceFromRight := 7 - file

	nw := uint8Min(distanceFromTop, distanceFromRight)
	ne := uint8Min(distanceFromTop, file)
	sw := uint8Min(rank, distanceFromRight)
	se := uint8Min(rank, file)
	return [4]uint8{nw, ne, sw, se}
}

func uint8Min(a, b uint8) uint8 {
	return uint8(math.Min(float64(a), float64(b)))
}
