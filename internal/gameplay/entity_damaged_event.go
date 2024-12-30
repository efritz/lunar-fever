package gameplay

import (
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity"
	"github.com/efritz/lunar-fever/internal/engine/event"
)

type (
	EntityDamagedEvent        struct{ Entity entity.Entity }
	EntityDamagedListener     interface{ OnEntityDamaged(e EntityDamagedEvent) }
	EntityDamagedEventManager = event.TypedManager[EntityDamagedEvent, entityDamagedEventType, EntityDamagedListener]

	entityDamagedEventType struct{}
)

var NewEntityDamagedEventManager = event.NewTypedManager[EntityDamagedEvent, entityDamagedEventType, EntityDamagedListener]

func (e EntityDamagedEvent) EventType() entityDamagedEventType { return entityDamagedEventType{} }
func (e EntityDamagedEvent) Notify(l EntityDamagedListener)    { l.OnEntityDamaged(e) }
