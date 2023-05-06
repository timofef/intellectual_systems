package sudoku

func (s *Sudoku) verticalConstraint(idx int) uint32 {
	var res uint32

	j := idx % s.size

	for k := 0; k < s.size; k++ {
		res |= s.field[s.size*k+j]
	}

	return res
}

func (s *Sudoku) horizontalConstraint(idx int) uint32 {
	var res uint32

	i := idx / s.size

	for k := 0; k < s.size; k++ {
		res |= s.field[s.size*i+k]
	}

	return res
}

func (s *Sudoku) blockConstraint(idx int) uint32 {
	var res uint32

	// Calculate block indexes
	i := idx / (s.subSize * s.size)
	j := (idx % s.size) / s.subSize

	for k := 0; k < s.subSize; k++ {
		for l := 0; l < s.subSize; l++ {
			res |= s.field[i*s.subSize*s.size+k*s.size+j*s.subSize+l]
		}
	}

	return res
}

func (s *Sudoku) heuristic() int {
	var res int
	// horizontal
	for i := 0; i < s.size; i++ {
		var heuristic uint32
		for j := 0; j < s.size; j++ {
			heuristic |= s.field[i*s.size+j]
		}
		res += countZeros(heuristic, s.size)
	}

	// vertical
	for j := 0; j < s.size; j++ {
		var heuristic uint32
		for i := 0; i < s.size; i++ {
			heuristic |= s.field[i*s.size+j]
		}
		res += countZeros(heuristic, s.size)
	}

	// block
	for i := 0; i < s.subSize; i++ {
		for j := 0; j < s.subSize; j++ {
			var heuristic uint32
			for k := 0; k < s.subSize; k++ {
				for l := 0; l < s.subSize; l++ {
					heuristic |= s.field[i*s.subSize*s.size+k*s.size+j*s.subSize+l]
				}
			}
			res += countZeros(heuristic, s.size)
		}
	}

	return res
}
