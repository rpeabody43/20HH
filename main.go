package main

import (
	"20hh/engine"
	"strings"

	"bufio"
	"fmt"
	"os"
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
	fmt.Print(asciiArt)
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
