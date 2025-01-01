package gameplay

import (
	"github.com/efritz/lunar-fever/internal/engine/ecs/system"
)

type doorOpenerSystem struct {
	*GameContext
}

func NewDoorOpenerSystem(ctx *GameContext) system.System {
	return &doorOpenerSystem{GameContext: ctx}
}

func (s *doorOpenerSystem) Init() {}
func (s *doorOpenerSystem) Exit() {}

func (s *doorOpenerSystem) Process(elapsedMs int64) {

	for _, entity := range s.DoorCollection.Entities() {
		component, ok := s.PhysicsComponentManager.GetComponent(entity)
		if !ok {
			return
		}

		collisionsDisabled := false
		for _, entity := range s.ScientistCollection.Entities() {
			scientistPhysicsComponent, ok := s.PhysicsComponentManager.GetComponent(entity)
			if !ok {
				return
			}

			if component.Body.Position.Sub(scientistPhysicsComponent.Body.Position).Len() < 50 {
				collisionsDisabled = true
				break
			}
		}

		component.CollisionsDisabled = collisionsDisabled
	}
}
