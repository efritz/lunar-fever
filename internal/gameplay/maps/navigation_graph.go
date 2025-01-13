package maps

import "github.com/efritz/lunar-fever/internal/common/math"

type NavigationGraph struct {
	Nodes     map[int]*NavigationNode
	Edges     []*NavigationEdge
	Obstacles []Edge
}

type NavigationNode struct {
	Door   bool
	Bound  Bound
	Center math.Vector
}

func newNavigationNode(bound Bound, door bool) *NavigationNode {
	center := math.Vector{}
	for _, vertex := range bound.Vertices {
		center = center.Add(vertex)
	}

	return &NavigationNode{
		Door:   door,
		Bound:  bound,
		Center: center.Divs(float32(len(bound.Vertices))),
	}
}

type NavigationEdge struct {
	From int
	To   int
}

type doorBound struct {
	edge  Edge
	bound Bound
}

// constructNavigationGraph returns a graph where each node is a unique bound and each edge
// denotes two bounds that share an edge without an obstacle between them.
func constructNavigationGraph(rooms []Room, walls []Edge, doors []Edge) *NavigationGraph {
	nodes := map[int]*NavigationNode{}
	for _, room := range rooms {
		for _, bound := range room.Bounds {
			nodes[bound.ID] = newNavigationNode(bound, false)
		}
	}

	var doorBounds []doorBound
	for _, door := range doors {
		bound := expandObstacleEdge(door)
		doorBounds = append(doorBounds, doorBound{door, bound})
		nodes[bound.ID] = newNavigationNode(bound, true)
	}

	return &NavigationGraph{
		Nodes: nodes,
		Edges: append(
			findAdjacentBoundsWithinSameRoom(rooms, walls),
			findAjacentBoundsConnectedByDoor(rooms, doorBounds)...,
		),
		Obstacles: walls,
	}
}

func findAdjacentBoundsWithinSameRoom(rooms []Room, walls []Edge) []*NavigationEdge {
	var edges []*NavigationEdge
	for _, room := range rooms {
		for i := 0; i < len(room.Bounds); i++ {
			for j := i + 1; j < len(room.Bounds); j++ {
				b1 := room.Bounds[i]
				b2 := room.Bounds[j]

				if boundsShareFreeEdge(b1, b2, walls) {
					edges = append(edges, &NavigationEdge{
						From: b1.ID,
						To:   b2.ID,
					})
				}
			}
		}
	}

	return edges
}

func boundsShareFreeEdge(a, b Bound, walls []Edge) bool {
	n := len(a.Vertices)
	m := len(b.Vertices)

	for i := 0; i < len(a.Vertices); i++ {
		v1 := a.Vertices[i]
		v2 := a.Vertices[nextVertexIndex(i, n)]

		for j := 0; j < len(b.Vertices); j++ {
			v3 := b.Vertices[j]
			v4 := b.Vertices[nextVertexIndex(j, m)]

			if (v1.Equal(v3) && v2.Equal(v4)) || (v1.Equal(v4) && v2.Equal(v3)) {
				for _, wall := range walls {
					if axisAlignedSegmentsOverlap(v1, v2, wall.From, wall.To) {
						return false
					}
				}

				return true
			}
		}
	}

	return false
}

func findAjacentBoundsConnectedByDoor(rooms []Room, doorBounds []doorBound) []*NavigationEdge {
	type indexPair struct {
		roomIndex  int
		boundIndex int
	}
	overlappingBoundsByDoorIndex := map[int][]indexPair{}

	for i, room := range rooms {
		for j, bound := range room.Bounds {
			for k, doorBound := range doorBounds {
				if boundsShareFreeEdge(bound, doorBound.bound, nil) {
					// if edgeExistsOnBound(bound, doorBound.edge) {
					overlappingBoundsByDoorIndex[k] = append(overlappingBoundsByDoorIndex[k], indexPair{
						roomIndex:  i,
						boundIndex: j,
					})
				}
			}
		}
	}

	var edges []*NavigationEdge
	for k, overlappingBounds := range overlappingBoundsByDoorIndex {
		for i := 0; i < len(overlappingBounds); i++ {
			b1 := overlappingBounds[i]

			edges = append(
				edges,
				&NavigationEdge{
					From: doorBounds[k].bound.ID,
					To:   rooms[b1.roomIndex].Bounds[b1.boundIndex].ID,
				},
			)
		}
	}

	return edges
}

func edgeExistsOnBound(a Bound, b Edge) bool {
	n := len(a.Vertices)

	for i, v1 := range a.Vertices {
		v2 := a.Vertices[nextVertexIndex(i, n)]

		if (v1.Equal(b.From) && v2.Equal(b.To)) || (v1.Equal(b.To) && v2.Equal(b.From)) {
			return true
		}
	}

	return false
}
