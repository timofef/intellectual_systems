package a_search

import (
	"container/heap"
	"fmt"
	"src/board"
)

func A(start, terminal board.Board) ([]board.Board, int, bool) {
	allNodes := nodeMap{}
	openedList := &priorityQueue{}
	heap.Init(openedList)

	// Init OPENED list
	fromNode := allNodes.get(start)
	fromNode.opened = true
	heap.Push(openedList, fromNode)

	for {
		// If there's no path -> failure
		if openedList.Len() == 0 {
			return nil, 0, false
		}

		// Close best node
		current := heap.Pop(openedList).(*node)
		current.opened = false
		current.closed = true

		// If found end -> trace back and return path
		if current.pather.Board == allNodes.get(terminal).pather.Board {
			var p []board.Board
			curr := current
			for curr != nil {
				p = append(p, curr.pather)
				curr = curr.parent
			}
			fmt.Printf("Opened: %d\n", openedList.Len())
			return p, current.cost, true
		}

		for _, neighbour := range current.pather.GetNeighbours() {

			cost := current.cost + 1
			neighborNode := allNodes.get(neighbour)
			// If already in OPENED -> check if cost is lower
			if cost < neighborNode.cost {
				if neighborNode.opened {
					heap.Remove(openedList, neighborNode.index)
				}
			}
			// If completely new node -> add to OPENED
			if !neighborNode.opened && !neighborNode.closed {
				neighborNode.cost = cost
				neighborNode.opened = true
				neighborNode.rank = cost + neighbour.Heuristic(terminal)
				neighborNode.parent = current
				heap.Push(openedList, neighborNode)
			}
		}

	}
}
