package entity

import "github.com/efritz/lunar-fever/internal/engine/event"

type (
	EntityRemovedEvent        struct{ Entity Entity }
	EntityRemovedListener     interface{ OnEntityRemoved(e EntityRemovedEvent) }
	EntityRemovedEventManager = event.TypedManager[EntityRemovedEvent, entityRemovedEventType, EntityRemovedListener]

	entityRemovedEventType struct{}
)

var NewEntityRemovedEventManager = event.NewTypedManager[EntityRemovedEvent, entityRemovedEventType, EntityRemovedListener]

func (e EntityRemovedEvent) EventType() entityRemovedEventType { return entityRemovedEventType{} }
func (e EntityRemovedEvent) Notify(l EntityRemovedListener)    { l.OnEntityRemoved(e) }
