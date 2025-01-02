package editor

import (
	"fmt"

	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
	"github.com/efritz/lunar-fever/internal/engine/view"
	"github.com/efritz/lunar-fever/internal/gameplay/maps"
	"github.com/efritz/lunar-fever/internal/gameplay/maps/loader"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Editor struct {
	*engine.Context
	texture      rendering.Texture
	tileMap      *maps.TileMap // No need to store
	baseRenderer *maps.BaseRenderer
	executor     *MapCommandExecutor

	x, y             int
	selected         Palette
	performingAction bool
	isValidSelection bool
}

func NewEditor(engineCtx *engine.Context) view.View {
	return &Editor{
		Context: engineCtx,
	}
}

func (e *Editor) Init() {
	tm, err := loader.ReadTileMap()
	if err != nil {
		tm = maps.NewTileMap(100, 100, 64)
	}
	e.tileMap = tm

	e.texture = e.TextureLoader.Load("base").Region(7*32, 1*32, 32, 32)
	e.baseRenderer = maps.NewBaseRenderer(e.SpriteBatch, e.TextureLoader, e.tileMap, true)
	e.executor = NewMapCommandExecutor(e.tileMap)
	e.selected = FLOOR_TOOL
}

func (e *Editor) Exit() {}

func (e *Editor) Update(elapsedMs int64, hasFocus bool) {
	//
	// Palette selection

	if e.Keyboard.IsKeyNewlyDown(glfw.Key1) {
		e.selected = FLOOR_TOOL
		fmt.Printf("Selected floor tool\n")
	}
	if e.Keyboard.IsKeyNewlyDown(glfw.Key2) {
		e.selected = HWALL_TOOL
		fmt.Printf("Selected horizontal wall tool\n")
	}
	if e.Keyboard.IsKeyNewlyDown(glfw.Key3) {
		e.selected = VWALL_TOOL
		fmt.Printf("Selected vertical wall tool\n")
	}
	if e.Keyboard.IsKeyNewlyDown(glfw.Key4) {
		e.selected = HDOOR_TOOL
		fmt.Printf("Selected horizontal door tool\n")
	}
	if e.Keyboard.IsKeyNewlyDown(glfw.Key5) {
		e.selected = VDOOR_TOOL
		fmt.Printf("Selected vertical door tool\n")
	}

	//
	// Update position

	size := 64
	x := int(e.Mouse.X()) / size
	y := int(e.Mouse.Y()) / size

	oldX := e.x
	oldY := e.y
	e.x = x
	e.y = y

	row := y
	col := x
	e.isValidSelection = e.executor.HasAction(e.selected, row, col)

	//
	// Fire actions

	if e.Mouse.LeftButton() {
		if e.Mouse.LeftButtonNewlyDown() {
			e.executor.PrepareAction(e.selected, row, col)
		}

		if e.Mouse.LeftButtonNewlyDown() || e.x != oldX || e.y != oldY {
			e.executor.ExecuteAction(e.selected, row, col)
		}

		e.performingAction = true
	} else {
		e.performingAction = false
	}

	//
	// Fire undo/redo actions

	if e.Keyboard.IsKeyDown(glfw.KeyLeftSuper) && e.Keyboard.IsKeyNewlyDown(glfw.KeyZ) {
		e.executor.Undo()
	}
	if e.Keyboard.IsKeyDown(glfw.KeyLeftSuper) && e.Keyboard.IsKeyNewlyDown(glfw.KeyY) {
		e.executor.Redo()
	}

	//
	// Save tile map

	if e.Keyboard.IsKeyDown(glfw.KeyLeftSuper) && e.Keyboard.IsKeyNewlyDown(glfw.KeyS) {
		if err := loader.Write(e.tileMap); err != nil {
			panic(err.Error())
		}
	}
}

func (e *Editor) Render(elapsedMs int64) {
	tileSize := float32(64)

	// TODO - grid lines?
	// g.setColor(Color.lightGray);
	// g.setStroke(new BasicStroke(1f, BasicStroke.CAP_BUTT, BasicStroke.JOIN_MITER, 100f, new float[]{10f}, 0f));
	// int tileWidth = 4 * zoom;
	// for (int i = 0; i < getHeight(); i += tileWidth) {
	// 	g.drawLine(0, i, getWidth(), i);
	// }
	// for (int j = 0; j < getWidth(); j += tileWidth) {
	// 	g.drawLine(j, 0, j, getHeight());
	// }

	e.baseRenderer.Render(0, 0, rendering.DisplayWidth, rendering.DisplayHeight, nil, nil, false)

	var color rendering.Color
	if e.performingAction {
		color = rendering.Color{0, 0, 0, 0.5}
	} else if e.isValidSelection {
		color = rendering.Color{0, 1, 0, 0.5}
	} else {
		color = rendering.Color{1, 0, 0, 0.5}
	}

	e.SpriteBatch.Begin()
	e.SpriteBatch.Draw(e.texture, float32(e.x)*tileSize, float32(e.y)*tileSize, tileSize, tileSize, rendering.WithColor(color))
	e.SpriteBatch.End()
}

func (e *Editor) IsOverlay() bool {
	return false
}
