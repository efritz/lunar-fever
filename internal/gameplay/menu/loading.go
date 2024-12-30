package menu

import (
	"math"
	"sync"
	"sync/atomic"

	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
	"github.com/efritz/lunar-fever/internal/engine/view"
)

type Loading struct {
	*engine.Context
	beginExiting   func()
	target         view.View
	once           sync.Once
	done           atomic.Bool
	loadFn         func()
	texture        rendering.Texture
	totalElapsedMs int64
}

func Load(engineCtx *engine.Context, target view.View, loadFn func()) {
	inner := &Loading{
		Context: engineCtx,
		target:  target,
		loadFn:  loadFn,
	}

	loading := view.NewTransitionView(inner, engineCtx.ViewManager)
	inner.beginExiting = loading.BeginExiting

	engineCtx.ViewManager.Clear()
	engineCtx.ViewManager.Add(loading)
}

func (l *Loading) Init() {
	l.texture = l.TextureLoader.Load("base").Region(7*32, 1*32, 32, 32)
}

func (l *Loading) Exit() {
	// TODO - cancel loading func
}

func (l *Loading) Update(elapsedMs int64, hasFocus bool) {
	if l.ViewManager.NumViews() != 1 {
		return
	}

	l.once.Do(func() {
		go func() {
			l.loadFn()
			l.done.Store(true)
		}()
	})

	if l.done.Load() {
		l.beginExiting()
		l.ViewManager.Add(l.target)
	}
}

func (l *Loading) Render(elapsedMs int64) {
	l.totalElapsedMs += elapsedMs
	rot := float32(l.totalElapsedMs/2%360) * math.Pi / 180

	l.SpriteBatch.Begin()
	l.SpriteBatch.Draw(
		l.texture,                  // texture
		rendering.DisplayWidth-96,  // x
		rendering.DisplayHeight-96, // y
		64,                         // w
		64,                         // h
		rendering.WithRotation(rot),
		rendering.WithOrigin(32, 32),
	)
	l.SpriteBatch.End()
}

func (l *Loading) IsOverlay() bool {
	return true
}
