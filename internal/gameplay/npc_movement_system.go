package gameplay

import (
	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine/ecs/system"
)

type npcMovementSystem struct {
	*GameContext
	target *math.Vector
}

func NewNpcMovementSystem(ctx *GameContext) system.System {
	return &npcMovementSystem{GameContext: ctx}
}

func (s *npcMovementSystem) Init() {}
func (s *npcMovementSystem) Exit() {}

func (s *npcMovementSystem) Process(elapsedMs int64) {
	if s.Mouse.LeftButtonNewlyDown() {
		mx := s.Camera.Unprojectx(float32(s.Mouse.X()))
		my := s.Camera.UnprojectY(float32(s.Mouse.Y()))
		s.target = &math.Vector{mx, my}
	}

	for _, entity := range s.NpcCollection.Entities() {
		physicsComponent, ok := s.PhysicsComponentManager.GetComponent(entity)
		if !ok {
			continue
		}

		mod := float32(1000)
		speed := float32(.35)
		transitionSpeed := float32(4)

		if s.target != nil {
			physicsComponent.Body.LinearVelocity =
				physicsComponent.Body.LinearVelocity.Muls(1 - (float32(elapsedMs) / mod * transitionSpeed)).Add(
					s.target.Sub(physicsComponent.Body.Position).Normalize().Muls(speed * float32(elapsedMs) / mod * transitionSpeed),
				)

			if s.target.Sub(physicsComponent.Body.Position).Len() < 2 {
				s.target = nil
			}
		} else {
			physicsComponent.Body.LinearVelocity = math.Vector{0, 0}
		}
	}
}
