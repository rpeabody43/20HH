package board

import (
	"fmt"
	"strings"

	"20hh/engine/util/collections"
)

type Square = uint8

// Used for values that store a single square index
// i.e. en passant square when no en passant is possible
const NoSq = -1

type SquareOrNone = int16

type Piece = uint8

var initialized = false

func Init() {
	if initialized {
		return
	}
	SetupTables()
	ZobristInit()
	initialized = true
}

type Board struct {
	// 1d array of what piece types are where
	pieces [64]Piece

	// Bitboard for each piece type
	pieceBitboards [7]Bitboard
	// Bitboard for each color
	colorBitboards [2]Bitboard

	whoseTurn int

	castleRights uint8

	enPassantSq SquareOrNone

	halfMoveClock int // used for 50-move draw rule
	fullMoves     int
	halfMoves     int

	inCheck     bool
	doubleCheck bool
	checkMask   Bitboard

	capturedPieces collections.ArrayStack[Piece]
	rollbacks      collections.ArrayStack[Rollback]

	hash            uint64
	positionHistory [150]uint64 // used for repetitions
}

// Doesn't set actual board state, just initializes data structures
func NewBoard() Board {
	return Board{
		capturedPieces: collections.NewArrayStack[Piece](30),
		rollbacks:      collections.NewArrayStack[Rollback](100),
	}
}

func ConvertRankFile(rank, file uint8) Square {
	return Square(rank*8 + file)
}

func StartPos() Board {
	return FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
}

type Rollback struct {
	castleRights  uint8
	inCheck       bool
	doubleCheck   bool
	checkMask     Bitboard
	enPassantSq   SquareOrNone
	halfMoveClock int
	hash          uint64
}

func (board Board) rollback() Rollback {
	return Rollback{
		castleRights:  board.castleRights,
		inCheck:       board.inCheck,
		doubleCheck:   board.doubleCheck,
		checkMask:     board.checkMask,
		enPassantSq:   board.enPassantSq,
		halfMoveClock: board.halfMoveClock,
		hash:          board.hash,
	}
}

func (board *Board) InCheck() bool {
	return board.inCheck
}

func (board *Board) ColorBitboards() (Bitboard, Bitboard) {
	return board.colorBitboards[White], board.colorBitboards[Black]
}

func (board *Board) PieceArray() *[64]Piece {
	return &board.pieces
}

func (board *Board) WhiteToMove() bool {
	return board.whoseTurn == White
}

func (board *Board) BlackToMove() bool {
	return board.whoseTurn == Black
}

func (board *Board) HalfMoveClock() int {
	return board.halfMoveClock
}

func (board *Board) TotalHalfMoves() int {
	return board.halfMoves
}

func (board *Board) Hash() uint64 {
	return board.hash
}

func (board *Board) PosAtNthPly(ply int) uint64 {
	return board.positionHistory[ply]
}

func (board *Board) String() string {
	out := ""
	for row := 7; row >= 0; row-- {
		out += fmt.Sprintf("%d ", row+1)
		for col := 0; col < 8; col++ {
			idx := ConvertRankFile(uint8(row), uint8(col))
			pieceAtIdx := board.pieces[idx]
			val := string(" PNBRQK"[pieceAtIdx])
			if board.colorBitboards[Black].QuerySquare(idx) {
				val = strings.ToLower(val)
			}
			out += fmt.Sprintf("[%s]", val)
		}
		out += "\n"
	}
	out += "   A  B  C  D  E  F  G  H "
	return out
}
