package entity

import (
	"sync/atomic"

	"github.com/efritz/lunar-fever/internal/engine/event"
)

type Manager struct {
	id                        int64
	entityCreatedEventManager *EntityCreatedEventManager
	entityRemovedEventManager *EntityRemovedEventManager
}

func NewManager(eventManager *event.Manager) *Manager {
	return &Manager{
		entityCreatedEventManager: NewEntityCreatedEventManager(eventManager),
		entityRemovedEventManager: NewEntityRemovedEventManager(eventManager),
	}
}

func (m *Manager) Create() Entity {
	id := atomic.AddInt64(&m.id, 1)
	entity := Entity{ID: id}
	m.entityCreatedEventManager.Dispatch(EntityCreatedEvent{Entity: entity})
	return entity
}

func (m *Manager) Remove(entity Entity) {
	m.entityRemovedEventManager.Dispatch(EntityRemovedEvent{Entity: entity})
}
