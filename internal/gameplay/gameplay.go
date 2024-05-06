package gameplay

import (
	"fmt"
	"os"

	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/ecs/component"
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity"
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity/group"
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity/tag"
	"github.com/efritz/lunar-fever/internal/engine/ecs/system"
	"github.com/efritz/lunar-fever/internal/engine/event"
	"github.com/efritz/lunar-fever/internal/engine/physics"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
	"github.com/efritz/lunar-fever/internal/engine/view"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Gameplay struct {
	*engine.Context
	updateSystemManager     *system.Manager
	renderSystemManager     *system.Manager
	playerCollection        *entity.Collection
	physicsComponentManager *component.TypedManager[*physics.PhysicsComponent, physics.PhysicsComponentType]
	director                *CameraDirector

	updateMss     []int64
	updateMsTotal int64
	renderMss     []int64
	renderMsTotal int64
}

// setTransitionOnTime(250);

func NewGameplay(engineCtx *engine.Context) view.View {
	eventManager := event.NewManager()
	entityManager := entity.NewManager(eventManager)
	componentManager := component.NewManager(eventManager)
	tagManager := tag.NewManager(eventManager)
	groupManager := group.NewManager(eventManager)
	playerCollection := entity.NewCollection(tag.NewEntityMatcher(tagManager, "player"), eventManager)
	roverCollection := entity.NewCollection(tag.NewEntityMatcher(tagManager, "rover"), eventManager)
	npcCollection := entity.NewCollection(group.NewEntityMatcher(groupManager, "npc"), eventManager)
	physicsCollection := entity.NewCollection(group.NewEntityMatcher(groupManager, "physics"), eventManager)
	physicsComponentManager := component.NewTypedManager[*physics.PhysicsComponent, physics.PhysicsComponentType](componentManager, eventManager)
	director := &CameraDirector{Context: engineCtx}

	tileMap, err := readTileMap()
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}

		tileMap = NewTileMap(100, 100, 64)
	}

	updateSystemManager := system.NewManager()
	updateSystemManager.Add(physics.NewPhysicsComponentSystem(eventManager, componentManager), 0)
	updateSystemManager.Add(physics.NewCollisionResolution(eventManager, componentManager), 0)
	updateSystemManager.Add(&playerMovementSystem{Context: engineCtx, playerCollection: playerCollection, physicsComponentManager: physicsComponentManager}, 0)
	updateSystemManager.Add(&roverMovementSystem{Context: engineCtx, roverCollection: roverCollection, physicsComponentManager: physicsComponentManager}, 0)
	updateSystemManager.Add(&cameraMovementSystem{Context: engineCtx}, 0)
	updateSystemManager.Add(director, 0)

	renderSystemManager := system.NewManager()
	renderSystemManager.Add(&regolithRenderSystem{Context: engineCtx}, 0)
	renderSystemManager.Add(&baseRenderSystem{Context: engineCtx, tileMap: tileMap}, 1)
	renderSystemManager.Add(&playerRenderSystem{Context: engineCtx, playerCollection: playerCollection, physicsComponentManager: physicsComponentManager}, 2)
	renderSystemManager.Add(&roverRenderSystem{Context: engineCtx, roverCollection: roverCollection, physicsComponentManager: physicsComponentManager}, 2)
	renderSystemManager.Add(&npcRenderSystem{Context: engineCtx, npcCollection: npcCollection, physicsComponentManager: physicsComponentManager}, 2)
	renderSystemManager.Add(&physicsRenderSystem{Context: engineCtx, entityCollection: physicsCollection, physicsComponentManager: physicsComponentManager}, 3)

	createPlayer(entityManager, tagManager, groupManager, physicsComponentManager)
	createRover(entityManager, tagManager, groupManager, physicsComponentManager)
	createNPC(entityManager, groupManager, physicsComponentManager)
	createWalls(entityManager, groupManager, physicsComponentManager, tileMap)

	return &Gameplay{
		Context:                 engineCtx,
		updateSystemManager:     updateSystemManager,
		renderSystemManager:     renderSystemManager,
		playerCollection:        playerCollection,
		physicsComponentManager: physicsComponentManager,
		director:                director,
	}
}

func (g *Gameplay) Init() {
	g.updateSystemManager.Init()
	g.renderSystemManager.Init()
}

func (g *Gameplay) Exit() {
	g.updateSystemManager.Exit()
	g.renderSystemManager.Exit()
}

func (g *Gameplay) Update(elapsedMs int64, hasFocus bool) {
	g.updateMss = append(g.updateMss, elapsedMs)
	g.updateMsTotal += elapsedMs
	for g.updateMsTotal > 1000 {
		g.updateMsTotal -= g.updateMss[0]
		g.updateMss = g.updateMss[1:]
	}

	// Menu management
	if g.Keyboard.IsKeyNewlyDown(glfw.KeyEscape) {
		g.ViewManager.Add(NewPauseMenu(g.Context))
	}
	if g.Keyboard.IsKeyNewlyDown(glfw.KeyTab) {
		g.ViewManager.Add(NewObjectiveMenu(g.Context))
	}

	// Center on player
	if g.Keyboard.IsKeyNewlyDown(glfw.KeySpace) {
		for _, entity := range g.playerCollection.Entities() {
			component, ok := g.physicsComponentManager.GetComponent(entity)
			if !ok {
				continue
			}

			x1, y1, x2, y2 := component.Body.CoverBound()
			g.director.LookAt(x1+(x2-x1)/2, y1+(y2-y1)/2, 1000)
		}
	}

	// Explosion camera
	if g.Keyboard.IsKeyNewlyDown(glfw.KeyP) {
		g.director.AddShake(5)
	}

	g.updateSystemManager.Process(elapsedMs)
}

func (g *Gameplay) Render(elapsedMs int64) {
	g.renderMss = append(g.renderMss, elapsedMs)
	g.renderMsTotal += elapsedMs
	for g.renderMsTotal > 1000 {
		g.renderMsTotal -= g.renderMss[0]
		g.renderMss = g.renderMss[1:]
	}

	g.SpriteBatch.SetViewMatrix(g.Camera.ViewMatrix())
	g.renderSystemManager.Process(elapsedMs)
	g.SpriteBatch.SetViewMatrix(math.IdentityMatrix)

	font.Printf(
		rendering.DisplayWidth-200,
		30,
		fmt.Sprintf("%d ups, %d fps", len(g.updateMss), len(g.renderMss)),
		rendering.WithTextScale(0.5),
		rendering.WithTextColor(rendering.Color{0, 0, 0, 1}),
	)
}

func (g *Gameplay) IsOverlay() bool {
	return false
}
