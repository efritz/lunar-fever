package maps

import (
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/ecs/system"
)

type baseRenderSystem struct {
	*engine.Context
	tileMap      *TileMap
	baseRenderer *BaseRenderer
}

func NewBaseRenderSystem(engineCtx *engine.Context, tileMap *TileMap) system.System {
	return &baseRenderSystem{
		Context: engineCtx,
		tileMap: tileMap,
	}
}

func (s *baseRenderSystem) Init() {
	s.baseRenderer = NewBaseRenderer(s.SpriteBatch, s.TextureLoader, s.tileMap, false)
}

func (s *baseRenderSystem) Exit() {}

func (s *baseRenderSystem) Process(elapsedMs int64) {
	s.baseRenderer.Render(s.Camera.Bounds())
}
