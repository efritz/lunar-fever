package gameplay

import (
	stdmath "math"

	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine/ecs/system"
	"github.com/efritz/lunar-fever/internal/gameplay/maps"
)

type npcMovementSystem struct {
	*GameContext
	target []math.Vector
}

func NewNpcMovementSystem(ctx *GameContext) system.System {
	return &npcMovementSystem{GameContext: ctx}
}

func (s *npcMovementSystem) Init() {}
func (s *npcMovementSystem) Exit() {}

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

func (s *npcMovementSystem) Process(elapsedMs int64) {
	mx := s.Camera.Unprojectx(float32(s.Mouse.X()))
	my := s.Camera.UnprojectY(float32(s.Mouse.Y()))

	// TODO - separate for each entity
	for _, entity := range s.NpcCollection.Entities() {
		physicsComponent, ok := s.PhysicsComponentManager.GetComponent(entity)
		if !ok {
			continue
		}

		if s.Mouse.LeftButtonNewlyDown() {
			target := math.Vector{mx, my}

			var from, to maps.Bound
			for _, room := range s.Base.Rooms {
				for _, bound := range room.Bounds {
					if contains(bound, physicsComponent.Body.Position) {
						from = bound
					}

					if contains(bound, target) {
						to = bound
					}
				}
			}

			s.target = nil
			for _, id := range search(s.Base.NavigationGraph, from.ID, to.ID) {
				s.target = append(s.target, math.Vector{s.Base.NavigationGraph.Nodes[id].X, s.Base.NavigationGraph.Nodes[id].Y})
			}
		}

		mod := float32(1000)
		speed := float32(.35)
		transitionSpeed := float32(4)

		if len(s.target) > 0 {
			angle := math.Atan232(s.target[0].Y-physicsComponent.Body.Position.Y, s.target[0].X-physicsComponent.Body.Position.X)
			if angle < 0 {
				angle = (2 * stdmath.Pi) - (-angle)
			}
			angle -= float32(stdmath.Pi / 2)

			if physicsComponent.Body.Orient != angle {
				physicsComponent.Body.SetOrient(angle)
			}

			physicsComponent.Body.LinearVelocity =
				physicsComponent.Body.LinearVelocity.Muls(1 - (float32(elapsedMs) / mod * transitionSpeed)).Add(
					s.target[0].Sub(physicsComponent.Body.Position).Normalize().Muls(speed * float32(elapsedMs) / mod * transitionSpeed),
				)

			if s.target[0].Sub(physicsComponent.Body.Position).Len() < 2 {
				s.target = s.target[1:]
			}
		} else {
			physicsComponent.Body.LinearVelocity = math.Vector{0, 0}
		}
	}
}
