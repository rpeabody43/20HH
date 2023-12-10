package board

// Piece indices
const (
	EmptySquare = Piece(iota)
	Pawn        = Piece(iota)
	Knight      = Piece(iota)
	Bishop      = Piece(iota)
	Rook        = Piece(iota)
	Queen       = Piece(iota)
	King        = Piece(iota)
)

// Color indices
const (
	White = 0
	Black = 1
)

// Castle rights masks
const (
	K = uint8(0b0001)
	Q = uint8(0b0010)
	k = uint8(0b0100)
	q = uint8(0b1000)
)

// Files
const (
	A = iota
	B
	C
	D
	E
	F
	G
	H
)

const (
	A8, B8, C8, D8, E8, F8, G8, H8 = 56, 57, 58, 59, 60, 61, 62, 63
	A7, B7, C7, D7, E7, F7, G7, H7 = 48, 49, 50, 51, 52, 53, 54, 55
	A6, B6, C6, D6, E6, F6, G6, H6 = 40, 41, 42, 43, 44, 45, 46, 47
	A5, B5, C5, D5, E5, F5, G5, H5 = 32, 33, 34, 35, 36, 37, 38, 39
	A4, B4, C4, D4, E4, F4, G4, H4 = 24, 25, 26, 27, 28, 29, 30, 31
	A3, B3, C3, D3, E3, F3, G3, H3 = 16, 17, 18, 19, 20, 21, 22, 23
	A2, B2, C2, D2, E2, F2, G2, H2 = 8, 9, 1, 11, 12, 13, 14, 15
	A1, B1, C1, D1, E1, F1, G1, H1 = 0, 1, 2, 3, 4, 5, 6, 7
)
