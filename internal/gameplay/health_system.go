package gameplay

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type healthSystem struct {
	*GameContext
	entityDamagedEventManager *EntityDamagedEventManager
	entityDeathEventManager   *EntityDeathEventManager
}

func NewHealthSystem(ctx *GameContext) *healthSystem {
	return &healthSystem{
		GameContext:               ctx,
		entityDamagedEventManager: NewEntityDamagedEventManager(ctx.EventManager),
		entityDeathEventManager:   NewEntityDeathEventManager(ctx.EventManager),
	}
}

func (s *healthSystem) Init() {
	s.entityDamagedEventManager.AddListener(s)
}

func (s *healthSystem) Exit() {}

func (s *healthSystem) Process(elapsedMs int64) {
	// Temporary implementation
	for _, entity := range s.PlayerCollection.Entities() {
		healthComponent, ok := s.HealthComponentManager.GetComponent(entity)
		if !ok {
			return
		}

		if s.Keyboard.IsKeyNewlyDown(glfw.KeyB) {
			healthComponent.Health -= 15
			s.entityDamagedEventManager.Dispatch(EntityDamagedEvent{entity})
		}
	}
}

func (s *healthSystem) OnEntityDamaged(e EntityDamagedEvent) {
	healthComponent, ok := s.HealthComponentManager.GetComponent(e.Entity)
	if !ok {
		return
	}

	if healthComponent.Health <= 0 {
		s.entityDeathEventManager.Dispatch(EntityDeathEvent{e.Entity})
	}
}
