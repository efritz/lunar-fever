package maps

import (
	"fmt"
	stdmath "math"

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
	Color  rendering.Color
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
	Nodes map[int]*NavigationNode
	Edges []*NavigationEdge
}

type NavigationNode struct {
	X, Y  float32
	Bound Bound
}

type NavigationEdge struct {
	From int
	To   int
}

func ConstructBase(tileMap *TileMap) *Base {
	rooms, walls, doors := partitionRooms(tileMap)
	navigationGraph := constructNavigationGraph(rooms, walls, doors)

	return &Base{
		Rooms:           rooms,
		Doors:           doors,
		NavigationGraph: navigationGraph,
	}
}

//
//
//

func partitionRooms(tileMap *TileMap) ([]Room, []Door, []Door) {
	var rooms []Room
	var doors []Door
	var walls []Door
	visited := map[int]any{}

	for col := 0; col < tileMap.Width(); col++ {
		for row := 0; row < tileMap.Height(); row++ {
			if _, ok := visited[mapIndex(tileMap, row, col)]; ok {
				continue
			}
			if !tileMap.GetBit(row, col, FLOOR_BIT) {
				continue
			}

			regionCells, board := tiles(tileMap, row, col, visited)
			// bounds := buildBounds2(board)
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

					if tileMap.GetBit(row, col, INTERIOR_WALL_N_BIT) {
						boundID++
						walls = append(walls, Door{
							Bound{
								ID: boundID,
								Vertices: []math.Vector{
									{X: float32(col) * 64, Y: float32(row) * 64},
									{X: float32(col+1) * 64, Y: float32(row) * 64},
								},
								Color: rendering.Color{R: 0, G: 0, B: 0, A: 1},
							},
						})
					}

					if tileMap.GetBit(row, col, INTERIOR_WALL_E_BIT) {
						boundID++
						walls = append(walls, Door{
							Bound{
								ID: boundID,
								Vertices: []math.Vector{
									{X: float32(col+1) * 64, Y: float32(row) * 64},
									{X: float32(col+1) * 64, Y: float32(row+1) * 64},
								},
								Color: rendering.Color{R: 0, G: 0, B: 0, A: 1},
							},
						})
					}
				}
			}

			// room := Room{
			// 	Bounds: bounds,
			// }

			// rooms = append(rooms, room)

			// 1) BFS to collect all floor tiles in this region
			// regionCells, _ := tiles(tileMap, row, col, visited)
			// fmt.Printf("> regionCells=%v\n", regionCells)
			// regionCells is []Point

			// 2) Build boundary edges from these cells
			// boundaryEdges := buildBoundaryEdges(tileMap, regionCells)
			// fmt.Printf("> boundaryEdges=%v\n", boundaryEdges)

			// 3) Convert boundary edges to corner graph
			// cornerGraph := buildCornerGraph(boundaryEdges)
			// fmt.Printf("> cornerGraph=%v\n", cornerGraph)

			// 4) Trace out polygons
			// cornerLoops := tracePolygons(cornerGraph)
			// fmt.Printf("> cornerLoops=%v\n", cornerLoops)

			// 5) Convert each corner loop to a Bound (polygon)
			var polygons []Bound
			// for _, loop := range cornerLoops {
			// 	polygons = append(polygons, cornersToBound(loop))
			// }

			_ = regionCells
			for row, cols := range board {
				for col, ok := range cols {
					if !ok {
						continue
					}

					boundID++
					polygons = append(polygons, Bound{
						ID: boundID,
						Vertices: []math.Vector{
							{X: float32(col) * 64, Y: float32(row) * 64},
							{X: float32(col+1) * 64, Y: float32(row) * 64},
							{X: float32(col+1) * 64, Y: float32(row+1) * 64},
							{X: float32(col) * 64, Y: float32(row+1) * 64},
						},
						Color: randomColor(),
					})
				}
			}

			// 6) Construct the Room
			room := Room{
				Bounds: simplify3(polygons, doors, walls),
				Color:  randomColor(),
			}
			rooms = append(rooms, room)
		}
	}

	for i, room := range rooms {
		// fmt.Printf("Simplifying...\n")
		rooms[i].Bounds = simplify2(room.Bounds)
		// fmt.Printf("Done...\n")
	}

	for _, room := range rooms {
	outer1:
		for {
			for i, bound1 := range room.Bounds {
				for j, bound2 := range room.Bounds {
					if bound1.ID == bound2.ID {
						continue
					}

					if reduceTJunctions(&bound1, &bound2) {
						room.Bounds[i] = bound1
						room.Bounds[j] = bound2
						continue outer1
					}
				}
			}

			break
		}
	}

	for _, room := range rooms {
	outer2:
		for {
			for i, bound1 := range room.Bounds {
				for j, door := range doors {
					bound2 := door.Bound

					if reduceTJunctions(&bound1, &bound2) {
						room.Bounds[i] = bound1
						doors[j].Bound = bound2
						continue outer2
					}
				}
			}

			break
		}
	}

	for i, room := range rooms {
		// fmt.Printf("Triangulating...\n")
		rooms[i].Bounds = triangulate(room.Bounds)
		// fmt.Printf("Done...\n")
	}

	return rooms, walls, doors
}

func simplify3(bounds []Bound, doors, walls []Door) []Bound {
outer:
	for {
		for i, b1 := range bounds {
			for j, b2 := range bounds {
				if i == j {
					continue
				}

				if m, ok := tryMerge(b1, b2, doors, walls); ok {
					bounds[i] = m
					bounds[j] = bounds[len(bounds)-1]
					bounds = bounds[:len(bounds)-1]
					continue outer
				}
			}
		}

		break
	}

	return bounds
}

func tryMerge(b1, b2 Bound, doors, walls []Door) (Bound, bool) {
	i, j, c := findSharedEdges(b1, b2, doors, walls)
	if c == 0 {
		return Bound{}, false
	}

	var result []math.Vector
	for k := 0; k < len(b1.Vertices); k++ {
		if k == i {
			s := (j + c) % len(b2.Vertices)
			curr := s
			for {
				if curr == j {
					break
				}
				result = append(result, b2.Vertices[curr])
				curr = (curr + 1) % len(b2.Vertices)
			}
			k += c - 1
		} else {
			result = append(result, b1.Vertices[k])
		}
	}

	boundID++
	return Bound{
		ID:       boundID,
		Vertices: result,
		Color:    randomColor(),
	}, true
}

func findSharedEdges(b1, b2 Bound, doors, walls []Door) (int, int, int) {
	var is []int
	var js []int

	var vertices []math.Vector

	for i := 0; i < len(b1.Vertices); i++ {
		for j := 0; j < len(b2.Vertices); j++ {
			if b1.Vertices[i].Equal(b2.Vertices[j]) {
				is = append(is, i)
				js = append(js, j)
				vertices = append(vertices, b1.Vertices[i])
			}
		}
	}

	if len(is) < 2 {
		return 0, 0, 0
	}

	nextI := func(v int) int {
		return (v + 1) % len(b1.Vertices)
	}
	nextJ := func(v int) int {
		return (v + 1) % len(b2.Vertices)
	}

	// Ensure is are consecutive
	for i, v := range is {
		if i > 0 && v != nextI(is[i-1]) {
			return 0, 0, 0
		}
	}

	// Ensure js are consecutive
	for i, v := range js {
		if i > 0 && js[i-1] != nextJ(v) {
			return 0, 0, 0
		}
	}

	// Ensure vertices do not overlap with wall or door
	for _, s := range [][]Door{walls, doors} {
		for _, wall := range s {
			for k := 0; k < len(vertices)-1; k++ {
				if segmentsOverlap(vertices[k], vertices[k+1], wall.Bound.Vertices[0], wall.Bound.Vertices[1]) {
					// }
					// }
					// if boundsOverlap(Bound{Vertices: vertices}, wall.Bound) {
					// fmt.Printf("Overlap with wall or door\n")
					// fmt.Printf("> %v\n", vertices)
					// fmt.Printf("> %v\n", wall.Bound.Vertices)
					return 0, 0, 0
				}
			}
		}
	}

	return is[0], js[len(js)-1], len(is) - 1
}

//
//
//

func simplify2(bounds []Bound) []Bound {
	var f []Bound
	for _, b := range bounds {
		f = append(f, Bound{
			ID:       b.ID,
			Vertices: SimplifyPolygon(b.Vertices),
			Color:    b.Color,
		})
	}

	return f
}

func SimplifyPolygon(vertices []math.Vector) []math.Vector {
	if len(vertices) <= 2 {
		return vertices
	}

	// Helper function to check if two edges form a straight line
	isCollinear := func(a, b, c math.Vector) bool {
		// For axis-aligned edges, points are collinear if:
		// 1. a->b and b->c are both horizontal (same Y coordinates)
		// 2. a->b and b->c are both vertical (same X coordinates)
		abHorizontal := a.Y == b.Y
		bcHorizontal := b.Y == c.Y
		abVertical := a.X == b.X
		bcVertical := b.X == c.X

		return (abHorizontal && bcHorizontal) || (abVertical && bcVertical)
	}

	// Helper function to perform one pass of simplification
	simplifyOnce := func(points []math.Vector) ([]math.Vector, bool) {
		if len(points) <= 2 {
			return points, false
		}

		simplified := make([]math.Vector, 0, len(points))
		changed := false

		// Process all points including the last one, considering the wrap-around
		n := len(points)
		for i := 0; i < n; i++ {
			prev := points[(i-1+n)%n]
			curr := points[i]
			next := points[(i+1)%n]

			if !isCollinear(prev, curr, next) {
				simplified = append(simplified, curr)
			} else {
				changed = true
			}
		}

		return simplified, changed
	}

	// Repeatedly simplify until no more points can be removed
	current := vertices
	for {
		simplified, changed := simplifyOnce(current)
		if !changed || len(simplified) == len(current) {
			return simplified
		}
		current = simplified
	}
}

// If there's an overlapping edge between b1 and b2 but the edges are
// not exactly the same, we need to reduce the T-junctions by inserting
// an additional vertx at the intersection point.
func reduceTJunctions(b1 *Bound, b2 *Bound) bool {
	for i := 0; i < len(b1.Vertices); i++ {
		for j := 0; j < len(b2.Vertices); j++ {
			ii := (i + 1) % len(b1.Vertices)
			jj := (j + 1) % len(b2.Vertices)

			if !segmentsOverlap(b1.Vertices[i], b1.Vertices[ii], b2.Vertices[j], b2.Vertices[jj]) {
				continue
			}

			// If edges i->ii and j->jj are equal, skip
			if (b1.Vertices[i].Equal(b2.Vertices[j]) && b1.Vertices[ii].Equal(b2.Vertices[jj])) ||
				(b1.Vertices[i].Equal(b2.Vertices[jj]) && b1.Vertices[ii].Equal(b2.Vertices[j])) {
				continue
			}

			// Find intersection point between i->ii and j->jj
			if segmentIntersection(b1, b2, i, ii, j, jj) {
				return true
			}
			// fmt.Printf("Intersection of %v %v and %v %v at %v\n", b1.Vertices[i], b1.Vertices[ii], b2.Vertices[j], b2.Vertices[jj], intersection)
		}
	}

	return false
}

func insertVertex(vertices []math.Vector, i int, v math.Vector) ([]math.Vector, bool) {
	v2 := vertices[(i+1)%len(vertices)]
	if vertices[i].Equal(v) || v2.Equal(v) {
		return vertices, false
	}

	// Skip if v is not between vertices[i] and vertices[(i+1)%len(vertices)]
	// Note that vertices[i] and v2 might be in either relation to each other
	if vertices[i].X == v2.X {
		if vertices[i].Y < v2.Y {
			if v.Y < vertices[i].Y || v.Y > v2.Y {
				return vertices, false
			}
		} else {
			if v.Y > vertices[i].Y || v.Y < v2.Y {
				return vertices, false
			}
		}
	} else {
		if vertices[i].X < v2.X {
			if v.X < vertices[i].X || v.X > v2.X {
				return vertices, false
			}
		} else {
			if v.X > vertices[i].X || v.X < v2.X {
				return vertices, false
			}
		}
	}

	// fmt.Printf("Inserting %v in bewteen %v and %v\n", v, vertices[i], vertices[(i+1)%len(vertices)])
	res := append([]math.Vector(nil), append(vertices[:i+1], append([]math.Vector{v}, vertices[i+1:]...)...)...)
	// fmt.Printf("> %v ->\n> %v\n", vertices, res)
	return res, true
}

//
//
//

// isReflex returns true if the vertex at index i forms a reflex angle
func isReflex(vertices []math.Vector, i int) bool {
	prev := vertices[(i-1+len(vertices))%len(vertices)]
	curr := vertices[i]
	next := vertices[(i+1)%len(vertices)]

	// Calculate cross product to determine if angle is reflex
	v1 := math.Vector{X: curr.X - prev.X, Y: curr.Y - prev.Y}
	v2 := math.Vector{X: next.X - curr.X, Y: next.Y - curr.Y}
	cross := v1.X*v2.Y - v1.Y*v2.X

	return cross < 0
}

// pointInTriangle returns true if point p lies inside the triangle formed by a, b, c
func pointInTriangle(p, a, b, c math.Vector) bool {
	// Using barycentric coordinates
	denominator := ((b.Y-c.Y)*(a.X-c.X) + (c.X-b.X)*(a.Y-c.Y))
	if stdmath.Abs(float64(denominator)) < 1e-10 {
		return false
	}

	alpha := ((b.Y-c.Y)*(p.X-c.X) + (c.X-b.X)*(p.Y-c.Y)) / denominator
	beta := ((c.Y-a.Y)*(p.X-c.X) + (a.X-c.X)*(p.Y-c.Y)) / denominator
	gamma := 1.0 - alpha - beta

	// Point is inside if all barycentric coordinates are between 0 and 1
	return alpha > -1e-10 && beta > -1e-10 && gamma > -1e-10
}

// lineIntersects returns true if line segments (a,b) and (c,d) intersect
func lineIntersects(a, b, c, d math.Vector) bool {
	// Calculate direction vectors
	r := math.Vector{X: b.X - a.X, Y: b.Y - a.Y}
	s := math.Vector{X: d.X - c.X, Y: d.Y - c.Y}

	// Calculate cross products
	denominator := r.X*s.Y - r.Y*s.X
	if stdmath.Abs(float64(denominator)) < 1e-10 {
		return false // Lines are parallel
	}

	// Calculate parameters of intersection point
	t := ((c.X-a.X)*s.Y - (c.Y-a.Y)*s.X) / denominator
	u := ((c.X-a.X)*r.Y - (c.Y-a.Y)*r.X) / denominator

	return t >= 0 && t <= 1 && u >= 0 && u <= 1
}

// diagonalIntersectsEdge returns true if the diagonal intersects any non-adjacent polygon edge
func diagonalIntersectsEdge(vertices []math.Vector, i, j int) bool {
	a := vertices[i]
	b := vertices[j]

	n := len(vertices)
	for k := 0; k < n; k++ {
		if k == i || k == j || k == (i-1+n)%n || k == (i+1)%n ||
			k == (j-1+n)%n || k == (j+1)%n {
			continue // Skip adjacent edges
		}

		if lineIntersects(a, b, vertices[k], vertices[(k+1)%n]) {
			return true
		}
	}
	return false
}

// isEar returns true if the vertex at index i forms a valid ear
func isEar(vertices []math.Vector, i int, reflexVertices map[int]bool) bool {
	if len(vertices) < 3 {
		return false
	}

	n := len(vertices)
	prev := vertices[(i-1+n)%n]
	curr := vertices[i]
	next := vertices[(i+1)%n]

	// Check if vertex is reflex
	if isReflex(vertices, i) {
		return false
	}

	// Check if diagonal intersects any non-adjacent edges
	if diagonalIntersectsEdge(vertices, (i-1+n)%n, (i+1)%n) {
		return false
	}

	// Check if any reflex vertex lies inside the potential ear triangle
	for j := range reflexVertices {
		if j >= len(vertices) {
			continue // Skip invalid indices
		}
		if j == i || j == (i-1+n)%n || j == (i+1)%n {
			continue
		}
		if pointInTriangle(vertices[j], prev, curr, next) {
			return false
		}
	}

	// Check triangle area is not degenerate
	area := stdmath.Abs(float64((prev.X*(curr.Y-next.Y) + curr.X*(next.Y-prev.Y) + next.X*(prev.Y-curr.Y)) / 2))
	return area > 1e-10
}

func triangulate3(bound Bound) []Bound {
	if len(bound.Vertices) < 3 {
		panic("Not enough vertices")
	}

	if len(bound.Vertices) == 3 {
		return []Bound{bound}
	}

	for i, v := range bound.Vertices {
		pi := (i - 1 + len(bound.Vertices)) % len(bound.Vertices)
		ni := (i + 1) % len(bound.Vertices)
		prev := bound.Vertices[pi]
		next := bound.Vertices[ni]

		if isConvex(prev, v, next) && !isColinear(prev, v, next) {
			bad := false
			for j, v2 := range bound.Vertices {
				if j == pi || j == i || j == ni {
					continue
				}

				if insideTriangle(prev, v, next, v2) {
					bad = true
				}
			}
			if bad {
				continue
			}

			// triangleArea := func(a, b, c math.Vector) float32 {
			// 	return math.Abs32((b.X-a.X)*(c.Y-a.Y)-(c.X-a.X)*(b.Y-a.Y)) / 2
			// }
			// fmt.Printf("> %.3f\n", triangleArea(prev, v, next))
			if prev.X == v.X && v.X == next.X {
				fmt.Printf("BAD 1\n")
			}
			if prev.Y == v.Y && v.Y == next.Y {
				fmt.Printf("BAD 2\n")
			}

			// Create triangle
			boundID++
			triangle := Bound{
				ID:       boundID,
				Vertices: []math.Vector{prev, v, next},
				Color:    bound.Color,
			}

			// Recursively triangulate remaining polygon
			remaining := append([]math.Vector(nil), append(bound.Vertices[:i], bound.Vertices[i+1:]...)...)
			boundID++
			triangles := triangulate3(Bound{
				ID:       boundID,
				Vertices: remaining,
				Color:    bound.Color,
			})

			return append(triangles, triangle)
		}
	}

	panic("Failed to triangulate")
}

func isColinear(a, b, c math.Vector) bool {
	return (a.X == b.X && b.X == c.X) || (a.Y == b.Y && b.Y == c.Y)
}

func insideTriangle(a, b, c, p math.Vector) bool {
	v0 := c.Sub(a)
	v1 := b.Sub(a)
	v2 := p.Sub(a)

	dot00 := v0.Dot(v0)
	dot01 := v0.Dot(v1)
	dot02 := v0.Dot(v2)
	dot11 := v1.Dot(v1)
	dot12 := v1.Dot(v2)

	denom := dot00*dot11 - dot01*dot01
	if stdmath.Abs(float64(denom)) < 1e-20 {
		return true
	}
	invDenom := 1.0 / denom
	u := (dot11*dot02 - dot01*dot12) * invDenom
	v := (dot00*dot12 - dot01*dot02) * invDenom

	return (u >= 0) && (v >= 0) && (u+v < 1)
}

// func cw(a, b, c math.Vector) bool {
// 	//
// }

func isConvex(a, b, c math.Vector) bool {
	angle := angle(c.Sub(b), a.Sub(b))
	if angle == 0 {
		return false
	}
	return angle <= stdmath.Pi
}

func angle(a, b math.Vector) float32 {
	// v1 := b.Sub(a)
	// v2 := b.Sub(c)

	// dot := v1.Dot(v2)
	// mag := v1.Len() * v2.Len()

	// return stdmath.Acos(float64(dot / mag))

	angle := math.Atan232(a.Cross(b), a.Dot(b))
	if angle < 0 {
		angle += 2 * stdmath.Pi
	}

	return angle
}

func triangulate2(bounds []Bound) []Bound {
	var triangles []Bound
	for _, b := range bounds {
		triangles = append(triangles, triangulate3(b)...)
	}

	return triangles
}

func triangulate(bounds []Bound) []Bound {
	if true {
		return triangulate2(bounds)
	}

	result := make([]Bound, 0)

	// Process each polygon separately
	for _, bound := range bounds {
		if len(bound.Vertices) < 3 {
			continue
		}

		// Create working copy of vertices
		vertices := make([]math.Vector, len(bound.Vertices))
		copy(vertices, bound.Vertices)

		// Initialize map of reflex vertices
		reflexVertices := make(map[int]bool)
		for i := 0; i < len(vertices); i++ {
			if isReflex(vertices, i) {
				reflexVertices[i] = true
			}
		}

		// Main ear clipping loop
		for len(vertices) > 3 {
			found := false

			// Find and clip an ear
			for i := 0; i < len(vertices); i++ {
				if isEar(vertices, i, reflexVertices) {
					// Create triangle from ear
					prev := vertices[(i-1+len(vertices))%len(vertices)]
					curr := vertices[i]
					next := vertices[(i+1)%len(vertices)]

					boundID++
					triangle := Bound{
						ID:       boundID,
						Vertices: []math.Vector{prev, curr, next},
						Color:    randomColor(),
					}
					result = append(result, triangle)

					// Remove ear vertex
					vertices = append(vertices[:i], vertices[i+1:]...)

					// Update reflex vertices map
					newReflexVertices := make(map[int]bool)
					for j := 0; j < len(vertices); j++ {
						if isReflex(vertices, j) {
							newReflexVertices[j] = true
						}
					}
					reflexVertices = newReflexVertices

					found = true
					break
				}
			}

			if !found {
				// If no ear is found, there might be numerical precision issues
				// Add the remaining polygon as a triangle if possible
				boundID++
				if len(vertices) == 3 {
					triangle := Bound{
						ID:       boundID,
						Vertices: []math.Vector{vertices[0], vertices[1], vertices[2]},
						Color:    randomColor(),
					}
					result = append(result, triangle)
				}
				break
			}
		}

		// Add final triangle
		if len(vertices) == 3 {
			boundID++
			triangle := Bound{
				ID:       boundID,
				Vertices: vertices,
				Color:    randomColor(),
			}
			result = append(result, triangle)
		}
	}

	return result
}

//
//
//

func segmentIntersection(b1, b2 *Bound, i, ii, j, jj int) bool {
	a := b1.Vertices[i]
	b := b1.Vertices[ii]
	c := b2.Vertices[j]
	d := b2.Vertices[jj]

	var ok1, ok2 bool

	ix := math.Min(i, ii)
	if ii == 0 {
		ix = i // wraparound
	}
	jx := math.Min(j, jj)
	if jj == 0 {
		jx = j // wraparound
	}

	// fmt.Printf("a %v b %v c %v d %v\n", a, b, c, d)

	if a.X == b.X {
		if c.X == d.X && a.X == c.X {
			ok := false
			for _, intersection := range []math.Vector{
				{X: a.X, Y: a.Y},
				{X: a.X, Y: b.Y},
				{X: a.X, Y: c.Y},
				{X: a.X, Y: d.Y},
			} {
				b1.Vertices, ok1 = insertVertex(b1.Vertices, ix, intersection)
				b2.Vertices, ok2 = insertVertex(b2.Vertices, jx, intersection)
				ok = ok || ok1 || ok2
			}

			return ok
		}
	}

	if a.Y == b.Y {
		if c.Y == d.Y && a.Y == c.Y {
			// minA := math.Min(a.X, b.X)
			// maxA := math.Max(a.X, b.X)
			// minC := math.Min(c.X, d.X)
			// maxC := math.Max(c.X, d.X)

			// if minA < maxC && maxA > minC {
			// 	// Intersection
			// 	intersection := math.Vector{X: math.Max(minA, minC), Y: a.Y}
			// 	b1.Vertices, ok1 = insertVertex(b1.Vertices, ix, intersection)
			// 	b2.Vertices, ok2 = insertVertex(b2.Vertices, jx, intersection)
			// 	return ok1 || ok2
			// }

			ok := false
			for _, intersection := range []math.Vector{
				{X: a.X, Y: a.Y},
				{X: b.X, Y: a.Y},
				{X: c.X, Y: a.Y},
				{X: d.X, Y: a.Y},
			} {
				b1.Vertices, ok1 = insertVertex(b1.Vertices, ix, intersection)
				b2.Vertices, ok2 = insertVertex(b2.Vertices, jx, intersection)
				ok = ok || ok1 || ok2
				if ok {
					break
				}
			}

			return ok
		}
	}

	// fmt.Printf("a %v b %v c %v d %v\n", a, b, c, d)
	panic("Not overlapping!")
}

func constructNavigationGraph(rooms []Room, walls []Door, doors []Door) *NavigationGraph {
	var edges []*NavigationEdge
	nodes := map[int]*NavigationNode{}

	type X struct{ i, j int }
	x := map[int][]X{}

	for i, room := range rooms {
		for _, bound := range room.Bounds {
			center := math.Vector{}
			for _, vertex := range bound.Vertices {
				center = center.Add(vertex)
			}
			center = center.Divs(float32(len(bound.Vertices)))

			node := &NavigationNode{
				X:     center.X,
				Y:     center.Y,
				Bound: bound,
			}

			nodes[bound.ID] = node
		}

		for j, bound := range room.Bounds {
		outer:
			for k, otherBound := range room.Bounds {
				if k <= j {
					continue
				}

				// fmt.Printf("%d %d %d %d\n", i, j, bound.ID, otherBound.ID)

				if i, ii, _, _, ok := boundsOverlap2(bound, otherBound); ok {
					for _, wall := range walls {
						if segmentsOverlap(i, ii, wall.Bound.Vertices[0], wall.Bound.Vertices[1]) {
							continue outer
						}
					}

					edges = append(edges, &NavigationEdge{
						From: bound.ID,
						To:   otherBound.ID,
					})
				}
			}

			for _, door := range doors {
				if boundsOverlap(bound, door.Bound) {
					x[door.Bound.ID] = append(x[door.Bound.ID], X{i, j})
				}
			}
		}
	}

	for _, xx := range x {
		// fmt.Printf("> %v\n", xx)

		for i, x1 := range xx {
			for j, x2 := range xx {
				if i <= j {
					continue
				}

				edges = append(edges, &NavigationEdge{
					From: rooms[x1.i].Bounds[x1.j].ID,
					To:   rooms[x2.i].Bounds[x2.j].ID,
				})
			}
		}
	}

	// for _, door := range doors {
	// 	center := math.Vector{}
	// 	for _, vertex := range door.Bound.Vertices {
	// 		center = center.Add(vertex)
	// 	}
	// 	center = center.Divs(float32(len(door.Bound.Vertices)))

	// 	node := &NavigationNode{
	// 		X:     center.X,
	// 		Y:     center.Y,
	// 		Bound: door.Bound,
	// 	}

	// 	nodes[door.Bound.ID] = node

	// 	for _, room := range rooms {
	// 		for _, bound := range room.Bounds {
	// 			if boundsOverlap(bound, door.Bound) {
	// 				edges = append(edges, &NavigationEdge{
	// 					From: bound.ID,
	// 					To:   door.Bound.ID,
	// 				})
	// 			}
	// 		}
	// 	}
	// }

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

func tiles(tileMap *TileMap, row, col int, visited map[int]any) ([]Point, [][]bool) {
	queue := []Point{{Row: row, Col: col}}
	visited[mapIndex(tileMap, row, col)] = struct{}{}

	board := make([][]bool, tileMap.Width())
	for i := range board {
		board[i] = make([]bool, tileMap.Height())
	}

	var points []Point
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		board[p.Row][p.Col] = true
		points = append(points, p)

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

	return points, board
}

func mapIndex(tileMap *TileMap, row, col int) int {
	return col*tileMap.Width() + row
}

//
//
//

func boundsOverlap2(a, b Bound) (math.Vector, math.Vector, math.Vector, math.Vector, bool) {
	for i := 0; i < len(a.Vertices); i++ {
		ii := (i + 1) % len(a.Vertices)

		for j := 0; j < len(b.Vertices); j++ {
			jj := (j + 1) % len(b.Vertices)

			if (a.Vertices[i].Equal(b.Vertices[j]) && a.Vertices[ii].Equal(b.Vertices[jj])) || (a.Vertices[i].Equal(b.Vertices[jj]) && a.Vertices[ii].Equal(b.Vertices[j])) {
				return a.Vertices[i], a.Vertices[ii], b.Vertices[j], b.Vertices[jj], true
			}
		}
	}

	return math.Vector{}, math.Vector{}, math.Vector{}, math.Vector{}, false
}

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
		if j.X == jj.X && i.X == j.X {
			minI := math.Min(i.Y, ii.Y)
			maxI := math.Max(i.Y, ii.Y)
			minJ := math.Min(j.Y, jj.Y)
			maxJ := math.Max(j.Y, jj.Y)

			return minI < maxJ && maxI > minJ
		}
	}

	if i.Y == ii.Y {
		if j.Y == jj.Y && i.Y == j.Y {
			minI := math.Min(i.X, ii.X)
			maxI := math.Max(i.X, ii.X)
			minJ := math.Min(j.X, jj.X)
			maxJ := math.Max(j.X, jj.X)

			return minI < maxJ && maxI > minJ
		}
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

//
//
//
//
//

func cornersToBound(corners []Corner) Bound {
	var vertices []math.Vector
	for _, corner := range corners {
		// corner.R is a row => in Y if row=0 at top
		// corner.C is a col => in X
		x := float32(corner.C) * 64
		y := float32(corner.R) * 64
		vertices = append(vertices, math.Vector{X: x, Y: y})
	}
	boundID++
	return Bound{
		ID:       boundID,
		Vertices: vertices,
		Color:    randomColor(),
	}
}

func tracePolygons(graph map[Corner]map[Corner]bool) [][]Corner {
	visitedEdges := make(map[[2]Corner]bool)
	var polygons [][]Corner

	for startCorner := range graph {
		for neighborCorner := range graph[startCorner] {
			edgeKey := makeEdgeKey2(startCorner, neighborCorner)
			if visitedEdges[edgeKey] {
				continue
			}
			// We found a new edge => start a boundary walk
			poly := walkBoundary(graph, startCorner, neighborCorner, visitedEdges)
			polygons = append(polygons, poly)
		}
	}

	return polygons
}

func walkBoundary(
	graph map[Corner]map[Corner]bool,
	start Corner,
	next Corner,
	visitedEdges map[[2]Corner]bool,
) []Corner {
	var result []Corner
	current := start

	prev := start // Might just track the previous corner
	for {
		// Mark both directions visited
		visitedEdges[makeEdgeKey2(current, next)] = true
		visitedEdges[makeEdgeKey2(next, current)] = true

		// Remove it from adjacency so we can't pick it again
		delete(graph[current], next)
		delete(graph[next], current)

		// Add 'current' to the polygon
		result = append(result, current)

		// Advance
		prev = current
		current = next

		// Pick the next corner
		n, ok := pickNextCorner(graph, current, visitedEdges, prev)
		if !ok {
			// no next => done
			break
		}

		// If we've looped around back to start edges, stop
		if current == start && n == result[0] {
			break
		}
		next = n
	}

	return result
}

func pickNextCorner(
	graph map[Corner]map[Corner]bool,
	current Corner,
	visitedEdges map[[2]Corner]bool,
	prev Corner,
) (Corner, bool) {
	// skip going back to 'prev'
	for neigh := range graph[current] {
		if neigh == prev {
			continue
		}
		if !visitedEdges[makeEdgeKey2(current, neigh)] {
			return neigh, true
		}
	}
	return Corner{}, false
}

func makeEdgeKey2(a, b Corner) [2]Corner {
	return [2]Corner{a, b}
}

func buildCornerGraph(edges map[EdgeKey]bool) map[Corner]map[Corner]bool {
	graph := make(map[Corner]map[Corner]bool)

	for ek := range edges {
		c1, c2 := edgeKeyToCorners(ek)
		// Insert the adjacency in *both* directions
		// if you prefer an undirected approach:
		addCornerEdge(graph, c1, c2)
		addCornerEdge(graph, c2, c1)
	}

	return graph
}

func addCornerEdge(graph map[Corner]map[Corner]bool, a, b Corner) {
	if graph[a] == nil {
		graph[a] = make(map[Corner]bool)
	}
	graph[a][b] = true
}

func edgeKeyToCorners(ek EdgeKey) (Corner, Corner) {
	// For example:
	// Dir = 0 => north edge: from (row, col) to (row, col+1)
	switch ek.Dir {
	case 0: // North
		return Corner{R: ek.Row, C: ek.Col}, Corner{R: ek.Row, C: ek.Col + 1}
	case 1: // East
		return Corner{R: ek.Row, C: ek.Col + 1}, Corner{R: ek.Row + 1, C: ek.Col + 1}
	case 2: // South
		return Corner{R: ek.Row + 1, C: ek.Col}, Corner{R: ek.Row + 1, C: ek.Col + 1}
	case 3: // West
		return Corner{R: ek.Row, C: ek.Col}, Corner{R: ek.Row + 1, C: ek.Col}
	}
	// fallback
	return Corner{}, Corner{}
}

type Corner struct {
	R, C int // corner in row-col corner space
}

const (
	N = 0
	E = 1
	S = 2
	W = 3
)

func edgeBlockedByWallOrDoor(tileMap *TileMap, row, col int, flag int) bool {
	switch flag {
	case N:
		return tileMap.GetBit(row, col, INTERIOR_WALL_N_BIT) || tileMap.GetBit(row, col, DOOR_N_BIT)
	case E:
		return tileMap.GetBit(row, col, INTERIOR_WALL_E_BIT) || tileMap.GetBit(row, col, DOOR_E_BIT)
	case S:
		return tileMap.GetBit(row, col, INTERIOR_WALL_S_BIT) || tileMap.GetBit(row, col, DOOR_S_BIT)
	case W:
		return tileMap.GetBit(row, col, INTERIOR_WALL_W_BIT) || tileMap.GetBit(row, col, DOOR_W_BIT)
	}

	panic("bad flag")
}

func buildBoundaryEdges(tileMap *TileMap, cells []Point) map[EdgeKey]bool {
	edges := make(map[EdgeKey]bool)

	// Convert slice of Points to a quick lookup set for membership
	cellSet := make(map[Point]bool, len(cells))
	for _, c := range cells {
		cellSet[c] = true
	}

	for _, c := range cells {
		r, co := c.Row, c.Col

		// Each cell has 4 potential edges: top, right, bottom, left.
		// We'll check if crossing that edge would remain inside
		// this same floor region or not. If not, it's a boundary edge.

		// Top edge => between (r,c) and (r-1,c)
		if r == 0 || !cellSet[Point{Row: r - 1, Col: co}] ||
			edgeBlockedByWallOrDoor(tileMap, r, co /* which edge? */, N) {
			ek := EdgeKey{Row: r, Col: co, Dir: 0} // north
			edges[ek] = true
		}

		// Right edge => between (r,c) and (r,c+1)
		if co == tileMap.Height()-1 || !cellSet[Point{Row: r, Col: co + 1}] ||
			edgeBlockedByWallOrDoor(tileMap, r, co /* which edge? */, E) {
			ek := EdgeKey{Row: r, Col: co, Dir: 1} // east
			edges[ek] = true
		}

		// Bottom edge => between (r,c) and (r+1,c)
		if r == tileMap.Width()-1 || !cellSet[Point{Row: r + 1, Col: co}] ||
			edgeBlockedByWallOrDoor(tileMap, r, co /* which edge? */, S) {
			ek := EdgeKey{Row: r, Col: co, Dir: 2} // south
			edges[ek] = true
		}

		// Left edge => between (r,c) and (r,c-1)
		if co == 0 || !cellSet[Point{Row: r, Col: co - 1}] ||
			edgeBlockedByWallOrDoor(tileMap, r, co /* which edge? */, W) {
			ek := EdgeKey{Row: r, Col: co, Dir: 3} // west
			edges[ek] = true
		}
	}

	return edges
}

// We'll define a small struct to keep track of a "directed edge"
// in the grid boundary graph.
type EdgeKey struct {
	Row, Col int
	Dir      int // 0 = north edge, 1 = east, 2 = south, 3 = west, for example
}
