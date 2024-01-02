package gameplay

import (
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
	"github.com/efritz/lunar-fever/internal/engine/view"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type ObjectiveMenu struct {
	*engine.Context
	beginExiting func()
	texture      rendering.Texture
}

func NewObjectiveMenu(engineCtx *engine.Context) view.View {
	m := &ObjectiveMenu{Context: engineCtx}
	v := view.NewTransitionView(m, engineCtx.ViewManager)
	m.beginExiting = v.BeginExiting
	return v
}

func (m *ObjectiveMenu) Init() {
	m.texture = m.TextureLoader.Load("base").Region(7*32, 0*32, 32, 32)
}

func (m *ObjectiveMenu) Exit() {}

func (m *ObjectiveMenu) Update(elapsedMs int64, hasFocus bool) {
	if !m.Keyboard.IsKeyDown(glfw.KeyTab) {
		m.beginExiting()
		return
	}
}

func (m *ObjectiveMenu) Render(elapsedMs int64) {
	m.SpriteBatch.Begin()
	m.SpriteBatch.Draw(
		m.texture,                 // texture
		rendering.DisplayWidth/4,  // x
		rendering.DisplayHeight/4, // y
		rendering.DisplayWidth/2,  // w
		rendering.DisplayHeight/2, // h
		rendering.WithColor(rendering.Color{1, 1, 1, .75}))
	m.SpriteBatch.End()
}

func (m *ObjectiveMenu) IsOverlay() bool {
	return true
}
