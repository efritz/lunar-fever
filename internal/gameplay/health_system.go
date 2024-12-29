package gameplay

import (
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/ecs/component"
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity"
	"github.com/efritz/lunar-fever/internal/engine/event"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type HealthComponent struct {
	Health    float32
	MaxHealth float32
}

type HealthComponentType struct{}

var healthComponentType = HealthComponentType{}

func (c *HealthComponent) ComponentType() HealthComponentType {
	return healthComponentType
}

//
//

type (
	EntityDamagedEvent        struct{ entity entity.Entity }
	EntityDamagedListener     interface{ OnEntityDamaged(e EntityDamagedEvent) }
	EntityDamagedEventManager = event.TypedManager[EntityDamagedEvent, entityDamagedEventType, EntityDamagedListener]

	entityDamagedEventType struct{}
)

var NewEntityDamagedEventManager = event.NewTypedManager[EntityDamagedEvent, entityDamagedEventType, EntityDamagedListener]

func (e EntityDamagedEvent) EventType() entityDamagedEventType { return entityDamagedEventType{} }
func (e EntityDamagedEvent) Notify(l EntityDamagedListener)    { l.OnEntityDamaged(e) }

//
//

type (
	EntityDeathEvent        struct{ entity entity.Entity }
	EntityDeathListener     interface{ OnEntityDeath(e EntityDeathEvent) }
	EntityDeathEventManager = event.TypedManager[EntityDeathEvent, entityDeathEventType, EntityDeathListener]

	entityDeathEventType struct{}
)

var NewEntityDeathEventManager = event.NewTypedManager[EntityDeathEvent, entityDeathEventType, EntityDeathListener]

func (e EntityDeathEvent) EventType() entityDeathEventType { return entityDeathEventType{} }
func (e EntityDeathEvent) Notify(l EntityDeathListener)    { l.OnEntityDeath(e) }

//

//
//

type healthSystem struct {
	*engine.Context
	entityDamagedEventManager *event.TypedManager[EntityDamagedEvent, entityDamagedEventType, EntityDamagedListener]
	entityDeathEventManager   *event.TypedManager[EntityDeathEvent, entityDeathEventType, EntityDeathListener]
	healthCollection          *entity.Collection
	healthComponentManager    *component.TypedManager[*HealthComponent, HealthComponentType]
}

func NewHealthSystem(
	ctx *engine.Context,
	eventManager *event.Manager,
	healthCollection *entity.Collection,
	healthComponentManager *component.TypedManager[*HealthComponent, HealthComponentType],
) *healthSystem {
	return &healthSystem{
		Context:                   ctx,
		entityDamagedEventManager: NewEntityDamagedEventManager(eventManager),
		entityDeathEventManager:   NewEntityDeathEventManager(eventManager),
		healthCollection:          healthCollection,
		healthComponentManager:    healthComponentManager,
	}
}

func (s *healthSystem) Init() {
	s.entityDamagedEventManager.AddListener(s)
}

func (s *healthSystem) Exit() {}

func (s *healthSystem) Process(elapsedMs int64) {
	// Temporary implementation
	for _, entity := range s.healthCollection.Entities() {
		healthComponent, ok := s.healthComponentManager.GetComponent(entity)
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
	healthComponent, ok := s.healthComponentManager.GetComponent(e.entity)
	if !ok {
		return
	}

	if healthComponent.Health <= 0 {
		s.entityDeathEventManager.Dispatch(EntityDeathEvent{e.entity})
	}
}
