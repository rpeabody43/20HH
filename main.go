package main

import (
	"20hh/engine"
	"strings"

	"bufio"
	"fmt"
	"os"
)

const (
	CYAN  = "\033[36m"
	RESET = "\033[0m"
)

const asciiArt = `
              $$$
         $$$$$$$$$$$$$
     $$$$$$$$$$$$$$$$$$$$$
  $$$$$$$$$       $$$$$$$$$$$    /$$$$$$   /$$$$$$  /$$   /$$ /$$   /$$
  $$$$$$               $$$$$$   /$$__  $$ /$$$_  $$| $$  | $$| $$  | $$
  $$$$$$        $      $$$$$$  |__/  \ $$| $$$$\ $$| $$  | $$| $$  | $$
  $$$$       $$$        $$$$$    /$$$$$$/| $$ $$ $$| $$$$$$$$| $$$$$$$$
  $$$$$   $$$$$         $$$$$   /$$____/ | $$\ $$$$| $$__  $$| $$__  $$
  $$$$$$$$$$$$           $$$$  | $$      | $$ \ $$$| $$  | $$| $$  | $$
  $$$$$$$$$$             $$$$  | $$$$$$$$|  $$$$$$/| $$  | $$| $$  | $$
     $$$$$$$$$$$$$$$$$$$$$     |________/ \______/ |__/  |__/|__/  |__/
         $$$$$$$$$$$$$
              $$$
`

func printBanner() {
	fmt.Print(CYAN)
	fmt.Print(asciiArt)
	fmt.Print(RESET)
	fmt.Println()
	fmt.Println("Enter 'uci' to start or 'quit' to exit")
	fmt.Println()
}

func main() {
	printBanner()
	engine.Init()

	reader := bufio.NewReader(os.Stdin)
	for {
		command, _ := reader.ReadString('\n')
		fields := strings.Fields(command)
		if len(fields) == 0 {
			continue
		}
		command = strings.Fields(command)[0]
		switch command {
		case "uci":
			engine.UCILoop()
			return
		case "quit":
			return
		default:
			fmt.Printf("\"%s\" not recognized\n", command)
		}
	}
}
