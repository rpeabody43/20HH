package board

import (
	"fmt"
	"testing"
)

func bitboardToString(bb Bitboard) string {
	out := ""
	for idx := 0; idx < 64; idx++ {
		bit := bb & (1 << idx)
		val := " "
		if bit > 0 {
			val = "X"
		}
		out += fmt.Sprintf("[%s]", val)
		if idx > 0 && idx < 64 && (idx+1)%8 == 0 {
			out += "\n"
		}
	}
	return out
}

func compareBitboard(t *testing.T, actual Bitboard, expected Bitboard) {
	if actual != expected {
		t.Errorf(
			"Generated bitboard:\n%s\n!= expected:\n%s",
			bitboardToString(actual),
			bitboardToString(expected),
		)
	}
}

func TestStartPosBitboards(t *testing.T) {
	board := fromStartPos()

	blackPiecesExpected := Bitboard(0xFFFF)
	whitePiecesExpected := blackPiecesExpected << 48

	compareBitboard(t, board.colorPieces[Black], blackPiecesExpected)
	compareBitboard(t, board.colorPieces[White], whitePiecesExpected)

	pieceBitboards := make(map[int]Bitboard)
	pieceBitboards[Pawn] = Bitboard(0xFF<<48 | 0xFF00)
	pieceBitboards[Knight] = Bitboard(0x42<<56 | 0x42)
	pieceBitboards[Bishop] = Bitboard(0x24<<56 | 0x24)
	pieceBitboards[Rook] = Bitboard(0x81<<56 | 0x81)
	pieceBitboards[Queen] = Bitboard(0x8<<56 | 0x8)
	pieceBitboards[King] = Bitboard(0x10<<56 | 0x10)

	for piece, bitboard := range pieceBitboards {
		t.Log("Currently testing", piece)
		compareBitboard(t, board.pieceTypes[piece], bitboard)
	}
}

func TestStartPosOtherState(t *testing.T) {
	board := fromStartPos()

	if board.whoseTurn != White {
		t.Errorf("Instead of white to start, parsing %d", board.whoseTurn)
	}

	for i, castleRight := range board.castleRights {
		if !castleRight {
			t.Errorf("%dth castle right should be true but is false", i)
		}
	}

	if board.enPassantSq != NoSq {
		t.Errorf("Parsing %d as en passant square instead of none", board.enPassantSq)
	}

	if board.halfMoveClock != 0 {
		t.Errorf("Incorrectly parsing half move clock as %d", board.halfMoveClock)
	}
	if board.fullMoveClock != 1 {
		t.Errorf("Incorrectly parsing full move clock as %d", board.fullMoveClock)
	}
}
