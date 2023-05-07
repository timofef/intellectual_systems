package sudoku

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
)

var startFields = make(map[int]struct{})

type Sudoku struct {
	size    int
	subSize int
	field   []uint32
}

func NewSudoku(path string) *Sudoku {
	csvConf, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = csvConf.Close(); err != nil {
			panic(err)
		}
	}()

	reader := csv.NewReader(csvConf)
	reader.Comma = ' '
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	size, _ := strconv.Atoi(data[0][0])
	sudoku := &Sudoku{size: size, subSize: int(math.Sqrt(float64(size))), field: make([]uint32, size*size)}
	data = data[1:]
	for i, row := range data {
		for j, el := range row {
			val, _ := strconv.Atoi(el)
			sudoku.field[i*size+j] = getBinaryFromInt(val, size)
			if val != 0 {
				startFields[i*size+j] = struct{}{}
			}
		}
	}

	return sudoku
}

func (s *Sudoku) Solve() *Sudoku {
	var stack []*Sudoku
	var count int

	stack = append(stack, s)

	for len(stack) != 0 {
		fmt.Print("Opened: ", count)
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		count++

		if curr.heuristic() == 0 {
			fmt.Println()
			return curr
		}

		neighbours := curr.getNeighbours()
		stack = append(stack, neighbours...)
		fmt.Print("\033[1K\r")
	}

	fmt.Println()

	return nil
}

func (s *Sudoku) getNeighbours() []*Sudoku {
	var neighbourhood []*Sudoku
	
	// Get undefined variable with the smallest domain
	var idx int
	var smallestDomain []uint32
	smallestDomainLen := math.MaxInt
	for i := 0; i < len(s.field); i++ {
		if v := getIntFromBinary(s.field[i], s.size); v == 0 {
			d := extractDomain(
				s.horizontalConstraint(i)|s.verticalConstraint(i)|s.blockConstraint(i),
				s.size)
			if len(d) < smallestDomainLen {
				smallestDomain = d
				smallestDomainLen = len(smallestDomain)
				idx = i
			}
		}
	}
	domain := smallestDomain

	// Generate neighbours with forward checking
	for i := 0; i < len(domain); i++ {
		neighbour := &Sudoku{
			size:    s.size,
			subSize: s.subSize,
		}
		neighbour.field = append(neighbour.field, s.field...)
		neighbour.field[idx] = domain[i]

		neighbour.forwardCheck()

		neighbourhood = append(neighbourhood, neighbour)
	}

	return neighbourhood
}

func (s *Sudoku) forwardCheck() {
	for i := 0; i < len(s.field); i++ {
		if v := getIntFromBinary(s.field[i], s.size); v == 0 {
			domain := extractDomain(
				s.horizontalConstraint(i)|s.verticalConstraint(i)|s.blockConstraint(i),
				s.size)
			if len(domain) == 1 {
				s.field[i] = domain[0]
				i = 0
			}
		}
	}
}
