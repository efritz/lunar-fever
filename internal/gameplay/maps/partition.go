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

	fixtures := []Bound{}
	for col := 0; col < tileMap.Width(); col++ {
		for row := 0; row < tileMap.Height(); row++ {
			if fixture, ok := tileMap.GetFixture(row, col); ok {
				fixtures = append(fixtures, newBound(
					math.Vector{float32(col * 64), float32(row * 64)},
					math.Vector{float32(col+fixture.TileWidth) * 64, float32(row * 64)},
					math.Vector{float32(col+fixture.TileWidth) * 64, float32(row+fixture.TileHeight) * 64},
					math.Vector{float32(col * 64), float32(row+fixture.TileHeight) * 64},
				))
			}
		}
	}

	// Next, we'll convert each of these connected components into a list of bounds.
	// This creates one bound per tile, which we'll transform in the next steps.

	boundsByID := map[int][]Bound{}
	for row, cols := range board {
		for col, id := range cols {
			if id == 0 {
				continue
			}

			boundsByID[id] = append(boundsByID[id], newBound(
				math.Vector{float32(col * 64), float32(row * 64)},
				math.Vector{float32(col+1) * 64, float32(row * 64)},
				math.Vector{float32(col+1) * 64, float32(row+1) * 64},
				math.Vector{float32(col * 64), float32(row+1) * 64},
			))

			// vec := func(col, row int) math.Vector {
			// 	return math.Vector{
			// 		X: float32(col),
			// 		Y: float32(row),
			// 	}
			// }

			// ul := newBound(vec(col*64+0*32, row*64+0*32), vec(col*64+1*32, row*64+0*32), vec(col*64+1*32, row*64+1*32), vec(col*64+0*32, row*64+1*32))
			// ur := newBound(vec(col*64+1*32, row*64+0*32), vec(col*64+2*32, row*64+0*32), vec(col*64+2*32, row*64+1*32), vec(col*64+1*32, row*64+1*32))
			// ll := newBound(vec(col*64+0*32, row*64+1*32), vec(col*64+1*32, row*64+1*32), vec(col*64+1*32, row*64+2*32), vec(col*64+0*32, row*64+2*32))
			// lr := newBound(vec(col*64+1*32, row*64+1*32), vec(col*64+2*32, row*64+1*32), vec(col*64+2*32, row*64+2*32), vec(col*64+1*32, row*64+2*32))

			// if tileMap.GetBit(row, col, DOOR_N_BIT) || tileMap.GetBit(row, col, DOOR_W_BIT) ||
			// 	(!tileMap.GetBit(row, col, INTERIOR_WALL_N_BIT) && !tileMap.GetBit(row, col, INTERIOR_WALL_W_BIT) &&
			// 		!tileMap.GetBit(row, col-1, INTERIOR_WALL_N_BIT) && !tileMap.GetBit(row-1, col, INTERIOR_WALL_W_BIT)) {
			// 	boundsByID[id] = append(boundsByID[id], ul)
			// }

			// if tileMap.GetBit(row, col, DOOR_N_BIT) || tileMap.GetBit(row, col, DOOR_E_BIT) ||
			// 	(!tileMap.GetBit(row, col, INTERIOR_WALL_N_BIT) && !tileMap.GetBit(row, col, INTERIOR_WALL_E_BIT) &&
			// 		!tileMap.GetBit(row, col+1, INTERIOR_WALL_N_BIT) && !tileMap.GetBit(row-1, col, INTERIOR_WALL_E_BIT)) {
			// 	boundsByID[id] = append(boundsByID[id], ur)
			// }

			// if tileMap.GetBit(row, col, DOOR_S_BIT) || tileMap.GetBit(row, col, DOOR_W_BIT) ||
			// 	(!tileMap.GetBit(row, col, INTERIOR_WALL_S_BIT) && !tileMap.GetBit(row, col, INTERIOR_WALL_W_BIT) &&
			// 		!tileMap.GetBit(row, col-1, INTERIOR_WALL_S_BIT) && !tileMap.GetBit(row+1, col, INTERIOR_WALL_W_BIT)) {
			// 	boundsByID[id] = append(boundsByID[id], ll)
			// }

			// if tileMap.GetBit(row, col, DOOR_S_BIT) || tileMap.GetBit(row, col, DOOR_E_BIT) ||
			// 	(!tileMap.GetBit(row, col, INTERIOR_WALL_S_BIT) && !tileMap.GetBit(row, col, INTERIOR_WALL_E_BIT) &&
			// 		!tileMap.GetBit(row, col+1, INTERIOR_WALL_S_BIT) && !tileMap.GetBit(row+1, col, INTERIOR_WALL_E_BIT)) {
			// 	boundsByID[id] = append(boundsByID[id], lr)
			// }
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

		rooms = append(rooms, newRoom(
			triangulate(
				splitBoundsAtIntersections(
					subtract(
						simplifyBounds(
							mergeBounds(
								bounds,
								obstacles,
							),
						),
						walls,
						doors,
						fixtures,
					),
					doors,
				),
			),
		))
	}

	return rooms, walls, doors
}

func extractWallsAndDoors(tileMap *TileMap) (walls []Edge, doors []Edge) {
	for col := 0; col < tileMap.Width(); col++ {
		for row := 0; row < tileMap.Height(); row++ {
			if tileMap.GetBit(row, col, INTERIOR_WALL_N_BIT) {
				walls = append(walls, newEdge(vec(col, row), vec(col+1, row)))
			}
			if tileMap.GetBit(row, col, INTERIOR_WALL_E_BIT) {
				walls = append(walls, newEdge(vec(col+1, row), vec(col+1, row+1)))
			}
			if tileMap.GetBit(row, col, INTERIOR_WALL_S_BIT) {
				walls = append(walls, newEdge(vec(col, row+1), vec(col+1, row+1)))
			}
			if tileMap.GetBit(row, col, INTERIOR_WALL_W_BIT) {
				walls = append(walls, newEdge(vec(col, row), vec(col, row+1)))
			}

			if tileMap.GetBit(row, col, DOOR_N_BIT) {
				doors = append(doors, newEdge(vec(col, row), vec(col+1, row)))
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
