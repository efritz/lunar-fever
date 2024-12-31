package gameplay

import (
	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine/ecs/system"
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
	for _, entity := range s.NpcCollection.Entities() {
		physicsComponent, ok := s.PhysicsComponentManager.GetComponent(entity)
		if !ok {
			continue
		}

		xDir := float32(1)
		yDir := float32(1)
		mod := float32(1000)
		speed := float32(.35)
		transitionSpeed := float32(4)

		physicsComponent.Body.LinearVelocity = physicsComponent.Body.LinearVelocity.
			Muls(1 - (float32(elapsedMs) / mod * transitionSpeed)).
			Add(math.Vector{xDir, yDir}.Muls(speed * float32(elapsedMs) / mod * transitionSpeed))
	}
}
