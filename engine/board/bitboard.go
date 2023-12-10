package board

import (
	"fmt"
	"math"
	"math/bits"
)

// Based on https://pages.cs.wisc.edu/~psilord/blog/data/chess-pages/rep.html
type Bitboard uint64

func (bb *Bitboard) SetSquare(idx Square) {
	*bb |= Bitboard(1 << idx)
}

func (bb *Bitboard) ClearSquare(idx Square) {
	*bb &^= Bitboard(1 << idx)
}

func (bb *Bitboard) ToggleSquare(idx Square) {
	*bb ^= Bitboard(1 << idx)
}

func (bb Bitboard) QuerySquare(idx Square) bool {
	return bb>>Bitboard(idx)&1 == 1
}

func (bb *Bitboard) PopLSB() uint8 {
	n := uint8(bits.TrailingZeros64(uint64(*bb)))
	*bb &= (*bb - 1)
	return n
}

func (bb Bitboard) RShift(amt uint8) Bitboard {
	return bb >> Bitboard(amt)
}

func (bb Bitboard) LShift(amt uint8) Bitboard {
	return bb << Bitboard(amt)
}

// When passed a negative number, shifts right instead of left
// Nice when generating bishop / rook move masks
func (bb Bitboard) ShiftWithNegative(amt int16) Bitboard {
	shiftAmt := uint8(math.Abs(float64(amt)))
	if amt < 0 {
		return bb >> Bitboard(shiftAmt)
	} else {
		return bb << Bitboard(shiftAmt)
	}
}

// to string for debugging
func (bb Bitboard) String() string {
	out := ""
	for row := 7; row >= 0; row-- {
		out += fmt.Sprintf("%d ", row+1)
		for col := 0; col < 8; col++ {
			idx := row*8 + col
			bit := bb & (1 << idx)
			val := " "
			if bit > 0 {
				val = "X"
			}
			out += fmt.Sprintf("[%s]", val)
		}
		out += "\n"
	}
	out += "   A  B  C  D  E  F  G  H "
	return out
}
