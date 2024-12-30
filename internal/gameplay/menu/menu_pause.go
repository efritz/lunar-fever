package menu

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
	delegate := &PauseMenu{Context: engineCtx}
	menu := NewMenu(engineCtx, delegate)
	v := view.NewTransitionView(menu, engineCtx.ViewManager)
	delegate.beginExiting = v.BeginExiting

	menu.AddEntry("Resume", &resumeMenuEntry{beginExiting: v.BeginExiting})
	menu.AddEntry("Exit", &exitMenuEntry{exit: engineCtx.Game.Stop})

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

	// font.Printf(
	// 	float32(128),
	// 	float32(128),
	// 	strings.ToUpper("PAUSED"),
	// 	rendering.WithTextScale(0.4),
	// 	rendering.WithTextColor(rendering.Color{1, 1, 1, 1}),
	// )
}

func (m *PauseMenu) IsOverlay() bool {
	return false
}

type resumeMenuEntry struct {
	*engine.Context
	beginExiting func()
}

func (e *resumeMenuEntry) OnSelect() {
	e.beginExiting()
}
