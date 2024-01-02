package tag

import "github.com/efritz/lunar-fever/internal/engine/ecs/entity"

type EntityMatcher struct {
	manager *Manager
	tag     string
}

func NewEntityMatcher(manager *Manager, tag string) entity.Matcher {
	return &EntityMatcher{
		manager: manager,
		tag:     tag,
	}
}

func (m *EntityMatcher) Matches(e entity.Entity) bool {
	return m.manager.HasTag(e, m.tag)
}
