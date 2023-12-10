package engine

import (
	"20hh/engine/board"

	"bufio"
	"fmt"
	"os"
	"strings"
)

type UCIState struct {
	engine Engine
}

func (state *UCIState) newEngineState() {
	state.engine = Engine{}
	state.engine.GameFromStartPos()
}

func (state *UCIState) newEngineStateFromFEN(fen string) {
	state.engine = Engine{}
	state.engine.GameFromFENString(fen)
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
			state.goCommand(line)
		case "stop":
			fmt.Printf("bestmove %s\n", state.engine.GetBestMove())
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
	// TODO handle b/wtime, b/winc, movestogo, nodes, depth, movetime
	waitForStop := false
	for _, field := range strings.Fields(command) {
		if field == "infinite" {
			waitForStop = true
		}
	}
	if !waitForStop {
		fmt.Printf("bestmove %s\n", state.engine.GetBestMove())
	}
}
