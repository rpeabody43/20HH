package engine

import (
	"time"

	"20hh/engine/board"
	"20hh/engine/search"
)

type Engine struct {
	currentBoard board.Board
	search       search.Searcher
}

func Init() {
	board.Init()
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

func (engine *Engine) EndSearch() {
	engine.search.CancelSearch()
}

func (engine *Engine) GetBestMove(
	timeRemaining, timeInc, maxNodes int, infiniteTime bool,
) board.Move {
	moveChan := make(chan board.Move)
	moveTime := timeRemaining/40 + timeInc/2

	// Spawn a search thread
	go engine.search.StartSearch(&engine.currentBoard, moveChan, maxNodes)
	if infiniteTime || maxNodes < int((^uint(0))>>1) {
		// If infiniteTime is true, EndSearch() is called by passing "stop" to UCI
		// or if the engine hits its max nodes
		// Wait for the search thread to finish up
		return <-moveChan
	}
	// wait for it to either end on its own or cancel it when time runs out
	select {
	case bestMove := <-moveChan:
		return bestMove
	case <-time.After(time.Duration(moveTime) * time.Millisecond):
		engine.search.CancelSearch()
		return <-moveChan
	}
}
