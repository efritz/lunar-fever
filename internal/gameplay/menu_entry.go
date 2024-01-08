package gameplay

import (
	"strings"

	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
)

type MenuEntrySelectionDelegate interface {
	OnSelect()
}

type MenuEntry struct {
	*engine.Context
	delegate MenuEntrySelectionDelegate
	text     string
	texture  rendering.Texture
}

func NewMenuEntry(engineCtx *engine.Context, text string, delegate MenuEntrySelectionDelegate) *MenuEntry {
	return &MenuEntry{
		Context:  engineCtx,
		text:     text,
		delegate: delegate,
	}
}

func (e *MenuEntry) OnSelect() {
	e.delegate.OnSelect()
}

func (e *MenuEntry) Init() {
	e.texture = e.TextureLoader.Load("base").Region(7*32, 1*32, 32, 32)
}

func (e *MenuEntry) Exit() {}

func (e *MenuEntry) Update(elapsedMs int64, selected bool) {}

func (e *MenuEntry) Render(elapsedMs int64, index int, selected bool) {
	a := float32(0.5)
	if selected {
		a = 1
	}

	// e.SpriteBatch.Begin()
	// e.SpriteBatch.Draw(
	// 	e.texture,
	// 	float32(e.xpos), float32(e.ypos)+8,
	// 	256, 2,
	// 	rendering.WithColor(rendering.Color{1, 1, 1, a}),
	// )
	// e.SpriteBatch.End()

	font.Printf(
		float32(128),
		float32(128+32*index),
		strings.ToUpper(e.text),
		rendering.WithTextScale(0.4),
		rendering.WithTextColor(rendering.Color{1, 1, 1, a}),
	)
}
