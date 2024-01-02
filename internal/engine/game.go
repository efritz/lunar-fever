package engine

import (
	"github.com/efritz/lunar-fever/internal/engine/game"
	"github.com/efritz/lunar-fever/internal/engine/view"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type delegate struct {
	*Context
	initialView func(engineCtx *Context) view.View
}

type DelegateOption func(*delegate)

func WithInitialView(f func(engineCtx *Context) view.View) DelegateOption {
	return func(d *delegate) { d.initialView = f }
}

func InitGame(opts ...DelegateOption) (*game.Game, error) {
	delegate := &delegate{}
	for _, opt := range opts {
		opt(delegate)
	}

	game := game.New(delegate)
	engineCtx, err := InitContext(game)
	if err != nil {
		return nil, err
	}

	delegate.Context = engineCtx
	return game, nil
}

func (d *delegate) Init() {
	d.ViewManager.Add(d.initialView(d.Context))
	d.ViewManager.Init()
}

func (d *delegate) Exit() {
	d.ViewManager.Exit()
	glfw.Terminate()
}

func (d *delegate) Active() bool {
	return true
}

func (d *delegate) Update(elapsedMs int64) {
	if d.Window.ShouldClose() {
		d.Game.Stop()
		return
	}

	d.Keyboard.Update()
	d.Mouse.Update()
	d.ViewManager.Update(elapsedMs)
}

func (d *delegate) Render(elapsedMs int64) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.LoadIdentity()

	d.ViewManager.Render(elapsedMs)
	d.Window.SwapBuffers()

	glfw.PollEvents()
	glfw.GetCurrentContext()
}
