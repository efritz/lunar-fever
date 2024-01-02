package component

import "github.com/efritz/lunar-fever/internal/engine/ecs/entity"

type EntityMatcher struct {
	componentManager *Manager
	componentTypes   []ComponentType
}

func NewEntityMatcher(componentManager *Manager, componentTypes ...ComponentType) entity.Matcher {
	return &EntityMatcher{
		componentManager: componentManager,
		componentTypes:   componentTypes,
	}
}

func (m *EntityMatcher) Matches(e entity.Entity) bool {
	for _, componentType := range m.componentTypes {
		if !m.componentManager.HasComponent(e, componentType) {
			return false
		}
	}

	return true
}
