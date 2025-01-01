package maps

import (
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/ecs/system"
)

type baseRenderSystem struct {
	*engine.Context
	tileMap         *TileMap
	baseRenderer    *BaseRenderer
	rooms           []Room
	doors           []Door
	navigationGraph *NavigationGraph
}

func NewBaseRenderSystem(engineCtx *engine.Context, tileMap *TileMap) system.System {
	return &baseRenderSystem{
		Context: engineCtx,
		tileMap: tileMap,
	}
}

func (s *baseRenderSystem) Init() {
	s.baseRenderer = NewBaseRenderer(s.SpriteBatch, s.TextureLoader, s.tileMap, false)
	s.rooms, s.doors = PartitionRooms(s.tileMap)
	s.navigationGraph = ConstructNavigationGraph(s.rooms, s.doors)
}

func (s *baseRenderSystem) Exit() {}

func (s *baseRenderSystem) Process(elapsedMs int64) {
	x1, y1, x2, y2 := s.Camera.Bounds()
	s.baseRenderer.Render(x1, y1, x2, y2, s.rooms, s.doors, s.navigationGraph)
}
