package board

type Bitboard uint64

func (bitboard *Bitboard) SetSquare(idx int) {
	*bitboard |= Bitboard(1 << idx)
}

func (bitboard *Bitboard) ClearSquare(idx int) {
	*bitboard &= Bitboard(^(1 << idx))
}
