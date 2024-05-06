package physics

import (
	"github.com/efritz/lunar-fever/internal/engine/ecs/component"
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity"
	"github.com/efritz/lunar-fever/internal/engine/ecs/system"
	"github.com/efritz/lunar-fever/internal/engine/event"
)

type CollisionResolutionSystem struct {
	eventManager            *event.Manager
	physicsComponentManager *component.TypedManager[*PhysicsComponent, PhysicsComponentType]
	entityCollection        *entity.Collection
}

func NewCollisionResolution(eventManager *event.Manager, componentManager *component.Manager) system.System {
	physicsComponentManager := component.NewTypedManager[*PhysicsComponent, PhysicsComponentType](componentManager, eventManager)
	physicsComponentMatcher := component.NewEntityMatcher(componentManager, physicsComponentType)
	collection := entity.NewCollection(physicsComponentMatcher, eventManager)

	return &CollisionResolutionSystem{
		eventManager:            eventManager,
		physicsComponentManager: physicsComponentManager,
		entityCollection:        collection,
	}
}

func (d *CollisionResolutionSystem) Init() {}
func (d *CollisionResolutionSystem) Exit() {}

const iterations = 10 // TODO - rename

func (d *CollisionResolutionSystem) Process(elapsedMs int64) {
	var contacts []*Contact

	entities := d.entityCollection.Entities()
	for i := 0; i < len(entities); i++ {
		for j := i + 1; j < len(entities); j++ {
			component1, _ := d.physicsComponentManager.GetComponent(entities[i])
			component2, _ := d.physicsComponentManager.GetComponent(entities[j])

			if component1.CollisionsDisabled || component2.CollisionsDisabled {
				continue
			}

			body1 := component1.Body
			body2 := component2.Body

			if body1.inverseMass == 0 && body2.inverseMass == 0 {
				continue
			}

			x1a, y1a, x2a, y2a := body1.CoverBound()
			x1b, y1b, x2b, y2b := body2.CoverBound()

			if !intersects(x1a, y1a, x2a, y2a, x1b, y1b, x2b, y2b) {
				continue
			}

			for _, fixture1 := range body1.Fixtures {
				for _, fixture2 := range body2.Fixtures {
					if contact := NewContact(fixture1, body1, fixture2, body2); contact != nil {
						contacts = append(contacts, contact)
					}
				}
			}
		}
	}

	for j := 0; j < iterations; j++ {
		for _, contact := range contacts {
			contact.ApplyImpulse()
		}
	}

	for _, contact := range contacts {
		contact.Correct()
	}
}

func intersects(x1a, y1a, x2a, y2a, x1b, y1b, x2b, y2b float32) bool {
	if x1a > x2b || x2a < x1b {
		return false
	}

	if y1a > y2b || y2a < y1b {
		return false
	}

	return true
}
