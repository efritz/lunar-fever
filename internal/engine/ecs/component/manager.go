package component

import (
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity"
	"github.com/efritz/lunar-fever/internal/engine/event"
)

type Manager struct {
	entityChangedEventManager *entity.EntityChangedEventManager
	componentsByEntityID      map[int64]map[ComponentType]Component[any]
}

func NewManager(eventManager *event.Manager) *Manager {
	m := &Manager{
		entityChangedEventManager: entity.NewEntityChangedEventManager(eventManager),
		componentsByEntityID:      map[int64]map[ComponentType]Component[any]{},
	}

	entityRemovedEventManager := entity.NewEntityRemovedEventManager(eventManager)
	entityRemovedEventManager.AddListener(m)
	return m
}

func (m *Manager) HasComponent(e entity.Entity, componentType ComponentType) bool {
	components, ok := m.componentsByEntityID[e.ID]
	if !ok {
		return false
	}

	_, ok = components[componentType]
	return ok
}

func (m *Manager) GetComponent(e entity.Entity, componentType ComponentType) (component Component[any], _ bool) {
	components, ok := m.componentsByEntityID[e.ID]
	if !ok {
		return
	}

	component, ok = components[componentType]
	return component, ok
}

func (m *Manager) AddComponent(e entity.Entity, component Component[any], componentType ComponentType) {
	components, ok := m.componentsByEntityID[e.ID]
	if !ok {
		components = map[ComponentType]Component[any]{}
		m.componentsByEntityID[e.ID] = components
	}

	if _, ok := components[componentType]; !ok {
		components[componentType] = component
		m.entityChangedEventManager.Dispatch(entity.EntityChangedEvent{Entity: e})
	}
}

func (m *Manager) RemoveComponent(e entity.Entity, componentType ComponentType) {
	components, ok := m.componentsByEntityID[e.ID]
	if !ok {
		return
	}

	if _, ok := components[componentType]; ok {
		delete(components, componentType)
		m.entityChangedEventManager.Dispatch(entity.EntityChangedEvent{Entity: e})
	}
}

func (m *Manager) OnEntityRemoved(e entity.EntityRemovedEvent) {
	delete(m.componentsByEntityID, e.Entity.ID)
}
