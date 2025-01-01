package gameplay

import (
	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/gameplay/maps"
)

type PathfindingComponent struct {
	Target []math.Vector
}

type PathfindingComponentType struct{}

var pathfindingComponentType = PathfindingComponentType{}

func (c *PathfindingComponent) ComponentType() PathfindingComponentType {
	return pathfindingComponentType
}

//
//
//

func contains(bounds maps.Bound, point math.Vector) bool {
	vertices := bounds.Vertices
	n := len(vertices)
	if n < 3 {
		return false
	}

	inside := false
	for i, j := 0, n-1; i < n; i, j = i+1, i {
		if (vertices[i].Y > point.Y) != (vertices[j].Y > point.Y) &&
			point.X < (vertices[j].X-vertices[i].X)*(point.Y-vertices[i].Y)/(vertices[j].Y-vertices[i].Y)+vertices[i].X {
			inside = !inside
		}
	}

	return inside
}

type NodeInfo struct {
	id     int
	g      float32 // Cost from start to this node
	h      float32 // Heuristic estimate to goal
	f      float32 // f = g + h
	parent int     // Parent node in the path
}

func search(navigationGraph *maps.NavigationGraph, from, to int) []int {
	openSet := make(map[int]*NodeInfo)
	closedSet := make(map[int]bool)
	allNodes := make(map[int]*NodeInfo)

	// Initialize start node with g=0 and h=1 (just counting nodes)
	startNode := &NodeInfo{
		id:     from,
		g:      0,
		h:      1, // One step away is cost of 1
		parent: -1,
	}
	startNode.f = startNode.g + startNode.h
	openSet[from] = startNode
	allNodes[from] = startNode

	for len(openSet) > 0 {
		// Find node with lowest f score in open set
		var current *NodeInfo
		var currentId int
		for id, node := range openSet {
			if current == nil || node.f < current.f {
				current = node
				currentId = id
			}
		}

		// If we reached the goal, reconstruct and return the path
		if currentId == to {
			return reconstructPath(current, allNodes)
		}

		// Move current node from open to closed set
		delete(openSet, currentId)
		closedSet[currentId] = true

		// Check all neighbors through edges
		for _, edge := range navigationGraph.Edges {
			var neighborId int

			if edge.From == currentId {
				neighborId = edge.To
			} else if edge.To == currentId {
				neighborId = edge.From
			} else {
				continue
			}

			if closedSet[neighborId] {
				continue
			}

			// Cost to reach neighbor is current cost plus 1
			tentativeG := current.g + 1

			neighbor, exists := openSet[neighborId]
			if !exists {
				// New node discovered
				neighbor = &NodeInfo{
					id:     neighborId,
					g:      tentativeG,
					h:      1, // Always 1 step cost estimate
					parent: currentId,
				}
				neighbor.f = neighbor.g + neighbor.h
				openSet[neighborId] = neighbor
				allNodes[neighborId] = neighbor
			} else if tentativeG < neighbor.g {
				// Found a better path to neighbor
				neighbor.g = tentativeG
				neighbor.f = tentativeG + neighbor.h
				neighbor.parent = currentId
			}
		}
	}

	return nil
}

func reconstructPath(goal *NodeInfo, nodes map[int]*NodeInfo) []int {
	path := []int{goal.id}
	current := goal

	for current.parent != -1 {
		path = append([]int{current.parent}, path...)
		current = nodes[current.parent]
	}

	return path
}
