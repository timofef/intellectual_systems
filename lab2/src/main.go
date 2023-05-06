package main

import (
	"fmt"
	"os"
	"src/sudoku"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		println("usage: ./main <path_to_csv>")
		return
	}
	s := sudoku.NewSudoku(os.Args[1])
	fmt.Print("Unsolved sudoku:\n")
	s.PrintSudoku(true)

	start := time.Now()
	solution := s.Solve()
	finish := time.Since(start)

	fmt.Println("Time elapsed: ", finish)
	fmt.Print("Solved sudoku:\n")
	solution.PrintSudoku(false)
}
