package gameplay

import (
	"github.com/efritz/lunar-fever/internal/engine/ecs/system"
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
	*GameContext
}

func NewInteractionSystem(ctx *GameContext) system.System {
	return &interactionSystem{GameContext: ctx}
}

func (s *interactionSystem) Init() {}
func (s *interactionSystem) Exit() {}

var interactionCooldown = 0.5

func (s *interactionSystem) Process(elapsedMs int64) {
	for _, entity := range s.PlayerCollection.Entities() {
		physicsComponent, ok := s.PhysicsComponentManager.GetComponent(entity)
		if !ok {
			return
		}

		interactionComponent, ok := s.InteractionComponentManager.GetComponent(entity)
		if !ok {
			return
		}

		healthComponent, ok := s.HealthComponentManager.GetComponent(entity)
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
