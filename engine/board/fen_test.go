package board

import (
	"testing"
)

func compareBitboard(t *testing.T, actual, expected Bitboard, name string) {
	if actual != expected {
		t.Errorf(
			"Generated %s bitboard:\n%s\n!= expected:\n%s",
			name,
			actual,
			expected,
		)
	}
}

func testPosBitboards(t *testing.T, whiteExpected, blackExpected Bitboard,
	pieceBitboards [7]Bitboard, board Board) {
	compareBitboard(t, board.colorBitboards[Black], blackExpected, "black pieces")
	compareBitboard(t, board.colorBitboards[White], whiteExpected, "white pieces")

	for piece, expected := range pieceBitboards {
		name := ""
		switch uint8(piece) {
		case Pawn:
			name = "pawn"
		case Knight:
			name = "knight"
		case Bishop:
			name = "bishop"
		case Rook:
			name = "rook"
		case Queen:
			name = "queen"
		case King:
			name = "king"
		}
		actual := board.pieceBitboards[piece]
		compareBitboard(t, actual, expected, name)
	}
}

func testOtherState(t *testing.T, whoseTurn int, castleRights uint8,
	enPassantSq SquareOrNone, halfMoveClock, fullMoves int, board Board) {

	if board.whoseTurn != whoseTurn {
		expectedTurn := "white"
		actualTurn := "black"
		if whoseTurn == Black {
			expectedTurn = "black"
			actualTurn = "white"
		}
		t.Errorf("Instead of %s to start, parsing %s", expectedTurn, actualTurn)
	}

	if board.castleRights != castleRights {
		t.Errorf("castle rights should be %b but is actually %b", castleRights, board.castleRights)
	}

	if board.enPassantSq != enPassantSq {
		t.Errorf("Parsing %d as en passant square instead of %d", board.enPassantSq, enPassantSq)
	}

	if board.halfMoveClock != halfMoveClock {
		t.Errorf("Incorrectly parsing half move clock as %d instead of %d", board.halfMoveClock, halfMoveClock)
	}
	if board.fullMoves != fullMoves {
		t.Errorf("Incorrectly parsing full move clock as %d instead of %d", board.fullMoves, fullMoves)
	}
}

func TestStartPos(t *testing.T) {
	SetupTables()
	board := FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	whitePiecesExpected := Bitboard(0xFFFF)
	blackPiecesExpected := whitePiecesExpected << 48

	var pieceBitboards [King + 1]Bitboard
	pieceBitboards[Pawn] = Bitboard(0xFF<<48 | 0xFF00)
	pieceBitboards[Knight] = Bitboard(0x42<<56 | 0x42)
	pieceBitboards[Bishop] = Bitboard(0x24<<56 | 0x24)
	pieceBitboards[Rook] = Bitboard(0x81<<56 | 0x81)
	pieceBitboards[Queen] = Bitboard(1<<59 | 0x8)
	pieceBitboards[King] = Bitboard(1<<60 | 0x10)

	testPosBitboards(t, whitePiecesExpected, blackPiecesExpected, pieceBitboards, board)
	testOtherState(t, White, 0b1111, NoSq, 0, 1, board)
}

func TestOtherFENs(t *testing.T) {
	board := FromFEN("4k2r/6r1/8/8/8/8/3R4/R3K3 w Qk - 0 1")
	whitePiecesExpected := Bitboard(1 | 1<<11 | 1<<4)
	blackPiecesExpected := Bitboard(1<<60 | 1<<54 | 1<<63)

	var pieceBitboards [7]Bitboard
	pieceBitboards[Rook] = Bitboard(1<<63 | 1<<54 | 1<<11 | 1)
	pieceBitboards[King] = Bitboard(1<<60 | 0x10)

	testPosBitboards(t, whitePiecesExpected, blackPiecesExpected, pieceBitboards, board)
	testOtherState(t, White, Q|k, NoSq, 0, 1, board)
}

func TestWeirdKnightThing(t *testing.T) {
	SetupTables()
	board1 := FromFEN("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1")
	board1.MakeMove(NewMove(G2, F1, KnightPromo|Capture))
	board2 := FromFEN("n1n5/PPPk4/8/8/8/8/4Kp1p/5n1N w - - 0 2")
	compareBitboard(t, board1.pieceBitboards[Knight], board2.pieceBitboards[Knight], "knight")
	compareBitboard(t, board1.colorBitboards[White], board2.colorBitboards[White], "white pieces")
	compareBitboard(t, board1.colorBitboards[Black], board2.colorBitboards[Black], "black pieces")
}

func TestEnPassant(t *testing.T) {
	board := FromFEN("rnbqkbnr/pp1p2pp/8/2p1Pp2/8/1P6/P1P1PPPP/RNBQKBNR w KQkq f6 0 2")
	enPassantSq := board.enPassantSq
	if enPassantSq != SquareOrNone(F6) {
		t.Errorf("En passant set to %d instead of F6", enPassantSq)
	}
}
