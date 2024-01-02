package entity

import "github.com/efritz/lunar-fever/internal/engine/event"

type Collection struct {
	matcher  Matcher
	entities map[int64]Entity
}

func NewCollection(matcher Matcher, eventManager *event.Manager) *Collection {
	c := &Collection{
		matcher:  matcher,
		entities: map[int64]Entity{},
	}

	entityChangedEventManager := NewEntityChangedEventManager(eventManager)
	entityChangedEventManager.AddListener(c)
	entityRemovedEventManager := NewEntityRemovedEventManager(eventManager)
	entityRemovedEventManager.AddListener(c)
	return c
}

func (c *Collection) Entities() []Entity {
	entities := make([]Entity, 0, len(c.entities))
	for _, e := range c.entities {
		entities = append(entities, e)
	}

	return entities
}

func (c *Collection) OnEntityChanged(e EntityChangedEvent) {
	if c.matcher.Matches(e.Entity) {
		c.entities[e.Entity.ID] = e.Entity
	} else {
		delete(c.entities, e.Entity.ID)
	}
}

func (c *Collection) OnEntityRemoved(e EntityRemovedEvent) {
	delete(c.entities, e.Entity.ID)
}
