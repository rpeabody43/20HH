package board

const (
	whiteQCastleMask = Bitboard(1<<B1) | Bitboard(1<<C1) | Bitboard(1<<D1)
	whiteKCastleMask = Bitboard(1<<F1) | Bitboard(1<<G1)
	blackQCastleMask = Bitboard(1<<B8) | Bitboard(1<<C8) | Bitboard(1<<D8)
	blackKCastleMask = Bitboard(1<<F8) | Bitboard(1<<G8)
)

func (board *Board) addMove(moves *[218]Move, moveIdx *int,
	from, to uint8, flag uint8) {
	// Prevents some (not all) non blocking moves from being added
	if board.checkMask > 0 && Bitboard(1<<to)&board.checkMask == 0 {
		return
	}
	moves[*moveIdx] = NewMove(from, to, flag)
	*moveIdx++
}

// Generates pseudo-legal moves
// MakeMove undoes things like pins that aren't checked here
func (board *Board) GenMoves() ([]Move, int) {
	var moves [218]Move
	moveIdx := 0

	// First, check for check
	board.handleCheck()

	whoseTurn := board.whoseTurn
	friendlyBitboard := board.colorBitboards[whoseTurn]
	enemyBitboard := board.colorBitboards[(whoseTurn+1)%2]
	allPieces := friendlyBitboard | enemyBitboard

	// King moves
	friendlyKingMask := board.pieceBitboards[King] & friendlyBitboard
	friendlyKingSq := friendlyKingMask.PopLSB()
	kingAttacks := KingAttacks[friendlyKingSq]
	kingAttacks &^= friendlyBitboard
	for kingAttacks > 0 {
		to := kingAttacks.PopLSB()
		if board.squareAttacked(to, SquareOrNone(friendlyKingSq)) {
			continue
		}

		flag := NoFlag
		if board.pieces[to] > EmptySquare {
			flag = Capture
		}
		moves[moveIdx] = NewMove(friendlyKingSq, to, flag)
		moveIdx++
	}
	// when in double check, only king moves are allowed
	if board.doubleCheck {
		return moves[:moveIdx-1], moveIdx - 1
	}

	board.genCastleMoves(&moves, &moveIdx, friendlyKingSq, allPieces)

	board.genPawnMoves(&moves, &moveIdx)

	// Knight moves
	friendlyKnights := board.pieceBitboards[Knight] & friendlyBitboard
	for friendlyKnights > 0 {
		from := friendlyKnights.PopLSB()
		attacks := KnightAttacks[from]
		attacks &^= friendlyBitboard
		for attacks > 0 {
			to := attacks.PopLSB()
			flag := NoFlag
			if board.pieces[to] > EmptySquare {
				flag = Capture
			}
			board.addMove(&moves, &moveIdx, from, to, flag)
		}
	}

	// Slider Moves

	friendlyOrthoPieces := (board.pieceBitboards[Rook] |
		board.pieceBitboards[Queen]) & friendlyBitboard
	for friendlyOrthoPieces > 0 {
		from := friendlyOrthoPieces.PopLSB()
		attacks := rookAttackBitboard(from, allPieces)
		attacks &^= friendlyBitboard
		for attacks > 0 {
			to := attacks.PopLSB()
			flag := NoFlag
			if board.pieces[to] > EmptySquare {
				flag = Capture
			}
			board.addMove(&moves, &moveIdx, from, to, flag)
		}
	}

	friendlyDiagPieces := (board.pieceBitboards[Bishop] |
		board.pieceBitboards[Queen]) & friendlyBitboard
	for friendlyDiagPieces > 0 {
		from := friendlyDiagPieces.PopLSB()
		attacks := bishopAttackBitboard(from, allPieces)
		attacks &^= friendlyBitboard
		for attacks > 0 {
			to := attacks.PopLSB()
			flag := NoFlag
			if board.pieces[to] > EmptySquare {
				flag = Capture
			}
			board.addMove(&moves, &moveIdx, from, to, flag)
		}
	}

	return moves[:moveIdx-1], moveIdx - 1
}

func (board *Board) genCastleMoves(moves *[218]Move, moveIdx *int,
	friendlyKingSq Square, allPieces Bitboard) {
	if board.castleRights == 0 || board.inCheck {
		return
	}
	var hasRight bool
	var isClear bool
	var isSafe bool
	whoseTurn := board.whoseTurn
	if whoseTurn == White {
		hasRight = board.castleRights&K > 0
		isClear = whiteKCastleMask&allPieces == 0
		isSafe = !board.squareAttacked(G1, NoSq) && !board.squareAttacked(F1, NoSq)
		if hasRight && isClear && isSafe {
			moves[*moveIdx] = NewMove(friendlyKingSq, G1, Castle)
			*moveIdx++
		}
		hasRight = board.castleRights&Q > 0
		isClear = whiteQCastleMask&allPieces == 0
		isSafe = !board.squareAttacked(C1, NoSq) && !board.squareAttacked(D1, NoSq)
		if hasRight && isClear && isSafe {
			moves[*moveIdx] = NewMove(friendlyKingSq, C1, QueenCastle)
			*moveIdx++
		}
	} else {
		hasRight = board.castleRights&k > 0
		isClear = blackKCastleMask&allPieces == 0
		isSafe = !board.squareAttacked(G8, NoSq) && !board.squareAttacked(F8, NoSq)
		if hasRight && isClear && isSafe {
			moves[*moveIdx] = NewMove(friendlyKingSq, G8, Castle)
			*moveIdx++
		}
		hasRight = board.castleRights&q > 0
		isClear = blackQCastleMask&allPieces == 0
		isSafe = !board.squareAttacked(C8, NoSq) && !board.squareAttacked(D8, NoSq)
		if hasRight && isClear && isSafe {
			moves[*moveIdx] = NewMove(friendlyKingSq, C8, QueenCastle)
			*moveIdx++
		}
	}
}

func (board *Board) genPawnMoves(moves *[218]Move, moveIdx *int) {
	whoseTurn := board.whoseTurn
	friendlyBitboard := board.colorBitboards[whoseTurn]
	enemyBitboard := board.colorBitboards[(whoseTurn+1)%2]
	allPieces := friendlyBitboard | enemyBitboard

	promotionFlags := [4]uint8{
		KnightPromo, BishopPromo, RookPromo, QueenPromo,
	}

	friendlyPawns := board.pieceBitboards[Pawn] & friendlyBitboard
	for friendlyPawns > 0 {
		sq := friendlyPawns.PopLSB()
		rank := sq / 8
		promotion := (rank == 6 && whoseTurn == White) ||
			(rank == 1 && whoseTurn == Black)

		quietMoves := PawnQuietMoves[whoseTurn][sq]
		quietMoves &^= enemyBitboard | friendlyBitboard

		attacks := PawnAttacks[whoseTurn][sq]
		enemies := enemyBitboard
		if board.enPassantSq > NoSq {
			enemies |= Bitboard(1) << board.enPassantSq
		}
		attacks &= enemies

		allPawnMoves := quietMoves | attacks
		for allPawnMoves > 0 {
			to := allPawnMoves.PopLSB()

			flag := NoFlag
			moveDelta := int16(to) - int16(sq)
			if moveDelta == 16 || moveDelta == -16 {
				// Prevent pawn from hopping over another piece
				possibleBlockerSq := Square(int16(sq) + moveDelta/2)
				if allPieces.QuerySquare(possibleBlockerSq) {
					continue
				}
				flag = DblPawnMove
			} else if board.pieces[to] > EmptySquare {
				flag = Capture
			} else if SquareOrNone(to) == board.enPassantSq {
				flag = EnPassant
			}

			if !promotion {
				board.addMove(moves, moveIdx, sq, to, flag)
			} else {
				for _, promoFlag := range promotionFlags {
					board.addMove(moves, moveIdx, sq, to, flag|promoFlag)
				}
			}
		}
	}
}

func rookAttackBitboard(sq Square, blockers Bitboard) Bitboard {
	blockers &= RookBlockerMasks[sq]
	index := (uint64(blockers) * RookMagics[sq]) >> RookShifts[sq]
	return RookAttacks[sq][index]
}

func bishopAttackBitboard(sq Square, blockers Bitboard) Bitboard {
	blockers &= BishopBlockerMasks[sq]
	index := (uint64(blockers) * BishopMagics[sq]) >> BishopShifts[sq]
	return BishopAttacks[sq][index]
}

func (board *Board) handleCheck() {
	friendlyKingMask := board.colorBitboards[board.whoseTurn] &
		board.pieceBitboards[King]
	sq := friendlyKingMask.PopLSB()

	board.inCheck = false
	board.doubleCheck = false
	board.checkMask = 0

	oppositeColor := (board.whoseTurn + 1) % 2
	enemyBitboard := board.colorBitboards[oppositeColor]
	blockers := board.colorBitboards[White] | board.colorBitboards[Black]

	rookAttacks := rookAttackBitboard(sq, blockers)
	enemyOrthoPieces := (board.pieceBitboards[Rook] |
		board.pieceBitboards[Queen]) & enemyBitboard
	if enemyOrthoPieces&rookAttacks > 0 {
		board.inCheck = true
		board.checkMask = rookAttacks
	}

	bishopAttacks := bishopAttackBitboard(sq, blockers)
	enemyDiagPieces := (board.pieceBitboards[Bishop] |
		board.pieceBitboards[Queen]) & enemyBitboard
	if enemyDiagPieces&bishopAttacks > 0 {
		if board.inCheck {
			board.doubleCheck = true
			board.checkMask = 0
			return
		} else {
			board.inCheck = true
			board.checkMask = bishopAttacks
		}
	}

	knightAttacks := KnightAttacks[sq]
	enemyKnights := board.pieceBitboards[Knight] & enemyBitboard
	if enemyKnights&knightAttacks > 0 {
		if board.inCheck {
			board.doubleCheck = true
			board.checkMask = 0
		} else {
			board.inCheck = true
		}
	}

	if board.inCheck {
		return
	}

	pawnAttacks := PawnAttacks[board.whoseTurn][sq]
	enemyPawns := board.pieceBitboards[Pawn] & enemyBitboard
	if enemyPawns&pawnAttacks > 0 {
		board.inCheck = true
	}
}
