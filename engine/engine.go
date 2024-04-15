package engine

import (
	"time"

	"20hh/engine/board"
	"20hh/engine/search"
	"20hh/engine/util"
)

// I don't like how this file is a lot of 1 line functions - seems like a lot
// of useless encapsulation but doing it this way makes logical sense

type Engine struct {
	currentBoard board.Board
	search       search.Searcher
	ttSizeMb     uint16
}

func Init() {
	util.RandInit(0xfdfc7fd283aac769)
	board.Init()
}

func (engine *Engine) ResetSearch() {
	engine.ttSizeMb = 1024 // TODO fix with a UCI opt
	engine.search.Init(engine.ttSizeMb)
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

type SearchOpts struct {
	timeRemaining int
	timeInc       int
	maxNodes      int
	infiniteTime  bool
}

func (engine *Engine) GetBestMove(
	opts SearchOpts, loggingCallback search.LogCallback,
) board.Move {
	moveChan := make(chan board.Move)
	moveTime := opts.timeRemaining/40 + opts.timeInc/2

	// Spawn a search thread
	go engine.search.StartSearch(&engine.currentBoard, moveChan,
		loggingCallback, opts.maxNodes)

	// If infinite time is enabled, search stops when UCI tells it to
	// If a node limit is enabled, the search stops when it reaches max nodes
	// Either way wait on the search to return a move
	if opts.infiniteTime || opts.maxNodes < int((^uint(0))>>1) {
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
