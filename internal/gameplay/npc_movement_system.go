package gameplay

import (
	stdmath "math"

	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine/ecs/system"
	"github.com/efritz/lunar-fever/internal/gameplay/maps"
)

type npcMovementSystem struct {
	*GameContext
}

func NewNpcMovementSystem(ctx *GameContext) system.System {
	return &npcMovementSystem{GameContext: ctx}
}

func (s *npcMovementSystem) Init() {}
func (s *npcMovementSystem) Exit() {}

func (s *npcMovementSystem) Process(elapsedMs int64) {
	mx := s.Camera.Unprojectx(float32(s.Mouse.X()))
	my := s.Camera.UnprojectY(float32(s.Mouse.Y()))

	for _, entity := range s.NpcCollection.Entities() {
		physicsComponent, ok := s.PhysicsComponentManager.GetComponent(entity)
		if !ok {
			continue
		}

		pathfindingComponent, ok := s.PathfindingComponentManager.GetComponent(entity)
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

			pathfindingComponent.Target = SmoothPath(s.Base.NavigationGraph, search(s.Base.NavigationGraph, from.ID, to.ID), physicsComponent.Body.Position, target)
		}

		mod := float32(1000)
		speed := float32(.35)
		transitionSpeed := float32(4)

		if len(pathfindingComponent.Target) > 0 {
			angle := math.Atan232(pathfindingComponent.Target[0].Y-physicsComponent.Body.Position.Y, pathfindingComponent.Target[0].X-physicsComponent.Body.Position.X)
			if angle < 0 {
				angle = (2 * stdmath.Pi) - (-angle)
			}
			angle -= float32(stdmath.Pi / 2)

			if physicsComponent.Body.Orient != angle {
				physicsComponent.Body.SetOrient(angle)
			}

			physicsComponent.Body.LinearVelocity =
				physicsComponent.Body.LinearVelocity.Muls(1 - (float32(elapsedMs) / mod * transitionSpeed)).Add(
					pathfindingComponent.Target[0].Sub(physicsComponent.Body.Position).Normalize().Muls(speed * float32(elapsedMs) / mod * transitionSpeed),
				)

			if pathfindingComponent.Target[0].Sub(physicsComponent.Body.Position).Len() < 30 {
				pathfindingComponent.Target = pathfindingComponent.Target[1:]
			}
		} else {
			physicsComponent.Body.LinearVelocity = math.Vector{0, 0}
		}
	}
}

func contains(bounds maps.Bound, point math.Vector) bool {
	if len(bounds.Vertices) != 3 {
		return false
	}

	return maps.PointInTriangle(bounds.Vertices[0], bounds.Vertices[1], bounds.Vertices[2], point)
}
