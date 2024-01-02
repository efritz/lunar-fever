package gameplay

import (
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
	"github.com/efritz/lunar-fever/internal/engine/view"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type PauseMenu struct {
	*engine.Context
	beginExiting func()
	texture      rendering.Texture
}

// setTransitionOnTime(250);

func NewPauseMenu(engineCtx *engine.Context) view.View {
	m := &PauseMenu{Context: engineCtx}
	v := view.NewTransitionView(m, engineCtx.ViewManager)
	m.beginExiting = v.BeginExiting
	return v
}

func (m *PauseMenu) Init() {
	m.texture = m.TextureLoader.Load("base").Region(6*32, 0*32, 32, 32)
}

func (m *PauseMenu) Exit() {}

func (m *PauseMenu) Update(elapsedMs int64, hasFocus bool) {
	if m.Keyboard.IsKeyNewlyDown(glfw.KeySpace) {
		m.beginExiting()
	}
}

func (m *PauseMenu) Render(elapsedMs int64) {
	color := rendering.White // color.a = (float) getTransitionPosition();

	m.SpriteBatch.Begin()
	m.SpriteBatch.Draw(
		m.texture,               // texture
		0,                       // x
		0,                       // y
		rendering.DisplayWidth,  // w
		rendering.DisplayHeight, // h
		rendering.WithColor(color),
	)
	m.SpriteBatch.End()
}

func (m *PauseMenu) IsOverlay() bool {
	return false
}
