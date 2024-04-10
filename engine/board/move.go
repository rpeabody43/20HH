package board

import "fmt"

type Move uint16

// Be able to do bitwise between flags and full moves

// MOVE STRUCTURE :
// (copied from https://www.chessprogramming.org/Encoding_Moves)
//
// flags to     from
// 0000  000000 000000

// FLAGS
const (
	NoFlag      = uint8(0b0000)
	DblPawnMove = uint8(0b0001) // Have to also check starting rank
	Castle      = uint8(0b0010) // Either king or queenside
	QueenCastle = uint8(0b0011)
	Capture     = uint8(0b0100)
	EnPassant   = uint8(0b0101)

	Promotion   = uint8(0b1000)
	KnightPromo = uint8(0b1000)
	BishopPromo = uint8(0b1001)
	RookPromo   = uint8(0b1010)
	QueenPromo  = uint8(0b1011)
)

// MASKS
const (
	SqMask   = 0x3F
	FlagMask = 0xF
)

const NullMove = Move(0)

func NewMove(from, to Square, flag uint8) Move {
	from16 := uint16(from)
	toShifted := uint16(to) << 6
	flagShifted := uint16(flag) << 12
	return Move(from16 | toShifted | flagShifted)
}

func (move Move) GetFrom() Square {
	return Square(move & SqMask)
}

func (move Move) GetTo() Square {
	return Square((move >> 6) & SqMask)
}

func (move Move) GetFlag() uint8 {
	return uint8((move >> 12) & FlagMask)
}

func (move Move) HasFlag(flag uint8) bool {
	return (uint16(move)>>12)&uint16(flag) == uint16(flag)
}

func fileName(sq Square) string {
	file := sq % 8
	return string("abcdefgh"[file])
}

func (move Move) String() string {
	from := move.GetFrom()
	fromRank := from/8 + 1
	fromFile := fileName(from)
	to := move.GetTo()
	toRank := to/8 + 1
	toFile := fileName(to)

	flag := ""
	if move.HasFlag(Promotion) {
		flag = string("nbrq"[move.GetFlag()&0b11])
	}

	return fmt.Sprintf("%s%d%s%d%s", fromFile, fromRank, toFile, toRank, flag)
}
