package a_search

type Pather interface {
	GetNeighbours() []Pather
	GetCost(to Pather) int
	Heuristic(to Pather) int
}
