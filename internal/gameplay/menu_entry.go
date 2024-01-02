package gameplay

import (
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
)

type MenuEntrySelectionDelegate interface {
	OnSelect()
}

type MenuEntry struct {
	*engine.Context
	delegate MenuEntrySelectionDelegate
	texture  rendering.Texture
	xpos     int64
	ypos     int64
}

func NewMenuEntry(engineCtx *engine.Context, delegate MenuEntrySelectionDelegate) *MenuEntry {
	return &MenuEntry{
		Context:  engineCtx,
		delegate: delegate,
	}
}

func (e *MenuEntry) SetPosition(xpos, ypos int64) {
	e.xpos = xpos
	e.ypos = ypos
}

func (e *MenuEntry) OnSelect() {
	e.delegate.OnSelect()
}

func (e *MenuEntry) Init() {
	e.texture = e.TextureLoader.Load("base").Region(7*32, 1*32, 32, 32)
}

func (e *MenuEntry) Exit() {}

func (e *MenuEntry) Update(elapsedMs int64, selected bool) {}

func (e *MenuEntry) Render(elapsedMs int64, selected bool) {
	a := float32(0.5)
	if selected {
		a = 1
	}

	e.SpriteBatch.Begin()
	e.SpriteBatch.Draw(e.texture, float32(e.xpos), float32(e.ypos), 128, 64,
		rendering.WithColor(rendering.Color{1, 0, 1, a}),
		rendering.WithRotation(0.25),
		rendering.WithOrigin(0, 25),
		rendering.WithScale(1, 1.5),
	)
	e.SpriteBatch.End()
}
