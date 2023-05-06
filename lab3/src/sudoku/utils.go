package sudoku

import "fmt"

func (s *Sudoku) PrintSudoku(isUnsolved bool) {
	for i := 0; i < s.size; i++ {
		if i%s.subSize == 0 {
			fmt.Print(" ")
			for k := 0; k < (s.size+1)*3+s.subSize+1; k++ {
				fmt.Printf("-")
			}
			fmt.Println()
		}
		for j := 0; j < s.size; j++ {
			if j%s.subSize == 0 {
				fmt.Print(" |")
			}
			n := getIntFromBinary(s.field[i*s.size+j], s.size)
			if _, isStatic := startFields[i*s.size+j]; isStatic {
				fmt.Print("\033[32m") // green
			}
			if isUnsolved {
				if n == 0 {
					fmt.Printf("  *")
				} else {
					fmt.Printf(" %2d", n)
				}
			} else {
				fmt.Printf(" %2d", n)
			}
			if _, isStatic := startFields[i*s.size+j]; isStatic {
				fmt.Print("\033[0m") // green
			}
		}
		fmt.Print(" |")
		fmt.Print("\n")
	}
	fmt.Print(" ")
	for k := 0; k < (s.size+1)*3+s.subSize+1; k++ {
		fmt.Printf("-")
	}
	fmt.Println()
}

func extractDomain(bin uint32, max int) []uint32 {
	var res []uint32

	for i := 0; i < max; i++ {
		isFilled := bin & 1
		if isFilled == 0 {
			res = append(res, getBinaryFromInt(i+1, max))
		}
		bin >>= 1
	}

	return res
}

// Functions to work with binary
func getIntFromBinary(b uint32, max int) int {
	var res int
	for b != 0 {
		res++
		b >>= 1
	}

	return res
}

func getBinaryFromInt(n int, max int) uint32 {
	if n == 0 {
		return 0
	}

	var res uint32 = 1
	for i := 0; i < n-1; i++ {
		res <<= 1
	}

	return res
}

func countZeros(b uint32, max int) int {
	var res int

	for i := 0; i < max; i++ {
		if b&1 != 1 {
			res++
		}
		b >>= 1
	}

	return res
}
