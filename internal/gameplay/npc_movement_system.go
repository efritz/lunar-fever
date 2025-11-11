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
			pathfindingComponent.Target = &target
		}

		if pathfindingComponent.Target != nil {
			var from, to maps.Bound
			var minFromDist, minToDist float32
			minFromDist, minToDist = stdmath.MaxFloat32, stdmath.MaxFloat32

			for _, room := range s.Base.Rooms {
				for _, bound := range room.Bounds {
					// Ensure the bound is a triangle
					if len(bound.Vertices) != 3 {
						continue
					}

					// Calculate distance from the body's position to this triangle
					fromDist := pointToTriangleDistance(physicsComponent.Body.Position, bound.Vertices[0], bound.Vertices[1], bound.Vertices[2])
					if fromDist < minFromDist {
						minFromDist = fromDist
						from = bound
					}

					// Calculate distance from the target to this triangle
					toDist := pointToTriangleDistance(*pathfindingComponent.Target, bound.Vertices[0], bound.Vertices[1], bound.Vertices[2])
					if toDist < minToDist {
						minToDist = toDist
						to = bound
					}
				}
			}

			if true {
				path := smoothPath(s.Base.NavigationGraph, search(s.Base.NavigationGraph, from.ID, to.ID), physicsComponent.Body.Position, *pathfindingComponent.Target)
				pathfindingComponent.Waypoints = path[1:]
			}
		} else {
			pathfindingComponent.Waypoints = nil
		}

		const (
			desiredSpeed  = float32(0.35)
			maxSpeed      = float32(0.35)
			maxForce      = float32(2.0)
			slowingRadius = float32(80.0)
			maxTurnRate   = float32(6.0)
		)

		dt := float32(elapsedMs) / 1000.0

		if len(pathfindingComponent.Waypoints) > 0 {
			manager := NewSteeringManager(maxForce)

			if len(pathfindingComponent.Waypoints) == 1 {
				manager.AddBehavior(&Arrival{
					Position:      &physicsComponent.Body.Position,
					Velocity:      &physicsComponent.Body.LinearVelocity,
					Target:        pathfindingComponent.Waypoints[0],
					DesiredSpeed:  desiredSpeed,
					SlowingRadius: slowingRadius,
				}, 1.0)
			} else {
				manager.AddBehavior(&Seek{
					Position:     &physicsComponent.Body.Position,
					Velocity:     &physicsComponent.Body.LinearVelocity,
					Target:       pathfindingComponent.Waypoints[0],
					DesiredSpeed: desiredSpeed,
				}, 1.0)
			}

			physicsComponent.Body.LinearVelocity = physicsComponent.Body.LinearVelocity.Add(manager.Calculate().Muls(dt))

			if speed := physicsComponent.Body.LinearVelocity.Len(); speed > maxSpeed {
				physicsComponent.Body.LinearVelocity = physicsComponent.Body.LinearVelocity.Normalize().Muls(maxSpeed)
			}

			if v := physicsComponent.Body.LinearVelocity; v.Len() > 0.0001 {
				targetAngle := math.Atan232(v.Y, v.X)
				if targetAngle < 0 {
					targetAngle = (2 * stdmath.Pi) - (-targetAngle)
				}
				targetAngle -= float32(stdmath.Pi / 2)

				maxDelta := maxTurnRate * dt
				current := physicsComponent.Body.Orient
				delta := shortestAngleDiff(current, targetAngle)
				if delta > maxDelta {
					delta = maxDelta
				} else if delta < -maxDelta {
					delta = -maxDelta
				}
				newAngle := normalizeAngle(current + delta)
				if physicsComponent.Body.Orient != newAngle {
					physicsComponent.Body.SetOrient(newAngle)
				}
			}

			if len(pathfindingComponent.Waypoints) == 1 {
				if dist := pathfindingComponent.Waypoints[0].Sub(physicsComponent.Body.Position).Len(); dist < 30 {
					pathfindingComponent.Target = nil
				}
			}
		} else {
			physicsComponent.Body.LinearVelocity = math.Vector{0, 0}
		}
	}
}

func pointToTriangleDistance(point, a, b, c math.Vector) float32 {
	if maps.PointInTriangle(a, b, c, point) {
		return 0
	}

	// Calculate distances to each edge
	edgeDist1 := pointToSegmentDistance(point, a, b)
	edgeDist2 := pointToSegmentDistance(point, b, c)
	edgeDist3 := pointToSegmentDistance(point, c, a)

	// Calculate distances to each vertex
	vertexDist1 := point.Sub(a).Len()
	vertexDist2 := point.Sub(b).Len()
	vertexDist3 := point.Sub(c).Len()

	// Return the smallest distance
	return math.Min(math.Min(edgeDist1, edgeDist2), math.Min(edgeDist3, math.Min(vertexDist1, math.Min(vertexDist2, vertexDist3))))
}

func normalizeAngle(a float32) float32 {
	for a <= 0 {
		a += 2 * float32(stdmath.Pi)
	}
	for a > 2*float32(stdmath.Pi) {
		a -= 2 * float32(stdmath.Pi)
	}
	return a
}

func shortestAngleDiff(from, to float32) float32 {
	diff := normalizeAngle(to) - normalizeAngle(from)

	if diff > float32(stdmath.Pi) {
		diff -= 2 * float32(stdmath.Pi)
	} else if diff < -float32(stdmath.Pi) {
		diff += 2 * float32(stdmath.Pi)
	}

	return diff
}

func pointToSegmentDistance(p, a, b math.Vector) float32 {
	ab := b.Sub(a)
	ap := p.Sub(a)

	// Project p onto ab, but limit to [0,1] to stay on the segment
	t := ap.Dot(ab) / ab.Dot(ab)
	t = math.Max(0, math.Min(1, t))

	// Find the closest point on the segment to p
	closest := a.Add(ab.Muls(t))

	// Return the distance to the closest point
	return p.Sub(closest).Len()
}
