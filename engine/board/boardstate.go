package board

import (
	"20hh/engine/collections"
)

type Square = uint8

// Used for values that store a single square index
// i.e. en passant square when no en passant is possible
const NoSq = -1

type SquareOrNone = int16

type Piece = uint8

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

	inCheck     bool
	doubleCheck bool
	checkMask   Bitboard

	capturedPieces collections.ArrayStack[Piece]

	rollbacks collections.ArrayStack[Rollback]
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
}

func (board Board) Rollback() Rollback {
	return Rollback{
		castleRights:  board.castleRights,
		inCheck:       board.inCheck,
		doubleCheck:   board.doubleCheck,
		checkMask:     board.checkMask,
		enPassantSq:   board.enPassantSq,
		halfMoveClock: board.halfMoveClock,
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
