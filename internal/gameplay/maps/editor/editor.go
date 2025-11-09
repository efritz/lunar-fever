package editor

import (
	stdmath "math"

	"github.com/efritz/lunar-fever/internal/common/math"
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

	offsetRow int
	offsetCol int
}

func NewEditor(engineCtx *engine.Context) view.View {
	return &Editor{
		Context: engineCtx,
	}
}

func (e *Editor) Init() {
	tm, err := loader.ReadTileMap()
	if err != nil {
		tm = maps.NewTileMap(50, 50, 64)
	}
	e.tileMap = tm

	e.texture = e.TextureLoader.Load("base").Region(7*32, 1*32, 32, 32)
	e.baseRenderer = maps.NewBaseRenderer(e.SpriteBatch, e.TextureLoader, e.tileMap, true)
	e.executor = NewMapCommandExecutor(e.tileMap)
	e.selected = FLOOR_TOOL
	initFonts()
}

func (e *Editor) Exit() {}

func (e *Editor) ensureMapAccommodates(row, col, padding int) {
	expandLeft := 0
	if minCol := col - padding; minCol < 0 {
		expandLeft = -minCol
	}

	expandTop := 0
	if minRow := row - padding; minRow < 0 {
		expandTop = -minRow
	}

	expandRight := 0
	if maxCol := col + padding; maxCol >= e.tileMap.Width() {
		expandRight = maxCol - e.tileMap.Width() + 1
	}

	expandBottom := 0
	if maxRow := row + padding; maxRow >= e.tileMap.Height() {
		expandBottom = maxRow - e.tileMap.Height() + 1
	}

	if expandLeft == 0 && expandTop == 0 && expandRight == 0 && expandBottom == 0 {
		return
	}

	oldMap := e.tileMap
	newMap := maps.NewTileMap(
		oldMap.Width()+expandLeft+expandRight,
		oldMap.Height()+expandTop+expandBottom,
		oldMap.GridSize(),
	)

	for col := 0; col < oldMap.Width(); col++ {
		for row := 0; row < oldMap.Height(); row++ {
			newMap.SetBits(row+expandTop, col+expandLeft, oldMap.GetBits(row, col))
		}
	}

	e.tileMap = newMap
	e.offsetRow += expandTop
	e.offsetCol += expandLeft

	e.baseRenderer = maps.NewBaseRenderer(e.SpriteBatch, e.TextureLoader, e.tileMap, true)
	e.executor = NewMapCommandExecutor(e.tileMap)
}

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
	// Camera controls

	cameraXDir := int64(0)
	cameraYDir := int64(0)
	if e.Keyboard.IsKeyDown(glfw.KeyUp) {
		cameraYDir++
	}
	if e.Keyboard.IsKeyDown(glfw.KeyDown) {
		cameraYDir--
	}
	if e.Keyboard.IsKeyDown(glfw.KeyLeft) {
		cameraXDir++
	}
	if e.Keyboard.IsKeyDown(glfw.KeyRight) {
		cameraXDir--
	}

	mod := float32(500)
	e.Camera.Translate(float32(cameraXDir*elapsedMs), float32(cameraYDir*elapsedMs))
	e.Camera.Zoom(float32(e.Mouse.ScrollDelta()) / mod)

	//
	// Update position

	size := 64
	mx := e.Camera.Unprojectx(float32(e.Mouse.X()))
	my := e.Camera.UnprojectY(float32(e.Mouse.Y()))
	x := int(stdmath.Floor(float64(mx) / float64(size)))
	y := int(stdmath.Floor(float64(my) / float64(size)))

	oldX := e.x
	oldY := e.y
	e.x = x
	e.y = y

	row := y + e.offsetRow
	col := x + e.offsetCol
	e.ensureMapAccommodates(row, col, 2) // padding

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
		if err := loader.Write(e.tileMap.Trim()); err != nil {
			panic(err.Error())
		}
	}
}

func (e *Editor) Render(elapsedMs int64) {
	tileSize := float32(64)

	gridSize := float32(e.tileMap.GridSize())
	offsetX := -float32(e.offsetCol) * gridSize
	offsetY := -float32(e.offsetRow) * gridSize

	offsetMatrix := math.IdentityMatrix.Translate(offsetX, offsetY)
	combinedMatrix := e.Camera.ViewMatrix().Multiply(offsetMatrix)
	e.SpriteBatch.SetViewMatrix(combinedMatrix)

	x1, y1, x2, y2 := e.Camera.Bounds()
	e.baseRenderer.Render(x1-offsetX, y1-offsetY, x2-offsetX, y2-offsetY, nil, nil, false)

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
		row := e.y + e.offsetRow
		col := e.x + e.offsetCol
		e.SpriteBatch.Draw(e.texture, float32(col)*tileSize, float32(row)*tileSize, tileSize, tileSize, rendering.WithColor(color))
	}
	e.SpriteBatch.End()
	e.SpriteBatch.SetViewMatrix(math.IdentityMatrix)

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
