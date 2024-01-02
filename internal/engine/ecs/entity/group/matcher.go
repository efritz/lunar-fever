package group

import "github.com/efritz/lunar-fever/internal/engine/ecs/entity"

type EntityMatcher struct {
	manager *Manager
	group   string
}

func NewEntityMatcher(manager *Manager, group string) entity.Matcher {
	return &EntityMatcher{
		manager: manager,
		group:   group,
	}
}

func (m *EntityMatcher) Matches(e entity.Entity) bool {
	return m.manager.HasGroup(e, m.group)
}
