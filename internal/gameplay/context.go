package gameplay

import (
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/ecs/component"
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity"
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity/group"
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity/tag"
	"github.com/efritz/lunar-fever/internal/engine/event"
	"github.com/efritz/lunar-fever/internal/engine/physics"
	"github.com/efritz/lunar-fever/internal/gameplay/maps"
)

type GameContext struct {
	*engine.Context

	TileMap        *maps.TileMap
	CameraDirector *CameraDirector

	EventManager     *event.Manager
	EntityManager    *entity.Manager
	ComponentManager *component.Manager
	TagManager       *tag.Manager
	GroupManager     *group.Manager

	PhysicsComponentManager     *component.TypedManager[*physics.PhysicsComponent, physics.PhysicsComponentType]
	HealthComponentManager      *component.TypedManager[*HealthComponent, HealthComponentType]
	InteractionComponentManager *component.TypedManager[*InteractionComponent, InteractionComponentType]

	PlayerCollection    *entity.Collection
	ScientistCollection *entity.Collection
	RoverCollection     *entity.Collection
	NpcCollection       *entity.Collection
	DoorCollection      *entity.Collection
	PhysicsCollection   *entity.Collection
}

func NewGameContext(engineCtx *engine.Context, tileMap *maps.TileMap) *GameContext {
	eventManager := event.NewManager()
	entityManager := entity.NewManager(eventManager)
	componentManager := component.NewManager(eventManager)
	tagManager := tag.NewManager(eventManager)
	groupManager := group.NewManager(eventManager)

	return &GameContext{
		Context:        engineCtx,
		TileMap:        tileMap,
		CameraDirector: &CameraDirector{Context: engineCtx},

		EventManager:     eventManager,
		EntityManager:    entityManager,
		ComponentManager: componentManager,
		TagManager:       tagManager,
		GroupManager:     groupManager,

		PhysicsComponentManager:     component.NewTypedManager[*physics.PhysicsComponent](componentManager, eventManager),
		HealthComponentManager:      component.NewTypedManager[*HealthComponent](componentManager, eventManager),
		InteractionComponentManager: component.NewTypedManager[*InteractionComponent](componentManager, eventManager),

		PlayerCollection:    entity.NewCollection(tag.NewEntityMatcher(tagManager, "player"), eventManager),
		ScientistCollection: entity.NewCollection(group.NewEntityMatcher(groupManager, "scientist"), eventManager),
		RoverCollection:     entity.NewCollection(tag.NewEntityMatcher(tagManager, "rover"), eventManager),
		NpcCollection:       entity.NewCollection(group.NewEntityMatcher(groupManager, "npc"), eventManager),
		DoorCollection:      entity.NewCollection(group.NewEntityMatcher(groupManager, "door"), eventManager),
		PhysicsCollection:   entity.NewCollection(group.NewEntityMatcher(groupManager, "physics"), eventManager),
	}
}
