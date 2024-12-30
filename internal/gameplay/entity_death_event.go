package gameplay

import (
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity"
	"github.com/efritz/lunar-fever/internal/engine/event"
)

type (
	EntityDeathEvent        struct{ Entity entity.Entity }
	EntityDeathListener     interface{ OnEntityDeath(e EntityDeathEvent) }
	EntityDeathEventManager = event.TypedManager[EntityDeathEvent, entityDeathEventType, EntityDeathListener]

	entityDeathEventType struct{}
)

var NewEntityDeathEventManager = event.NewTypedManager[EntityDeathEvent, entityDeathEventType, EntityDeathListener]

func (e EntityDeathEvent) EventType() entityDeathEventType { return entityDeathEventType{} }
func (e EntityDeathEvent) Notify(l EntityDeathListener)    { l.OnEntityDeath(e) }
