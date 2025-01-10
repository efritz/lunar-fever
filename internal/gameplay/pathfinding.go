package gameplay

import (
	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/gameplay/maps"
)

type PathfindingComponent struct {
	// NextWaypoint []math.Vector
	Target    *math.Vector
	Waypoints []math.Vector
}

type PathfindingComponentType struct{}

var pathfindingComponentType = PathfindingComponentType{}

func (c *PathfindingComponent) ComponentType() PathfindingComponentType {
	return pathfindingComponentType
}

//
//
//

type nodeInfo struct {
	id     int
	parent int     // Parent node in the path
	g      float32 // Cost from start to this node
	h      float32 // Heuristic estimate to goal
}

const travelCost = 1

func search(navigationGraph *maps.NavigationGraph, from, to int) []int {
	edges := map[int][]int{}
	for _, edge := range navigationGraph.Edges {
		edges[edge.From] = append(edges[edge.From], edge.To)
		edges[edge.To] = append(edges[edge.To], edge.From)
	}

	openSet := map[int]*nodeInfo{}
	closedSet := map[int]any{}
	allNodes := map[int]*nodeInfo{}

	// Initialize open set with start node
	startNode := &nodeInfo{id: from, parent: -1, g: 0, h: travelCost}
	openSet[from] = startNode
	allNodes[from] = startNode

	for len(openSet) > 0 {
		// Find lowest cost node in open set
		var current *nodeInfo
		var currentID int
		for id, node := range openSet {
			if current == nil || node.g+node.h < current.g+current.h {
				current = node
				currentID = id
			}
		}

		if currentID == to {
			// Found goal, reconstruct path
			return reconstructPath(current, allNodes)
		}

		// Move from open set to closed set
		delete(openSet, currentID)
		closedSet[currentID] = struct{}{}

		// Expand neighbors
		for _, neighborID := range edges[currentID] {
			if _, ok := closedSet[neighborID]; ok {
				continue
			}

			tentativeG := current.g + travelCost

			neighbor, exists := openSet[neighborID]
			if !exists {
				// Discovered neighbor for first time
				neighbor = &nodeInfo{id: neighborID, parent: currentID, g: tentativeG, h: travelCost}
				openSet[neighborID] = neighbor
				allNodes[neighborID] = neighbor
			} else if tentativeG < neighbor.g {
				// Discovered better path for existing neighbor
				neighbor.g = tentativeG
				neighbor.parent = currentID
			}
		}
	}

	// No path exists
	return nil
}

func reconstructPath(goal *nodeInfo, nodes map[int]*nodeInfo) []int {
	path := []int{goal.id}
	for current := goal; current.parent != -1; current = nodes[current.parent] {
		path = append([]int{current.parent}, path...)
	}

	return path
}

//
//

func smoothPath(navigationGraph *maps.NavigationGraph, path []int, startPoint, endPoint math.Vector) []math.Vector {
	waypoints := smoothPathsBetweenDoors(navigationGraph, path, startPoint, endPoint)
	// for i, waypoint := range waypoints {
	// 	var collisions []maps.Edge
	// 	for _, obstacle := range navigationGraph.Obstacles {
	// 		if waypoint.Equal(obstacle.From) || waypoint.Equal(obstacle.To) {
	// 			collisions = append(collisions, obstacle)
	// 		}
	// 	}

	// 	if len(collisions) > 0 {
	// 		var totalNormal math.Vector
	// 		for _, obstacle := range collisions {
	// 			if obstacle.From.X == obstacle.To.X {
	// 				if waypoint.Y == math.Min(obstacle.From.Y, obstacle.To.Y) {
	// 					totalNormal = totalNormal.Add(math.Vector{0, -1})
	// 				} else {
	// 					totalNormal = totalNormal.Add(math.Vector{0, +1})
	// 				}
	// 			} else {
	// 				if waypoint.X == math.Min(obstacle.From.X, obstacle.To.X) {
	// 					totalNormal = totalNormal.Add(math.Vector{-1, 0})
	// 				} else {
	// 					totalNormal = totalNormal.Add(math.Vector{+1, 0})
	// 				}
	// 			}
	// 		}

	// 		waypoints[i] = waypoint.Add(totalNormal.Normalize().Muls(offsetDistance))
	// 	}
	// }

	return waypoints
}

func smoothPathsBetweenDoors(navigationGraph *maps.NavigationGraph, path []int, startPoint, endPoint math.Vector) []math.Vector {
	for i, v := range path {
		if navigationGraph.Nodes[v].Door {
			first := smoothPathsBetweenDoors(navigationGraph, path[:i], startPoint, navigationGraph.Nodes[v].Center)
			second := smoothPathsBetweenDoors(navigationGraph, path[i+1:], navigationGraph.Nodes[v].Center, endPoint)
			return append(first, second...)
		}
	}

	return smoothPathSegment(navigationGraph, path, startPoint, endPoint)
}

func smoothPathSegment(navigationGraph *maps.NavigationGraph, path []int, startPoint, endPoint math.Vector) []math.Vector {
	portals := constructPortals(navigationGraph, path, startPoint, endPoint)
	if len(portals) == 0 {
		return []math.Vector{startPoint, endPoint}
	}

	// Initialize the funnel
	var points []math.Vector
	points = append(points, startPoint)

	// Initialize funnel state
	apex := startPoint
	left := startPoint
	right := startPoint
	apexIndex := 0
	leftIndex := 0
	rightIndex := 0

	// Process each portal
	for i := 1; i < len(portals); i++ {
		portalLeft := portals[i].Left
		portalRight := portals[i].Right

		// Update right vertex
		if triarea2(apex, right, portalRight) <= 0 {
			if apex.Equal(right) || triarea2(apex, left, portalRight) > 0 {
				// Tighten the funnel
				right = portalRight
				rightIndex = i
			} else {
				// Right over left, insert left to path and restart scan from portal left point
				points = append(points, left)
				apex = left
				apexIndex = leftIndex
				left = apex
				right = apex
				leftIndex = apexIndex
				rightIndex = apexIndex
				i = apexIndex
				continue
			}
		}

		// Update left vertex
		if triarea2(apex, left, portalLeft) >= 0 {
			if apex.Equal(left) || triarea2(apex, right, portalLeft) < 0 {
				// Tighten the funnel
				left = portalLeft
				leftIndex = i
			} else {
				// Left over right, insert right to path and restart scan from portal right point
				points = append(points, right)
				apex = right
				apexIndex = rightIndex
				left = apex
				right = apex
				leftIndex = apexIndex
				rightIndex = apexIndex
				i = apexIndex
				continue
			}
		}
	}

	// Add the end point
	points = append(points, endPoint)
	return points
}

func constructPortals(navigationGraph *maps.NavigationGraph, path []int, startPoint, endPoint math.Vector) []Portal {
	if len(path) == 0 {
		return nil
	}

	// Initialize portals list with start portal
	portals := []Portal{{
		Left:  startPoint,
		Right: startPoint,
	}}

	// Find shared edges between consecutive triangles
	for i := 0; i < len(path)-1; i++ {
		currentNode := navigationGraph.Nodes[path[i]]
		nextNode := navigationGraph.Nodes[path[i+1]]

		// Find the shared edge between these nodes
		sharedEdge := findSharedEdge(currentNode.Bound.Vertices, nextNode.Bound.Vertices)
		if sharedEdge == nil {
			continue
		}

		// Orient the edge relative to the path direction
		left, right := orientPortalPoints(sharedEdge[0], sharedEdge[1], currentNode.Center, nextNode.Center)
		portals = append(portals, Portal{
			Left:  left,
			Right: right,
		})
	}

	// Add end portal
	portals = append(portals, Portal{
		Left:  endPoint,
		Right: endPoint,
	})

	return portals
}

// findSharedEdge finds the common edge between two triangles
func findSharedEdge(tri1, tri2 []math.Vector) []math.Vector {
	if len(tri1) != 3 || len(tri2) != 3 {
		panic("findSharedEdge called with non-triangle")
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if (tri1[i].Equal(tri2[j]) && tri1[(i+1)%3].Equal(tri2[(j+1)%3])) ||
				(tri1[i].Equal(tri2[(j+1)%3]) && tri1[(i+1)%3].Equal(tri2[j])) {
				return []math.Vector{tri1[i], tri1[(i+1)%3]}
			}
		}
	}
	return nil
}

// orientPortalPoints ensures the portal points are oriented correctly relative to the path direction
func orientPortalPoints(p1, p2, from, to math.Vector) (left, right math.Vector) {
	// Use the cross product to determine which point should be left/right
	path := to.Sub(from)
	edge := p2.Sub(p1)
	cross := path.X*edge.Y - path.Y*edge.X

	if cross < 0 {
		return p1, p2
	}
	return p2, p1
}

// triarea2 returns twice the signed area of the triangle abc
// positive if ccw, negative if cw, 0 if collinear
func triarea2(a, b, c math.Vector) float32 {
	return (c.X-a.X)*(b.Y-a.Y) - (b.X-a.X)*(c.Y-a.Y)
}

type Portal struct {
	Left  math.Vector
	Right math.Vector
}
