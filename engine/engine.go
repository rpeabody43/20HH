package engine

import (
	"20hh/engine/board"
	"20hh/engine/search"
)

type Engine struct {
	currentBoard board.Board
}

func Init() {
	board.SetupTables()
}

func (engine *Engine) GameFromFENString(fen string) {
	engine.currentBoard = board.FromFEN(fen)
}

func (engine *Engine) GameFromStartPos() {
	engine.currentBoard = board.StartPos()
}

func (engine *Engine) PlayMoveFromUCI(moveString string) {
	engine.currentBoard.UCIMakeMove(moveString)
}

func (engine *Engine) GetBestMove() board.Move {
	return search.RandomMove(&engine.currentBoard)
}
