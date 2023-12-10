package board

func rookStartIdx(queenside bool, whoseTurn int) Square {
	rookIdx := H1 // White kingside
	if whoseTurn == Black {
		if queenside {
			rookIdx = A8 // Black queenside
		} else {
			rookIdx = H8 // Black kingside
		}
	} else if queenside {
		rookIdx = A1 // White queenside
	}
	return Square(rookIdx)
}

func rookCastleIdx(startIdx Square) Square {
	ret := startIdx - 2
	if startIdx == A1 || startIdx == A8 {
		ret = startIdx + 3
	}
	return ret
}

func (board *Board) handleCastling(move Move) {
	// Take away both castle rights
	rightsMask := K | Q
	if board.whoseTurn == Black {
		rightsMask = k | q
	}
	board.castleRights &^= rightsMask

	friendlyBB := &board.colorBitboards[board.whoseTurn]

	queenside := move.HasFlag(QueenCastle)
	rookIdx := rookStartIdx(queenside, board.whoseTurn)
	// Castling is a king move, so just move the rook here
	board.pieces[rookIdx] = EmptySquare
	board.pieceBitboards[Rook].ClearSquare(rookIdx)
	friendlyBB.ClearSquare(rookIdx)
	newRookIdx := rookCastleIdx(rookIdx)
	board.pieces[newRookIdx] = Rook
	board.pieceBitboards[Rook].SetSquare(newRookIdx)
	friendlyBB.SetSquare(newRookIdx)
}

func (board *Board) updateCastleRights(from, to Square) {
	if board.castleRights != 0 {
		if to == H1 || from == H1 {
			board.castleRights &^= K
		} else if to == A1 || from == A1 {
			board.castleRights &^= Q
		}

		if to == H8 || from == H8 {
			board.castleRights &^= k
		} else if to == A8 || from == A8 {
			board.castleRights &^= q
		}

		if from == E1 {
			board.castleRights &^= K | Q
		} else if from == E8 {
			board.castleRights &^= k | q
		}
	}
}

func (board *Board) swapTurn() {
	board.whoseTurn = (board.whoseTurn + 1) % 2
}

// similar to handleCheck but returns as soon as possible
func (board *Board) squareAttacked(sq Square, clearSq SquareOrNone) bool {
	friendlyColor := board.whoseTurn
	oppositeColor := (friendlyColor + 1) % 2
	enemyBitboard := board.colorBitboards[oppositeColor]
	blockers := board.colorBitboards[White] | board.colorBitboards[Black]
	if clearSq > NoSq {
		blockers.ClearSquare(Square(clearSq))
	}

	knightAttacks := KnightAttacks[sq]
	enemyKnights := board.pieceBitboards[Knight] & enemyBitboard
	if enemyKnights&knightAttacks > 0 {
		return true
	}

	pawnAttacks := PawnAttacks[friendlyColor][sq]
	enemyPawns := board.pieceBitboards[Pawn] & enemyBitboard
	if enemyPawns&pawnAttacks > 0 {
		return true
	}

	rookAttacks := rookAttackBitboard(sq, blockers)
	enemyOrthoPieces := (board.pieceBitboards[Rook] |
		board.pieceBitboards[Queen]) & enemyBitboard
	if enemyOrthoPieces&rookAttacks > 0 {
		return true
	}

	bishopAttacks := bishopAttackBitboard(sq, blockers)
	enemyDiagPieces := (board.pieceBitboards[Bishop] |
		board.pieceBitboards[Queen]) & enemyBitboard
	if enemyDiagPieces&bishopAttacks > 0 {
		return true
	}

	// This method is used to compute check/pins so this isn't really needed
	kingAttacks := KingAttacks[sq]
	enemyKings := board.pieceBitboards[King] & enemyBitboard
	if enemyKings&kingAttacks > 0 {
		return true
	}

	return false
}
