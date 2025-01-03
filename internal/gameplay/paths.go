package gameplay

import (
	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/gameplay/maps"
)

type PathfindingComponent struct {
	Portals    []Portal
	Path       []math.Vector
	TargetCopy []math.Vector
	Target     []math.Vector
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
//

func SmoothPath(navigationGraph *maps.NavigationGraph, path []int, startPoint, endPoint math.Vector) []math.Vector {
	var points []math.Vector
	points = append(points, startPoint)

	for _, id := range path {
		points = append(points, navigationGraph.Nodes[id].Center)
	}

	return append(points, endPoint)
}

func constructPortals(navigationGraph *maps.NavigationGraph, path []int, startPoint, endPoint math.Vector) []Portal {
	return nil
}

type Portal struct {
	Left  math.Vector
	Right math.Vector
}
