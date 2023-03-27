package main

import (
	"fmt"
	"os"
	"src/a_search"
	"src/board"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		println("usage: ./main <path_to_csv>")
		return
	}
	graphPathStart := os.Args[1]

	start := board.NewBoardFromFile(graphPathStart)
	terminal := start.Inverse()

	fmt.Println("Start:")
	start.Print()
	fmt.Println("Terminal:")
	terminal.Print()
	fmt.Println()

	s := time.Now()
	path, dur, ok := a_search.A(*start, *terminal)
	f := time.Since(s)

	if ok {
		fmt.Println("Path:")
		for i := len(path) - 1; i >= 0; i-- {
			path[i].Print()
			fmt.Println()
		}
		fmt.Printf("Len: %d\n", dur)
		fmt.Printf("Elapsed time: %s\n", f)
	} else {
		fmt.Println("Couldn't find path")
	}
}
