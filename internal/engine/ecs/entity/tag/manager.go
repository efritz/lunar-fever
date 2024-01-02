package tag

import (
	"github.com/efritz/lunar-fever/internal/common/datastructures"
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity"
	"github.com/efritz/lunar-fever/internal/engine/event"
)

type Manager struct {
	entityChangedEventManager *entity.EntityChangedEventManager
	tagByEntityID             map[int64]string
	tags                      datastructures.Set[string]
}

func NewManager(eventManager *event.Manager) *Manager {
	m := &Manager{
		entityChangedEventManager: entity.NewEntityChangedEventManager(eventManager),
		tagByEntityID:             map[int64]string{},
		tags:                      datastructures.Set[string]{},
	}

	entityRemovedEventManager := entity.NewEntityRemovedEventManager(eventManager)
	entityRemovedEventManager.AddListener(m)
	return m
}

func (m *Manager) HasTag(e entity.Entity, tag string) bool {
	return m.tagByEntityID[e.ID] == tag
}

func (m *Manager) SetTag(e entity.Entity, tag string) {
	if _, ok := m.tags[tag]; ok {
		panic("tag already set")
	}

	m.tags[tag] = struct{}{}
	m.tagByEntityID[e.ID] = tag
	m.entityChangedEventManager.Dispatch(entity.EntityChangedEvent{Entity: e})
}

func (m *Manager) RemoveTag(e entity.Entity) {
	if m.removeTag(e) {
		m.entityChangedEventManager.Dispatch(entity.EntityChangedEvent{Entity: e})
	}
}

func (m *Manager) OnEntityRemoved(e entity.EntityRemovedEvent) {
	_ = m.removeTag(e.Entity)
}

func (m *Manager) removeTag(e entity.Entity) bool {
	tag, ok := m.tagByEntityID[e.ID]
	if !ok {
		return false
	}

	delete(m.tagByEntityID, e.ID)
	delete(m.tags, tag)
	return true
}
