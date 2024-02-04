package board

import (
	"strings"
	"unicode"
)

func pieceNumFromLetter(pieceLetter rune) Piece {
	switch pieceLetter {
	case 'p':
		return Pawn
	case 'n':
		return Knight
	case 'b':
		return Bishop
	case 'r':
		return Rook
	case 'q':
		return Queen
	case 'k':
		return King
	default:
		return EmptySquare
	}
}

func FromFEN(fen string) Board {
	boardState := NewBoard()

	fields := strings.Fields(fen)
	rankStrings := strings.Split(fields[0], "/")

	// Piece Placement
	for rank := 7; rank >= 0; rank-- {
		rankString := rankStrings[7-rank]
		file := uint8(0)
		for _, pieceLetter := range rankString {
			if unicode.IsDigit(pieceLetter) {
				emptyPiecesAmt := uint8(pieceLetter - '0')
				file += emptyPiecesAmt
			} else {
				idx := ConvertRankFile(uint8(rank), file)
				piece := pieceNumFromLetter(unicode.ToLower(pieceLetter))
				boardState.pieceBitboards[piece].SetSquare(idx)
				boardState.pieces[idx] = piece

				color := Black
				if unicode.IsUpper(pieceLetter) {
					color = White
				}
				boardState.colorBitboards[color].SetSquare(idx)
				file++
			}
		}
	}

	// Whose turn it is
	if fields[1] == "w" {
		boardState.whoseTurn = White
	} else {
		boardState.whoseTurn = Black
	}

	// Castling rights
	if fields[2] != "-" {
		for _, letter := range fields[2] {
			castleMask := uint8(0)
			switch letter {
			case 'K':
				castleMask = K
			case 'Q':
				castleMask = Q
			case 'k':
				castleMask = k
			case 'q':
				castleMask = q
			}
			boardState.castleRights |= castleMask
		}
	}

	// En passant square
	if fields[3] == "-" {
		boardState.enPassantSq = NoSq
	} else {
		squareString := fields[3]
		file := uint8(squareString[0] - 'a')
		rank := uint8(squareString[1] - '1')
		boardState.enPassantSq = SquareOrNone(ConvertRankFile(rank, file))
	}

	// Half & full moves
	boardState.halfMoveClock = int(fields[4][0] - '0')
	boardState.fullMoves = int(fields[5][0] - '0')
	boardState.halfMoves = boardState.fullMoves * 2
	if boardState.whoseTurn == Black {
		boardState.halfMoves++
	}

	boardState.handleCheck()

	boardState.genHash()

	return boardState
}
