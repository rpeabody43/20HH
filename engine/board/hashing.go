package board

import (
	"20hh/engine/util"
)

// Random numbers used for zobrist hashing
var zVals struct {
	pieceSquares   [768]uint64
	castleRights   [16]uint64
    // one for each file + empty
	enPassantFiles [9]uint64 
	blackToMove    uint64
}

func ZobristInit() {
	for i := range zVals.pieceSquares {
		zVals.pieceSquares[i] = util.RandU64()
	}

	for i := range zVals.castleRights {
		zVals.castleRights[i] = util.RandU64()
	}

	for i := range zVals.enPassantFiles {
		zVals.enPassantFiles[i] = util.RandU64()
	}

	zVals.blackToMove = util.RandU64()
}

func (b *Board) genHash() {
	workingHash := uint64(0)

	whiteBB, blackBB := b.ColorBitboards()
	for whiteBB > 0 {
		sq := whiteBB.PopLSB()
		idx := pieceNumIdx(sq, b.pieces[sq], White)
		workingHash ^= zVals.pieceSquares[idx]
	}
	for blackBB > 0 {
		sq := blackBB.PopLSB()
		idx := pieceNumIdx(sq, b.pieces[sq], Black)
		workingHash ^= zVals.pieceSquares[idx]
	}

	workingHash ^= zVals.castleRights[b.castleRights]

	if b.enPassantSq > -1 {
		workingHash ^= zVals.enPassantFiles[b.enPassantSq%8]
	} else {
		workingHash ^= zVals.enPassantFiles[8]
	}

	if b.whoseTurn == Black {
		workingHash ^= zVals.blackToMove
	}

	b.hash = workingHash
}

func (b *Board) updateHash(m Move, moving, captured Piece, oldCastleRights uint8) {
	newHash := b.hash

	if captured > 0 {
		capturedColor := b.whoseTurn
		newHash ^= zVals.pieceSquares[pieceNumIdx(m.GetTo(), captured, capturedColor)]
	}
	movingColor := (b.whoseTurn + 1) % 2
	newHash ^= zVals.pieceSquares[pieceNumIdx(m.GetFrom(), moving, movingColor)]
	newHash ^= zVals.pieceSquares[pieceNumIdx(m.GetTo(), moving, movingColor)]

	newHash ^= zVals.castleRights[oldCastleRights]
	newHash ^= zVals.castleRights[b.castleRights]
	if m.HasFlag(Castle) && !m.HasFlag(Promotion) {
		rookFrom := rookStartIdx(m.HasFlag(QueenCastle), movingColor)
		rookTo := rookCastleIdx(rookFrom)
		newHash ^= zVals.pieceSquares[pieceNumIdx(rookFrom, Rook, movingColor)]
		newHash ^= zVals.pieceSquares[pieceNumIdx(rookTo, Rook, movingColor)]
	}

	if m.HasFlag(EnPassant) && !m.HasFlag(Promotion) {
		file := m.GetTo() % 8
		newHash ^= zVals.enPassantFiles[file]
	} else if b.enPassantSq > -1 {
		file := b.enPassantSq % 8
		newHash ^= zVals.enPassantFiles[file]
	}

	newHash ^= zVals.blackToMove

	b.hash = newHash
}

func pieceNumIdx(sq Square, piece Piece, color int) int {
	return (int(sq)*12 + int(piece-1)*2) + color
}
