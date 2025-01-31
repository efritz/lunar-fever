package physics

import (
	"github.com/efritz/lunar-fever/internal/engine/ecs/component"
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity"
	"github.com/efritz/lunar-fever/internal/engine/ecs/system"
	"github.com/efritz/lunar-fever/internal/engine/event"
)

type PhysicsComponentSystemDelegate struct {
	entityMovedEventManager *EntityMovedEventManager
	physicsComponentManager *component.TypedManager[*PhysicsComponent, PhysicsComponentType]
}

func NewPhysicsComponentSystem(eventManager *event.Manager, componentManager *component.Manager) system.System {
	physicsComponentManager := component.NewTypedManager[*PhysicsComponent](componentManager, eventManager)
	physicsComponentMatcher := component.NewEntityMatcher(componentManager, physicsComponentType)
	collection := entity.NewCollection(physicsComponentMatcher, eventManager)

	return entity.NewSystem(&PhysicsComponentSystemDelegate{
		entityMovedEventManager: NewEntityMovedEventManager(eventManager),
		physicsComponentManager: physicsComponentManager,
	}, collection)
}

func (d *PhysicsComponentSystemDelegate) Init() {}
func (d *PhysicsComponentSystemDelegate) Exit() {}

func (d *PhysicsComponentSystemDelegate) Process(entity entity.Entity, elapsedMs int64) {
	physicsComponent, ok := d.physicsComponentManager.GetComponent(entity)
	if !ok {
		return
	}

	body := physicsComponent.Body

	if body.inverseMass == 0 {
		return
	}

	body.LinearVelocity = body.LinearVelocity.Add(body.force.Muls(body.inverseMass * float32(elapsedMs)))
	body.AngularVelocity = body.AngularVelocity + (body.torque * body.inverseInertia * float32(elapsedMs))

	body.Position = body.Position.Add(body.LinearVelocity.Muls(float32(elapsedMs)))
	body.SetOrient(body.Orient + body.AngularVelocity*float32(elapsedMs))

	decayRate := float32(0.99)
	body.LinearVelocity = body.LinearVelocity.Muls(decayRate)
	body.AngularVelocity = body.AngularVelocity * decayRate

	// TODO - only if position or orient actually changed
	d.entityMovedEventManager.Dispatch(EntityMovedEvent{entity})

	// TODO - can do before dispatch?
	body.ClearForces()
}
