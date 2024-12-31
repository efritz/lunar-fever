package maps

import (
	"sort"

	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
)

type Room struct {
	Bounds []Bound
	Color  rendering.Color
}

type Bound struct {
	MinX  int
	MaxX  int
	MinY  int
	MaxY  int
	Color rendering.Color
}

func PartitionRooms(tileMap *TileMap) []Room {
	var rooms []Room
	visited := map[int]any{}

	for col := 0; col < tileMap.Width(); col++ {
		for row := 0; row < tileMap.Height(); row++ {
			if _, ok := visited[mapIndex(tileMap, row, col)]; ok {
				continue
			}
			if !tileMap.GetBit(row, col, FLOOR_BIT) {
				continue
			}

			rooms = append(rooms, Room{
				Bounds: buildBounds(tiles(tileMap, row, col, visited)),
				Color:  randomColor(),
			})
		}
	}

	return rooms
}

type Point struct {
	Row int
	Col int
}

type ColumnSegment struct {
	StartCol int
	EndCol   int
}

func buildBounds(points []Point) []Bound {
	colsByRow := map[int][]int{}
	for _, point := range points {
		if _, ok := colsByRow[point.Row]; !ok {
			colsByRow[point.Row] = []int{}
		}

		colsByRow[point.Row] = append(colsByRow[point.Row], point.Col)
	}

	segmentsByRow := map[int][]ColumnSegment{}
	for row, cols := range colsByRow {
		sort.Ints(cols)

		var segments []ColumnSegment

		activeSegment := ColumnSegment{cols[0], cols[0]}
		for i := 1; i < len(cols); i++ {
			if cols[i] == activeSegment.EndCol+1 {
				activeSegment.EndCol = cols[i]
				continue
			}

			segments = append(segments, activeSegment)
			activeSegment = ColumnSegment{cols[i], cols[i]}
		}

		segments = append(segments, activeSegment)
		segmentsByRow[row] = segments
	}

	var rows []int
	for row := range segmentsByRow {
		rows = append(rows, row)
	}
	sort.Ints(rows)

	var bounds []Bound
	for i := 0; i < len(rows); i++ {
		row := rows[i]
		for _, segment := range segmentsByRow[row] {
			bound := Bound{
				MinX:  segment.StartCol,
				MaxX:  segment.EndCol,
				MinY:  rows[i],
				MaxY:  rows[i],
				Color: randomColor(),
			}

			for j := 1; ; j++ {
				nextSegments, ok := segmentsByRow[row+j]
				if !ok {
					break
				}

				found := false
				for k, s := range nextSegments {
					if s.StartCol == segment.StartCol && s.EndCol == segment.EndCol {
						found = true
						bound.MaxY = rows[i+j]
						segmentsByRow[row+j] = append(segmentsByRow[row+j][:k], segmentsByRow[row+j][k+1:]...)
						break
					}
				}

				if !found {
					break
				}
			}

			bounds = append(bounds, bound)
		}
	}

	return bounds
}

type PointWithForbiddenBits struct {
	Point
	ForbiddenSelfBits     int64
	ForbiddenNeighborBits int64
}

func tiles(tileMap *TileMap, row, col int, visited map[int]any) []Point {
	queue := []Point{{Row: row, Col: col}}
	visited[mapIndex(tileMap, row, col)] = struct{}{}

	var tiles []Point

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		tiles = append(tiles, p)

		for _, neighbor := range []PointWithForbiddenBits{
			{Point: Point{Row: p.Row - 1, Col: p.Col}, ForbiddenSelfBits: (1 << INTERIOR_WALL_N_BIT) | (1 << DOOR_N_BIT), ForbiddenNeighborBits: (1 << INTERIOR_WALL_S_BIT) | (1 << DOOR_S_BIT)},
			{Point: Point{Row: p.Row + 1, Col: p.Col}, ForbiddenSelfBits: (1 << INTERIOR_WALL_S_BIT) | (1 << DOOR_S_BIT), ForbiddenNeighborBits: (1 << INTERIOR_WALL_N_BIT) | (1 << DOOR_N_BIT)},
			{Point: Point{Row: p.Row, Col: p.Col - 1}, ForbiddenSelfBits: (1 << INTERIOR_WALL_W_BIT) | (1 << DOOR_W_BIT), ForbiddenNeighborBits: (1 << INTERIOR_WALL_E_BIT) | (1 << DOOR_E_BIT)},
			{Point: Point{Row: p.Row, Col: p.Col + 1}, ForbiddenSelfBits: (1 << INTERIOR_WALL_E_BIT) | (1 << DOOR_E_BIT), ForbiddenNeighborBits: (1 << INTERIOR_WALL_W_BIT) | (1 << DOOR_W_BIT)},
		} {
			if neighbor.Row < 0 || neighbor.Row >= tileMap.Width() || neighbor.Col < 0 || neighbor.Col >= tileMap.Height() {
				continue
			}
			if _, ok := visited[mapIndex(tileMap, neighbor.Row, neighbor.Col)]; ok {
				continue
			}
			if !tileMap.GetBit(neighbor.Row, neighbor.Col, FLOOR_BIT) {
				continue
			}
			if tileMap.GetBits(p.Row, p.Col)&neighbor.ForbiddenSelfBits != 0 {
				continue
			}
			if tileMap.GetBits(neighbor.Row, neighbor.Col)&neighbor.ForbiddenNeighborBits != 0 {
				continue
			}

			queue = append(queue, neighbor.Point)
			visited[mapIndex(tileMap, neighbor.Row, neighbor.Col)] = struct{}{}
		}
	}

	return tiles
}

func mapIndex(tileMap *TileMap, row, col int) int {
	return col*tileMap.Width() + row
}

func randomColor() rendering.Color {
	return rendering.Color{
		R: math.Random(0, 1),
		G: math.Random(0, 1),
		B: math.Random(0, 1),
		A: 0.25,
	}
}
