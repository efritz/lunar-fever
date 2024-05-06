package gameplay

import (
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/ecs/component"
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity"
	"github.com/efritz/lunar-fever/internal/engine/physics"
)

type doorOpenerSystem struct {
	*engine.Context
	doorCollection          *entity.Collection
	playerCollection        *entity.Collection
	physicsComponentManager *component.TypedManager[*physics.PhysicsComponent, physics.PhysicsComponentType]
}

func (s *doorOpenerSystem) Init() {}
func (s *doorOpenerSystem) Exit() {}

func (s *doorOpenerSystem) Process(elapsedMs int64) {
	var playerPhysicsComponent *physics.PhysicsComponent
	for _, entity := range s.playerCollection.Entities() {
		component, ok := s.physicsComponentManager.GetComponent(entity)
		if !ok {
			return
		}

		playerPhysicsComponent = component
	}

	for _, entity := range s.doorCollection.Entities() {
		component, ok := s.physicsComponentManager.GetComponent(entity)
		if !ok {
			return
		}

		component.CollisionsDisabled = component.Body.Position.Sub(playerPhysicsComponent.Body.Position).Len() < 50
	}
}
