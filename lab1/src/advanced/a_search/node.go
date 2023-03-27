package a_search

import "src/board"

type node struct {
	pather board.Board
	parent *node
	cost   int
	opened bool
	closed bool
	index  int
	rank   int
}

// Collection of nodes, indexed by their Pather component
type nodeMap map[board.Board]*node

// Method to get pather from collection
// or add new element and return pointer
func (nm nodeMap) get(p board.Board) *node {
	n, ok := nm[p]
	if !ok {
		n = &node{
			pather: p,
		}
		nm[p] = n
	}
	return n
}
