package maps

import (
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/ecs/system"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type baseRenderSystem struct {
	*engine.Context
	tileMap      *TileMap
	baseRenderer *BaseRenderer
	base         *Base
	debugging    bool
}

func NewBaseRenderSystem(engineCtx *engine.Context, tileMap *TileMap, base *Base) system.System {
	return &baseRenderSystem{
		Context: engineCtx,
		tileMap: tileMap,
		base:    base,
	}
}

func (s *baseRenderSystem) Init() {
	s.baseRenderer = NewBaseRenderer(s.SpriteBatch, s.TextureLoader, s.tileMap, false)
}

func (s *baseRenderSystem) Exit() {}

func (s *baseRenderSystem) Process(elapsedMs int64) {
	if s.Keyboard.IsKeyNewlyDown(glfw.KeyL) {
		s.debugging = !s.debugging
	}

	x1, y1, x2, y2 := s.Camera.Bounds()
	s.baseRenderer.Render(x1, y1, x2, y2, s.base.Rooms, s.base.NavigationGraph, s.debugging)
}
