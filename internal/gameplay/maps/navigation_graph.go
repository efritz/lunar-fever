package maps

import "github.com/efritz/lunar-fever/internal/common/math"

type NavigationGraph struct {
	Nodes map[int]*NavigationNode
	Edges []*NavigationEdge
}

type NavigationNode struct {
	Bound  Bound
	Center math.Vector
}

func newNavigationNode(bound Bound) *NavigationNode {
	center := math.Vector{}
	for _, vertex := range bound.Vertices {
		center = center.Add(vertex)
	}

	return &NavigationNode{
		Bound:  bound,
		Center: center.Divs(float32(len(bound.Vertices))),
	}
}

type NavigationEdge struct {
	From int
	To   int
}

// constructNavigationGraph returns a graph where each node is a unique bound and each edge
// denotes two bounds that share an edge without an obstacle between them.
func constructNavigationGraph(rooms []Room, walls []Edge, doors []Edge) *NavigationGraph {
	nodes := map[int]*NavigationNode{}
	for _, room := range rooms {
		for _, bound := range room.Bounds {
			nodes[bound.ID] = newNavigationNode(bound)
		}
	}

	return &NavigationGraph{
		Nodes: nodes,
		Edges: append(
			findAdjacentBoundsWithinSameRoom(rooms, walls),
			findAjacentBoundsConnectedByDoor(rooms, doors)...,
		),
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
					if (v1.Equal(wall.From) && v2.Equal(wall.To)) || (v1.Equal(wall.To) && v2.Equal(wall.From)) {
						return false
					}
				}

				return true
			}
		}
	}

	return false
}

func findAjacentBoundsConnectedByDoor(rooms []Room, doors []Edge) []*NavigationEdge {
	type indexPair struct {
		roomIndex  int
		boundIndex int
	}
	overlappingBoundsByDoorIndex := map[int][]indexPair{}

	for i, room := range rooms {
		for j, bound := range room.Bounds {
			for k, door := range doors {
				if edgeExistsOnBound(bound, door) {
					overlappingBoundsByDoorIndex[k] = append(overlappingBoundsByDoorIndex[k], indexPair{
						roomIndex:  i,
						boundIndex: j,
					})
				}
			}
		}
	}

	var edges []*NavigationEdge
	for _, overlappingBounds := range overlappingBoundsByDoorIndex {
		for i := 0; i < len(overlappingBounds); i++ {
			for j := i + 1; j < len(overlappingBounds); j++ {
				b1 := overlappingBounds[i]
				b2 := overlappingBounds[j]

				edges = append(edges, &NavigationEdge{
					From: rooms[b1.roomIndex].Bounds[b1.boundIndex].ID,
					To:   rooms[b2.roomIndex].Bounds[b2.boundIndex].ID,
				})
			}
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
