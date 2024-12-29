package gameplay

import (
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/ecs/component"
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity"
	"github.com/efritz/lunar-fever/internal/engine/physics"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type InteractionComponent struct {
	Interacting   bool
	CooldownTimer int64
}

type InteractionComponentType struct{}

var interactionComponentType = InteractionComponentType{}

func (c *InteractionComponent) ComponentType() InteractionComponentType {
	return interactionComponentType
}

//
//

type interactionSystem struct {
	*engine.Context
	playerCollection            *entity.Collection
	physicsComponentManager     *component.TypedManager[*physics.PhysicsComponent, physics.PhysicsComponentType]
	interactionComponentManager *component.TypedManager[*InteractionComponent, InteractionComponentType]
	healthComponentManager      *component.TypedManager[*HealthComponent, HealthComponentType]
}

func (s *interactionSystem) Init() {}
func (s *interactionSystem) Exit() {}

var interactionCooldown = 0.5

func (s *interactionSystem) Process(elapsedMs int64) {
	for _, entity := range s.playerCollection.Entities() {
		physicsComponent, ok := s.physicsComponentManager.GetComponent(entity)
		if !ok {
			return
		}

		interactionComponent, ok := s.interactionComponentManager.GetComponent(entity)
		if !ok {
			return
		}

		healthComponent, ok := s.healthComponentManager.GetComponent(entity)
		if !ok {
			return
		}

		interactionComponent.CooldownTimer -= elapsedMs

		if s.Keyboard.IsKeyNewlyDown(glfw.KeyE) && canInteract(physicsComponent, interactionComponent, healthComponent) {
			// TODO - should interact _against_ another object
			interactionComponent.Interacting = true
			interactionComponent.CooldownTimer = int64(interactionCooldown * 1000)
		} else {
			interactionComponent.Interacting = false
		}
	}
}

func canInteract(physicsComponent *physics.PhysicsComponent, interactionComponent *InteractionComponent, healthComponent *HealthComponent) bool {
	if healthComponent.Health <= 0 {
		return false
	}

	if physicsComponent.Body.LinearVelocity.Len() >= minStartAnimationSpeed {
		return false
	}

	if interactionComponent.CooldownTimer > 0 {
		return false
	}

	return true
}
