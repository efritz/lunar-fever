package group

import (
	"github.com/efritz/lunar-fever/internal/common/datastructures"
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity"
	"github.com/efritz/lunar-fever/internal/engine/event"
)

type Manager struct {
	entityChangedEventManager *entity.EntityChangedEventManager
	groupsByEntityID          map[int64]datastructures.Set[string]
}

func NewManager(eventManager *event.Manager) *Manager {
	m := &Manager{
		entityChangedEventManager: entity.NewEntityChangedEventManager(eventManager),
		groupsByEntityID:          map[int64]datastructures.Set[string]{},
	}

	entityRemovedEventManager := entity.NewEntityRemovedEventManager(eventManager)
	entityRemovedEventManager.AddListener(m)
	return m
}

func (m *Manager) HasGroup(e entity.Entity, group string) bool {
	_, ok := m.groupsByEntityID[e.ID][group]
	return ok
}

func (m *Manager) AddGroup(e entity.Entity, group string) {
	groups, ok := m.groupsByEntityID[e.ID]
	if !ok {
		groups = datastructures.Set[string]{}
		m.groupsByEntityID[e.ID] = groups
	}

	if _, ok := groups[group]; !ok {
		groups[group] = struct{}{}
		m.entityChangedEventManager.Dispatch(entity.EntityChangedEvent{Entity: e})
	}
}

func (m *Manager) RemoveGroup(e entity.Entity, group string) {
	groups, ok := m.groupsByEntityID[e.ID]
	if !ok {
		return
	}

	if _, ok := groups[group]; ok {
		delete(groups, group)
		m.entityChangedEventManager.Dispatch(entity.EntityChangedEvent{Entity: e})
	}
}

func (m *Manager) OnEntityRemoved(e entity.EntityRemovedEvent) {
	delete(m.groupsByEntityID, e.Entity.ID)
}
