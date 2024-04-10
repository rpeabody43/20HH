package engine

import (
	"20hh/engine/board"

	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type UCIState struct {
	engine Engine
}

func (state *UCIState) newEngineState() {
	state.engine = Engine{}
	state.engine.GameFromStartPos()
	state.engine.ResetSearch()
}

func (state *UCIState) newEngineStateFromFEN(fen string) {
	state.engine = Engine{}
	state.engine.GameFromFENString(fen)
	state.engine.ResetSearch()
}

func UCILoop() {
	board.SetupTables()
	state := UCIState{}
	initUCI()

	reader := bufio.NewReader(os.Stdin)

	for {
		line, _ := reader.ReadString('\n')
		// fmt.Println(line)
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		switch command := fields[0]; command {
		case "uci":
			initUCI()
		case "isready":
			fmt.Println("readyok")
		case "setoption":
			state.setOption(fields[2:])
		case "ucinewgame":
			state.newEngineState()
		case "position":
			state.positionCommand(line)
		case "go":
			// Spawn a thread so the search doesn't clog the UCI
			go state.goCommand(line)
		case "stop":
			state.engine.EndSearch()
		case "printboard":
			fmt.Println(state.engine.currentBoard.String())
		case "zobrist":
			fmt.Printf("0x%x\n", state.engine.currentBoard.Hash())
		case "quit":
			return
		}
	}
}

func initUCI() {
	fmt.Println("id name 20HH")
	fmt.Println("id author Ryan Peabody")

	// TODO : print options
	fmt.Println("uciok")
}

func (state *UCIState) setOption(fields []string) {
	parsingValue := false
	optionName := ""
	optionVal := ""
	for _, field := range fields {
		if field == "value" {
			parsingValue = true
		} else {
			if parsingValue {
				optionVal += field + " "
			} else {
				optionName += field + " "
			}
		}
	}
	optionName = optionName[:len(optionName)-1]
	optionVal = optionVal[:len(optionVal)-1]
}

func (state *UCIState) positionCommand(command string) {
	command = strings.TrimPrefix(command, "position ")
	if strings.HasPrefix(command, "fen") {
		command = strings.TrimPrefix(command, "fen ")
		fields := strings.Fields(command)
		fenString := strings.Join(fields[:6], " ")
		state.engine.GameFromFENString(fenString)
		command = strings.Join(fields[6:], " ")
	} else if strings.HasPrefix(command, "startpos") {
		command = strings.TrimPrefix(command, "startpos ")
		state.engine.GameFromStartPos()
	}

	if !strings.HasPrefix(command, "moves") {
		return
	}
	command = strings.TrimPrefix(command, "moves ")
	moveStrings := strings.Fields(command)
	for _, moveString := range moveStrings {
		state.engine.PlayMoveFromUCI(moveString)
	}
}

func (state *UCIState) goCommand(command string) {
	// TODO handle  movestogo, depth, movetime
	infinite := false
	fields := strings.Fields(command)
	timeRemaining := 60000 // 1 minute default
	timeInc := 0
	maxNodes := int((^uint(0)) >> 1)
	whiteToMove := state.engine.currentBoard.WhiteToMove()
	for idx, field := range fields {
		if field == "infinite" {
			infinite = true
		} else if field == "nodes" {
			maxNodes, _ = (strconv.Atoi(fields[idx+1]))
		}
		if whiteToMove {
			if field == "wtime" {
				timeRemaining, _ = strconv.Atoi(fields[idx+1])
			} else if field == "winc" {
				timeInc, _ = strconv.Atoi(fields[idx+1])
			}
		} else {
			if field == "btime" {
				timeRemaining, _ = strconv.Atoi(fields[idx+1])
			} else if field == "binc" {
				timeInc, _ = strconv.Atoi(fields[idx+1])
			}
		}
	}
	fmt.Printf("bestmove %s\n",
		state.engine.GetBestMove(timeRemaining, timeInc, maxNodes, infinite),
	)
}
