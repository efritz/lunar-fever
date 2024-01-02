package entity

import "github.com/efritz/lunar-fever/internal/engine/event"

type (
	EntityChangedEvent        struct{ Entity Entity }
	EntityChangedListener     interface{ OnEntityChanged(e EntityChangedEvent) }
	EntityChangedEventManager = event.TypedManager[EntityChangedEvent, entityChangedEventType, EntityChangedListener]

	entityChangedEventType struct{}
)

var NewEntityChangedEventManager = event.NewTypedManager[EntityChangedEvent, entityChangedEventType, EntityChangedListener]

func (e EntityChangedEvent) EventType() entityChangedEventType { return entityChangedEventType{} }
func (e EntityChangedEvent) Notify(l EntityChangedListener)    { l.OnEntityChanged(e) }
