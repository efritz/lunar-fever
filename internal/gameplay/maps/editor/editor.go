package editor

import (
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
	"github.com/efritz/lunar-fever/internal/engine/view"
	"github.com/efritz/lunar-fever/internal/gameplay/maps"
	"github.com/efritz/lunar-fever/internal/gameplay/maps/editor/commands"
	"github.com/efritz/lunar-fever/internal/gameplay/maps/loader"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Editor struct {
	*engine.Context
	texture      rendering.Texture
	tileMap      *maps.TileMap // No need to store
	baseRenderer *maps.BaseRenderer
	executor     *MapCommandExecutor

	x, y                int
	selected            Palette
	performingAction    bool
	affectedTileIndexes []commands.TileIndex
	isRemoveAction      bool
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
	initFonts()
}

func (e *Editor) Exit() {}

func (e *Editor) Update(elapsedMs int64, hasFocus bool) {
	//
	// Palette selection

	if e.Keyboard.IsKeyNewlyDown(glfw.Key1) {
		e.selected = FLOOR_TOOL
	}
	if e.Keyboard.IsKeyNewlyDown(glfw.Key2) {
		if e.selected == VWALL_TOOL {
			e.selected = HWALL_TOOL
		} else {
			e.selected = VWALL_TOOL
		}
	}
	if e.Keyboard.IsKeyNewlyDown(glfw.Key3) {
		if e.selected == VDOOR_TOOL {
			e.selected = HDOOR_TOOL
		} else {
			e.selected = VDOOR_TOOL
		}
	}
	if e.Keyboard.IsKeyNewlyDown(glfw.Key4) {
		e.selected = FIXTURE_TOOL
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
	e.affectedTileIndexes, e.isRemoveAction = e.executor.HasAction(e.selected, row, col)

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
	if len(e.affectedTileIndexes) > 0 {
		if e.performingAction {
			color = rendering.Color{0, 0, 0, 0.5}
		} else if e.isRemoveAction {
			color = rendering.Color{1, 0, 0, 0.5}
		} else {
			color = rendering.Color{0, 1, 0, 0.5}
		}
	} else {
		color = rendering.Color{1, 1, 1, 0.5}
	}

	e.SpriteBatch.Begin()
	for _, tileIndex := range e.affectedTileIndexes {
		e.SpriteBatch.Draw(e.texture, float32(tileIndex.Col)*tileSize, float32(tileIndex.Row)*tileSize, tileSize, tileSize, rendering.WithColor(color))
	}
	if len(e.affectedTileIndexes) == 0 {
		e.SpriteBatch.Draw(e.texture, float32(e.x)*tileSize, float32(e.y)*tileSize, tileSize, tileSize, rendering.WithColor(color))
	}
	e.SpriteBatch.End()

	text := ""
	switch e.selected {
	case FLOOR_TOOL:
		text = "Floor"
	case HWALL_TOOL:
		text = "Horizontal wall"
	case VWALL_TOOL:
		text = "Vertical wall"
	case HDOOR_TOOL:
		text = "Horizontal door"
	case VDOOR_TOOL:
		text = "Vertical door"
	case FIXTURE_TOOL:
		text = "Fixture"
	}

	font.Printf(10, 20, text+" tool selected", rendering.WithTextColor(rendering.White), rendering.WithTextScale(0.25))
}

func (e *Editor) IsOverlay() bool {
	return false
}

var font *rendering.Font

func initFonts() {
	var err error
	if font, err = rendering.LoadFont("Roboto-Light"); err != nil {
		panic(err)
	}
}
