package board

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	FREE      rune = '0'
	BLACK     rune = '1'
	WHITE     rune = '2'
	TABOO     rune = '3'
	UNDEFINED rune = '4'
)

type Board struct {
	size     int
	Board    string
	currMove rune
}

func NewBoardFromFile(path string) *Board {
	csvConf, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := csvConf.Close(); err != nil {
			panic(err)
		}
	}()

	reader := csv.NewReader(csvConf)
	reader.Comma = ' '
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	var board Board
	board.size, _ = strconv.Atoi(data[0][0])
	data = data[1:]
	sb := strings.Builder{}
	for _, row := range data {
		for _, el := range row {
			sb.WriteString(el)
		}
	}

	board.Board = sb.String()
	board.currMove = WHITE

	return &board
}

func (b Board) Inverse() *Board {
	inv := Board{size: b.size, currMove: UNDEFINED}

	sb := strings.Builder{}
	for i := 0; i < inv.size*inv.size; i++ {
		switch rune(b.Board[i]) {
		case BLACK:
			sb.WriteRune(WHITE)
		case WHITE:
			sb.WriteRune(BLACK)
		default:
			sb.WriteRune(rune(b.Board[i]))
		}
	}
	inv.Board = sb.String()

	return &inv
}

func (b Board) GetNeighbours() []Board {
	var moves []Board
	for j := 0; j < b.size; j++ {
		for i := 0; i < b.size; i++ {
			if rune(b.Board[j*b.size+i]) == b.currMove {
				if h := b.getHorizontalMoves(i, j); len(h) != 0 {
					moves = append(moves, h...)
				}
				if v := b.getVerticalMoves(i, j); len(v) != 0 {
					moves = append(moves, v...)
				}
			}
		}
	}

	if len(moves) == 0 {
		return nil
	}

	return moves
}

func (b Board) getHorizontalMoves(x, y int) []Board {
	var moves []Board
	// right
	if b.isOnBoard(x+1, y) {
		if rune(b.Board[y*b.size+x+1]) == FREE {
			moves = append(moves, b.getMove(x, y, x+1, y))
		}
	}
	if b.isOnBoard(x+2, y) &&
		rune(b.Board[y*b.size+x+1]) != FREE && rune(b.Board[y*b.size+x+1]) != TABOO {
		if rune(b.Board[y*b.size+x+2]) == FREE {
			moves = append(moves, b.getMove(x, y, x+2, y))
		}
	}
	// left
	if b.isOnBoard(x-1, y) {
		if rune(b.Board[y*b.size+x-1]) == FREE {
			moves = append(moves, b.getMove(x, y, x-1, y))
		}
	}
	if b.isOnBoard(x-2, y) {
		if rune(b.Board[y*b.size+x-2]) == FREE &&
			rune(b.Board[y*b.size+x-1]) != FREE && rune(b.Board[y*b.size+x-1]) != TABOO {
			moves = append(moves, b.getMove(x, y, x-2, y))
		}
	}

	return moves
}

func (b Board) getVerticalMoves(x, y int) []Board {
	var moves []Board
	// top
	if b.isOnBoard(x, y+1) {
		if rune(b.Board[(y+1)*b.size+x]) == FREE {
			moves = append(moves, b.getMove(x, y, x, y+1))
		}
	}
	if b.isOnBoard(x, y+2) &&
		rune(b.Board[(y+1)*b.size+x]) != FREE &&
		rune(b.Board[(y+1)*b.size+x]) != TABOO {
		if rune(b.Board[(y+2)*b.size+x]) == FREE {
			moves = append(moves, b.getMove(x, y, x, y+2))
		}
	}
	// bottom
	if b.isOnBoard(x, y-1) {
		if rune(b.Board[(y-1)*b.size+x]) == FREE {
			moves = append(moves, b.getMove(x, y, x, y-1))
		}
	}
	if b.isOnBoard(x, y-2) {
		if rune(b.Board[(y-2)*b.size+x]) == FREE &&
			rune(b.Board[(y-1)*b.size+x]) != FREE &&
			rune(b.Board[(y-1)*b.size+x]) != TABOO {
			moves = append(moves, b.getMove(x, y, x, y-2))
		}
	}

	return moves
}

func (b Board) isOnBoard(x, y int) bool {
	if x < 0 || y < 0 || x > b.size-1 || y > b.size-1 {
		return false
	}

	return true
}

func (b Board) getMove(oldX, oldY, newX, newY int) Board {
	m := Board{size: b.size}
	if b.currMove == WHITE {
		m.currMove = BLACK
	} else {
		m.currMove = WHITE
	}
	r := []rune(b.Board)
	r[oldY*m.size+oldX], r[newY*m.size+newX] = FREE, r[oldY*m.size+oldX]
	m.Board = string(r)

	return m
}

// Sum of manh distances of each checker to corner
func (b Board) Heuristic(to Board) int {
	res := 0
	for j := 0; j < b.size; j++ {
		for i := 0; i < b.size; i++ {
			switch rune(b.Board[j*b.size+i]) {
			case BLACK:
				res += b.size - 1 - i + b.size - 1 - j
			case WHITE:
				res += i + j
			}
		}
	}

	return res
}

// Print pseudographic of board
var Symbols = map[rune]string{BLACK: "◎", WHITE: "◉", TABOO: "✕", FREE: " "}

func (b Board) Print() {
	// First row
	fmt.Print("⎾" + Symbols[rune(b.Board[0])])
	for i := 1; i < b.size; i++ {
		fmt.Print("⏉" + Symbols[rune(b.Board[i])])
	}
	fmt.Print("⏋\n")

	// Middle
	for j := 1; j < b.size-1; j++ {
		fmt.Print("⎾" + Symbols[rune(b.Board[j*b.size])])
		for i := 1; i < b.size; i++ {
			fmt.Print("⏉" + Symbols[rune(b.Board[j*b.size+i])])
		}
		fmt.Print("⏋\n")
	}

	// Last row
	fmt.Print("⎾" + Symbols[rune(b.Board[b.size*(b.size-1)])])
	for i := 1; i < b.size; i++ {
		fmt.Print("⏉" + Symbols[rune(b.Board[b.size*(b.size-1)+i])])
	}
	fmt.Print("⏋\n")
}
