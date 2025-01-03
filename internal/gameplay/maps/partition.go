package maps

import (
	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
)

type Room struct {
	Bounds []Bound
	Color  rendering.Color
}

func newRoom(bounds []Bound) Room {
	return Room{
		Bounds: bounds,
		Color:  randomColor(),
	}
}

// partitionRooms takes a tile map and returns a rooms and lists of edges representing
// the bounds of walls and doors. Each room's bounds are triangulated.
func partitionRooms(tileMap *TileMap) (_ []Room, walls []Edge, doors []Edge) {
	// First, we'll traverse the tile map to find connected components. Each group of
	// mutually navigable floor tiles will be given a unique integer ID in the following
	// two-dimensional board.

	board := make([][]int, tileMap.Width())
	for i := range board {
		board[i] = make([]int, tileMap.Height())
	}

	id := 1
	for col := 0; col < tileMap.Width(); col++ {
		for row := 0; row < tileMap.Height(); row++ {
			if traverse(tileMap, board, point{row: row, col: col}, id) {
				id++
			}
		}
	}

	// Next, we'll convert each of these connected components into a list of bounds.
	// This creates one bound per tile, which we'll transform in the next steps.

	boundsByID := map[int][]Bound{}
	for row, cols := range board {
		for col, id := range cols {
			if id != 0 {
				boundsByID[id] = append(boundsByID[id], newBound(
					vec(col, row),
					vec(col+1, row),
					vec(col+1, row+1),
					vec(col, row+1),
				))
			}
		}
	}

	walls, doors = extractWallsAndDoors(tileMap)
	obstacles := append(append([]Edge(nil), walls...), doors...)

	var rooms []Room
	for _, bounds := range boundsByID {
		// Transform the tiles of each connected component into the bounds of a room by:
		//
		// (1) Merging the set of single-tile bounds into more complex polygons (being cautious of holes).
		// (2) Simplifying the vertex list of the merged polygons by removing collinear points.
		// (3) Adding back vertices that denote the extent of overlap with doors and other bounds (to help with triangulation).
		// (4) Triangulating the resulting polygons.

		rooms = append(rooms, newRoom(triangulate(splitBoundsAtIntersections(simplifyBounds(mergeBounds(bounds, obstacles)), doors))))
	}

	return rooms, walls, doors
}

func extractWallsAndDoors(tileMap *TileMap) (walls []Edge, doors []Edge) {
	for col := 0; col < tileMap.Width(); col++ {
		for row := 0; row < tileMap.Height(); row++ {
			if tileMap.GetBit(row, col, INTERIOR_WALL_N_BIT) {
				walls = append(walls, newEdge(vec(col, row), vec(col+1, row)))
			}
			if tileMap.GetBit(row, col, DOOR_N_BIT) {
				doors = append(doors, newEdge(vec(col, row), vec(col+1, row)))
			}

			if tileMap.GetBit(row, col, INTERIOR_WALL_E_BIT) {
				walls = append(walls, newEdge(vec(col+1, row), vec(col+1, row+1)))
			}
			if tileMap.GetBit(row, col, DOOR_E_BIT) {
				doors = append(doors, newEdge(vec(col+1, row), vec(col+1, row+1)))
			}
		}
	}

	return walls, doors
}

func vec(col, row int) math.Vector {
	return math.Vector{
		X: float32(col) * 64,
		Y: float32(row) * 64,
	}
}
