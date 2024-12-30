package gameplay

import (
	"github.com/efritz/lunar-fever/internal/engine/physics"
)

type doorOpenerSystem struct {
	*GameContext
}

func NewDoorOpenerSystem(ctx *GameContext) *doorOpenerSystem {
	return &doorOpenerSystem{GameContext: ctx}
}

func (s *doorOpenerSystem) Init() {}
func (s *doorOpenerSystem) Exit() {}

func (s *doorOpenerSystem) Process(elapsedMs int64) {
	var playerPhysicsComponent *physics.PhysicsComponent
	for _, entity := range s.PlayerCollection.Entities() {
		component, ok := s.PhysicsComponentManager.GetComponent(entity)
		if !ok {
			return
		}

		playerPhysicsComponent = component
	}

	for _, entity := range s.DoorCollection.Entities() {
		component, ok := s.PhysicsComponentManager.GetComponent(entity)
		if !ok {
			return
		}

		component.CollisionsDisabled = component.Body.Position.Sub(playerPhysicsComponent.Body.Position).Len() < 50
	}
}
