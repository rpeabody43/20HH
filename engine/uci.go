package engine

import (
	"20hh/engine/board"
	"20hh/engine/search"

	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func UCILoop() {
	board.SetupTables()
	var engine Engine
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
			setOption(&engine, fields[2:])
		case "ucinewgame":
			engine = Engine{}
			engine.GameFromStartPos()
			engine.ResetSearch()
		case "position":
			positionCommand(&engine, line)
		case "go":
			// Spawn a thread so the search doesn't clog the UCI
			go goCommand(&engine, line)
		case "stop":
			engine.EndSearch()
		case "printboard":
			fmt.Println(engine.currentBoard.String())
		case "zobrist":
			fmt.Printf("0x%x\n", engine.currentBoard.Hash())
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

func setOption(engine *Engine, fields []string) {
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

func positionCommand(engine *Engine, command string) {
	command = strings.TrimPrefix(command, "position ")
	if strings.HasPrefix(command, "fen") {
		command = strings.TrimPrefix(command, "fen ")
		fields := strings.Fields(command)
		fenString := strings.Join(fields[:6], " ")
		engine.GameFromFENString(fenString)
		command = strings.Join(fields[6:], " ")
	} else if strings.HasPrefix(command, "startpos") {
		command = strings.TrimPrefix(command, "startpos ")
		engine.GameFromStartPos()
	}

	if !strings.HasPrefix(command, "moves") {
		return
	}
	command = strings.TrimPrefix(command, "moves ")
	moveStrings := strings.Fields(command)
	for _, moveString := range moveStrings {
		engine.PlayMoveFromUCI(moveString)
	}
}

func goCommand(engine *Engine, command string) {
	// TODO handle movestogo, depth, movetime
	infinite := false
	fields := strings.Fields(command)
	timeRemaining := 60000 // 1 minute default
	timeInc := 0
	maxNodes := int((^uint(0)) >> 1)
	whiteToMove := engine.currentBoard.WhiteToMove()
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

	opts := SearchOpts{
		timeRemaining,
		timeInc,
		maxNodes,
		infinite,
	}
	fmt.Printf("bestmove %s\n",
		engine.GetBestMove(opts, printInfo),
	)
}

// Print incremental updates to UCI
func printInfo(log search.SearchLog) {
	// Format principle variation
	pvString := ""
	for i := 0; i < 15 && log.PV[i] != board.NullMove; i++ {
		pvString += fmt.Sprintf(" %s", log.PV[i])
	}

	// Format score display
	scoreString := fmt.Sprintf("cp %d", log.Score)
	if log.CheckmateScore {
		scoreString = fmt.Sprintf("mate %d", log.Score)
	}

	fmt.Printf(
		"info depth %d score %s nodes %d nps %d time %d hashfull %d pv%s\n",
		log.Depth, scoreString, log.TotalNodes, uint64(log.NPS),
		log.Elapsed, log.TTPermillFull, pvString,
	)
}
