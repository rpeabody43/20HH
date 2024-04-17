package board

// Parses a string fed by UCI to make a move
func (board *Board) UCIMakeMove(moveString string) {
	startFile := uint8(moveString[0] - 'a')
	startRank := uint8(moveString[1] - '1')
	endFile := uint8(moveString[2] - 'a')
	endRank := uint8(moveString[3] - '1')

	startSq := ConvertRankFile(startRank, startFile)
	endSq := ConvertRankFile(endRank, endFile)
	flag := NoFlag

	var delta uint8
	if endSq > startSq {
		delta = endSq - startSq
	} else {
		delta = startSq - endSq
	}
	movingPiece := board.pieces[startSq]
	capturedPiece := board.pieces[endSq]
	if capturedPiece > EmptySquare {
		flag |= Capture
	}

	if movingPiece == Pawn {
		if capturedPiece == 0 && startFile != endFile { // En passant
			flag |= EnPassant
		} else if delta == 16 { // Double pawn move
			flag |= DblPawnMove
		} else if len(moveString) == 5 { // Promotion
			switch moveString[4] {
			case 'n':
				flag |= KnightPromo
			case 'b':
				flag |= BishopPromo
			case 'r':
				flag |= RookPromo
			case 'q':
				flag |= QueenPromo
			}
		}
	}

	if movingPiece == King {
		var fileDelta uint8
		if endFile > startFile {
			fileDelta = endFile - startFile
		} else {
			fileDelta = startFile - endFile
		}

		if fileDelta > 1 {
			flag |= Castle
			if startFile > endFile {
				flag |= QueenCastle
			}
		}
	}

	board.MakeMove(NewMove(startSq, endSq, flag))
}

// Returns false if the move is illegal
func (board *Board) MakeMove(move Move) bool {
	// Need to save this to update hash
	// The position hash only uses the en passant
	// square if there there was a nearby pawn
	hashedEPSq := NoSq
	if board.enPassantPossible() {
		hashedEPSq = board.enPassantSq
	}

	rollback := board.rollback()
	board.rollbacks.Push(rollback)

	whoseTurn := board.whoseTurn

	from := move.GetFrom()
	to := move.GetTo()

	movingPiece := board.pieces[from]
	capturedPiece := board.pieces[to]

	toRank := to / 8
	isEnPassant := (toRank != 0 && toRank != 7) && move.HasFlag(EnPassant)
	promoPiece := EmptySquare // Will remain empty if not a promotion

	movingPieceBB := &board.pieceBitboards[movingPiece]
	capturedPieceBB := &board.pieceBitboards[capturedPiece]

	friendlyBB := &board.colorBitboards[whoseTurn]
	enemyBB := &board.colorBitboards[(whoseTurn+1)%2]

	// Remove the moved piece from its previous position
	movingPieceBB.ClearSquare(from)
	friendlyBB.ClearSquare(from)
	board.pieces[from] = EmptySquare

	if isEnPassant {
		passantedPawnPos := to - 8
		if whoseTurn == Black {
			passantedPawnPos = to + 8
		}
		board.capturedPieces.Push(Pawn)
		board.pieceBitboards[Pawn].ClearSquare(passantedPawnPos)
		enemyBB.ClearSquare(passantedPawnPos)
		board.pieces[passantedPawnPos] = 0
	} else if capturedPiece != EmptySquare {
		board.capturedPieces.Push(capturedPiece)
		capturedPieceBB.ClearSquare(to)
		enemyBB.ClearSquare(to)
	}

	if move.HasFlag(Promotion) {
		// Conveniently, we can chop off everything but the last 2 bits
		// of the flag and it's equal to the promotion piece - 2
		promoPiece = Piece(2 + (move.GetFlag() & 0b11))
		board.pieces[to] = promoPiece
		board.pieceBitboards[promoPiece].SetSquare(to)
	} else {
		board.pieces[to] = movingPiece
		movingPieceBB.SetSquare(to)
	}
	friendlyBB.SetSquare(to)

	if movingPiece == King && move.HasFlag(Castle) {
		board.handleCastling(move)
	} else {
		board.updateCastleRights(from, to)
	}

	if movingPiece == Pawn && !isEnPassant &&
		!move.HasFlag(Promotion) && move.HasFlag(DblPawnMove) {
		board.enPassantSq = SquareOrNone(from - 8)
		if board.whoseTurn == White {
			board.enPassantSq = SquareOrNone(from + 8)
		}
	} else {
		board.enPassantSq = NoSq
	}

	// If move puts king in danger, undo it
	friendlyKingMask := *friendlyBB & board.pieceBitboards[King]
	friendlyKingPos := friendlyKingMask.PopLSB()

	board.halfMoves++
	if whoseTurn == Black {
		board.fullMoves++
	}
	if move.HasFlag(Capture) || movingPiece == Pawn {
		board.halfMoveClock = 0
	} else {
		board.halfMoveClock++
	}

	// King move safety is checked in GenMoves
	if movingPiece != King && board.squareAttacked(friendlyKingPos, NoSq) {
		board.swapTurn()
		board.UndoMove(move)
		return false
	}
	board.swapTurn()
	board.updateHash(move, movingPiece, capturedPiece,
		rollback.castleRights, hashedEPSq)
	board.positionHistory[board.halfMoves] = board.hash
	return true
}

func (board *Board) UndoMove(move Move) {
	rollback := board.rollbacks.Pop()

	board.swapTurn()

	from := move.GetFrom()
	to := move.GetTo()

	movedPiece := board.pieces[to]
	movedPieceBB := &board.pieceBitboards[movedPiece]

	fromRank := from / 8
	isEnPassant := (fromRank == 3 || fromRank == 4) && move.HasFlag(EnPassant)

	friendlyBB := &board.colorBitboards[board.whoseTurn]

	// Remove the moved piece from ending square
	board.pieces[to] = 0
	movedPieceBB.ClearSquare(to)
	friendlyBB.ClearSquare(to)

	// Put moved piece on start square
	if move.HasFlag(Promotion) {
		board.pieceBitboards[Pawn].SetSquare(from)
		board.pieces[from] = Pawn
	} else {
		movedPieceBB.SetSquare(from)
		board.pieces[from] = movedPiece
	}
	friendlyBB.SetSquare(from)

	// Replace captured piece
	if move.HasFlag(Capture) {
		capturedPiece := board.capturedPieces.Pop()
		capturedPieceBB := &board.pieceBitboards[capturedPiece]
		enemyBB := &board.colorBitboards[(board.whoseTurn+1)%2]
		if isEnPassant {
			passantedPawnPos := to - 8
			if board.whoseTurn == Black {
				passantedPawnPos = to + 8
			}
			capturedPieceBB.SetSquare(passantedPawnPos)
			board.pieces[passantedPawnPos] = Pawn
			enemyBB.SetSquare(passantedPawnPos)
		} else {
			capturedPieceBB.SetSquare(to)
			board.pieces[to] = capturedPiece
			enemyBB.SetSquare(to)
		}
	}

	if movedPiece == King && move.HasFlag(Castle) {
		queenside := move.HasFlag(QueenCastle)
		rookStart := rookStartIdx(queenside, board.whoseTurn)
		rookEnd := rookCastleIdx(rookStart)

		board.pieces[rookStart] = Rook
		board.pieces[rookEnd] = 0

		board.pieceBitboards[Rook].SetSquare(rookStart)
		friendlyBB.SetSquare(rookStart)
		board.pieceBitboards[Rook].ClearSquare(rookEnd)
		friendlyBB.ClearSquare(rookEnd)
	}

	board.halfMoves--
	if board.whoseTurn == Black {
		board.fullMoves--
	}

	board.castleRights = rollback.castleRights
	board.inCheck = rollback.inCheck
	board.doubleCheck = rollback.doubleCheck
	board.checkMask = rollback.checkMask
	board.enPassantSq = rollback.enPassantSq
	board.halfMoveClock = rollback.halfMoveClock
	board.hash = rollback.hash
}
