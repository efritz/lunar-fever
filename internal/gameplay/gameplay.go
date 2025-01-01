package gameplay

import (
	"fmt"
	"os"

	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/ecs/system"
	"github.com/efritz/lunar-fever/internal/engine/physics"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
	"github.com/efritz/lunar-fever/internal/engine/view"
	"github.com/efritz/lunar-fever/internal/gameplay/maps"
	"github.com/efritz/lunar-fever/internal/gameplay/maps/loader"
	"github.com/efritz/lunar-fever/internal/gameplay/menu"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Gameplay struct {
	*GameContext

	updateSystemManager *system.Manager
	renderSystemManager *system.Manager

	updateMss     []int64
	updateMsTotal int64
	renderMss     []int64
	renderMsTotal int64
}

// setTransitionOnTime(250);

func NewGameplay(engineCtx *engine.Context) view.View {
	tileMap, err := loader.ReadTileMap()
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}

		tileMap = maps.NewTileMap(100, 100, 64)
	}
	base := maps.ConstructBase(tileMap)

	gameCtx := NewGameContext(engineCtx, tileMap, base)

	updateSystemManager := system.NewManager()
	updateSystemManager.Add(physics.NewPhysicsComponentSystem(gameCtx.EventManager, gameCtx.ComponentManager), 0)
	updateSystemManager.Add(physics.NewCollisionResolution(gameCtx.EventManager, gameCtx.ComponentManager), 0)
	updateSystemManager.Add(NewPlayerMovementSystem(gameCtx), 0)
	updateSystemManager.Add(NewRoverMovementSystem(gameCtx), 0)
	updateSystemManager.Add(NewCameraMovementSystem(gameCtx), 0)
	updateSystemManager.Add(NewDoorOpenerSystem(gameCtx), 0)
	updateSystemManager.Add(NewInteractionSystem(gameCtx), 0)
	updateSystemManager.Add(NewHealthSystem(gameCtx), 0)
	updateSystemManager.Add(NewNpcMovementSystem(gameCtx), 0)
	updateSystemManager.Add(gameCtx.CameraDirector, 0)

	renderSystemManager := system.NewManager()
	renderSystemManager.Add(NewRegolithRenderSystem(gameCtx), 0)
	renderSystemManager.Add(maps.NewBaseRenderSystem(engineCtx, tileMap, base), 1)
	renderSystemManager.Add(NewScientistRenderSystem(gameCtx), 2)
	renderSystemManager.Add(NewRoverRenderSystem(gameCtx), 2)
	renderSystemManager.Add(NewPhysicsRenderSystem(gameCtx), 3)
	renderSystemManager.Add(NewDoorRenderSystem(gameCtx), 4)
	renderSystemManager.Add(NewInteractionRenderSystem(gameCtx), 5)

	createPlayer(gameCtx)
	createScientist(gameCtx)
	createRover(gameCtx)
	createWalls(gameCtx)

	return &Gameplay{
		GameContext:         gameCtx,
		updateSystemManager: updateSystemManager,
		renderSystemManager: renderSystemManager,
	}
}

func (g *Gameplay) Init() {
	initFonts()
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
		g.ViewManager.Add(menu.NewPauseMenu(g.Context))
	}
	if g.Keyboard.IsKeyNewlyDown(glfw.KeyTab) {
		g.ViewManager.Add(menu.NewObjectiveMenu(g.Context))
	}

	// Center on player
	if g.Keyboard.IsKeyNewlyDown(glfw.KeySpace) {
		for _, entity := range g.PlayerCollection.Entities() {
			component, ok := g.GameContext.PhysicsComponentManager.GetComponent(entity)
			if !ok {
				continue
			}

			x1, y1, x2, y2 := component.Body.CoverBound()
			g.GameContext.CameraDirector.LookAt(x1+(x2-x1)/2, y1+(y2-y1)/2, 1000)
		}
	}

	// Explosion camera
	// if g.Keyboard.IsKeyNewlyDown(glfw.KeyP) {
	// 	g.director.AddShake(5)
	// }

	// Toggle debug flag
	if g.Keyboard.IsKeyNewlyDown(glfw.KeyP) {
		debug = !debug
	}

	g.updateSystemManager.Process(elapsedMs)
}

var debug = false

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

	if debug {
		font.Printf(
			rendering.DisplayWidth-200,
			30,
			fmt.Sprintf("%d ups, %d fps", len(g.updateMss), len(g.renderMss)),
			rendering.WithTextScale(0.5),
			rendering.WithTextColor(rendering.Color{0, 0, 0, 1}),
		)
	}
}

func (g *Gameplay) IsOverlay() bool {
	return false
}
