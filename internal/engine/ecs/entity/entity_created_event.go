package entity

import "github.com/efritz/lunar-fever/internal/engine/event"

type (
	EntityCreatedEvent        struct{ Entity Entity }
	EntityCreatedListener     interface{ OnEntityCreated(e EntityCreatedEvent) }
	EntityCreatedEventManager = event.TypedManager[EntityCreatedEvent, entityCreatedEventType, EntityCreatedListener]

	entityCreatedEventType struct{}
)

var NewEntityCreatedEventManager = event.NewTypedManager[EntityCreatedEvent, entityCreatedEventType, EntityCreatedListener]

func (e EntityCreatedEvent) EventType() entityCreatedEventType { return entityCreatedEventType{} }
func (e EntityCreatedEvent) Notify(l EntityCreatedListener)    { l.OnEntityCreated(e) }
