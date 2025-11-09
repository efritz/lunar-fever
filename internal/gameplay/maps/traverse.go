package maps

type point struct {
	row int
	col int
}

type neighborMeta struct {
	delta                point // row/col delta from self -> neighbor
	obstacleBitsSelf     int64 // impassable obstacle bits set on self
	obstacleBitsNeighbor int64 // impassable obstacle bits set on neighbor
}

var neighborMetas = []neighborMeta{
	{delta: point{row: -1, col: +0}, obstacleBitsSelf: bits(INTERIOR_WALL_N_BIT, DOOR_N_BIT), obstacleBitsNeighbor: bits(INTERIOR_WALL_S_BIT, DOOR_S_BIT)},
	{delta: point{row: +1, col: +0}, obstacleBitsSelf: bits(INTERIOR_WALL_S_BIT, DOOR_S_BIT), obstacleBitsNeighbor: bits(INTERIOR_WALL_N_BIT, DOOR_N_BIT)},
	{delta: point{row: +0, col: -1}, obstacleBitsSelf: bits(INTERIOR_WALL_W_BIT, DOOR_W_BIT), obstacleBitsNeighbor: bits(INTERIOR_WALL_E_BIT, DOOR_E_BIT)},
	{delta: point{row: +0, col: +1}, obstacleBitsSelf: bits(INTERIOR_WALL_E_BIT, DOOR_E_BIT), obstacleBitsNeighbor: bits(INTERIOR_WALL_W_BIT, DOOR_W_BIT)},
}

// traverse assigns the given identifier into the given two-dimensional boolean array representing
// the connected floor tiles reachable from the given starting point. The board is shared across
// all traversals so that the entire tile map can be efficiently processed into connected components.
// If the starting point is not a floor tile or has already been visited, a false-valued flag is
// returned.
func traverse(tileMap *TileMap, board [][]int, p point, id int) bool {
	// Skip the traversal entirely if the source tile is not a floor tile or has been visited
	if !tileMap.GetBit(p.row, p.col, FLOOR_BIT) || board[p.row][p.col] != 0 {
		return false
	}

	queue := []point{p}
	board[p.row][p.col] = id

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]

		for _, neighborMeta := range neighborMetas {
			neighbor := point{
				row: p.row + neighborMeta.delta.row,
				col: p.col + neighborMeta.delta.col,
			}

			// Check for out of bounds
			if neighbor.row < 0 || neighbor.row >= tileMap.Height() || neighbor.col < 0 || neighbor.col >= tileMap.Width() {
				continue
			}

			// Skip this neighbor if it has already been visited
			if board[neighbor.row][neighbor.col] != 0 {
				continue
			}

			selfBits := tileMap.GetBits(p.row, p.col)
			neighborBits := tileMap.GetBits(neighbor.row, neighbor.col)

			// If neighbor is a floor tile and there is no obstacle between the two tiles, continue traversal
			if neighborBits&bits(FLOOR_BIT) != 0 && selfBits&neighborMeta.obstacleBitsSelf == 0 && neighborBits&neighborMeta.obstacleBitsNeighbor == 0 {
				queue = append(queue, neighbor)
				board[neighbor.row][neighbor.col] = id
			}
		}
	}

	return true
}
