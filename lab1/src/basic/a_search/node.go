package a_search

import "src/board"

type node struct {
	board  board.Board
	parent *node
	cost   int
	opened bool
	closed bool
	index  int
	rank   int
}

type nodeMap map[board.Board]*node

func (nm nodeMap) get(p board.Board) *node {
	n, ok := nm[p]
	if !ok {
		n = &node{
			board: p,
		}
		nm[p] = n
	}
	return n
}
