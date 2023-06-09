package sudoku

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
)

const SEED = 123

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
	var isFilled bool
	for i, row := range data {
		for j, el := range row {
			val, _ := strconv.Atoi(el)
			if val != 0 {
				isFilled = true
			} else {
				isFilled = false
			}
			sudoku.field[i*size+j] = getBinaryFromInt(val, isFilled, size)
		}
	}

	return sudoku
}

func (s *Sudoku) Solve() *Sudoku {
	r := rand.New(rand.NewSource(SEED))

	// Fill sub-grids
	s.initField()

	h := s.heuristic()
	neighbourBlocks := 1
	currMethod := 0
	methods := []func(int, int, int) *Sudoku{s.insert, s.swap, s.megaswap}

	shackesNum := 0

	for h != 0 {
		fmt.Printf("Heuristic: %3d     Blocks: %d", h, neighbourBlocks)

		var best *Sudoku
		oldH := h
		for i := 0; i < neighbourBlocks; i++ {
			s1 := methods[currMethod](i%s.subSize, i/s.subSize, h)
			if h1 := s1.heuristic(); h1 < h {
				best = s1
				h = h1
			}
		}

		if h == oldH {
			neighbourBlocks++
			if neighbourBlocks > s.size {
				neighbourBlocks = s.size

				currMethod++
				if currMethod == len(methods) {
					currMethod = 0
					neighbourBlocks = 1
					s.shake(r)
					shackesNum++
				}
			}
		} else {
			s.field = best.field
		}
		h = s.heuristic()
		fmt.Print("\033[1K\r")
	}

	fmt.Printf("Shakes: %d\n", shackesNum)

	return s
}

func (s *Sudoku) shake(r *rand.Rand) {
	for i := 0; i < s.subSize; i++ {
		for j := 0; j < s.subSize; j++ {
			alreadyInserted := make(map[int]struct{})
			for k := 0; k < s.subSize; k++ {
				for l := 0; l < s.subSize; l++ {
					if isStatic(s.field[i*s.subSize*s.size+k*s.size+j*s.subSize+l], s.size) {
						alreadyInserted[getIntFromBinary(s.field[i*s.subSize*s.size+k*s.size+j*s.subSize+l], s.size)] = struct{}{}
					}
				}
			}
			//count := 1
			for k := 0; k < s.subSize; k++ {
				for l := 0; l < s.subSize; l++ {

					if isStatic(s.field[i*s.subSize*s.size+k*s.size+j*s.subSize+l], s.size) {
						continue
					}
					for {
						v := r.Int()%s.size + 1
						_, exists := alreadyInserted[v]
						if !exists {
							s.field[i*s.subSize*s.size+k*s.size+j*s.subSize+l] = getBinaryFromInt(v, false, s.size)
							alreadyInserted[v] = struct{}{}
							break
						}
					}
				}
			}
		}
	}
}

// Invert
func (s *Sudoku) invert(i, j int) *Sudoku {
	// Get copy of field
	res := &Sudoku{size: s.size, subSize: s.subSize, field: make([]uint32, 0)}
	res.field = append(res.field, s.field...)

	// Get indexes of non-fixed elements
	a := make([]struct {
		k int
		l int
	}, 0)
	for k := 0; k < s.subSize; k++ {
		for l := 0; l < s.subSize; l++ {
			if !isStatic(s.field[i*s.subSize*s.size+k*s.size+j*s.subSize+l], s.size) {
				a = append(a, struct {
					k int
					l int
				}{k: k, l: l})
			}
		}
	}

	start := 0
	finish := len(a) - 1

	for start < finish {
		sk, sl := a[start].k, a[start].l
		fk, fl := a[finish].k, a[finish].l
		tmp := res.field[i*s.subSize*s.size+sk*s.size+j*s.subSize+sl]
		res.field[i*s.subSize*s.size+sk*s.size+j*s.subSize+sl] = res.field[i*s.subSize*s.size+fk*s.size+j*s.subSize+fl]
		res.field[i*s.subSize*s.size+fk*s.size+j*s.subSize+fl] = tmp
		start++
		finish--
	}

	return res
}

func (s *Sudoku) insert(i, j int, target int) *Sudoku {
	a := make([]struct {
		k int
		l int
	}, 0)
	for k := 0; k < s.subSize; k++ {
		for l := 0; l < s.subSize; l++ {
			if !isStatic(s.field[i*s.subSize*s.size+k*s.size+j*s.subSize+l], s.size) {
				a = append(a, struct {
					k int
					l int
				}{k: k, l: l})
			}
		}
	}

	res := &Sudoku{size: s.size, subSize: s.subSize}
	res.field = append(res.field, s.field...)
	tmp := &Sudoku{size: s.size, subSize: s.subSize}

	for m := 0; m < len(a)-1; m++ {
		for n := m + 1; n < len(a); n++ {
			tmp.field = make([]uint32, 0)
			tmp.field = append(tmp.field, res.field...)
			start, finish := m, n
			for start < finish {
				sk, sl := a[start].k, a[start].l
				fk, fl := a[finish].k, a[finish].l
				t := tmp.field[i*s.subSize*s.size+sk*s.size+j*s.subSize+sl]
				tmp.field[i*s.subSize*s.size+sk*s.size+j*s.subSize+sl] =
					tmp.field[i*s.subSize*s.size+fk*s.size+j*s.subSize+fl]
				tmp.field[i*s.subSize*s.size+fk*s.size+j*s.subSize+fl] = t
				start++
				finish--
			}
			if tmp.heuristic() < target {
				res.field = tmp.field
			}
		}
	}

	return res
}

func (s *Sudoku) swap(i, j int, target int) *Sudoku {
	a := make([]struct {
		k int
		l int
	}, 0)
	for k := 0; k < s.subSize; k++ {
		for l := 0; l < s.subSize; l++ {
			if !isStatic(s.field[i*s.subSize*s.size+k*s.size+j*s.subSize+l], s.size) {
				a = append(a, struct {
					k int
					l int
				}{k: k, l: l})
			}
		}
	}

	// Create copy of field
	res := &Sudoku{size: s.size, subSize: s.subSize}
	res.field = append(res.field, s.field...)
	tmp := &Sudoku{size: s.size, subSize: s.subSize}

	for m := 0; m < len(a)-1; m++ {
		for n := m + 1; n < len(a); n++ {
			tmp.field = make([]uint32, 0)
			tmp.field = append(tmp.field, s.field...)
			start, finish := m, n
			sk, sl := a[start].k, a[start].l
			fk, fl := a[finish].k, a[finish].l
			t := tmp.field[i*s.subSize*s.size+sk*s.size+j*s.subSize+sl]
			tmp.field[i*s.subSize*s.size+sk*s.size+j*s.subSize+sl] =
				tmp.field[i*s.subSize*s.size+fk*s.size+j*s.subSize+fl]
			tmp.field[i*s.subSize*s.size+fk*s.size+j*s.subSize+fl] = t
			if tmp.heuristic() < target {
				res.field = tmp.field
			}
		}
	}

	return res
}

func (s *Sudoku) megaswap(i, j int, target int) *Sudoku {
	a := make([]struct {
		k int
		l int
	}, 0)
	for k := 0; k < s.subSize; k++ {
		for l := 0; l < s.subSize; l++ {
			if !isStatic(s.field[i*s.subSize*s.size+k*s.size+j*s.subSize+l], s.size) {
				a = append(a, struct {
					k int
					l int
				}{k: k, l: l})
			}
		}
	}

	res := &Sudoku{size: s.size, subSize: s.subSize}
	res.field = append(res.field, s.field...)
	tmp := &Sudoku{size: s.size, subSize: s.subSize}

	for m := 1; m < len(a)-1; m++ {
		tmp.field = make([]uint32, 0)
		tmp.field = append(tmp.field, s.field...)
		start, finish := m-1, m+1
		for start >= 0 && finish <= len(a)-1 {
			sk, sl := a[start].k, a[start].l
			fk, fl := a[finish].k, a[finish].l
			t := tmp.field[i*s.subSize*s.size+sk*s.size+j*s.subSize+sl]
			tmp.field[i*s.subSize*s.size+sk*s.size+j*s.subSize+sl] =
				tmp.field[i*s.subSize*s.size+fk*s.size+j*s.subSize+fl]
			tmp.field[i*s.subSize*s.size+fk*s.size+j*s.subSize+fl] = t
			start--
			finish++
		}
		if tmp.heuristic() < target {
			res.field = tmp.field
		}
	}

	return res
}

// Fill empty spaces to satisfy sub-grid constraint
func (s *Sudoku) initField() {
	for i := 0; i < s.subSize; i++ {
		for j := 0; j < s.subSize; j++ {
			alreadyInserted := make(map[int]struct{})
			for k := 0; k < s.subSize; k++ {
				for l := 0; l < s.subSize; l++ {
					v := getIntFromBinary(s.field[i*s.subSize*s.size+k*s.size+j*s.subSize+l], s.size)
					if v != 0 {
						alreadyInserted[v] = struct{}{}
					}
				}
			}
			count := 1
			for k := 0; k < s.subSize; k++ {
				for l := 0; l < s.subSize; l++ {

					v := getIntFromBinary(s.field[i*s.subSize*s.size+k*s.size+j*s.subSize+l], s.size)
					if v != 0 {
						continue
					}
					for {
						_, exists := alreadyInserted[count]
						if !exists {
							s.field[i*s.subSize*s.size+k*s.size+j*s.subSize+l] = getBinaryFromInt(count, false, s.size)
							alreadyInserted[count] = struct{}{}
							break
						}
						count++
					}
				}
			}
		}
	}
}

func (s *Sudoku) heuristic() int {
	var res int
	var mask uint32
	for i := 0; i < s.size; i++ {
		mask <<= 1
		mask++
	}

	// horizontal
	for i := 0; i < s.size; i++ {
		var heuristic uint32
		for j := 0; j < s.size; j++ {
			heuristic |= s.field[i*s.size+j]
		}
		res += countZeros(heuristic, mask)
	}

	// vertical
	for j := 0; j < s.size; j++ {
		var heuristic uint32
		for i := 0; i < s.size; i++ {
			heuristic |= s.field[i*s.size+j]
		}
		res += countZeros(heuristic, mask)
	}

	return res
}

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
			if isStatic(s.field[i*s.size+j], s.size) {
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
			if isStatic(s.field[i*s.size+j], s.size) {
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

// Functions to work with binary
func getIntFromBinary(b uint32, max int) int {
	b &= ^(1 << max)
	var res int
	for b != 0 {
		res++
		b >>= 1
	}

	return res
}

func getBinaryFromInt(n int, isStatic bool, max int) uint32 {
	if n == 0 {
		return 0
	}

	var res uint32 = 1
	for i := 0; i < n-1; i++ {
		res <<= 1
	}
	if isStatic {
		res |= 1 << max
	}

	return res
}

func isStatic(b uint32, max int) bool {
	return b&(1<<max) != 0
}

func countZeros(b uint32, mask uint32) int {
	var res int
	b &= mask
	for b != 0 {
		if b&1 != 1 {
			res++
		}
		b >>= 1
	}

	return res
}
