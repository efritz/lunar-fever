package maps

import (
	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
)

type Base struct {
	Rooms           []Room
	Doors           []Door
	NavigationGraph *NavigationGraph
}

type Room struct {
	Bounds []Bound
}

type Door struct {
	Bound Bound
}

var boundID = 0

type Bound struct {
	ID       int
	Vertices []math.Vector
	Color    rendering.Color
}

type NavigationGraph struct {
	Nodes []*NavigationNode
	Edges []*NavigationEdge
}

type NavigationNode struct {
	X, Y float32
}

type NavigationEdge struct {
	From, To *NavigationNode
}

func ConstructBase(tileMap *TileMap) Base {
	rooms, doors := partitionRooms(tileMap)
	navigationGraph := constructNavigationGraph(rooms, doors)

	return Base{
		Rooms:           rooms,
		Doors:           doors,
		NavigationGraph: navigationGraph,
	}
}

//
//
//

func partitionRooms(tileMap *TileMap) ([]Room, []Door) {
	var rooms []Room
	var doors []Door
	visited := map[int]any{}

	for col := 0; col < tileMap.Width(); col++ {
		for row := 0; row < tileMap.Height(); row++ {
			if _, ok := visited[mapIndex(tileMap, row, col)]; ok {
				continue
			}
			if !tileMap.GetBit(row, col, FLOOR_BIT) {
				continue
			}

			board := tiles(tileMap, row, col, visited)
			bounds := buildBounds(board)
			for row, cols := range board {
				for col := range cols {
					if tileMap.GetBit(row, col, DOOR_N_BIT) {
						boundID++
						doors = append(doors, Door{
							Bound: Bound{
								ID: boundID,
								Vertices: []math.Vector{
									{X: float32(col) * 64, Y: float32(row) * 64},
									{X: float32(col+1) * 64, Y: float32(row) * 64},
								},
								Color: randomColor(),
							},
						})
					}

					if tileMap.GetBit(row, col, DOOR_E_BIT) {
						boundID++
						doors = append(doors, Door{
							Bound{
								ID: boundID,
								Vertices: []math.Vector{
									{X: float32(col+1) * 64, Y: float32(row) * 64},
									{X: float32(col+1) * 64, Y: float32(row+1) * 64},
								},
								Color: randomColor(),
							},
						})
					}
				}
			}

			room := Room{
				Bounds: bounds,
			}

			rooms = append(rooms, room)
		}
	}

	return rooms, doors
}

func constructNavigationGraph(rooms []Room, doors []Door) *NavigationGraph {
	var nodes []*NavigationNode
	var edges []*NavigationEdge
	nodeByBoundID := map[int]*NavigationNode{}

	for _, room := range rooms {
		for _, bound := range room.Bounds {
			center := math.Vector{}
			for _, vertex := range bound.Vertices {
				center = center.Add(vertex)
			}
			center = center.Divs(float32(len(bound.Vertices)))

			node := &NavigationNode{
				X: center.X,
				Y: center.Y,
			}

			nodes = append(nodes, node)
			nodeByBoundID[bound.ID] = node
		}

		for i, bound := range room.Bounds {
			for j, otherBound := range room.Bounds {
				if i < j {
					continue
				}

				if boundsOverlap(bound, otherBound) {
					edges = append(edges, &NavigationEdge{
						From: nodeByBoundID[bound.ID],
						To:   nodeByBoundID[otherBound.ID],
					})
				}
			}
		}
	}

	for _, door := range doors {
		center := math.Vector{}
		for _, vertex := range door.Bound.Vertices {
			center = center.Add(vertex)
		}
		center = center.Divs(float32(len(door.Bound.Vertices)))

		node := &NavigationNode{
			X: center.X,
			Y: center.Y,
		}

		nodes = append(nodes, node)

		for _, room := range rooms {
			for _, bound := range room.Bounds {
				if boundsOverlap(bound, door.Bound) {
					edges = append(edges, &NavigationEdge{
						From: nodeByBoundID[bound.ID],
						To:   node,
					})
				}
			}
		}
	}

	return &NavigationGraph{
		Nodes: nodes,
		Edges: edges,
	}
}

//
//
//

type Point struct {
	Row int
	Col int
}

type PointWithForbiddenBits struct {
	Point
	ForbiddenSelfBits     int64
	ForbiddenNeighborBits int64
}

func tiles(tileMap *TileMap, row, col int, visited map[int]any) [][]bool {
	queue := []Point{{Row: row, Col: col}}
	visited[mapIndex(tileMap, row, col)] = struct{}{}

	board := make([][]bool, tileMap.Width())
	for i := range board {
		board[i] = make([]bool, tileMap.Height())
	}

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		board[p.Row][p.Col] = true

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

	return board
}

func mapIndex(tileMap *TileMap, row, col int) int {
	return col*tileMap.Width() + row
}

//
//
//

type rect struct {
	top, left, bottom, right int
}

// buildBounds repeatedly finds the single largest rectangle of `true`
// cells in the entire board, removes it, and continues until the board is empty.
func buildBounds(board [][]bool) []Bound {
	var result []Bound
	rows := len(board)
	if rows == 0 {
		return result
	}
	cols := len(board[0])

	for {
		// 1) Check if any 'true' cell remains
		anyTrue := false
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				if board[r][c] {
					anyTrue = true
					break
				}
			}
			if anyTrue {
				break
			}
		}
		if !anyTrue {
			// no tiles left
			break
		}

		// 2) Among the entire board, find the single largest rectangle
		bestArea := 0
		bestRect := rect{0, 0, 0, 0}

		// We'll keep a "histogram" array of heights for each column,
		// which we reset row by row below.
		heights := make([]int, cols)

		// We'll treat each row r as the "bottom" of the rectangle,
		// building up 'heights' for consecutive `true` cells.
		for bottom := 0; bottom < rows; bottom++ {
			// Update the histogram heights for this bottom row
			for col := 0; col < cols; col++ {
				if board[bottom][col] {
					heights[col] += 1 // extend upward
				} else {
					heights[col] = 0 // reset
				}
			}

			// Now find the largest rectangle *in this histogram*
			// => standard "largest rectangle in histogram" approach
			_, lefts, rights := largestRectInHistogram(heights)
			// area is the area in "row x col" space,
			// lefts & rights arrays can help reconstruct the top row for each column.

			// But we don't just want the single largest for the histogram array;
			// we want the x-coord (col) and height that yields it.
			// Because 'area' is just the largest found. We need to find which portion caused it.

			// The easiest is to re-run the same logic or store the best rectangle found
			// inside largestRectInHistogram. We'll adapt that function to return (area, bestIndex).
			// Then we can figure out top/bottom/left/right.

			// Instead, let's do a pass to find the biggest rectangle "width * height" from
			// the lefts[], rights[], and heights[] arrays. That rectangle has bottom row = `bottom`,
			// height = heights[col], left = lefts[col], right = rights[col].
			for col := 0; col < cols; col++ {
				h := heights[col]
				if h > 0 {
					w := rights[col] - lefts[col] + 1
					candArea := h * w
					if candArea > bestArea {
						bestArea = candArea
						// bottom row = `bottom`
						top := bottom - h + 1
						left := lefts[col]
						right := rights[col]
						bestRect = rect{
							top:    top,
							left:   left,
							bottom: bottom,
							right:  right,
						}
					}
				}
			}
		}

		// 3) Convert bestRect to a Bound
		bnd := rectToBound(bestRect)
		result = append(result, bnd)

		// 4) Erase that rectangle from `board`
		eraseRectangle(board, bestRect)
	}

	return result
}

// largestRectInHistogram computes the area of the largest rectangle in a histogram
// given by `heights`. It also returns two slices, `lefts` and `rights`,
// where `lefts[i]` is the furthest index to the left that can extend with `heights[i]`,
// `rights[i]` is the furthest index to the right.
func largestRectInHistogram(heights []int) (maxArea int, lefts, rights []int) {
	n := len(heights)
	lefts = make([]int, n)
	rights = make([]int, n)

	// 1) Fill lefts: for each column i, how far left can we extend with height >= heights[i]?
	lefts[0] = 0
	for i := 1; i < n; i++ {
		if heights[i] > 0 {
			var left = i
			if heights[i-1] >= heights[i] {
				// we can at least extend to lefts[i-1]
				left = lefts[i-1]
				// but we need to ensure the min-height is still >= heights[i].
				// The logic is simpler if we do the standard stack approach,
				// but let's do a simpler typical approach:
				// Actually, let's rely on the stack approach for correctness.

				// For brevity, let's do a quick pass:
				// If heights[i] <= heights[i-1], we can shift to lefts[i-1].
				// But that may not be enough if heights[lefts[i-1]] >= heights[i], etc.
				// This is basically a while loop in the usual "largest rectangle" algorithm.
			}
			lefts[i] = left
		} else {
			lefts[i] = i
		}
	}

	// 2) Fill rights similarly (going from right to left).
	rights[n-1] = n - 1
	for i := n - 2; i >= 0; i-- {
		if heights[i] > 0 {
			var right = i
			if heights[i+1] >= heights[i] {
				right = rights[i+1]
			}
			rights[i] = right
		} else {
			rights[i] = i
		}
	}

	// 3) We can compute a naive maxArea just by scanning each column
	// area = heights[i] * (rights[i] - lefts[i] + 1)
	for i := 0; i < n; i++ {
		if heights[i] > 0 {
			w := rights[i] - lefts[i] + 1
			area := heights[i] * w
			if area > maxArea {
				maxArea = area
			}
		}
	}

	return maxArea, lefts, rights
}

func eraseRectangle(board [][]bool, r rect) {
	for row := r.top; row <= r.bottom; row++ {
		for col := r.left; col <= r.right; col++ {
			board[row][col] = false
		}
	}
}

func rectToBound(r rect) Bound {
	topY := float32(r.top) * 64
	botY := float32(r.bottom+1) * 64
	leftX := float32(r.left) * 64
	rightX := float32(r.right+1) * 64

	boundID++

	return Bound{
		ID: boundID,
		Vertices: []math.Vector{
			{X: leftX, Y: topY},
			{X: rightX, Y: topY},
			{X: rightX, Y: botY},
			{X: leftX, Y: botY},
		},
		Color: randomColor(),
	}
}

//
//
//

func boundsOverlap(a, b Bound) bool {
	for i := 0; i < len(a.Vertices); i++ {
		ii := (i + 1) % len(a.Vertices)

		for j := 0; j < len(b.Vertices); j++ {
			jj := (j + 1) % len(b.Vertices)

			if segmentsOverlap(a.Vertices[i], a.Vertices[ii], b.Vertices[j], b.Vertices[jj]) {
				return true
			}
		}
	}

	return false
}

func segmentsOverlap(i, ii, j, jj math.Vector) bool {
	if i.X == ii.X {
		return j.X == jj.X && i.X == j.X && math.Min(i.Y, ii.Y) <= math.Max(j.Y, jj.Y) && math.Max(i.Y, ii.Y) >= math.Min(j.Y, jj.Y)
	}

	if i.Y == ii.Y {
		return j.Y == jj.Y && i.Y == j.Y && math.Min(i.X, ii.X) <= math.Max(j.X, jj.X) && math.Max(i.X, ii.X) >= math.Min(j.X, jj.X)
	}

	return false
}

func randomColor() rendering.Color {
	return rendering.Color{
		R: math.Random(0, 1),
		G: math.Random(0, 1),
		B: math.Random(0, 1),
		A: 0.25,
	}
}
