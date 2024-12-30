package physics

import (
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity"
	"github.com/efritz/lunar-fever/internal/engine/event"
)

type (
	EntityMovedEventType    struct{}
	EntityMovedEvent        struct{ Entity entity.Entity }
	EntityMovedListener     interface{ OnEntityMoved(e EntityMovedEvent) }
	EntityMovedEventManager = event.TypedManager[EntityMovedEvent, EntityMovedEventType, EntityMovedListener]
)

var (
	entityCreatedEventType     = EntityMovedEventType{}
	NewEntityMovedEventManager = event.NewTypedManager[EntityMovedEvent, EntityMovedEventType, EntityMovedListener]
)

func (e EntityMovedEvent) EventType() EntityMovedEventType { return entityCreatedEventType }
func (e EntityMovedEvent) Notify(l EntityMovedListener)    { l.OnEntityMoved(e) }
