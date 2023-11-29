package board

import (
	"strings"
	"unicode"
)

// Used for values that store a single square index
// i.e. en passant square when no en passant is possible
const NoSq = -1

// Piece indices
const (
	EmptySquare = 0
	Pawn        = 1
	Knight      = 2
	Bishop      = 3
	Rook        = 4
	Queen       = 5
	King        = 6
)

// Color indices
const (
	White = 0
	Black = 1
)

// Castle rights indices
const (
	K = 0
	Q = 1
	k = 2
	q = 3
)

type Board struct {
	pieceTypes  [7]Bitboard
	colorPieces [2]Bitboard

	whoseTurn int

	castleRights [4]bool

	enPassantSq int

	halfMoveClock int // used for 50-move draw rule
	fullMoveClock int
}

func squareIdx(rank, file int) int {
	return rank*8 + file
}

func pieceNumFromLetter(pieceLetter rune) int {
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

func fromFEN(fen string) Board {
	boardState := Board{}

	fields := strings.Fields(fen)
	ranks := strings.Split(fields[0], "/")

	// Piece Placement
	for rank, rankString := range ranks {
		file := 0
		for _, pieceLetter := range rankString {
			if unicode.IsDigit(pieceLetter) {
				emptyPiecesAmt := int(pieceLetter - '0')
				file += emptyPiecesAmt
			} else {
				idx := squareIdx(rank, file)
				piece := pieceNumFromLetter(unicode.ToLower(pieceLetter))
				boardState.pieceTypes[piece].SetSquare(idx)

				color := Black
				if unicode.IsUpper(pieceLetter) {
					color = White
				}
				boardState.colorPieces[color].SetSquare(idx)
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
			castleIdx := 0
			switch letter {
			case 'K':
				castleIdx = K
			case 'Q':
				castleIdx = Q
			case 'k':
				castleIdx = k
			case 'q':
				castleIdx = q
			}
			boardState.castleRights[castleIdx] = true
		}
	}

	// En passant square
	if fields[3] == "-" {
		boardState.enPassantSq = NoSq
	} else {
		squareString := fields[3]
		fileLetter := squareString[0]
		rank := int(squareString[1])
		boardState.enPassantSq = squareIdx(rank, int(fileLetter-'a'))
	}

	// Half & full move clocks
	boardState.halfMoveClock = int(fields[4][0] - '0')
	boardState.fullMoveClock = int(fields[5][0] - '0')

	return boardState
}

func fromStartPos() Board {
	return fromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
}
